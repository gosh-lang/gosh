// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build ignore

// run_fuzzer runs go-fuzz for the given package.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func run(args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "GO111MODULE=off")
	log.Print(strings.Join(cmd.Args, " "))
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetFlags(0)
	packages := []string{"scanner", "parser", "interpreter"}
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: go run misc/run_fuzzer.go [package]")
		fmt.Fprintln(flag.CommandLine.Output(), "       [package] is one of: "+strings.Join(packages, ", ")+".")
		flag.CommandLine.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	pack := flag.Arg(0)
	var found bool
	for _, p := range packages {
		if p == pack {
			found = true
			break
		}
	}
	if !found {
		flag.Usage()
		os.Exit(1)
	}

	importPath := "gosh-lang.org/gosh/" + pack
	workdir := filepath.Join("internal", "gofuzz", "workdir")
	file := filepath.Join(workdir, pack+"_fuzz.zip")

	if _, err := os.Stat(workdir); err != nil {
		log.Fatalf("%s.\nRun tests with `make` to create working directory and initial go-fuzz corpus", err)
	}

	run("go", "install", "-v", "-tags", "gofuzz", "./...")
	run("go-fuzz-build", "-o="+file, importPath)
	run("go-fuzz", "-bin="+file, "-workdir="+workdir)
}
