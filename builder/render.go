package builder

import (
	"strings"
	"text/template"
)

// Render renders a template with the given content.
func (b *Builder) Render(tplName string, funcMap map[string]any, bs []byte, content *RenderContext) (string, error) {
	tpl, err := template.New(tplName).Funcs(funcMap).Option("missingkey=error").Parse(string(bs))
	if err != nil {
		return "", err
	}
	sb := strings.Builder{}
	err = tpl.Execute(&sb, content)
	return sb.String(), err
}

// RenderFile renders a template from a file with the given content.
func (b *Builder) RenderFile(tplPath string, funcMap map[string]any, content *RenderContext) (string, error) {
	bs, err := b.files.ReadFile(tplPath)
	if err != nil {
		return "", err
	}

	return b.Render(tplPath, funcMap, bs, content)
}

// RenderSource renders a source template with the given content.
func (b *Builder) RenderSource(tplPath string, content *RenderContext) (string, error) {
	return b.RenderFile(tplPath, b.getSourceFuncMap(), content)
}

// RenderTheme renders a template with the given content.
func (b *Builder) RenderTheme(tplPath string, content *RenderContext) (string, error) {
	return b.RenderFile(tplPath, b.getThemeFuncMap(), content)
}
