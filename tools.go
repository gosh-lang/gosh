// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build tools

package main

import (
	_ "github.com/dvyukov/go-fuzz/go-fuzz"
	_ "github.com/dvyukov/go-fuzz/go-fuzz-build"
	_ "golang.org/x/tools/cmd/stringer"
)
