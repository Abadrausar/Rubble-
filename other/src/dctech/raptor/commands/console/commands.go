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

// Raptor Console IO Commands.
package console

import "dctech/raptor"

// Adds the console IO commands to the state.
// The console IO commands are:
//	console:print
func Setup(state *raptor.State) {
	state.NewNameSpace("console")
	state.NewNativeCommand("console:print", CommandConsole_Print)
}

// Print a list of values.
// 	print [value...]
// Returns unchanged.
func CommandConsole_Print(state *raptor.State, params []*raptor.Value) {
	for _, val := range params {
		state.Print(val.String())
	}
}
