// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build gofuzz

package interpreter

import (
	"context"
	"io/ioutil"

	"gosh-lang.org/gosh/objects"
	"gosh-lang.org/gosh/parser"
	"gosh-lang.org/gosh/scanner"
)

func Fuzz(data []byte) int {
	s, err := scanner.New(string(data), &scanner.Config{
		SkipShebang: true,
	})
	if err != nil {
		return 0
	}

	p := parser.New(s, nil)
	program := p.ParseProgram()
	i := New(nil)
	i.Eval(context.TODO(), program, objects.NewScope(objects.Builtin(ioutil.Discard)))
	return 0
}
