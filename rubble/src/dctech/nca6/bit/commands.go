// NCA v6 Bitwise Commands.
package bit

import "dctech/nca6"

// Adds the bitwise commands to the state.
// The bitwise commands are:
//	lshift
//	rshift
//	and
//	or
//	xor
//	not
func Setup(state *nca6.State) {
	state.NewNativeCommand("lshift", CommandLShift)
	state.NewNativeCommand("rshift", CommandRShift)
	state.NewNativeCommand("and", CommandAnd)
	state.NewNativeCommand("or", CommandOr)
	state.NewNativeCommand("xor", CommandXor)
	state.NewNativeCommand("not", CommandNot)
}

// Shifts the value left.
// Count values of 0 or less are given the value 1.
// 	lshift a [count]
// Returns the param shifted "count" places
func CommandLShift(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic("Wrong number of params to lshift.")
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

	state.RetVal = nca6.NewValueInt64(params[0].Int64() << opp2)
	return
}

// Shifts the value right.
// Count values of 0 or less are given the value 1.
// 	rshift a [count]
// Returns the param shifted "count" places
func CommandRShift(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic("Wrong number of params to rshift.")
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

	state.RetVal = nca6.NewValueInt64(params[0].Int64() >> opp2)
	return
}

// Ands two values.
// 	and a b
// Returns a & b
func CommandAnd(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to and.")
	}

	state.RetVal = nca6.NewValueInt64(params[0].Int64() & params[1].Int64())
	return
}

// Ors two values.
// 	or a b
// Returns a | b
func CommandOr(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to or.")
	}

	state.RetVal = nca6.NewValueInt64(params[0].Int64() | params[1].Int64())
	return
}

// Xors two values.
// 	xor a b
// Returns a ^ b
func CommandXor(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to xor.")
	}

	state.RetVal = nca6.NewValueInt64(params[0].Int64() ^ params[1].Int64())
	return
}

// Inverts a value.
// 	not a
// Returns a ^ -1
func CommandNot(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to not.")
	}

	state.RetVal = nca6.NewValueInt64(params[0].Int64() ^ -1)
	return
}
