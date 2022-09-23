package confluence

import (
	_ "embed"
)

//go:embed template.html
var htmlTemplate string

func NewHTMLTemplateRetriver() TemplateRetriver {
	return TemplateRetriverFunc(func() string {
		return htmlTemplate
	})
}
