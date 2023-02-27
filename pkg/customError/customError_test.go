package customError

import (
	"errors"
	"testing"
)

type ErrorList struct {
	level int // уровень ошибки (warn, error, fatal etc)
	deep  *ErrorList
}

type ErrorErr struct {
	level int // уровень ошибки (warn, error, fatal etc)
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
	// была откуда-то получена такая ошибка,
	// нужно обернуть её в новую
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
		// fmt.Println(e.level, e.err)
		e = e.err.(ErrorErr)
	}
	// fmt.Println(e)
	return e
}

// func TestNewErrorErr(t *testing.T) {
// 	t.Run("test", func(t *testing.T) {
// 		e := NewErrorErr(3)
// 		// fmt.Println(e)

// 		fmt.Println(GetErrorErr(e))
// 	})
// }

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
