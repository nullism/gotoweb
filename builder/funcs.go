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
		"add":     b.add,
		"ftime":   b.ftime,
		"haskey":  b.mapHasKey,
		"href":    b.href,
		"tpl":     b.tplTheme,
		"map":     b.toMap,
		"mapset":  b.mapSet,
		"sub":     b.sub,
		"list":    b.list,
		"listadd": b.listAdd,
	}
}

func (b *Builder) getSourceFuncMap() map[string]any {
	return map[string]any{
		"add":     b.add,
		"ftime":   b.ftime,
		"href":    b.href,
		"sub":     b.sub,
		"map":     b.toMap,
		"mapset":  b.mapSet,
		"plink":   b.postLink,
		"list":    b.list,
		"listadd": b.listAdd,
		"tpl":     b.tplSource,
	}
}

func (b *Builder) add(num, amount int) int {
	return num + amount
}

func (b *Builder) ftime(t time.Time, fs ...string) string {
	for _, f := range fs {
		return t.Format(f)
	}
	return t.Format(b.site.TimeFormat)
}

func (b *Builder) mapHasKey(m map[string]any, key string) bool {
	_, ok := m[key]
	return ok
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

func (b *Builder) listAdd(items []any, item any) []any {
	return append(items, item)
}

func (b *Builder) mapSet(m map[string]any, key string, value any) map[string]any {
	m[key] = value
	return m
}

func (b *Builder) postLink(name string) (string, error) {
	path := b.files.Join(b.site.SourceDir, name+".md")
	if !b.files.Exists(path) {
		return "", fmt.Errorf("post %s does not exist", path)
	}
	pconf, _, err := b.postConfigFromPath(path)
	if err != nil {
		return "", err
	}
	title := pconf.Title
	if title == "" {
		title = name
	}

	href, err := b.prefix(name + ".html")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`<a href="%v">%v</a>`, href, title), nil
}

func (b *Builder) prefix(href string) (string, error) {
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
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

func (b *Builder) tplSource(name string, pairs ...any) (string, error) {
	m, err := b.toMap(pairs...)
	if err != nil {
		return "", err
	}
	path := b.files.Join(b.site.SourceDir, config.HelpersDir, name)
	b.context.Args = m
	return b.RenderSource(path, b.context)
}

func (b *Builder) tplTheme(name string, pairs ...any) (string, error) {
	m, err := b.toMap(pairs...)
	if err != nil {
		return "", err
	}
	path := b.files.Join(b.site.ThemeDir(), config.HelpersDir, name)
	b.context.Args = m
	return b.RenderTheme(path, b.context)
}

func (b *Builder) toMap(pairs ...any) (map[string]any, error) {

	if len(pairs)%2 != 0 {
		return nil, errors.New("misaligned map")
	}

	m := make(map[string]any)

	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)

		if !ok {
			return nil, fmt.Errorf("cannot use type %T as map key", pairs[i])
		}
		m[key] = pairs[i+1]
	}
	return m, nil
}
