package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"
	BOOL   = "BOOL"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

	// Keywords
	LET   = "bhai_sun"
	PRINT = "bol_bhai"
	INPUT = "suna_bhai"
	IF    = "agar"
	ELSE  = "nahi_to"
	WHILE = "jak_tak"
	FOR   = "shuru_kar"
	TRUE  = "sach"
	FALSE = "jhuth"
)

var keywords = map[string]TokenType{
	"bhai_sun":  LET,
	"bol_bhai":  PRINT,
	"suna_bhai": INPUT,
	"agar":      IF,
	"nahi_to":   ELSE,
	"jak_tak":   WHILE,
	"shuru_kar": FOR,
	"sach":      TRUE,
	"jhuth":     FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
