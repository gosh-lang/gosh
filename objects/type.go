// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package objects // import "gosh-lang.org/gosh/objects"

//go:generate stringer -type Type

// Type is the set of object types of the Gosh programming language.
type Type int

// The list of object types.
const (
	IntegerType Type = iota
	BooleanType
	StringType
	FunctionType
	GoFunctionType
	ContinueType
)
