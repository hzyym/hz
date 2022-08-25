package object

type Error struct {
	Msg string
}

func (e *Error) Type() ObjectType {
	return ERROR
}

func (e *Error) Inspect() string {
	return e.Msg
}
