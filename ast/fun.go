package ast

import (
	"bytes"
	"hek/token"
	"strings"
)

type FunExpression struct {
	Token  token.Token
	Params []*Identifier
	Block  *BlockStatement
	Name   *Identifier
}

func (f *FunExpression) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FunExpression) String() string {
	var out bytes.Buffer
	out.WriteString("fun " + f.TokenLiteral())
	var params []string
	for _, param := range f.Params {
		params = append(params, param.String())
	}
	out.WriteString("(" + strings.Join(params, ",") + ")")
	out.WriteString(f.Block.String())
	return out.String()
}

func (f FunExpression) expressionNode() {

}
