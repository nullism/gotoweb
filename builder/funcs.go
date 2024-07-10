package builder

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/nullism/gotoweb/config"
)

func (b *Builder) getThemeFuncMap() map[string]any {
	return map[string]any{
		"tpl":   b.tplTheme,
		"map":   b.toMap,
		"sub":   b.sub,
		"add":   b.add,
		"href":  b.href,
		"list":  b.list,
		"ftime": b.time,
	}
}

func (b *Builder) getSourceFuncMap() map[string]any {
	return map[string]any{
		"add":   b.add,
		"sub":   b.sub,
		"map":   b.toMap,
		"plink": b.postLink,
		"href":  b.href,
		"list":  b.list,
	}
}

func (b *Builder) add(num, amount int) int {
	return num + amount
}

func (b *Builder) href(href string, params ...any) (string, error) {
	var uri string
	var err error
	uri, err = b.prefix(href)
	if err != nil {
		return "", err
	}

	if len(params) > 0 {
		uri = uri + "?"
		qparams, err := b.toMap(params...)
		if err != nil {
			return "", err
		}
		for k, v := range qparams {
			uri = uri + fmt.Sprintf("%v=%v&", k, v)
		}
		uri = strings.TrimSuffix(uri, "&")
	}
	return uri, err
}

func (b *Builder) list(items ...any) ([]any, error) {
	return items, nil
}

func (b *Builder) postLink(name string) (string, error) {
	path := b.files.Join(b.site.SourceDir, name+".md")
	if !b.files.Exists(path) {
		return "", fmt.Errorf("post %s does not exist", path)
	}
	p, err := b.postFromSource(path)
	if err != nil {
		return "", err
	}
	title := p.Title
	href, err := b.prefix(name + ".html")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`<a href="%v">%v</a>`, href, title), nil
}

func (b *Builder) prefix(href string) (string, error) {
	if strings.HasPrefix(href, "http") {
		return href, nil
	}
	if !strings.HasPrefix(href, "/") {
		href = "/" + href
	}

	if b.site.Prefix == "" {
		return href, nil
	}

	return url.JoinPath(b.site.Prefix, href)
}

func (b *Builder) sub(num, amount int) int {
	return num - amount
}

func (b *Builder) time(t time.Time) string {
	return t.Format(b.site.TimeFormat)
}

func (b *Builder) tplTheme(name string, pairs ...any) (string, error) {
	m, err := b.toMap(pairs...)
	if err != nil {
		return "", err
	}
	path := b.files.Join(b.site.ThemeDir(), config.HelpersDir, name+config.TemplateExt)
	b.context.Args = m
	return b.RenderTheme(path, b.context)
}

func (b *Builder) toMap(pairs ...any) (map[string]any, error) {

	if len(pairs)%2 != 0 {
		return nil, errors.New("misaligned map")
	}

	m := make(map[string]any, len(pairs)/2)

	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)

		if !ok {
			return nil, fmt.Errorf("cannot use type %T as map key", pairs[i])
		}
		m[key] = pairs[i+1]
	}
	return m, nil
}
