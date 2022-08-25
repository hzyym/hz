package ast

import (
	"bytes"
	"hek/token"
)

type HashExpression struct {
	Token token.Token
	Value map[Expression]Expression
}

func (h *HashExpression) TokenLiteral() string {
	return h.Token.Literal
}

func (h *HashExpression) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for index, val := range h.Value {
		out.WriteString(index.String() + ":" + val.String() + ",")
	}
	out.WriteString("}")
	return out.String()
}

func (h *HashExpression) expressionNode() {

}
