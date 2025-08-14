package models

import (
	"html/template"
	"log/slog"
	"sync"
)

const pattern = "views/*"

var (
	templ *template.Template
	once  sync.Once
)

func NewTemplate() *template.Template {
	if templ == nil {
		once.Do(
			func() {
				slog.Info("Loading templates...")

				templ = template.Must(template.ParseGlob(pattern))
			},
		)
	} else {
		slog.Info("Using existing template compile")
	}

	return templ
}
