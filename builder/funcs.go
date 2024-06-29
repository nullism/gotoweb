package builder

import (
	"errors"
	"fmt"
)

func (b *Builder) getFuncMap() map[string]any {
	return map[string]any{
		"include": b.include,
		"map":     b.toMap,
	}
}

func (b *Builder) include(name string, pairs ...any) (string, error) {
	m, err := b.toMap(pairs...)
	if err != nil {
		return "", err
	}

	return b.Render(name, m)
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
