/*
Copyright 2014 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

// Rex Boolean Commands.
package boolean

import "dctech/rex"

// Adds the boolean commands to the state.
// The boolean commands are:
//	bool:and
//	bool:or
//	bool:sand
//	bool:sor
//	bool:not
func Setup(state *rex.State) (err error) {
	defer func() {
		if !state.NoRecover {
			if x := recover(); x != nil {
				if y, ok := x.(rex.ScriptError); ok {
					err = y
					return
				}
				panic(x)
			}
		}
	}()
	
	mod := state.RegisterModule("bool")
	mod.RegisterCommand("and", Command_And)
	mod.RegisterCommand("or", Command_Or)
	mod.RegisterCommand("sand", Command_SAnd)
	mod.RegisterCommand("sor", Command_SOr)
	mod.RegisterCommand("not", Command_Not)
	
	return nil
}

// Ands two values.
// 	bool:and a b
// Returns a && b
func Command_And(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("bool:and", "2")
	}

	script.RetVal = rex.NewValueBool(params[0].Bool() && params[1].Bool())
	return
}

// Ors two values.
// 	bool:or a b
// Returns a || b
func Command_Or(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("bool:or", "2")
	}

	script.RetVal = rex.NewValueBool(params[0].Bool() || params[1].Bool())
	return
}

// Ands two values using short-circuit evaluation.
// 	bool:sand a b
// a and b must both be code blocks!
// Returns a && b
func Command_SAnd(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("bool:sand", "2")
	}
	
	if params[0].Type != rex.TypCode {
		rex.ErrorGeneralCmd("bool:sand", "Param 0 is not type code.")
	}
	
	if params[1].Type != rex.TypCode {
		rex.ErrorGeneralCmd("bool:sand", "Param 1 is not type code.")
	}
	
	block := params[0].Data.(*rex.Code)
	script.Locals.Add(script.Host, block)
	script.Exec(block)
	script.Return = false
	script.Locals.Remove()
	if !script.RetVal.Bool() {
		return
	}
	
	block = params[1].Data.(*rex.Code)
	script.Locals.Add(script.Host, block)
	script.Exec(block)
	script.Return = false
	script.Locals.Remove()
}

// Ors two values using short-circuit evaluation.
// 	bool:sor a b
// a and b must both be code blocks!
// Returns a || b
func Command_SOr(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("bool:sor", "2")
	}
	
	if params[0].Type != rex.TypCode {
		rex.ErrorGeneralCmd("bool:sor", "Param 0 is not type code.")
	}
	
	if params[1].Type != rex.TypCode {
		rex.ErrorGeneralCmd("bool:sor", "Param 1 is not type code.")
	}
	
	block := params[0].Data.(*rex.Code)
	script.Locals.Add(script.Host, block)
	script.Exec(block)
	script.Return = false
	script.Locals.Remove()
	if script.RetVal.Bool() {
		return
	}
	
	block = params[1].Data.(*rex.Code)
	script.Locals.Add(script.Host, block)
	script.Exec(block)
	script.Return = false
	script.Locals.Remove()
}

// Inverts a value.
// 	bool:not a
// Returns !a
func Command_Not(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("bool:not", "1")
	}

	script.RetVal = rex.NewValueBool(!params[0].Bool())
	return
}
