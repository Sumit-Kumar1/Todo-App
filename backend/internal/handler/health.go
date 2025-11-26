package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"
	"todoapp/internal/models"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	DB *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		DB: db,
	}
}

// HealthHandler handles health check requests
func (s *HealthHandler) HealthHandler(c *gin.Context) {
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Database:  true,
		Service:   true,
	}

	// Check database health
	if err := s.checkDatabaseHealth(c); err != nil {
		response.Status = "unhealthy"
		response.Database = false
		response.Message = "Database connection failed: " + err.Error()
	}

	if response.Status == "unhealthy" {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, response)
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

// ReadyHandler handles readiness checks
func (s *HealthHandler) ReadyHandler(c *gin.Context) {
	response := models.HealthResponse{
		Status:    "ready",
		Timestamp: time.Now(),
		Database:  true,
		Service:   true,
	}

	// Check database readiness
	if err := s.checkDatabaseHealth(c); err != nil {
		response.Status = "not ready"
		response.Database = false
		response.Message = "Database not ready: " + err.Error()
	}

	if response.Status == "not ready" {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, response)
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

// LivenessHandler handles liveness checks
func (s *HealthHandler) LivenessHandler(c *gin.Context) {
	response := map[string]any{
		"status":    "alive",
		"timestamp": time.Now(),
	}

	c.IndentedJSON(http.StatusOK, response)
}

// checkDatabaseHealth verifies database connectivity
func (s *HealthHandler) checkDatabaseHealth(c *gin.Context) error {
	if s.DB == nil {
		return context.Canceled
	}

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	// Simple query to check database connectivity
	var result int
	err := s.DB.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return err
	}

	if result != 1 {
		return context.Canceled
	}

	return nil
}
