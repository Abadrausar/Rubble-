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

// Rex Integer Commands.
package integer

import "dctech/rex"

// Adds the integer commands to the state.
// The integer commands are:
//	int:add
//	int:sub
//	int:div
//	int:mod
//	int:mul
//	int:gt
//	int:lt
//	int:eq
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
	
	mod := state.RegisterModule("int")
	mod.RegisterCommand("add", Command_Add)
	mod.RegisterCommand("sub", Command_Sub)
	mod.RegisterCommand("div", Command_Div)
	mod.RegisterCommand("mod", Command_Mod)
	mod.RegisterCommand("mul", Command_Mul)
	mod.RegisterCommand("gt", Command_Gt)
	mod.RegisterCommand("lt", Command_Lt)
	mod.RegisterCommand("eq", Command_Eq)
	
	return nil
}

// Adds two values.
// 	int:add a b
// Returns a + b
func Command_Add(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("int:add", "2")
	}

	script.RetVal = rex.NewValueInt64(params[0].Int64() + params[1].Int64())
	return
}

// Subtracts two values.
// 	int:sub a b
// Returns a - b
func Command_Sub(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("int:sub", "2")
	}

	script.RetVal = rex.NewValueInt64(params[0].Int64() - params[1].Int64())
	return
}

// Divides two values.
// 	int:div a b
// Returns a / b
func Command_Div(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("int:div", "2")
	}

	script.RetVal = rex.NewValueInt64(params[0].Int64() / params[1].Int64())
	return
}

// Gives the remainder of dividing two values.
// 	int:mod a b
// Returns a % b
func Command_Mod(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("int:mod", "2")
	}

	script.RetVal = rex.NewValueInt64(params[0].Int64() % params[1].Int64())
	return
}

// Multiplies two values.
// 	int:mul a b
// Returns a * b
func Command_Mul(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("int:mul", "2")
	}

	script.RetVal = rex.NewValueInt64(params[0].Int64() * params[1].Int64())
	return
}

// Integer greater than.
// 	int:gt a b
// Returns true or false.
func Command_Gt(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("int:gt", "2")
	}

	result := params[0].Int64() > params[1].Int64()
	if result {
		script.RetVal = rex.NewValueBool(true)
		return
	}
	script.RetVal = rex.NewValueBool(false)
}

// Integer less than.
// 	int:lt a b
// Returns true or false.
func Command_Lt(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("int:lt", "2")
	}

	result := params[0].Int64() < params[1].Int64()
	if result {
		script.RetVal = rex.NewValueBool(true)
		return
	}
	script.RetVal = rex.NewValueBool(false)
}

// Integer equal.
// 	int:eq a b
// Returns true or false.
func Command_Eq(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("int:eq", "2")
	}

	result := params[0].Int64() == params[1].Int64()
	if result {
		script.RetVal = rex.NewValueBool(true)
		return
	}
	script.RetVal = rex.NewValueBool(false)
}
