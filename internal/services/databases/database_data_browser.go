package databases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"codedock.run/codedock/internal/models"
)

func (s *DatabaseService) GetSchemas(ctx context.Context, id string) ([]models.TableSchema, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, errors.New("database not found")
	}

	switch db.Engine {
	case "postgresql", "postgres":
		return getPostgresSchemas(db)
	case "mysql", "mariadb":
		return getMySQLSchemas(db)
	default:
		return nil, fmt.Errorf("schema introspection not supported for engine: %s", db.Engine)
	}
}

func getPostgresSchemas(db *models.Database) ([]models.TableSchema, error) {
	host := "localhost"
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, db.Port, db.Username, db.Password, db.DatabaseName)
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

	schemaMap := make(map[string][]models.ColumnSchema)
	for colRows.Next() {
		var tName, cName, cType, cNullable string
		if err := colRows.Scan(&tName, &cName, &cType, &cNullable); err != nil {
			return nil, fmt.Errorf("failed to scan column row: %w", err)
		}
		schemaMap[tName] = append(schemaMap[tName], models.ColumnSchema{
			Name:       cName,
			Type:       cType,
			IsNullable: cNullable == "YES",
			IsPrimary:  false,
		})
	}

	if err := colRows.Err(); err != nil {
		return nil, fmt.Errorf("column rows iteration error: %w", err)
	}

	var schemas []models.TableSchema
	for _, t := range tableNames {
		cols := schemaMap[t]
		if cols == nil {
			cols = []models.ColumnSchema{}
		}
		schemas = append(schemas, models.TableSchema{Name: t, Columns: cols})
	}
	if schemas == nil {
		schemas = []models.TableSchema{}
	}
	return schemas, nil
}

func getMySQLSchemas(db *models.Database) ([]models.TableSchema, error) {
	host := "localhost"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db.Username, db.Password, host, db.Port, db.DatabaseName)
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

	schemaMap := make(map[string][]models.ColumnSchema)
	for colRows.Next() {
		var tName, cField, cType, cNull, cKey string
		if err := colRows.Scan(&tName, &cField, &cType, &cNull, &cKey); err != nil {
			return nil, fmt.Errorf("failed to scan column row: %w", err)
		}
		schemaMap[tName] = append(schemaMap[tName], models.ColumnSchema{
			Name:       cField,
			Type:       cType,
			IsNullable: cNull == "YES",
			IsPrimary:  cKey == "PRI",
		})
	}

	if err := colRows.Err(); err != nil {
		return nil, fmt.Errorf("column rows iteration error: %w", err)
	}

	var schemas []models.TableSchema
	for _, t := range tableNames {
		cols := schemaMap[t]
		if cols == nil {
			cols = []models.ColumnSchema{}
		}
		schemas = append(schemas, models.TableSchema{Name: t, Columns: cols})
	}
	if schemas == nil {
		schemas = []models.TableSchema{}
	}
	return schemas, nil
}

func (s *DatabaseService) GetTableData(ctx context.Context, id, table string, limit, offset int) (*models.DatabaseQueryResponse, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, errors.New("database not found")
	}

	switch db.Engine {
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

func escapeSQLValue(v any) string {
	if v == nil {
		return "NULL"
	}
	str := fmt.Sprintf("%v", v)
	str = strings.ReplaceAll(str, "'", "''")
	return fmt.Sprintf("'%s'", str)
}

func (s *DatabaseService) InsertTableRow(ctx context.Context, id, table string, data map[string]any) (*models.DatabaseQueryResponse, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, errors.New("database not found")
	}

	query, err := buildInsertSQL(db.Engine, table, data)
	if err != nil {
		return nil, err
	}
	return s.QueryDatabase(ctx, id, query, db.ProjectID)
}

type UpdateTableRowOpts struct {
	ID    string
	Table string
	Keys  map[string]any
	Data  map[string]any
}

func (s *DatabaseService) UpdateTableRow(ctx context.Context, opts UpdateTableRowOpts) (*models.DatabaseQueryResponse, error) {
	db, err := s.repo.GetByID(ctx, opts.ID)
	if err != nil || db == nil {
		return nil, errors.New("database not found")
	}

	query, err := buildUpdateSQL(db.Engine, opts.Table, opts.Keys, opts.Data)
	if err != nil {
		return nil, err
	}
	return s.QueryDatabase(ctx, opts.ID, query, db.ProjectID)
}

func (s *DatabaseService) DeleteTableRow(ctx context.Context, id, table string, keys map[string]any) (*models.DatabaseQueryResponse, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, errors.New("database not found")
	}

	query, err := buildDeleteSQL(db.Engine, table, keys)
	if err != nil {
		return nil, err
	}
	return s.QueryDatabase(ctx, id, query, db.ProjectID)
}

func buildInsertSQL(engine, table string, data map[string]any) (string, error) {
	var cols []string
	var vals []string
	for k, v := range data {
		if engine == "postgresql" || engine == "postgres" {
			cols = append(cols, fmt.Sprintf("\"%s\"", k))
		} else {
			cols = append(cols, fmt.Sprintf("`%s`", k))
		}
		vals = append(vals, escapeSQLValue(v))
	}

	switch engine {
	case "postgresql", "postgres":
		return fmt.Sprintf("INSERT INTO \"%s\" (%s) VALUES (%s)", table, strings.Join(cols, ", "), strings.Join(vals, ", ")), nil
	case "mysql", "mariadb":
		return fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", table, strings.Join(cols, ", "), strings.Join(vals, ", ")), nil
	default:
		return "", fmt.Errorf("inserts not supported for engine: %s", engine)
	}
}

func buildUpdateSQL(engine, table string, keys, data map[string]any) (string, error) {
	if len(keys) == 0 {
		return "", errors.New("at least one primary key is required for updates")
	}
	var sets []string
	var wheres []string
	isPostgres := engine == "postgresql" || engine == "postgres"

	for k, v := range data {
		if isPostgres {
			sets = append(sets, fmt.Sprintf("\"%s\"=%s", k, escapeSQLValue(v)))
		} else {
			sets = append(sets, fmt.Sprintf("`%s`=%s", k, escapeSQLValue(v)))
		}
	}
	for k, v := range keys {
		if isPostgres {
			wheres = append(wheres, fmt.Sprintf("\"%s\"=%s", k, escapeSQLValue(v)))
		} else {
			wheres = append(wheres, fmt.Sprintf("`%s`=%s", k, escapeSQLValue(v)))
		}
	}

	switch engine {
	case "postgresql", "postgres":
		return fmt.Sprintf("UPDATE \"%s\" SET %s WHERE %s", table, strings.Join(sets, ", "), strings.Join(wheres, " AND ")), nil
	case "mysql", "mariadb":
		return fmt.Sprintf("UPDATE `%s` SET %s WHERE %s", table, strings.Join(sets, ", "), strings.Join(wheres, " AND ")), nil
	default:
		return "", fmt.Errorf("updates not supported for engine: %s", engine)
	}
}

func buildDeleteSQL(engine, table string, keys map[string]any) (string, error) {
	if len(keys) == 0 {
		return "", errors.New("at least one primary key is required for deletes")
	}
	var wheres []string
	isPostgres := engine == "postgresql" || engine == "postgres"

	for k, v := range keys {
		if isPostgres {
			wheres = append(wheres, fmt.Sprintf("\"%s\"=%s", k, escapeSQLValue(v)))
		} else {
			wheres = append(wheres, fmt.Sprintf("`%s`=%s", k, escapeSQLValue(v)))
		}
	}

	switch engine {
	case "postgresql", "postgres":
		return fmt.Sprintf("DELETE FROM \"%s\" WHERE %s", table, strings.Join(wheres, " AND ")), nil
	case "mysql", "mariadb":
		return fmt.Sprintf("DELETE FROM `%s` WHERE %s", table, strings.Join(wheres, " AND ")), nil
	default:
		return "", fmt.Errorf("deletes not supported for engine: %s", engine)
	}
}
