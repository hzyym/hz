package ast

import (
	"bytes"
	"hek/token"
)

type InfixExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(" + i.Left.String() + " ")
	out.WriteString(i.Token.Type.ToString())
	out.WriteString(" " + i.Right.String() + ")")
	return out.String()
}

func (i *InfixExpression) expressionNode() {

}
