// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast // import "gosh-lang.org/gosh/ast"

import (
	"strings"

	"gosh-lang.org/gosh/tokens"
)

// Expression is a common interface for all AST expression nodes.
type Expression interface {
	Node
	expression()
}

// Identifier represents an identifier expression.
type Identifier struct {
	Token tokens.Token // tokens.IDENT
	Value string
}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) node()       {}
func (i *Identifier) expression() {}

// IntegerLiteral represents an integer literal expression.
type IntegerLiteral struct {
	Token tokens.Token // tokens.INT
	Value int
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) node()       {}
func (il *IntegerLiteral) expression() {}

// StringLiteral represents a string literal expression.
type StringLiteral struct {
	Token tokens.Token // tokens.INT
	Value string
}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) node()       {}
func (sl *StringLiteral) expression() {}

// BooleanLiteral represents a boolean literal expression.
type BooleanLiteral struct {
	Token tokens.Token // tokens.TRUE or tokens.FALSE
	Value bool
}

func (bl *BooleanLiteral) String() string {
	return bl.Token.Literal
}

func (bl *BooleanLiteral) node()       {}
func (bl *BooleanLiteral) expression() {}

// PrefixExpression represents prefix expression (e.g. `!x`).
type PrefixExpression struct {
	Token tokens.Token // tokens.NOT, tokens.SUB
	Right Expression
}

func (pe *PrefixExpression) String() string {
	var res strings.Builder
	res.WriteString("(")
	res.WriteString(pe.Token.Literal)
	res.WriteString(pe.Right.String())
	res.WriteString(")")
	return res.String()
}

func (pe *PrefixExpression) node()       {}
func (pe *PrefixExpression) expression() {}

// InfixExpression represents infix expression (e.g. `x + y`).
type InfixExpression struct {
	Token tokens.Token // tokens.ADD, tokens.SUB, etc.
	Left  Expression
	Right Expression
}

func (ie *InfixExpression) String() string {
	var res strings.Builder
	// res.WriteString("(")
	res.WriteString(ie.Left.String())
	res.WriteString(" ")
	res.WriteString(ie.Token.Literal)
	res.WriteString(" ")
	res.WriteString(ie.Right.String())
	// res.WriteString(")")
	return res.String()
}

func (ie *InfixExpression) node()       {}
func (ie *InfixExpression) expression() {}

// FunctionLiteral represents a function literal expression.
type FunctionLiteral struct {
	Token      tokens.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) String() string {
	params := make([]string, len(fl.Parameters))
	for i, p := range fl.Parameters {
		params[i] = p.String()
	}

	var res strings.Builder
	res.WriteString("func(")
	res.WriteString(strings.Join(params, ", "))
	res.WriteString(") ")
	res.WriteString(fl.Body.String())
	return res.String()
}

func (fl *FunctionLiteral) node()       {}
func (fl *FunctionLiteral) expression() {}

// CallExpression represents a call expression.
type CallExpression struct {
	Token     tokens.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) String() string {
	args := make([]string, len(ce.Arguments))
	for i, a := range ce.Arguments {
		args[i] = a.String()
	}

	var res strings.Builder
	res.WriteString(ce.Function.String())
	res.WriteString("(")
	res.WriteString(strings.Join(args, ", "))
	res.WriteString(")")
	return res.String()
}

func (ce *CallExpression) node()       {}
func (ce *CallExpression) expression() {}

// check interfaces
var (
	_ Expression = (*Identifier)(nil)
	_ Expression = (*IntegerLiteral)(nil)
	_ Expression = (*BooleanLiteral)(nil)
	_ Expression = (*PrefixExpression)(nil)
	_ Expression = (*InfixExpression)(nil)
	_ Expression = (*FunctionLiteral)(nil)
	_ Expression = (*CallExpression)(nil)
)
