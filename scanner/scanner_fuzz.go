// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build gofuzz

package scanner

import (
	"fmt"

	"gosh-lang.org/gosh/tokens"
)

func logTokens(tokens []tokens.Token) {
	for _, t := range tokens {
		fmt.Println(t)
	}
}

func Fuzz(data []byte) int {
	input := string(data)
	l := len([]rune(input))
	s, err := New(input, &Config{
		SkipShebang: true,
	})
	if err != nil {
		if err != errNulCharacter {
			panic(err)
		}
		return 0
	}

	// check that 0 tokens are never returned
	t := s.allTokens()
	if len(t) == 0 {
		panic("should not return 0 tokens")
	}

	// check that offsets are increasing
	offset := -1
	for _, tok := range t {
		if offset >= tok.Offset {
			logTokens(t)
			panic(fmt.Sprintf("unexpected offset for token %s (previous offset: %d)", tok, offset))
		}
		offset = tok.Offset
	}

	// check that last token is either Illegal in the middle of input, or EOF at the end
	last := t[len(t)-1]
	switch last.Type {
	case tokens.Illegal:
		if last.Offset == l {
			logTokens(t)
			panic(fmt.Sprintf("unexpected last illegal token offset: %d", last.Offset))
		}
		return 0

	case tokens.EOF:
		if last.Offset != l {
			logTokens(t)
			panic(fmt.Sprintf("unexpected last token offset: %d (expected: %d)", last.Offset, l))
		}
		return 1 // correct input

	default:
		logTokens(t)
		panic(fmt.Sprintf("unexpected last token: %s", last))
	}
}
