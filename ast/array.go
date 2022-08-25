package ast

import (
	"bytes"
	"hek/token"
	"strings"
)

type ArrayExpression struct {
	Token token.Token
	Value []Expression
}

func (a *ArrayExpression) TokenLiteral() string {
	return a.Token.Literal
}

func (a *ArrayExpression) String() string {
	var out bytes.Buffer
	var str []string
	for _, expression := range a.Value {
		str = append(str, expression.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(str, ","))
	out.WriteString("]")
	return out.String()
}

func (a *ArrayExpression) expressionNode() {

}
