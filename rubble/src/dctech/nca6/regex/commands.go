// NCA v6 Regular Expression Commands.
package regex

import "dctech/nca6"
//import "fmt"
import "regexp"

// Adds the regex commands to the state.
// The regex commands are:
//	regex:replace
func Setup(state *nca6.State) {
	state.NewNameSpace("regex")
	state.NewNativeCommand("regex:replace", CommandRegEx_Replace)
}

// Runs a regular expression search and replace.
// 	regex:replace regex input replace
// Returns input with all strings matching regex replaced with replace.
func CommandRegEx_Replace(state *nca6.State, params []*nca6.Value) {
	if len(params) != 3 {
		panic("Wrong number of params to regex:replace.")
	}

	regEx := regexp.MustCompile(params[0].String())
	state.RetVal = nca6.NewValueString(regEx.ReplaceAllString(params[1].String(), params[2].String()))
}
