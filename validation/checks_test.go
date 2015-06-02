package validation_test

import (
	"github.com/wlMalk/ergo"
	"github.com/wlMalk/ergo/validation"

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
	expectNot(t, nil, err)
}
