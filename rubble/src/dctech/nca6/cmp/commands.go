// NCA v6 Integer Comparison Commands.
package cmp

import "dctech/nca6"

// Adds the integer comparison commands to the state.
// The integer comparison commands are:
//	gt
//	lt
//	eq
func Setup(state *nca6.State) {
	state.NewNativeCommand("gt", CommandGt)
	state.NewNativeCommand("lt", CommandLt)
	state.NewNativeCommand("eq", CommandEq)
}

// Integer greater than.
// 	gt a b
// Returns 0 or -1
func CommandGt(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to gt.")
	}

	result := params[0].Int64() > params[1].Int64()
	if result {
		state.RetVal = nca6.NewValueInt64(-1)
		return
	}
	state.RetVal = nca6.NewValueInt64(0)
	return
}

// Integer less than.
// 	lt a b
// Returns 0 or -1
func CommandLt(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to lt.")
	}

	result := params[0].Int64() < params[1].Int64()
	if result {
		state.RetVal = nca6.NewValueInt64(-1)
		return
	}
	state.RetVal = nca6.NewValueInt64(0)
	return
}

// Integer equal.
// 	eq a b
// Returns 0 or -1
func CommandEq(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to eq.")
	}

	result := params[0].Int64() == params[1].Int64()
	if result {
		state.RetVal = nca6.NewValueInt64(-1)
		return
	}
	state.RetVal = nca6.NewValueInt64(0)
	return
}
