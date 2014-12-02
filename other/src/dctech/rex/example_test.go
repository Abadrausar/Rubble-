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

package rex_test

import "dctech/rex"
import "dctech/rex/commands/base"
import "fmt"

func Example() {
	// Everything starts with a state and a script.
	// You shouldn't need more than one state in most cases, as it can be shared by many scripts.
	// The script is where all the "local" data is. You will need one for each concurrent block of
	// code. Scripts are reusable.
	state := rex.NewState()
	script := rex.NewScript()

	// Load the base commands.
	base.Setup(state)

	// If you want the standard indexables you would register them here
	
	// Most of the following could be replaced with a single call to state.CompileAndRun
	// (in which case the script above would not be needed)
	
	// Compile the code, this function returns a script value.
	// It should be possible to inject code into a running script with this function,
	// so long as the running script knows where to find the injected code :)
	val, err := state.CompileToValue(`(ret true)`, rex.NewPosition(1, 1, ""))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	// Now run the value.
	// Note that the procedure for an interactive shell is a little different
	// as you need to worry about preserving local variables across calls.
	ret, err := state.Run(script, val)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Ret:", ret)

	// Output: Ret: true
}
