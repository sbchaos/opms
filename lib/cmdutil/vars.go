package cmdutil

import (
	"errors"

	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/util"
)

func GetArgFromVar[T any](cfg *config.Config, cmd, proj, arg string) (T, error) {
	p1 := cfg.GetCurrentProfile()
	key1 := proj + ":" + arg
	v1, found := util.ConfigAs[T](p1.Variables, key1)
	if found {
		return v1, nil
	}

	key2 := cmd + ":" + arg
	v2, found := util.ConfigAs[T](p1.Variables, key2)
	if found {
		return v2, nil
	}

	v3, found := util.ConfigAs[T](p1.Variables, arg)
	if found {
		return v3, nil
	}

	return v3, errors.New("no value for " + arg + " found")
}
