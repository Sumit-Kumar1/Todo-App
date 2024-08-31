package models

import "html/template"

func NewTemplate() *template.Template {
	pattern := "views/*.html"

	return template.Must(template.ParseGlob(pattern))
}
