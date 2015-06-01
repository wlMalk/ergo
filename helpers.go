package ergo

import (
	"strings"
)

func prepareArgsSlice(args []string, f func(s string) bool) []string {
	if len(args) == 0 {
		return args
	}
	var nArgs []string
	for _, a := range args {
		a = strings.ToLower(a)
		duplicate := false
		for _, b := range nArgs {
			if a == b {
				duplicate = true
			}
		}
		if !duplicate && f(a) {
			nArgs = append(nArgs, a)
		}
	}

	return nArgs
}

type schemer interface {
	GetSchemes() []string
	setSchemes([]string)
}

func schemes(s schemer, schemes []string) {
	schemes = prepareArgsSlice(schemes, func(scheme string) bool {
		if scheme == SCHEME_HTTP ||
			scheme == SCHEME_HTTPS {
			return true
		}
		return false
	})
	if len(schemes) > 0 {
		s.setSchemes(schemes)
	}
}

