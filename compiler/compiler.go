package compiler

import "C"
import (
	"errors"
	"fmt"
	"hek/ast"
	"hek/code"
	"hek/object"
	"hek/token"
)

func NewCompile() *Compiler {
	c := &Compiler{symbolTable: NewSymbolTable(nil)}
	c.createScope()
	return c
}
func NewCompileCache(symbolTable *SymbolTable, constants []object.Object) *Compiler {
	c := &Compiler{symbolTable: symbolTable, constants: constants}
	c.createScope()
	return c
}
func (c *Compiler) Compile(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Program:
		for _, statement := range n.Statements {
			e := c.callBack(statement)
			if e != nil {
				return e
			}
		}
	case *ast.ExpressionStatement:
		err := c.callBack(n.Expression)
		if err != nil {
			return err
		}
		c.WherePop(n.Expression)
		return err
	case *ast.InfixExpression:
		return c.infixExpression(n)
	case *ast.IntegerLiteral:
		c.IntegerLiteral(n)
	case *ast.BoolExpression:
		if n.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.PrefixExpression:
		err := c.callBack(n.Right)
		if err != nil {
			return err
		}
		return c.prefixExpression(n.Token.Type)
	case *ast.IFExpression:
		return c.ifExpression(n)
	case *ast.BlockStatement:
		for _, statement := range n.Statements {
			e := c.callBack(statement)
			if e != nil {
				return e
			}
		}
	case *ast.LetStatement:
		symbol := c.symbolTable.SetSymbol(n.Name.Value)
		err := c.callBack(n.Value)
		if err != nil {
			return err
		}

		if symbol.types == Global {
			c.emit(code.OpSetGlobal, symbol.index)
		} else {
			c.emit(code.OpSetLocal, symbol.index)
		}
	case *ast.Identifier:
		symbol, ok := c.symbolTable.GetSymbol(n.Value)
		if !ok {
			pos, ok_ := object.GetNameIndex(n.Value)
			if ok_ {
				c.emit(code.OpInternalFun, pos)
				return nil
			}
			return errors.New(fmt.Sprintf("使用了未定义的变量 %s", n.Value))
		}
		c.symbolEmitGet(symbol)
	case *ast.StringExpression:
		obj := &object.String{Value: n.Value}
		c.emit(code.OpConstant, c.addConstant(obj))
	case *ast.ArrayExpression:
		for _, expression := range n.Value {
			err := c.callBack(expression)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpArray, len(n.Value))
	case *ast.IndexExpression:
		err := c.callBack(n.Left)
		if err != nil {
			return err
		}
		err = c.callBack(n.Index)
		if err != nil {
			return err
		}
		c.emit(code.OpIndex)
	case *ast.AssigExpression:
		return c.Assig(n)
	case *ast.FunExpression:
		return c.fun(n)
	case *ast.ReturnStatement:
		err := c.callBack(n.Value)
		if err != nil {
			return err
		}
		c.emit(code.OpReturnValue)
	case *ast.CallExpression:
		for _, param := range n.Params {
			err := c.callBack(param)
			if err != nil {
				return err
			}
		}
		err := c.callBack(n.Fun)
		if err != nil {
			return err
		}
		c.emit(code.OpCall, len(n.Params))
	case *ast.ForExpression:
		return c.forExpression(n)
	case *ast.SuffixExpression:
		symbol, ok := c.symbolTable.GetSymbol(n.Left.Value)
		if !ok {
			return errors.New(fmt.Sprintf("不能对没有定义的变量赋值 %s", n.Left.Value))
		}
		c.symbolEmitGet(symbol)
		err := c.infixOperator(n.Token.Type)
		c.symbolEmitSet(symbol)
		return err
	}

	return nil
}
func (c *Compiler) callBack(node ast.Node) error {
	err := c.Compile(node)
	if err != nil {
		return err
	}
	return nil
}
func (c *Compiler) infixExpression(infix *ast.InfixExpression) error {
	err := c.Compile(infix.Left)
	if err != nil {
		return err
	}
	err = c.Compile(infix.Right)
	if err != nil {
		return err
	}
	err = c.infixOperator(infix.Token.Type)
	return err
}
func (c *Compiler) IntegerLiteral(integer *ast.IntegerLiteral) {
	obj := &object.Integer{Value: integer.Value}
	c.emit(code.OpConstant, c.addConstant(obj))
}
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	c.setEmitted(op, pos)
	return pos
}
func (c *Compiler) setEmitted(op code.Opcode, pos int) {
	prev := c.scopes[c.scopeIndex].previous
	c.scopes[c.scopeIndex].last = EmittedInstruction{Op: op, Pos: pos}
	c.scopes[c.scopeIndex].previous = prev
}
func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}
	return c.scopes[c.scopeIndex].last.Op == op
}
func (c *Compiler) delLastPop() {
	last := c.scopes[c.scopeIndex].last
	previous := c.scopes[c.scopeIndex].previous
	old := c.currentInstructions()

	c.scopes[c.scopeIndex].instructions = old[:last.Pos]
	c.scopes[c.scopeIndex].last = previous
}
func (c *Compiler) infixOperator(tok token.Type) error {
	switch tok {
	case token.PLUS:
		c.emit(code.OpAdd)
	case token.SLASH:
		c.emit(code.OpDiv)
	case token.ASTERISK:
		c.emit(code.OpMul)
	case token.MINUS:
		c.emit(code.OpSub)
	case token.EQ:
		c.emit(code.OpEqual)
	case token.NotEq:
		c.emit(code.OpNotEqual)
	case token.GT:
		c.emit(code.OpGT)
	case token.LT:
		c.emit(code.OpLT)
	case token.TwoPlus:
		c.emit(code.OpTwoAdd)
	case token.TwoMinus:
		c.emit(code.OpTwoSub)
	default:
		return errors.New(fmt.Sprintf("unknown operator %s", tok.ToString()))
	}
	return nil
}
func (c *Compiler) addInstruction(int []byte) int {
	pos := len(c.currentInstructions())
	c.scopes[c.scopeIndex].instructions = append(c.currentInstructions(), int...)
	return pos
}
func (c *Compiler) addConstant(object_ object.Object) int {
	c.constants = append(c.constants, object_)
	return len(c.constants) - 1
}
func (c *Compiler) ByteCode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}
func (c *Compiler) prefixExpression(tok token.Type) error {
	switch tok {
	case token.BANG:
		c.emit(code.OpBang)
	case token.MINUS:
		c.emit(code.OpMinus)
	default:
		return errors.New(fmt.Sprintf("unknown operator %s", tok.ToString()))
	}
	return nil
}
func (c *Compiler) String() string {
	return c.currentInstructions().String()
}
func (c *Compiler) ifExpression(if_ *ast.IFExpression) error {
	err := c.callBack(if_.Condition) //生成条件部位
	if err != nil {
		return err
	}
	jumpNotPos := c.emit(code.OpJumpNotTrueThy, 999)
	err = c.callBack(if_.Consequence) //生成ture 语法
	if err != nil {
		return err
	}

	if c.lastInstructionIs(code.OpPop) {
		c.delLastPop()
	}
	c.changOperand(jumpNotPos, len(c.currentInstructions()))
	if if_.Alternative != nil {
		pos := c.emit(code.OpJump, 999)
		c.changOperand(jumpNotPos, len(c.currentInstructions()))

		err = c.callBack(if_.Alternative)
		if err != nil {
			return err
		}
		if c.lastInstructionIs(code.OpPop) {
			c.delLastPop()
		}
		c.changOperand(pos, len(c.currentInstructions()))
	} else {
		//平栈
		c.emit(code.OpNull)
	}

	return err
}
func (c *Compiler) changOperand(opPos, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
}
func (c *Compiler) replaceInstruction(opPost int, operand []byte) {
	for i := 0; i < len(operand); i++ {
		c.scopes[c.scopeIndex].instructions[opPost+i] = operand[i]
	}
}
func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}
func (c *Compiler) enterScope() {
	c.scopes = append(c.scopes, &CompilationScope{})
	c.scopeIndex++
}
func (c *Compiler) leaveScope() code.Instructions {
	ins := c.currentInstructions()
	c.scopes = c.scopes[:c.scopeIndex]
	c.scopeIndex--
	return ins
}
func (c *Compiler) createScope() {
	c.scopes = append(c.scopes, &CompilationScope{})
}
func (c *Compiler) fun(fun_ *ast.FunExpression) error {
	var tmpTSymbol *Symbol
	if fun_.Name != nil {
		tmpTSymbol = c.symbolTable.SetSymbol(fun_.Name.Value)
	}
	symbol := NewSymbolTable(c.symbolTable)
	c.enterScope() //开启新的作用域
	c.symbolTable = symbol

	for _, param := range fun_.Params {
		sy := c.symbolTable.SetSymbol(param.Value)
		c.emit(code.OpSetLocal, sy.index)
	}
	err := c.callBack(fun_.Block)
	if err != nil {
		return err
	}
	if c.lastInstructionIs(code.OpPop) {
		c.replaceLastPosWithReturn()
	}
	ins := c.leaveScope() //恢复作用域
	fmt.Println("fun")
	fmt.Println(ins.String())
	c.symbolTable = symbol.top
	for _, free := range symbol.free {
		c.symbolEmitGet(free)
	}
	c.emit(code.OpLoadFun, c.addConstant(&object.CompliedFun{Instructions: ins, NumLocal: symbol.index}), len(symbol.free))
	if fun_.Name != nil {
		if tmpTSymbol.types == Global {
			c.emit(code.OpSetGlobal, tmpTSymbol.index)
		} else {
			c.emit(code.OpSetLocal, tmpTSymbol.index)
		}
	}
	return nil
}
func (c *Compiler) replaceLastPosWithReturn() {
	pos := c.scopes[c.scopeIndex].last.Pos

	c.replaceInstruction(pos, code.Make(code.OpReturn))
	c.scopes[c.scopeIndex].last.Op = code.OpReturn
}
func (c *Compiler) WherePop(exp ast.Expression) {
	switch exp.(type) {
	//case *ast.CallExpression:
	//	return
	case *ast.AssigExpression:
		return
	case *ast.FunExpression:
		return
	case *ast.ForExpression:
		return
	default:
		c.emit(code.OpPop)
	}
}
func (c *Compiler) symbolEmitGet(symbol *Symbol) int {
	var pos int
	switch symbol.types {
	case Global:
		pos = c.emit(code.OpGetGlobal, symbol.index)
	case Local:
		pos = c.emit(code.OpGetLocal, symbol.index)
	case Free:
		pos = c.emit(code.OpGetFree, symbol.index)
	}
	if c.getEmitPos == true {
		c.pos = pos
		c.getEmitPos = false
	}
	return pos
}
func (c *Compiler) symbolEmitSet(symbol *Symbol) int {
	switch symbol.types {
	case Global:
		return c.emit(code.OpSetGlobal, symbol.index)
	case Local:
		return c.emit(code.OpSetLocal, symbol.index)
	}
	return -1
}
func (c *Compiler) symbolEmitDel(symbol *Symbol) int {
	switch symbol.types {
	case Global:
		return c.emit(code.OpDelGlobal, symbol.index)
	case Local:
		return c.emit(code.OpDelLocal, symbol.index)
	}
	return -1
}
func (c *Compiler) forExpression(for_ *ast.ForExpression) error {
	//0 let a = 0  //这个位置 //get 声明 //
	//1  get
	//2 a < 10
	//3 OpJumpNotTrueThy 7
	//4 a++  //set opset
	//5 block
	//6 OpJump 1
	//7 ...

	err := c.callBack(for_.Left)
	if err != nil {
		return err
	}

	//LetSymbol, ok := c.symbolTable.GetSymbol(for_.Left.Name.Value)
	//if !ok {
	//	return errors.New(fmt.Sprintf("不能引用没有定义的变量 %s", for_.Left.Name.Value))
	//}
	//getIndex := c.symbolEmitGet(LetSymbol)
	c.EmitGetPos()
	err = c.callBack(for_.Mid)
	jumpIndex, errs := c.getEmitPos_()
	if errs != nil {
		return err
	}
	if err != nil {
		return err
	}
	index := c.emit(code.OpJumpNotTrueThy, 999)
	err = c.callBack(for_.Block)
	if err != nil {
		return err
	}
	err = c.callBack(for_.Right)
	if err != nil {
		return err
	}
	c.emit(code.OpJump, jumpIndex)
	if c.lastInstructionIs(code.OpPop) {
		c.delLastPop()
	}
	c.changOperand(index, len(c.currentInstructions()))
	//c.symbolEmitDel(LetSymbol)
	c.symbolTable.DelSymbol(for_.Left.Name.Value)
	return nil
}
func (c *Compiler) EmitGetPos() {
	c.getEmitPos = true
}
func (c *Compiler) getEmitPos_() (int, error) {
	if c.pos > -1 {
		tmp := c.pos
		c.pos = tmp
		return c.pos, nil
	}
	return -1, errors.New("回退位置失败")
}
func (c *Compiler) Assig(n *ast.AssigExpression) error {
	//普通赋值
	//数组赋值
	//...
	switch n.Name.(type) {
	case *ast.IndexExpression:
		return c.AssigArray(n)
	case *ast.Identifier:
		return c.AssigOrdinary(n)
	}
	return nil
}
func (c *Compiler) AssigArray(node *ast.AssigExpression) error {
	indexNode := node.Name.(*ast.IndexExpression)
	name := indexNode.Left.(*ast.Identifier)
	symbol, ok := c.symbolTable.GetSymbol(name.Value)
	if !ok {
		return errors.New(fmt.Sprintf("不能对一个没有声明的变量赋值 %s", name))
	}
	err := c.callBack(node.Value)
	if err != nil {
		return err
	}
	err = c.callBack(indexNode.Index)
	if err != nil {
		return nil
	}
	if symbol.types == Global {
		c.emit(code.OpSetIndexGlobal, symbol.index)
	} else {
		c.emit(code.OpSetIndexLocal, symbol.index)
	}
	return nil
}
func (c *Compiler) AssigOrdinary(node *ast.AssigExpression) error {
	name := node.Name.(*ast.Identifier)
	symbol, ok := c.symbolTable.GetSymbol(name.Value)
	if !ok {
		return errors.New(fmt.Sprintf("不能对一个没有声明的变量赋值 %s", name))
	}
	err := c.callBack(node.Value)
	if err != nil {
		return err
	}
	c.symbolEmitSet(symbol)
	return nil
}
