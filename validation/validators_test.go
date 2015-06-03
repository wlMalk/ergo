package validation_test

import (
	"github.com/wlMalk/ergo"
	. "github.com/wlMalk/ergo/validation"

	"regexp"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %+v to equal %+v", a, b)
	}
}

func expectNot(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Expected %+v to not equal %+v", a, b)
	}
}

func TestEq(t *testing.T) {
	eq := Eq("ergo")
	pv := ergo.NewParamValue("param", "ergo", "query")
	req := ergo.NewRequest(nil)
	err := eq.Validate(pv, req)
	expect(t, nil, err)
	pv = ergo.NewParamValue("param", "rgo", "query")
	err = eq.Validate(pv, req)
	expect(t, ErrEq.Fmt("param", "ergo").Error(), err.Error())
}

func TestRegexp(t *testing.T) {
	re := Regexp(regexp.MustCompile("[ergo]{4}"))
	pv := ergo.NewParamValue("param", "ergo", "query")
	req := ergo.NewRequest(nil)
	err := re.Validate(pv, req)
	expect(t, nil, err)
	pv = ergo.NewParamValue("param", "rgo", "query")
	err = re.Validate(pv, req)
	expect(t, ErrRegexp.Fmt("param", "[ergo]{4}").Error(), err.Error())
}

func TestIf(t *testing.T) {
	i := If(func(Valuer, Requester) bool {
		return false
	}, Regexp(regexp.MustCompile("[ergo]{4}")))
	pv := ergo.NewParamValue("param", "rgo", "query")
	req := ergo.NewRequest(nil)
	err := i.Validate(pv, req)
	expect(t, nil, err)
	i = If(func(Valuer, Requester) bool {
		return true
	}, Regexp(regexp.MustCompile("[ergo]{4}")))
	err = i.Validate(pv, req)
	expect(t, ErrRegexp.Fmt("param", "[ergo]{4}").Error(), err.Error())
}

func TestIfElse(t *testing.T) {
	i := IfElse(func(Valuer, Requester) bool {
		return true
	}, []Validator{Regexp(regexp.MustCompile("[ergo]{4}"))},
		[]Validator{Eq("ergo")})
	pv := ergo.NewParamValue("param", "rgo", "query")
	req := ergo.NewRequest(nil)
	err := i.Validate(pv, req)
	expect(t, ErrRegexp.Fmt("param", "[ergo]{4}").Error(), err.Error())
	i = IfElse(func(Valuer, Requester) bool {
		return false
	}, []Validator{Regexp(regexp.MustCompile("[ergo]{4}"))},
		[]Validator{Eq("ergo")})
	err = i.Validate(pv, req)
	expect(t, ErrEq.Fmt("param", "ergo").Error(), err.Error())
}
