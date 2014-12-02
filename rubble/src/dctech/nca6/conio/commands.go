// NCA v6 Console IO Commands.
package conio

import "fmt"
import "dctech/nca6"

// Adds the console IO commands to the state.
// The console IO commands are:
//	console:print
func Setup(state *nca6.State) {
	state.NewNameSpace("console")
	state.NewNativeCommand("console:print", CommandConsole_Print)
}

// Print a list of values.
// 	print [value...]
// Returns unchanged.
func CommandConsole_Print(state *nca6.State, params []*nca6.Value) {
	for _, val := range params {
		fmt.Print(val.String())
	}
}
