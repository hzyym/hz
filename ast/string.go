package ast

import "hek/token"

type StringExpression struct {
	Token token.Token
	Value string
}

func (s *StringExpression) TokenLiteral() string {
	return s.Token.Literal
}

func (s *StringExpression) String() string {
	return s.Value
}

func (s *StringExpression) expressionNode() {

}
