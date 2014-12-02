PACKAGE DOCUMENTATION

package float
    import "dctech/rex/commands/float"



FUNCTIONS

func Command_Add(script *rex.Script, params []*rex.Value)
    Adds two values.

	float:add a b

    Returns a + b

func Command_Div(script *rex.Script, params []*rex.Value)
    Divides two values.

	float:div a b

    Returns a / b

func Command_Eq(script *rex.Script, params []*rex.Value)
    Floating point equal (basically useless, but included anyway).

	float:eq a b

    Returns true or false.

func Command_Gt(script *rex.Script, params []*rex.Value)
    Floating point greater than.

	float:gt a b

    Returns true or false.

func Command_Lt(script *rex.Script, params []*rex.Value)
    Floating point less than.

	float:lt a b

    Returns true or false.

func Command_Mul(script *rex.Script, params []*rex.Value)
    Multiplies two values.

	float:mul a b

    Returns a * b

func Command_Sub(script *rex.Script, params []*rex.Value)
    Subtracts two values.

	float:sub a b

    Returns a - b


