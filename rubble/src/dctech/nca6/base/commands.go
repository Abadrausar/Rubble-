// NCA v6 Base Commands.
package base

//import "fmt"
import "dctech/nca6"

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
//	map
//	array
//	delvar
//	set
//	exists
//	len
//	string
//	int
//	run
//	eval
//	evalinparent
//	evalinnew
//	if
//	loop
//	foreach
// The base commands are more or less required.
func Setup(state *nca6.State) {
	state.NewNativeCommand("nop", CommandNop)
	state.NewNativeCommand("ret", CommandRet)
	state.NewNativeCommand("exit", CommandExit)
	state.NewNativeCommand("break", CommandBreak)
	state.NewNativeCommand("seterror", CommandSetError)
	state.NewNativeCommand("geterror", CommandGetError)
	state.NewNativeCommand("command", CommandCommand)
	state.NewNativeCommand("namespace", CommandNamespace)
	state.NewNativeCommand("var", CommandVar)
	state.NewNativeCommand("map", CommandMap)
	state.NewNativeCommand("array", CommandArray)
	state.NewNativeCommand("delvar", CommandDelVar)
	state.NewNativeCommand("set", CommandSet)
	state.NewNativeCommand("exists", CommandExists)
	state.NewNativeCommand("len", CommandLen)
	state.NewNativeCommand("string", CommandString)
	state.NewNativeCommand("int", CommandInt)
	state.NewNativeCommand("run", CommandRun)
	state.NewNativeCommand("eval", CommandEval)
	state.NewNativeCommand("evalinparent", CommandEvalInParent)
	state.NewNativeCommand("evalinnew", CommandEvalInNew)
	state.NewNativeCommand("if", CommandIf)
	state.NewNativeCommand("loop", CommandLoop)
	state.NewNativeCommand("foreach", CommandForEach)
}

// Does nothing.
// 	nop
// Returns unchanged.
func CommandNop(state *nca6.State, params []*nca6.Value) {
}

// Return from current command.
// 	ret [value]
// Some commands will be bypassed like if and loop, for example calling ret from inside a 
// loop will not return from the loop, it will return from the command that called loop.
// See break.
// Returns value or unchanged.
func CommandRet(state *nca6.State, params []*nca6.Value) {
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
func CommandExit(state *nca6.State, params []*nca6.Value) {
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
func CommandBreak(state *nca6.State, params []*nca6.Value) {
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
// If you pass no params the error flag will be set, to unset pass a false value
// Returns unchanged.
func CommandSetError(state *nca6.State, params []*nca6.Value) {
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
// Returns -1 or 0.
func CommandGetError(state *nca6.State, params []*nca6.Value) {
	if len(params) != 0 {
		panic("Wrong number of params to geterror.")
	}

	if state.Error {
		state.RetVal = nca6.NewValueInt64(-1)
		return
	}
	state.RetVal = nca6.NewValueInt64(0)
}

// Creates a new user command.
// 	command name code [paramName...]
// Returns unchanged.
func CommandCommand(state *nca6.State, params []*nca6.Value) {
	if len(params) < 2 {
		panic("Wrong number of params to command.")
	}
	if len(params) == 2 {
		// no params
		state.NewUserCommand(params[0].String(), params[1], make([]*nca6.Value, 0, 0))
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
func CommandNamespace(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to namespace.")
	}

	state.NewNameSpace(params[0].String())
}

// Creates a new variable setting the value to value if present.
// 	var name [value]
// Returns value or unchanged.
func CommandVar(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic("Wrong number of params to var.")
	}

	if len(params) > 1 {
		state.NewVar(params[0].String(), params[1])
		state.RetVal = params[1]
		return
	}
	state.NewVar(params[0].String(), nca6.NewValueString(""))
}

// Creates a new map.
// 	map
// Returns the new map
func CommandMap(state *nca6.State, params []*nca6.Value) {
	state.RetVal = nca6.NewValueObject(NewScriptMap())
}

// Creates a new array.
// 	array [size]
// Returns the new array
func CommandArray(state *nca6.State, params []*nca6.Value) {
	if len(params) == 0 {
		state.RetVal = nca6.NewValueObject(NewScriptArray())
		return
	}
	state.RetVal = nca6.NewValueObject(NewScriptArraySized(int(params[0].Int64())))
}

// Deletes a variable.
// 	delvar name
// Returns the deleted vars value.
func CommandDelVar(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to delvar.")
	}
	state.RetVal = state.DeleteVar(params[0].String())
}

// Sets the value of variable "name" to value or sets the value of the map or array at index to value.
// 	set name value
// 	set objectvalue index value
// Returns value.
func CommandSet(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 && len(params) != 3 {
		panic("Wrong number of params to set.")
	}

	if len(params) == 2 {
		state.SetValue(params[0].String(), params[1])
		state.RetVal = params[1]
		return
	}
	
	if params[0].Type != nca6.TypObject {
		panic("Non-Object Param 0 passed to set.")
	}
	if !params[0].IsIndexable() {
		panic("Non-Indexable object Param 0 passed to set.")
	}
	val := params[0].Data.(nca6.Indexable)
	
	val.Set(params[1].String(), params[2])
	state.RetVal = params[2]
}

// Returns true (-1) if variable exists or if a index exists in a map or array.
// 	exists name
//	exists value index
// Returns 0 or -1.
func CommandExists(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic("Wrong number of params to exists.")
	}
	
	if len(params) == 1 {
		if state.VarExists(params[0].String()) {
			state.RetVal = nca6.NewValueInt64(-1)
			return
		}
		state.RetVal = nca6.NewValueInt64(0)
		return
	}
	
	if params[0].Type != nca6.TypObject {
		panic("Non-Object Param 0 passed to exists.")
	}
	if !params[0].IsIndexable() {
		panic("Non-Indexable object Param 0 passed to exists.")
	}
	val := params[0].Data.(nca6.Indexable)
	
	if val.Exists(params[1].String()) {
		state.RetVal = nca6.NewValueInt64(-1)
		return
	}
	state.RetVal = nca6.NewValueInt64(0)
}

// Fetches the element count of an Indexable.
// 	len value
// Returns the element count.
func CommandLen(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to len.")
	}
	
	if params[0].Type != nca6.TypObject {
		panic("Non-Object Param 0 passed to len.")
	}
	if !params[0].IsIndexable() {
		panic("Non-Indexable object Param 0 passed to len.")
	}
	val := params[0].Data.(nca6.Indexable)
	
	state.RetVal = nca6.NewValueInt64(val.Len())
}

// Forces a value's internal representation to be a string.
// 	string value
// Returns 0 or -1.
func CommandString(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to string.")
	}
	nca6.NewValueString(params[0].String())
}

// Forces a value's internal representation to be an integer.
// 	int value
// Returns the new value.
func CommandInt(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to int.")
	}
	nca6.NewValueInt64(params[0].Int64())
}

// Runs code as a user command.
// 	run code [params...]
// Returns the return value of the last command in the code it runs.
func CommandRun(state *nca6.State, params []*nca6.Value) {
	if len(params) < 1 {
		panic("Wrong number of params to run.")
	}

	state.Envs.Add(nca6.NewEnvironment())

	state.AddParamsValue(params[1:]...)

	state.Code.AddCodeSource(params[0].AsLexer())
	state.Exec()
	state.Envs.Remove()
	state.Return = false
}

// Evaluates code in the current environment.
// 	eval code
// Returns the return value of the last command in the code it runs.
func CommandEval(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to eval.")
	}

	state.Code.AddCodeSource(params[0].AsLexer())
	state.Exec()
}

// Evaluates code in the current environment's parent.
// 	evalinparent code
// Returns the return value of the last command in the code it runs.
func CommandEvalInParent(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to evalinparent.")
	}
	if len(*state.Envs) <= 1 {
		panic("Call to evalinparent from code running in root env.")
	}

	state.Code.AddCodeSource(params[0].AsLexer())
	tempEnv := state.Envs.Remove()
	state.Exec()
	state.Envs.Add(tempEnv)
}

// Evaluates code in a new environment.
// 	evalinnew code
// Returns the return value of the last command in the code it runs.
func CommandEvalInNew(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to evalinnew.")
	}

	state.Code.AddCodeSource(params[0].AsLexer())
	state.Envs.Add(nca6.NewEnvironment())
	state.Exec()
	state.Envs.Remove()
}

// If the condition is true run true code else if false code exists call false code.
// 	if condition truecode [falsecode]
// Returns the return value of the last command in the code it runs or unchanged.
func CommandIf(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 && len(params) != 3 {
		panic("Wrong number of params to if.")
	}

	runtrue := false
	if params[0].Bool() {
		runtrue = true
	}

	if runtrue {
		state.Code.AddCodeSource(params[1].AsLexer())
		state.Exec()
		return
	}

	if len(params) > 2 {
		state.Code.AddCodeSource(params[2].AsLexer())
		state.Exec()
		return
	}
}

// Runs code for as long as the code returns true.
// 	loop code
// Returns the return value of the last command in the code it runs, always "0" or "" unless loop
// exited with ret (In which case the return value is unusable by the command calling loop anyway).
func CommandLoop(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to loop.")
	}

	code := nca6.Compile(params[0].String(), params[0].Line)
	for {
		state.Code.AddCodeSource(nca6.NewCompiledLexer(code))
		state.Exec()
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
// Returns the return value of the last command in code.
func CommandForEach(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to foreach.")
	}

	if params[0].Type != nca6.TypObject {
		panic("Non-Object Param 0 passed to foreach.")
	}
	if !params[0].IsIndexable() {
		panic("Non-Indexable object Param 0 passed to foreach.")
	}
	val := params[0].Data.(nca6.Indexable)
	
	code := nca6.Compile(params[1].String(), params[1].Line)
	
	for _, i := range val.Keys() {
		state.Code.AddCodeSource(nca6.NewCompiledLexer(code))
		state.Envs.Add(nca6.NewEnvironment())
		state.AddParamsValue(nca6.NewValueString(i), val.Get(i))
		state.Exec()
		if !state.RetVal.Bool() {
			return
		}
		state.Envs.Remove()
	}
}
