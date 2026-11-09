package migrations

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"codedock/internal/modules/databases"
	systemdb "codedock/internal/platform/database"
	"codedock/internal/platform/networking"
)

const bundleManifestVersion = "1"

type BundleManifest struct {
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
	Databases   []string  `json:"databases"`
	HasSystemDB bool      `json:"hasSystemDb"`
}

type DatabaseLister interface {
	ListAll(ctx context.Context, search string) ([]*databases.Database, error)
}

type BundleServiceImplementation struct {
	db          *sql.DB
	dbRepo      DatabaseLister
	dataDir     string
	databaseURL string
}

func NewBundleService(db *sql.DB, dbRepo DatabaseLister, dataDir, databaseURL string) *BundleServiceImplementation {
	return &BundleServiceImplementation{db: db, dbRepo: dbRepo, dataDir: dataDir, databaseURL: databaseURL}
}

func (s *BundleServiceImplementation) Export(ctx context.Context, passphrase string) ([]byte, error) {
	files := make(map[string][]byte)

	systemData, err := s.dumpSystemDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("system database dump failed: %w", err)
	}
	files["systemdb.sql"] = systemData

	var dbs []*databases.Database
	if s.dbRepo != nil {
		dbs, err = s.dbRepo.ListAll(ctx, "")
		if err != nil {
			return nil, fmt.Errorf("failed to list databases: %w", err)
		}
	}

	manifest := BundleManifest{
		Version:     bundleManifestVersion,
		CreatedAt:   time.Now().UTC(),
		HasSystemDB: true,
	}

	for _, db := range dbs {
		dump, err := s.dumpDatabase(ctx, db)
		if err != nil {
			continue
		}
		ext := dumpExtension(string(db.Engine))
		filename := fmt.Sprintf("databases/%s%s", db.Name, ext)
		files[filename] = dump
		manifest.Databases = append(manifest.Databases, db.Name)
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manifest: %w", err)
	}
	files["manifest.json"] = manifestData

	var tarBuf bytes.Buffer
	if err := networking.CreateTarGz(&tarBuf, files); err != nil {
		return nil, fmt.Errorf("failed to create bundle archive: %w", err)
	}

	var encBuf bytes.Buffer
	if err := networking.EncryptBundle(&tarBuf, &encBuf, passphrase); err != nil {
		return nil, fmt.Errorf("failed to encrypt bundle: %w", err)
	}

	return encBuf.Bytes(), nil
}

func (s *BundleServiceImplementation) Import(ctx context.Context, bundleData []byte, passphrase string) (any, error) {
	var decBuf bytes.Buffer
	if err := networking.DecryptBundle(bytes.NewReader(bundleData), &decBuf, passphrase); err != nil {
		return nil, err
	}

	files, err := networking.ExtractTarGz(&decBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to extract bundle: %w", err)
	}

	manifestData, ok := files["manifest.json"]
	if !ok {
		return nil, fmt.Errorf("bundle is missing manifest.json")
	}
	var manifest BundleManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	if sqlData, ok := files["systemdb.sql"]; ok {
		if err := s.restoreSystemDB(ctx, sqlData); err != nil {
			return nil, fmt.Errorf("system database restore failed: %w", err)
		}
	}

	var dbs []*databases.Database
	if s.dbRepo != nil {
		dbs, _ = s.dbRepo.ListAll(ctx, "")
	}
	dbsByName := make(map[string]*databases.Database, len(dbs))
	for _, db := range dbs {
		dbsByName[db.Name] = db
	}

	for _, dbName := range manifest.Databases {
		db, found := dbsByName[dbName]
		if !found {
			continue
		}
		for _, ext := range []string{".sql", ".rdb", ".dump"} {
			key := fmt.Sprintf("databases/%s%s", dbName, ext)
			if data, ok := files[key]; ok {
				_ = s.restoreDatabase(ctx, db, data, ext)
				break
			}
		}
	}

	return &manifest, nil
}

func (s *BundleServiceImplementation) dumpSystemDB(ctx context.Context) ([]byte, error) {
	if s.databaseURL == "" {
		password, err := systemdb.ReadPassword(s.dataDir)
		if err != nil {
			return nil, err
		}
		return systemdb.DumpContainer(ctx, systemdb.ContainerName, systemdb.DBUser, systemdb.DBName, password)
	}
	return systemdb.DumpURL(ctx, s.databaseURL)
}

func (s *BundleServiceImplementation) restoreSystemDB(ctx context.Context, sqlData []byte) error {
	if s.db == nil {
		return fmt.Errorf("system database handle is required for restore")
	}
	if _, err := s.db.ExecContext(ctx, `DROP SCHEMA public CASCADE`); err != nil {
		return fmt.Errorf("drop schema for restore: %w", err)
	}
	if _, err := s.db.ExecContext(ctx, `CREATE SCHEMA public`); err != nil {
		return fmt.Errorf("recreate schema for restore: %w", err)
	}
	if s.databaseURL == "" {
		password, err := systemdb.ReadPassword(s.dataDir)
		if err != nil {
			return err
		}
		return systemdb.RestoreContainer(ctx, systemdb.ContainerName, systemdb.DBUser, systemdb.DBName, password, sqlData)
	}
	return systemdb.RestoreURL(ctx, s.databaseURL, sqlData)
}

func (s *BundleServiceImplementation) dumpDatabase(_ context.Context, db *databases.Database) ([]byte, error) {
	containerName := db.InternalDNS
	if containerName == "" {
		return nil, fmt.Errorf("no internal DNS for database %s", db.Name)
	}

	switch db.Engine {
	case "postgres", "timescaledb":
		cmd := exec.Command("docker", "exec", containerName,
			"pg_dump", "-U", db.Username, db.DatabaseName)
		return cmd.Output()
	case "mysql", "mariadb":
		cmd := exec.Command("docker", "exec", "-e", "MYSQL_PWD="+db.Password, containerName,
			"mysqldump", "-u", db.Username, db.DatabaseName)
		return cmd.Output()
	case "mongodb":
		cmd := exec.Command("docker", "exec", containerName,
			"mongodump", "--archive", "--authenticationDatabase=admin",
			fmt.Sprintf("--username=%s", db.Username), fmt.Sprintf("--password=%s", db.Password))
		return cmd.Output()
	case "clickhouse":
		cmd := exec.Command("docker", "exec", containerName,
			"clickhouse-client", "--user", db.Username, "--password", db.Password,
			"--query", fmt.Sprintf("SELECT * FROM %s FORMAT Native", db.DatabaseName))
		return cmd.Output()
	case "redis":
		cmd := exec.Command("docker", "exec", containerName,
			"redis-cli", "--rdb", "/tmp/dump.rdb")
		if err := cmd.Run(); err != nil {
			return nil, err
		}
		return exec.Command("docker", "exec", containerName, "cat", "/tmp/dump.rdb").Output()
	default:
		return nil, fmt.Errorf("unsupported engine for dump: %s", db.Engine)
	}
}

func (s *BundleServiceImplementation) restoreDatabase(_ context.Context, db *databases.Database, data []byte, ext string) error {
	containerName := db.InternalDNS
	if containerName == "" {
		return fmt.Errorf("no internal DNS for database %s", db.Name)
	}

	switch db.Engine {
	case "postgres", "timescaledb":
		cmd := exec.Command("docker", "exec", "-i", containerName,
			"psql", "-U", db.Username, db.DatabaseName)
		cmd.Stdin = bytes.NewReader(data)
		return cmd.Run()
	case "mysql", "mariadb":
		cmd := exec.Command("docker", "exec", "-i", "-e", "MYSQL_PWD="+db.Password, containerName,
			"mysql", "-u", db.Username, db.DatabaseName)
		cmd.Stdin = bytes.NewReader(data)
		return cmd.Run()
	case "mongodb":
		cmd := exec.Command("docker", "exec", "-i", containerName,
			"mongorestore", "--archive", "--authenticationDatabase=admin",
			fmt.Sprintf("--username=%s", db.Username), fmt.Sprintf("--password=%s", db.Password))
		cmd.Stdin = bytes.NewReader(data)
		return cmd.Run()
	case "redis":
		if ext != ".rdb" {
			return fmt.Errorf("redis restore only supports .rdb format")
		}
		copyCmd := exec.Command("docker", "cp", "-", fmt.Sprintf("%s:/tmp/codedock-import.rdb", containerName))
		copyCmd.Stdin = bytes.NewReader(data)
		if err := copyCmd.Run(); err != nil {
			return err
		}
		return exec.Command("docker", "exec", containerName,
			"redis-cli", "DEBUG", "RELOAD").Run()
	default:
		return fmt.Errorf("unsupported engine for restore: %s", db.Engine)
	}
}

func dumpExtension(engine string) string {
	switch engine {
	case "redis":
		return ".rdb"
	default:
		return ".sql"
	}
}
