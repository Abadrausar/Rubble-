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

type NativeTemplate func(*State, []*Value) *Value

type TemplateParam struct {
	Name    string
	Default string
}

type Template struct {
	// The native template handler.
	Handler NativeTemplate

	// The script template code.
	Code *rex.Value

	// The user template text.
	Text *Value

	// Param names/defaults for user templates.
	Params []*TemplateParam
}

var repTokenNameRegEx = regexp.MustCompile("%[a-zA-Z_]+")

//var varNameRegEx = regexp.MustCompile("^\\$\\{?([a-zA-Z_][a-zA-Z0-9_]*)\\}?$")

// Call runs a template. The output is run through the stage parser for (non script) user templates.
func (this *Template) Call(state *State, params []*Value) *Value {
	// Handle the ... param.
	if len(params) != 0 {
		if params[len(params)-1].Data == "..." {
			params = append(params[:len(params)-1], state.PrevParams...)
		} else {
			state.PrevParams = make([]*Value, len(params))
			copy(state.PrevParams, params)
		}
	} else {
		state.PrevParams = make([]*Value, 0)
	}

	for i := range params {
		params[i].Data = strings.TrimSpace(params[i].Data)

		// just expand everything
		params[i].Data = state.ExpandVars(params[i].Data)
	}

	// Native template
	if this.Handler != nil {
		return this.Handler(state, params)
	}

	// Script template
	if this.Code != nil {
		script := rex.NewScript()
		tmp := make([]*rex.Value, len(params))
		for i := range params {
			tmp[i] = params[i].Script()
		}

		rtn, err := state.ScriptState.RunCommand(script, this.Code, tmp)
		if err != nil {
			panic(Abort("Script Error: " + err.Error()))
		}
		return NewValueScript(rtn)
	}

	// User template
	out := state.ExpandVars(this.Text.Data)

	// Named replacement tokens
	for i := range this.Params {
		if len(params) <= i {
			out = strings.Replace(out, "%{"+this.Params[i].Name+"}", this.Params[i].Default, -1)
		} else {
			out = strings.Replace(out, "%{"+this.Params[i].Name+"}", params[i].Data, -1)
		}
	}
	out = repTokenNameRegEx.ReplaceAllStringFunc(out, func(in string) string {
		in = strings.TrimLeft(in, "%")
		for i := range this.Params {
			if this.Params[i].Name == in {
				if len(params) > i {
					return params[i].Data
				}
				return this.Params[i].Default
			}
		}
		return "%" + in
	})

	return NewValuePos(string(state.ParseFile([]byte(out), StgUseCurrent, this.Text.Pos)), this.Text.Pos)
}

func (state *State) NewNativeTemplate(name string, handler NativeTemplate) {
	rtn := new(Template)
	rtn.Handler = handler
	state.Templates[name] = rtn
}

func (state *State) NewUserTemplate(name string, text *Value, params []*TemplateParam) {
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

// Built-in template

func tempTemplate(state *State, params []*Value) *Value {
	if len(params) < 2 {
		panic(Abort("Wrong number of params to !TEMPLATE."))
	}

	name := params[0].Data
	text := params[len(params)-1]
	paramNames := params[1 : len(params)-1]

	parsedParams := make([]*TemplateParam, 0, len(paramNames))

	for _, val := range paramNames {
		rtn := new(TemplateParam)
		if strings.Contains(val.Data, "=") {
			parts := strings.SplitN(val.Data, "=", 2)
			rtn.Name = parts[0]
			rtn.Default = parts[1]
			parsedParams = append(parsedParams, rtn)
			continue
		}
		rtn.Name = val.Data
		parsedParams = append(parsedParams, rtn)
	}

	state.NewUserTemplate(name, text, parsedParams)

	return NewValue("")
}
