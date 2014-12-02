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

// Rex Floating Point Commands.
package float

import "dctech/rex"

// Adds the floating point commands to the state.
// The floating point commands are:
//	float:add
//	float:sub
//	float:div
//	float:mul
//	float:gt
//	float:lt
//	float:eq
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
	
	mod := state.RegisterModule("float")
	mod.RegisterCommand("add", Command_Add)
	mod.RegisterCommand("sub", Command_Sub)
	mod.RegisterCommand("div", Command_Div)
	mod.RegisterCommand("mul", Command_Mul)
	mod.RegisterCommand("gt", Command_Gt)
	mod.RegisterCommand("lt", Command_Lt)
	mod.RegisterCommand("eq", Command_Eq)
	
	return nil
}

// Adds two values.
// 	float:add a b
// Returns a + b
func Command_Add(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("float:add", "2")
	}

	script.RetVal = rex.NewValueFloat64(params[0].Float64() + params[1].Float64())
	return
}

// Subtracts two values.
// 	float:sub a b
// Returns a - b
func Command_Sub(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("float:sub", "2")
	}

	script.RetVal = rex.NewValueFloat64(params[0].Float64() - params[1].Float64())
	return
}

// Divides two values.
// 	float:div a b
// Returns a / b
func Command_Div(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("float:div", "2")
	}

	script.RetVal = rex.NewValueFloat64(params[0].Float64() / params[1].Float64())
	return
}

// Multiplies two values.
// 	float:mul a b
// Returns a * b
func Command_Mul(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("float:mul", "2")
	}

	script.RetVal = rex.NewValueFloat64(params[0].Float64() * params[1].Float64())
	return
}

// Floating point greater than.
// 	float:gt a b
// Returns true or false.
func Command_Gt(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("float:gt", "2")
	}

	result := params[0].Float64() > params[1].Float64()
	if result {
		script.RetVal = rex.NewValueBool(true)
		return
	}
	script.RetVal = rex.NewValueBool(false)
}

// Floating point less than.
// 	float:lt a b
// Returns true or false.
func Command_Lt(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("float:lt", "2")
	}

	result := params[0].Float64() < params[1].Float64()
	if result {
		script.RetVal = rex.NewValueBool(true)
		return
	}
	script.RetVal = rex.NewValueBool(false)
}

// Floating point equal.
// 	float:eq a b
// Returns true or false.
func Command_Eq(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("float:eq", "2")
	}

	result := params[0].Float64() == params[1].Float64()
	if result {
		script.RetVal = rex.NewValueBool(true)
		return
	}
	script.RetVal = rex.NewValueBool(false)
}
