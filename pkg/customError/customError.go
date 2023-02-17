package customError

const (
	WarnLevel = iota
	ErrorLevel
	FatalLevel
)

type IError interface {
	SetLevel()
}

type Error struct {
	level int    // уровень ошибки (warn, error, fatal etc)
	code  string // код ошибки
	err   error
	desc  string // напр, эта ошибка вызвана отказоустойчивостью бд
}

// TODO: здесь будет формироваться красивенький вывод
func (e *Error) Error() string {
	return e.err.Error()
}

func NewError(level int, code, desc string) *Error {
	return &Error{
		level: level,
		code:  code,
		desc:  desc,
	}
}

func (e *Error) SetLevel(level int) *Error {
	e.level = level
	return e
}

func (e *Error) SetError(err error) *Error {
	e.err = err
	return e
}

func (e *Error) SetCode(code string) *Error {
	e.code = code
	return e
}

func (e *Error) SetDesc(desc string) *Error {
	e.desc = desc
	return e
}

func (e *Error) Marshal() {}

func (e *Error) UnMarshal() {}
