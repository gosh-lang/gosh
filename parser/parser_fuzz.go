// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build gofuzz

package parser

import (
	"gosh-lang.org/gosh/scanner"
)

func Fuzz(data []byte) int {
	s, err := scanner.New(string(data), &scanner.Config{
		SkipShebang: true,
	})
	if err != nil {
		return 0
	}

	p := New(s, nil)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		return 0
	}
	program.String()
	return 1
}
