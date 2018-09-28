// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast

import (
	"strings"

	"gosh-lang.org/gosh/tokens"
)

// Statement is a common interface for all AST statement nodes.
type Statement interface {
	Node
	statement()
}

// IncrementDecrementStatement represents increment or decrement statement (e.g. `x++`, `x--`).
type IncrementDecrementStatement struct {
	Token tokens.Token // tokens.Increment or tokens.Decrement
	Name  *Identifier  // TODO it can be a more complex expression
}

func (ids *IncrementDecrementStatement) String() string {
	var res strings.Builder
	res.WriteString(ids.Name.String())
	res.WriteString(ids.Token.Literal)
	return res.String()
}

func (ids *IncrementDecrementStatement) node()      {}
func (ids *IncrementDecrementStatement) statement() {}

// VarStatement represents a var statement.
type VarStatement struct {
	Token tokens.Token // tokens.Var
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) String() string {
	var res strings.Builder
	res.WriteString("var ")
	res.WriteString(vs.Name.String())
	if vs.Value != nil {
		res.WriteString(" = ")
		res.WriteString(vs.Value.String())
	}
	return res.String()
}

func (vs *VarStatement) node()      {}
func (vs *VarStatement) statement() {}

// AssignStatement represents an assign statement.
type AssignStatement struct {
	Token tokens.Token // tokens.Assignment or tokens.XXXAssignment
	Name  *Identifier  // TODO it can be a more complex expression
	Value Expression
}

func (as *AssignStatement) String() string {
	var res strings.Builder
	res.WriteString(as.Name.String())
	res.WriteString(" ")
	res.WriteString(as.Token.Literal)
	res.WriteString(" ")
	res.WriteString(as.Value.String())
	return res.String()
}

func (as *AssignStatement) node()      {}
func (as *AssignStatement) statement() {}

// ReturnStatement represents a return statement.
type ReturnStatement struct {
	Token tokens.Token // tokens.Return
	Value Expression
}

func (rs *ReturnStatement) String() string {
	var res strings.Builder
	res.WriteString("return")
	if rs.Value != nil {
		res.WriteString(" ")
		res.WriteString(rs.Value.String())
	}
	return res.String()
}

func (rs *ReturnStatement) node()      {}
func (rs *ReturnStatement) statement() {}

// ContinueStatement represents a continue statement.
type ContinueStatement struct {
	Token tokens.Token // tokens.Continue
}

func (cs *ContinueStatement) String() string {
	return "continue"
}

func (cs *ContinueStatement) node()      {}
func (cs *ContinueStatement) statement() {}

// IfStatement represent if/else statement.
type IfStatement struct {
	Token tokens.Token // tokens.If
	// Init  *AssignStatement // initialization statement; or nil // TODO it also can be a define statement
	Cond Expression // condition; or nil
	Body *BlockStatement
	// Else Statement // else branch; or nil
}

func (is *IfStatement) String() string {
	var res strings.Builder
	res.WriteString("if (")
	res.WriteString(is.Cond.String())
	res.WriteString(") ")
	res.WriteString(is.Body.String())
	return res.String()
}

func (is *IfStatement) node()      {}
func (is *IfStatement) statement() {}

// ForStatement represent a for statement.
type ForStatement struct {
	Token tokens.Token     // tokens.For
	Init  *AssignStatement // initialization statement; or nil // TODO it also can be a define statement
	Cond  Expression       // condition; or nil
	Post  Statement        // post iteration statement; or nil
	Body  *BlockStatement
}

func (fs *ForStatement) String() string {
	var res strings.Builder
	res.WriteString("for ")
	if fs.Init != nil {
		res.WriteString(fs.Init.String())
	}
	res.WriteString("; ")
	if fs.Cond != nil {
		res.WriteString(fs.Cond.String())
	}
	res.WriteString("; ")
	if fs.Post != nil {
		res.WriteString(fs.Post.String())
	}
	res.WriteString(" ")
	if fs.Body != nil {
		res.WriteString(fs.Body.String())
	}
	return res.String()
}

func (fs *ForStatement) node()      {}
func (fs *ForStatement) statement() {}

// ExpressionStatement represents an expression when it is used as a statement.
type ExpressionStatement struct {
	Token      tokens.Token // first token of expression
	Expression Expression
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (es *ExpressionStatement) node()      {}
func (es *ExpressionStatement) statement() {}

// BlockStatement represents a block statement.
type BlockStatement struct {
	Token      tokens.Token // tokens.LBRACE
	Statements []Statement
}

func (bs *BlockStatement) String() string {
	var res strings.Builder
	res.WriteString("{\n")
	for _, s := range bs.Statements {
		res.WriteString(s.String() + ";\n")
	}
	res.WriteString("}")
	return res.String()
}

func (bs *BlockStatement) node()      {}
func (bs *BlockStatement) statement() {}

// check interfaces
var (
	_ Statement = (*IncrementDecrementStatement)(nil)
	_ Statement = (*VarStatement)(nil)
	_ Statement = (*AssignStatement)(nil)
	_ Statement = (*ReturnStatement)(nil)
	_ Statement = (*ForStatement)(nil)
	_ Statement = (*ExpressionStatement)(nil)
	_ Statement = (*BlockStatement)(nil)
)
