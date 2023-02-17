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

	SetLevel()
	SetError(err error) *Error
	SetCode(code string)
	SetDesc(desc string)
}

type Error struct {
	level int    // уровень ошибки (warn, error, fatal etc)
	code  string // код ошибки
	desc  string // напр, эта ошибка вызвана отказоустойчивостью бд
	err   error  // текст ошибки
	deep  *Error // ссылка на структуру вложенной ошибки
}

// TODO: здесь будет формироваться красивенький вывод
func (e Error) Error() string {
	output := strings.Builder{}

	for e.deep != nil {

		output.WriteString(fmt.Sprintf("level: %d, code: %s, description: %s, error: %v;\n",
			e.level, e.code, e.desc, e.err))

		e = *e.deep
	}

	output.WriteString(fmt.Sprintf("level: %d, code: %s, description: %s, error: %v",
		e.level, e.code, e.desc, e.err))

	return output.String()
}

func NewError(level int, code, desc string) *Error {
	return &Error{
		level: level,
		code:  code,
		desc:  desc,
	}
}

// NextError создаёт новую ошибку, наследуя переданную deep
// level, code, desc, err задаются отдельно функциями SetXXX
func (e *Error) NextError(deep *Error) *Error {
	e.deep = deep
	return e
}

func (e *Error) SetError(err error) *Error {
	e.err = err
	return e
}

func (e *Error) SetLevel(level int) {
	e.level = level
}

func (e *Error) SetCode(code string) {
	e.code = code
}

func (e *Error) SetDesc(desc string) {
	e.desc = desc
}

func (e *Error) Marshal() {}

func (e *Error) UnMarshal() {}
