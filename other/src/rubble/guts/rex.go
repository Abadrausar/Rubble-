/*
Copyright 2013-2014 by Milo Christiansen

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

package guts

import "dctech/rex"
import "dctech/rex/genii"
import "dctech/rex/commands/base"
import "dctech/rex/commands/boolean"
import "dctech/rex/commands/console"
import "dctech/rex/commands/convert"
import "dctech/rex/commands/debug"
//import "dctech/rex/commands/env"
import "dctech/rex/commands/expr"
//import "dctech/rex/commands/file"
//import "dctech/rex/commands/fileio"
import "dctech/rex/commands/float"
import "dctech/rex/commands/integer"
import "dctech/rex/commands/png"
import "dctech/rex/commands/regex"
import "dctech/rex/commands/sort"
import "dctech/rex/commands/str"

import "dctech/dfrex/dfraw"

import "dctech/axis/axisrex"

import "regexp"
import "os"

var GlobalScriptState *rex.State

func InitScriptingEarly() {
	state := rex.NewState()

	state.NoRecover = NoRecover

	// Load most commands
	rex.SetupArrays(state)
	rex.SetupMaps(state)
	base.Setup(state)
	boolean.Setup(state)
	console.Setup(state)
	convert.Setup(state)
	debug.Setup(state)
	//env.Setup(state)
	expr.Setup(state)
	//file.Setup(state)
	//fileio.Setup(state)
	float.Setup(state)
	integer.Setup(state)
	png.Setup(state)
	regex.Setup(state)
	sort.Setup(state)
	str.Setup(state)

	genii.Setup(state)

	dfraw.Setup(state)

	axisrex.Setup(state)

	rbl := state.RegisterModule("rubble")

	rbl.RegisterVar("version", rex.NewValueString(RubbleVersions[0]))
	versions := make(map[string]*rex.Value, len(RubbleVersions))
	for i := range RubbleVersions {
		versions[RubbleVersions[i]] = rex.NewValueBool(true)
	}
	rbl.RegisterVar("versions", rex.NewValueIndex(rex.NewStaticMap(versions)))

	rbl.RegisterVar("fs", rex.NewValueUser(FS))

	rbl.RegisterVar("addons", genii.New(&Addons))
	
	rbl.RegisterCommand("abort", Command_Abort)

	rbl.RegisterCommand("setvar", Command_SetVar)
	rbl.RegisterCommand("getvar", Command_GetVar)

	rbl.RegisterCommand("filetag", Command_FileTag)
	
	rbl.RegisterCommand("newaddon", Command_NewAddon)
	rbl.RegisterCommand("newfile", Command_NewFile)
	
	// Redirect output to logger
	state.Output = logFile

	GlobalScriptState = state
}

func InitScriptingLate() {
	rbl := GlobalScriptState.FetchModule("rubble")

	rbl.RegisterVar("files", genii.New(&Files))

	rbl.RegisterVar("raws", rex.NewValueIndex(NewIndexableRaws()))

	rbl.RegisterCommand("currentfile", Command_CurrentFile)

	rbl.RegisterCommand("template", Command_Template)
	rbl.RegisterCommand("stageparse", Command_Parse)
	rbl.RegisterCommand("calltemplate", Command_CallTemplate)
	rbl.RegisterCommand("expandvars", Command_ExpandVars)
}

// Causes rubble to abort with an error, use for correctable errors like configuration problems.
// 	rubble:abort msg
// Returns unchanged.
func Command_Abort(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("rubble:abort", "1")
	}

	LogPrintln("Abort:", params[0].String())
	os.Exit(1) // Evil, evil, evil, but it works
}

// Returns the name of the current file.
// 	rubble:currentfile
// Returns the file name.
func Command_CurrentFile(script *rex.Script, params []*rex.Value) {
	script.RetVal = rex.NewValueString(CurrentFile)
}

var varNameValidateRegEx = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")

// Sets a Rubble variable.
// 	rubble:setvar name value
// Returns unchanged.
func Command_SetVar(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("rubble:setvar", "2")
	}

	if !varNameValidateRegEx.MatchString(params[0].String()) {
		rex.ErrorGeneralCmd("rubble:setvar", "Config variable name is invalid.")
	}

	VariableData[params[0].String()] = params[1].String()
}

// Gets the value of a Rubble variable.
// 	rubble:getvar name
// Returns the value.
func Command_GetVar(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("rubble:getvar", "1")
	}

	if _, ok := VariableData[params[0].String()]; !ok {
		script.RetVal = rex.NewValueString("")
	}

	script.RetVal = rex.NewValueString(VariableData[params[0].String()])
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
func Command_Parse(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 && len(params) != 2 {
		rex.ErrorParamCount("rubble:stageparse", "1 or 2")
	}

	if len(params) == 2 {
		stage := int(params[0].Int64())
		script.RetVal = rex.NewValueString(string(Parse([]byte(params[0].String()), stage,
			NewPositionScript(params[0].Pos))))
		return
	}
	script.RetVal = rex.NewValueString(string(Parse([]byte(params[0].String()), stgUseCurrent,
		NewPositionScript(params[0].Pos))))
}

// Calls a Rubble template.
// 	rubble:calltemplate name [params...]
// Returns the templates return value.
func Command_CallTemplate(script *rex.Script, params []*rex.Value) {
	if len(params) < 1 {
		rex.ErrorParamCount("rubble:calltemplate", ">0")
	}
	name := params[0].String()

	strParams := make([]*Value, 0, len(params)-1)

	for _, val := range params[1:] {
		strParams = append(strParams, NewValueScript(val))
	}

	if _, ok := Templates[name]; !ok {
		rex.ErrorGeneralCmd("rubble:calltemplate", "Invalid template: " + name)
	}

	script.RetVal = Templates[name].Call(strParams).Script()
}

// Expands Rubble variables.
// 	rubble:expandvars raws
// Returns the raws with all Rubble variables expanded.
func Command_ExpandVars(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("rubble:expandvars", "1")
	}

	pos := params[0].Pos
	val := rex.NewValueString(ExpandVars(params[0].String()))
	val.Pos = pos
	script.RetVal = val
}

// Manages file tags.
// 	rubble:filetag filename tag [value]
// Returns the tag's state (if called without a value) or returns unchanged.
func Command_FileTag(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 && len(params) != 3 {
		rex.ErrorParamCount("rubble:filetag", "2 or 3")
	}
	
	if len(params) == 2 {
		
		if file, ok := Files.Data[params[0].String()]; ok {
			script.RetVal = rex.NewValueBool(file.Tags[params[1].String()])
		} else {
			script.RetVal = rex.NewValueBool(false)
		}
		return
	}

	if file, ok := Files.Data[params[0].String()]; ok {
		file.Tags[params[1].String()] = params[2].Bool()
	}
}

// Creates an empty addon object and adds it to the addon list.
// 	rubble:newaddon name
// Does nothing if the addon already exists.
// Returns the addon.
func Command_NewAddon(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("rubble:newaddon", "1")
	}

	name := params[0].String()
	for _, addon := range Addons {
		if addon.Name == name {
			return
		}
	}
	Addons = append(Addons, NewAddon(name))
}

// Adds a new file to an addon.
// 	rubble:newfile addon name contents
// Fails silently if the addon does not exist.
// Returns unchanged.
func Command_NewFile(script *rex.Script, params []*rex.Value) {
	if len(params) != 3 {
		rex.ErrorParamCount("rubble:newfile", "3")
	}

	name := params[0].String()
	for _, addon := range Addons {
		if addon.Name == name {
			addon.Files[params[1].String()] = NewAddonFile(params[1].String(), "(from script: " + CurrentFile + ") ", []byte(params[2].String()))
			return
		}
	}
	
}

// Defines a Rubble script template.
// 	rubble:template name code
// code MUST be a block created via a block declaration!
// Parameter names, count, and default values is determined by the block meta-data.
// Returns unchanged.
func Command_Template(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("rubble:template", "2")
	}

	NewScriptTemplate(params[0].String(), params[1])
}
