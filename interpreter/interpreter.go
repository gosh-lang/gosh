// Gosh programming language.
// Copyright (c) 2018 Alexey Palazhchenko and contributors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package interpreter // import "gosh-lang.org/gosh/interpreter"

import (
	"context"
	"fmt"

	"gosh-lang.org/gosh/ast"
	"gosh-lang.org/gosh/objects"
	"gosh-lang.org/gosh/tokens"
)

// Interpreter evaluates Gosh AST nodes.
type Interpreter struct {
	config *Config
}

// Config configures interpreter.
type Config struct {
}

// New creates a new interpreter.
func New(config *Config) *Interpreter {
	if config == nil {
		config = new(Config)
	}

	return &Interpreter{
		config: config,
	}
}

func (i *Interpreter) crash(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	panic(msg)
}

// Eval evaluates given node in the given scope.
func (i *Interpreter) Eval(ctx context.Context, node ast.Node, scope *objects.Scope) objects.Object {
	if ctx.Err() != nil {
		// FIXME return error
		return nil
	}

	switch node := node.(type) {
	case *ast.Program:
		var res objects.Object
		for _, s := range node.Statements {
			res = i.Eval(ctx, s, scope)
		}
		return res

	case *ast.BlockStatement:
		var res objects.Object
		for _, s := range node.Statements {
			res = i.Eval(ctx, s, scope)
			if res != nil {
				switch res.Type() {
				case objects.ContinueType:
					return res
				}
			}
		}
		return res

	case *ast.ExpressionStatement:
		return i.Eval(ctx, node.Expression, scope)

	case *ast.ReturnStatement:
		return i.Eval(ctx, node.Value, scope)

	case *ast.VarStatement:
		val := i.Eval(ctx, node.Value, scope)
		scope.Set(node.Name.Value, val)
		return nil

	case *ast.AssignStatement:
		return i.evalAssignStatement(ctx, node, scope)

	case *ast.ForStatement:
		return i.evalForStatement(ctx, node, scope)

	case *ast.IfStatement:
		return i.evalIfStatement(ctx, node, scope)

	case *ast.IncrementDecrementStatement:
		return i.evalIncrementDecrementStatement(node, scope)

	case *ast.ContinueStatement:
		return &objects.Continue{}

	case *ast.Identifier:
		val, ok := scope.Lookup(node.Value)
		if !ok {
			i.crash("identifier not found: %s", node.Value)
		}
		return val

	case *ast.PrefixExpression:
		right := i.Eval(ctx, node.Right, scope)
		return i.evalPrefixExpression(node.Token.Literal, right)

	case *ast.InfixExpression:
		left := i.Eval(ctx, node.Left, scope)
		right := i.Eval(ctx, node.Right, scope)
		return i.evalInfixExpression(node.Token.Literal, left, right)

	case *ast.IntegerLiteral:
		return &objects.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &objects.Float{Value: node.Value}

	case *ast.BooleanLiteral:
		return &objects.Boolean{Value: node.Value}

	case *ast.StringLiteral:
		return &objects.String{Value: node.Value}

	case *ast.FunctionLiteral:
		return &objects.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Scope:      scope,
		}

	case *ast.CallExpression:
		return i.evalCallExpression(ctx, node, scope)

	default:
		i.crash("unexpected node %T:\n%#v", node, node)
		panic("not reached")
	}
}

func (i *Interpreter) evalPrefixExpression(operator string, right objects.Object) objects.Object {
	switch operator {
	case "!":
		if b, ok := right.(*objects.Boolean); ok {
			return &objects.Boolean{Value: !b.Value}
		}
		i.crash("prefix expression operator ! on %T:\n%#v", right, right)

	case "-":
		if i, ok := right.(*objects.Integer); ok {
			return &objects.Integer{Value: -i.Value}
		}
		i.crash("prefix expression operator - on %T:\n%#v", right, right)

	default:
		i.crash("unhandled prefix expression operator %s", operator)
	}
	panic("not reached")
}

func (i *Interpreter) evalInfixIntegerExpression(operator string, left, right int) objects.Object {
	switch operator {
	case "+":
		return &objects.Integer{Value: left + right}
	case "-":
		return &objects.Integer{Value: left - right}
	case "*":
		return &objects.Integer{Value: left * right}
	case "/":
		return &objects.Integer{Value: left / right}
	case "%":
		return &objects.Integer{Value: left % right}

	case "<":
		return &objects.Boolean{Value: left < right}
	case "<=":
		return &objects.Boolean{Value: left <= right}
	case ">":
		return &objects.Boolean{Value: left > right}
	case ">=":
		return &objects.Boolean{Value: left >= right}
	case "==":
		return &objects.Boolean{Value: left == right}
	case "!=":
		return &objects.Boolean{Value: left != right}

	default:
		i.crash("unhandled infix expression operator %s for two Integers", operator)
		panic("not reached")
	}
}

func (i *Interpreter) evalInfixFloatExpression(operator string, left, right float64) objects.Object {
	switch operator {
	case "+":
		return &objects.Float{Value: left + right}
	case "-":
		return &objects.Float{Value: left - right}
	case "*":
		return &objects.Float{Value: left * right}
	case "/":
		return &objects.Float{Value: left / right}

	case "<":
		return &objects.Boolean{Value: left < right}
	case "<=":
		return &objects.Boolean{Value: left <= right}
	case ">":
		return &objects.Boolean{Value: left > right}
	case ">=":
		return &objects.Boolean{Value: left >= right}
	case "==":
		return &objects.Boolean{Value: left == right}
	case "!=":
		return &objects.Boolean{Value: left != right}

	default:
		i.crash("unhandled infix expression operator %s for two Floats", operator)
		panic("not reached")
	}
}


func (i *Interpreter) evalInfixBooleanExpression(operator string, left, right bool) objects.Object {
	switch operator {
	case "==":
		return &objects.Boolean{Value: left == right}
	case "!=":
		return &objects.Boolean{Value: left != right}
	case "&&":
		return &objects.Boolean{Value: left && right}
	case "||":
		return &objects.Boolean{Value: left || right}
	default:
		i.crash("unhandled infix expression operator %s for two Booleans", operator)
		panic("not reached")
	}
}

func (i *Interpreter) evalInfixExpression(operator string, left, right objects.Object) objects.Object {
	switch left.Type() {
	case objects.IntegerType:
		switch right.Type() {
		case objects.IntegerType:
			l := left.(*objects.Integer).Value
			r := right.(*objects.Integer).Value
			return i.evalInfixIntegerExpression(operator, l, r)
		}

	case objects.FloatType:
		switch right.Type() {
		case objects.FloatType:
			l := left.(*objects.Float).Value
			r := right.(*objects.Float).Value
			return i.evalInfixFloatExpression(operator, l, r)
		}

	case objects.BooleanType:
		switch right.Type() {
		case objects.BooleanType:
			l := left.(*objects.Boolean).Value
			r := right.(*objects.Boolean).Value
			return i.evalInfixBooleanExpression(operator, l, r)
		}
	}

	i.crash("unhandled combination: %T %s %T", left, operator, right)
	panic("not reached")
}

func (i *Interpreter) evalExpressions(ctx context.Context, exps []ast.Expression, scope *objects.Scope) []objects.Object {
	res := make([]objects.Object, len(exps))
	for n, e := range exps {
		res[n] = i.Eval(ctx, e, scope)
	}
	return res
}

func (i *Interpreter) evalAssignStatement(ctx context.Context, node *ast.AssignStatement, scope *objects.Scope) objects.Object {
	val := i.Eval(ctx, node.Value, scope)
	switch node.Token.Type {
	case tokens.Assignment:
		// nothing
	default:
		i.crash("unhandled token %s", node.Token)
	}
	scope.Set(node.Name.Value, val)
	return nil
}

func (i *Interpreter) evalForStatement(ctx context.Context, node *ast.ForStatement, scope *objects.Scope) objects.Object {
	i.Eval(ctx, node.Init, scope)
	for {
		cond := i.Eval(ctx, node.Cond, scope)
		var b *objects.Boolean
		var ok bool
		if b, ok = cond.(*objects.Boolean); !ok {
			i.crash("expected boolean, got %T %s", cond, cond)
		}
		if !b.Value {
			return nil
		}

		body := i.Eval(ctx, node.Body, scope)

		i.Eval(ctx, node.Post, scope)

		if body != nil {
			switch body.Type() {
			case objects.ContinueType:
				continue
			}
		}
	}
}

func (i *Interpreter) evalIfStatement(ctx context.Context, node *ast.IfStatement, scope *objects.Scope) objects.Object {
	cond := i.Eval(ctx, node.Cond, scope)
	var b *objects.Boolean
	var ok bool
	if b, ok = cond.(*objects.Boolean); !ok {
		i.crash("expected boolean, got %T %s", cond, cond)
	}
	if !b.Value {
		return nil
	}

	body := i.Eval(ctx, node.Body, scope)
	if body != nil {
		switch body.Type() {
		case objects.ContinueType:
			return body
		}
	}
	return nil
}

func (i *Interpreter) evalIncrementDecrementStatement(node *ast.IncrementDecrementStatement, scope *objects.Scope) objects.Object {
	name := node.Name.Value
	val, ok := scope.Lookup(name)
	if !ok {
		i.crash("failed to lookup %s", name)
	}

	v := val.(*objects.Integer).Value

	switch node.Token.Type {
	case tokens.Increment:
		v++
	case tokens.Decrement:
		v--
	default:
		i.crash("unexpected token")
	}

	scope.Set(name, &objects.Integer{Value: v})
	return nil
}

func (i *Interpreter) evalCallExpression(ctx context.Context, node *ast.CallExpression, scope *objects.Scope) objects.Object {
	f := i.Eval(ctx, node.Function, scope)
	args := i.evalExpressions(ctx, node.Arguments, scope)
	switch f := f.(type) {
	case *objects.Function:
		newScope := objects.NewScope(scope)
		for i, name := range f.Parameters {
			newScope.Set(name.Value, args[i])
		}
		return i.Eval(ctx, f.Body, newScope)
	case *objects.GoFunction:
		return f.Func(args...)
	default:
		i.crash("unexpected node %T:\n%#v", node, node)
		panic("not reached")
	}
}
