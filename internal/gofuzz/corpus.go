// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gofuzz

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

var prefix string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller(0) failed")
	}
	prefix = filepath.Join(filepath.Dir(file), "workdir", "corpus")
	if err := os.MkdirAll(prefix, 0750); err != nil {
		panic(err)
	}
}

// AddFileToCorpus adds named Gosh source code fragment to the go-fuzz corpus.
func AddFileToCorpus(name string, data []byte) {
	path := filepath.Join(prefix, name)
	if err := ioutil.WriteFile(path, data, 0640); err != nil {
		panic(err)
	}
}

// AddDataToCorpus adds unnamed Gosh source code fragment to the go-fuzz corpus.
func AddDataToCorpus(prefix string, data []byte) {
	name := fmt.Sprintf("%s-%040x.gosh", prefix, sha1.Sum(data))
	AddFileToCorpus(name, data)
}
