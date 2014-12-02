package nca5_test

import "dctech/nca5"
import "dctech/nca5/base"
import "dctech/nca5/conio"
import "fmt"

func Example() {
	
	// Create a new state
	state := nca4.NewState()
	
	// Load a few of the predefined commands
	// Note that each set of commands are in their own package.
	base.Setup(state)
	conio.Setup(state)
	
	// Add the program code to the state.
	// console:print is from dctech/nca4/conio
	// ret is from dctech/nca4/base
	// You will always want to load base as it has the commands
	// needed to make nca usable.
	state.Code.Add(`(console:print Testing... "\n")(ret 0)`)
	
	// Run!
	// Note that Run and its subset Exec remove the code they execute from the BlockStore 
	// so there is no need to clean up before using the state again.
	rtn, err := state.Run()
	if err != nil {
		// If we fell flat on our face make sure to make a note of it.
		fmt.Println("Error:", err)
	}
	// And in any case write the scripts return value.
	fmt.Println("Ret:", rtn)
	
	// Output: Testing...
	// Ret: 0
}
