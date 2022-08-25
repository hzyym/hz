package object

import (
	"fmt"
	"hek/lexer"
	"hek/parser"
	"testing"
)

func TestEval(t *testing.T) {
	input := "let a = fun (x) { return x + 5;}(5)"

	l := lexer.NewLexer(input)

	p := parser.NewParser(l)

	program := p.ParseProgram()

	env := NewEnv(nil)
	result := Eval(program, env)

	fmt.Println(result)
}
