package token

// TokenType defines the type of the token that should be processed
// TokenType can be EOF, INT, LPAREN, etc...
type TokenType string

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343456
	STRING = "STRING" // "foobar"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"
)

// Token defines the unit of the tokenization process
// It is actually a struct that contains the Type of the token and its value
type Token struct {
	// Type is the token type
	Type TokenType
	// Literal is the value of the token
	Literal string
}

// LookupIdent return the token type of a specific token
// It is meant to be used when we encounter a string literal
// If it's a keyword we return it, otherwise it is simply an identifier
func LookupIdent(identifier string) TokenType {
	if t, ok := keywords[identifier]; ok {
		return t
	}
	return IDENT
}

func NewToken(tokenType TokenType, value byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(value),
	}
}
