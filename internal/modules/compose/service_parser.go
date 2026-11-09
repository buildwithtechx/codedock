package compose

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type ComposeParserService struct{}

func NewComposeParserService() *ComposeParserService {
	return &ComposeParserService{}
}

type ParsedComposeResult struct {
	AppServices []CreateAppServiceRequest `json:"appServices"`
	Databases   []CreateDatabaseRequest   `json:"databases"`
}

func (s *ComposeParserService) Parse(composeData []byte, projectID string) (*ParsedComposeResult, error) {
	var userCompose UserComposeFile
	if err := yaml.Unmarshal(composeData, &userCompose); err != nil {
		return nil, err
	}

	result := &ParsedComposeResult{
		AppServices: make([]CreateAppServiceRequest, 0),
		Databases:   make([]CreateDatabaseRequest, 0),
	}

	for name, svc := range userCompose.Services {
		imageLower := strings.ToLower(svc.Image)
		engine := detectDatabaseEngine(imageLower)

		if engine != "" {
			dbReq := CreateDatabaseRequest{
				ProjectID: projectID,
				Name:      name,
				Engine:    engine,
				Version:   extractVersion(imageLower),
				Port:      getDefaultPort(engine),
			}
			result.Databases = append(result.Databases, dbReq)
		} else {
			appReq := CreateAppServiceRequest{
				ProjectID:   projectID,
				Name:        name,
				ImageRef:    svc.Image,
				RuntimeMode: RuntimeModeWeb,
			}
			if svc.Build != nil {
				appReq.BuildEngine = string(BuildEngineDockerfile)
			}
			result.AppServices = append(result.AppServices, appReq)
		}
	}

	return result, nil
}

func detectDatabaseEngine(image string) DatabaseEngine {
	if strings.Contains(image, "postgres") {
		return DatabaseEnginePostgres
	}
	if strings.Contains(image, "mysql") {
		return DatabaseEngineMySQL
	}
	if strings.Contains(image, "redis") {
		return DatabaseEngineRedis
	}
	if strings.Contains(image, "mongo") {
		return DatabaseEngineMongoDB
	}
	if strings.Contains(image, "mariadb") {
		return DatabaseEngineMariaDB
	}
	if strings.Contains(image, "clickhouse") {
		return DatabaseEngineClickhouse
	}
	return ""
}

func extractVersion(image string) string {
	parts := strings.Split(image, ":")
	if len(parts) > 1 {
		return parts[1]
	}
	return "latest"
}

func getDefaultPort(engine DatabaseEngine) int {
	switch engine {
	case DatabaseEnginePostgres:
		return 5432
	case DatabaseEngineMySQL:
		return 3306
	case DatabaseEngineRedis:
		return 6379
	case DatabaseEngineMongoDB:
		return 27017
	case DatabaseEngineMariaDB:
		return 3306
	case DatabaseEngineClickhouse:
		return 9000
	default:
		return 0
	}
}
