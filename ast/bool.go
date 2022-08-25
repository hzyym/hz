package ast

import "hek/token"

type BoolExpression struct {
	Token token.Token
	Value bool
}

func (b *BoolExpression) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BoolExpression) String() string {
	return b.Token.Literal
}

func (b *BoolExpression) expressionNode() {

}
