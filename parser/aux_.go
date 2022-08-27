package parser

import (
	"fmt"
	"hek/token"
)

var precedence = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NotEq:    EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: LBRACKET,
}

func (p *Parser) peekPrecedence() int {
	if val, ok := precedence[p.peekToken.Type]; ok {
		return val
	}
	return LOWEST
}
func (p *Parser) curPrecedence() int {
	if val, ok := precedence[p.curToken.Type]; ok {
		return val
	}
	return LOWEST
}
func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}
func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekErrors(t)
	return false
}
func (p *Parser) peekErrors(t token.Type) {
	p.errors = append(p.errors, fmt.Sprintf("peek token is '%s' not is '%s'", p.peekToken.Type.ToString(), t.ToString()))
}
func (p *Parser) integerLiteralErrors(t token.Token) {
	p.errors = append(p.errors, fmt.Sprintf("%s not int type", t.Literal))
}
func (p *Parser) Errors() []string {
	return p.errors
}
func (p *Parser) registerPrefixFun(types token.Type, fun prefixParseFun) {
	p.prefixParseFus[types] = fun
}
func (p *Parser) registerInfixFun(types token.Type, fun infixParseFun) {
	p.infixParseFus[types] = fun
}
func (p *Parser) registerFunALL() {
	//prefix
	p.registerPrefixFun(token.IDENT, p.parseIdentifier)
	p.registerPrefixFun(token.INT, p.parseIntegerLiteral)
	p.registerPrefixFun(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixFun(token.BANG, p.parsePrefixExpression)
	p.registerPrefixFun(token.TRUE, p.parseBoolExpression)
	p.registerPrefixFun(token.FALSE, p.parseBoolExpression)
	p.registerPrefixFun(token.LPAREN, p.parseGroupExpression)
	p.registerPrefixFun(token.IF, p.parseIFExpression)
	p.registerPrefixFun(token.FUNCTION, p.parseFunExpression)
	p.registerPrefixFun(token.String, p.parseString)
	p.registerPrefixFun(token.LBRACKET, p.parseArrayExpression)
	p.registerPrefixFun(token.LBRACE, p.parseHashExpression)
	p.registerPrefixFun(token.FOR, p.parseForExpression)

	//infix
	p.registerInfixFun(token.SLASH, p.parseInfixExpression)
	p.registerInfixFun(token.LT, p.parseInfixExpression)
	p.registerInfixFun(token.GT, p.parseInfixExpression)
	p.registerInfixFun(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFun(token.EQ, p.parseInfixExpression)
	p.registerInfixFun(token.NotEq, p.parseInfixExpression)
	p.registerInfixFun(token.MINUS, p.parseInfixExpression)
	p.registerInfixFun(token.PLUS, p.parseInfixExpression)
	p.registerInfixFun(token.LPAREN, p.parseInfixCallExpression)
	p.registerInfixFun(token.LBRACKET, p.parseIndexExpression)
}
