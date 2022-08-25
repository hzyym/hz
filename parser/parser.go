package parser

import (
	"hek/ast"
	"hek/lexer"
	"hek/token"
	"strconv"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string

	prefixParseFus map[token.Type]prefixParseFun
	infixParseFus  map[token.Type]infixParseFun
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		prefixParseFus: make(map[token.Type]prefixParseFun),
		infixParseFus:  make(map[token.Type]infixParseFun),
	}

	p.registerFunALL()
	p.nextToken()
	p.nextToken()
	return p
}
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToke()
}
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}
	return stmt
}
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}
	return stmt
}
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFus[p.curToken.Type]

	if prefix == nil {
		return nil
	}
	left := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFus[p.peekToken.Type]
		if infix == nil {
			return left
		}
		p.nextToken()
		left = infix(left)
	}
	return left
}
func (p *Parser) parseIdentifier() ast.Expression {
	id := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if p.peekTokenIs(token.ASSIGN) {
		exp := &ast.AssigExpression{Name: id}
		p.nextToken()
		p.nextToken()
		exp.Value = p.parseExpression(LOWEST)
		return exp
	} else if p.peekTokenIs(token.TwoPlus) || p.peekTokenIs(token.TwoPlus) {
		exp := &ast.SuffixExpression{Left: &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}}
		p.nextToken()
		exp.Token = p.curToken
		return exp
	}
	return id
}
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.integerLiteralErrors(p.curToken)
		return nil
	}
	lit.Value = val
	return lit
}
func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token: p.curToken,
	}
	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)

	return exp
}
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{Token: p.curToken, Left: left}

	precedence_ := p.curPrecedence()
	p.nextToken()

	exp.Right = p.parseExpression(precedence_)
	return exp
}
func (p *Parser) parseBoolExpression() ast.Expression {
	return &ast.BoolExpression{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}
func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}
func (p *Parser) parseIFExpression() ast.Expression {
	exp := &ast.IFExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	exp.Consequence = p.parseBlockStatement()

	if !p.peekTokenIs(token.ELSE) {
		return exp
	}
	p.nextToken()
	//p.nextToken()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	exp.Alternative = p.parseBlockStatement()
	return exp
}
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}

	p.nextToken()
	for !p.curTokenIs(token.EOF) && !p.curTokenIs(token.RBRACE) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}
func (p *Parser) parseFunExpression() ast.Expression {
	exp := &ast.FunExpression{Token: p.curToken}

	if p.peekTokenIs(token.IDENT) {
		//函数命名
		p.nextToken()
		exp.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
	} else {
		if !p.peekTokenIs(token.LPAREN) {
			return nil
		}
		p.nextToken()
	}

	exp.Params = p.parseFunParams()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	exp.Block = p.parseBlockStatement()
	return exp
}
func (p *Parser) parseFunParams() []*ast.Identifier {
	var ids []*ast.Identifier

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return nil
	}
	p.nextToken()
	ids = append(ids, &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() //把符号 , 移到 当前token位置
		p.nextToken() //把 参数 移到 当前token位置
		ids = append(ids, &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		})
	}
	if !p.peekTokenIs(token.RPAREN) {
		return nil
	}
	p.nextToken()
	return ids
}
func (p *Parser) parseCallExpression() ast.Expression {
	return nil
}
func (p *Parser) parseInfixCallExpression(exp ast.Expression) ast.Expression {
	callExp := &ast.CallExpression{
		Token: p.curToken,
		Fun:   exp,
	}
	callExp.Params = p.parseCallParams()
	return callExp
}
func (p *Parser) parseCallParams() []ast.Expression {
	var exps []ast.Expression

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return nil
	}
	p.nextToken()
	exps = append(exps, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() //把符号 , 移到 当前token位置
		p.nextToken() //把 参数 移到 当前token位置
		exps = append(exps, p.parseExpression(LOWEST))
	}
	if !p.peekTokenIs(token.RPAREN) {
		return nil
	}
	p.nextToken()
	return exps
}
func (p *Parser) parseString() ast.Expression {
	return &ast.StringExpression{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}
func (p *Parser) parseArrayExpression() ast.Expression {
	exp := &ast.ArrayExpression{Token: p.curToken}
	exp.Value = p.parseParamList(token.RBRACKET)
	return exp
}
func (p *Parser) parseParamList(end token.Type) []ast.Expression {
	var exps []ast.Expression

	if p.peekTokenIs(end) {
		p.nextToken()
		return nil
	}
	p.nextToken()
	exps = append(exps, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() //把符号 , 移到 当前token位置
		p.nextToken() //把 参数 移到 当前token位置
		exps = append(exps, p.parseExpression(LOWEST))
	}
	if !p.peekTokenIs(end) {
		return nil
	}
	p.nextToken()
	return exps
}
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}
	p.nextToken()

	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	if p.peekTokenIs(token.ASSIGN) {
		//赋值表达式
		p.nextToken()
		p.nextToken()
		ass := &ast.AssigExpression{Name: exp}
		ass.Value = p.parseExpression(LOWEST)
		return ass
	}
	return exp
}
func (p *Parser) parseHashExpression() ast.Expression {
	exp := &ast.HashExpression{
		Token: p.curToken,
		Value: map[ast.Expression]ast.Expression{},
	}
	p.nextToken()
	for !p.peekTokenIs(token.RBRACE) {
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(LOWEST)
		exp.Value[key] = value

		if p.peekTokenIs(token.RBRACE) {
			break
		}
		if !p.peekTokenIs(token.COMMA) || p.peekTokenIs(token.EOF) {
			return nil
		}
		p.nextToken()
		p.nextToken()
	}
	return exp
}
func (p *Parser) parseForExpression() ast.Expression {
	exp := &ast.ForExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	left := &ast.LetStatement{Token: p.curToken}
	p.nextToken()
	left.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	left.Value = p.parseExpression(LOWEST)
	exp.Left = left
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	p.nextToken()
	exp.Mid = p.parseExpression(LOWEST)
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	p.nextToken()
	exp.Right = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	exp.Block = p.parseBlockStatement()
	return exp
}
