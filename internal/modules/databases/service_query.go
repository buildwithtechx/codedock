package databases

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

func (s *Service) QueryDatabase(ctx context.Context, id, query, _ string) (*DatabaseQueryResponse, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, ErrDatabaseNotFound
	}

	switch strings.ToLower(string(db.Engine)) {
	case "postgresql", "postgres":
		return executePostgresQuery(db, query)
	case "mysql", "mariadb":
		return executeMySQLQuery(db, query)
	case "redis":
		return executeRedisCommand(ctx, db, query)
	default:
		return nil, fmt.Errorf("interactive console not supported for engine: %s", db.Engine)
	}
}

func executePostgresQuery(db *Database, query string) (*DatabaseQueryResponse, error) {
	dsn := fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.Port, db.Username, db.Password, db.DatabaseName)
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return runSQLQuery(conn, query)
}

func executeMySQLQuery(db *Database, query string) (*DatabaseQueryResponse, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:%d)/%s", db.Username, db.Password, db.Port, db.DatabaseName)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return runSQLQuery(conn, query)
}

func runSQLQuery(conn *sql.DB, query string) (*DatabaseQueryResponse, error) {
	start := time.Now()
	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]any
	for rows.Next() {
		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]any)
		for i, colName := range cols {
			val := columnPointers[i].(*any)
			if b, ok := (*val).([]byte); ok {
				m[colName] = string(b)
			} else {
				m[colName] = *val
			}
		}
		results = append(results, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	elapsed := time.Since(start).Milliseconds()
	if results == nil {
		results = []map[string]any{}
	}
	if cols == nil {
		cols = []string{}
	}

	return &DatabaseQueryResponse{
		Columns:         cols,
		Rows:            results,
		RowCount:        len(results),
		ExecutionTimeMs: elapsed,
	}, nil
}

func executeRedisCommand(ctx context.Context, db *Database, command string) (*DatabaseQueryResponse, error) {
	addr := fmt.Sprintf("localhost:%d", db.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: db.Password,
		DB:       0,
	})
	defer rdb.Close()

	parts := strings.Fields(command)
	if len(parts) == 0 {
		return &DatabaseQueryResponse{}, nil
	}

	args := make([]any, len(parts))
	for i, v := range parts {
		args[i] = v
	}

	start := time.Now()
	res, err := rdb.Do(ctx, args...).Result()
	elapsed := time.Since(start).Milliseconds()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	return &DatabaseQueryResponse{
		Result:          fmt.Sprintf("%v", res),
		ExecutionTimeMs: elapsed,
	}, nil
}
