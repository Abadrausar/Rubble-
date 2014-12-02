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

// Rex Console IO Commands.
package console

import "dctech/rex"

// Adds the console IO commands to the state.
// The console IO commands are:
//	console:print
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
	
	mod := state.RegisterModule("console")
	mod.RegisterCommand("print", Command_Print)
	
	return nil
}

// Print a list of values.
// 	console:print [value...]
// Returns unchanged.
func Command_Print(script *rex.Script, params []*rex.Value) {
	for _, val := range params {
		script.Print(val.String())
	}
}
