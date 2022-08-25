package lexer

import (
	"fmt"
	"hek/token"
	"testing"
)

func TestNewLexer(t *testing.T) {
	input := "\"test\""

	l := NewLexer(input)

	for {
		tok := l.NextToke()
		if tok.Type != token.EOF {
			fmt.Println(tok)
		} else {
			break
		}
	}
}
