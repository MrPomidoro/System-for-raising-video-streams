package customError

import (
	"fmt"
	"strings"
)

const (
	WarnLevel = iota
	ErrorLevel
	FatalLevel
)

type IError interface {
	Error() string
	SetError(err error) *Error
}

type Error struct {
	// уровень ошибки (warn, error, fatal etc)
	level int
	// код ошибки
	code string
	// описание ошибки, например, "ошибка вызвана отказоустойчивостью бд"
	desc string
	// текст ошибки
	err error
}

// NewError инициализирует новую ошибку
func NewError(level int, code, desc string) IError {
	return &Error{
		level: level,
		code:  code,
		desc:  desc,
	}
}

// SetError настривает новый текст поля err,
// возвращает структуру типа Error
func (e *Error) SetError(err error) *Error {
	e.err = err
	return e
}

func (e *Error) Error() string {
	output := strings.Builder{}
	output.WriteString("\n")

	output.WriteString(fmt.Sprintf("\tlevel: %s, code: %s, description: %s, error: %v",
		e.defineLevel(), e.code, e.desc, e.err))

	return output.String()
}

// defineLevel возвращает уровень ошибки в виде строки
func (e *Error) defineLevel() string {
	switch e.level {
	case 0:
		return "warn"
	case 1:
		return "error"
	case 2:
		return "fatal"
	default:
		return "unknown level"
	}
}

//// NextError создаёт новую ошибку, наследуя переданную deep
//// level, code, desc, err задаются отдельно функциями SetXXX
//func (e *Error) NextError(deep IError) {
//	e.deep = deep.getError()
//}
//
//func (e *Error) getError() *Error {
//	return e
//}

//// SetLevel настривает новый текст поля level
//func (e *Error) SetLevel(level int) {
//	e.level = level
//}
//
//// SetCode настривает новый текст поля code
//func (e *Error) SetCode(code string) {
//	e.code = code
//}
//
//// SetDesc настривает новый текст поля desc
//func (e *Error) SetDesc(desc string) {
//	e.desc = desc
//}

// func (e *Error) Marshal() {}

// func (e *Error) UnMarshal() {}
