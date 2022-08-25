package ast

import (
	"bytes"
	"hek/token"
)

type PrefixExpression struct {
	Token token.Token
	Right Expression
}

func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Token.Type.ToString())
	out.WriteString(" " + p.Right.String())
	out.WriteString(")")

	return out.String()
}

func (p *PrefixExpression) expressionNode() {

}
