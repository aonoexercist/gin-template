package server

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/yourname/gin-template/internal/health"
	"github.com/yourname/gin-template/internal/platform/config"
	"github.com/yourname/gin-template/internal/platform/middleware"
)

// Registrar is implemented by every feature that exposes HTTP routes.
type Registrar interface {
	RegisterRoutes(rg *gin.RouterGroup)
}

func NewRouter(cfg *config.Config, log *slog.Logger, h *health.Handler, features ...Registrar) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Logger(log), middleware.Recovery(log), middleware.CORS())

	r.GET("/healthz", h.Liveness)
	r.GET("/readyz", h.Readiness)

	v1 := r.Group("/api/v1")
	for _, f := range features {
		f.RegisterRoutes(v1)
	}
	return r
}
