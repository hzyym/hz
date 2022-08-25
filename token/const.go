package token

const (
	EOF Type = iota
	ILLEGAL

	// IDENT 标识符 add i z ...
	IDENT

	//运算符

	ASSIGN   //=
	PLUS     //+
	MINUS    //-
	BANG     //!
	ASTERISK //*
	SLASH    //  /
	EQ       // ==
	NotEq    // !=
	TwoPlus
	TwoMinus

	LT //<
	GT // >
	//分隔符
	COMMA     //,
	SEMICOLON //;

	LPAREN //(
	RPAREN //)

	LBRACE //{
	RBRACE //}

	LBRACKET //[
	RBRACKET //]
	COLON    //:
	//关键字

	FUNCTION //fun
	LET      //let
	IF       //if
	ELSE     //else
	RETURN   //return
	FALSE    //false
	TRUE     //true
	FOR      //for
	BREAK    //break
	CONTINUE //continue
	//类型
	INT
	String
)
