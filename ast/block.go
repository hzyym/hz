package ast

import (
	"bytes"
	"hek/token"
)

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	for _, statement := range b.Statements {
		out.WriteString(statement.String())
	}
	out.WriteString("}")
	return out.String()
}

func (b BlockStatement) statementNode() {

}
