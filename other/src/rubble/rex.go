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

package rubble

import "dctech/rex"
import "dctech/patch"
import "rubble/rblutil"

import "fmt"
import "strings"

// Causes rubble to abort with an error, use for correctable errors like configuration problems.
// 	rubble:abort msg
// Returns unchanged.
func Command_Abort(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("rubble:abort", "1")
	}

	RaiseAbort(params[0].String())
}

// Returns the name of the current file.
// 	rubble:currentfile
// Returns the file name.
func Command_CurrentFile(script *rex.Script, params []*rex.Value) {
	state := script.Host.FetchVariable("rubble:state").Data.(*State)
	script.RetVal = rex.NewValueString(state.CurrentFile)
}

// Manages Rubble variables.
// 	rubble:configvar name [value]
// Returns unchanged or the variables value.
func Command_ConfigVar(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 && len(params) != 2 {
		rex.ErrorParamCount("rubble:configvar", "1 or 2")
	}
	state := script.Host.FetchVariable("rubble:state").Data.(*State)

	if len(params) == 1 {
		script.RetVal = rex.NewValueString(state.VariableData[params[0].String()])
		return
	}
	state.VariableData[params[0].String()] = params[1].String()
}

// Parses Rubble code.
// 	rubble:stageparse code [stage]
// Note that how code is parsed depends on the parse stage.
// Valid values for stage are:
//	-1 (or just leave it off) to use the current stage
//	3 for preparse
//	4 for parse
//	5 for postparse
// The other stage numbers are not valid for the stage parser.
// Returns the result of running code through the stage parser.
func Command_Parse(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 && len(params) != 2 {
		rex.ErrorParamCount("rubble:stageparse", "1 or 2")
	}
	state := script.Host.FetchVariable("rubble:state").Data.(*State)

	if len(params) == 2 {
		stage := ParseStage(int(params[0].Int64()))
		script.RetVal = rex.NewValueString(string(state.ParseFile([]byte(params[0].String()), stage, params[0].Pos)))
		return
	}
	script.RetVal = rex.NewValueString(string(state.ParseFile([]byte(params[0].String()), StgUseCurrent, params[0].Pos)))
}

// Calls a Rubble template.
// 	rubble:calltemplate name [params...]
// Returns the templates return value.
func Command_CallTemplate(script *rex.Script, params []*rex.Value) {
	if len(params) < 1 {
		rex.ErrorParamCount("rubble:calltemplate", ">0")
	}
	state := script.Host.FetchVariable("rubble:state").Data.(*State)

	name := params[0].String()
	if _, ok := state.Templates[name]; !ok {
		rex.ErrorGeneralCmd("rubble:calltemplate", "Invalid template: "+name)
	}

	script.RetVal = state.Templates[name].Call(state, params[1:])
}

// Expands Rubble variables.
// 	rubble:expandvars raws
// Returns the raws with all Rubble variables expanded.
func Command_ExpandVars(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("rubble:expandvars", "1")
	}
	state := script.Host.FetchVariable("rubble:state").Data.(*State)

	pos := params[0].Pos
	val := rex.NewValueString(state.ExpandVars(params[0].String()))
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
	state := script.Host.FetchVariable("rubble:state").Data.(*State)

	if len(params) == 2 {
		if file, ok := state.Files.Data[params[0].String()]; ok {
			script.RetVal = rex.NewValueBool(file.Tags[params[1].String()])
		} else {
			script.RetVal = rex.NewValueBool(false)
		}
		return
	}

	if file, ok := state.Files.Data[params[0].String()]; ok {
		file.Tags[params[1].String()] = params[2].Bool()
	}
}

// Manages global script file tags.
// 	rubble:gfiletag filename tag [value]
// Returns the tag's state (if called without a value) or returns unchanged.
func Command_GFileTag(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 && len(params) != 3 {
		rex.ErrorParamCount("rubble:gfiletag", "2 or 3")
	}
	state := script.Host.FetchVariable("rubble:state").Data.(*State)

	if len(params) == 2 {
		if file, ok := state.GlobalFiles.Data[params[0].String()]; ok {
			script.RetVal = rex.NewValueBool(file.Tags[params[1].String()])
		} else {
			script.RetVal = rex.NewValueBool(false)
		}
		return
	}

	if file, ok := state.GlobalFiles.Data[params[0].String()]; ok {
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
	state := script.Host.FetchVariable("rubble:state").Data.(*State)

	name := params[0].String()
	for _, addon := range state.Addons {
		if addon.Name == name {
			return
		}
	}
	state.Addons = append(state.Addons, NewAddon(name))
}

// Adds a new file to an addon.
// 	rubble:newfile addon name contents
// Fails silently if the addon does not exist.
// Returns unchanged.
func Command_NewFile(script *rex.Script, params []*rex.Value) {
	if len(params) != 3 {
		rex.ErrorParamCount("rubble:newfile", "3")
	}
	state := script.Host.FetchVariable("rubble:state").Data.(*State)

	name := params[0].String()
	for _, addon := range state.Addons {
		if addon.Name == name {
			addon.Files[params[1].String()] = NewAddonFile(params[1].String(), "(from script: "+state.CurrentFile+") ", []byte(params[2].String()))
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
	state := script.Host.FetchVariable("rubble:state").Data.(*State)

	state.NewScriptTemplate(params[0].String(), params[1])
}

// Defines a Rubble user template.
// 	rubble:usertemplate name [params...] code
// Returns unchanged.
func Command_UserTemplate(script *rex.Script, params []*rex.Value) {
	if len(params) < 1 {
		rex.ErrorParamCount("rubble:usertemplate", ">1")
	}
	state := script.Host.FetchVariable("rubble:state").Data.(*State)

	name := params[0].String()
	text := params[len(params)-1]
	paramNames := params[1 : len(params)-1]

	parsedParams := make([]*TemplateParam, 0, len(paramNames))
	for _, val := range paramNames {
		rtn := new(TemplateParam)
		if strings.Contains(val.String(), "=") {
			parts := strings.SplitN(val.String(), "=", 2)
			rtn.Name = parts[0]
			rtn.Default = parts[1]
			parsedParams = append(parsedParams, rtn)
			continue
		}
		rtn.Name = val.String()
		parsedParams = append(parsedParams, rtn)
	}

	state.NewUserTemplate(name, text, parsedParams)
}

// Exactly like rubble:usertemplate, just it returns an empty string.
// Cannot be called directly by scripts.
// Returns "".
func userTemplateWrap(script *rex.Script, params []*rex.Value) {
	Command_UserTemplate(script, params)
	script.RetVal = rex.NewValueString("")
}

// Script commands related to patching and loading addons from scripts.

// Applies a patch to a string.
// 	rubble:patch string patch
// Returns the patched text.
func Command_Patch(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("rubble:patch", "2")
	}

	patches := patch.FromText(params[1].String())
	if patches == nil {
		rex.ErrorGeneralCmd("rubble:patch", "Invalid patch string.")
	}

	text, applied := patch.Apply(patches, params[0].String())

	count := 0
	for i := range applied {
		if applied[i] {
			count++
		}
	}
	if len(patches) != count {
		rex.ErrorGeneralCmd("rubble:patch", fmt.Sprintf("Not all patches applied: %v out of %v", count, len(patches)))
	}

	script.RetVal = rex.NewValueString(text)
}

// Decompresses a Rubble encoded string.
// 	rubble:decompress string
// Returns the decoded and decompressed text.
func Command_Decompress(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("rubble:decompress", "1")
	}

	defer func() {
		if x := recover(); x != nil {
			if a, ok := x.(error); ok {
				rex.ErrorGeneralCmd("rubble:decompress", a.Error())
			}
			if a, ok := x.(string); ok {
				rex.ErrorGeneralCmd("rubble:decompress", a)
			}
		}
	}()

	data := string(rblutil.Decompress(rblutil.Decode(rblutil.StripWS([]byte(params[0].String())))))
	script.RetVal = rex.NewValueString(data)
}

// Compresses a string using the Rubble encoding.
// 	rubble:compress string
// Returns the encoded and compressed text.
func Command_Compress(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("rubble:compress", "1")
	}

	defer func() {
		if x := recover(); x != nil {
			if a, ok := x.(error); ok {
				rex.ErrorGeneralCmd("rubble:compress", a.Error())
			}
			if a, ok := x.(string); ok {
				rex.ErrorGeneralCmd("rubble:compress", a)
			}
		}
	}()

	data := string(rblutil.Split(rblutil.Encode(rblutil.Compress([]byte(params[0].String())))))
	script.RetVal = rex.NewValueString(data)
}
