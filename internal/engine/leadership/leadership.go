package leadership

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"
	"time"
)

const lockNamespace = 727425696

type Elector struct {
	db       *sql.DB
	key      string
	retry    time.Duration
	watchdog time.Duration
}

func NewElector(db *sql.DB, key string) *Elector {
	return &Elector{db: db, key: key, retry: 5 * time.Second, watchdog: 15 * time.Second}
}

func (e *Elector) Run(ctx context.Context, serves ...func(context.Context)) {
	for {
		if ctx.Err() != nil {
			return
		}
		conn, err := e.db.Conn(ctx)
		if err != nil {
			slog.Warn("leadership connection failed", "key", e.key, "error", err)
			e.sleep(ctx)
			continue
		}
		if !e.tryAcquire(ctx, conn) {
			_ = conn.Close()
			e.sleep(ctx)
			continue
		}
		slog.Info("leadership acquired", "key", e.key)
		e.serve(ctx, conn, serves)
		slog.Info("leadership released", "key", e.key)
	}
}

func (e *Elector) serve(ctx context.Context, conn *sql.Conn, serves []func(context.Context)) {
	defer func() { _ = conn.Close() }()
	serveCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	for _, serve := range serves {
		wg.Add(1)
		go func() {
			defer wg.Done()
			serve(serveCtx)
		}()
	}
	finished := make(chan struct{})
	go func() {
		wg.Wait()
		close(finished)
	}()
	ticker := time.NewTicker(e.watchdog)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			cancel()
			<-finished
			e.release(context.Background(), conn)
			return
		case <-finished:
			e.release(context.Background(), conn)
			return
		case <-ticker.C:
			if err := conn.PingContext(ctx); err != nil {
				slog.Warn("leadership connection lost", "key", e.key, "error", err)
				cancel()
				<-finished
				return
			}
		}
	}
}

func (e *Elector) tryAcquire(ctx context.Context, conn *sql.Conn) bool {
	var held bool
	if err := conn.QueryRowContext(ctx, `SELECT pg_try_advisory_lock($1, hashtext($2))`, lockNamespace, e.key).Scan(&held); err != nil {
		return false
	}
	return held
}

func (e *Elector) release(ctx context.Context, conn *sql.Conn) {
	_, _ = conn.ExecContext(ctx, `SELECT pg_advisory_unlock($1, hashtext($2))`, lockNamespace, e.key)
}

func (e *Elector) sleep(ctx context.Context) {
	timer := time.NewTimer(e.retry)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}
