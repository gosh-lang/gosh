// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser // import "gosh-lang.org/gosh/parser"

import (
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gosh-lang.org/gosh/ast"
	"gosh-lang.org/gosh/internal/gofuzz"
	"gosh-lang.org/gosh/internal/golden"
	"gosh-lang.org/gosh/scanner"
	"gosh-lang.org/gosh/tokens"
)

var spewConfig = &spew.ConfigState{
	Indent:                  "  ",
	DisableMethods:          true,
	DisablePointerMethods:   true,
	DisablePointerAddresses: true,
	DisableCapacities:       true,
	ContinueOnMethod:        true,
}

func assertEqual(t *testing.T, expected, actual []ast.Statement) {
	t.Helper()

	if !assert.Equal(t, expected, actual) {
		e := spewConfig.Sdump(expected)
		a := spewConfig.Sdump(actual)
		diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        difflib.SplitLines(e),
			B:        difflib.SplitLines(a),
			FromFile: "Expected",
			FromDate: "",
			ToFile:   "Actual",
			ToDate:   "",
			Context:  1,
		})
		require.NoError(t, err)
		t.Logf("\n%s", diff)
	}
}

func formatErrors(errors []error) string {
	var res strings.Builder
	for _, err := range errors {
		res.WriteString(err.Error())
		res.WriteString("\n")
	}
	return res.String()
}

func TestGolden(t *testing.T) {
	for _, f := range golden.Data {
		t.Run(f.File, func(t *testing.T) {
			s, err := scanner.New(f.Source, &scanner.Config{
				SkipShebang: true,
			})
			require.NoError(t, err)
			p := New(s, &Config{
				crashOnError: true,
			})
			program := p.ParseProgram()
			require.Nil(t, p.Errors(), "%s", formatErrors(p.Errors()))

			actual := spewConfig.Sdump(program)
			assert.Equal(t, f.AST, strings.Split(actual, "\n"), "actual:\n%s", actual)

			actual = program.String()
			assert.Equal(t, f.Text, strings.Split(actual, "\n"), "actual:\n%s", actual)
		})
	}
}

func TestParser(t *testing.T) {
	for source, expected := range map[string]ast.Statement{
		"var answer = 42": &ast.VarStatement{
			Token: tokens.Token{Offset: 0, Type: tokens.Var, Literal: "var"},
			Name: &ast.Identifier{
				Token: tokens.Token{Offset: 4, Type: tokens.Identifier, Literal: "answer"},
				Value: "answer",
			},
			Value: &ast.IntegerLiteral{
				Token: tokens.Token{Offset: 13, Type: tokens.Integer, Literal: "42"},
				Value: 42,
			},
		},

		"answer = 42": &ast.AssignStatement{
			Token: tokens.Token{Offset: 7, Type: tokens.Assignment, Literal: "="},
			Name: &ast.Identifier{
				Token: tokens.Token{Offset: 0, Type: tokens.Identifier, Literal: "answer"},
				Value: "answer",
			},
			Value: &ast.IntegerLiteral{
				Token: tokens.Token{Offset: 9, Type: tokens.Integer, Literal: "42"},
				Value: 42,
			},
		},

		"answer == 42": &ast.ExpressionStatement{
			Token: tokens.Token{Offset: 0, Type: tokens.Identifier, Literal: "answer"},
			Expression: &ast.InfixExpression{
				Token: tokens.Token{Offset: 7, Type: tokens.Equal, Literal: "=="},
				Left: &ast.Identifier{
					Token: tokens.Token{Offset: 0, Type: tokens.Identifier, Literal: "answer"},
					Value: "answer",
				},
				Right: &ast.IntegerLiteral{
					Token: tokens.Token{Offset: 10, Type: tokens.Integer, Literal: "42"},
					Value: 42,
				},
			},
		},

		"answer += 42": &ast.AssignStatement{
			Token: tokens.Token{Offset: 7, Type: tokens.SumAssignment, Literal: "+="},
			Name: &ast.Identifier{
				Token: tokens.Token{Offset: 0, Type: tokens.Identifier, Literal: "answer"},
				Value: "answer",
			},
			Value: &ast.IntegerLiteral{
				Token: tokens.Token{Offset: 10, Type: tokens.Integer, Literal: "42"},
				Value: 42,
			},
		},

		"answer++": &ast.IncrementDecrementStatement{
			Token: tokens.Token{Offset: 6, Type: tokens.Increment, Literal: "++"},
			Name: &ast.Identifier{
				Token: tokens.Token{Offset: 0, Type: tokens.Identifier, Literal: "answer"},
				Value: "answer",
			},
		},

		"return 42": &ast.ReturnStatement{
			Token: tokens.Token{Offset: 0, Type: tokens.Return, Literal: "return"},
			Value: &ast.IntegerLiteral{
				Token: tokens.Token{Offset: 7, Type: tokens.Integer, Literal: "42"},
				Value: 42,
			},
		},

		"if (6 * 9 == 42) {\ntrue;\nfalse;\n}": &ast.IfStatement{
			Token: tokens.Token{Offset: 0, Type: tokens.If, Literal: "if"},
			Cond: &ast.InfixExpression{
				Token: tokens.Token{Offset: 10, Type: tokens.Equal, Literal: "=="},
				Left: &ast.InfixExpression{
					Token: tokens.Token{Offset: 6, Type: tokens.Product, Literal: "*"},
					Left: &ast.IntegerLiteral{
						Token: tokens.Token{Offset: 4, Type: tokens.Integer, Literal: "6"},
						Value: 6,
					},
					Right: &ast.IntegerLiteral{
						Token: tokens.Token{Offset: 8, Type: tokens.Integer, Literal: "9"},
						Value: 9,
					},
				},
				Right: &ast.IntegerLiteral{
					Token: tokens.Token{Offset: 13, Type: tokens.Integer, Literal: "42"},
					Value: 42,
				},
			},
			Body: &ast.BlockStatement{
				Token: tokens.Token{Offset: 17, Type: tokens.LBRACE, Literal: "{"},
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: tokens.Token{Offset: 19, Type: tokens.True, Literal: "true"},
						Expression: &ast.BooleanLiteral{
							Token: tokens.Token{Offset: 19, Type: tokens.True, Literal: "true"},
							Value: true,
						},
					},
					&ast.ExpressionStatement{
						Token: tokens.Token{Offset: 25, Type: tokens.False, Literal: "false"},
						Expression: &ast.BooleanLiteral{
							Token: tokens.Token{Offset: 25, Type: tokens.False, Literal: "false"},
							Value: false,
						},
					},
				},
			},
		},

		"for i = 1; i <= 100; i++ {\n}": &ast.ForStatement{
			Token: tokens.Token{Offset: 0, Type: tokens.For, Literal: "for"},
			Init: &ast.AssignStatement{
				Token: tokens.Token{Offset: 6, Type: tokens.Assignment, Literal: "="},
				Name: &ast.Identifier{
					Token: tokens.Token{Offset: 4, Type: tokens.Identifier, Literal: "i"},
					Value: "i",
				},
				Value: &ast.IntegerLiteral{
					Token: tokens.Token{Offset: 8, Type: tokens.Integer, Literal: "1"},
					Value: 1,
				},
			},
			Cond: &ast.InfixExpression{
				Token: tokens.Token{Offset: 13, Type: tokens.LessOrEqual, Literal: "<="},
				Left: &ast.Identifier{
					Token: tokens.Token{Offset: 11, Type: tokens.Identifier, Literal: "i"},
					Value: "i",
				},
				Right: &ast.IntegerLiteral{
					Token: tokens.Token{Offset: 16, Type: tokens.Integer, Literal: "100"},
					Value: 100,
				},
			},
			Post: &ast.IncrementDecrementStatement{
				Token: tokens.Token{Offset: 22, Type: tokens.Increment, Literal: "++"},
				Name: &ast.Identifier{
					Token: tokens.Token{Offset: 21, Type: tokens.Identifier, Literal: "i"},
					Value: "i",
				},
			},
			Body: &ast.BlockStatement{
				Token:      tokens.Token{Offset: 25, Type: tokens.LBRACE, Literal: "{"},
				Statements: []ast.Statement{},
			},
		},

		`println("answer")`: &ast.ExpressionStatement{
			Token: tokens.Token{Offset: 0, Type: tokens.Identifier, Literal: "println"},
			Expression: &ast.CallExpression{
				Token: tokens.Token{Offset: 7, Type: tokens.LPAREN, Literal: "("},
				Function: &ast.Identifier{
					Token: tokens.Token{Offset: 0, Type: tokens.Identifier, Literal: "println"},
					Value: "println",
				},
				Arguments: []ast.Expression{
					&ast.StringLiteral{
						Token: tokens.Token{Offset: 8, Type: tokens.String, Literal: `"answer"`},
						Value: "answer",
					},
				},
			},
		},
	} {
		t.Run(source, func(t *testing.T) {
			formal := source + ";\n"
			for _, input := range []string{source, source + ";", source + "\n", formal, source + ";\n\n;;"} {
				gofuzz.AddDataToCorpus("parser", []byte(input))

				s, err := scanner.New(input, &scanner.Config{
					SkipShebang: true,
				})
				require.NoError(t, err)
				p := New(s, &Config{
				// crashOnError: true,
				})
				program := p.ParseProgram()
				require.Nil(t, p.Errors(), "%s", formatErrors(p.Errors()))
				require.NotNil(t, program)
				assertEqual(t, []ast.Statement{expected}, program.Statements)
				assert.Equal(t, formal, program.String())
				assert.Equal(t, tokens.Token{Offset: len(input), Type: tokens.EOF}, p.curToken)
			}
		})
	}
}

func TestErrors(t *testing.T) {
	for input, errors := range map[string][]error{
		`(`: {
			&Error{Err: "no prefix parse function for EOF found (token [ 1: EOF ])"},
			&Error{Err: "expected next token to be RPAREN, got [ 2: EOF ] instead"},
		},
	} {
		t.Run(input, func(t *testing.T) {
			gofuzz.AddDataToCorpus("parser", []byte(input))

			s, err := scanner.New(input, &scanner.Config{
				SkipShebang: true,
			})
			require.NoError(t, err)
			p := New(s, nil)
			program := p.ParseProgram()
			assert.Nil(t, program)
			assert.Equal(t, errors, p.Errors())
		})
	}
}
