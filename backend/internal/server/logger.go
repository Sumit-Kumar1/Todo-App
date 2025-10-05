package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	reset = "\033[0m"

	cyan         = 36
	darkGray     = 90
	lightRed     = 91
	lightYellow  = 93
	lightMagenta = 95
	white        = 97

	logDir = "logs"
)

var (
	levelColors = map[slog.Level]int{
		slog.LevelDebug: darkGray,
		slog.LevelInfo:  cyan,
		slog.LevelWarn:  lightYellow,
		slog.LevelError: lightRed,
	}

	levelIcons = map[slog.Level]string{
		slog.LevelDebug: "🔍",
		slog.LevelInfo:  "ℹ️",
		slog.LevelWarn:  "⚠️",
		slog.LevelError: "❌",
	}
)

// Buffer pool to reduce allocations
var bufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

type Handler struct {
	h            slog.Handler
	m            *sync.Mutex
	enableColors bool
	showIcons    bool
	showSource   bool
}

type HandlerOptions struct {
	EnableColors bool
	ShowIcons    bool
	ShowSource   bool
	Level        slog.Leveler
	AddSource    bool
	ReplaceAttr  func(groups []string, a slog.Attr) slog.Attr
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		h:            h.h.WithAttrs(attrs),
		m:            h.m,
		enableColors: h.enableColors,
		showIcons:    h.showIcons,
		showSource:   h.showSource,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		h:            h.h.WithGroup(name),
		m:            h.m,
		enableColors: h.enableColors,
		showIcons:    h.showIcons,
		showSource:   h.showSource,
	}
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	var levelText string
	if h.showIcons {
		levelText = levelIcons[r.Level] + " " + r.Level.String()
	} else {
		levelText = r.Level.String() + ":"
	}

	if h.enableColors {
		if colorCode, ok := levelColors[r.Level]; ok {
			levelText = colorize(colorCode, levelText)
		}
	}

	timeText := r.Time.Format("2006-01-02 15:04:05.000")
	if h.enableColors {
		timeText = colorize(lightMagenta, timeText)
	}

	messageText := r.Message
	if h.enableColors {
		messageText = colorize(white, messageText)
	}

	var sourceText string
	if h.showSource && r.PC != 0 {
		if fs := runtime.CallersFrames([]uintptr{r.PC}); fs != nil {
			if frame, more := fs.Next(); more {
				sourceText = fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
				if h.enableColors {
					sourceText = colorize(darkGray, sourceText)
				}
			}
		}
	}

	attrs, err := h.computeAttrs(ctx, &r)
	if err != nil {
		return err
	}

	var logLine string
	if sourceText != "" {
		logLine = fmt.Sprintf("%s %s %s [%s]", timeText, levelText, messageText, sourceText)
	} else {
		logLine = fmt.Sprintf("%s %s %s", timeText, levelText, messageText)
	}

	if len(attrs) > 0 {
		logData, err := json.Marshal(attrs)
		if err != nil {
			return fmt.Errorf("error when marshaling attrs: %w", err)
		}

		if h.enableColors {
			logLine += " " + colorize(darkGray, string(logData))
		} else {
			logLine += " " + string(logData)
		}
	}

	_, err = fmt.Println(logLine)
	return err
}

func (h *Handler) computeAttrs(ctx context.Context, r *slog.Record) (map[string]any, error) {
	h.m.Lock()
	defer h.m.Unlock()

	buf := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufferPool.Put(buf)
	}()

	var level slog.Level
	if h.h.Enabled(ctx, r.Level) {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	tempHandler := slog.NewJSONHandler(buf, &slog.HandlerOptions{
		Level: level,
	})

	if err := tempHandler.Handle(ctx, *r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	if err := json.Unmarshal(buf.Bytes(), &attrs); err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}

	// Remove standard fields that we're already displaying
	delete(attrs, "time")
	delete(attrs, "level")
	delete(attrs, "msg")
	delete(attrs, "source")

	return attrs, nil
}

func suppressDefaults(
	next func([]string, slog.Attr) slog.Attr,
) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}

		if next == nil {
			return a
		}

		return next(groups, a)
	}
}

func newHandler(opts *HandlerOptions) *Handler {
	if opts == nil {
		opts = &HandlerOptions{
			EnableColors: true,
			ShowIcons:    true,
			ShowSource:   true,
		}
	}

	innerOpts := &slog.HandlerOptions{
		Level:       opts.Level,
		AddSource:   opts.AddSource,
		ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
	}

	return &Handler{
		h:            slog.NewJSONHandler(&bytes.Buffer{}, innerOpts),
		m:            &sync.Mutex{},
		enableColors: opts.EnableColors,
		showIcons:    opts.ShowIcons,
		showSource:   opts.ShowSource,
	}
}

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

type LoggerConfig struct {
	Level      string
	EnableJSON bool
	ShowSource bool
	ShowIcons  bool
	UseColors  bool
}

func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:      "INFO",
		EnableJSON: true,
		ShowSource: true,
		ShowIcons:  true,
		UseColors:  true,
	}
}

func newLoggerWithConfig(cfg *LoggerConfig) (*slog.Logger, *os.File) {
	if cfg == nil {
		cfg = NewLoggerConfig()
	}

	// Create logs directory
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		log.Printf("failed to create logs directory: %v", err)
	}

	// Create log file with timestamp
	logFileName := filepath.Join(logDir, fmt.Sprintf("%s.log", time.Now().Format("20060102_150405")))
	logFile, err := os.Create(logFileName)
	if err != nil {
		log.Fatal("Failed to create log file:", err)
	}

	// Parse log level
	var leveler slog.Level
	switch cfg.Level {
	case "DEBUG":
		leveler = slog.LevelDebug
	case "WARN":
		leveler = slog.LevelWarn
	case "ERROR":
		leveler = slog.LevelError
	default:
		leveler = slog.LevelInfo
	}

	// Common handler options
	handlerOpts := &slog.HandlerOptions{
		AddSource: cfg.ShowSource,
		Level:     leveler,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if timeVal, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(timeVal.Format(time.RFC3339))
				}
			}
			return a
		},
	}

	// File handler (always JSON)
	fileHandler := slog.NewJSONHandler(logFile, handlerOpts)

	// Terminal handler
	termHandlerOpts := &HandlerOptions{
		EnableColors: cfg.UseColors && isTerminal(),
		ShowIcons:    cfg.ShowIcons,
		ShowSource:   cfg.ShowSource,
		Level:        leveler,
		AddSource:    cfg.ShowSource,
	}
	termHandler := newHandler(termHandlerOpts)

	// Create multi-handler logger
	logger := slog.New(NewMultiHandler(fileHandler, termHandler))
	slog.SetDefault(logger)

	return logger, logFile
}

func newLogger() (*slog.Logger, *os.File) {
	cfg := NewLoggerConfig()
	cfg.Level = GetEnvOrDefault("LOG_LEVEL", "INFO")

	// Check if we should disable colors (e.g., in CI environments)
	if GetEnvOrDefault("NO_COLOR", false) || GetEnvOrDefault("CI", false) {
		cfg.UseColors = false
	}

	if GetEnvOrDefault("SHOW_ICONS", true) {
		cfg.ShowIcons = true
	}

	return newLoggerWithConfig(cfg)
}

func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}

	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	var lastErr error
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r); err != nil {
				lastErr = err
			}
		}
	}
	return lastErr
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return NewMultiHandler(handlers...)
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return NewMultiHandler(handlers...)
}

func LogWithFields(logger *slog.Logger, level slog.Level, msg string, fields map[string]any) {
	if logger == nil {
		logger = slog.Default()
	}

	attrs := make([]slog.Attr, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}

	logger.LogAttrs(context.Background(), level, msg, attrs...)
}

func LogError(logger *slog.Logger, msg string, err error, fields ...map[string]any) {
	allFields := map[string]any{"error": err.Error()}
	if len(fields) > 0 {
		for k, v := range fields[0] {
			allFields[k] = v
		}
	}

	LogWithFields(logger, slog.LevelError, msg, allFields)
}

func LogInfo(logger *slog.Logger, msg string, fields ...map[string]any) {
	var allFields map[string]any
	if len(fields) > 0 {
		allFields = fields[0]
	}

	LogWithFields(logger, slog.LevelInfo, msg, allFields)
}
