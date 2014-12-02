PACKAGE DOCUMENTATION

package random
    import "dctech/rex/commands/random"



FUNCTIONS

func Command_Float(script *rex.Script, params []*rex.Value)
    Reads a float from the specified generator.

	rand:float gen

    Returns a float value between 0.0 and 1.0.

func Command_Int(script *rex.Script, params []*rex.Value)
    Reads an int from the specified generator.

	rand:int gen

    Returns a non-negative int value.

func Command_New(script *rex.Script, params []*rex.Value)
    Creates a new random number generator.

	rand:new seed

    Returns a handle to the new random number generator.


