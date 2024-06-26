package parser

import (
	"fmt"
	"strconv"

	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.curToken = l.GetToken()

	if p.curToken.Type != "EOF" {
		p.peekToken = l.GetToken()
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.RegisterPrefix(token.IDENTIFIER, p.ParseIdentifier)
	p.RegisterPrefix(token.INT, p.ParseIntegerLiteral)
	p.RegisterPrefix(token.BANG, p.ParsePrefixExpression)
	p.RegisterPrefix(token.MINUS, p.ParsePrefixExpression)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.RegisterInfix(token.PLUS, p.ParseInfixExpression)
	p.RegisterInfix(token.MINUS, p.ParseInfixExpression)
	p.RegisterInfix(token.SLASH, p.ParseInfixExpression)
	p.RegisterInfix(token.ASTERISK, p.ParseInfixExpression)
	p.RegisterInfix(token.EQ, p.ParseInfixExpression)
	p.RegisterInfix(token.NOT_EQ, p.ParseInfixExpression)
	p.RegisterInfix(token.LT, p.ParseInfixExpression)
	p.RegisterInfix(token.GT, p.ParseInfixExpression)
	p.RegisterPrefix(token.TRUE, p.ParseBoolean)
	p.RegisterPrefix(token.FALSE, p.ParseBoolean)
	p.RegisterPrefix(token.LPAREN, p.ParseGroupedExpression)
	p.RegisterPrefix(token.IF, p.ParseIfExpression)
	p.RegisterPrefix(token.FUNCTION, p.ParseFunctionLiteral)
	p.RegisterInfix(token.LPAREN, p.ParseCallFunction)
	p.RegisterPrefix(token.STRING, p.ParseStringLiteral)

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) PeekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) NextToken() {
	p.curToken = p.peekToken

	if p.curToken.Type != "EOF" {
		p.peekToken = p.l.GetToken()
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.ParseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}

	return program
}

func (p *Parser) PeekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) CurPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) ParseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.ParseLetStatement()
	case token.RETURN:
		return p.ParseReturnStatement()
	default:
		return p.ParseExpressionStatement()
	}
}

func (p *Parser) ParseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.ExpectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: string(p.curToken.Literal)}

	if !p.ExpectPeek(token.ASSIGN) {
		return nil
	}

	p.NextToken()

	stmt.Value = p.ParseExpression(LOWEST)

	if p.PeekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) ParseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.NextToken()

	stmt.ReturnValue = p.ParseExpression(LOWEST)

	if p.PeekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) CurTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) PeekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) ExpectPeek(t token.TokenType) bool {
	if p.PeekTokenIs(t) {
		p.NextToken()
		return true
	} else {
		p.PeekError(t)
		return false
	}
}

func (p *Parser) RegisterPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) RegisterInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) ParseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.ParseExpression(LOWEST)

	if p.PeekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) ParseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.PeekTokenIs(token.SEMICOLON) && precedence < p.PeekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.NextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) ParseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: string(p.curToken.Literal)}
}

func (p *Parser) ParseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(string(p.curToken.Literal), 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", string(p.curToken.Literal))
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: string(p.curToken.Literal),
	}
	p.NextToken()
	expression.Right = p.ParseExpression(PREFIX)
	return expression
}

func (p *Parser) ParseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: string(p.curToken.Literal),
		Left:     left,
	}

	precedence := p.CurPrecedence()
	p.NextToken()
	expression.Right = p.ParseExpression(precedence)

	return expression
}

func (p *Parser) ParseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.CurTokenIs(token.TRUE)}

}

func (p *Parser) ParseGroupedExpression() ast.Expression {
	p.NextToken()

	exp := p.ParseExpression(LOWEST)

	if !p.ExpectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) ParseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.ExpectPeek(token.LPAREN) {
		return nil
	}

	p.NextToken()
	expression.Condition = p.ParseExpression(LOWEST)

	if !p.ExpectPeek(token.RPAREN) {
		return nil
	}

	if !p.ExpectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.ParseBlockStatement()

	if p.PeekTokenIs(token.ELSE) {
		p.NextToken()

		if !p.ExpectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.ParseBlockStatement()
	}

	return expression
}

func (p *Parser) ParseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.NextToken()

	for !p.CurTokenIs(token.RBRACE) && !p.CurTokenIs(token.EOF) {
		stmt := p.ParseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.NextToken()
	}

	return block
}

func (p *Parser) ParseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.ExpectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.ParseFunctionParameters()

	if !p.ExpectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.ParseBlockStatement()

	return lit
}

func (p *Parser) ParseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.PeekTokenIs(token.RPAREN) {
		p.NextToken()
		return identifiers
	}

	p.NextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: string(p.curToken.Literal)}
	identifiers = append(identifiers, ident)

	for p.PeekTokenIs(token.COMMA) {
		p.NextToken()
		p.NextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: string(p.curToken.Literal)}
		identifiers = append(identifiers, ident)
	}

	if !p.ExpectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) ParseCallFunction(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.ParseCallArguments()
	return exp
}

func (p *Parser) ParseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.PeekTokenIs(token.RPAREN) {
		p.NextToken()
		return args
	}

	p.NextToken()
	args = append(args, p.ParseExpression(LOWEST))

	for p.PeekTokenIs(token.COMMA) {
		p.NextToken()
		p.NextToken()
		args = append(args, p.ParseExpression(LOWEST))
	}

	if !p.ExpectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) ParseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: string(p.curToken.Literal)}
}
