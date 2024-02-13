package ast

import (
	"bytes"
	"strings"

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
	Token      token.Token // The first token of the statement
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}
func (e *ExpressionStatement) String() string {
	var buffer bytes.Buffer
	if e.Expression != nil {
		buffer.WriteString(e.Expression.String())
	}
	return buffer.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode() {}
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}
func (p *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
	Left     Expression
}

func (i *InfixExpression) expressionNode() {}
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}
func (i *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" ")
	out.WriteString(i.Operator)
	out.WriteString(" ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return b.Token.Literal
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) expressionNode() {}
func (s *StringLiteral) TokenLiteral() string {
	return s.Token.Literal
}
func (s *StringLiteral) String() string {
	return s.Value
}

type IfExpression struct {
	Token       token.Token // The IF Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode() {}
func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(i.Condition.String())
	out.WriteString("{")
	if i.Consequence != nil {
		out.WriteString(i.Consequence.String())
	}
	out.WriteString("}")

	if i.Alternative != nil {
		out.WriteString("else {")
		out.WriteString(i.Alternative.String())
		out.WriteString("}")
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) expressionNode() {}
func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}
func (b *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range b.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // The fn token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode() {}
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString("{")
	out.WriteString(f.Body.String())
	out.WriteString("}")

	return out.String()
}

type CallExpression struct {
	Token     token.Token // The fn token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, p := range ce.Arguments {
		args = append(args, p.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type Array struct {
	Token    token.Token // The [ token
	Elements []Expression
}

func (a *Array) expressionNode() {}
func (a *Array) TokenLiteral() string {
	return a.Token.Literal
}
func (a *Array) String() string {
	var out bytes.Buffer
	elements := []string{}

	for _, el := range a.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
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

type IndexExpression struct {
	Token token.Token // The token.LBRACKET token
	Left  Expression
	Index Expression
}

func (i *IndexExpression) expressionNode() {}
func (i *IndexExpression) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString("[")
	out.WriteString(i.Index.String())
	out.WriteString("])")
	return out.String()
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
