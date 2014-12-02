// NCA v5 Bitwise Commands.
package bit

import "dctech/nca5"

// Adds the bitwise commands to the state.
// The bitwise commands are:
//	lshift
//	rshift
//	and
//	or
//	xor
//	not
func Setup(state *nca5.State) {
	state.NewNativeCommand("lshift", CommandLShift)
	state.NewNativeCommand("rshift", CommandRShift)
	state.NewNativeCommand("and", CommandAnd)
	state.NewNativeCommand("or", CommandOr)
	state.NewNativeCommand("xor", CommandXor)
	state.NewNativeCommand("not", CommandNot)
}

// Shifts the value left.
// Opperand is converted to a 64 bit integer. Invalid strings are given the value "0"
// Count is converted to an uint8 with an invalid string getting the value "1".
// 	lshift a [count]
// Returns the param shifted "count" places
func CommandLShift(state *nca5.State, params []*nca5.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic("Wrong number of params to lshift.")
	}

	opp2 := uint8(1)
	if len(params) > 1 {
		tmp := params[1].UInt8()
		if tmp < 1 {
			opp2 = 1
		} else {
			opp2 = tmp
		}
	}

	state.RetVal = nca5.NewValueFromI64(params[0].Int64() << opp2)
	return
}

// Shifts the value right.
// Opperand is converted to a 64 bit integer. Invalid strings are given the value "0"
// Count is converted to an uint8 with an invalid string getting the value "1".
// 	rshift a [count]
// Returns the param shifted "count" places
func CommandRShift(state *nca5.State, params []*nca5.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic("Wrong number of params to rshift.")
	}

	opp2 := uint8(1)
	if len(params) > 1 {
		tmp := params[1].UInt8()
		if tmp < 1 {
			opp2 = 1
		} else {
			opp2 = tmp
		}
	}

	state.RetVal = nca5.NewValueFromI64(params[0].Int64() >> opp2)
	return
}

// Ands two values.
// Opperands are converted to 64 bit integers. Invalid strings are given the value "0"
// 	and a b
// Returns a & b
func CommandAnd(state *nca5.State, params []*nca5.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to and.")
	}

	state.RetVal = nca5.NewValueFromI64(params[0].Int64() & params[1].Int64())
	return
}

// Ors two values.
// Opperands are converted to 64 bit integers. Invalid strings are given the value "0"
// 	or a b
// Returns a | b
func CommandOr(state *nca5.State, params []*nca5.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to or.")
	}

	state.RetVal = nca5.NewValueFromI64(params[0].Int64() | params[1].Int64())
	return
}

// Xors two values.
// Opperands are converted to 64 bit integers. Invalid strings are given the value "0"
// 	xor a b
// Returns a ^ b
func CommandXor(state *nca5.State, params []*nca5.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to xor.")
	}

	state.RetVal = nca5.NewValueFromI64(params[0].Int64() ^ params[1].Int64())
	return
}

// Inverts a value.
// Opperand is converted to a 64 bit integer. Invalid strings are given the value "0"
// 	not a
// Returns a ^ -1
func CommandNot(state *nca5.State, params []*nca5.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to not.")
	}

	state.RetVal = nca5.NewValueFromI64(params[0].Int64() ^ -1)
	return
}
