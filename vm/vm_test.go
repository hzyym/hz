package vm

import (
	"fmt"
	"hek/compiler"
	"hek/lexer"
	"hek/parser"
	"os"
	"testing"
)

func TestVM(t *testing.T) {
	f, _ := os.Open("./test.txt")
	defer f.Close()
	fileInfo, _ := f.Stat()
	buf := make([]byte, fileInfo.Size())
	_, _ = f.Read(buf)
	input := string(buf)

	lexer_ := lexer.NewLexer(input)

	p := parser.NewParser(lexer_)

	program := p.ParseProgram()

	compile := compiler.NewCompile()

	err := compile.Compile(program)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(compile.String())

	vm_ := NewVM(compile.ByteCode())

	err = vm_.Run()
	if err != nil {
		fmt.Println("vm err:", err)
	}

	//fmt.Println(vm_.LastPoppedStackElem().Inspect())
}
