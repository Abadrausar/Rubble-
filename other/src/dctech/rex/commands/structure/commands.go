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

// Rex "struct" Type.
// 
// The struct type is basically a map indexable that has a fixed set of keys.
// The keys are determined based on the value passed to the "proto" key (which
// MUST be first, and MUST be an indexable), any keys which are not given a value
// take their value from the prototype indexable.
// 
// To put it simply a struct is a copy of a map that cannot have new keys added.
// 
// At some point I may allow more advanced behavior via special indexables for
// use in the prototype, but for now things are really simple.
package structure

import "dctech/rex"

// Adds the struct type to the state.
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
	
	state.RegisterType("struct", NewStructureFromLit)
	
	return nil
}
