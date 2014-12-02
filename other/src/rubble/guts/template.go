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

import "strings"
import "regexp"
import "dctech/rex"

type NativeTemplate func([]*Value) *Value

var Templates = make(map[string]*Template)

type TemplateParam struct {
	Name    string
	Default string
}

type Template struct {
	// Is this a native template?
	Native bool

	// Is this template's body made up of script code?
	Script bool

	// The native template handler
	Handler NativeTemplate

	// The script template code
	Code *rex.Value

	// The user template text
	Text *Value

	// Param names/defaults for user templates.
	Params []*TemplateParam
}

var repTokenNameRegEx = regexp.MustCompile("%[a-zA-Z_]+")
var repTokenNumberRegEx = regexp.MustCompile("%[1-9]+")

var varNameRegEx = regexp.MustCompile("^\\$\\{?([a-zA-Z_][a-zA-Z0-9_]*)\\}?$")

// Call runs a template. The output is run through the stage parser for (non script) user templates.
func (this *Template) Call(params []*Value) *Value {
	// Handle the ... param.
	if len(params) != 0 {
		if params[len(params)-1].Data == "..." {
			params = append(params[:len(params)-1], PrevParams...)
		} else {
			PrevParams = make([]*Value, len(params))
			copy(PrevParams, params)
		}
	} else {
		PrevParams = make([]*Value, 0)
	}

	for i := range params {
		params[i].Data = strings.TrimSpace(params[i].Data)

		// just expand everything
		params[i].Data = ExpandVars(params[i].Data)
	}

	// Native template
	if this.Native {
		return this.Handler(params)
	}

	// Script template
	if this.Script {
		script := rex.NewScript()
		tmp := make([]*rex.Value, len(params))
		for i := range params {
			tmp[i] = params[i].Script()
		}
		
		rtn, err := GlobalScriptState.RunCommand(script, this.Code, tmp)
		if err != nil {
			panic("Script Error: " + err.Error())
		}
		return NewValueScript(rtn)
	}

	// User template
	out := ExpandVars(this.Text.Data)

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

	return NewValuePos(string(Parse([]byte(out), stgUseCurrent, this.Text.Pos)), this.Text.Pos)

}

func NewNativeTemplate(name string, handler NativeTemplate) {
	rtn := new(Template)
	rtn.Native = true
	rtn.Handler = handler
	Templates[name] = rtn
}

func NewUserTemplate(name string, text *Value, params []*TemplateParam) {
	rtn := new(Template)
	rtn.Text = text
	rtn.Params = params
	Templates[name] = rtn
}

func NewScriptTemplate(name string, code *rex.Value) {
	rtn := new(Template)
	rtn.Script = true
	rtn.Code = code
	Templates[name] = rtn
}

// Built-in template

func tempTemplate(params []*Value) *Value {
	if len(params) < 2 {
		panic("Wrong number of params to !TEMPLATE.")
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

	NewUserTemplate(name, text, parsedParams)

	return NewValue("")
}
