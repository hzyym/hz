package object

import (
	"fmt"
	"strings"
)

type Array struct {
	Value []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY
}

func (a *Array) Inspect() string {
	var str []string
	for _, object := range a.Value {
		str = append(str, object.Inspect())
	}
	return fmt.Sprintf("[%s]", strings.Join(str, ","))
}
