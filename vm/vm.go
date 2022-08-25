package vm

import (
	"errors"
	"fmt"
	"hek/code"
	"hek/compiler"
	"hek/object"
	"strings"
)

var True = &object.Bool{Value: true}
var False = &object.Bool{Value: false}
var Null = &object.Null{}

const StackSiz = 2048
const GlobalSiz = 2048

type VM struct {
	constants []object.Object
	global    []object.Object
	stack     []object.Object
	sp        int

	err    []string
	errLen int

	frame      []*Frame
	frameIndex int
}

func NewVM(byteCode *compiler.Bytecode) *VM {
	main_ := NewFrame(&object.CompliedFun{Instructions: byteCode.Instructions})
	return &VM{
		constants:  byteCode.Constants,
		sp:         0,
		frame:      []*Frame{main_},
		frameIndex: 1,
		stack:      make([]object.Object, StackSiz),
		global:     make([]object.Object, GlobalSiz),
	}
}
func NewVMCache(byteCode *compiler.Bytecode, global []object.Object) *VM {
	main_ := NewFrame(&object.CompliedFun{Instructions: byteCode.Instructions})
	return &VM{
		constants:  byteCode.Constants,
		frame:      []*Frame{main_},
		frameIndex: 1,
		global:     global,
		stack:      make([]object.Object, StackSiz),
	}
}
func (v *VM) LastPoppedStackElem() object.Object {
	return v.stack[v.sp]
}
func (v *VM) Run() error {
	var index int
	for v.currentFrame().ip < len(v.currentFrame().Instructions())-1 {
		v.currentFrame().ip++
		index++
		frame := v.currentFrame()

		op := code.Opcode(frame.Instructions()[frame.ip])

		switch op {
		case code.OpConstant:
			constOIndex := v.getUint()
			v.push(v.constants[constOIndex])
		case code.OpAdd, code.OpSub, code.OpDiv, code.OpMul:
			v.operation(op)
		case code.OpPop:
			v.pop()
		case code.OpTrue, code.OpFalse:
			v.bool(op)
		case code.OpGT, code.OpLT, code.OpEqual, code.OpNotEqual:
			v.compare(op)
		case code.OpBang, code.OpMinus, code.OpTwoSub, code.OpTwoAdd:
			v.prefix(op)
		case code.OpJumpNotTrueThy:
			jumpIndex := v.getUint()
			if !v.IF(v.pop()) {
				v.currentFrame().ip = int(jumpIndex) - 1
			}
		case code.OpJump:
			jumpIndex := v.getUint()
			v.currentFrame().ip = int(jumpIndex) - 1
		case code.OpNull:
			v.push(Null)
		case code.OpSetGlobal:
			index := v.getUint()
			val := v.pop()
			v.global[index] = val
		case code.OpGetGlobal:
			index := v.getUint()
			val := v.global[index]
			if val == nil {
				v.push(Null)
			} else {
				v.push(val)
			}
		case code.OpArray:
			v.array()
		case code.OpIndex:
			v.index()
		case code.OpCall:
			v.call()
		case code.OpReturnValue, code.OpReturn:
			v.returnValue(op)
		case code.OpSetLocal:
			index := int(v.getUint())
			frame.Push(v.pop(), index)
		case code.OpGetLocal:
			index := int(v.getUint())
			v.push(frame.Pop(index))
		case code.OpInternalFun:
			v.internalFun()
		case code.OpLoadFun:
			v.loadFun()
		case code.OpGetFree:
			index := int(v.getUint())
			v.push(frame.PopFree(index))
		case code.OpSetIndexGlobal, code.OpSetIndexLocal:
			v.indexSet(op)
		default:
			v.errors("VM Op err")
		}

		if v.errLen > 0 {
			return errors.New(v.echoError())
		}
	}
	return nil
}
func (v *VM) push(object_ object.Object) {
	if v.sp >= StackSiz {
		v.errors("stack overflow")
		return
	}
	v.stack[v.sp] = nil
	v.stack[v.sp] = object_
	v.sp++
}
func (v *VM) pop() object.Object {
	if v.sp == 0 {
		v.errors("栈区已经无数据...")
		return Null
	}
	obj := v.stack[v.sp-1]
	v.sp--
	return obj
}
func (v *VM) operation(op code.Opcode) {
	right := v.pop()
	left := v.pop()
	var obj object.Object
	if left.Type() == object.STRING && right.Type() == object.STRING {
		obj = v.operationString(op, left.(*object.String), right.(*object.String))
	} else if left.Type() == object.INT && right.Type() == object.INT {
		obj = v.operationINT(op, left.(*object.Integer), right.(*object.Integer))
	} else {
		v.errors(fmt.Sprintf("操作类型不一致 %s - %s", left.Type().String(), right.Type().String()))
	}
	v.push(obj)
}
func (v *VM) operationString(op code.Opcode, left, right *object.String) object.Object {
	if op != code.OpAdd {
		v.errors("字符串类型只支持 '+' 的操作方式")
		return Null
	}
	return &object.String{Value: fmt.Sprintf("%s%s", left.Value, right.Value)}
}
func (v *VM) operationINT(op code.Opcode, left, right *object.Integer) object.Object {
	var obj object.Object
	switch op {
	case code.OpAdd:
		obj = &object.Integer{Value: left.Value + right.Value}
	case code.OpDiv:
		obj = &object.Integer{Value: left.Value / right.Value}
	case code.OpMul:
		obj = &object.Integer{Value: left.Value * right.Value}
	case code.OpSub:
		obj = &object.Integer{Value: left.Value - right.Value}
	}
	return obj
}
func (v *VM) errors(msg string) {
	v.err = append(v.err, msg)
	v.errLen++
}
func (v *VM) echoError() string {
	str := strings.Join(v.err, "\n")
	v.err = []string{}
	v.errLen = 0
	return str
}
func (v *VM) bool(op code.Opcode) {
	if op == code.OpTrue {
		v.push(True)
		return
	}
	v.push(False)
}
func (v *VM) compare(op code.Opcode) {
	var right object.Object
	var left object.Object
	right = v.pop()
	left = v.pop()
	if op == code.OpLT || op == code.OpGT {
		if !v.compareLGCheck(left, right) {
			v.errors("< 和 > 运算 必须是数字类型")
			return
		}
	}
	var obj object.Object
	switch op {
	case code.OpEqual:
		obj = v.compareBool(left.Inspect() == right.Inspect())
	case code.OpNotEqual:
		obj = v.compareBool(left.Inspect() != right.Inspect())
	case code.OpGT:
		obj = v.compareBool(left.(*object.Integer).Value > right.(*object.Integer).Value)
	case code.OpLT:
		obj = v.compareBool(left.(*object.Integer).Value < right.(*object.Integer).Value)
	}
	v.push(obj)
}
func (v *VM) compareBool(is bool) object.Object {
	if is {
		return True
	}
	return False
}
func (v *VM) compareLGCheck(left object.Object, right object.Object) bool {
	if left.Type() != object.INT || right.Type() != object.INT {
		return false
	}
	return true
}
func (v *VM) prefix(op code.Opcode) {
	val := v.pop()
	var obj object.Object

	switch op {
	case code.OpBang:
		obj = v.prefixBang(val)
	case code.OpMinus:
		obj = &object.Integer{Value: -val.(*object.Integer).Value}
	case code.OpTwoSub:
		obj = &object.Integer{Value: val.(*object.Integer).Value - 1}
	case code.OpTwoAdd:
		obj = &object.Integer{Value: val.(*object.Integer).Value + 1}
	}

	v.push(obj)
}
func (v *VM) prefixBang(obj object.Object) object.Object {
	switch obj {
	case True:
		return False
	case False:
		return True
	case Null:
		return True
	default:
		return False
	}
}
func (v *VM) IF(object_ object.Object) bool {
	if object_ == True {
		return true
	}
	return false
}
func (v *VM) getUint() uint16 {
	frame := v.currentFrame()
	val := code.ReadUint16(frame.Instructions()[frame.ip+1:])
	frame.ip += 2
	return val
}
func (v *VM) array() {
	num := int(v.getUint())
	obj := &object.Array{Value: make([]object.Object, num)}
	for i := 0; i < num; i++ {
		obj.Value[num-i-1] = v.pop()
	}
	v.push(obj)
}
func (v *VM) index() {
	index := v.pop()
	value := v.pop()
	if value.Type() != object.ARRAY {
		v.errors("不能操作数组一样操作普通变量")
		return
	}
	arr := value.(*object.Array)
	if index.Type() != object.INT {
		v.errors("数组的索引只能int类型")
		return
	}
	if len(arr.Value) == 0 {
		v.push(Null)
		return
	}
	i := index.(*object.Integer).Value
	if int(i) >= len(arr.Value) || i < 0 {
		v.push(Null)
		return
	}
	v.push(arr.Value[i])
}
func (v *VM) currentFrame() *Frame {
	return v.frame[v.frameIndex-1]
}
func (v *VM) pushFrame(f *Frame) {
	v.frame = append(v.frame, f)
	v.frameIndex++
}
func (v *VM) popFrame() *Frame {
	v.frameIndex--
	tmp := v.frame[v.frameIndex]
	v.frame = v.frame[:v.frameIndex]
	return tmp
}
func (v *VM) call() {
	val := int(v.getUint())
	fun := v.pop()
	if fun.Type() == object.BUILTFun {
		f := fun.(*object.InternalFun)
		var prams []object.Object
		if val == 1 {
			prams = v.stack[v.sp-val : v.sp]
		} else {
			prams = v.stack[v.sp-val : v.sp]
		}

		v.sp -= val
		obj := f.Fun_(prams...)
		if obj != nil {
			v.push(obj)
		} else {
			v.push(Null)
		}
		return
	}
	if fun.Type() != object.CompiledFun {
		v.errors("调用的不是一个方法")
		return
	}
	f := fun.(*object.CompliedFun)
	v.pushFrame(NewFrame(f))
}
func (v *VM) returnValue(op code.Opcode) {
	var obj object.Object
	if op == code.OpReturn {
		obj = Null
	} else {
		obj = v.pop()
	}
	v.popFrame()
	v.push(obj)
}
func (v *VM) internalFun() {
	fun := object.GetFun(int(v.getUint()))
	v.push(fun)
}
func (v *VM) loadFun() {
	constantsIndex := v.getUint()
	freeNum := int(v.getUint())
	fun := v.constants[constantsIndex].(*object.CompliedFun)

	if freeNum > 0 {
		params := make([]object.Object, freeNum)
		copy(params, v.stack[v.sp-freeNum:v.sp])
		fun = &object.CompliedFun{Instructions: fun.Instructions, Free: params}
		v.sp -= freeNum
	}
	v.push(fun)
}
func (v *VM) indexSet(op code.Opcode) {
	//区分
	//全局
	//局部
	index := int(v.getUint())
	arrIndex := v.pop()
	value := v.pop()
	if op == code.OpSetIndexGlobal {
		//全局状态设置
		val := v.global[index]
		arr := val.(*object.Array)

		if arrIndex.Type() != object.INT {
			v.errors("数组索引只能是数字类型")
			return
		}
		arr.Value[arrIndex.(*object.Integer).Value] = value
		v.constants[index] = arr
	} else {
		val := v.currentFrame().local[index]
		arr := val.(*object.Array)
		if arrIndex.Type() != object.INT {
			v.errors("数组索引只能是数字类型")
			return
		}
		arr.Value[arrIndex.(*object.Integer).Value] = value
		v.currentFrame().local[index] = arr
	}
}
