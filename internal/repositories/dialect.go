package repositories

import (
	"database/sql"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	DriverSQLite   = "sqlite"
	DriverPostgres = "pgx"
)

func Rebind(driver, query string) string {
	if driver == DriverPostgres {
		return sqlx.Rebind(sqlx.DOLLAR, query)
	}
	return query
}

func DriverName(db *sql.DB) string {
	if db == nil {
		return ""
	}
	if _, ok := db.Driver().(*stdlib.Driver); ok {
		return DriverPostgres
	}
	return DriverSQLite
}
