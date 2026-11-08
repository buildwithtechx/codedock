package databases

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (s *Service) GetSchemas(ctx context.Context, id string) ([]TableSchema, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, ErrDatabaseNotFound
	}

	switch strings.ToLower(string(db.Engine)) {
	case "postgresql", "postgres":
		return getPostgresSchemas(db)
	case "mysql", "mariadb":
		return getMySQLSchemas(db)
	default:
		return nil, fmt.Errorf("schema introspection not supported for engine: %s", db.Engine)
	}
}

func getPostgresSchemas(db *Database) ([]TableSchema, error) {
	dsn := fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.Port, db.Username, db.Password, db.DatabaseName)
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, fmt.Errorf("failed to scan table name row: %w", err)
		}
		tableNames = append(tableNames, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("table rows iteration error: %w", err)
	}

	colRows, err := conn.Query("SELECT table_name, column_name, data_type, is_nullable FROM information_schema.columns WHERE table_schema='public'")
	if err != nil {
		return nil, err
	}
	defer colRows.Close()

	schemaMap := make(map[string][]ColumnSchema)
	for colRows.Next() {
		var tName, cName, cType, cNullable string
		if err := colRows.Scan(&tName, &cName, &cType, &cNullable); err != nil {
			return nil, fmt.Errorf("failed to scan column row: %w", err)
		}
		schemaMap[tName] = append(schemaMap[tName], ColumnSchema{
			Name:       cName,
			Type:       cType,
			IsNullable: cNullable == "YES",
			IsPrimary:  false,
		})
	}
	if err := colRows.Err(); err != nil {
		return nil, fmt.Errorf("column rows iteration error: %w", err)
	}

	var schemas []TableSchema
	for _, t := range tableNames {
		cols := schemaMap[t]
		if cols == nil {
			cols = []ColumnSchema{}
		}
		schemas = append(schemas, TableSchema{Name: t, Columns: cols})
	}
	if schemas == nil {
		schemas = []TableSchema{}
	}
	return schemas, nil
}

func getMySQLSchemas(db *Database) ([]TableSchema, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:%d)/%s", db.Username, db.Password, db.Port, db.DatabaseName)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, fmt.Errorf("failed to scan table name row: %w", err)
		}
		tableNames = append(tableNames, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("table rows iteration error: %w", err)
	}

	colRows, err := conn.Query("SELECT TABLE_NAME, COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_KEY FROM information_schema.columns WHERE table_schema=?", db.DatabaseName)
	if err != nil {
		return nil, err
	}
	defer colRows.Close()

	schemaMap := make(map[string][]ColumnSchema)
	for colRows.Next() {
		var tName, cField, cType, cNull, cKey string
		if err := colRows.Scan(&tName, &cField, &cType, &cNull, &cKey); err != nil {
			return nil, fmt.Errorf("failed to scan column row: %w", err)
		}
		schemaMap[tName] = append(schemaMap[tName], ColumnSchema{
			Name:       cField,
			Type:       cType,
			IsNullable: cNull == "YES",
			IsPrimary:  cKey == "PRI",
		})
	}
	if err := colRows.Err(); err != nil {
		return nil, fmt.Errorf("column rows iteration error: %w", err)
	}

	var schemas []TableSchema
	for _, t := range tableNames {
		cols := schemaMap[t]
		if cols == nil {
			cols = []ColumnSchema{}
		}
		schemas = append(schemas, TableSchema{Name: t, Columns: cols})
	}
	if schemas == nil {
		schemas = []TableSchema{}
	}
	return schemas, nil
}

func (s *Service) GetTableData(ctx context.Context, id, table string, limit, offset int) (*DatabaseQueryResponse, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, ErrDatabaseNotFound
	}

	switch strings.ToLower(string(db.Engine)) {
	case "postgresql", "postgres":
		query := fmt.Sprintf("SELECT * FROM \"%s\" LIMIT %d OFFSET %d", table, limit, offset)
		return s.QueryDatabase(ctx, id, query, db.ProjectID)
	case "mysql", "mariadb":
		query := fmt.Sprintf("SELECT * FROM `%s` LIMIT %d OFFSET %d", table, limit, offset)
		return s.QueryDatabase(ctx, id, query, db.ProjectID)
	default:
		return nil, fmt.Errorf("data browsing not supported for engine: %s", db.Engine)
	}
}
