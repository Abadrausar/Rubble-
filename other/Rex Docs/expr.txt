PACKAGE DOCUMENTATION

package expr
    import "dctech/rex/commands/expr"



FUNCTIONS

func Command_Expr(script *rex.Script, params []*rex.Value)
    Evaluates an expression.

	expr expression values...

    Supported operators (in precedence order):

	( )
	/ * %
	+ -
	== != > < >= <=
	!
	&& ||

    Values are specified by any non-operator character and are filled in by
    position. Returns the value of the evaluated expression.


