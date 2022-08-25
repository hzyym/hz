package parser

import (
	"fmt"
	"hek/lexer"
	"testing"
)

func TestNewParser(t *testing.T) {
	input := "arr[0] = 1"

	l := lexer.NewLexer(input)

	p := NewParser(l)

	program := p.ParseProgram()

	fmt.Println(program.String())
	//if len(program.Statements) < 3 {
	//	t.Error("program len < 3")
	//	return
	//}
	for _, stmt := range program.Statements {
		if stmt != nil {
			fmt.Println(stmt.String())
		}

	}
	for _, err := range p.Errors() {
		fmt.Println(err)
	}
}
