package ast

import (
	"bytes"
	"hek/token"
	"strings"
)

type CallExpression struct {
	Token  token.Token
	Fun    Expression
	Params []Expression
}

func (c *CallExpression) TokenLiteral() string {
	return c.Token.Literal
}

func (c *CallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(c.Fun.String())
	var params []string
	for _, param := range c.Params {
		params = append(params, param.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	return out.String()
}

func (c *CallExpression) expressionNode() {

}
