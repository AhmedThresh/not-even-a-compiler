package lexer

import "github.com/AhmedThresh/not-even-a-compiler/pkg/token"

type Lexer struct {
	code            string
	currentPosition int  // current position in input
	readPosition    int  // current reading position
	currentCh       byte // current char under examination
}

func NewLexer(code string) *Lexer {
	l := &Lexer{
		code: code,
	}
	l.readCh()
	return l
}

func (l *Lexer) readCh() {
	if l.readPosition >= len(l.code) {
		l.currentCh = 0
	} else {
		l.currentCh = l.code[l.readPosition]
	}

	l.currentPosition = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token
	l.skipWhitespace()
	switch l.currentCh {
	case '=':
		if l.peekChar() == '=' {
			ch := l.currentCh
			l.readCh()
			t = token.Token{
				Literal: string(ch) + string(l.currentCh),
				Type:    token.EQ,
			}
		} else {
			t = token.NewToken(token.ASSIGN, l.currentCh)
		}
	case '+':
		t = token.NewToken(token.PLUS, l.currentCh)
	case '-':
		t = token.NewToken(token.MINUS, l.currentCh)
	case '!':
		if l.peekChar() == '=' {
			ch := l.currentCh
			l.readCh()
			t = token.Token{
				Literal: string(ch) + string(l.currentCh),
				Type:    token.NOT_EQ,
			}
		} else {
			t = token.NewToken(token.BANG, l.currentCh)
		}
	case '/':
		t = token.NewToken(token.SLASH, l.currentCh)
	case '*':
		t = token.NewToken(token.ASTERISK, l.currentCh)
	case '<':
		t = token.NewToken(token.LT, l.currentCh)
	case '>':
		t = token.NewToken(token.GT, l.currentCh)
	case ';':
		t = token.NewToken(token.SEMICOLON, l.currentCh)
	case ',':
		t = token.NewToken(token.COMMA, l.currentCh)
	case '(':
		t = token.NewToken(token.LPAREN, l.currentCh)
	case ')':
		t = token.NewToken(token.RPAREN, l.currentCh)
	case '{':
		t = token.NewToken(token.LBRACE, l.currentCh)
	case '}':
		t = token.NewToken(token.RBRACE, l.currentCh)
	case 0:
		t.Literal = ""
		t.Type = token.EOF
	default:
		if isLetter(l.currentCh) {
			t.Literal = l.readIdentifier()
			t.Type = token.LookupIdent(t.Literal)
			return t
		} else if isDigit(l.currentCh) {
			t.Type = token.INT
			t.Literal = l.readNumber()
			return t
		}
		t = token.NewToken(token.ILLEGAL, l.currentCh)
	}

	l.readCh()
	return t
}

func (l *Lexer) skipWhitespace() {
	for l.currentCh == ' ' || l.currentCh == '\t' || l.currentCh == '\n' || l.currentCh == '\r' {
		l.readCh()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.code) {
		return 0
	}
	return l.code[l.readPosition]
}

func (l *Lexer) readIdentifier() string {
	position := l.currentPosition
	for isLetter(l.currentCh) {
		l.readCh()
	}

	return l.code[position:l.currentPosition]
}

func (l *Lexer) readNumber() string {
	position := l.currentPosition
	for isDigit(l.currentCh) {
		l.readCh()
	}

	return l.code[position:l.currentPosition]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
