package customError

import (
	"errors"
	"reflect"
	"testing"
)

func TestDefineError(t *testing.T) {
	tests := []struct {
		name   string
		arg    int
		expect string
	}{
		{
			name:   "testWarn",
			arg:    0,
			expect: "warn",
		},
		{
			name:   "testError",
			arg:    1,
			expect: "error",
		},
		{
			name:   "testFatal",
			arg:    2,
			expect: "fatal",
		},
		{
			name:   "testUnknown",
			arg:    10,
			expect: "unknown level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewError(tt.arg, "", "")
			if err.defineLevel() != tt.expect {
				t.Errorf("want %s, got %s", err.defineLevel(), tt.expect)
			}
		})
	}
}

func TestSetError(t *testing.T) {
	err := NewError(0, "", "")
	e := errors.New("fail")
	expect := &Error{level: 0, desc: "", code: "", err: errors.New("fail")}
	if !reflect.DeepEqual(err.SetError(e), expect) {
		t.Errorf("want %v, got %v", err, expect)
	}
}

func TestError(t *testing.T) {
	err := &Error{level: 0, desc: "", code: "", err: errors.New("fail")}
	expect := "\n\tlevel: warn, code: , description: , error: fail"
	if err.Error() != expect {
		t.Errorf("want %s, got %s", err.Error(), expect)
	}
}

type ErrorList struct {
	level int
	deep  *ErrorList
}

type ErrorErr struct {
	err error
}

func (e ErrorErr) Error() string {
	return e.err.Error()
}

func (e ErrorList) AddNewErrList(i int) *ErrorList {
	newErr := ErrorList{level: i}
	newErr.deep = &e
	return &newErr
}

func NewErrorList(nDeep int) *ErrorList {
	e := &ErrorList{level: 0, deep: nil}
	for i := 1; i <= nDeep; i++ {
		e = e.AddNewErrList(i)
	}
	return e
}

func GetDeepErrList(e ErrorList) {
	for e.deep != nil {
		e = *e.deep
	}
}
