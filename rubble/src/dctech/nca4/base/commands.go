// NCA v4 Base Commands.
package base

import "dctech/nca4"

// Setup adds the base commands to the state.
// The base commands are:
//	nop
//	ret
//	exit
//	break
//	seterror
//	geterror
//	command
//	namespace
//	var
//	delvar
//	map
//	delmap
//	set
//	run
//	eval
//	evalinparent
//	if
//	loop
//	foreach
// The base commands are more or less required.
func Setup(state *nca4.State) {
	state.NewNativeCommand("nop", CommandNop)
	state.NewNativeCommand("ret", CommandRet)
	state.NewNativeCommand("exit", CommandExit)
	state.NewNativeCommand("break", CommandBreak)
	state.NewNativeCommand("seterror", CommandSetError)
	state.NewNativeCommand("geterror", CommandGetError)
	state.NewNativeCommand("command", CommandCommand)
	state.NewNativeCommand("namespace", CommandNamespace)
	state.NewNativeCommand("var", CommandVar)
	state.NewNativeCommand("delvar", CommandDelVar)
	state.NewNativeCommand("map", CommandMap)
	state.NewNativeCommand("delmap", CommandDelMap)
	state.NewNativeCommand("set", CommandSet)
	state.NewNativeCommand("run", CommandRun)
	state.NewNativeCommand("eval", CommandEval)
	state.NewNativeCommand("evalinparent", CommandEvalInParent)
	state.NewNativeCommand("if", CommandIf)
	state.NewNativeCommand("loop", CommandLoop)
	state.NewNativeCommand("foreach", CommandForEach)
}

// Does nothing.
// 	nop
// VERY useful as a replacement for short lived variables.
// Example:
//	(fileio:read "somefile.ini")
//	(ini:parse (nop))
//	(ini:getvalue (nop) "somesection" "somekey")
//	(console:print (nop))
// This showcases the rule that if a command returns nothing the return 
// value will be the same as the return value of the previous command to run.
// Returns unchanged.
func CommandNop(state *nca4.State, params []*nca4.Value) {
}

// Return from current command.
// 	ret [value]
// Some commands will be bypassed like if and loop, for example calling ret from inside a 
// loop will not return from the loop, it will return from the command that called loop.
// See break.
// Returns value or unchanged.
func CommandRet(state *nca4.State, params []*nca4.Value) {
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
func CommandExit(state *nca4.State, params []*nca4.Value) {
	if len(params) != 0 && len(params) != 1 {
		panic("Wrong number of params to exit.")
	}

	if len(params) > 0 {
		panic(nca4.ExitScript(params[0].String()))
	}
	panic(nca4.ExitScript(state.RetVal.String()))
}

// A "soft" return, break will never return more than one level.
// 	break [value]
// Calling break inside a loop or if command will return from the current BLOCK not the command itself,
// this makes break good for ensuring if returns a specific value and/or "returning" a value to loop.
// Returns value or unchanged.
func CommandBreak(state *nca4.State, params []*nca4.Value) {
	if len(params) != 0 && len(params) != 1 {
		panic("Wrong number of params to break.")
	}

	state.Break = true
	if len(params) > 0 {
		state.RetVal = params[0]
	}
}

// Sets (or unsets) the error flag.
// 	seterror [value]
// If you pass no params the error flag will be set, to unset pass "0" or ""
// Returns unchanged.
func CommandSetError(state *nca4.State, params []*nca4.Value) {
	if len(params) != 0 && len(params) != 1 {
		panic("Wrong number of params to seterror.")
	}

	state.Error = true
	if len(params) > 0 {
		state.Error = params[0].Bool()
	}
}

// Gets the error flag.
// 	geterror
// Returns "-1" or "0".
func CommandGetError(state *nca4.State, params []*nca4.Value) {
	if len(params) != 0 {
		panic("Wrong number of params to geterror.")
	}

	if state.Error {
		state.RetVal = nca4.NewValue("-1")
		return
	}
	state.RetVal = nca4.NewValue("0")
}

// Creates a new user command.
// 	command name code [paramName...]
// Returns unchanged.
func CommandCommand(state *nca4.State, params []*nca4.Value) {
	if len(params) < 2 {
		panic("Wrong number of params to command.")
	}
	if len(params) == 2 {
		// no params
		state.NewUserCommand(params[0].String(), params[1], make([]*nca4.Value, 0, 0))
	} else if len(params) == 3 && params[1].String() == "..." {
		// variable params
		state.NewUserCommand(params[0].String(), params[2], nil)
	} else {
		// fixed param count
		state.NewUserCommand(params[0].String(), params[len(params)-1], params[1:len(params)-1])
	}
}

// Creates a new namespace.
// 	namespace name 
// Returns unchanged.
func CommandNamespace(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to namespace.")
	}

	state.NewNameSpace(params[0].String())
}

// Creates a new variable setting the value to value if present.
// 	var name [value]
// Returns value or unchanged.
func CommandVar(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic("Wrong number of params to var.")
	}

	if len(params) > 1 {
		state.NewVar(params[0].String(), params[1])
		state.RetVal = params[1]
		return
	}
	state.NewVar(params[0].String(), new(nca4.Value))
}

// Deletes a variable.
// 	delvar name
// Returns the deleted vars value.
func CommandDelVar(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to delvar.")
	}
	state.RetVal = state.DeleteVar(params[0].String())
}

// Creates a new map.
// 	map name
// Returns unchanged
func CommandMap(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to array.")
	}

	state.NewMap(params[0].String())
}

// Deletes a map.
// 	delmap name
// Returns unchanged
func CommandDelMap(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to delarray.")
	}

	state.DeleteMap(params[0].String())
}

// Sets the value of variable "name" to value or sets the value of map "name" at index to value.
// 	set name [index] value
// Returns value.
func CommandSet(state *nca4.State, params []*nca4.Value) {
	if len(params) != 2 && len(params) != 3 {
		panic("Wrong number of params to set.")
	}

	if len(params) > 2 {
		state.SetValueMap(params[0].String(), params[1].String(), params[2])
		state.RetVal = params[2]
		return
	}
	state.SetValue(params[0].String(), params[1])
	state.RetVal = params[1]
}

// Runs code as a user command.
// 	run code [params...]
// Returns the return value of the last command in the code it runs.
func CommandRun(state *nca4.State, params []*nca4.Value) {
	if len(params) < 1 {
		panic("Wrong number of params to run.")
	}

	state.Envs.Add(nca4.NewEnvironment())

	state.AddParamsValue(params[1:]...)

	state.Code.AddLexer(params[0].AsLexer())
	state.Exec()
	state.Envs.Remove()
	state.Return = false
}

// Evaluates code in the current environment.
// 	eval code
// Returns the return value of the last command in the code it runs.
func CommandEval(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to eval.")
	}

	state.Code.AddLexer(params[0].AsLexer())
	state.Exec()
}

// Evaluates code in the current environment's parent.
// 	evalinparent code
// Returns the return value of the last command in the code it runs.
func CommandEvalInParent(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to evalinparent.")
	}
	if len(*state.Envs) <= 1 {
		panic("Call to evalinparent from code running in root env.")
	}

	state.Code.AddLexer(params[0].AsLexer())
	tempEnv := state.Envs.Remove()
	state.Exec()
	state.Envs.Add(tempEnv)
}

// If the condition is true run true code else if false code exists call false code.
// 	if condition truecode [falsecode]
// Returns the return value of the last command in the code it runs or unchanged.
func CommandIf(state *nca4.State, params []*nca4.Value) {
	if len(params) != 2 && len(params) != 3 {
		panic("Wrong number of params to if.")
	}

	runtrue := false
	if params[0].Bool() {
		runtrue = true
	}

	if runtrue {
		state.Code.AddLexer(params[1].AsLexer())
		state.Exec()
		return
	}

	if len(params) > 2 {
		state.Code.AddLexer(params[2].AsLexer())
		state.Exec()
		return
	}
}

// Runs code for as long as the code returns true.
// 	loop code
// Returns the return value of the last command in the code it runs, always "0" or "" unless loop
// exited with ret (In which case the return value is unusable by the command calling loop anyway).
func CommandLoop(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to loop.")
	}

	for {
		state.Code.AddLexer(params[0].AsLexer())
		state.Exec()
		if !state.RetVal.Bool() {
			return
		}
	}
	panic("CommandLoop: unreachable")
}

// Runs code as command for each entry in a map.
// 	foreach map code
// Params for code:
//	code key value
// If code returns false foreach aborts.
// Returns the return value of the last command in code.
func CommandForEach(state *nca4.State, params []*nca4.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to foreach.")
	}

	for i, val := range state.GetMap(params[0].String()) {
		state.Code.AddLexer(params[1].AsLexer())
		state.Envs.Add(nca4.NewEnvironment())
		state.AddParams(i, val.String())
		state.Exec()
		if !state.RetVal.Bool() {
			return
		}
		state.Envs.Remove()
	}
	panic("CommandLoop: unreachable")
}
