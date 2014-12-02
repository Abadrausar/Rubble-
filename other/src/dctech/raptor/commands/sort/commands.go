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

// Raptor Ordered Map Commands.
package sort

import "dctech/raptor"

// Adds the ordered map commands to the state.
// The ordered map commands are:
//	sort:map
//	sort:new
func Setup(state *raptor.State) {
	state.NewNameSpace("sort")
	state.NewNativeCommand("sort:map", CommandSort_Map)
	state.NewNativeCommand("sort:new", CommandSort_New)
}

// Sorts a map alphabetically.
// 	sort:map map
// Returns a new ordered map with the same keys/values as the old one.
func CommandSort_Map(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	value := params[0].Indexable()
	if value == nil {
		panic(script.GeneralCmdError("Parameter 0 is not Indexable."))
	}
	script.RetVal = raptor.NewValueObject(NewOrderedMapFromIndexable(value))
}

// Creates a new ordered map.
// 	sort:new
// Ordered maps are also thread safe.
// Returns a new (empty) ordered map.
func CommandSort_New(script *raptor.Script, params []*raptor.Value) {
	script.RetVal = raptor.NewValueObject(NewOrderedMap())
}
