/*
Copyright 2012-2014 by Milo Christiansen

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
Raptor scripting language.

raptor n.
A dromaeosaurid dinosaur, especially a velociraptor or utahraptor.

Raptor is a simple command language similar in concept to TCL.
Raptor is based on the syntax of NCA but is faster. Not all NCA syntax exists in Raptor, for example 
curly brackets '{}' are not string delimiters but are used for the new type, code blocks. 
There are also many new features like new types and some basic OOP support.

The following is a series of semi-random notes that didn't really fit anywhere else.

Error message position info may (will) be wrong. The line number is almost always correct
but the column will mostly be a little after the problem.
Tabs will tend to throw column info off a little, just remember that a tab is counted as one column.
If lexing a double-quote string (that is being used as code) that has \n escape sequences
the position will not match the source.

Error positions generated from compiled (binary) scripts are almost useless as position 
is given as a number of tokens from the beginning of the current block. Unless you know what 
block the error is in the position is nearly useless. Best bet is to restart with a source 
version of the script and re-induce the error.

Valuable debugging information can sometimes be acquired by causing your error to happen both normally and 
with State.NoRecover set to false. Note that this does require some knowledge of the Raptor internals.

Note that any commands or internal functions that say that they may "panic with an error" have all such 
panics recovered before they can reach any user code (unless State.NoRecover is true).

*/
package raptor
