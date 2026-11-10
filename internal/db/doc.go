// Package db encapsulates PostgreSQL schema migrations and schema initialization.
//
// # Responsibilities:
//   - migrations/: Embedded SQL migration files ordered sequentially (001_initial.sql, etc.).
//   - migrations.go: Robust transactional migration runner applying pending schemas on startup.
//
// All database access from domain modules targets the PostgreSQL connection pool managed here.
package db
