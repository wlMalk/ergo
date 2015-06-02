package validation_test

import (
	"github.com/wlMalk/ergo"
	"github.com/wlMalk/ergo/validation"

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
	eq := validation.Eq("ergo")
	pv := ergo.NewParamValue("param", "ergo", "query")
	req := ergo.NewRequest(nil)
	err := eq.Validate(pv, req)
	expect(t, nil, err)
	pv = ergo.NewParamValue("param", "rgo", "query")
	err = eq.Validate(pv, req)
	expect(t, validation.ErrEq.Fmt("param", "ergo").Error(), err.Error())
}

func TestRegexp(t *testing.T) {
	re := validation.Regexp(regexp.MustCompile("[ergo]{4}"))
	pv := ergo.NewParamValue("param", "ergo", "query")
	req := ergo.NewRequest(nil)
	err := re.Validate(pv, req)
	expect(t, nil, err)
	pv = ergo.NewParamValue("param", "rgo", "query")
	err = re.Validate(pv, req)
	expectNot(t, validation.ErrRegexp.Fmt("param", "[ergo]{4}").Error(), err)
}
