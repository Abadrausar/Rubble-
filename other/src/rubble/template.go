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

import "strings"
import "regexp"
import "strconv"
import "dctech/raptor"

type NativeTemplate func([]string) string

var Templates = make(map[string]*Template)

type TemplateParam struct {
	Name    string
	Default string
}

type Template struct {
	// Is this a native template?
	Native bool

	// Is this template's body made up of NCA code?
	NCA bool

	// The native template handler
	Handler NativeTemplate

	// The user template text
	Text string

	// Param names/defaults for user commands
	Params []*TemplateParam
}

var repTokenNameRegEx = regexp.MustCompile("%[a-zA-Z_]+")
var repTokenNumberRegEx = regexp.MustCompile("%[1-9]+")

var varNameRegEx = regexp.MustCompile("^\\$\\{?([a-zA-Z_][a-zA-Z0-9_]*)\\}?$")

// Call runs a template. The output is run through the stage parser.
func (this *Template) Call(params []string) string {
	// prep params
	if len(params) != 0 {
		if params[len(params)-1] == "..." {
			params = append(params[:len(params)-1], PrevParams...)
		} else {
			PrevParams = make([]string, len(params))
			copy(PrevParams, params)
		}
	} else {
		PrevParams = make([]string, 0)
	}

	for i := range params {
		params[i] = strings.TrimSpace(params[i])

		// If the param is a variable name replace with value
		// Does not replace variables embeded in other text
		name := varNameRegEx.FindStringSubmatch(params[i])
		if name != nil {
			params[i] = varNameRegEx.ReplaceAllString(params[i], VariableData[name[1]])
		}
	}

	// Native template
	if this.Native {
		return this.Handler(params)
	}

	// Script template
	if this.NCA {
		GlobalRaptorState.Code.Add(this.Text)
		GlobalRaptorState.Envs.Add(raptor.NewEnvironment())

		// Handle params
		if len(this.Params) == 1 && this.Params[0].Name == "..." {
			GlobalRaptorState.AddParams(params...)
		} else {
			for i := range this.Params {
				if len(params) <= i {
					GlobalRaptorState.NewVar(this.Params[i].Name, raptor.NewValueString(this.Params[i].Default))
				} else {
					GlobalRaptorState.NewVar(this.Params[i].Name, raptor.NewValueString(params[i]))
				}
			}
		}

		rtn, err := GlobalRaptorState.Run()
		if err != nil {
			panic("Script Error: " + err.Error())
		}

		GlobalRaptorState.Envs.Remove()

		if rtn == nil {
			return ""
		}
		return rtn.String()
	}

	// User template
	out := ExpandVars(this.Text)

	// Named replacement tokens
	for i := range this.Params {
		if len(params) <= i {
			out = strings.Replace(out, "%{"+this.Params[i].Name+"}", this.Params[i].Default, -1)
		} else {
			out = strings.Replace(out, "%{"+this.Params[i].Name+"}", params[i], -1)
		}
	}
	out = repTokenNameRegEx.ReplaceAllStringFunc(out, func(in string) string {
		in = strings.TrimLeft(in, "%")
		for i := range this.Params {
			if this.Params[i].Name == in {
				if len(params) > i {
					return params[i]
				}
				return this.Params[i].Default
			}
		}
		return "%" + in
	})

	// Numbered replacement tokens
	out = repTokenNumberRegEx.ReplaceAllStringFunc(out, func(in string) string {
		in = strings.TrimLeft(in, "%")
		i, err := strconv.ParseInt(in, 10, 8)
		if err != nil {
			panic(err)
		}
		if len(params) >= int(i) {
			return params[i-1]
		}
		if len(this.Params) >= int(i) {
			return this.Params[i-1].Default
		}
		return "%" + in
	})

	return string(Parse([]byte(out), stgUseCurrent))
}

func NewNativeTemplate(name string, handler NativeTemplate) {
	rtn := new(Template)
	rtn.Native = true
	rtn.Handler = handler
	Templates[name] = rtn
}

func NewUserTemplate(name string, text string, params []*TemplateParam) {
	rtn := new(Template)
	rtn.Text = text
	rtn.Params = params
	Templates[name] = rtn
}

func NewScriptTemplate(name string, text string, params []*TemplateParam) {
	rtn := new(Template)
	rtn.NCA = true
	rtn.Text = text
	rtn.Params = params
	Templates[name] = rtn
}
