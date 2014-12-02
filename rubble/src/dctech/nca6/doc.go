// NCA Version 6.
// 
// No Clever Acronym :p.
// 
// NCA is a simple command launguage similar in concept to TCL.
// 
// What's new from v5:
//	Script Values may now hold arbitrary data, just make sure that your commands know how to handle it.
//	Script code may be "compiled" into a pre-lexed version for storage, this is used by user commands and the like.
//	Error messages no longer have the column, this was (almost) always wrong anyways.
//	The Map and Array Indexables are now part of the base command package.
//	The special array "params" is now write protected.
//	The way strings are converted to booleans has changed.
//	The way the element count of an Indexable is fetched has changed.
// 
package nca6
