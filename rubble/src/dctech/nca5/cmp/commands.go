// NCA v5 Comparison Commands.
package cmp

import "dctech/nca5"
import "strconv"

// Adds the comparison commands to the state.
// The comparison commands are:
//	gt
//	lt
//	eq
func Setup(state *nca5.State) {
	state.NewNativeCommand("gt", CommandGt)
	state.NewNativeCommand("lt", CommandLt)
	state.NewNativeCommand("eq", CommandEq)
}

// Greater than.
// If the value can be converted to int64 that is the value used otherwise the length is used.
// 	gt a b
// Returns "0" or "-1"
func CommandGt(state *nca5.State, params []*nca5.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to gt.")
	}

	opp1, err := strconv.ParseInt(params[0].String(), 0, 64)
	if err != nil {
		opp1 = int64(len(params[0].String()))
	}
	opp2, err := strconv.ParseInt(params[1].String(), 0, 64)
	if err != nil {
		opp2 = int64(len(params[1].String()))
	}

	result := opp1 > opp2
	if result {
		state.RetVal = nca5.NewValue("-1")
		return
	}
	state.RetVal = nca5.NewValue("0")
	return
}

// Less than.
// If the value can be converted to int64 that is the value used otherwise the length is used.
// 	lt a b
// Returns "0" or "-1"
func CommandLt(state *nca5.State, params []*nca5.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to lt.")
	}

	opp1, err := strconv.ParseInt(params[0].String(), 0, 64)
	if err != nil {
		opp1 = int64(len(params[0].String()))
	}
	opp2, err := strconv.ParseInt(params[1].String(), 0, 64)
	if err != nil {
		opp2 = int64(len(params[1].String()))
	}

	result := opp1 < opp2
	if result {
		state.RetVal = nca5.NewValue("-1")
		return
	}
	state.RetVal = nca5.NewValue("0")
	return
}

// Equal.
// Does a simple string compare.
// Make sure numeric values are normalized. If you try to load a hexadecimal constant 
// and compare it to a existing number it will not work correctly. 
// A simple (set numbervar (add [numbervar] 0)) should work.
// 	eq a b
// Returns "0" or "-1"
func CommandEq(state *nca5.State, params []*nca5.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to eq.")
	}

	result := params[0].String() == params[1].String()
	if result {
		state.RetVal = nca5.NewValue("-1")
		return
	}
	state.RetVal = nca5.NewValue("0")
	return
}
