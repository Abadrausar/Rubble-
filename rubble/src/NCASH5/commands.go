package main

import "dctech/nca5"

// Set the return value to nil.
// 	clrret
// Returns nil.
func CommandClrRet(state *nca5.State, params []*nca5.Value) {
	state.RetVal = nil
}
