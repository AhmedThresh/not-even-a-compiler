package parser

import (
	"fmt"

	"github.com/AhmedThresh/not-even-a-compiler/pkg/ast"
	"github.com/AhmedThresh/not-even-a-compiler/pkg/lexer"
	"github.com/AhmedThresh/not-even-a-compiler/pkg/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	Errors       []string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		Errors: []string{},
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) addError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", p.peekToken.Type, t)
	p.Errors = append(p.Errors, msg)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.currentToken.Type != token.EOF {
		if p.currentToken.Type == token.LET {
			letStatement := p.parseLetStatement()
			if letStatement != nil {
				program.Statements = append(program.Statements, letStatement)
			}
		} else if p.currentToken.Type == token.RETURN {
			returnStatement := p.parseReturnStatement()
			if returnStatement != nil {
				program.Statements = append(program.Statements, returnStatement)
			}
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	identifier := p.parseIdentifier()
	if !p.expectPeek(token.ASSIGN) {
		return nil
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
	return &ast.ReturnStatement{
		Token: p.currentToken,
	}
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
