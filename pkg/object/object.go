package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/AhmedThresh/not-even-a-compiler/pkg/ast"
)

const (
	INTEGER          = "INTEGER"
	STRING           = "STRING"
	BOOLEAN          = "BOOLEAN"
	ARRAY            = "ARRAY"
	FUNCTION         = "FUNCTION"
	RETURN_VALUE_OBJ = "RETURN_VAL"
	ERROR_OBJ        = "ERROR"
	NULL             = "NULL"
	BUILTIN          = "BUILTIN"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER
}

type String struct {
	Value string
}

func (s *String) Inspect() string {
	return fmt.Sprintf("%s", s.Value)
}

func (s *String) Type() ObjectType {
	return STRING
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%v", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN
}

type Array struct {
	Elements []Object
}

func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}

	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

func (a *Array) Type() ObjectType {
	return ARRAY
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

type Error struct {
	Message string
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

func (f *Function) Type() ObjectType {
	return FUNCTION
}

type BuiltinFunction func(...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN
}

type Null struct{}

func (n *Null) Inspect() string {
	return "NULL"
}

func (n *Null) Type() ObjectType {
	return NULL
}
