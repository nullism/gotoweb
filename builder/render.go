package builder

import (
	"os"
	"strings"
	"text/template"
)

// Render renders a template with the given content.
func (b *Builder) Render(tplPath string, content *RenderContext) (string, error) {
	bs, err := os.ReadFile(tplPath)
	if err != nil {
		return "", err
	}

	tpl, err := template.New(tplPath).Funcs(b.getFuncMap()).Option("missingkey=error").Parse(string(bs))

	if err != nil {
		return "", err
	}

	sb := strings.Builder{}
	err = tpl.Execute(&sb, content)
	return sb.String(), err
}
