// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package golden // import "gosh-lang.org/gosh/internal/golden"

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	"gosh-lang.org/gosh/internal/gofuzz"
)

// File describes a single golden file.
type File struct {
	File   string
	Source string
	Tokens string
	AST    []string
	Text   []string
	Output []string
}

// Data is a golden data shared by all golden tests.
var Data []File

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller(0) failed")
	}
	files, err := filepath.Glob(filepath.Join(filepath.Dir(file), "*.gosh"))
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		source, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}

		gofuzz.AddFileToCorpus(filepath.Base(f), source)

		tokens, err := ioutil.ReadFile(f + ".tokens")
		if err != nil {
			panic(err)
		}

		ast, err := ioutil.ReadFile(f + ".ast")
		if err != nil {
			panic(err)
		}

		text, err := ioutil.ReadFile(f + ".text")
		if err != nil {
			panic(err)
		}

		output, err := ioutil.ReadFile(f + ".output")
		if err != nil {
			panic(err)
		}

		Data = append(Data, File{
			File:   filepath.Base(f),
			Source: string(source),
			Tokens: strings.TrimSpace(string(tokens)),
			AST:    strings.Split(string(ast), "\n"),
			Text:   strings.Split(string(text), "\n"),
			Output: strings.Split(string(output), "\n"),
		})
	}

	// double check
	const expected = 1
	if len(Data) != expected {
		panic(fmt.Sprintf("expected %d files, read %d", expected, len(Data)))
	}
}
