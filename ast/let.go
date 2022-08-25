package ast

import (
	"bytes"
	"hek/token"
)

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

func (l *LetStatement) statementNode() {}

func (l *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")

	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	return out.String()
}
