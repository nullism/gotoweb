package builder

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/nullism/gotoweb/config"
)

func (b *Builder) getFuncMap() map[string]any {
	return map[string]any{
		"tpl": b.tpl,
		"map": b.toMap,
		"sub": b.sub,
		"add": b.add,
	}
}

func (b *Builder) add(num, amount int) int {
	return num + amount
}

func (b *Builder) sub(num, amount int) int {
	return num - amount
}

func (b *Builder) tpl(name string, pairs ...any) (string, error) {
	m, err := b.toMap(pairs...)
	if err != nil {
		return "", err
	}
	path := filepath.Join(b.site.ThemeDir(), config.HelpersDir, name+config.TemplateExt)
	b.context.Args = m
	return b.Render(path, b.context)
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
