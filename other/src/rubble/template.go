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

import "strings"
import "regexp"
import "dctech/rex"

type TemplateParam struct {
	Name    string
	Default string
}

type Template struct {
	// The script template code (must be type code or type command).
	Code *rex.Value

	// The user template text.
	Text *rex.Value

	// Param names/defaults for user templates.
	Params []*TemplateParam
}

var repTokenNameRegEx = regexp.MustCompile("%[a-zA-Z_]+")

//var varNameRegEx = regexp.MustCompile("^\\$\\{?([a-zA-Z_][a-zA-Z0-9_]*)\\}?$")

// Call runs a template. The output is run through the stage parser for (non script) user templates.
func (temp *Template) Call(state *State, params []*rex.Value) *rex.Value {
	// Handle the ... param.
	if len(params) != 0 {
		if params[len(params)-1].Data == "..." {
			params = append(params[:len(params)-1], state.PrevParams...)
		} else {
			state.PrevParams = make([]*rex.Value, len(params))
			copy(state.PrevParams, params)
		}
	} else {
		state.PrevParams = make([]*rex.Value, 0)
	}

	for i := range params {
		if params[i].Type == rex.TypString {
			// Setting the data field like this is kinda evil...
			params[i].Data = strings.TrimSpace(params[i].String())
			params[i].Data = state.ExpandVars(params[i].String())
		}
	}

	// Script template
	if temp.Code != nil {
		script := rex.NewScript()
		rtn, err := state.ScriptState.RunCommand(script, temp.Code, params)
		if err != nil {
			RaiseError("Script Error: " + err.Error())
		}
		return rtn
	}

	// User template
	out := state.ExpandVars(temp.Text.String())

	// Named replacement tokens
	for i := range temp.Params {
		if len(params) <= i {
			out = strings.Replace(out, "%{" + temp.Params[i].Name + "}", temp.Params[i].Default, -1)
		} else {
			out = strings.Replace(out, "%{" + temp.Params[i].Name + "}", params[i].String(), -1)
		}
	}
	out = repTokenNameRegEx.ReplaceAllStringFunc(out, func(in string) string {
		in = strings.TrimLeft(in, "%")
		for i := range temp.Params {
			if temp.Params[i].Name == in {
				if len(params) > i {
					return params[i].String()
				}
				return temp.Params[i].Default
			}
		}
		return "%" + in
	})

	rtn := rex.NewValueString(string(state.ParseFile([]byte(out), StgUseCurrent, temp.Text.Pos)))
	rtn.Pos = temp.Text.Pos.Copy()
	return rtn
}

func (state *State) NewUserTemplate(name string, text *rex.Value, params []*TemplateParam) {
	rtn := new(Template)
	rtn.Text = text
	rtn.Params = params
	state.Templates[name] = rtn
}

func (state *State) NewScriptTemplate(name string, code *rex.Value) {
	rtn := new(Template)
	rtn.Code = code
	state.Templates[name] = rtn
}
