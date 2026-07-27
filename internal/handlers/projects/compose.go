package projects

import (
	"context"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	databaseservices "codedock.run/codedock/internal/services/databases"
	projectservices "codedock.run/codedock/internal/services/projects"
	"codedock.run/codedock/internal/utils"
)

type ComposeHandler struct {
	projectService  *projectservices.ProjectService
	appService      *projectservices.AppService
	databaseService *databaseservices.DatabaseService
	envRepo         repositories.EnvironmentRepository
	appRepo         repositories.AppServiceRepository
	composeParser   *projectservices.ComposeParserService
}

func NewComposeHandler(
	ps *projectservices.ProjectService,
	as *projectservices.AppService,
	ds *databaseservices.DatabaseService,
	er repositories.EnvironmentRepository,
	ar repositories.AppServiceRepository,
	cp *projectservices.ComposeParserService,
) *ComposeHandler {
	return &ComposeHandler{
		projectService:  ps,
		appService:      as,
		databaseService: ds,
		envRepo:         er,
		appRepo:         ar,
		composeParser:   cp,
	}
}

type ComposeAnalyzeRequest struct {
	ComposeContent string `json:"composeContent"`
	ProjectID      string `json:"projectId"`
}

func (h *ComposeHandler) Analyze(c echo.Context) error {
	var req ComposeAnalyzeRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}

	if req.ComposeContent == "" {
		return utils.Error(c, http.StatusBadRequest, "compose content is required")
	}

	result, err := h.composeParser.Parse([]byte(req.ComposeContent), req.ProjectID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "failed to parse docker-compose: "+err.Error())
	}

	return utils.Success(c, "Compose analyzed", result)
}

func (h *ComposeHandler) Deploy(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	projectID := c.FormValue("projectId")
	if projectID == "" {
		projectID = c.FormValue("project_id")
	}
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "projectId parameter is required")
	}
	if user.Role != "admin" {
		if !h.projectService.HasPermission(c.Request().Context(), projectID, user.UserID, models.UserRole(user.Role), "") {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions for this project")
		}
	}

	composeBytes, err := h.readUploadedFile(c)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	result, err := h.composeParser.Parse(composeBytes, projectID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "compose deploy failed: "+err.Error())
	}

	createdCount, err := h.provisionComposeResources(c.Request().Context(), result)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Compose file deployed", map[string]any{
		"count": createdCount,
	})
}

func (h *ComposeHandler) readUploadedFile(c echo.Context) ([]byte, error) {
	const maxUploadBytes = 50 * 1024 * 1024
	file, err := c.FormFile("file")
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "compose file is required")
	}

	if file.Size > maxUploadBytes {
		return nil, echo.NewHTTPError(http.StatusRequestEntityTooLarge, "compose file exceeds 50 MB limit")
	}

	src, err := file.Open()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to open uploaded file")
	}
	defer src.Close()

	return io.ReadAll(io.LimitReader(src, maxUploadBytes))
}

func (h *ComposeHandler) provisionComposeResources(ctx context.Context, result *projectservices.ParsedComposeResult) (int, error) {
	var createdCount int

	for _, dbReq := range result.Databases {
		db := &models.Database{
			ProjectID:    dbReq.ProjectID,
			Name:         dbReq.Name,
			Engine:       dbReq.Engine,
			Version:      dbReq.Version,
			Port:         dbReq.Port,
			Username:     dbReq.Username,
			Password:     dbReq.Password,
			DatabaseName: dbReq.DatabaseName,
		}
		if _, err := h.databaseService.CreateDatabase(ctx, db); err != nil {
			return createdCount, echo.NewHTTPError(http.StatusInternalServerError, "failed to create database "+dbReq.Name+": "+err.Error())
		}
		createdCount++
	}

	for _, appReq := range result.AppServices {
		app := &models.AppService{
			ProjectID:      appReq.ProjectID,
			Name:           appReq.Name,
			RuntimeMode:    appReq.RuntimeMode,
			DockerfilePath: appReq.DockerfilePath,
			InstallCommand: appReq.InstallCommand,
			BuildCommand:   appReq.BuildCommand,
			StartCommand:   appReq.StartCommand,
			RepositoryURL:  appReq.RepositoryURL,
			ImageRef:       appReq.ImageRef,
		}
		if appReq.BuildEngine != "" {
			app.BuildEngine = models.BuildEngine(appReq.BuildEngine)
		}

		if _, err := h.appService.CreateAppService(ctx, app); err != nil {
			return createdCount, echo.NewHTTPError(http.StatusInternalServerError, "failed to create app service "+app.Name+": "+err.Error())
		}
		createdCount++
	}

	return createdCount, nil
}
