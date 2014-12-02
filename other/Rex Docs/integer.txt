PACKAGE DOCUMENTATION

package integer
    import "dctech/rex/commands/integer"



FUNCTIONS

func Command_Add(script *rex.Script, params []*rex.Value)
    Adds two values.

	int:add a b

    Returns a + b

func Command_Dec(script *rex.Script, params []*rex.Value)
    Decrements a value "in place".

	int:-- a

    Returns unchanged.

func Command_Div(script *rex.Script, params []*rex.Value)
    Divides two values.

	int:div a b

    Returns a / b

func Command_Eq(script *rex.Script, params []*rex.Value)
    Integer equal.

	int:eq a b

    Returns true or false.

func Command_Gt(script *rex.Script, params []*rex.Value)
    Integer greater than.

	int:gt a b

    Returns true or false.

func Command_Inc(script *rex.Script, params []*rex.Value)
    Increments a value "in place".

	int:++ a

    Returns unchanged.

func Command_Lt(script *rex.Script, params []*rex.Value)
    Integer less than.

	int:lt a b

    Returns true or false.

func Command_Mod(script *rex.Script, params []*rex.Value)
    Gives the remainder of dividing two values.

	int:mod a b

    Returns a % b

func Command_Mul(script *rex.Script, params []*rex.Value)
    Multiplies two values.

	int:mul a b

    Returns a * b

func Command_Sub(script *rex.Script, params []*rex.Value)
    Subtracts two values.

	int:sub a b

    Returns a - b


