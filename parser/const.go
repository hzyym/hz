package parser

const (
	_ int = iota
	LOWEST
	EQUALS      //== or !=
	LESSGREATER // > or <
	SUM         //+ -
	PRODUCT     //* /
	PREFIX      //-x os !x
	CALL        //CALL test(x,y)
	LBRACKET    // [
)
