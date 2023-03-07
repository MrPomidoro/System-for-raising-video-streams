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
	level int
	err   error
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

func BenchmarkErrorList10(b *testing.B) {
	e := NewErrorList(10)
	b.Run("ErrorList-10", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			GetDeepErrList(*e)
		}
	})
}
func BenchmarkErrorList100(b *testing.B) {
	e := NewErrorList(100)
	b.Run("ErrorList-100", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			GetDeepErrList(*e)
		}
	})
}
func BenchmarkErrorList1000(b *testing.B) {
	e := NewErrorList(1000)
	b.Run("ErrorList-1000", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			GetDeepErrList(*e)
		}
	})
}
func BenchmarkErrorList10000(b *testing.B) {
	e := NewErrorList(10000)
	b.Run("ErrorList-10000", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			GetDeepErrList(*e)
		}
	})
}

func NewErrorErr(nErr int) ErrorErr {
	e := ErrorErr{level: 0, err: errors.New("blablabla")}
	for i := 1; i <= nErr; i++ {
		e = ErrorErr{err: e, level: i}
	}
	return e
}

func GetErrorErr(e ErrorErr) error {
	for errors.As(e.err, &ErrorErr{}) {
		e = e.err.(ErrorErr)
	}
	return e
}

func BenchmarkErrorErr10(b *testing.B) {
	e := NewErrorErr(10)
	b.Run("ErrorErr-10", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			GetErrorErr(e)
		}
	})
}
func BenchmarkErrorErr100(b *testing.B) {
	e := NewErrorErr(100)
	b.Run("ErrorErr-100", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			GetErrorErr(e)
		}
	})
}
func BenchmarkErrorErr1000(b *testing.B) {
	e := NewErrorErr(1000)
	b.Run("ErrorErr-1000", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			GetErrorErr(e)
		}
	})
}
func BenchmarkErrorErr10000(b *testing.B) {
	e := NewErrorErr(10000)
	b.Run("ErrorErr-10000", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			GetErrorErr(e)
		}
	})
}
