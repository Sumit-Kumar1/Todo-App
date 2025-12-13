package models

import "time"

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Database  bool      `json:"database"`
	Service   bool      `json:"service"`
	Message   string    `json:"message,omitempty"`
}
