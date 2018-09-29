// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/davecgh/go-spew/spew"
	"github.com/peterh/liner"
	"gopkg.in/alecthomas/kingpin.v2"

	"gosh-lang.org/gosh/interpreter"
	"gosh-lang.org/gosh/objects"
	"gosh-lang.org/gosh/parser"
	"gosh-lang.org/gosh/scanner"
	"gosh-lang.org/gosh/tokens"
)

var (
	// Version of the interpreter.
	Version = "0.0.1-dev"
)

// flags
var (
	DebugScannerF *bool
	DebugASTF     *bool
	DebugParserF  *bool
)

var versionRE = regexp.MustCompile(`go(\S+)`)

func extractGoVersion(s string) string {
	res := versionRE.FindStringSubmatch(s)
	if len(res) == 2 {
		return res[1]
	}
	return ""
}

func goVersion() string {
	cmd := exec.Command("go", "version") //nolint:gas
	b, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return string(b)
}

func eval(line string, scope *objects.Scope) {
	s, err := scanner.New(line, &scanner.Config{
		SkipShebang: true,
	})
	if err != nil {
		log.Printf("Scanner error: %s.", err)
		return
	}
	if *DebugScannerF {
		log.Print("Tokens:")
		for {
			t := s.NextToken()
			log.Print(t)
			switch t.Type {
			case tokens.EOF, tokens.Illegal:
				return
			}
		}
	}

	p := parser.New(s, nil)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		log.Print("Parser errors:\n")
		for _, e := range p.Errors() {
			log.Printf("\t%s", e)
		}
		return
	}
	if *DebugASTF {
		cfg := &spew.ConfigState{
			Indent:                  "  ",
			DisableMethods:          true,
			DisablePointerMethods:   true,
			DisablePointerAddresses: true,
			DisableCapacities:       true,
			ContinueOnMethod:        true,
		}
		b := cfg.Sdump(program)
		log.Printf("AST:\n%s", b)
		return
	}
	if *DebugParserF {
		log.Printf("Parsed program:\n%s", program.String())
		return
	}

	i := interpreter.New(nil)
	res := i.Eval(context.TODO(), program, scope)
	if res != nil {
		fmt.Println(res.String())
	}
}

func evalFile(filename string) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	scope := objects.NewScope(objects.Builtin(os.Stdout))
	eval(string(b), scope)
}

// readREPLHistory reads REPL history from from file and returns file name where it should be wrote at exit.
// It returns empty string if history can't be wrote at exit.
func readREPLHistory(liner *liner.State) string {
	var historyFilename string
	u, err := user.Current()
	if err == nil && u.HomeDir != "" {
		historyFilename = filepath.Join(u.HomeDir, ".gosh_history")
	}
	if historyFilename == "" {
		// we failed to detect user's home directory, so we can't read and write history file
		return ""
	}

	f, err := os.Open(historyFilename)
	if err != nil && os.IsNotExist(err) {
		// history file does not exist, but we can create it at exit
		return historyFilename
	}
	if err != nil {
		log.Printf("Warning: failed to open history file %s: %s.", historyFilename, err)
		return ""
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Printf("Warning: failed to close history file %s: %s.", historyFilename, err)
		}
	}()

	if _, err = liner.ReadHistory(f); err != nil {
		log.Printf("Warning: failed to read history file %s: %s.", historyFilename, err)
	}
	return historyFilename
}

func writeREPLHistory(liner *liner.State, historyFilename string) {
	f, err := os.Create(historyFilename)
	if err != nil {
		log.Printf("Warning: failed to create history file %s: %s.", historyFilename, err)
		return
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Printf("Warning: failed to close history file %s: %s.", historyFilename, err)
		}
	}()

	if _, err = liner.WriteHistory(f); err != nil {
		log.Printf("Warning: failed to write history file %s: %s.", historyFilename, err)
	}
}

func runREPL() {
	fmt.Printf("Gosh v%s. https://gosh-lang.org/\n", Version)
	fmt.Printf("Built with Go v%s.\n", extractGoVersion(runtime.Version()))
	fmt.Printf("Runtime Go v%s.\n", extractGoVersion(goVersion()))

	liner := liner.NewLiner()
	historyFilename := readREPLHistory(liner)
	defer func() {
		fmt.Println()
		writeREPLHistory(liner, historyFilename)
		if err := liner.Close(); err != nil {
			log.Printf("Failed to close liner: %s.", err)
		}
	}()

	scope := objects.NewScope(objects.Builtin(os.Stdout))
	for {
		line, err := liner.Prompt(`\ʕ•ϖ•ʔ/ >> `)
		switch err {
		case nil:
			liner.AppendHistory(line)
			eval(line, scope)
		case io.EOF:
			return
		default:
			log.Fatal(err)
		}
	}
}

func main() {
	log.SetFlags(0)

	DebugScannerF = kingpin.Flag("debug-scanner", "Print tokens and exit.").Bool()
	DebugASTF = kingpin.Flag("debug-ast", "Print AST and exit.").Bool()
	DebugParserF = kingpin.Flag("debug-parser", "Print parsed program and exit.").Bool()
	fileArg := kingpin.Arg("file", "Gosh program file.").String()
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	if *fileArg == "" {
		runREPL()
	} else {
		evalFile(*fileArg)
	}
}
