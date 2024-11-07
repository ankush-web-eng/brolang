package parser

import (
	"fmt"
	"strconv"

	"github.com/ankush-web-eng/brolang/ast"
	"github.com/ankush-web-eng/brolang/lexer"
	"github.com/ankush-web-eng/brolang/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	// fmt.Printf("Parsing token: %s (%s)\n", p.curToken.Type, p.curToken.Literal)
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.IDENT:
		if p.peekTokenIs(token.ASSIGN) {
			return p.parseAssignStatement()
		}
		return p.parseExpressionStatement()
	case token.PRINT:
		return p.parsePrintStatement()
	case token.IF:
		return p.parseExpressionStatement()
	case token.WHILE:
		return p.parseExpressionStatement()
	case token.FOR:
		return p.parseExpressionStatement()
	case token.BREAK:
		stmt := &ast.BreakStatement{Token: p.curToken}
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		return stmt
	case token.CONTINUE:
		stmt := &ast.ContinueStatement{Token: p.curToken}
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		return stmt
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	stmt := &ast.AssignStatement{Token: p.curToken}

	// The variable name (IDENT) is the current token
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Expect '=' after the identifier
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression()

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parsePrintStatement() *ast.PrintStatement {
	stmt := &ast.PrintStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken() // moving past '('

	stmt.Expression = p.parseExpression()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression()
	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	var leftExp ast.Expression

	// Parse left-hand side of the expression
	switch p.curToken.Type {
	case token.INT:
		leftExp = p.parseIntegerLiteral()
	case token.STRING:
		leftExp = p.parseStringLiteral()
	case token.TRUE, token.FALSE:
		leftExp = p.parseBoolean()

	case token.IF:
		leftExp = p.parseIfExpression()
	case token.WHILE:
		leftExp = p.parseWhileExpression()
	case token.FOR:
		leftExp = p.parseForExpression()

	case token.LBRACKET:
		leftExp = p.parseArrayLiteral()
	case token.IDENT:
		if p.peekTokenIs(token.LPAREN) {
			leftExp = p.parseCallExpression(p.parseIdentifier())
		} else {
			leftExp = p.parseIdentifier()
		}

	default:
		return nil
	}

	// If the next token is an infix operator, parse it as an infix expression
	for p.peekTokenIs(token.PLUS) || p.peekTokenIs(token.MINUS) ||
		p.peekTokenIs(token.ASTERISK) || p.peekTokenIs(token.SLASH) || p.peekTokenIs(token.MOD) {
		p.nextToken() // Move to the operator
		leftExp = p.parseInfixExpression(leftExp)
	}

	if p.peekTokenIs(token.LBRACKET) {
		p.nextToken()
		leftExp = p.parseIndexExpression(leftExp)
	}

	return leftExp
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression()

	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression())

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression())
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token:  p.curToken,
		ElseIf: []*ast.IfExpression{},
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseSimpleExpression()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	// Parse else-if statements
	for p.peekTokenIs(token.ELSE_IF) {
		p.nextToken() // move to ELSE_IF
		elseIfExp := &ast.IfExpression{Token: p.curToken}

		if !p.expectPeek(token.LPAREN) {
			return nil
		}

		p.nextToken()
		elseIfExp.Condition = p.parseSimpleExpression()

		if !p.expectPeek(token.RPAREN) {
			return nil
		}

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		elseIfExp.Consequence = p.parseBlockStatement()
		expression.ElseIf = append(expression.ElseIf, elseIfExp)
	}

	// Parse else statement
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseSimpleExpression() ast.Expression {
	// Parse the initial expression
	left := p.parseExpression()
	if left == nil {
		return nil
	}

	// If the next token is a comparison operator, create an infix expression
	if p.peekTokenIs(token.GT) || p.peekTokenIs(token.LT) ||
		p.peekTokenIs(token.EQ) || p.peekTokenIs(token.NOT_EQ) ||
		p.peekTokenIs(token.GTE) || p.peekTokenIs(token.LTE) {
		p.nextToken() // move to the operator
		operator := p.curToken.Literal

		p.nextToken() // move to the right side
		right := p.parseExpression()
		if right == nil {
			return nil
		}

		return &ast.InfixExpression{
			Token:    p.curToken,
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseWhileExpression() ast.Expression {
	expression := &ast.WhileExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseSimpleExpression()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// Create new scope for the loop body
	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseForExpression() ast.Expression {
	expression := &ast.ForExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// Parse initialization
	p.nextToken() // Move past '('
	if !p.curTokenIs(token.SEMICOLON) {
		expression.Init = p.parseStatement()
		if expression.Init == nil {
			return nil
		}
	}

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	// Parse condition
	p.nextToken() // Move past semicolon
	if !p.curTokenIs(token.SEMICOLON) {
		expression.Condition = p.parseSimpleExpression()
		if expression.Condition == nil {
			return nil
		}
	}

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	// Parse update
	p.nextToken() // Move past semicolon
	if !p.curTokenIs(token.RPAREN) {
		expression.Update = p.parseStatement()
		if expression.Update == nil {
			return nil
		}
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	p.nextToken()
	expression.Right = p.parseExpression()
	return expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Arguments = []ast.Expression{}

	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
		return exp
	}

	exp.Arguments = append(exp.Arguments, p.parseExpression())

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		exp.Arguments = append(exp.Arguments, p.parseExpression())
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken() // Move past '['
	exp.Index = p.parseExpression()

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Sahi se code likhna bhi nahi aa raha tere se! %s kaha se aa gaya %s se pehle!!!!",
		p.peekToken.Literal, t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}
