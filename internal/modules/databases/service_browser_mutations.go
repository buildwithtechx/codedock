package databases

import (
	"context"
	"fmt"
	"strings"
)

func escapeSQLValue(v any) string {
	if v == nil {
		return "NULL"
	}
	str := fmt.Sprintf("%v", v)
	str = strings.ReplaceAll(str, "'", "''")
	return fmt.Sprintf("'%s'", str)
}

func (s *Service) InsertTableRow(ctx context.Context, id, table string, data map[string]any) (*DatabaseQueryResponse, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, ErrDatabaseNotFound
	}

	query, err := buildInsertSQL(string(db.Engine), table, data)
	if err != nil {
		return nil, err
	}
	return s.QueryDatabase(ctx, id, query, db.ProjectID)
}

func (s *Service) UpdateTableRow(ctx context.Context, opts UpdateTableRowOpts) (*DatabaseQueryResponse, error) {
	db, err := s.repo.GetByID(ctx, opts.ID)
	if err != nil || db == nil {
		return nil, ErrDatabaseNotFound
	}

	query, err := buildUpdateSQL(string(db.Engine), opts.Table, opts.Keys, opts.Data)
	if err != nil {
		return nil, err
	}
	return s.QueryDatabase(ctx, opts.ID, query, db.ProjectID)
}

func (s *Service) DeleteTableRow(ctx context.Context, id, table string, keys map[string]any) (*DatabaseQueryResponse, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, ErrDatabaseNotFound
	}

	query, err := buildDeleteSQL(string(db.Engine), table, keys)
	if err != nil {
		return nil, err
	}
	return s.QueryDatabase(ctx, id, query, db.ProjectID)
}

func buildInsertSQL(engine, table string, data map[string]any) (string, error) {
	var cols []string
	var vals []string
	isPostgres := engine == "postgresql" || engine == "postgres"
	for k, v := range data {
		if isPostgres {
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
		return "", ErrPrimaryKeyRequired
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
		return "", ErrPrimaryKeyRequired
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
