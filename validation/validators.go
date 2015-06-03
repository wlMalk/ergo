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

func If(f func(Valuer, Requester) bool, validators ...Validator) Validator {
	return ValidatorFunc(func(v Valuer, r Requester) error {
		if f(v, r) {
			for _, va := range validators {
				err := va.Validate(v, r)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func IfElse(f func(Valuer, Requester) bool, tValidators []Validator, fValidators []Validator) Validator {
	return ValidatorFunc(func(v Valuer, r Requester) error {
		if f(v, r) {
			for _, va := range tValidators {
				err := va.Validate(v, r)
				if err != nil {
					return err
				}
			}
		} else {
			for _, va := range fValidators {
				err := va.Validate(v, r)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
