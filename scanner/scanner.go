// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package scanner

import (
	"fmt"

	"gosh-lang.org/gosh/tokens"
)

// Scanner extracts tokens from Gosh source code.
type Scanner struct {
	config *Config
	input  []rune

	rPos            int  // current rune position in input
	r               rune // current rune; the same as input[rPos]
	insertSemicolon bool // return next \n as semicolon
}

// Config configures scanner.
type Config struct {
	SkipShebang bool // if true, scanner will skip the first line of input if it starts with #!

	crashOnError        bool // crash scanner on any illegal token, for testing only
	dontInsertSemicolon bool // disable automatic semicolon insertion, for testing only
}

var keywords = map[string]tokens.Type{
	"break":       tokens.Break,
	"case":        tokens.Case,
	"chan":        tokens.Chan,
	"const":       tokens.Const,
	"continue":    tokens.Continue,
	"default":     tokens.Default,
	"defer":       tokens.Defer,
	"else":        tokens.Else,
	"fallthrough": tokens.Fallthrough,
	"for":         tokens.For,
	"func":        tokens.Func,
	"go":          tokens.Go,
	"goto":        tokens.Goto,
	"if":          tokens.If,
	"import":      tokens.Import,
	"interface":   tokens.Interface,
	"map":         tokens.Map,
	"package":     tokens.Package,
	"range":       tokens.Range,
	"return":      tokens.Return,
	"select":      tokens.Select,
	"struct":      tokens.Struct,
	"switch":      tokens.Switch,
	"var":         tokens.Var,

	// TODO remove - those are not keywords
	"true":  tokens.True,
	"false": tokens.False,
}

var errNulCharacter = fmt.Errorf("input contains NUL character (U+0000)")

// New creates new scanner for the given Gosh source code.
func New(input string, config *Config) (*Scanner, error) {
	if config == nil {
		config = new(Config)
	}

	runes := []rune(input)
	for _, r := range runes {
		if r == 0 {
			return nil, errNulCharacter
		}
	}

	l := &Scanner{
		config: config,
		input:  runes,
		rPos:   -1,
	}
	l.readRune()
	return l, nil
}

func isLetter(r rune) bool {
	switch {
	case 'a' <= r && r <= 'z':
		return true
	case 'A' <= r && r <= 'Z':
		return true
	case r == '_':
		return true
	default:
		return false
	}
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func (s *Scanner) crash(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	panic(fmt.Errorf("%s\nrPos: %d\nr: %q", msg, s.rPos, s.r))
}

func (s *Scanner) peekRune() rune {
	peekPos := s.rPos + 1
	if peekPos >= len(s.input) {
		return 0
	}
	return s.input[peekPos]
}

func (s *Scanner) readRune() {
	s.r = s.peekRune()
	s.rPos++
}

func (s *Scanner) skipWhitespace() {
	for {
		switch s.r {
		case ' ', '\t', '\r':
			s.readRune()
		case '\n':
			if s.insertSemicolon {
				return
			}
			s.readRune()
		default:
			return
		}
	}
}

// readLine reads and returns line up to '\n' or EOF.
func (s *Scanner) readLine() string {
	pos := s.rPos
	for {
		s.readRune()
		if s.r == '\n' || s.r == 0 {
			break
		}
	}
	return string(s.input[pos:s.rPos])
}

func (s *Scanner) readInt() (string, bool) {
	pos := s.rPos
	for isDigit(s.r) {
		s.readRune()
	}
	ok := true
	if isLetter(s.r) {
		ok = false
		for isLetter(s.r) || isDigit(s.r) {
			s.readRune()
		}
	}
	return string(s.input[pos:s.rPos]), ok
}

func (s *Scanner) readString() (string, bool) {
	pos := s.rPos
	for {
		s.readRune()
		if s.r == '"' || s.r == 0 {
			break
		}
	}
	ok := s.r == '"'
	if ok {
		s.readRune()
	}
	return string(s.input[pos:s.rPos]), ok
}

func (s *Scanner) readIdentifier() string {
	pos := s.rPos
	for isLetter(s.r) || isDigit(s.r) {
		s.readRune()
	}
	return string(s.input[pos:s.rPos])
}

func (s *Scanner) lookupIdentifier(ident string) tokens.Type {
	if t, ok := keywords[ident]; ok {
		return t
	}
	return tokens.Identifier
}

// NextToken returns next scanned tokens.
// Once it returns tokens.EOF, it will continue to do so.
//nolint:gocyclo
func (s *Scanner) NextToken() tokens.Token {
	s.skipWhitespace()
	tok := tokens.Token{Offset: s.rPos, Type: tokens.Illegal}

	if s.config.crashOnError {
		defer func() {
			if tok.Type == tokens.Illegal {
				s.crash("illegal token: %s", tok)
			}
		}()
	}

	var insertSemicolon bool
	defer func() {
		if !s.config.dontInsertSemicolon {
			s.insertSemicolon = insertSemicolon
		}
	}()

	switch s.r {
	case 0:
		tok.Type = tokens.EOF
	case '\n':
		// s.skipWhitespace() exited on \n, insert semicolon
		tok.Type = tokens.Semicolon
		tok.Literal = "\n"
	case '#':
		if s.rPos == 0 && s.peekRune() == '!' && s.config.SkipShebang {
			s.readLine()
			return s.NextToken()
		}
		tok.Literal = string(s.r)

	case '=':
		switch s.peekRune() {
		case '=':
			s.readRune()
			tok.Type = tokens.Equal
			tok.Literal = "=="
		default:
			tok.Type = tokens.Assignment
			tok.Literal = "="
		}
	case ':':
		switch s.peekRune() {
		case '=':
			s.readRune()
			tok.Type = tokens.Define
			tok.Literal = ":="
		default:
			tok.Type = tokens.Colon
			tok.Literal = ":"
		}

	case '+':
		switch s.peekRune() {
		case '+':
			s.readRune()
			tok.Type = tokens.Increment
			tok.Literal = "++"
			insertSemicolon = true
		case '=':
			s.readRune()
			tok.Type = tokens.SumAssignment
			tok.Literal = "+="
		default:
			tok.Type = tokens.Sum
			tok.Literal = "+"
		}
	case '-':
		switch s.peekRune() {
		case '-':
			s.readRune()
			tok.Type = tokens.Decrement
			tok.Literal = "--"
			insertSemicolon = true
		case '=':
			s.readRune()
			tok.Type = tokens.DifferenceAssignment
			tok.Literal = "-="
		default:
			tok.Type = tokens.Difference
			tok.Literal = "-"
		}
	case '*':
		switch s.peekRune() {
		case '=':
			s.readRune()
			tok.Type = tokens.ProductAssignment
			tok.Literal = "*="
		default:
			tok.Type = tokens.Product
			tok.Literal = "*"
		}
	case '/':
		switch s.peekRune() {
		case '/':
			tok.Type = tokens.Comment
			tok.Literal = s.readLine()
			return tok // l.readRune() already called by l.readLine(), so exit early
		case '=':
			s.readRune()
			tok.Type = tokens.QuotientAssignment
			tok.Literal = "/="
		default:
			tok.Type = tokens.Quotient
			tok.Literal = "/"
		}
	case '%':
		switch s.peekRune() {
		case '=':
			s.readRune()
			tok.Type = tokens.RemainderAssignment
			tok.Literal = "%="
		default:
			tok.Type = tokens.Remainder
			tok.Literal = "%"
		}

	case '&':
		switch s.peekRune() {
		case '&':
			s.readRune()
			tok.Type = tokens.LogicalAnd
			tok.Literal = "&&"
		default:
			tok.Type = tokens.BitwiseAnd
			tok.Literal = "&"
		}
	case '|':
		switch s.peekRune() {
		case '|':
			s.readRune()
			tok.Type = tokens.LogicalOr
			tok.Literal = "||"
		default:
			tok.Type = tokens.BitwiseOr
			tok.Literal = "|"
		}
	case '^':
		tok.Type = tokens.BitwiseXor
		tok.Literal = "^"

	case '!':
		switch s.peekRune() {
		case '=':
			s.readRune()
			tok.Type = tokens.NotEqual
			tok.Literal = "!="
		default:
			tok.Type = tokens.Not
			tok.Literal = "!"
		}

	case '<':
		switch s.peekRune() {
		case '=':
			s.readRune()
			tok.Type = tokens.LessOrEqual
			tok.Literal = "<="
		default:
			tok.Type = tokens.Less
			tok.Literal = "<"
		}
	case '>':
		switch s.peekRune() {
		case '=':
			s.readRune()
			tok.Type = tokens.GreaterOrEqual
			tok.Literal = ">="
		default:
			tok.Type = tokens.Greater
			tok.Literal = ">"
		}

	case ';':
		tok.Type = tokens.Semicolon
		tok.Literal = ";"
	case ',':
		tok.Type = tokens.Comma
		tok.Literal = ","
	case '.':
		tok.Type = tokens.Period
		tok.Literal = "."

	case '(':
		tok.Type = tokens.LPAREN
		tok.Literal = "("
	case ')':
		tok.Type = tokens.RPAREN
		tok.Literal = ")"
		insertSemicolon = true
	case '{':
		tok.Type = tokens.LBRACE
		tok.Literal = "{"
	case '}':
		tok.Type = tokens.RBRACE
		tok.Literal = "}"
		insertSemicolon = true

	case '"':
		lit, ok := s.readString()
		tok.Literal = lit
		if ok {
			tok.Type = tokens.String
		}
		insertSemicolon = true
		return tok // l.readRune() already called by l.readString(), so exit early

	default:
		switch {
		case isLetter(s.r):
			tok.Literal = s.readIdentifier()
			tok.Type = s.lookupIdentifier(tok.Literal)
			switch tok.Type {
			case tokens.Identifier, tokens.Break, tokens.Continue, tokens.Fallthrough, tokens.Return:
				fallthrough
			case tokens.True, tokens.False: // TODO remove
				insertSemicolon = true
			}
			return tok // l.readRune() already called by l.readIdentifier(), so exit early

		case isDigit(s.r):
			lit, ok := s.readInt()
			tok.Literal = lit
			if ok {
				tok.Type = tokens.Integer
			}

			if s.r == '.' {
				// Consume '.' and see if we can parse a float
				s.readRune()
				tok.Type = tokens.Float

				s.peekRune()

				lit, ok = s.readInt()
				switch {
				case ok:
					tok.Literal = fmt.Sprintf("%v.%v", tok.Literal, lit)
				case len(lit) > 0:
					tok.Literal = fmt.Sprintf("%v.%v", tok.Literal, lit)
					tok.Type = tokens.Illegal
				}
			}

			insertSemicolon = true
			return tok // l.readRune() already called by l.readInt(), so exit early

		default:
			// TODO insertSemicolon?
			tok.Literal = string(s.r)
		}
	}

	s.readRune()
	return tok
}

// allTokens returns all tokens until tokens.EOF or tokens.ILLEGAL.
func (s *Scanner) allTokens() []tokens.Token {
	var res []tokens.Token
	for {
		tok := s.NextToken()
		res = append(res, tok)
		switch tok.Type {
		case tokens.EOF, tokens.Illegal:
			return res
		}
	}
}
