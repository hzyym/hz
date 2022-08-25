package ast

import (
	"bytes"
	"hek/token"
)

type SuffixExpression struct {
	Token token.Token
	Left  *Identifier
}

func (s *SuffixExpression) TokenLiteral() string {
	return s.Token.Literal
}

func (s *SuffixExpression) String() string {
	var out bytes.Buffer
	out.WriteString(s.Left.String())
	out.WriteString(s.Token.Type.ToString())
	return out.String()
}

func (s *SuffixExpression) expressionNode() {
	//TODO implement me
	panic("implement me")
}
