package object

import (
	"hek/ast"
	"hek/code"
)

type Fun struct {
	Params []*ast.Identifier
	Block  *ast.BlockStatement
	Env    *Env
}

func (f *Fun) Type() ObjectType {
	return FUN
}

func (f Fun) Inspect() string {
	return "fun"
}

type CompliedFun struct {
	Instructions code.Instructions
	Free         []Object
	NumLocal     int
}

func (c *CompliedFun) Type() ObjectType {
	return CompiledFun
}

func (c *CompliedFun) Inspect() string {
	return "Complied fun"
}
