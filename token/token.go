package token

type Type int
type Token struct {
	Type    Type
	Literal string
}

var keywords = map[string]Type{
	"fun":      FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"return":   RETURN,
	"else":     ELSE,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
}
var typeWords = map[Type]string{
	LET:       "let",
	INT:       "int",
	IF:        "if",
	ELSE:      "else",
	PLUS:      "+",
	MINUS:     "-",
	ASSIGN:    "=",
	BANG:      "!",
	ASTERISK:  "*",
	SLASH:     "/",
	EQ:        "==",
	NotEq:     "!=",
	LT:        "<",
	GT:        ">",
	COMMA:     ",",
	SEMICOLON: ";",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	FUNCTION:  "fun",
	RETURN:    "return",
	FALSE:     "false",
	TRUE:      "true",
	EOF:       "EOF",
	ILLEGAL:   "ILLEGAL",
	IDENT:     "IDENT",
	String:    "string",
	FOR:       "for",
	BREAK:     "break",
	CONTINUE:  "continue",
	TwoPlus:   "++",
	TwoMinus:  "--",
}

func LookupIdent(ident string) Type {
	if t, ok := keywords[ident]; ok {
		return t
	}
	return IDENT
}
func (T Type) ToString() string {
	return typeWords[T]
}
