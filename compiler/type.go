package compiler

import (
	"hek/code"
	"hek/object"
)

type EmittedInstruction struct {
	Op  code.Opcode
	Pos int
}
type CompilationScope struct {
	instructions code.Instructions
	last         EmittedInstruction
	previous     EmittedInstruction
}
type Compiler struct {
	constants   []object.Object
	scopes      []*CompilationScope
	scopeIndex  int
	symbolTable *SymbolTable

	getEmitPos bool
	pos        int
}
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
