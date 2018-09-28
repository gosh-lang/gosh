// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package interpreter

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gosh-lang.org/gosh/internal/gofuzz"
	"gosh-lang.org/gosh/internal/golden"
	"gosh-lang.org/gosh/objects"
	"gosh-lang.org/gosh/parser"
	"gosh-lang.org/gosh/scanner"
)

func TestGolden(t *testing.T) {
	for _, f := range golden.Data {
		t.Run(f.File, func(t *testing.T) {
			s, err := scanner.New(f.Source, &scanner.Config{
				SkipShebang: true,
			})
			require.NoError(t, err)
			p := parser.New(s, nil)
			program := p.ParseProgram()
			var buf bytes.Buffer
			i := New(nil)
			res := i.Eval(context.Background(), program, objects.NewScope(objects.Builtin(&buf)))
			t.Log(res)
			assert.Equal(t, f.Output, strings.Split(buf.String(), "\n"))
		})
	}
}

func eval(t *testing.T, input string) (objects.Object, *bytes.Buffer) {
	t.Helper()

	s, err := scanner.New(input, nil)
	require.NoError(t, err)

	p := parser.New(s, nil)
	program := p.ParseProgram()
	require.Nil(t, p.Errors(), "%s", p.Errors())
	require.NotNil(t, program)

	i := New(nil)
	var buf bytes.Buffer
	res := i.Eval(context.Background(), program, objects.NewScope(objects.Builtin(&buf)))
	return res, &buf
}

func TestInfixExpression(t *testing.T) {
	for input, expected := range map[string]bool{
		`7 < 42`:  true,
		`7 <= 42`: true,
		`7 > 42`:  false,
		`7 >= 42`: false,
		`7 == 42`: false,
		`7 != 42`: true,
	} {
		t.Run(input, func(t *testing.T) {
			gofuzz.AddDataToCorpus("interpreter", []byte(input))

			actual, res := eval(t, input)
			require.IsType(t, &objects.Boolean{}, actual)
			assert.Equal(t, expected, actual.(*objects.Boolean).Value)
			assert.Empty(t, res.String())
		})
	}

	for input, expected := range map[string]int{
		`42 + 7`: 49,
		`42 - 7`: 35,
		`42 * 7`: 294,
		`42 / 7`: 6,
		`42 % 7`: 0,
	} {
		t.Run(input, func(t *testing.T) {
			gofuzz.AddDataToCorpus("interpreter", []byte(input))

			actual, res := eval(t, input)
			require.IsType(t, &objects.Integer{}, actual)
			assert.Equal(t, expected, actual.(*objects.Integer).Value)
			assert.Empty(t, res.String())
		})
	}
}

func TestLen(t *testing.T) {
	for input, expected := range map[string]int{
		`len("FizzBuzz")`: 8,
	} {
		t.Run(input, func(t *testing.T) {
			gofuzz.AddDataToCorpus("interpreter", []byte(input))

			actual, res := eval(t, input)
			require.IsType(t, &objects.Integer{}, actual)
			assert.Equal(t, expected, actual.(*objects.Integer).Value)
			assert.Empty(t, res.String())
		})
	}
}

func TestIf(t *testing.T) {
	for input, output := range map[string]string{
		`if (true) { print(true) }`:          "true",
		`if (false && true) { print(true) }`: "",
		`if (true && false) { print(true) }`: "",
	} {
		t.Run(input, func(t *testing.T) {
			gofuzz.AddDataToCorpus("interpreter", []byte(input))

			res, buf := eval(t, input)
			assert.Nil(t, res)
			assert.Equal(t, output, buf.String())
		})
	}
}
