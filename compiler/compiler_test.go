package compiler

import (
	"fmt"
	"hek/lexer"
	"hek/parser"
	"testing"
)

func TestNewCompile(t *testing.T) {
	input := "let arr = [\"arr1\",\"arr2\",\"arr3\",\"arr4\"];arr[0] = \"test\""

	lexer_ := lexer.NewLexer(input)

	p := parser.NewParser(lexer_)

	program := p.ParseProgram()

	compile := NewCompile()

	err := compile.Compile(program)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(compile.String())

}
