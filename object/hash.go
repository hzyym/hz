package object

import (
	"bytes"
	"strings"
)

type Hash struct {
	Value map[Object]Object
}

func (h *Hash) Type() ObjectType {
	return HASH
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer
	var arr []string
	out.WriteString("{")
	for object, o := range h.Value {
		arr = append(arr, object.Inspect()+":"+o.Inspect())
	}
	out.WriteString(strings.Join(arr, ","))
	out.WriteString("]")
	return out.String()
}
