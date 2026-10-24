package projects

import (
	"codedock/internal/models"
	"codedock/internal/utils"
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/url"
)

func (h *AppHandler) prepareApplicationDomain(ctx context.Context, app *models.AppService) error {
	if app.ID == "" {
		app.ID = uuid.NewString()
	}
	if app.RuntimeMode == models.RuntimeModeWorker {
		if app.Domain != "" {
			return fmt.Errorf("workers cannot have an HTTP domain")
		}
		return nil
	}
	if app.Domain == "" {
		parsed, err := url.Parse(utils.GenerateAppDomain(app.ID+"-"+app.Name, "", ""))
		if err != nil || parsed.Hostname() == "" {
			return fmt.Errorf("generate application domain")
		}
		app.Domain = parsed.Hostname()
	}
	canonical, err := h.envService.ValidateApplicationDomain(ctx, app.ID, app.Domain)
	if err != nil {
		return err
	}
	app.Domain = canonical
	return nil
}

func (h *AppHandler) ensureApplicationDomain(ctx context.Context, app *models.AppService) error {
	domains, err := h.envService.ListDomainsByService(ctx, app.ID)
	if err != nil {
		return fmt.Errorf("load application domains: %w", err)
	}
	for _, domain := range domains {
		if domain.DomainName == app.Domain {
			return nil
		}
	}
	if _, err := h.envService.CreateGeneratedDomain(ctx, &models.DomainConfig{ServiceID: app.ID, DomainName: app.Domain}); err != nil {
		return fmt.Errorf("save application domain: %w", err)
	}
	return nil
}
