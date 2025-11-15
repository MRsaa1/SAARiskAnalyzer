package handlers

import (
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/reserveone/saa-risk-analyzer/internal/domain"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, domain.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	})
}
