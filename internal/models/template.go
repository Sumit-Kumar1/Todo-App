package models

import "html/template"

func NewTemplate() *template.Template {
	pattern := "views/*"

	return template.Must(template.ParseGlob(pattern))
}
