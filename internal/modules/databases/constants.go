package databases

type DatabaseEngine string
type DatabaseStatus string

const (
	DatabaseEnginePostgres   DatabaseEngine = "postgres"
	DatabaseEngineMySQL      DatabaseEngine = "mysql"
	DatabaseEngineRedis      DatabaseEngine = "redis"
	DatabaseEngineMongoDB    DatabaseEngine = "mongodb"
	DatabaseEngineMariaDB    DatabaseEngine = "mariadb"
	DatabaseEngineClickhouse DatabaseEngine = "clickhouse"
)

const (
	DatabaseStatusCreated DatabaseStatus = "created"
	DatabaseStatusRunning DatabaseStatus = "running"
	DatabaseStatusStopped DatabaseStatus = "stopped"
	DatabaseStatusError   DatabaseStatus = "error"
)
