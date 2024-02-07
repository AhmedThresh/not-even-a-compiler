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

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

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

	val, err := strconv.ParseInt(p.currentToken.Literal, 0, 10)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("cannot parse integer %s", p.currentToken.Literal))
		return nil
	}

	integerLiteral.Value = val
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
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	identifier := p.parseIdentifier()
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for p.currentToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return &ast.LetStatement{
		Token: token.Token{
			Type:    token.LET,
			Literal: "let",
		},
		Name: identifier,
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{
		Token: p.currentToken,
	}

	p.nextToken()

	for p.currentToken.Type != token.SEMICOLON {
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

func (p *Parser) parseIdentifierExpression() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}
