/*
Copyright 2014 by Milo Christiansen

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

package genii_test

import "dctech/rex"
import "dctech/rex/genii"
import "fmt"

type T struct {
	a int
	B string
}

func Example() {
	state := rex.NewState()
	
	// Create our test data
	testStruct := &T{a: 100, B: "Test Data"}
	testII := genii.New(testStruct)
	
	// and add it to the state as "test:T"
	mod := state.RegisterModule("test")
	mod.RegisterVar("T", testII)
	
	// The script code.
	source := `
	# Should return 100, no real way to check without loading a bunch of commands.
	[test:T a]
	
	# Note that as "a" is not exported it's value is readonly.
	# It is an error to try to change the value of a readonly field.
	
	# Change the value of B.
	[test:T B = "Other Data"]
	
	# Clear the return register to make sure we get good results
	nil
	
	# and then get the current value of B.
	[test:T B]
	
	# Should return "Other Data".
	`
	
	// Just the normal compile-and-run song and dance...
	ret, err := state.CompileAndRun(source, "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Ret:", ret)
	
	// Check to make sure the script really wrote to the original struct.
	if testStruct.B != "Other Data" {
		fmt.Println("Error: Script changes not written to struct properly.")
	}

	// Output: Ret: Other Data
}
