/*
Copyright 2012-2013 by Milo Christiansen

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

// Raptor Base Commands.
package base

import "dctech/raptor"

// Setup adds the base commands to the state.
// The base commands are:
//	nop
//	ret
//	exit
//	break
//	breakloop
//	error
//	command
//	delcommand
//	getcommand
//	namespace
//	delnamespace
//	var
//	delvar
//	this
//	set
//	exists
//	len
//	run
//	eval
//	evalinparent
//	evalinnew
//	if
//	loop
//	foreach
// Also registers the following indexables:
//	map
//	array
// The base commands are more or less required.
func Setup(state *raptor.State) {
	state.NewNativeCommand("nop", CommandNop)
	state.NewNativeCommand("ret", CommandRet)
	state.NewNativeCommand("exit", CommandExit)
	state.NewNativeCommand("break", CommandBreak)
	state.NewNativeCommand("breakloop", CommandBreakLoop)
	state.NewNativeCommand("error", CommandError)
	state.NewNativeCommand("command", CommandCommand)
	state.NewNativeCommand("delcommand", CommandDelCommand)
	state.NewNativeCommand("getcommand", CommandGetCommand)
	state.NewNativeCommand("namespace", CommandNamespace)
	state.NewNativeCommand("delnamespace", CommandDelNamespace)
	state.NewNativeCommand("var", CommandVar)
	state.NewNativeCommand("delvar", CommandDelVar)
	state.NewNativeCommand("this", CommandThis)
	state.NewNativeCommand("set", CommandSet)
	state.NewNativeCommand("exists", CommandExists)
	state.NewNativeCommand("len", CommandLen)
	state.NewNativeCommand("run", CommandRun)
	state.NewNativeCommand("eval", CommandEval)
	state.NewNativeCommand("evalinparent", CommandEvalInParent)
	state.NewNativeCommand("evalinnew", CommandEvalInNew)
	state.NewNativeCommand("if", CommandIf)
	state.NewNativeCommand("loop", CommandLoop)
	state.NewNativeCommand("foreach", CommandForEach)

	state.RegisterType("map", NewScriptMapFromLit)
	state.RegisterType("array", NewScriptArrayFromLit)
}

// Does nothing.
// 	nop
// Returns unchanged.
func CommandNop(state *raptor.State, params []*raptor.Value) {
}

// Return from current command.
// 	ret [value]
// Some commands will be bypassed like if and loop, for example calling ret from inside a 
// loop will not return from the loop, it will return from the command that called loop.
// See break.
// Returns value or unchanged.
func CommandRet(state *raptor.State, params []*raptor.Value) {
	if len(params) != 0 && len(params) != 1 {
		panic("Wrong number of params to ret.")
	}

	state.Return = true
	if len(params) > 0 {
		state.RetVal = params[0]
	}
}

// Exit the script.
// 	exit [value]
// Returns value or unchanged.
func CommandExit(state *raptor.State, params []*raptor.Value) {
	if len(params) != 0 && len(params) != 1 {
		panic("Wrong number of params to exit.")
	}

	if len(params) > 0 {
		state.RetVal = params[0]
	}
	state.Exit = true
}

// A "soft" return, break will never return more than one level.
// 	break [value]
// Calling break inside a loop or if command will return from the current BLOCK not the command itself,
// this makes break good for ensuring if returns a specific value and/or "returning" a value to loop.
// Returns value or unchanged.
func CommandBreak(state *raptor.State, params []*raptor.Value) {
	if len(params) != 0 && len(params) != 1 {
		panic("Wrong number of params to break.")
	}

	state.Break = true
	if len(params) > 0 {
		state.RetVal = params[0]
	}
}

// Forces a return until it hits a loop or foreach command or the script exits.
// 	breakloop [value]
// Returns value or unchanged.
func CommandBreakLoop(state *raptor.State, params []*raptor.Value) {
	if len(params) != 0 && len(params) != 1 {
		panic("Wrong number of params to breakloop.")
	}

	state.BreakLoop = true
	if len(params) > 0 {
		state.RetVal = params[0]
	}
}

// Manipulates the error flag.
// 	error [value]
// If you pass no params the error flag will returned, to set or unset the flag pass a boolean value.
// Returns unchanged or the value of the error flag.
func CommandError(state *raptor.State, params []*raptor.Value) {
	if len(params) != 0 && len(params) != 1 {
		panic("Wrong number of params to error.")
	}

	if len(params) == 0 {
		state.RetVal = raptor.NewValueBool(state.Error)
		return
	}
	state.Error = params[0].Bool()
}

// Creates a new user command.
// 	command name code [paramName...]
// Returns unchanged.
func CommandCommand(state *raptor.State, params []*raptor.Value) {
	if len(params) < 2 {
		panic("Wrong number of params to command.")
	}
	if len(params) == 2 {
		// no params
		state.NewUserCommand(params[0].String(), params[1], make([]*raptor.Value, 0, 0))
	} else if len(params) == 3 && params[1].String() == "..." {
		// variable params
		state.NewUserCommand(params[0].String(), params[2], nil)
	} else {
		// fixed param count
		state.NewUserCommand(params[0].String(), params[len(params)-1], params[1:len(params)-1])
	}
}

// Deletes a command.
// Be careful with this one! Some actions are not reversable from a script.
// 	delcommand name
// Returns unchanged.
func CommandDelCommand(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to delcommand.")
	}

	state.DeleteCommand(params[0].String())
}

// Gets a reference to a command.
// Note that command references are just strings with a special type. This type is needed to make this work correctly.
// 	getcommand name
// Returns a reference to the command.
func CommandGetCommand(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to getcommand.")
	}

	state.RetVal = raptor.NewValueCommand(params[0].String())
}

// Creates a new namespace.
// 	namespace name 
// Returns unchanged.
func CommandNamespace(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to namespace.")
	}

	state.NewNameSpace(params[0].String())
}

// Deletes a namespace. 
// Be careful with this one! Some actions are not reversable from a script.
// For example it may be a very bad idea to delete debug or int.
// 	delnamespace name 
// Returns unchanged.
func CommandDelNamespace(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to delnamespace.")
	}

	state.DeleteNameSpace(params[0].String())
}

// Creates a new variable setting the value to value if present.
// 	var name [value]
// Returns value or unchanged.
func CommandVar(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic("Wrong number of params to var.")
	}

	if len(params) > 1 {
		state.NewVar(params[0].String(), params[1])
		state.RetVal = params[1]
		return
	}
	state.NewVar(params[0].String(), raptor.NewValueString(""))
}

// Deletes a variable.
// 	delvar name
// Returns the deleted vars value.
func CommandDelVar(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to delvar.")
	}
	state.RetVal = state.DeleteVar(params[0].String())
}

// Retrieves the current "this" value.
// The value returned by this command will only be useful under certain circumstaces!
// 	this
// Returns the value of the This register.
func CommandThis(state *raptor.State, params []*raptor.Value) {
	state.RetVal = state.This
}

// Sets the value of variable "name" to value or sets the value of the map or array at index to value.
// 	set name value
// 	set objectvalue index value
// Returns value.
func CommandSet(state *raptor.State, params []*raptor.Value) {
	if len(params) != 2 && len(params) != 3 {
		panic("Wrong number of params to set.")
	}

	if len(params) == 2 {
		state.SetValue(params[0].String(), params[1])
		state.RetVal = params[1]
		return
	}

	val := params[0].EditIndexable()
	if val == nil {
		panic("Non-EditIndexable object passed to set.")
	}

	if !val.Set(params[1].String(), params[2]) {
		panic("Attempt to write to readonly index with set.")
	}
	state.RetVal = params[2]
}

// Returns true (-1) if variable exists or if a index exists in a map or array.
// 	exists name
//	exists value index
// Returns true or false.
func CommandExists(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic("Wrong number of params to exists.")
	}

	if len(params) == 1 {
		if state.VarExists(params[0].String()) {
			state.RetVal = raptor.NewValueInt64(-1)
			return
		}
		state.RetVal = raptor.NewValueInt64(0)
		return
	}

	val := params[0].Indexable()
	if val == nil {
		panic("Non-Indexable object passed to exists.")
	}

	if val.Exists(params[1].String()) {
		state.RetVal = raptor.NewValueBool(true)
		return
	}
	state.RetVal = raptor.NewValueBool(false)
}

// Fetches the element count of an Indexable.
// 	len value
// Returns the element count.
func CommandLen(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to len.")
	}

	val := params[0].Indexable()
	if val == nil {
		panic("Non-Indexable object passed to exists.")
	}

	state.RetVal = raptor.NewValueInt64(val.Len())
}

// Runs code as a user command.
// 	run code [params...]
// Returns the return value of the last command in the code it runs.
func CommandRun(state *raptor.State, params []*raptor.Value) {
	if len(params) < 1 {
		panic("Wrong number of params to run.")
	}

	state.Envs.Add(raptor.NewEnvironment())

	state.AddParamsValue(params[1:]...)

	state.Code.AddCodeSource(params[0].CodeSource())
	state.Exec()
	state.Envs.Remove()
	state.Return = false
}

// Evaluates code in the current environment.
// 	eval code
// Returns the return value of the last command in the code it runs.
func CommandEval(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to eval.")
	}

	state.Code.AddCodeSource(params[0].CodeSource())
	state.Exec()
}

// Evaluates code in the current environment's parent.
// 	evalinparent code
// Returns the return value of the last command in the code it runs.
func CommandEvalInParent(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to evalinparent.")
	}
	if len(*state.Envs) <= 1 {
		panic("Call to evalinparent from code running in root env.")
	}

	state.Code.AddCodeSource(params[0].CodeSource())
	tempEnv := state.Envs.Remove()
	state.Exec()
	state.Envs.Add(tempEnv)
}

// Evaluates code in a new environment.
// 	evalinnew code
// Returns the return value of the last command in the code it runs.
func CommandEvalInNew(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to evalinnew.")
	}

	state.Code.AddCodeSource(params[0].CodeSource())
	state.Envs.Add(raptor.NewEnvironment())
	state.Exec()
	state.Envs.Remove()
}

// If the condition is true run true code else if false code exists call false code.
// 	if condition truecode [falsecode]
// Returns the return value of the last command in the code it runs or unchanged.
func CommandIf(state *raptor.State, params []*raptor.Value) {
	if len(params) != 2 && len(params) != 3 {
		panic("Wrong number of params to if.")
	}

	if params[0].Bool() {
		state.Code.AddCompiledScript(params[1].CompiledScript())
		state.Exec()
		return
	}

	if len(params) > 2 {
		state.Code.AddCompiledScript(params[2].CompiledScript())
		state.Exec()
		return
	}
}

// Runs code for as long as the code returns true.
// 	loop code
// Returns the return value of the last command in the code it runs, always false unless loop
// exited with ret (In which case the return value is unusable by the command calling loop anyway).
func CommandLoop(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to loop.")
	}

	code := params[0].CompiledScript()
	for {
		state.Code.AddCodeSource(raptor.NewCompiledLexer(code))
		state.Exec()
		state.BreakLoop = false
		if !state.RetVal.Bool() {
			return
		}
	}
	panic("CommandLoop: unreachable")
}

// Runs code as command for each entry in a map or array value.
// 	foreach objectvalue code
// Params for code:
//	code key value
// If code returns false foreach aborts.
// Does not stop returns, but does work with breakloop.
// Returns the return value of the last command in code.
func CommandForEach(state *raptor.State, params []*raptor.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to foreach.")
	}

	val := params[0].Indexable()
	if val == nil {
		panic("Non-Indexable object passed to foreach.")
	}

	code := params[1].CompiledScript()

	for _, i := range val.Keys() {
		state.Code.AddCodeSource(raptor.NewCompiledLexer(code))
		state.Envs.Add(raptor.NewEnvironment())
		state.AddParamsValue(raptor.NewValueString(i), val.Get(i))
		state.Exec()
		state.BreakLoop = false
		if !state.RetVal.Bool() {
			return
		}
		state.Envs.Remove()
	}
}
