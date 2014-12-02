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
import "dctech/raptor/commands/sort"
import "dctech/raptor/commands/str"
import "dctech/raptor/commands/thread"

import "regexp"
import "strings"
import "os"

var GlobalRaptorState *raptor.State

func InitScripting() {
	state := raptor.NewState()

	state.NoRecover = NoRecover

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
	sort.Setup(state)
	str.Setup(state)
	thread.Setup(state)
	
	state.NewNameSpace("rubble")
	state.NewNamespacedVar("rubble:version", raptor.NewValueString(RubbleVersion))
	state.NewNamespacedVar("rubble:dfdir", raptor.NewValueString(DFDir))
	state.NewNamespacedVar("rubble:outputdir", raptor.NewValueString(OutputDir))
	state.NewNamespacedVar("rubble:addonsdir", raptor.NewValueString(AddonsDir))
	state.NewNamespacedVar("rubble:raws", raptor.NewValueObject(NewIndexableRaws()))
	
	array := make([]*raptor.Value, 0, 20)
	for i := range Addons {
		if Addons[i].Active == true {
			array = append(array, raptor.NewValueString(Addons[i].Name))
		}
	}
	state.NewNamespacedVar("rubble:activeaddons", raptor.NewValueObject(raptor.NewParamsArray(array)))
	
	state.NewNativeCommand("panic", CommandPanic)
	
	state.NewNativeCommand("rubble:abort", CommandRubble_Abort)
	
	state.NewNativeCommand("rubble:skipfile", CommandRubble_SkipFile)
	state.NewNativeCommand("rubble:setvar", CommandRubble_SetVar)
	state.NewNativeCommand("rubble:getvar", CommandRubble_GetVar)
	
	state.NewNativeCommand("rubble:stageparse", CommandRubble_Parse)
	state.NewNativeCommand("rubble:calltemplate", CommandRubble_CallTemplate)
	state.NewNativeCommand("rubble:expandvars", CommandRubble_ExpandVars)
	
	state.NewNativeCommand("rubble:template", CommandRubble_Template)
	
	state.NewNativeCommand("rubble:addonactive", CommandRubble_AddonActive)

	// Redirect output to logger
	state.Output = logFile

	GlobalRaptorState = state
}

// Causes a panic.
// 	panic value
// Returns unchanged.
func CommandPanic(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to panic (how ironic).")
	}

	panic(params[0].String())
}

// Causes rubble to abort with an error, use for correctable errors like configuration problems.
// 	rubble:abort msg
// Returns unchanged.
func CommandRubble_Abort(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:abort (how ironic).")
	}
	
	LogPrintln("Abort:", params[0].String())
	os.Exit(1) // Evil, evil, evil, but it works
}

// Makes Rubble skip a file.
// 	rubble:skipfile name
// name is the file's BASE NAME not it's path!
// Returns unchanged.
func CommandRubble_SkipFile(script *raptor.Script, params []*raptor.Value) {
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
func CommandRubble_SetVar(script *raptor.Script, params []*raptor.Value) {
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
func CommandRubble_GetVar(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:getvar.")
	}

	if _, ok := VariableData[params[0].String()]; !ok {
		panic("Rubble variable " + params[0].String() + " does not exist.")
	}

	script.RetVal = raptor.NewValueString(VariableData[params[0].String()])
}

// Parses Rubble code.
// 	rubble:stageparse code [stage]
// Note that how code is parsed depends on the parse stage.
// Valid values for stage are:
//	0 (or just leave it off) to use the current stage
//	3 for preparse
//	4 for parse
//	5 for postparse
// The other stage numbers are not valid for the stage parser.
// Returns the result of running code through the stage parser.
func CommandRubble_Parse(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 && len(params) != 2 {
		panic("Wrong number of params to rubble:stageparse.")
	}

	if len(params) == 2 {
		stage := int(params[0].Int64())
		script.RetVal = raptor.NewValueString(string(Parse([]byte(params[0].String()), stage)))
		return
	}
	script.RetVal = raptor.NewValueString(string(Parse([]byte(params[0].String()), stgUseCurrent)))
}

// Calls a Rubble template.
// 	rubble:calltemplate name [params...]
// Returns the templates return value.
func CommandRubble_CallTemplate(script *raptor.Script, params []*raptor.Value) {
	if len(params) < 1 {
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

	script.RetVal = raptor.NewValueString(Templates[name].Call(strParams))
}

// Expands Rubble variables.
// 	rubble:expandvars raws
// Returns the raws with all Rubble variables expanded.
func CommandRubble_ExpandVars(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:expandvars.")
	}

	script.RetVal = raptor.NewValueString(ExpandVars(params[0].String()))
}

// Defines a Rubble script template.
// 	rubble:template name [params...] code
// Returns unchanged.
func CommandRubble_Template(script *raptor.Script, params []*raptor.Value) {
	if len(params) < 2 {
		panic("Wrong number of params to rubble:template.")
	}
	
	name := params[0].String()
	code := params[len(params)-1].CompiledScript()
	paramNames := params[1:len(params)-1]

	parsedParams := make([]*TemplateParam, 0, len(paramNames))

	for i := range paramNames {
		val := paramNames[i].String()
		rtn := new(TemplateParam)
		if strings.Contains(val, "=") {
			parts := strings.SplitN(val, "=", 2)
			rtn.Name = parts[0]
			rtn.Default = parts[1]
			parsedParams = append(parsedParams, rtn)
			continue
		}
		rtn.Name = val
		parsedParams = append(parsedParams, rtn)
	}

	NewScriptTemplate(name, code, parsedParams)
}

// Checks if a named addon is active.
// 	rubble:addonactive name
// name needs to be the FULL name of the addon! (for example "Base/Files")
// The name check is case-sensitive!
// Returns true or false.
func CommandRubble_AddonActive(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to rubble:addonactive.")
	}
	
	name := params[0].String()
	for _, val := range Addons {
		if val.Name == name {
			if val.Active {
				script.RetVal = raptor.NewValueBool(true)
				return
			}
			script.RetVal = raptor.NewValueBool(false)
			return
		}
		
	}
	script.RetVal = raptor.NewValueBool(false)
}
