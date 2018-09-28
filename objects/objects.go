// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package objects

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gosh-lang.org/gosh/ast"
)

// Object is a common interface for all Gosh runtime objects.
type Object interface {
	// Type returns object's type.
	Type() Type
	fmt.Stringer
}

// Integer represents integer runtime object.
type Integer struct {
	Value int
}

// Type returns INTEGER.
func (i *Integer) Type() Type { return INTEGER }

func (i *Integer) String() string { return strconv.Itoa(i.Value) }

// Boolean represents boolean runtime object.
type Boolean struct {
	Value bool
}

// Type returns BOOLEAN.
func (b *Boolean) Type() Type { return BOOLEAN }

func (b *Boolean) String() string { return strconv.FormatBool(b.Value) }

// String represents string runtime object.
type String struct {
	Value string
}

// Type returns STRING.
func (s *String) Type() Type { return STRING }

func (s *String) String() string { return s.Value }

// Continue represents continue runtime object.
type Continue struct{}

// Type returns CONTINUE.
func (c *Continue) Type() Type { return CONTINUE }

func (c *Continue) String() string {
	return "continue"
}

// Function represents function runtime object.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Scope      *Scope
}

// Type returns FUNCTION.
func (f *Function) Type() Type { return FUNCTION }

func (f *Function) String() string {
	params := make([]string, len(f.Parameters))
	for i, p := range f.Parameters {
		params[i] = p.String()
	}

	var res strings.Builder
	res.WriteString("func(")
	res.WriteString(strings.Join(params, ", "))
	res.WriteString(") ")
	res.WriteString(f.Body.String())
	return res.String()
}

// GoFunction represents Go function.
type GoFunction struct {
	Func func(args ...Object) Object
}

// Type returns GOFUNCTION.
func (gf *GoFunction) Type() Type { return GOFUNCTION }

func (gf *GoFunction) String() string { return reflect.ValueOf(gf.Func).String() }

// check interfaces
var (
	_ Object = (*Integer)(nil)
	_ Object = (*Boolean)(nil)
	_ Object = (*Function)(nil)
	_ Object = (*GoFunction)(nil)
)
