package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
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
)

type Handler struct {
	h slog.Handler
	b *bytes.Buffer
	m *sync.Mutex
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name), b: h.b, m: h.m}
}

// nolint:gocritic // can't change it as it has this definition in slog library
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = colorize(darkGray, level)
	case slog.LevelInfo:
		level = colorize(cyan, level)
	case slog.LevelWarn:
		level = colorize(lightYellow, level)
	case slog.LevelError:
		level = colorize(lightRed, level)
	}

	attrs, err := h.computeAttrs(ctx, &r)
	if err != nil {
		return err
	}

	logData, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshaling attrs: %w", err)
	}

	_, err = fmt.Println(
		colorize(lightMagenta, r.Time.Format(time.RFC3339Nano)),
		level,
		colorize(white, r.Message),
		colorize(darkGray, string(logData)),
	)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) computeAttrs(ctx context.Context, r *slog.Record) (map[string]any, error) {
	h.m.Lock()

	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()

	if err := h.h.Handle(ctx, *r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any

	if err := json.Unmarshal(h.b.Bytes(), &attrs); err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}

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

func newHandler(opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	b := &bytes.Buffer{}

	return &Handler{
		b: b,
		m: &sync.Mutex{},
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
	}
}

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

func newLogger() *slog.Logger {
	var (
		leveler slog.Level
		level   = os.Getenv("LOG_LEVEL")
	)

	switch level {
	case "ERROR":
		leveler = slog.LevelError
	case "DEBUG":
		leveler = slog.LevelDebug
	case "WARN":
		leveler = slog.LevelWarn
	default:
		leveler = slog.LevelInfo
	}

	replaceFn := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == "time" {
			a.Value = slog.StringValue(time.Now().Format(time.RFC1123))
		}

		return a
	}

	logger := slog.New(newHandler(&slog.HandlerOptions{
		AddSource:   false,
		Level:       leveler,
		ReplaceAttr: replaceFn,
	}))

	slog.SetDefault(logger)

	return logger
}
