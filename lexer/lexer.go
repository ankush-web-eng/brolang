package lexer

type TokenType string

const (
	EOF     = "EOF"
	INT     = "INT"
	PLUS    = "+"
	UNKNOWN = "UNKNOWN"
)

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() Token {
	var tok Token

	switch l.ch {
	case '+':
		tok = Token{Type: PLUS, Literal: string(l.ch)}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		tok.Literal = l.readNumber()
		tok.Type = INT
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		tok = Token{Type: UNKNOWN, Literal: string(l.ch)}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readNumber() string {
	position := l.position
	for l.ch >= '0' && l.ch <= '9' {
		l.readChar()
	}
	return l.input[position:l.position]
}
