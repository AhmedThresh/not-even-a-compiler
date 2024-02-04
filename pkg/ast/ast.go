package ast

import (
	"bytes"

	"github.com/AhmedThresh/not-even-a-compiler/pkg/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type LetStatement struct {
	Token token.Token // The token.LET token
	Name  *Identifier
	Value Expression
}

func (l *LetStatement) statementNode() {}
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

func (l *LetStatement) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(l.Token.Literal + " ")
	buffer.WriteString(l.Name.String() + " = ")

	if l.Value != nil {
		buffer.WriteString(l.Value.String())
	}

	buffer.WriteString(";")
	return buffer.String()
}

type ReturnStatement struct {
	Token token.Token // The token.RETURN token
	Value Expression
}

func (r *ReturnStatement) statementNode() {}
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}
func (r *ReturnStatement) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(r.Token.Literal + " ")

	if r.Value != nil {
		buffer.WriteString(r.Value.String())
	}

	buffer.WriteString(";")
	return buffer.String()
}

type ExpressionStatement struct {
	Token token.Token // The first token of the statement
	Value Expression
}

func (e *ExpressionStatement) statementNode() {}
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}
func (e *ExpressionStatement) String() string {
	var buffer bytes.Buffer
	if e.Value != nil {
		buffer.WriteString(e.Value.String())
	}
	return buffer.String()
}

type Identifier struct {
	Token token.Token // The token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var buffer bytes.Buffer
	for _, s := range p.Statements {
		// TODO: Should handle error here
		_, _ = buffer.WriteString(s.String())
	}

	return buffer.String()
}
