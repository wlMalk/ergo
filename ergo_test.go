package ergo

import (
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %+v to equal %+v", b, a)
	}
}

func expectNot(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Expected %+v to not equal %+v", a, b)
	}
}

func TestNewRoute(t *testing.T) {
	e := New()
	v1 := e.New("/v1//")
	usersRoute := v1.New("users")
	postRoute := v1.New("posts/{id}")
	expect(t, "/v1/users", usersRoute.GetFullPath())
	expect(t, "/v1/posts/{id}", postRoute.GetFullPath())
}
