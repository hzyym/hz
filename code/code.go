package code

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
	OpPop
	OpSub
	OpMul
	OpDiv
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGT
	OpLT
	OpMinus
	OpBang
	OpJumpNotTrueThy
	OpJump
	OpNull
	OpSetGlobal
	OpGetGlobal
	OpArray
	OpIndex
	OpCall
	OpReturnValue
	OpReturn
	OpSetLocal
	OpGetLocal
	OpInternalFun
	OpGetFree
	OpLoadFun
	OpTwoAdd
	OpTwoSub
	OpDelGlobal
	OpDelLocal
	OpSetIndexGlobal
	OpSetIndexLocal
)

type Definitions struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definitions{
	OpConstant:       {"opConstant", []int{2}},
	OpAdd:            {"opAdd", []int{}},
	OpPop:            {"opPop", []int{}},
	OpSub:            {"opSub", []int{}},
	OpDiv:            {"opDiv", []int{}},
	OpMul:            {"opMul", []int{}},
	OpTrue:           {"opTrue", []int{}},
	OpFalse:          {"opFalse", []int{}},
	OpEqual:          {"opEqual", []int{}},
	OpNotEqual:       {"opNotEqual", []int{}},
	OpLT:             {"opLT", []int{}},
	OpGT:             {"opGT", []int{}},
	OpBang:           {"opBang", []int{}},
	OpMinus:          {"opBang", []int{}},
	OpJumpNotTrueThy: {"opJumpNotTrueThy", []int{2}},
	OpJump:           {"opJump", []int{2}},
	OpNull:           {"opNull", []int{}},
	OpGetGlobal:      {"opGetGlobal", []int{2}},
	OpSetGlobal:      {"opSetGlobal", []int{2}},
	OpArray:          {"opArray", []int{2}},
	OpIndex:          {"opIndex", []int{}},
	OpCall:           {"opCall", []int{2}},
	OpReturn:         {"opReturn", []int{}},
	OpReturnValue:    {"opReturnValue", []int{}},
	OpSetLocal:       {"opSetLocal", []int{2}},
	OpGetLocal:       {"opGetLocal", []int{2}},
	OpInternalFun:    {"opInternalFun", []int{2}},
	OpGetFree:        {"opGetFree", []int{2}},
	OpLoadFun:        {"opLoadFun", []int{2, 2}},
	OpTwoAdd:         {"opTwoAdd", []int{}},
	OpTwoSub:         {"opTwoSub", []int{}},
	OpDelGlobal:      {"opDelGlobal", []int{2}},
	OpDelLocal:       {"opDelLocal", []int{2}},
	OpSetIndexGlobal: {"opSetIndexGlobal", []int{2}},
	OpSetIndexLocal:  {"opSetIndexLocal", []int{2}},
}

func Lookup(op byte) (*Definitions, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, errors.New(fmt.Sprintf("opcode %d undefined", op))
	}
	return def, nil
}
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1

	for _, width := range def.OperandWidths {
		instructionLen += width
	}
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1

	for i, operand := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		}
		offset += width
	}

	return instruction
}

func ReadOperands(def *Definitions, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
func (in Instructions) String() string {
	var out bytes.Buffer
	i := 0

	for i < len(in) {
		def, err := Lookup(in[i])
		if err != nil {
			_, _ = fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}
		operands, read := ReadOperands(def, in[i+1:])

		_, _ = fmt.Fprintf(&out, "%04d %s\n", i, in.fmtInstruction(def, operands))

		i += 1 + read
	}
	return out.String()
}
func (in Instructions) fmtInstruction(def *Definitions, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d \n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}
