/*
Copyright 2013 by Milo Christiansen

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

package main

import "dctech/nca7"
import "dctech/nca7/commands/base"
import "dctech/nca7/commands/bit"
import "dctech/nca7/commands/console"
import "dctech/nca7/commands/csv"
import "dctech/nca7/commands/debug"
import "dctech/nca7/commands/env"
import "dctech/nca7/commands/file"
import "dctech/nca7/commands/fileio"
import "dctech/nca7/commands/ini"
import "dctech/nca7/commands/integer"
import "dctech/nca7/commands/regex"
import "dctech/nca7/commands/stack"
import "dctech/nca7/commands/str"

import "regexp"

var GlobalNCAState *nca7.State

func InitNCA() {
	state := nca7.NewState()
	
	state.NoRecover = !NCARecover

	// Load all commands
	base.Setup(state)
	bit.Setup(state)
	console.Setup(state)
	csv.Setup(state)
	debug.Setup(state)
	env.Setup(state)
	file.Setup(state)
	fileio.Setup(state)
	ini.Setup(state)
	integer.Setup(state)
	regex.Setup(state)
	stack.Setup(state)
	str.Setup(state)
	
	state.NewNameSpace("rubble")
	state.NewVar("rubble:dfdir", nca7.NewValueString(DFDir))
	state.NewVar("rubble:outputdir", nca7.NewValueString(OutputDir))
	state.NewVar("rubble:configdir", nca7.NewValueString(ConfigDir))
	state.NewVar("rubble:basedir", nca7.NewValueString(BaseDir))
	state.NewVar("rubble:addonsdir", nca7.NewValueString(AddonsDir))
	
	state.NewNativeCommand("panic", CommandPanic)
	
	state.NewNativeCommand("rubble:skipfile", CommandRubble_SkipFile)
	state.NewNativeCommand("rubble:setvar", CommandRubble_SetVar)
	state.NewNativeCommand("rubble:getvar", CommandRubble_GetVar)
	
	state.NewNativeCommand("rubble:stageparse", CommandRubble_StageParse)
	state.NewNativeCommand("rubble:calltemplate", CommandRubble_CallTemplate)
	state.NewNativeCommand("rubble:expandvars", CommandRubble_ExpandVars)
	
	GlobalNCAState = state
}

// Causes a panic.
// 	panic value
// Returns unchanged.
func CommandPanic(state *nca7.State, params []*nca7.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to panic (how ironic).")
	}

	panic(params[0].String())
}

// Makes Rubble skip a file.
// 	rubble:skipfile name
// name is the file's BASE NAME not it's path!
// Returns unchanged.
func CommandRubble_SkipFile(state *nca7.State, params []*nca7.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:skipfile.")
	}

	if _, ok := RawFiles[params[0].String()]; !ok {
		panic("rubble:skipfile: \"" + params[0].String() + "\" is not the name of a loaded raw file.")
	}
	
	RawFiles[params[0].String()].Skip = true
}

var varNameValidateRegEx = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")
// Sets a Rubble variable.
// 	rubble:setvar name value
// Returns unchanged.
func CommandRubble_SetVar(state *nca7.State, params []*nca7.Value) {
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
// Returns the value.
func CommandRubble_GetVar(state *nca7.State, params []*nca7.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:getvar.")
	}
	
	if _, ok := VariableData[params[0].String()]; !ok {
		panic("Rubble variable " + params[0].String() + " does not exist.")
	}
	
	state.RetVal = nca7.NewValueString(VariableData[params[0].String()])
}

// Parses Rubble code.
// 	rubble:stageparse code
// Note that how code is parsed depends on the parse stage.
// Returns the result of running code through the stage parser.
func CommandRubble_StageParse(state *nca7.State, params []*nca7.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:stageparse.")
	}
	
	state.RetVal = nca7.NewValueString(StageParse(params[0].String()))
}

// Calles a Rubble template.
// 	rubble:calltemplate name [params...]
// Returns the templates return value.
func CommandRubble_CallTemplate(state *nca7.State, params []*nca7.Value) {
	if len(params) > 1 {
		panic("Wrong number of params to rubble:calltemplate.")
	}
	name := params[0].String()
	strParams := make([]string, 0, len(params)-1)
	for _, val := range params[1:] {
		strParams = append(strParams, val.String())
	}
	
	if _, ok := Templates[name]; !ok {
		panic("Invalid template: " + name)
	}
	
	state.RetVal = nca7.NewValueString(Templates[name].Call(strParams))
}

// Expands Rubble variables.
// 	rubble:expandvars raws
// Returns the raws with all Rubble variables expanded.
func CommandRubble_ExpandVars(state *nca7.State, params []*nca7.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:expandvars.")
	}
	
	state.RetVal = nca7.NewValueString(ExpandVars(params[0].String()))
}
