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

/*
Rex scripting language.

For a description of the Rex syntax see "Rex Basics.txt" in the docs directory.

Rex has two main parts, the Script and the State.

The Script handles:
	Internal execution of script code.
	Local variables.

The State handles:
	Compilation of scripts.
	External execution of script code.
	Storage of global data (modules, types, and commands)
	Error handling.

Generally you will only need one State, as it can be used by several scripts at once, the Script on the other hand
can only be used by one script at a time (as it's name suggests). The only thing the Script handles that the average
user will need to worry about is local variables, every other action is negotiated though the State unless you are 
writing a new native script command (in which case there are some methods exported from the Script that will be useful).

Code MUST be run on the same State that compiled it (or an identical one). This is because when a script is compiled it can change the State (command, module, and module variable declarations) and it needs to lookup indexes in the State. If
any of these indexes are different and/or if any expected global data is missing the script will almost certainly crash!

Most of the stuff that is exported is there for the use of external command packages and advanced applications, most of
the time if you don't know what it is for then keep your grubby fingers off!

Rex is still quite new, and so I am sure I have missed a bug or three, if you find something that does not work as
described feel free to let me know :)

If anyone want to do some bug fixing it can be very helpful to set State.NoRecover, this stops Rex from recovering
panics (and all errors start as a panic at some level) which makes it much easier to trace internal errors.
Of course NoRecover is only useful for debugging internal errors, script errors do not generally benefit from this 
setting (and it can actually make things harder, as then the state never gets a chance to fill in position information).

All panics are recovered and turned into errors before they reach user code except in a few cases:
	Command package Setup(state) functions only recover panics caused by Rex (type rex.ScriptError).

Send patches (*snort*) and bug reports (probably) to: milo.christiansen (at) gmail (dot) com

*/
package rex

/*
 Master TODO List:

	The debug commands need expanding.
	
	The ScriptError type should have some type flags that can be used to classify errors into categories.
	Some predefined ScriptError values may be useful.
	
	Look for cases where the error flag should be used but is not.
	
	I should write a good set of tests.
		Finish the base command tests.
		The various indexables need to be tested to make sure they are fully thread safe.
			I think there may be cases where I need to synchronize where I don't.
				The read-only indexables are not synchronized at all!
					As these are never written after creation thread safety should not be an issue.
					In any case test this assumption.
	
	There is a long list of stuff-to-do in the GenII doc.go file.
	
	It should be possible to make a generic function-to-command interface via reflection.
		Support for some complex types would be really hard to do, probably best to use the same rules as GenII.
		Such an interface would make writing new native commands much easier
			(provided performance doesn't suffer too much, which is a possibility)
		(low priority)
	
*/
