// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"gosh-lang.org/gosh/interpreter"
	"gosh-lang.org/gosh/objects"
	"gosh-lang.org/gosh/parser"
	"gosh-lang.org/gosh/scanner"
)

func Example() {
	code := `println("Hello, world!")`

	s, err := scanner.New(code, nil)
	if err != nil {
		log.Fatal(err)
	}

	p := parser.New(s, nil)
	program := p.ParseProgram()
	if p.Errors() != nil {
		log.Fatal(p.Errors())
	}

	scope := objects.NewScope(objects.Builtin(os.Stdout))
	i := interpreter.New(nil)
	res := i.Eval(context.Background(), program, scope)
	fmt.Println("Eval result:", res)
	// Output:
	// Hello, world!
	// Eval result: <nil>
}
