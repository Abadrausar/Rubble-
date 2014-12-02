/*
NCA SHell for NCA v5.

Run "ncash5 -h" for usage information.

NCASH's default mode is interactive mode, in which you can type NCA code and run it instantly, but there is also a batch mode where NCASH will run code from a file and then exit. 

If using NCASH to debug custom commands you may want to use "-recover=false" to cause it to print a stack trace on error.
Note that to load a custom command you will need to edit NCASH as there is no way (AFAIK) to dynamicly link a Go program at this time.

NCASH has some custom commands:
	clrret			Set the return value to <nil>. Useful for clearing the state.
	valueinspect	See detailed information about a script value.

*/
package documentation


