package object

import "fmt"

const (
	INTEGER          = "INTEGER"
	BOOLEAN          = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VAL"
	NULL             = "NULL"
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

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%v", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN
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

type Null struct{}

func (n *Null) Inspect() string {
	return "NULL"
}

func (n *Null) Type() ObjectType {
	return NULL
}
