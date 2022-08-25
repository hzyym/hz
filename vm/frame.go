package vm

import (
	"hek/code"
	"hek/object"
)

type Frame struct {
	fn    *object.CompliedFun
	ip    int
	local []object.Object
}

func NewFrame(fu *object.CompliedFun) *Frame {
	return &Frame{
		fn: fu,
		ip: -1,
	}
}
func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
func (f *Frame) Push(obj object.Object, index ...int) {
	if len(index) > 0 && len(f.local)-1 >= len(index) {
		f.local[index[0]] = obj
		return
	}
	f.local = append(f.local, obj)
}
func (f *Frame) Pop(index int) object.Object {
	return f.local[index]
}
func (f *Frame) PushFree(obj object.Object) {
	f.fn.Free = append(f.fn.Free, obj)
}
func (f *Frame) PopFree(index int) object.Object {
	return f.fn.Free[index]
}
