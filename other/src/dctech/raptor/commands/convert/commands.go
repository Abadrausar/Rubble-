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

// Raptor Type Conversion Commands.
package convert

import "dctech/raptor"

// Adds the type conversion commands to the state.
// The type conversion commands are:
//	convert:int
//	convert:float
//	convert:bool
//	convert:string
//	convert:command
//	convert:escape
func Setup(state *raptor.State) {
	state.NewNameSpace("convert")
	state.NewNativeCommand("convert:int", CommandConvert_Int)
	state.NewNativeCommand("convert:float", CommandConvert_Float)
	state.NewNativeCommand("convert:bool", CommandConvert_Bool)
	state.NewNativeCommand("convert:string", CommandConvert_String)
	state.NewNativeCommand("convert:command", CommandConvert_Command)
	state.NewNativeCommand("convert:escape", CommandConvert_Escape)
}

// Forces a value to type int.
// Not all values will produce useful results.
// 	convert:int value
// Returns the converted value.
func CommandConvert_Int(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to convert:int.")
	}

	state.RetVal = raptor.NewValueInt64(params[0].Int64())
}

// Forces a value to type float.
// Not all values will produce useful results.
// 	convert:float value
// Returns the converted value.
func CommandConvert_Float(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to convert:float.")
	}

	state.RetVal = raptor.NewValueFloat64(params[0].Float64())
}

// Forces a value to type bool.
// Not all values will produce useful results.
// 	convert:bool value
// Returns the converted value.
func CommandConvert_Bool(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to convert:bool.")
	}

	state.RetVal = raptor.NewValueBool(params[0].Bool())
}

// Forces a value to type string.
// Not all values will produce useful results.
// 	convert:string value
// Returns the converted value.
func CommandConvert_String(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to convert:string.")
	}

	state.RetVal = raptor.NewValueString(params[0].String())
}

// Forces a value to type command.
// Not all values will produce useful results.
// 	convert:command value
// Returns the converted value.
func CommandConvert_Command(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to convert:command.")
	}

	state.RetVal = raptor.NewValueCommand(params[0].String())
}

// Converts a value into a string that you can write out for parsing later.
// Not all values will produce useful results.
// 	convert:escape value
// Returns a properly escaped string representaion of value.
func CommandConvert_Escape(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to convert:escape.")
	}

	state.RetVal = raptor.NewValueString(params[0].CodeString())
}
