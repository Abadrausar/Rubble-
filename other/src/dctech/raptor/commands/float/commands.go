/*
Copyright 2012-2014 by Milo Christiansen

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

// Raptor Floating Point Commands.
package float

import "dctech/raptor"

// Adds the floating point commands to the state.
// The floating point commands are:
//	float:add
//	float:sub
//	float:div
//	float:mul
//	float:gt
//	float:lt
//	float:eq
func Setup(state *raptor.State) {
	state.NewNameSpace("float")
	state.NewNativeCommand("float:add", CommandFloat_Add)
	state.NewNativeCommand("float:sub", CommandFloat_Sub)
	state.NewNativeCommand("float:div", CommandFloat_Div)
	state.NewNativeCommand("float:mul", CommandFloat_Mul)
	state.NewNativeCommand("float:gt", CommandFloat_Gt)
	state.NewNativeCommand("float:lt", CommandFloat_Lt)
	state.NewNativeCommand("float:eq", CommandFloat_Eq)
}

// Adds two values.
// 	float:add a b
// Returns a + b
func CommandFloat_Add(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to float:add.")
	}

	script.RetVal = raptor.NewValueFloat64(params[0].Float64() + params[1].Float64())
	return
}

// Subtracts two values.
// 	float:sub a b
// Returns a - b
func CommandFloat_Sub(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to float:sub.")
	}

	script.RetVal = raptor.NewValueFloat64(params[0].Float64() - params[1].Float64())
	return
}

// Divides two values.
// 	float:div a b
// Returns a / b
func CommandFloat_Div(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to float:div.")
	}

	script.RetVal = raptor.NewValueFloat64(params[0].Float64() / params[1].Float64())
	return
}

// Multiplies two values.
// 	float:mul a b
// Returns a * b
func CommandFloat_Mul(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to float:mul.")
	}

	script.RetVal = raptor.NewValueFloat64(params[0].Float64() * params[1].Float64())
	return
}

// Floating point greater than.
// 	float:gt a b
// Returns true or false.
func CommandFloat_Gt(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to float:gt.")
	}

	result := params[0].Float64() > params[1].Float64()
	if result {
		script.RetVal = raptor.NewValueBool(true)
		return
	}
	script.RetVal = raptor.NewValueBool(false)
}

// Floating point less than.
// 	float:lt a b
// Returns true or false.
func CommandFloat_Lt(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to float:lt.")
	}

	result := params[0].Float64() < params[1].Float64()
	if result {
		script.RetVal = raptor.NewValueBool(true)
		return
	}
	script.RetVal = raptor.NewValueBool(false)
}

// Floating point equal.
// 	float:eq a b
// Returns true or false.
func CommandFloat_Eq(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to float:eq.")
	}

	result := params[0].Float64() == params[1].Float64()
	if result {
		script.RetVal = raptor.NewValueBool(true)
		return
	}
	script.RetVal = raptor.NewValueBool(false)
}
