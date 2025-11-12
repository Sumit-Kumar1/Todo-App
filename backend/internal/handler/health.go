package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"todoapp/internal/models"
	"todoapp/internal/server"
)

type HealthHandler struct {
	srv *server.Server
}

func NewHealthHandler(srv *server.Server) *HealthHandler {
	return &HealthHandler{srv: srv}
}

// HealthHandler handles health check requests
func (s *HealthHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Database:  true,
		Service:   true,
	}

	// Check database health
	if err := s.checkDatabaseHealth(r.Context()); err != nil {
		response.Status = "unhealthy"
		response.Database = false
		response.Message = "Database connection failed: " + err.Error()
		s.srv.Logger.Error("Database health check failed", "error", err)
	}

	w.Header().Set("Content-Type", "application/json")

	if response.Status == "unhealthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.srv.Logger.Error("Failed to encode health response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// ReadyHandler handles readiness checks
func (s *HealthHandler) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status:    "ready",
		Timestamp: time.Now(),
		Database:  true,
		Service:   true,
	}

	// Check database readiness
	if err := s.checkDatabaseHealth(r.Context()); err != nil {
		response.Status = "not ready"
		response.Database = false
		response.Message = "Database not ready: " + err.Error()
		s.srv.Logger.Error("Database readiness check failed", "error", err)
	}

	w.Header().Set("Content-Type", "application/json")

	if response.Status == "not ready" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.srv.Logger.Error("Failed to encode readiness response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// LivenessHandler handles liveness checks
func (s *HealthHandler) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.srv.Logger.Error("Failed to encode liveness response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// checkDatabaseHealth verifies database connectivity
func (s *HealthHandler) checkDatabaseHealth(ctx context.Context) error {
	if s.srv.DB == nil {
		return context.Canceled
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Simple query to check database connectivity
	var result int
	err := s.srv.DB.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return err
	}

	if result != 1 {
		return context.Canceled
	}

	return nil
}
