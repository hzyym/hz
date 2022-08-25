package ast

import (
	"bytes"
	"hek/token"
)

type ForExpression struct {
	Token token.Token
	Left  *LetStatement
	Mid   Expression
	Right Expression
	Block *BlockStatement
}

func (f *ForExpression) TokenLiteral() string {
	return f.Token.Literal
}

func (f *ForExpression) String() string {
	var out bytes.Buffer
	out.WriteString("for(" + f.Left.String() + ";" + f.Mid.String() + ";" + f.Right.String() + ") {\n")
	out.WriteString(f.Block.String() + "\n}")
	return out.String()
}

func (f *ForExpression) expressionNode() {
}
