package ast

import (
	"bytes"
)

type AssigExpression struct {
	Name  Expression
	Value Expression
}

func (a *AssigExpression) TokenLiteral() string {
	return a.Name.TokenLiteral()
}

func (a *AssigExpression) expressionNode() {}

func (a *AssigExpression) String() string {
	var out bytes.Buffer

	out.WriteString(a.Name.String())
	out.WriteString(" = ")
	if a.Value != nil {
		out.WriteString(a.Value.String())
	}
	return out.String()
}
