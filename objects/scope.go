// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package objects // import "gosh-lang.org/gosh/objects"

// A Scope maintains the set of named language entities declared in the scope
// and a link to the immediately surrounding (outer) scope.
type Scope struct {
	outer *Scope
	store map[string]Object
}

// NewScope creates a new scope nested in the outer scope.
func NewScope(outer *Scope) *Scope {
	return &Scope{
		outer: outer,
		store: make(map[string]Object),
	}
}

// Lookup return a named entity with this or outer scope (recursively).
func (e *Scope) Lookup(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Lookup(name)
	}
	return obj, ok
}

// Set adds or replaces a named entity in scope.
func (e *Scope) Set(name string, obj Object) {
	e.store[name] = obj
}
