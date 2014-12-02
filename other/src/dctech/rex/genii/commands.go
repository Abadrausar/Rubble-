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

package genii

import "dctech/rex"
import "reflect"

// Adds the GenII helper commands to the state.
// The GenII helper commands are:
//	genii:bytes_string
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
	
	mod := state.RegisterModule("genii")
	mod.RegisterCommand("bytes_string", Command_BytesToString)
	mod.RegisterCommand("string_bytes", Command_StringToBytes)
	
	return nil
}

// Converts a GenII byte slice to a string.
// 	genii:bytes_string GenII:Array
// Returns the string.
func Command_BytesToString(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("genii:bytes_string", "1")
	}
	if _, ok := params[0].Data.(*Array); !ok {
		rex.ErrorGeneralCmd("genii:bytes_string", "Parameter 0 is not a GenII:Array.")
	}
	val := reflect.Value(*params[0].Data.(*Array))
	var tmp []byte
	if val.Type() != reflect.TypeOf(tmp) {
		rex.ErrorGeneralCmd("genii:bytes_string", "Parameter 0 is not backed by a byte slice.")
	}
	
	script.RetVal = rex.NewValueString(string(val.Bytes()))
	return
}

// Converts a string to a byte slice.
// 	genii:string_bytes string
// Returns the byte slice in a user value.
func Command_StringToBytes(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("genii:string_bytes", "1")
	}
	
	script.RetVal = rex.NewValueUser([]byte(params[0].String()))
	return
}
