// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tokens

import (
	"fmt"
)

// Type is the set of lexical token types of the Gosh programming language.
type Type string

func (t Type) String() string {
	return string(t)
}

// The list of token types.
const (
	Illegal Type = "ILLEGAL"

	EOF     Type = "EOF"
	Comment Type = "COMMENT"

	Identifier Type = "IDENTIFIER"
	Integer    Type = "INTEGER"
	// Float      = "FLOAT"
	// Rune       = "RUNE"
	String Type = "STRING"

	Assignment Type = "ASSIGNMENT" // =
	Define     Type = "DEFINE"     // :=

	Sum        Type = "SUM"        // +
	Difference Type = "DIFFERENCE" // -
	Product    Type = "PRODUCT"    // *
	Quotient   Type = "QUOTIENT"   // /
	Remainder  Type = "REMAINDER"  // %

	SumAssignment        Type = "SUM_ASSIGNMENT"        // +=
	DifferenceAssignment Type = "DIFFERENCE_ASSIGNMENT" // -=
	ProductAssignment    Type = "PRODUCT_ASSIGNMENT"    // *=
	QuotientAssignment   Type = "QUOTIENT_ASSIGNMENT"   // /=
	RemainderAssignment  Type = "REMAINDER_ASSIGNMENT"  // %=

	Increment Type = "INCREMENT" // ++
	Decrement Type = "DECREMENT" // --

	BitwiseAnd Type = "BITWISE_AND" // &
	BitwiseOr  Type = "BITWISE_OR"  // |
	BitwiseXor Type = "BITWISE_XOR" // ^
	// TODO &^

	// TODO &=
	// TODO |=
	// TODO ^=
	// TODO &^=

	// TODO <<
	// TODO >>

	// TODO <<=
	// TODO >>=

	LogicalAnd Type = "LOGICAL_AND" // &&
	LogicalOr  Type = "LOGICAL_OR"  // ||

	Not Type = "NOT" // !

	// Ellipsis // ...

	Equal          Type = "EQUAL"            // ==
	NotEqual       Type = "NOT_EQUAL"        // !=
	Less           Type = "LESS"             // <
	LessOrEqual    Type = "LESS_OR_EQUAL"    // <=
	Greater        Type = "GREATER"          // >
	GreaterOrEqual Type = "GREATER_OR_EQUAL" // >=

	// TODO <-

	// delimiters
	Colon     Type = "COLON"     // :
	Semicolon Type = "SEMICOLON" // ;
	Comma     Type = "COMMA"     // ,
	Period    Type = "PERIOD"    // .

	// TODO rename those
	LPAREN Type = "LPAREN" // (
	RPAREN Type = "RPAREN" // )
	LBRACE Type = "LBRACE" // {
	RBRACE Type = "RBRACE" // }

	// keywords
	Break       Type = "BREAK"
	Case        Type = "CASE"
	Chan        Type = "CHAN"
	Const       Type = "CONST"
	Continue    Type = "CONTINUE"
	Default     Type = "DEFAULT"
	Defer       Type = "DEFER"
	Else        Type = "ELSE"
	Fallthrough Type = "FALLTHROUGH"
	For         Type = "FOR"
	Func        Type = "FUNC"
	Go          Type = "GO"
	Goto        Type = "GOTO"
	If          Type = "IF"
	Import      Type = "IMPORT"
	Interface   Type = "INTERFACE"
	Map         Type = "MAP"
	Package     Type = "PACKAGE"
	Range       Type = "RANGE"
	Return      Type = "RETURN"
	Select      Type = "SELECT"
	Struct      Type = "STRUCT"
	Switch      Type = "SWITCH"
	// TODO Type
	Var Type = "VAR"

	// TODO remove
	True  Type = "TRUE"
	False Type = "FALSE"

	// notwithstanding
	// thetruthofthematter
	// despiteallobjections
	// whereas
	// insofaras
)

// check interfaces
var (
	_ fmt.Stringer = Type("FOO")
)
