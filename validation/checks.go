package validation

func Eq(a interface{}) Validator {
	return ValidatorFunc(func(v Valuer, r Requester) error {
		if a != v.Value() {
			return ErrEq.Fmt([]interface{}{v.Name(), a})
		}
		return nil
	})
}

