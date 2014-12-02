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

/*
NCA Version 7.

No Clever Acronym :p.

NCA is a simple command launguage similar in concept to TCL.

Error message position info may be wrong. The line is almost always correct
but the column will mostly be a little after the problem.
If lexing a double-quote string (that is being used as code) that has \n escape sequences
the position will not match the source.
In any case the error reporting is light-years better than that of version 3 and before.

Note that any commands or internal functions that say that they may "panic with an error" have all such panics recovered before they can reach any user code (unless State.NoRecover is true).

Changelog:
This gets a little sketchy for the older versions, I didn't keep good records before version 6
	v7
		Better support for semi and fully readonly Indexables
			Indexables may now reject some or all writes
		Multiple level indexing syntax changed
			[test index1 index2]
		Better type support
			new type: bool
			All types (except object) now parsed when calling TokenToValue
		Indexable literal syntax
			Example:
				(var test <map, akey="a value", also_a_key=<array, one, two, three>>)
		New Lexer. The old lexer needed to be compatible with version 2 and that made it nearly impossible to overhaul.
		Any expression that yields a value may be used in the global scope, previously only commands were allowed.
		Column information restored to error messages, the column will display as 0 if data is not avalible.
	
	v6
		Script Values may now hold arbitrary data, just make sure that your commands know how to handle it.
		Script code may be "compiled" into a pre-lexed version for storage, this is used by user commands and the like.
		Error messages no longer have the column, this was (almost) always wrong anyways.
		The Map and Array Indexables are now part of the base command package.
		The special array "params" is now write protected.
		The way strings are converted to booleans has changed.
		The way the element count of an Indexable is fetched has changed.
		Exit now a flag in the state, this is far more flexable than the old way
	
	v5
		New array type
		Maps (and arrays) are now treated like Values
	
	v4
		GDM and Parser combined into State
		Native command handler prototype changed to fit above change
	
	v3
		return value handling changed (now stored in GDM)
		error flag added to GDM

	v2
		array renamed map
		optional param names
		simple param checking (the count is verifyed (only if param names are used))
		command handling reworked. all commands are now added to the global data manager, user commands are no longer normal variables
		namespace. A namespace is more-or-less just a named root environment

	v1
		The first version

*/
package nca7
