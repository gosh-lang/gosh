// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package ast declares the types used to represent syntax trees for Gosh packages.
package ast // import "gosh-lang.org/gosh/ast"

import (
	// import "gosh-lang.org/gosh/ast"
	"fmt"
	"strings"
)

// Node is a common interface for all AST nodes.
type Node interface {
	fmt.Stringer
	node()
}

// Program is a root of AST tree.
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var res strings.Builder
	for _, s := range p.Statements {
		res.WriteString(s.String())
		res.WriteString(";\n")
	}
	return res.String()
}

func (p *Program) node() {}

// check interfaces
var (
	_ Node = (*Program)(nil)
)
