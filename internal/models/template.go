package models

import "html/template"

func NewTemplate() *template.Template {
	const pattern = "views/*"

	return template.Must(template.ParseGlob(pattern))
}
