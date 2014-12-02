// NCA v5 Console IO Commands.
package conio

import "fmt"
import "dctech/nca5"

// Adds the console IO commands to the state.
// The console IO commands are:
//	console:print
func Setup(state *nca5.State) {
	state.NewNameSpace("console")
	state.NewNativeCommand("console:print", CommandConsole_Print)
}

// Print a list of values.
// 	print [value...]
// Returns unchanged.
func CommandConsole_Print(state *nca5.State, params []*nca5.Value) {
	for _, val := range params {
		fmt.Print(val.String())
	}
}
