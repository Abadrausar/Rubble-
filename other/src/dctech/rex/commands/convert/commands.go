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

// Rex Type Conversion Commands.
package convert

import "dctech/rex"

// Adds the type conversion commands to the state.
// The type conversion commands are:
//	convert:type
//	convert:int
//	convert:float
//	convert:bool
//	convert:string
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
	
	mod := state.RegisterModule("convert")
	mod.RegisterCommand("int", Command_Int)
	mod.RegisterCommand("float", Command_Float)
	mod.RegisterCommand("bool", Command_Bool)
	mod.RegisterCommand("string", Command_String)
	
	return nil
}

// Forces a value to type int.
// Not all values will produce useful results.
// 	convert:int value
// Returns the converted value.
func Command_Int(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("convert:int", "1")
	}

	script.RetVal = rex.NewValueInt64(params[0].Int64())
}

// Forces a value to type float.
// Not all values will produce useful results.
// 	convert:float value
// Returns the converted value.
func Command_Float(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("convert:float", "1")
	}

	script.RetVal = rex.NewValueFloat64(params[0].Float64())
}

// Forces a value to type bool.
// Not all values will produce useful results.
// 	convert:bool value
// Returns the converted value.
func Command_Bool(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("convert:bool", "1")
	}

	script.RetVal = rex.NewValueBool(params[0].Bool())
}

// Forces a value to type string.
// Not all values will produce useful results.
// 	convert:string value
// Returns the converted value.
func Command_String(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("convert:string", "1")
	}

	script.RetVal = rex.NewValueString(params[0].String())
}
