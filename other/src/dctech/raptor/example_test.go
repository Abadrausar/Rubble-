/*
Copyright 2012-2013 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

package raptor_test

import "dctech/raptor"
import "fmt"

func Example() {

	// Create a new state
	state := raptor.NewState()

	// At this point you would add any native commands or other such things to the state.

	// Add the program code to the state.
	state.Code.Add(`
	#nothing for now....
	`)

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

	// Output: Ret: <nil>
}
