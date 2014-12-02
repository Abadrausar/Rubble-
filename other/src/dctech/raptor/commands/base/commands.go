/*
Copyright 2012-2014 by Milo Christiansen

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

// Setup adds the base commands to the script.
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
//	copy
//	this
//	set
//	exists
//	len
//	run
//	eval
//	evalinparent
//	evalinnew
//	trap
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
	state.NewNativeCommand("copy", CommandCopy)
	state.NewNativeCommand("this", CommandThis)
	state.NewNativeCommand("set", CommandSet)
	state.NewNativeCommand("exists", CommandExists)
	state.NewNativeCommand("len", CommandLen)
	state.NewNativeCommand("run", CommandRun)
	state.NewNativeCommand("eval", CommandEval)
	state.NewNativeCommand("evalinparent", CommandEvalInParent)
	state.NewNativeCommand("evalinnew", CommandEvalInNew)
	state.NewNativeCommand("trap", CommandTrap)
	state.NewNativeCommand("if", CommandIf)
	state.NewNativeCommand("loop", CommandLoop)
	state.NewNativeCommand("foreach", CommandForEach)

	state.RegisterType("map", NewScriptMapFromLit)
	state.RegisterType("array", NewScriptArrayFromLit)
}

// Does nothing.
// 	nop
// Returns unchanged.
func CommandNop(script *raptor.Script, params []*raptor.Value) {
}

// Return from current command.
// 	ret [value]
// Some commands will be bypassed like if and loop, for example calling ret from inside a
// loop will not return from the loop, it will return from the command that called loop.
// See break.
// Returns value or unchanged.
func CommandRet(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 0 && len(params) != 1 {
		panic(script.BadParamCount("0 or 1"))
	}

	script.Return = true
	if len(params) > 0 {
		script.RetVal = params[0]
	}
}

// Exit the script.
// 	exit [value]
// Returns value or unchanged.
func CommandExit(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 0 && len(params) != 1 {
		panic(script.BadParamCount("0 or 1"))
	}

	script.Exit = true
	if len(params) > 0 {
		script.RetVal = params[0]
	}
}

// A "soft" return, break will never return more than one level.
// 	break [value]
// Calling break inside a loop or if command will return from the current BLOCK not the command itself,
// this makes break good for ensuring if returns a specific value and/or "returning" a value to loop.
// Returns value or unchanged.
func CommandBreak(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 0 && len(params) != 1 {
		panic(script.BadParamCount("0 or 1"))
	}

	script.Break = true
	if len(params) > 0 {
		script.RetVal = params[0]
	}
}

// Forces a return until it hits a loop or foreach command or the script exits.
// 	breakloop [value]
// Returns value or unchanged.
func CommandBreakLoop(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 0 && len(params) != 1 {
		panic(script.BadParamCount("0 or 1"))
	}

	script.BreakLoop = true
	if len(params) > 0 {
		script.RetVal = params[0]
	}
}

// Manipulates the error flag.
// 	error [value]
// If you pass no params the error flag will returned, to set or unset the flag pass a boolean value.
// Returns unchanged or the value of the error flag.
func CommandError(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 0 && len(params) != 1 {
		panic(script.BadParamCount("0 or 1"))
	}

	if len(params) == 0 {
		script.RetVal = raptor.NewValueBool(script.Error)
		return
	}
	script.Error = params[0].Bool()
}

// Creates a new user command.
// 	command name code [paramName...]
// Returns unchanged.
func CommandCommand(script *raptor.Script, params []*raptor.Value) {
	if len(params) < 2 {
		panic(script.BadParamCount(">=2"))
	}
	if len(params) == 2 {
		// no params
		script.NewUserCommand(params[0].String(), params[1], make([]*raptor.Value, 0, 0))
	} else if len(params) == 3 && params[1].String() == "..." {
		// variable params
		script.NewUserCommand(params[0].String(), params[2], nil)
	} else {
		// fixed param count
		script.NewUserCommand(params[0].String(), params[len(params)-1], params[1:len(params)-1])
	}
}

// Deletes a command.
// Be careful with this one! Some actions are not reversable from a script.
// 	delcommand name
// Returns unchanged.
func CommandDelCommand(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	script.DeleteCommand(params[0].String())
}

// Gets a reference to a command.
// Note that command references are just strings with a special type. This type is needed to make this work correctly.
// 	getcommand name
// Returns a reference to the command.
func CommandGetCommand(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	script.RetVal = raptor.NewValueCommand(params[0].String())
}

// Creates a new namespace.
// 	namespace name
// Returns unchanged.
func CommandNamespace(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	script.NewNameSpace(params[0].String())
}

// Deletes a namespace.
// Be careful with this one! Some actions are not reversable from a script.
// For example it may be a very bad idea to delete debug or int.
// 	delnamespace name
// Returns unchanged.
func CommandDelNamespace(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	script.DeleteNameSpace(params[0].String())
}

// Creates a new variable setting the value to value if present.
// 	var name [value]
// Returns value or unchanged.
func CommandVar(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic(script.BadParamCount("1 or 2"))
	}

	if len(params) > 1 {
		script.NewVar(params[0].String(), params[1])
		script.RetVal = params[1]
		return
	}
	script.NewVar(params[0].String(), raptor.NewValue())
}

// Deletes a variable.
// 	delvar name
// Returns the deleted vars value.
func CommandDelVar(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}
	script.RetVal = script.DeleteVar(params[0].String())
}

// Copies a value.
// 	copy value
// Note that this command may not be useful for types object and code.
// The new value has invalid position info.
// Returns the new value.
func CommandCopy(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	val := new(raptor.Value)
	val.Type = params[0].Type
	val.Pos = raptor.NewPosition(0, -1, "")

	switch val.Type {
	case raptor.TypString:
		val.Data = params[0].Data.(string) // Not sure if this is needed, but it should work at least.

	case raptor.TypInt:
		val.Data = params[0].Data.(int64)

	case raptor.TypFloat:
		val.Data = params[0].Data.(float64)

	case raptor.TypBool:
		val.Data = params[0].Data.(bool)

	case raptor.TypCode:
		val.Data = params[0].Data

	case raptor.TypCommand:
		val.Data = params[0].Data.(string)

	case raptor.TypObject:
		val.Data = params[0].Data

	case raptor.TypNil:
		val.Data = nil
	}

	script.RetVal = val
}

// Retrieves the current "this" value.
// The value returned by this command will only be useful under certain circumstaces!
// 	this
// Returns the value of the This register.
func CommandThis(script *raptor.Script, params []*raptor.Value) {
	script.RetVal = script.This
}

// Sets the value of variable "name" to value or sets the value of the map or array at index to value.
// 	set name value
// 	set objectvalue index value
// Returns value.
func CommandSet(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 && len(params) != 3 {
		panic(script.BadParamCount("2 or 3"))
	}

	if len(params) == 2 {
		script.SetValue(params[0].String(), params[1])
		script.RetVal = params[1]
		return
	}

	val := params[0].EditIndexable()
	if val == nil {
		panic(script.GeneralCmdError("Attempt to use non-EditIndexable object with indexing form."))
	}

	if !val.Set(params[1].String(), params[2]) {
		panic(script.GeneralCmdError("Attempt to write to readonly index."))
	}
	script.RetVal = params[2]
}

// Returns true if variable exists or if a index exists in a map or array.
// 	exists name
//	exists value index
// Returns true or false.
func CommandExists(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic(script.BadParamCount("1 or 2"))
	}

	if len(params) == 1 {
		if script.VarExists(params[0].String()) {
			script.RetVal = raptor.NewValueInt64(-1)
			return
		}
		script.RetVal = raptor.NewValueInt64(0)
		return
	}

	val := params[0].Indexable()
	if val == nil {
		panic(script.GeneralCmdError("Attempt to use non-Indexable object with indexing form."))
	}

	if val.Exists(params[1].String()) {
		script.RetVal = raptor.NewValueBool(true)
		return
	}
	script.RetVal = raptor.NewValueBool(false)
}

// Fetches the element count of an Indexable.
// 	len value
// Returns the element count.
func CommandLen(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	val := params[0].Indexable()
	if val == nil {
		panic(script.GeneralCmdError("Attempt to fetch non-Indexable object's size."))
	}

	script.RetVal = raptor.NewValueInt64(val.Len())
}

// Runs code as a user command.
// 	run code [params...]
// Returns the return value of the last command in the code it runs.
func CommandRun(script *raptor.Script, params []*raptor.Value) {
	if len(params) < 1 {
		panic(script.BadParamCount(">=1"))
	}

	script.Envs.Add(raptor.NewEnvironment())

	script.AddParamsValue(params[1:]...)

	script.Code.Add(params[0].CodeSource())
	script.Exec()
	script.Envs.Remove()
	script.Return = false
}

// Evaluates code in the current environment.
// 	eval code
// Returns the return value of the last command in the code it runs.
func CommandEval(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	script.Code.Add(params[0].CodeSource())
	script.Exec()
}

// Evaluates code in the current environment's parent.
// 	evalinparent code
// Returns the return value of the last command in the code it runs.
func CommandEvalInParent(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}
	if len(*script.Envs) <= 1 {
		panic(script.GeneralCmdError("The code is running in root environment."))
	}

	script.Code.Add(params[0].CodeSource())
	tempEnv := script.Envs.Remove()
	script.Exec()
	script.Envs.Add(tempEnv)
}

// Evaluates code in a new environment.
// 	evalinnew code
// Returns the return value of the last command in the code it runs.
func CommandEvalInNew(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	script.Code.Add(params[0].CodeSource())
	script.Envs.Add(raptor.NewEnvironment())
	script.Exec()
	script.Envs.Remove()
}

// Evaluates code in the current environment and halts any unrecoverable errors.
// WARNING!! This command may result in a messed up state/script if you are not careful!
// There is a reason "unrecoverable errors" are called that! In most cases you should be OK, but be careful.
// 	trap code
// Returns the return value of the last command in the code it runs or sets the error
// flag and returns a string describing the error..
func CommandTrap(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	script.Code.Add(params[0].CodeSource())
	err := script.SafeExec()
	if err != nil {
		script.Error = true
		script.RetVal = raptor.NewValueString(err.Error())
	}
}

// If the condition is true run true code else if false code exists call false code.
// 	if condition truecode [falsecode]
// Returns the return value of the last command in the code it runs or unchanged.
func CommandIf(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 && len(params) != 3 {
		panic(script.BadParamCount("2 or 3"))
	}

	if params[0].Bool() {
		script.Code.AddCode(params[1].Code())
		script.Exec()
		return
	}

	if len(params) > 2 {
		script.Code.AddCode(params[2].Code())
		script.Exec()
		return
	}
}

// Runs code for as long as the code returns true.
// 	loop code
// Returns the return value of the last command in the code it runs, always false unless loop
// exited with ret (In which case the return value is unusable by the command calling loop anyway).
func CommandLoop(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	code := params[0].Code()
	for {
		script.Code.Add(raptor.NewCodeReader(code))
		script.Exec()
		script.BreakLoop = false
		if !script.RetVal.Bool() {
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
func CommandForEach(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic(script.BadParamCount("2"))
	}

	val := params[0].Indexable()
	if val == nil {
		panic(script.GeneralCmdError("Attempt to range non-Indexable object."))
	}

	code := params[1].Code()

	for _, i := range val.Keys() {
		script.Code.Add(raptor.NewCodeReader(code))
		script.Envs.Add(raptor.NewEnvironment())
		script.AddParamsValue(raptor.NewValueString(i), val.Get(i))
		script.Exec()
		script.BreakLoop = false
		if !script.RetVal.Bool() {
			return
		}
		script.Envs.Remove()
	}
}
