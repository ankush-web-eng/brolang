package parser

import (
	"fmt"

	"github.com/ankush-web-eng/brolang/lexer"
)

type ASTNode interface{}

type BinaryExpression struct {
	Left     ASTNode
	Operator lexer.Token
	Right    ASTNode
}

type NumberLiteral struct {
	Value string
}

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken() // read two tokens, so curToken and peekToken are both set
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseExpression() ASTNode {
	left := p.parsePrimary() // parse the left operand

	for p.curToken.Type == lexer.PLUS || p.curToken.Type == lexer.ASSIGN {
		operator := p.curToken
		p.nextToken()
		right := p.parsePrimary() // parse the right operand
		left = &BinaryExpression{Left: left, Operator: operator, Right: right}
	}

	return left
}

// parsePrimary handles primary expressions like numbers.
func (p *Parser) parsePrimary() ASTNode {
	switch p.curToken.Type {
	case lexer.INT:
		value := &NumberLiteral{Value: p.curToken.Literal}
		p.nextToken()
		return value
	default:
		panic(fmt.Sprintf("Unexpected token: %s", p.curToken.Literal))
	}
}
