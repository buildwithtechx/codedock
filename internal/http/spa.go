package http

import (
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"

	"codedock/apps/dashboard"
	"codedock/internal/config"
)

// SetupSPAFallback sets up static dashboard serving or embedded dashboard handlers.
func SetupSPAFallback(e *echo.Echo) {
	staticDir := config.Get().Server.StaticDir

	if staticDir != "" {
		if stat, err := os.Stat(staticDir); err == nil && stat.IsDir() {
			e.GET("/*", func(c echo.Context) error {
				path := filepath.Join(staticDir, filepath.Clean(c.Request().URL.Path))
				if _, err := os.Stat(path); os.IsNotExist(err) {
					return c.File(filepath.Join(staticDir, "index.html"))
				}
				return c.File(path)
			})
			return
		}
	}

	dashboard.RegisterHandlers(e)
}
