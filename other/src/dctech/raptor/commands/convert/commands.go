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

// Raptor Type Conversion Commands.
package convert

import "dctech/raptor"

// Adds the type conversion commands to the state.
// The type conversion commands are:
//	convert:type
//	convert:int
//	convert:float
//	convert:bool
//	convert:string
//	convert:command
//	convert:escape
func Setup(state *raptor.State) {
	state.NewNameSpace("convert")
	state.NewNativeCommand("convert:type", CommandConvert_Type)
	state.NewNativeCommand("convert:int", CommandConvert_Int)
	state.NewNativeCommand("convert:float", CommandConvert_Float)
	state.NewNativeCommand("convert:bool", CommandConvert_Bool)
	state.NewNativeCommand("convert:string", CommandConvert_String)
	state.NewNativeCommand("convert:command", CommandConvert_Command)
	state.NewNativeCommand("convert:escape", CommandConvert_Escape)
}

// Returns a values internal type.
// 	convert:type value
// Returns the values type.
func CommandConvert_Type(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to convert:int.")
	}

	script.RetVal = raptor.NewValueString(params[0].TypeString())
}

// Forces a value to type int.
// Not all values will produce useful results.
// 	convert:int value
// Returns the converted value.
func CommandConvert_Int(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to convert:int.")
	}

	script.RetVal = raptor.NewValueInt64(params[0].Int64())
}

// Forces a value to type float.
// Not all values will produce useful results.
// 	convert:float value
// Returns the converted value.
func CommandConvert_Float(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to convert:float.")
	}

	script.RetVal = raptor.NewValueFloat64(params[0].Float64())
}

// Forces a value to type bool.
// Not all values will produce useful results.
// 	convert:bool value
// Returns the converted value.
func CommandConvert_Bool(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to convert:bool.")
	}

	script.RetVal = raptor.NewValueBool(params[0].Bool())
}

// Forces a value to type string.
// Not all values will produce useful results.
// 	convert:string value
// Returns the converted value.
func CommandConvert_String(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to convert:string.")
	}

	script.RetVal = raptor.NewValueString(params[0].String())
}

// Forces a value to type command.
// Not all values will produce useful results.
// 	convert:command value
// Returns the converted value.
func CommandConvert_Command(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to convert:command.")
	}

	script.RetVal = raptor.NewValueCommand(params[0].String())
}

// Converts a value into a string that you can write out for parsing later.
// Not all values will produce useful results.
// 	convert:escape value
// Returns a properly escaped string representaion of value.
func CommandConvert_Escape(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to convert:escape.")
	}

	script.RetVal = raptor.NewValueString(params[0].CodeString())
}
