package main

import "dctech/nca6"
import "dctech/nca6/base"
import "dctech/nca6/bit"
import "dctech/nca6/cmp"
import "dctech/nca6/conio"
import "dctech/nca6/csv"
import "dctech/nca6/env"
import "dctech/nca6/file"
import "dctech/nca6/fileio"
import "dctech/nca6/ini"
import "dctech/nca6/math"
import "dctech/nca6/ncash"
import "dctech/nca6/stack"
import "dctech/nca6/str"

import "regexp"

var GlobalNCAState *nca6.State

func InitNCA() {
	state := nca6.NewState()
	
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
	ncash.Setup(state)
	stack.Setup(state)
	str.Setup(state)
	
	state.NewNameSpace("rubble")
	state.NewVar("rubble:dfdir", nca6.NewValueString(DFDir))
	state.NewVar("rubble:outputdir", nca6.NewValueString(OutputDir))
	state.NewVar("rubble:configdir", nca6.NewValueString(ConfigDir))
	state.NewVar("rubble:basedir", nca6.NewValueString(BaseDir))
	state.NewVar("rubble:addonsdir", nca6.NewValueString(AddonsDir))
	
	state.NewNativeCommand("panic", CommandPanic)
	
	state.NewNativeCommand("valueinspect", nca6.CommandValueInspect)
	
	state.NewNameSpace("regex")
	state.NewNativeCommand("regex:replace", CommandRegEx_Replace)
	
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
func CommandPanic(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to panic (how ironic).")
	}

	panic(params[0].String())
}

// Runs a regular expression search and replace.
// 	regex:replace regex input replace
// Returns input with all strings matching regex replaced with replace.
func CommandRegEx_Replace(state *nca6.State, params []*nca6.Value) {
	if len(params) != 3 {
		panic("Wrong number of params to regex:replace.")
	}

	regEx := regexp.MustCompile(params[0].String())
	state.RetVal = nca6.NewValueString(regEx.ReplaceAllString(params[1].String(), params[2].String()))
}

// Makes Rubble skip a file.
// 	rubble:skipfile name
// name is the file's BASE NAME not it's path!
// Returns unchanged.
func CommandRubble_SkipFile(state *nca6.State, params []*nca6.Value) {
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
func CommandRubble_SetVar(state *nca6.State, params []*nca6.Value) {
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
func CommandRubble_GetVar(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:getvar.")
	}
	
	if _, ok := VariableData[params[0].String()]; !ok {
		panic("Rubble variable " + params[0].String() + " does not exist.")
	}
	
	state.RetVal = nca6.NewValueString(VariableData[params[0].String()])
}

// Parses Rubble code.
// 	rubble:stageparse code
// Note that how code is parsed depends on the parse stage.
// Returns the result of running code through the stage parser.
func CommandRubble_StageParse(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:stageparse.")
	}
	
	state.RetVal = nca6.NewValueString(StageParse(params[0].String()))
}

// Calles a Rubble template.
// 	rubble:calltemplate name [params...]
// Returns the templates return value.
func CommandRubble_CallTemplate(state *nca6.State, params []*nca6.Value) {
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
	
	state.RetVal = nca6.NewValueString(Templates[name].Call(strParams))
}

// Expands Rubble variables.
// 	rubble:expandvars raws
// Returns the raws with all Rubble variables expanded.
func CommandRubble_ExpandVars(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:expandvars.")
	}
	
	state.RetVal = nca6.NewValueString(ExpandVars(params[0].String()))
}
