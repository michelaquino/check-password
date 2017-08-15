package templates

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

type HtmlTemplate struct {
	templates *template.Template
}

func (h *HtmlTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return h.templates.ExecuteTemplate(w, name, data)
}

var ViewTemplates = &HtmlTemplate{
	templates: template.Must(template.ParseGlob("public/view/*.html")),
}
