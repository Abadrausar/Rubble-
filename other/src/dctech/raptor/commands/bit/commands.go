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

// Raptor Bitwise Commands.
package bit

import "dctech/raptor"

// Adds the bitwise commands to the state.
// The bitwise commands are:
//	bit:lshift
//	bit:rshift
//	bit:and
//	bit:or
//	bit:xor
//	bit:not
func Setup(state *raptor.State) {
	state.NewNameSpace("bit")
	state.NewNativeCommand("bit:lshift", CommandBit_LShift)
	state.NewNativeCommand("bit:rshift", CommandBit_RShift)
	state.NewNativeCommand("bit:and", CommandBit_And)
	state.NewNativeCommand("bit:or", CommandBit_Or)
	state.NewNativeCommand("bit:xor", CommandBit_Xor)
	state.NewNativeCommand("bit:not", CommandBit_Not)
}

// Shifts the value left.
// Count values of 0 or less are given the value 1.
// 	bit:lshift a [count]
// Returns the param shifted "count" places
func CommandBit_LShift(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic("Wrong number of params to bit:lshift.")
	}

	opp2 := uint64(1)
	if len(params) > 1 {
		tmp := uint64(params[1].Int64())
		if tmp < 1 {
			opp2 = 1
		} else {
			opp2 = tmp
		}
	}

	script.RetVal = raptor.NewValueInt64(params[0].Int64() << opp2)
	return
}

// Shifts the value right.
// Count values of 0 or less are given the value 1.
// 	bit:rshift a [count]
// Returns the param shifted "count" places
func CommandBit_RShift(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic("Wrong number of params to bit:rshift.")
	}

	opp2 := uint64(1)
	if len(params) > 1 {
		tmp := uint64(params[1].Int64())
		if tmp < 1 {
			opp2 = 1
		} else {
			opp2 = tmp
		}
	}

	script.RetVal = raptor.NewValueInt64(params[0].Int64() >> opp2)
	return
}

// Ands two values.
// 	bit:and a b
// Returns a & b
func CommandBit_And(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to bit:and.")
	}

	script.RetVal = raptor.NewValueInt64(params[0].Int64() & params[1].Int64())
	return
}

// Ors two values.
// 	bit:or a b
// Returns a | b
func CommandBit_Or(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to bit:or.")
	}

	script.RetVal = raptor.NewValueInt64(params[0].Int64() | params[1].Int64())
	return
}

// Xors two values.
// 	bit:xor a b
// Returns a ^ b
func CommandBit_Xor(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to bit:xor.")
	}

	script.RetVal = raptor.NewValueInt64(params[0].Int64() ^ params[1].Int64())
	return
}

// Inverts a value.
// 	bit:not a
// Returns a ^ -1
func CommandBit_Not(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to bit:not.")
	}

	script.RetVal = raptor.NewValueInt64(params[0].Int64() ^ -1)
	return
}
