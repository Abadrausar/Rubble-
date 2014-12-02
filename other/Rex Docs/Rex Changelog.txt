
v1.7
	Made it possible for a command created with the "command" keyword to refer to itself (recursive commands here I come).
	GenII indexables now return nil for non-existent keys (like all indexables should).
	Fixed a major bug with the IntEditIndexable type (had no GetI method).
	Added new command sort:array, sort an IntEditIndexable in place.
	Added int:++ and int:-- commands.

v1.6
	Fixed some incorrect error messages and incorrectly checked error conditions in indexing variable operator execution.
		There were no false positives, but some errors were not caught at the right spot (or at all) so error messages could make little to no sense.
	Added an extension to Indexable to allow faster value lookup in integer indexed Indexables.
		If you have array-like Indexables change them to implement IntIndexable or IntEditIndexable for improved performance.
	array and sarray now implement IntEditIndexable and IntIndexable.
	GenII array types now implement IntEditIndexable and IntIndexable.
	foreach now uses IntIndexable if the indexable implements it.
	The indexing variable operator uses IntIndexable or IntEditIndexable as appropriate.
		For the variable operator to use these the indexable must implement them AND the key must be type int.
	The debug:value command now prints the value's raw data with the go syntax format specifier.
	Added base commands modval, breakif and breakloopif.
	Added basic support for error type meta-data.
	Value types are now stored in a custom type (to encourage use of the type constants).

v1.5
	New command parameter handling, including two new bits of syntax!
		Variadic parameters are MUCH improved, it is now possible to specify a variadic param at the end of a list of normal params.
			Variadic params are now named! to get the old behavior use "params=..." rather than "..." for the parameter list.
		You may now specify (optional) parameters out of order (keyed by name).
			Required parameters still need to be filled by position, but it is now possible to choose arbitrary optional parameters without needing to specify any optional parameter that may have preceded them.
	Removed all support for "code strings", this includes removing the requirement for a CodeString method from indexables and the convert:escape command.
	All global data is now stored in the "global" module.
	Any failed local variable lookups will now fail over to looking in the "global" module.
		It is still possible to refer to global directly via name.
	Added bool:sand and bool:sor, short-circuit evaluation boolean AND and OR operations.
		The parameters must be code blocks to allow delaying evaluation.
		When used with simple values they are SLOWER than the normal versions!

v1.4
	"convert:type" removed, replaced with more flexible "type" command

v1.3
	GenII can now do direct replacements for type index values as well as user values.
	All types of global data may now be redeclared.

v1.2
	The command call syntax now has a by-value version.
	Removed the "run" command, it is obsolete with by-value command calls.
	Non-module commands are now stored in a predefined module named "global".
	Fixed a long standing issue that cased bad positions to be printed for runtime errors.
	Most error positions now point to the first token of the problem statement instead of the last.

v1.1
	It should now be possible to replace GenII objects with others of exactly the same type that are stored in user values.
	GenII now comes with a simple command to create a byte slice from a string.

v1.0
	First version.
