package ast

import (
	"bytes"
	"hek/token"
)

type IFExpression struct {
	Token       token.Token
	Condition   Expression      //条件表达式
	Consequence *BlockStatement //true  语法块
	Alternative *BlockStatement //false 语法块
}

func (I *IFExpression) TokenLiteral() string {
	return I.Token.Literal
}

func (I *IFExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(I.Condition.String())
	out.WriteString(" ")
	out.WriteString(I.Consequence.String())

	if I.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(I.Alternative.String())
	}
	return out.String()
}

func (I *IFExpression) expressionNode() {

}
