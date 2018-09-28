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
	l := len(input) // TODO len([]rune(input))
	s, err := New(input, &Config{
		SkipShebang: true,
	})
	if err != nil {
		return 0
	}

	t := s.allTokens()
	if len(t) == 0 {
		panic("should not return 0 tokens")
	}

	offset := -1
	for _, tok := range t {
		if offset >= tok.Offset {
			logTokens(t)
			panic(fmt.Sprintf("unexpected offset for token %s (previous offset: %d)", tok, offset))
		}
		offset = tok.Offset
	}

	last := t[len(t)-1]
	if last.Type == tokens.Illegal {
		if last.Offset == l {
			logTokens(t)
			panic(fmt.Sprintf("unexpected last illegal token offset: %d", last.Offset))
		}
		return 0
	}

	if last.Type != tokens.EOF {
		logTokens(t)
		panic(fmt.Sprintf("unexpected last token: %s", last))
	}
	if last.Offset != l {
		logTokens(t)
		panic(fmt.Sprintf("unexpected last token offset: %d (expected: %d)", last.Offset, l))
	}
	return 1
}
