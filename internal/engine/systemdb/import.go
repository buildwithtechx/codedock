package systemdb

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

type pgColumn struct {
	dataType   string
	isNullable bool
}

func ImportFromSQLite(ctx context.Context, pgdb *sql.DB, sqlitePath string) error {
	var migrated int
	if err := pgdb.QueryRowContext(ctx, `SELECT COUNT(*) FROM schema_migrations`).Scan(&migrated); err != nil {
		return fmt.Errorf("postgres must be migrated before import: %w", err)
	}
	litedb, err := sql.Open("sqlite", sqlitePath)
	if err != nil {
		return fmt.Errorf("failed to open sqlite database: %w", err)
	}
	defer litedb.Close()
	tables, err := pgBaseTables(ctx, pgdb)
	if err != nil {
		return err
	}
	empty, err := targetEmpty(ctx, pgdb, tables)
	if err != nil {
		return err
	}
	if !empty {
		return fmt.Errorf("refusing to import into a non-empty database")
	}
	liteTables, err := sqliteTables(litedb)
	if err != nil {
		return err
	}
	ordered, err := orderByForeignKeys(ctx, pgdb, tables)
	if err != nil {
		return err
	}
	total := 0
	for _, table := range ordered {
		if !liteTables[table] {
			continue
		}
		copied, err := copyTable(ctx, pgdb, litedb, table)
		if err != nil {
			return err
		}
		total += copied
	}
	slog.Info("sqlite import complete", "tables", len(ordered), "rows", total)
	return nil
}

func TargetEmpty(ctx context.Context, pgdb *sql.DB) (bool, error) {
	tables, err := pgBaseTables(ctx, pgdb)
	if err != nil {
		return false, err
	}
	return targetEmpty(ctx, pgdb, tables)
}

func targetEmpty(ctx context.Context, pgdb *sql.DB, tables []string) (bool, error) {
	for _, table := range tables {
		var nonEmpty bool
		if err := pgdb.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM "`+table+`")`).Scan(&nonEmpty); err != nil {
			return false, fmt.Errorf("failed to check %s: %w", table, err)
		}
		if nonEmpty {
			return false, nil
		}
	}
	return true, nil
}

func pgBaseTables(ctx context.Context, pgdb *sql.DB) ([]string, error) {
	rows, err := pgdb.QueryContext(ctx, `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE' AND table_name <> 'schema_migrations' ORDER BY table_name`)
	if err != nil {
		return nil, fmt.Errorf("failed to list postgres tables: %w", err)
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, table)
	}
	return tables, rows.Err()
}

func sqliteTables(litedb *sql.DB) (map[string]bool, error) {
	rows, err := litedb.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'`)
	if err != nil {
		return nil, fmt.Errorf("failed to list sqlite tables: %w", err)
	}
	defer rows.Close()
	tables := map[string]bool{}
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("failed to scan sqlite table: %w", err)
		}
		tables[table] = true
	}
	return tables, rows.Err()
}

func orderByForeignKeys(ctx context.Context, pgdb *sql.DB, tables []string) ([]string, error) {
	rows, err := pgdb.QueryContext(ctx, `
		SELECT tc.table_name, ccu.table_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.constraint_column_usage ccu ON ccu.constraint_name = tc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_schema = 'public'`)
	if err != nil {
		return nil, fmt.Errorf("failed to list foreign keys: %w", err)
	}
	defer rows.Close()
	depends := map[string]map[string]bool{}
	for _, table := range tables {
		depends[table] = map[string]bool{}
	}
	for rows.Next() {
		var table, ref string
		if err := rows.Scan(&table, &ref); err != nil {
			return nil, fmt.Errorf("failed to scan foreign key: %w", err)
		}
		if table == ref {
			continue
		}
		if _, ok := depends[table]; ok {
			depends[table][ref] = true
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var ordered []string
	placed := map[string]bool{}
	for len(ordered) < len(tables) {
		progress := false
		for _, table := range tables {
			if placed[table] {
				continue
			}
			ready := true
			for dep := range depends[table] {
				if !placed[dep] {
					ready = false
					break
				}
			}
			if ready {
				ordered = append(ordered, table)
				placed[table] = true
				progress = true
			}
		}
		if !progress {
			var rest []string
			for _, table := range tables {
				if !placed[table] {
					rest = append(rest, table)
				}
			}
			sort.Strings(rest)
			ordered = append(ordered, rest...)
			break
		}
	}
	return ordered, nil
}

func copyTable(ctx context.Context, pgdb, litedb *sql.DB, table string) (int, error) {
	pgCols, err := pgColumns(ctx, pgdb, table)
	if err != nil {
		return 0, err
	}
	liteCols, err := sqliteColumns(litedb, table)
	if err != nil {
		return 0, err
	}
	var columns []string
	for _, col := range liteCols {
		if _, ok := pgCols[col]; ok {
			columns = append(columns, col)
		}
	}
	if len(columns) == 0 {
		return 0, nil
	}
	quoted := make([]string, len(columns))
	for i, col := range columns {
		quoted[i] = `"` + col + `"`
	}
	rows, err := litedb.Query(`SELECT ` + strings.Join(quoted, ", ") + ` FROM "` + table + `"`)
	if err != nil {
		return 0, fmt.Errorf("failed to read sqlite %s: %w", table, err)
	}
	defer rows.Close()
	tx, err := pgdb.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin import tx for %s: %w", table, err)
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()
	_, _ = tx.ExecContext(ctx, `SET LOCAL session_replication_role = 'replica'`)
	copied := 0
	batch := make([][]any, 0, 250)
	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		n, err := insertBatch(ctx, tx, table, columns, batch)
		copied += n
		batch = batch[:0]
		return err
	}
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			return 0, fmt.Errorf("failed to scan sqlite %s: %w", table, err)
		}
		converted := make([]any, len(columns))
		for i, col := range columns {
			value, err := convertValue(values[i], pgCols[col])
			if err != nil {
				return 0, fmt.Errorf("failed to convert %s.%s: %w", table, col, err)
			}
			converted[i] = value
		}
		batch = append(batch, converted)
		if len(batch) >= 250 {
			if err := flush(); err != nil {
				return 0, fmt.Errorf("failed to insert into %s: %w", table, err)
			}
		}
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("failed to read sqlite %s: %w", table, err)
	}
	if err := flush(); err != nil {
		return 0, fmt.Errorf("failed to insert into %s: %w", table, err)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit import of %s: %w", table, err)
	}
	committed = true
	slog.Info("imported table", "table", table, "rows", copied)
	return copied, nil
}

func pgColumns(ctx context.Context, pgdb *sql.DB, table string) (map[string]pgColumn, error) {
	rows, err := pgdb.QueryContext(ctx, `SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_schema = 'public' AND table_name = $1`, table)
	if err != nil {
		return nil, fmt.Errorf("failed to describe %s: %w", table, err)
	}
	defer rows.Close()
	columns := map[string]pgColumn{}
	for rows.Next() {
		var name, dataType, nullable string
		if err := rows.Scan(&name, &dataType, &nullable); err != nil {
			return nil, fmt.Errorf("failed to scan column of %s: %w", table, err)
		}
		columns[name] = pgColumn{dataType: dataType, isNullable: nullable == "YES"}
	}
	return columns, rows.Err()
}

func sqliteColumns(litedb *sql.DB, table string) ([]string, error) {
	rows, err := litedb.Query(`SELECT name FROM pragma_table_info("` + table + `") ORDER BY cid`)
	if err != nil {
		return nil, fmt.Errorf("failed to describe sqlite %s: %w", table, err)
	}
	defer rows.Close()
	var columns []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan sqlite column: %w", err)
		}
		columns = append(columns, name)
	}
	return columns, rows.Err()
}

func insertBatch(ctx context.Context, tx *sql.Tx, table string, columns []string, batch [][]any) (int, error) {
	quoted := make([]string, len(columns))
	for i, col := range columns {
		quoted[i] = `"` + col + `"`
	}
	var query strings.Builder
	query.WriteString(`INSERT INTO "` + table + `" (` + strings.Join(quoted, ", ") + `) VALUES `)
	args := make([]any, 0, len(batch)*len(columns))
	for r, row := range batch {
		if r > 0 {
			query.WriteString(", ")
		}
		query.WriteString("(")
		for c, value := range row {
			if c > 0 {
				query.WriteString(", ")
			}
			query.WriteString("$")
			query.WriteString(strconv.Itoa(r*len(columns) + c + 1))
			args = append(args, value)
		}
		query.WriteString(")")
	}
	if _, err := tx.ExecContext(ctx, query.String(), args...); err != nil {
		return 0, err
	}
	return len(batch), nil
}

func convertValue(value any, col pgColumn) (any, error) {
	switch raw := value.(type) {
	case nil:
		return nil, nil
	case []byte:
		return convertValue(string(raw), col)
	case string:
		if raw == "" && col.isNullable && strings.HasPrefix(col.dataType, "timestamp") {
			return nil, nil
		}
		if col.dataType == "boolean" {
			return parseBool(raw)
		}
		return raw, nil
	case int64:
		if col.dataType == "boolean" {
			return raw != 0, nil
		}
		return raw, nil
	case float64:
		if col.dataType == "boolean" {
			return raw != 0, nil
		}
		return raw, nil
	default:
		return raw, nil
	}
}

func parseBool(raw string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "t", "yes", "y", "on":
		return true, nil
	case "0", "false", "f", "no", "n", "off", "":
		return false, nil
	default:
		return false, fmt.Errorf("unrecognized boolean %q", raw)
	}
}
