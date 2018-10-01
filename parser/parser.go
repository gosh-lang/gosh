// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// package parser // import "gosh-lang.org/gosh/parser" implements a parser for Gosh source files.
package parser // import "gosh-lang.org/gosh/parser"

import (
	"fmt"
	"strconv"
	"strings"

	"gosh-lang.org/gosh/ast"
	"gosh-lang.org/gosh/scanner"
	"gosh-lang.org/gosh/tokens"
)

// Parser implements parsing of Gosh source files.
type Parser struct {
	s      *scanner.Scanner
	config *Config
	errors []error

	curToken  tokens.Token
	peekToken tokens.Token

	prefixParseFns map[tokens.Type]prefixParseFn
	infixParseFns  map[tokens.Type]infixParseFn
}

// Config configures parser.
type Config struct {
	crashOnError bool // crash parser on any error, for testing only
}

// New creates a new parser.
func New(s *scanner.Scanner, config *Config) *Parser {
	if config == nil {
		config = new(Config)
	}

	p := &Parser{
		s:              s,
		config:         config,
		prefixParseFns: make(map[tokens.Type]prefixParseFn),
		infixParseFns:  make(map[tokens.Type]infixParseFn),
	}

	// groped just like tokens.Type constants

	for t, f := range map[tokens.Type]prefixParseFn{
		tokens.Comment: p.parseComment, // TODO really?!

		tokens.Integer:    p.parseIntegerLiteral,
		tokens.String:     p.parseStringLiteral,
		tokens.Identifier: p.parseIdentifier,

		tokens.Difference: p.parsePrefixExpression,

		tokens.Not: p.parsePrefixExpression,

		tokens.LPAREN: p.parseGroupedExpression,

		tokens.Func: p.parseFunctionLiteral,

		// TODO remove
		tokens.True:  p.parseBooleanLiteral,
		tokens.False: p.parseBooleanLiteral,
	} {
		p.registerPrefix(t, f)
	}

	for t, f := range map[tokens.Type]infixParseFn{
		tokens.Sum:        p.parseInfixExpression,
		tokens.Difference: p.parseInfixExpression,
		tokens.Product:    p.parseInfixExpression,
		tokens.Quotient:   p.parseInfixExpression,
		tokens.Remainder:  p.parseInfixExpression,

		// TODO AND, OR, XOR

		tokens.LogicalAnd: p.parseInfixExpression,
		tokens.LogicalOr:  p.parseInfixExpression,

		tokens.Equal:          p.parseInfixExpression,
		tokens.NotEqual:       p.parseInfixExpression,
		tokens.Less:           p.parseInfixExpression,
		tokens.LessOrEqual:    p.parseInfixExpression,
		tokens.Greater:        p.parseInfixExpression,
		tokens.GreaterOrEqual: p.parseInfixExpression,

		tokens.LPAREN: p.parseCallExpression,
	} {
		p.registerInfix(t, f)
	}

	// set both curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// List of precedences.
const (
	LowestPrec  = 0
	UnaryPrec   = 6 // -X, !X
	HighestPrec = 7 // foo(X)
)

var precedences = map[tokens.Type]int{
	// TODO make LowestPrec default, remove those
	tokens.Illegal:              LowestPrec,
	tokens.EOF:                  LowestPrec,
	tokens.Comment:              LowestPrec,
	tokens.LBRACE:               LowestPrec,
	tokens.RBRACE:               LowestPrec,
	tokens.RPAREN:               LowestPrec,
	tokens.Colon:                LowestPrec,
	tokens.Identifier:           LowestPrec,
	tokens.Assignment:           LowestPrec,
	tokens.Define:               LowestPrec,
	tokens.If:                   LowestPrec,
	tokens.Else:                 LowestPrec,
	tokens.Return:               LowestPrec,
	tokens.True:                 LowestPrec,
	tokens.False:                LowestPrec,
	tokens.Func:                 LowestPrec,
	tokens.Integer:              LowestPrec,
	tokens.Increment:            LowestPrec,
	tokens.Switch:               LowestPrec,
	tokens.Case:                 LowestPrec,
	tokens.Var:                  LowestPrec,
	tokens.For:                  LowestPrec,
	tokens.SumAssignment:        LowestPrec,
	tokens.DifferenceAssignment: LowestPrec,
	tokens.ProductAssignment:    LowestPrec,
	tokens.QuotientAssignment:   LowestPrec,
	tokens.RemainderAssignment:  LowestPrec,

	tokens.LogicalOr: 1,

	tokens.LogicalAnd: 2,

	tokens.Equal:          3,
	tokens.NotEqual:       3,
	tokens.Less:           3,
	tokens.LessOrEqual:    3,
	tokens.Greater:        3,
	tokens.GreaterOrEqual: 3,

	tokens.Sum:        4,
	tokens.Difference: 4,
	tokens.BitwiseOr:  4,
	tokens.BitwiseXor: 4,

	tokens.Product:    5,
	tokens.Quotient:   5,
	tokens.Remainder:  5,
	tokens.BitwiseAnd: 5,

	tokens.Not: UnaryPrec,

	tokens.LPAREN: HighestPrec,
}

func (p *Parser) crash(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	panic(fmt.Errorf("%s\ncurToken: %s\npeekToken: %s\nerrors: %s", msg, p.curToken, p.peekToken, p.errors))
}

func (p *Parser) peekPrecedence() int {
	t := p.peekToken.Type
	if p, ok := precedences[t]; ok {
		return p
	}

	p.crash("precedence for %s not found", t) // TODO remove
	return LowestPrec
}

func (p *Parser) curPrecedence() int {
	t := p.curToken.Type
	if p, ok := precedences[t]; ok {
		return p
	}

	p.crash("precedence for %s not found", t) // TODO remove
	return LowestPrec
}

func (p *Parser) addParsingError(format string, a ...interface{}) {
	err := fmt.Sprintf(format, a...)
	p.errors = append(p.errors, &Error{Err: err})

	if p.config.crashOnError {
		p.crash(format, a...)
	}
}

// Errors returns parsing errors, if any.
func (p *Parser) Errors() []error {
	return p.errors
}

// nextToken advances curToken and peekToken.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.s.NextToken()
}

func (p *Parser) expectCurrent(tt ...tokens.Type) bool {
	if len(tt) == 0 {
		p.crash("expectCurrent called with zero token types")
	}

	for _, t := range tt {
		if p.curToken.Type == t {
			return true
		}
	}

	switch l := len(tt); l {
	case 1:
		p.addParsingError("expected current token to be %s, got %s instead", tt[0], p.curToken.Type)
	default:
		expected := make([]string, l)
		for i, t := range tt {
			expected[i] = t.String()
		}
		exp := strings.Join(expected, ", ")
		p.addParsingError("expected current token to be one of %s, got %s instead", exp, p.curToken.Type)
	}

	return false

}

func (p *Parser) expectPeek(tt ...tokens.Type) bool {
	if len(tt) == 0 {
		p.crash("expectPeek called with zero token types")
	}

	for _, t := range tt {
		if p.peekToken.Type == t {
			p.nextToken()
			return true
		}
	}

	switch l := len(tt); l {
	case 1:
		p.addParsingError("expected next token to be %s, got %s instead", tt[0], p.peekToken)
	default:
		expected := make([]string, l)
		for i, t := range tt {
			expected[i] = t.String()
		}
		exp := strings.Join(expected, ", ")
		p.addParsingError("expected next token to be one of %s, got %s instead", exp, p.peekToken)
	}

	return false
}

func (p *Parser) registerPrefix(tokenType tokens.Type, fn prefixParseFn) {
	if _, ok := p.prefixParseFns[tokenType]; ok {
		p.crash("prefix function for %s already registered", tokenType)
	}
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType tokens.Type, fn infixParseFn) {
	if _, ok := p.infixParseFns[tokenType]; ok {
		p.crash("infix function for %s already registered", tokenType)
	}
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseComment() ast.Expression {
	// TODO
	return nil
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addParsingError("could not parse %q as integer", p.curToken.Literal)
		return nil
	}

	lit.Value = int(value)
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	s := p.curToken.Literal
	if !strings.HasPrefix(s, `"`) || !strings.HasSuffix(s, `"`) {
		p.addParsingError("could not parse %q as string", s)
		return nil
	}
	return &ast.StringLiteral{Token: p.curToken, Value: s[1 : len(s)-1]}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	t := p.curToken.Type == tokens.True
	return &ast.BooleanLiteral{
		Token: p.curToken,
		Value: t,
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.addParsingError("no prefix parse function for %s found (token %s)", p.curToken.Type, p.curToken)
		return nil
	}
	leftExp := prefix()

	for p.peekToken.Type != tokens.Semicolon && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token: p.curToken,
	}

	p.nextToken()
	expression.Right = p.parseExpression(UnaryPrec)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LowestPrec)
	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	if !p.expectCurrent(tokens.LBRACE) {
		return nil
	}
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = make([]ast.Statement, 0, 1)

	p.nextToken()

	for p.curToken.Type != tokens.RBRACE && p.curToken.Type != tokens.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(tokens.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	if p.peekToken.Type == tokens.RPAREN {
		p.nextToken()
		return []*ast.Identifier{}
	}

	p.nextToken()

	identifiers := []*ast.Identifier{{Token: p.curToken, Value: p.curToken.Literal}}

	for p.peekToken.Type == tokens.Comma {
		p.nextToken()
		p.nextToken()
		identifiers = append(identifiers, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token: p.curToken,
		Left:  left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekToken.Type == tokens.RPAREN {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LowestPrec))

	for p.peekToken.Type == tokens.Comma {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LowestPrec))
	}

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}
	if !p.expectPeek(tokens.Identifier) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(tokens.Assignment) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LowestPrec)

	for p.peekToken.Type == tokens.Semicolon {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	if !p.expectCurrent(tokens.If) {
		return nil
	}
	stmt := &ast.IfStatement{Token: p.curToken}

	if !p.expectPeek(tokens.LPAREN) {
		return nil
	}
	p.nextToken()
	stmt.Cond = p.parseExpression(LowestPrec)

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}
	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()

	for p.peekToken.Type == tokens.Semicolon {
		p.nextToken()
	}

	return stmt
}

var assignTokens = []tokens.Type{
	tokens.Assignment,
	tokens.SumAssignment,
	tokens.DifferenceAssignment,
	tokens.ProductAssignment,
	tokens.QuotientAssignment,
	tokens.RemainderAssignment,
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	if !p.expectCurrent(tokens.Identifier) {
		return nil
	}
	stmt := &ast.AssignStatement{
		Name: &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		},
	}

	if !p.expectPeek(assignTokens...) {
		return nil
	}
	stmt.Token = p.curToken

	p.nextToken()
	stmt.Value = p.parseExpression(LowestPrec)

	for p.peekToken.Type == tokens.Semicolon {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseIncrementDecrementStatement() *ast.IncrementDecrementStatement {
	if !p.expectCurrent(tokens.Identifier) {
		return nil
	}
	stmt := &ast.IncrementDecrementStatement{
		Name: &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		},
	}

	if !p.expectPeek(tokens.Increment, tokens.Decrement) {
		return nil
	}
	stmt.Token = p.curToken

	for p.peekToken.Type == tokens.Semicolon {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.Value = p.parseExpression(LowestPrec)

	for p.peekToken.Type == tokens.Semicolon {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	if !p.expectCurrent(tokens.Continue) {
		return nil
	}
	stmt := &ast.ContinueStatement{
		Token: p.curToken,
	}

	for p.peekToken.Type == tokens.Semicolon {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	if !p.expectCurrent(tokens.For) {
		return nil
	}
	stmt := &ast.ForStatement{Token: p.curToken}

	p.nextToken()
	stmt.Init = p.parseAssignStatement()

	if !p.expectCurrent(tokens.Semicolon) {
		return nil
	}
	p.nextToken()
	stmt.Cond = p.parseExpression(LowestPrec)

	if !p.expectPeek(tokens.Semicolon) {
		return nil
	}
	p.nextToken()
	stmt.Post = p.parseStatement()

	p.nextToken()
	stmt.Body = p.parseBlockStatement()

	for p.peekToken.Type == tokens.Semicolon {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionOrAssignmentStatement() ast.Statement {
	cur := p.curToken
	exp := p.parseExpression(LowestPrec)

	var stmt ast.Statement
	switch p.peekToken.Type {
	case tokens.Increment, tokens.Decrement:
		stmt = p.parseIncrementDecrementStatement()
	default:
		for _, t := range assignTokens {
			if p.peekToken.Type == t {
				stmt = p.parseAssignStatement()
				break
			}
		}
		if stmt == nil {
			stmt = &ast.ExpressionStatement{Token: cur, Expression: exp}
		}
	}

	for p.peekToken.Type == tokens.Semicolon {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case tokens.Var:
		return p.parseVarStatement()
	case tokens.If:
		return p.parseIfStatement()
	case tokens.Return:
		return p.parseReturnStatement()
	case tokens.Continue:
		return p.parseContinueStatement()
	case tokens.For:
		return p.parseForStatement()
	default:
		return p.parseExpressionOrAssignmentStatement()
	}
}

// ParseProgram parsers the whole program and returns root AST node, or nil, if error are encountered.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: make([]ast.Statement, 0, 8),
	}

	for p.curToken.Type != tokens.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	// TODO not sure about it
	if len(p.Errors()) != 0 {
		program = nil
	}

	return program
}
