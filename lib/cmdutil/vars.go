package cmdutil

import (
	"errors"

	"github.com/sbchaos/opms/lib/config"
)

func ConfigAs[T any](mapping map[string]any, key string) (T, bool) {
	var zero T
	val, ok := mapping[key]
	if ok {
		s, ok := val.(T)
		if ok {
			return s, true
		}
	}
	return zero, false
}

func GetArgFromVar[T any](cfg *config.Config, cmd, proj, arg string) (T, error) {
	p1 := cfg.GetCurrentProfile()
	key1 := proj + ":" + arg
	v1, found := ConfigAs[T](p1.Variables, key1)
	if found {
		return v1, nil
	}

	key2 := cmd + ":" + arg
	v2, found := ConfigAs[T](p1.Variables, key2)
	if found {
		return v2, nil
	}

	v3, found := ConfigAs[T](p1.Variables, arg)
	if found {
		return v3, nil
	}

	return v3, errors.New("no value for " + arg + " found")
}
