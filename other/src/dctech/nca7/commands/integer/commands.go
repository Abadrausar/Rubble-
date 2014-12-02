/*
Copyright 2012-2013 by Milo Christiansen

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

// NCA v7 Integer Commands.
package integer

import "dctech/nca7"

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
func Setup(state *nca7.State) {
	state.NewNameSpace("int")
	state.NewNativeCommand("int:add", CommandInt_Add)
	state.NewNativeCommand("int:sub", CommandInt_Sub)
	state.NewNativeCommand("int:div", CommandInt_Div)
	state.NewNativeCommand("int:mod", CommandInt_Mod)
	state.NewNativeCommand("int:mul", CommandInt_Mul)
	state.NewNativeCommand("int:gt", CommandInt_Gt)
	state.NewNativeCommand("int:lt", CommandInt_Lt)
	state.NewNativeCommand("int:eq", CommandInt_Eq)
}

// Adds two values.
// 	int:add a b
// Returns a + b
func CommandInt_Add(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to int:add.")
	}

	state.RetVal = nca7.NewValueInt64(params[0].Int64() + params[1].Int64())
	return
}

// Subtracts two values.
// 	int:sub a b
// Returns a - b
func CommandInt_Sub(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to int:sub.")
	}

	state.RetVal = nca7.NewValueInt64(params[0].Int64() - params[1].Int64())
	return
}

// Divides two values.
// 	int:div a b
// Returns a / b
func CommandInt_Div(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to int:div.")
	}

	state.RetVal = nca7.NewValueInt64(params[0].Int64() / params[1].Int64())
	return
}

// Gives the remainder of dividing two values.
// 	int:mod a b
// Returns a % b
func CommandInt_Mod(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to int:mod.")
	}

	state.RetVal = nca7.NewValueInt64(params[0].Int64() % params[1].Int64())
	return
}

// Multiplies two values.
// 	int:mul a b
// Returns a * b
func CommandInt_Mul(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to int:mul.")
	}

	state.RetVal = nca7.NewValueInt64(params[0].Int64() * params[1].Int64())
	return
}


// Integer greater than.
// 	int:gt a b
// Returns true or false.
func CommandInt_Gt(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to int:gt.")
	}

	result := params[0].Int64() > params[1].Int64()
	if result {
		state.RetVal = nca7.NewValueBool(true)
		return
	}
	state.RetVal = nca7.NewValueBool(false)
}

// Integer less than.
// 	int:lt a b
// Returns true or false.
func CommandInt_Lt(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to int:lt.")
	}

	result := params[0].Int64() < params[1].Int64()
	if result {
		state.RetVal = nca7.NewValueBool(true)
		return
	}
	state.RetVal = nca7.NewValueBool(false)
}

// Integer equal.
// 	int:eq a b
// Returns true or false.
func CommandInt_Eq(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to int:eq.")
	}

	result := params[0].Int64() == params[1].Int64()
	if result {
		state.RetVal = nca7.NewValueBool(true)
		return
	}
	state.RetVal = nca7.NewValueBool(false)
}

