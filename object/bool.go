package object

import "fmt"

type Bool struct {
	Value bool
}

func (b *Bool) Type() ObjectType {
	return BOOL
}

func (b *Bool) Inspect() string {
	return fmt.Sprintf("%v", b.Value)
}
