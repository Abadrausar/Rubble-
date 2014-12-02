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

// Rex Base Commands.
package base

import "dctech/rex"

// Setup adds the base commands to the script.
// The base commands are:
//	nop
//	modval
//	ret
//	exit
//	break
//	breakif
//	breakloop
//	breakloopif
//	eval
//	error
//	onerror
//	copy
//	exists
//	len
//	isnil
//	type
//	if
//	loop
//	for
//	foreach
// The base commands are more or less required.
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
	
	state.RegisterCommand("nop", Command_Nop)
	state.RegisterCommand("modval", Command_ModVal)
	state.RegisterCommand("ret", Command_Ret)
	state.RegisterCommand("exit", Command_Exit)
	state.RegisterCommand("break", Command_Break)
	state.RegisterCommand("breakif", Command_BreakIf)
	state.RegisterCommand("breakloop", Command_BreakLoop)
	state.RegisterCommand("breakloopif", Command_BreakLoopIf)
	state.RegisterCommand("eval", Command_Eval)
	state.RegisterCommand("error", Command_Error)
	state.RegisterCommand("onerror", Command_OnError)
	state.RegisterCommand("copy", Command_Copy)
	state.RegisterCommand("exists", Command_Exists)
	state.RegisterCommand("len", Command_Len)
	state.RegisterCommand("isnil", Command_IsNil)
	state.RegisterCommand("type", Command_Type)
	state.RegisterCommand("if", Command_If)
	state.RegisterCommand("loop", Command_Loop)
	state.RegisterCommand("for", Command_For)
	state.RegisterCommand("foreach", Command_ForEach)
	
	return nil
}

// Does nothing.
// 	nop
// Returns unchanged.
func Command_Nop(script *rex.Script, params []*rex.Value) {
}

// Copies the second value over the first, allows you to do some weird pointer-like things.
// In general do not use unless you know what you are doing.
// 	modval a b
// Returns unchanged.
func Command_ModVal(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("modval", "2")
	}
	
	params[0].Type = params[1].Type
	params[0].Data = params[1].Data
	params[0].Pos = params[1].Pos
}

// Return from current command.
// 	ret [value]
// Some commands will be bypassed like if and loop, for example calling ret from inside a
// loop will not return from the loop, it will return from the command that called loop.
// See break.
// Returns value or unchanged.
func Command_Ret(script *rex.Script, params []*rex.Value) {
	if len(params) != 0 && len(params) != 1 {
		rex.ErrorParamCount("ret", "0 or 1")
	}

	script.Return = true
	if len(params) > 0 {
		script.RetVal = params[0]
	}
}

// Exit the script.
// 	exit [value]
// Returns value or unchanged.
func Command_Exit(script *rex.Script, params []*rex.Value) {
	if len(params) != 0 && len(params) != 1 {
		rex.ErrorParamCount("exit", "0 or 1")
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
func Command_Break(script *rex.Script, params []*rex.Value) {
	if len(params) != 0 && len(params) != 1 {
		rex.ErrorParamCount("break", "0 or 1")
	}

	script.Break = true
	if len(params) > 0 {
		script.RetVal = params[0]
	}
}

// A conditional version of break.
// 	breakif condition [value]
// Returns value or unchanged.
func Command_BreakIf(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 && len(params) != 2 {
		rex.ErrorParamCount("breakif", "1 or 2")
	}

	script.Break = params[0].Bool()
	if len(params) > 1 {
		script.RetVal = params[1]
	}
}

// Forces a return until it hits a loop or foreach command or the script exits.
// 	breakloop [value]
// Returns value or unchanged.
func Command_BreakLoop(script *rex.Script, params []*rex.Value) {
	if len(params) != 0 && len(params) != 1 {
		rex.ErrorParamCount("breakloop", "0 or 1")
	}

	script.BreakLoop = true
	if len(params) > 0 {
		script.RetVal = params[0]
	}
}

// A conditional version of breakloop.
// 	breakloopif condition [value]
// Returns value or unchanged.
func Command_BreakLoopIf(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 && len(params) != 2 {
		rex.ErrorParamCount("breakloopif", "1 or 2")
	}

	script.BreakLoop = params[0].Bool()
	if len(params) > 1 {
		script.RetVal = params[1]
	}
}

// Runs a code block or (if the value is not code) converts the value to a string and compiles/runs that.
// 	eval code
// Does not halt any exit states except break and return.
// Returns result of running code.
func Command_Eval(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("eval", "1")
	}

	var block *rex.Code
	if params[0].Type == rex.TypCode {
		block = params[0].Data.(*rex.Code)
	} else {
		val, err := script.Host.CompileToValue(params[0].String(), params[0].Pos)
		if err != nil {
			rex.ErrorGeneralCmd("eval", "Compile error: " + err.Error())
		}
		block = val.Data.(*rex.Code)
	}
	
	script.Locals.Add(script.Host, block)
	script.Exec(block)
	script.Return = false
	script.Locals.Remove()
}

// Manipulates the error flag.
// 	error [value]
// If you pass no parameters the error flag will be returned, to set or unset the flag pass a boolean value.
// Returns unchanged or the value of the error flag.
func Command_Error(script *rex.Script, params []*rex.Value) {
	if len(params) != 0 && len(params) != 1 {
		rex.ErrorParamCount("error", "0 or 1")
	}

	if len(params) == 0 {
		script.RetVal = rex.NewValueBool(script.Error)
		return
	}
	script.Error = params[0].Bool()
}

// Runs code if the error flag is true, the error flag is cleared before the code is run.
// 	onerror code
// Returns unchanged.
func Command_OnError(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("onerror", "1")
	}

	if script.Error {
		script.Error = false
		
		if params[0].Type != rex.TypCode {
			rex.ErrorGeneralCmd("onerror", "Attempt to run non-executable Value.")
		}
		block := params[0].Data.(*rex.Code)
		script.Locals.Add(script.Host, block)
		script.Exec(block)
		script.Locals.Remove()
		return
	}
}

// Copies a value.
// 	copy value
// Note that this command may not be useful for some types.
// The new value has invalid position info.
// Returns the new value.
func Command_Copy(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("copy", "1")
	}

	val := new(rex.Value)
	val.Type = params[0].Type
	val.Pos = rex.NewPosition(0, -1, "")

	switch params[0].Type {
	case rex.TypNil:
		val.Data = nil
	case rex.TypString:
		val.Data = params[0].Data.(string)
	case rex.TypInt:
		val.Data = params[0].Data.(int64)
	case rex.TypFloat:
		val.Data = params[0].Data.(float64)
	case rex.TypBool:
		val.Data = params[0].Data.(bool)
	case rex.TypCode:
		val.Data = params[0].Data.(*rex.Code)
	case rex.TypIndex:
		val.Data = params[0].Data.(rex.Indexable)
	case rex.TypUser:
		val.Data = params[0].Data
	default:
		rex.ErrorGeneralCmd("copy", "Script Value has invalid Type.")
	}

	script.RetVal = val
}

// Returns true if an index exists in a map or array.
//	exists value index
// Returns true or false.
func Command_Exists(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("exists", "2")
	}

	if params[0].Type != rex.TypIndex {
		rex.ErrorGeneralCmd("exists", "Attempt to check key existence in non-Indexable object.")
	}

	if params[0].Data.(rex.Indexable).Exists(params[1].String()) {
		script.RetVal = rex.NewValueBool(true)
		return
	}
	script.RetVal = rex.NewValueBool(false)
}

// Fetches the element count of an Indexable.
// 	len value
// Returns the element count.
func Command_Len(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("len", "1")
	}

	if params[0].Type != rex.TypIndex {
		rex.ErrorGeneralCmd("len", "Attempt to fetch non-Indexable object's size.")
	}

	script.RetVal = rex.NewValueInt64(params[0].Data.(rex.Indexable).Len())
}

// Checks if a value is nil.
// 	isnil value
// Returns true or false.
func Command_IsNil(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("isnil", "1")
	}

	script.RetVal = rex.NewValueBool(params[0].Type == rex.TypNil)
}

// Reads or checks value types.
// 	type value [typ_string]
// Returns the type as a string (if called without a type string) else returns true or false.
func Command_Type(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 && len(params) != 2 {
		rex.ErrorParamCount("type", "1 or 2")
	}
	
	if len(params) == 1 {
		script.RetVal = rex.NewValueString(params[0].TypeString())
		return
	}
	
	script.RetVal = rex.NewValueBool(params[0].TypeString() == params[1].String())
}

// If the condition is true run true code else if false code exists call false code.
// 	if condition truecode [falsecode]
// Returns the return value of the last command in the code it runs or unchanged.
func Command_If(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 && len(params) != 3 {
		rex.ErrorParamCount("if", "2 or 3")
	}

	if params[0].Bool() {
		if params[1].Type != rex.TypCode {
			rex.ErrorGeneralCmd("if", "Attempt to run non-executable Value.")
		}
		block := params[1].Data.(*rex.Code)
		script.Locals.Add(script.Host, block)
		script.Exec(block)
		script.Locals.Remove()
		return
	}

	if len(params) > 2 {
		if params[2].Type != rex.TypCode {
			rex.ErrorGeneralCmd("if", "Attempt to run non-executable Value.")
		}
		block := params[2].Data.(*rex.Code)
		script.Locals.Add(script.Host, block)
		script.Exec(block)
		script.Locals.Remove()
		return
	}
}

// Runs code for as long as the code returns true.
// 	loop code
// Returns the return value of the last command in the code it runs, always false unless loop
// exited with ret (In which case the return value is unusable by the command calling loop anyway).
func Command_Loop(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("loop", "1")
	}

	if params[0].Type != rex.TypCode {
		rex.ErrorGeneralCmd("loop", "Attempt to run non-executable Value.")
	}
	block := params[0].Data.(*rex.Code)
	script.Locals.Add(script.Host, block)
	for {
		script.Exec(block)
		script.BreakLoop = false
		if !script.RetVal.Bool() || script.ExitFlags() {
			script.Locals.Remove()
			return
		}
	}
}

// Runs code (end-start)/step times.
// 	for start end step code
// Params for code:
//	index
// If code returns false for exits early.
// code MUST be a block created via a block declaration!
// Does not stop returns, but does work with breakloop.
// A step value of 0 will default to 1.
// Negative step values are allowed, in which case the loop will count down instead of up.
// Returns the return value of the last command in the code.
func Command_For(script *rex.Script, params []*rex.Value) {
	if len(params) != 4 {
		rex.ErrorParamCount("for", "4")
	}

	start := params[0].Int64()
	end := params[1].Int64()
	step := params[2].Int64()
	
	if params[3].Type != rex.TypCode {
		rex.ErrorGeneralCmd("for", "Attempt to run non-executable Value.")
	}
	block := params[3].Data.(*rex.Code)
	script.Locals.Add(script.Host, block)
	
	if step == 0 {
		step = 1
	}
	
	if step > 0 {
		for i := start; i <= end; i += step {
			script.Locals.Set(0, rex.NewValueInt64(i))
			script.Exec(block)
			script.BreakLoop = false
			if !script.RetVal.Bool() || script.ExitFlags() {
				break
			}
		}
	} else {
		for i := start; i >= end; i += step {
			script.Locals.Set(0, rex.NewValueInt64(i))
			script.Exec(block)
			script.BreakLoop = false
			if !script.RetVal.Bool() || script.ExitFlags() {
				break
			}
		}
	}
	script.Locals.Remove()
	return
}

// Runs code for each entry in a map or array value.
// 	foreach objectvalue code
// Params for code:
//	key value
// If code returns false foreach aborts.
// code MUST be a block created via a block declaration!
// Does not stop returns, but does work with breakloop.
// Returns the return value of the last command in code.
func Command_ForEach(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("foreach", "2")
	}

	if params[0].Type != rex.TypIndex {
		rex.ErrorGeneralCmd("foreach", "Attempt to range over non-Indexable.")
	}

	if params[1].Type != rex.TypCode {
		rex.ErrorGeneralCmd("foreach", "Attempt to run non-executable Value.")
	}

	block := params[1].Data.(*rex.Code)
	script.Locals.Add(script.Host, block)
	
	// If the value is an IntIndexable we can use this faster version
	// (hopefully requires less string conversions)
	if v, ok := params[0].Data.(rex.IntIndexable); ok {
		for i := int64(0); i < v.Len(); i++ {
			script.Locals.Set(0, rex.NewValueInt64(i))
			script.Locals.Set(1, v.GetI(i))
			script.Exec(block)
			script.BreakLoop = false
			if !script.RetVal.Bool() || script.ExitFlags() {
				break
			}
		}
		
		script.Locals.Remove()
		return
	}
	
	v := params[0].Data.(rex.Indexable)
	for _, i := range params[0].Data.(rex.Indexable).Keys() {
		script.Locals.Set(0, rex.NewValueString(i))
		script.Locals.Set(1, v.Get(i))
		script.Exec(block)
		script.BreakLoop = false
		if !script.RetVal.Bool() || script.ExitFlags() {
			break
		}
	}
	script.Locals.Remove()
}
