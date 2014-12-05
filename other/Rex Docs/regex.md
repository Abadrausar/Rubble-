PACKAGE DOCUMENTATION

package regex
    import "dctech/rex/commands/regex"



FUNCTIONS

func Command_Replace(script *rex.Script, params []*rex.Value)
    Runs a regular expression search and replace.

	regex:replace regex input replace

    Returns input with all strings matching regex replaced with replace.


