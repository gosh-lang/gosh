// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package scanner // import "gosh-lang.org/gosh/scanner"

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gosh-lang.org/gosh/internal/gofuzz"
	"gosh-lang.org/gosh/internal/golden"
	"gosh-lang.org/gosh/tokens"
)

func TestGolden(t *testing.T) {
	for _, f := range golden.Data {
		t.Run(f.File, func(t *testing.T) {
			var actual []string
			s, err := New(f.Source, &Config{
				SkipShebang: true,
			})
			require.NoError(t, err)
			for _, tok := range s.allTokens() {
				actual = append(actual, tok.String())
			}
			assert.Equal(t, strings.Split(f.Tokens, "\n"), actual, "actual:\n%s", strings.Join(actual, "\n"))
		})
	}
}

func TestScanner(t *testing.T) {
	// in order of tokens.Type constants
	testdata := map[string][]tokens.Token{
		`#`: {
			{Type: tokens.Illegal, Literal: `#`},
		},
		`…`: {
			{Type: tokens.Illegal, Literal: `…`},
		},
		`42foo`: {
			{Type: tokens.Illegal, Literal: `42foo`},
		},
		`"Invalid`: {
			{Type: tokens.Illegal, Literal: `"Invalid`},
		},
		``: {
			{Type: tokens.EOF},
		},

		// FIXME fix code and enable those tests
		"// Comment 1\n// Comment 2": {
			{Offset: 0, Type: tokens.Comment, Literal: "// Comment 1"},
			{Offset: 13, Type: tokens.Comment, Literal: "// Comment 2"},
			{Offset: 25, Type: tokens.EOF},
		},
		"// Comment 1\n// Comment 2\n": {
			{Offset: 0, Type: tokens.Comment, Literal: "// Comment 1"},
			{Offset: 13, Type: tokens.Comment, Literal: "// Comment 2"},
			{Offset: 26, Type: tokens.EOF},
		},

		`foo FOO _ foo42`: {
			{Offset: 0, Type: tokens.Identifier, Literal: `foo`},
			{Offset: 4, Type: tokens.Identifier, Literal: `FOO`},
			{Offset: 8, Type: tokens.Identifier, Literal: `_`},
			{Offset: 10, Type: tokens.Identifier, Literal: `foo42`},
			{Offset: 15, Type: tokens.EOF},
		},
		`42 042`: {
			{Offset: 0, Type: tokens.Integer, Literal: `42`},
			{Offset: 3, Type: tokens.Integer, Literal: `042`},
			{Offset: 6, Type: tokens.EOF},
		},
		// TODO Float
		// TODO Character, Rune, Byte?
		`"Hello, world!"`: {
			{Offset: 0, Type: tokens.String, Literal: `"Hello, world!"`},
			{Offset: 15, Type: tokens.EOF},
		},

		`=:=`: {
			{Offset: 0, Type: tokens.Assignment, Literal: `=`},
			{Offset: 1, Type: tokens.Define, Literal: `:=`},
			{Offset: 3, Type: tokens.EOF},
		},

		`+-*/%`: {
			{Offset: 0, Type: tokens.Sum, Literal: `+`},
			{Offset: 1, Type: tokens.Difference, Literal: `-`},
			{Offset: 2, Type: tokens.Product, Literal: `*`},
			{Offset: 3, Type: tokens.Quotient, Literal: `/`},
			{Offset: 4, Type: tokens.Remainder, Literal: `%`},
			{Offset: 5, Type: tokens.EOF},
		},

		`+=-=*=/=%=`: {
			{Offset: 0, Type: tokens.SumAssignment, Literal: `+=`},
			{Offset: 2, Type: tokens.DifferenceAssignment, Literal: `-=`},
			{Offset: 4, Type: tokens.ProductAssignment, Literal: `*=`},
			{Offset: 6, Type: tokens.QuotientAssignment, Literal: `/=`},
			{Offset: 8, Type: tokens.RemainderAssignment, Literal: `%=`},
			{Offset: 10, Type: tokens.EOF},
		},

		`++--`: {
			{Offset: 0, Type: tokens.Increment, Literal: `++`},
			{Offset: 2, Type: tokens.Decrement, Literal: `--`},
			{Offset: 4, Type: tokens.EOF},
		},

		`&|^`: {
			{Offset: 0, Type: tokens.BitwiseAnd, Literal: `&`},
			{Offset: 1, Type: tokens.BitwiseOr, Literal: `|`},
			{Offset: 2, Type: tokens.BitwiseXor, Literal: `^`},
			{Offset: 3, Type: tokens.EOF},
		},

		`&&||`: {
			{Offset: 0, Type: tokens.LogicalAnd, Literal: `&&`},
			{Offset: 2, Type: tokens.LogicalOr, Literal: `||`},
			{Offset: 4, Type: tokens.EOF},
		},

		`!`: {
			{Offset: 0, Type: tokens.Not, Literal: `!`},
			{Offset: 1, Type: tokens.EOF},
		},

		`==!=<=<>>=`: {
			{Offset: 0, Type: tokens.Equal, Literal: `==`},
			{Offset: 2, Type: tokens.NotEqual, Literal: `!=`},
			{Offset: 4, Type: tokens.LessOrEqual, Literal: `<=`},
			{Offset: 6, Type: tokens.Less, Literal: `<`},
			{Offset: 7, Type: tokens.Greater, Literal: `>`},
			{Offset: 8, Type: tokens.GreaterOrEqual, Literal: `>=`},
			{Offset: 10, Type: tokens.EOF},
		},

		`:;,.`: {
			{Offset: 0, Type: tokens.Colon, Literal: `:`},
			{Offset: 1, Type: tokens.Semicolon, Literal: `;`},
			{Offset: 2, Type: tokens.Comma, Literal: `,`},
			{Offset: 3, Type: tokens.Period, Literal: `.`},
			{Offset: 4, Type: tokens.EOF},
		},

		`(){}`: {
			{Offset: 0, Type: tokens.LPAREN, Literal: `(`},
			{Offset: 1, Type: tokens.RPAREN, Literal: `)`},
			{Offset: 2, Type: tokens.LBRACE, Literal: `{`},
			{Offset: 3, Type: tokens.RBRACE, Literal: `}`},
			{Offset: 4, Type: tokens.EOF},
		},

		`break case chan const continue default defer else fallthrough for func go ` +
			`goto if import interface map package range return select struct switch var`: {
			{Offset: 0, Type: tokens.Break, Literal: `break`},
			{Offset: 6, Type: tokens.Case, Literal: `case`},
			{Offset: 11, Type: tokens.Chan, Literal: `chan`},
			{Offset: 16, Type: tokens.Const, Literal: `const`},
			{Offset: 22, Type: tokens.Continue, Literal: `continue`},
			{Offset: 31, Type: tokens.Default, Literal: `default`},
			{Offset: 39, Type: tokens.Defer, Literal: `defer`},
			{Offset: 45, Type: tokens.Else, Literal: `else`},
			{Offset: 50, Type: tokens.Fallthrough, Literal: `fallthrough`},
			{Offset: 62, Type: tokens.For, Literal: `for`},
			{Offset: 66, Type: tokens.Func, Literal: `func`},
			{Offset: 71, Type: tokens.Go, Literal: `go`},
			{Offset: 74, Type: tokens.Goto, Literal: `goto`},
			{Offset: 79, Type: tokens.If, Literal: `if`},
			{Offset: 82, Type: tokens.Import, Literal: `import`},
			{Offset: 89, Type: tokens.Interface, Literal: `interface`},
			{Offset: 99, Type: tokens.Map, Literal: `map`},
			{Offset: 103, Type: tokens.Package, Literal: `package`},
			{Offset: 111, Type: tokens.Range, Literal: `range`},
			{Offset: 117, Type: tokens.Return, Literal: `return`},
			{Offset: 124, Type: tokens.Select, Literal: `select`},
			{Offset: 131, Type: tokens.Struct, Literal: `struct`},
			{Offset: 138, Type: tokens.Switch, Literal: `switch`},
			{Offset: 145, Type: tokens.Var, Literal: `var`},
			{Offset: 148, Type: tokens.EOF},
		},

		`true false`: {
			{Offset: 0, Type: tokens.True, Literal: `true`},
			{Offset: 5, Type: tokens.False, Literal: `false`},
			{Offset: 10, Type: tokens.EOF},
		},
	}

	for input, tokens := range testdata {
		t.Run(input, func(t *testing.T) {
			gofuzz.AddDataToCorpus("scanner", []byte(input))

			if strings.HasPrefix(input, "// ") {
				t.Skip("FIXME broken code")
			}

			offset := -1
			for _, tok := range tokens {
				require.True(t, offset < tok.Offset, "unexpected offset for token %s", tok)
				offset = tok.Offset
			}

			l, err := New(input, &Config{
				dontInsertSemicolon: true,
			})
			require.NoError(t, err)
			assert.Equal(t, tokens, l.allTokens(), "Input: %q", input)
		})
	}
}

func TestSemicolonInsertion(t *testing.T) {
	input := strings.TrimLeft(`
var
return
break;
continue
fallthrough;

true
false;

x
x += 1
x++

foo()
func() {}
`, "\n")
	l, err := New(input, nil)
	require.NoError(t, err)
	expected := []tokens.Token{
		{Offset: 0, Type: tokens.Var, Literal: "var"},
		{Offset: 4, Type: tokens.Return, Literal: "return"},
		{Offset: 10, Type: tokens.Semicolon, Literal: "\n"},
		{Offset: 11, Type: tokens.Break, Literal: "break"},
		{Offset: 16, Type: tokens.Semicolon, Literal: ";"},
		{Offset: 18, Type: tokens.Continue, Literal: "continue"},
		{Offset: 26, Type: tokens.Semicolon, Literal: "\n"},
		{Offset: 27, Type: tokens.Fallthrough, Literal: "fallthrough"},
		{Offset: 38, Type: tokens.Semicolon, Literal: ";"},

		{Offset: 41, Type: tokens.True, Literal: "true"},
		{Offset: 45, Type: tokens.Semicolon, Literal: "\n"},
		{Offset: 46, Type: tokens.False, Literal: "false"},
		{Offset: 51, Type: tokens.Semicolon, Literal: ";"},

		{Offset: 54, Type: tokens.Identifier, Literal: "x"},
		{Offset: 55, Type: tokens.Semicolon, Literal: "\n"},
		{Offset: 56, Type: tokens.Identifier, Literal: "x"},
		{Offset: 58, Type: tokens.SumAssignment, Literal: "+="},
		{Offset: 61, Type: tokens.Integer, Literal: "1"},
		{Offset: 62, Type: tokens.Semicolon, Literal: "\n"},
		{Offset: 63, Type: tokens.Identifier, Literal: "x"},
		{Offset: 64, Type: tokens.Increment, Literal: "++"},
		{Offset: 66, Type: tokens.Semicolon, Literal: "\n"},

		{Offset: 68, Type: tokens.Identifier, Literal: "foo"},
		{Offset: 71, Type: tokens.LPAREN, Literal: "("},
		{Offset: 72, Type: tokens.RPAREN, Literal: ")"},
		{Offset: 73, Type: tokens.Semicolon, Literal: "\n"},
		{Offset: 74, Type: tokens.Func, Literal: "func"},
		{Offset: 78, Type: tokens.LPAREN, Literal: "("},
		{Offset: 79, Type: tokens.RPAREN, Literal: ")"},
		{Offset: 81, Type: tokens.LBRACE, Literal: "{"},
		{Offset: 82, Type: tokens.RBRACE, Literal: "}"},
		{Offset: 83, Type: tokens.Semicolon, Literal: "\n"},

		{Offset: 84, Type: tokens.EOF},
	}
	assert.Equal(t, expected, l.allTokens())

	t.Run("Without newline", func(t *testing.T) {
		testdata := map[string][]tokens.Token{
			`0`: {
				{Offset: 0, Type: tokens.Integer, Literal: "0"},
				{Offset: 1, Type: tokens.EOF, Literal: ""},
			},
			`(`: {
				{Offset: 0, Type: tokens.LPAREN, Literal: "("},
				{Offset: 1, Type: tokens.EOF, Literal: ""},
			},
		}

		for input, tokens := range testdata {
			t.Run(input, func(t *testing.T) {
				gofuzz.AddDataToCorpus("scanner", []byte(input))

				offset := -1
				for _, tok := range tokens {
					require.True(t, offset < tok.Offset, "unexpected offset for token %s", tok)
					offset = tok.Offset
				}

				l, err := New(input, nil)
				require.NoError(t, err)
				assert.Equal(t, tokens, l.allTokens(), "Input: %q", input)
			})
		}
	})
}

func TestShebang(t *testing.T) {
	l, err := New("#!/usr/bin/env gosh\nfoo", &Config{
		SkipShebang: false,
	})
	require.NoError(t, err)
	expected := []tokens.Token{
		{Type: tokens.Illegal, Literal: `#`},
	}
	assert.Equal(t, expected, l.allTokens())

	l, err = New("#!/usr/bin/env gosh", &Config{
		SkipShebang: false,
	})
	require.NoError(t, err)
	expected = []tokens.Token{
		{Type: tokens.Illegal, Literal: `#`},
	}
	assert.Equal(t, expected, l.allTokens())

	l, err = New("#!/usr/bin/env gosh\nfoo", &Config{
		SkipShebang: true,
	})
	require.NoError(t, err)
	expected = []tokens.Token{
		{Offset: 20, Type: tokens.Identifier, Literal: `foo`},
		{Offset: 23, Type: tokens.EOF},
	}
	assert.Equal(t, expected, l.allTokens())

	l, err = New("#!/usr/bin/env gosh", &Config{
		SkipShebang: true,
	})
	require.NoError(t, err)
	expected = []tokens.Token{
		{Offset: 19, Type: tokens.EOF},
	}
	assert.Equal(t, expected, l.allTokens())
}

func TestZeroRune(t *testing.T) {
	for _, input := range []string{"12\x00", "1\x002"} {
		t.Run(input, func(t *testing.T) {
			gofuzz.AddDataToCorpus("scanner", []byte(input))

			l, err := New(input, nil)
			require.Error(t, err)
			require.Nil(t, l)
		})
	}
}
