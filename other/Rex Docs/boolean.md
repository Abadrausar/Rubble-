PACKAGE DOCUMENTATION

package boolean
    import "dctech/rex/commands/boolean"



FUNCTIONS

func Command_And(script *rex.Script, params []*rex.Value)
    Ands two values.

	bool:and a b

    Returns a && b

func Command_Not(script *rex.Script, params []*rex.Value)
    Inverts a value.

	bool:not a

    Returns !a

func Command_Or(script *rex.Script, params []*rex.Value)
    Ors two values.

	bool:or a b

    Returns a || b

func Command_SAnd(script *rex.Script, params []*rex.Value)
    Ands two values using short-circuit evaluation.

	bool:sand a b

    a and b must both be code blocks! Returns a && b

func Command_SOr(script *rex.Script, params []*rex.Value)
    Ors two values using short-circuit evaluation.

	bool:sor a b

    a and b must both be code blocks! Returns a || b


