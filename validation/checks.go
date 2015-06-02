package validation

import (
	"regexp"
)

func Eq(a interface{}) Validator {
	return ValidatorFunc(func(v Valuer, r Requester) error {
		if a != v.Value() {
			return ErrEq.Fmt([]interface{}{v.Name(), a}...)
		}
		return nil
	})
}

func Regexp(p *regexp.Regexp) Validator {
	return ValidatorFunc(func(v Valuer, r Requester) error {
		if !p.MatchString(v.String()) {
			return ErrRegexp.Fmt([]interface{}{v.Name(), p.String()}...)
		}
		return nil
	})
}

