// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package interpreter

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"gosh-lang.org/gosh/objects"
	"gosh-lang.org/gosh/parser"
	"gosh-lang.org/gosh/scanner"
)

var sink interface{}

func BenchmarkEval(b *testing.B) {
	input := `
	var i = 1
	for i = 1; i <= 100; i++ {
		var m3 = (i%3 == 0)
		var m5 = (i%5 == 0)

		if (m3 && m5) {
			println("FizzBuzz")
			continue
		}
		if (m3) {
			println("Fizz")
			continue
		}
		if (m5) {
			println("Buzz")
			continue
		}
		println(i)
	}`

	s, err := scanner.New(input, nil)
	require.NoError(b, err)

	p := parser.New(s, nil)
	program := p.ParseProgram()
	require.Nil(b, p.Errors())
	require.NotNil(b, program)

	i := New(nil)
	scope := objects.Builtin(ioutil.Discard)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sink = i.Eval(context.Background(), program, scope)
	}
}
