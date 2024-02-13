package parser

import (
	"fmt"
	"strconv"

	"github.com/AhmedThresh/not-even-a-compiler/pkg/ast"
	"github.com/AhmedThresh/not-even-a-compiler/pkg/lexer"
	"github.com/AhmedThresh/not-even-a-compiler/pkg/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var precedence = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()

	// Initialize parsing functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifierExpression)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionExpression)
	p.registerPrefix(token.STRING, p.parseStringLiteralExpression)
	p.registerPrefix(token.LBRACKET, p.parseArray)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) currentPrecedence() int {
	if val, ok := precedence[p.currentToken.Type]; ok {
		return val
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if val, ok := precedence[p.peekToken.Type]; ok {
		return val
	}

	return LOWEST
}

func (p *Parser) addError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", p.peekToken.Type, t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(token token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[token] = fn
}

func (p *Parser) registerInfix(token token.TokenType, fn infixParseFn) {
	p.infixParseFns[token] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := ast.ExpressionStatement{Token: p.currentToken}

	statement.Expression = p.parseExpression(LOWEST)
	// TODO: Here we should consider corner cases
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return &statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn := p.prefixParseFns[p.currentToken.Type]
	if prefixFn == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	left := prefixFn()
	for p.peekToken.Type != token.SEMICOLON && p.peekPrecedence() > precedence {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return left
		}

		p.nextToken()

		left = infix(left)
	}

	return left
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	integerLiteral := ast.IntegerLiteral{
		Token: p.currentToken,
	}

	val, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("cannot parse integer %s", p.currentToken.Literal))
		return nil
	}

	integerLiteral.Value = int64(val)
	return &integerLiteral
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return &expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := ast.InfixExpression{
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return &expression
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStatement := &ast.LetStatement{
		Token: token.Token{
			Type:    token.LET,
			Literal: "let",
		},
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	identifier := p.parseIdentifier()
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	letStatement.Name = identifier

	p.nextToken()

	letStatement.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return letStatement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{
		Token: p.currentToken,
	}

	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return statement
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type != t {
		p.addError(t)
		return false
	}

	p.nextToken()
	return true
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	i := &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	return i
}

func (p *Parser) parseBoolean() ast.Expression {
	val := false
	if p.currentToken.Type == token.TRUE {
		val = true
	}
	b := &ast.Boolean{
		Token: p.currentToken,
		Value: val,
	}

	return b
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		// TODO: Should handle error here
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := ast.IfExpression{
		Token: p.currentToken,
	}

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

	if p.peekToken.Type == token.ELSE {
		// skip the else token
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}

	return &exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blocks := &ast.BlockStatement{
		Token: p.currentToken,
	}

	p.nextToken()

	for p.currentToken.Type != token.EOF && p.currentToken.Type != token.RBRACE {
		s := p.parseStatement()
		if s != nil {
			blocks.Statements = append(blocks.Statements, s)
		}
		p.nextToken()
	}
	return blocks
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	expression := &ast.FunctionLiteral{
		Token: p.currentToken,
	}

	if !p.expectPeek(token.LPAREN) {
		// TODO: handle error here
		return nil
	}

	p.nextToken()

	expression.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		// TODO: handle error here
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	res := []*ast.Identifier{}

	for p.currentToken.Type != token.RPAREN {
		if p.currentToken.Type != token.COMMA {
			res = append(res, &ast.Identifier{
				Token: p.currentToken,
				Value: p.currentToken.Literal,
			})
		}

		p.nextToken()
	}

	return res
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := ast.CallExpression{
		Token:    p.currentToken,
		Function: function,
	}
	expression.Arguments = p.parseCallArguments()
	return &expression
}

func (p *Parser) parseCallArguments() []ast.Expression {
	res := []ast.Expression{}

	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return res
	}

	p.nextToken()

	res = append(res, p.parseExpression(LOWEST))

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		res = append(res, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return res
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	res := &ast.IndexExpression{
		Token: p.currentToken,
		Left:  left,
	}

	p.nextToken()

	res.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return res
}

func (p *Parser) parseIdentifierExpression() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseStringLiteralExpression() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseArray() ast.Expression {
	elements := []ast.Expression{}
	arr := &ast.Array{Elements: elements}

	if p.peekToken.Type == token.RBRACKET {
		p.nextToken()
		return arr
	}

	p.nextToken()
	arr.Elements = append(arr.Elements, p.parseExpression(LOWEST))

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		arr.Elements = append(arr.Elements, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return arr
}

func (p *Parser) Errors() []string {
	return p.errors
}
