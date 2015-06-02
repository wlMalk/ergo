package validation

import (
	"fmt"
)

type ErrFmt string

func (err ErrFmt) Fmt(vals ...interface{}) error {
	return ValidationErr(fmt.Sprintf(string(err), vals...))
}

type ValidationErr string

func (err ValidationErr) Error() string {
	return string(err)
}

const (
	ErrEq     = ErrFmt("Validation Error: %v should equal %v")
	ErrRegexp = ErrFmt("Validation Error: %v should match %v")
)
