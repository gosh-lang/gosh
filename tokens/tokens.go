// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tokens

import (
	"fmt"
)

// Token represents lexical token of the Gosh programming language.
type Token struct {
	Offset  int
	Type    Type
	Literal string
}

// String returns the string representation of the token.
func (tok Token) String() string {
	res := fmt.Sprintf("%d: %s", tok.Offset, tok.Type.String())
	if tok.Literal != "" {
		if tok.Type == Semicolon && tok.Literal == "\n" {
			res += " newline"
		} else {
			res += " " + tok.Literal
		}
	}
	return "[ " + res + " ]"
}

// check interfaces
var (
	_ fmt.Stringer = Token{}
)
