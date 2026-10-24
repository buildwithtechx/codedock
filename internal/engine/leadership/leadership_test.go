package leadership

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"codedock/internal/testdb"
)

func testConn(t *testing.T, db *sql.DB) *sql.Conn {
	t.Helper()
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Fatalf("dedicated connection: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

func TestAdvisoryLockMutualExclusion(t *testing.T) {
	db := testdb.Open(t)
	first := NewElector(db, "exclusion-probe")
	second := NewElector(db, "exclusion-probe")
	ctx := context.Background()
	held := testConn(t, db)
	if !first.tryAcquire(ctx, held) {
		t.Fatal("first elector did not acquire a free lock")
	}
	defer first.release(ctx, held)
	contender := testConn(t, db)
	if second.tryAcquire(ctx, contender) {
		t.Fatal("second elector acquired a held lock")
	}
	first.release(ctx, held)
	if !second.tryAcquire(ctx, contender) {
		t.Fatal("second elector did not acquire a released lock")
	}
	second.release(ctx, contender)
}

func TestAdvisoryLockDistinctKeysCoexist(t *testing.T) {
	db := testdb.Open(t)
	ctx := context.Background()
	left := NewElector(db, "coexist-left")
	right := NewElector(db, "coexist-right")
	leftConn := testConn(t, db)
	rightConn := testConn(t, db)
	if !left.tryAcquire(ctx, leftConn) {
		t.Fatal("left key did not acquire")
	}
	defer left.release(ctx, leftConn)
	if !right.tryAcquire(ctx, rightConn) {
		t.Fatal("right key did not acquire alongside left")
	}
	defer right.release(ctx, rightConn)
}

func TestElectorRunServesUntilCancelled(t *testing.T) {
	db := testdb.Open(t)
	elector := NewElector(db, "run-probe")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	serving := make(chan struct{})
	stopped := make(chan struct{})
	go elector.Run(ctx, func(serveCtx context.Context) {
		close(serving)
		<-serveCtx.Done()
		close(stopped)
	})
	select {
	case <-serving:
	case <-time.After(10 * time.Second):
		t.Fatal("elector did not start serving")
	}
	cancel()
	select {
	case <-stopped:
	case <-time.After(10 * time.Second):
		t.Fatal("elector did not stop serving after cancel")
	}
}

func TestElectorFailoverToContender(t *testing.T) {
	db := testdb.Open(t)
	first := NewElector(db, "failover-probe")
	first.retry = 100 * time.Millisecond
	second := NewElector(db, "failover-probe")
	second.retry = 100 * time.Millisecond
	firstCtx, cancelFirst := context.WithCancel(context.Background())
	defer cancelFirst()
	secondCtx, cancelSecond := context.WithCancel(context.Background())
	defer cancelSecond()
	firstServing := make(chan struct{})
	secondServing := make(chan struct{})
	go first.Run(firstCtx, func(serveCtx context.Context) {
		close(firstServing)
		<-serveCtx.Done()
	})
	select {
	case <-firstServing:
	case <-time.After(10 * time.Second):
		t.Fatal("first elector did not start serving")
	}
	go second.Run(secondCtx, func(serveCtx context.Context) {
		close(secondServing)
		<-serveCtx.Done()
	})
	select {
	case <-secondServing:
		t.Fatal("contender served while the lock was held")
	case <-time.After(500 * time.Millisecond):
	}
	cancelFirst()
	select {
	case <-secondServing:
	case <-time.After(10 * time.Second):
		t.Fatal("contender did not take over after release")
	}
}
