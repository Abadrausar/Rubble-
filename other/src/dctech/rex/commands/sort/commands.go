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

// Rex Ordered Map Commands.
package sort

import "dctech/rex"

// Adds the ordered map commands to the state.
// The ordered map commands are:
//	sort:map
//	sort:new
// In addition adds the following indexable type:
//	sort:map
// (command and type names do not conflict)
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
	
	mod := state.RegisterModule("sort")
	mod.RegisterCommand("map", Command_Map)
	mod.RegisterCommand("new", Command_New)
	
	mod.RegisterType("map", NewOrderedMapFromLit)
	
	return nil
}

// Sorts a map alphabetically.
// 	sort:map map
// Returns a new ordered map with the same keys/values as the old indexable.
func Command_Map(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("sort:map", "1")
	}

	if params[0].Type != rex.TypIndex {
		rex.ErrorGeneralCmd("sort:map", "Parameter 0 is not Indexable.")
	}
	script.RetVal = rex.NewValueIndex(NewOrderedMapFromIndexable(params[0].Data.(rex.Indexable)))
}

// Creates a new ordered map.
// 	sort:new
// Returns a new (empty) ordered map.
func Command_New(script *rex.Script, params []*rex.Value) {
	script.RetVal = rex.NewValueIndex(NewOrderedMap())
}
