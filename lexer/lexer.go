package lexer

import (
	"github.com/jumballaya/servo/token"
)

type Lexer struct {
	input        string
	position     int  // current position (current)
	readPosition int  // current reading position (after current)
	ch           byte // current char
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '%':
		tok = newToken(token.MODULO, l.ch)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString('"')
	case '\'':
		tok.Type = token.STRING
		tok.Literal = l.readString('\'')
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '#':
		tok.Type = token.COMMENT
		tok.Literal = l.readString('\n')
	default:
		if isLetter(l.ch) {
			if l.ch == 'i' && l.peekChar() == 'm' {
				tok = l.readImport()
				l.readChar()
				return tok
			}
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString(initial byte) string {
	position := l.position + 1
	if initial == '\'' || initial == '"' {
		for {
			l.readChar()
			if l.ch == '\'' || l.ch == '"' || l.ch == 0 {
				break
			}
		}
		return l.input[position:l.position]
	} else {
		for {
			l.readChar()
			if l.ch == initial || l.ch == 0 {
				break
			}
		}
		return l.input[position:l.position]
	}
}

func (l *Lexer) readImport() token.Token {
	if l.peekChar() != 'f' {

	}
	stmt := l.readString(' ')
	if stmt != "mport" {
		illegal := newToken(token.ILLEGAL, l.ch)
		return illegal
	}

	// Case: import map from 'Array';
	if l.ch != '\'' && l.ch != '"' {
		module := l.readString(' ')
		fromStmt := l.readString(' ')

		if fromStmt != "from" {
			illegal := newToken(token.ILLEGAL, l.ch)
			return illegal
		}

		l.skipWhiteSpace()

		if l.ch != '\'' && l.ch != '"' {
			illegal := newToken(token.ILLEGAL, l.ch)
			return illegal
		}

		object := l.readString('\'')

		return token.Token{Type: token.IMPORT, Literal: object + ":" + module}
	}

	illegal := newToken(token.ILLEGAL, l.ch)
	return illegal
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '\''
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
