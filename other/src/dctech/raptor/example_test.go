/*
For copyright/license see header in file "doc.go"
*/

package raptor_test

import "dctech/raptor"
import "fmt"

func Example() {

	// Create a new state and script
	state := raptor.NewState()
	script := raptor.NewScript()

	// At this point you would add any native commands or other such things to the state.

	// Add the program code to the script.
	script.Code.AddString(`
	# nothing for now....
	`, raptor.NewPosition(37, 28, "example_test.go"))

	// Run!
	// Note that Run removes the code it executes from the BlockStore 
	// so there is no need to clean up before using the script again.
	rtn, err := state.Run(script)
	if err != nil {
		// If we fell flat on our face make sure to make a note of it.
		fmt.Println("Error:", err)
	}
	// And in any case write the scripts return value.
	fmt.Println("Ret:", rtn)

	// Output: Ret: nil
}
