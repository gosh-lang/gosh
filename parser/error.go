// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

// Error is a parser error.
type Error struct {
	Err string
}

func (e *Error) Error() string {
	return e.Err
}

// check interfaces
var (
	_ error = (*Error)(nil)
)
