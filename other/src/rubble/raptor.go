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

import "dctech/raptor"
import "dctech/raptor/commands/base"
import "dctech/raptor/commands/bit"
import "dctech/raptor/commands/boolean"
import "dctech/raptor/commands/console"
import "dctech/raptor/commands/convert"
import "dctech/raptor/commands/csv"
import "dctech/raptor/commands/debug"
import "dctech/raptor/commands/env"
import "dctech/raptor/commands/file"
import "dctech/raptor/commands/fileio"
import "dctech/raptor/commands/float"
import "dctech/raptor/commands/ini"
import "dctech/raptor/commands/integer"
import "dctech/raptor/commands/raw"
import "dctech/raptor/commands/regex"
import "dctech/raptor/commands/str"

import "regexp"

var GlobalRaptorState *raptor.State

func InitScripting() {
	state := raptor.NewState()

	state.NoRecover = !NCARecover

	// Load all commands
	base.Setup(state)
	bit.Setup(state)
	boolean.Setup(state)
	console.Setup(state)
	convert.Setup(state)
	csv.Setup(state)
	debug.Setup(state)
	env.Setup(state)
	file.Setup(state)
	fileio.Setup(state)
	float.Setup(state)
	ini.Setup(state)
	integer.Setup(state)
	raw.Setup(state)
	regex.Setup(state)
	str.Setup(state)

	state.NewNameSpace("rubble")
	state.NewVar("rubble:dfdir", raptor.NewValueString(DFDir))
	state.NewVar("rubble:outputdir", raptor.NewValueString(OutputDir))
	state.NewVar("rubble:addonsdir", raptor.NewValueString(AddonsDir))
	state.NewVar("rubble:raws", raptor.NewValueObject(NewIndexableRaws()))

	state.NewNativeCommand("panic", CommandPanic)

	state.NewNativeCommand("rubble:skipfile", CommandRubble_SkipFile)
	state.NewNativeCommand("rubble:setvar", CommandRubble_SetVar)
	state.NewNativeCommand("rubble:getvar", CommandRubble_GetVar)

	state.NewNativeCommand("rubble:stageparse", CommandRubble_Parse)
	state.NewNativeCommand("rubble:calltemplate", CommandRubble_CallTemplate)
	state.NewNativeCommand("rubble:expandvars", CommandRubble_ExpandVars)

	// Redirect NCA output to logger
	state.Output = logFile

	GlobalRaptorState = state
}

// Causes a panic.
// 	panic value
// Returns unchanged.
func CommandPanic(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to panic (how ironic).")
	}

	panic(params[0].String())
}

// Makes Rubble skip a file.
// 	rubble:skipfile name
// name is the file's BASE NAME not it's path!
// Returns unchanged.
func CommandRubble_SkipFile(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:skipfile.")
	}

	if _, ok := Files.Files[params[0].String()]; !ok {
		panic("rubble:skipfile: \"" + params[0].String() + "\" is not the name of a loaded raw file.")
	}

	Files.Files[params[0].String()].Skip = true
}

var varNameValidateRegEx = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")

// Sets a Rubble variable.
// 	rubble:setvar name value
// Returns unchanged.
func CommandRubble_SetVar(state *raptor.State, params []*raptor.Value) {
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
func CommandRubble_GetVar(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:getvar.")
	}

	if _, ok := VariableData[params[0].String()]; !ok {
		panic("Rubble variable " + params[0].String() + " does not exist.")
	}

	state.RetVal = raptor.NewValueString(VariableData[params[0].String()])
}

// Parses Rubble code.
// 	rubble:stageparse code [stage]
// Note that how code is parsed depends on the parse stage.
// Valid values for stage are:
//	0 (or just leave it off) for use the current stage
//	3 for preparse
//	4 for parse
//	5 for postparse
// The other stage numbers are not valid for the stage parser.
// Returns the result of running code through the stage parser.
func CommandRubble_Parse(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic("Wrong number of params to rubble:stageparse.")
	}

	if len(params) == 2 {
		stage := int(params[0].Int64())
		state.RetVal = raptor.NewValueString(string(Parse([]byte(params[0].String()), stage)))
		return
	}
	state.RetVal = raptor.NewValueString(string(Parse([]byte(params[0].String()), stgUseCurrent)))
}

// Calles a Rubble template.
// 	rubble:calltemplate name [params...]
// Returns the templates return value.
func CommandRubble_CallTemplate(state *raptor.State, params []*raptor.Value) {
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

	state.RetVal = raptor.NewValueString(Templates[name].Call(strParams))
}

// Expands Rubble variables.
// 	rubble:expandvars raws
// Returns the raws with all Rubble variables expanded.
func CommandRubble_ExpandVars(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:expandvars.")
	}

	state.RetVal = raptor.NewValueString(ExpandVars(params[0].String()))
}
