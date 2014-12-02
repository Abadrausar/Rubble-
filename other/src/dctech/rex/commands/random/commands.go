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

// Rex Random Number Commands.
package random

import "dctech/rex"
import "math/rand"

// Adds the random number commands to the state.
// The random number commands are:
//	rand:new
//	rand:int
//	rand:float
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
	
	mod := state.RegisterModule("rand")
	mod.RegisterCommand("new", Command_New)
	mod.RegisterCommand("int", Command_Int)
	mod.RegisterCommand("float", Command_Float)
	
	return nil
}

// Creates a new random number generator.
// 	rand:new seed
// Returns a handle to the new random number generator.
func Command_New(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("rand:new", "1")
	}
	
	script.RetVal = rex.NewValueUser(rand.New(rand.NewSource(params[0].Int64())))
	return
}

// Reads an int from the specified generator.
// 	rand:int gen
// Returns a non-negative int value.
func Command_Int(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("rand:int", "1")
	}
	gen, ok := params[0].Data.(*rand.Rand)
	if !ok {
		rex.ErrorGeneralCmd("rand:int", "Param 0 is not an random number generator.")
	}
	
	script.RetVal = rex.NewValueInt64(gen.Int63())
	return
}

// Reads a float from the specified generator.
// 	rand:float gen
// Returns a float value between 0.0 and 1.0.
func Command_Float(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("rand:float", "1")
	}
	gen, ok := params[0].Data.(*rand.Rand)
	if !ok {
		rex.ErrorGeneralCmd("rand:float", "Param 0 is not an random number generator.")
	}
	
	script.RetVal = rex.NewValueFloat64(gen.Float64())
	return
}
