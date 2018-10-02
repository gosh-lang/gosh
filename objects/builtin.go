// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package objects // import "gosh-lang.org/gosh/objects"

import (
	"fmt"
	"io"
	"strings"
)

var (
	lenBuiltin = &GoFunction{Func: func(args ...Object) Object {
		if len(args) != 1 {
			panic(fmt.Errorf("len: expected 1 argument, got %d", len(args)))
		}
		arg := args[0]
		switch arg := arg.(type) {
		case *String:
			return &Integer{
				Value: len(arg.Value),
			}
		default:
			panic(fmt.Errorf("len: unexpected argument type %T", arg))
		}
	}}

	// TODO append
	// TODO cap
	// TODO close
	// TODO copy
	// TODO delete
	// TODO make
	// TODO new
	// TODO panic
	// TODO recover
)

func makePrintBuiltin(stdout io.Writer) *GoFunction {
	return &GoFunction{Func: func(args ...Object) Object {
		res := make([]string, len(args))
		for i, arg := range args {
			res[i] = arg.String()
		}
		fmt.Fprint(stdout, strings.Join(res, " "))
		return nil
	}}
}

func makePrintlnBuiltin(stdout io.Writer) *GoFunction {
	return &GoFunction{Func: func(args ...Object) Object {
		res := make([]string, len(args))
		for i, arg := range args {
			res[i] = arg.String()
		}
		fmt.Fprintln(stdout, strings.Join(res, " "))
		return nil
	}}
}

// Builtin returns a Scope of predeclared identifiers.
func Builtin(stdout io.Writer) *Scope {
	return &Scope{
		store: map[string]Object{
			"print":   makePrintBuiltin(stdout),
			"println": makePrintlnBuiltin(stdout),
			"len":     lenBuiltin,
		},
	}
}
