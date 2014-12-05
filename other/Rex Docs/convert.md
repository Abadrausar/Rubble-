PACKAGE DOCUMENTATION

package convert
    import "dctech/rex/commands/convert"



FUNCTIONS

func Command_Bool(script *rex.Script, params []*rex.Value)
    Forces a value to type bool. Not all values will produce useful results.

	convert:bool value

    Returns the converted value.

func Command_Float(script *rex.Script, params []*rex.Value)
    Forces a value to type float. Not all values will produce useful
    results.

	convert:float value

    Returns the converted value.

func Command_Int(script *rex.Script, params []*rex.Value)
    Forces a value to type int. Not all values will produce useful results.

	convert:int value

    Returns the converted value.

func Command_String(script *rex.Script, params []*rex.Value)
    Forces a value to type string. Not all values will produce useful
    results.

	convert:string value

    Returns the converted value.


