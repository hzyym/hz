package repl

import (
	"bufio"
	"fmt"
	"hek/compiler"
	"hek/lexer"
	"hek/object"
	"hek/parser"
	"hek/vm"
	"io"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	//env := object.NewEnv(nil)

	symbolTable := compiler.NewSymbolTable(nil)
	var consts []object.Object
	golbal := make([]object.Object, vm.GlobalSiz)
	for {
		fmt.Fprintf(out, PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.NewLexer(line)

		p := parser.NewParser(l)

		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				fmt.Println("err:", err)
			}
			continue
		}
		com := compiler.NewCompileCache(symbolTable, consts)

		err := com.Compile(program)
		if err != nil {
			fmt.Println("compile err:", err)
			continue
		}
		consts = com.ByteCode().Constants
		vm_ := vm.NewVMCache(com.ByteCode(), golbal)

		err = vm_.Run()
		if err != nil {
			fmt.Println("vm err:", err)
			continue
		}
	}
}
