package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rentoutdoor/api/internal/infrastructure/persistence/mysql"
	"github.com/rentoutdoor/api/pkg/response"
	"gorm.io/gorm"
)

// HealthHandler handles health check endpoints.
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health godoc
// @Summary Health check
// @Tags Health
// @Produce json
// @Success 200 {object} response.Response
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	checks := map[string]string{
		"status": "healthy",
	}

	if err := mysql.HealthCheck(h.db); err != nil {
		checks["database"] = "unhealthy"
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"checks": checks,
		})
		return
	}
	checks["database"] = "healthy"

	response.Success(c, http.StatusOK, "service healthy", checks)
}
