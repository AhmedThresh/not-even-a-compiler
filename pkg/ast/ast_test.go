package ast

import (
	"github.com/AhmedThresh/not-even-a-compiler/pkg/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVars"},
					Value: "anotherVars",
				},
			},
		},
	}
	if program.String() != "let myVar = anotherVars;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
