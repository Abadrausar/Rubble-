/*
Copyright 2014 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

// Rex Expression Commands.
package expr

import "dctech/rex"

// Setup adds the expression commands to the script.
// The expression commands are:
//	expr
func Setup(state *rex.State) (err error) {
	defer func() {
		if !state.NoRecover {
			if x := recover(); x != nil {
				if y, ok := x.(rex.ScriptError); ok {
					err = y
					return
				}
				panic(x)
			}
		}
	}()
	
	state.RegisterCommand("expr", Command_Expr)
	
	return nil
}

// Evaluates an expression.
// 	expr expression values...
// Supported operators (in precedence order):
//	( )
//	/ * %
//	+ -
//	== != > < >= <=
//	!
//	&& ||
// Values are specified by any non-operator character and are filled in by position.
// Returns the value of the evaluated expression.
func Command_Expr(script *rex.Script, params []*rex.Value) {
	if len(params) < 1 {
		rex.ErrorParamCount("expr", ">1")
	}

	lex := newLexer(params[0].String())
	params = params[1:]
	script.RetVal = expression(lex, &params)
}

// What follows is a simple expression parser with an implicit stack.

func expression(lex *lexer, vars *[]*rex.Value) *rex.Value {
	a := not(lex, vars)
	for lex.checkLook(tknAnd, tknOr) {
		lex.getToken(tknAnd, tknOr)
		tkn := lex.current
		b := not(lex, vars)
		switch tkn {
		case tknAnd:
			a = rex.NewValueBool(a.Bool() && b.Bool())
		case tknOr:
			a = rex.NewValueBool(a.Bool() || b.Bool())
		}
	}
	return a
}

func not(lex *lexer, vars *[]*rex.Value) *rex.Value {
	if lex.checkLook(tknNot) {
		lex.getToken(tknNot)
		return rex.NewValueBool(!relation(lex, vars).Bool())
	}
	return relation(lex, vars)
}

func relation(lex *lexer, vars *[]*rex.Value) *rex.Value {
	a := add(lex, vars)
	for lex.checkLook(tknEq, tknNotEq, tknGt, tknGtEq, tknLt, tknLtEq) {
		lex.getToken(tknEq, tknNotEq, tknGt, tknGtEq, tknLt, tknLtEq)
		tkn := lex.current
		b := add(lex, vars)
		switch tkn {
		case tknEq:
			if a.Type == rex.TypFloat {
				a = rex.NewValueBool(a.Float64() == b.Float64())
			} else if a.Type == rex.TypBool {
				a = rex.NewValueBool(a.Bool() == b.Bool())
			} else {
				a = rex.NewValueBool(a.Int64() == b.Int64())
			}
		case tknNotEq:
			if a.Type == rex.TypFloat {
				a = rex.NewValueBool(a.Float64() != b.Float64())
			} else if a.Type == rex.TypBool {
				a = rex.NewValueBool(a.Bool() != b.Bool())
			} else {
				a = rex.NewValueBool(a.Int64() != b.Int64())
			}
		case tknGt:
			if a.Type == rex.TypFloat {
				a = rex.NewValueBool(a.Float64() > b.Float64())
			} else {
				a = rex.NewValueBool(a.Int64() > b.Int64())
			}
		case tknGtEq:
			if a.Type == rex.TypFloat {
				a = rex.NewValueBool(a.Float64() >= b.Float64())
			} else {
				a = rex.NewValueBool(a.Int64() >= b.Int64())
			}
		case tknLt:
			if a.Type == rex.TypFloat {
				a = rex.NewValueBool(a.Float64() < b.Float64())
			} else {
				a = rex.NewValueBool(a.Int64() < b.Int64())
			}
		case tknLtEq:
			if a.Type == rex.TypFloat {
				a = rex.NewValueBool(a.Float64() <= b.Float64())
			} else {
				a = rex.NewValueBool(a.Int64() <= b.Int64())
			}
		}
	}
	return a
}

func add(lex *lexer, vars *[]*rex.Value) *rex.Value {
	a := mul(lex, vars)
	for lex.checkLook(tknAdd, tknSub) {
		lex.getToken(tknAdd, tknSub)
		tkn := lex.current
		b := mul(lex, vars)
		switch tkn {
		case tknAdd:
			if a.Type == rex.TypFloat {
				a = rex.NewValueFloat64(a.Float64() + b.Float64())
			} else {
				a = rex.NewValueInt64(a.Int64() + b.Int64())
			}
		case tknSub:
			if a.Type == rex.TypFloat {
				a = rex.NewValueFloat64(a.Float64() - b.Float64())
			} else {
				a = rex.NewValueInt64(a.Int64() - b.Int64())
			}
		}
	}
	return a
}

func mul(lex *lexer, vars *[]*rex.Value) *rex.Value {
	a := value(lex, vars)
	for lex.checkLook(tknMul, tknDiv, tknMod) {
		lex.getToken(tknMul, tknDiv, tknMod)
		tkn := lex.current
		b := value(lex, vars)
		switch tkn {
		case tknMul:
			if a.Type == rex.TypFloat {
				a = rex.NewValueFloat64(a.Float64() * b.Float64())
			} else {
				a = rex.NewValueInt64(a.Int64() * b.Int64())
			}
		case tknDiv:
			if a.Type == rex.TypFloat {
				a = rex.NewValueFloat64(a.Float64() / b.Float64())
			} else {
				a = rex.NewValueInt64(a.Int64() / b.Int64())
			}
		case tknMod:
			a = rex.NewValueInt64(a.Int64() % b.Int64())
		}
	}
	return a
}

func value(lex *lexer, vars *[]*rex.Value) *rex.Value {
	if lex.checkLook(tknOParen) {
		lex.getToken(tknOParen)
		a := expression(lex, vars)
		lex.getToken(tknCParen)
		return a
	}
	
	lex.getToken(tknVal)
	if len(*vars) <= 0 {
		rex.ErrorGeneralCmd("expr", "Insufficient values for expression.")
	}
	
	a := (*vars)[0]
	*vars = (*vars)[1:]
	return a
}
