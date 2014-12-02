package main

import "dctech/nca4"
import "dctech/nca4/base"
import "dctech/nca4/bit"
import "dctech/nca4/cmp"
import "dctech/nca4/conio"
import "dctech/nca4/csv"
import "dctech/nca4/env"
import "dctech/nca4/file"
import "dctech/nca4/fileio"
import "dctech/nca4/ini"
import "dctech/nca4/math"
import "dctech/nca4/stack"
import "dctech/nca4/str"

var GlobalNCAState *nca4.State

func InitNCA() {
	state := nca4.NewState()
	
	state.NoRecover = !Recover

	// Load all commands
	base.Setup(state)
	bit.Setup(state)
	cmp.Setup(state)
	conio.Setup(state)
	csv.Setup(state)
	env.Setup(state)
	file.Setup(state)
	fileio.Setup(state)
	ini.Setup(state)
	math.Setup(state)
	stack.Setup(state)
	str.Setup(state)
	
	state.NewNameSpace("rubble")
	state.NewVar("rubble:outputdir", nca4.NewValue(OutputDir))
	state.NewVar("rubble:configdir", nca4.NewValue(ConfigDir))
	state.NewVar("rubble:basedir", nca4.NewValue(BaseDir))
	state.NewVar("rubble:addonsdir", nca4.NewValue(AddonsDir))
	
	state.NewNativeCommand("panic", CommandPanic)
	
	state.NewNativeCommand("rubble:skipfile", CommandRubble_SkipFile)
	state.NewNativeCommand("rubble:setvar", CommandRubble_SetVar)
	state.NewNativeCommand("rubble:getvar", CommandRubble_GetVar)
	
	GlobalNCAState = state
}

// Causes a panic.
// 	panic value
// Returns unchanged.
func CommandPanic(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to panic (how ironic).")
	}

	panic(params[0].String())
}

// Makes Rubble skip a file.
// 	rubble:skipfile name
// name is the file's BASE NAME not it's path!
// Returns unchanged.
func CommandRubble_SkipFile(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:skipfile.")
	}

	if _, ok := RawFiles[params[0].String()]; !ok {
		panic("rubble:skipfile: \"" + params[0].String() + "\" is not the name of a loaded raw file.")
	}
	
	RawFiles[params[0].String()].Skip = true
}

// Sets a Rubble variable.
// 	rubble:setvar name value
// Returns unchanged.
func CommandRubble_SetVar(state *nca4.State, params []*nca4.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to rubble:setvar.")
	}
	
	if !varNameValidateRegEx.MatchString(params[0].String()) {
		panic("Variable name supplied to rubble:setvar is invalid.")
	}
	
	VariableData[params[0].String()] = params[1].String()
}

// Gets the value of a Rubble variable.
// 	rubble:getvar name
// Returns the vlaue.
func CommandRubble_GetVar(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:getvar.")
	}
	
	if _, ok := VariableData[params[0].String()]; !ok {
		panic("Rubble variable " + params[0].String() + " does not exist.")
	}
	
	state.RetVal = nca4.NewValue(VariableData[params[0].String()])
}
