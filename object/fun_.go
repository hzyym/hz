package object

import (
	"bytes"
	"fmt"
)

func Len(args ...Object) Object {
	if len(args) <= 0 {
		return &Error{Msg: "len param < 1"}
	}

	switch arg := args[0].(type) {
	case *String:
		return &Integer{Value: int64(len(arg.Value))}
	case *Array:
		return &Integer{Value: int64(len(arg.Value))}
	default:
		return nil
	}
}
func Put(args ...Object) Object {
	if len(args) < 2 {
		return nil
	}
	var arr *Array
	if v, ok := args[0].(*Array); !ok {
		return nil
	} else {
		arr = v
	}
	arr.Value = append(arr.Value, args[1])
	return arr
}
func Println(args ...Object) Object {
	var out bytes.Buffer
	for _, val := range args {
		out.WriteString(val.Inspect())
	}
	fmt.Println(out.String())
	return NULL_
}
func Echo(args ...Object) Object {
	if len(args) < 1 {
		return NULL_
	}
	fmt.Println(args[0].Inspect())
	return NULL_
}
func StringReversal(args ...Object) Object {
	if args[0].Type() != STRING {
		return nil
	}

	var out bytes.Buffer
	v := args[0].(*String)
	le := len(v.Value)
	for i := 0; i < le; i++ {
		out.WriteByte(v.Value[le-i-1])
	}

	return &String{Value: out.String()}
}
