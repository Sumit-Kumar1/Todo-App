package server

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"testing"
)

func TestLogger(t *testing.T) {
	// Test with default configuration
	logger, logFile := newLogger()
	defer logFile.Close()

	// Test different log levels
	logger.Debug("Debug message", "key", "value")
	logger.Info("Info message", "key", "value")
	logger.Warn("Warning message", "key", "value")
	logger.Error("Error message", "key", "value")

	// Test structured logging helpers
	LogInfo(logger, "Structured info", map[string]any{"user_id": 123, "action": "login"})
	LogError(logger, "Structured error",
		fmt.Errorf("database connection failed"),
		map[string]any{"attempt": 3, "timeout": "30s"})

	// Verify logger is not nil
	if logger == nil {
		t.Error("Logger should not be nil")
	}
}

func TestLoggerConfig(t *testing.T) {
	// Test with custom configuration
	cfg := &LoggerConfig{
		Level:      "DEBUG",
		EnableJSON: false,
		ShowSource: true,
		ShowIcons:  true,
		UseColors:  false, // Disable colors for testing
	}

	logger, logFile := newLoggerWithConfig(cfg)
	defer logFile.Close()

	// Test that debug level is enabled
	if !logger.Enabled(context.Background(), slog.LevelDebug) {
		t.Error("Debug level should be enabled")
	}

	logger.Debug("Debug message with custom config")
}

func TestMultiHandler(t *testing.T) {
	// Create a buffer to capture output
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	handler1 := slog.NewJSONHandler(buf1, nil)
	handler2 := slog.NewTextHandler(buf2, nil)

	multiHandler := NewMultiHandler(handler1, handler2)
	logger := slog.New(multiHandler)

	// Log a message
	logger.Info("Test message", "key", "value")

	// Check that both handlers received the message
	if buf1.Len() == 0 {
		t.Error("First handler should have received the message")
	}
	if buf2.Len() == 0 {
		t.Error("Second handler should have received the message")
	}
}

func TestColorize(t *testing.T) {
	result := colorize(lightRed, "error")
	expected := "\033[91merror\033[0m"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestIsTerminal(t *testing.T) {
	// This is a basic test - in practice, isTerminal() depends on the environment
	// We just verify it doesn't panic
	_ = isTerminal()
}

func TestLogHelpers(t *testing.T) {
	logger, logFile := newLogger()
	defer logFile.Close()

	// Test LogWithFields with nil logger (should use default)
	LogWithFields(nil, slog.LevelInfo, "Test with nil logger", map[string]any{"test": true})

	// Test LogInfo helper
	LogInfo(logger, "Info helper test", map[string]any{"count": 42})

	// Test LogError helper
	LogError(logger, "Error helper test", fmt.Errorf("test error"), map[string]any{"code": 500})
}
