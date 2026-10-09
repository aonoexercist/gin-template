package app

import (
	"database/sql"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yourname/gin-template/internal/health"
	"github.com/yourname/gin-template/internal/platform/config"
	"github.com/yourname/gin-template/internal/platform/server"
	"github.com/yourname/gin-template/internal/user"
)

// NewRouter builds every feature and returns the fully wired router.
// This is the only file you touch when adding a feature.
func NewRouter(cfg *config.Config, log *slog.Logger, db *gorm.DB, sqlDB *sql.DB) *gin.Engine {
	features := []server.Registrar{
		user.NewModule(db),
		// order.NewModule(db),
		// task.NewModule(db),
	}

	return server.NewRouter(cfg, log, health.NewHandler(sqlDB), features...)
}
