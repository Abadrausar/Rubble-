package main

import "strings"
import "regexp"
import "strconv"

type NativeTemplate func([]string) string

var Templates = make(map[string]*Template)

type TemplateParam struct {
	Name string
	Default string
}

type Template struct {
	// Is this a native template?
	Native bool

	// The native template handler
	Handler NativeTemplate

	// The user template text
	Text string

	// Param names/defaults for user commands
	Params []*TemplateParam
}

var repTokenNameRegEx = regexp.MustCompile("%[a-zA-Z_]+")
var repTokenNumberRegEx = regexp.MustCompile("%[1-9]+")

// Call runs a template. The output is run through the parser for the active stage.
func (this *Template) Call(params []string) string {
	// prep params
	for i := range params {
		params[i] = strings.TrimSpace(params[i])
	}
	
	// Native template
	if this.Native {
		return this.Handler(params)
	}

	// User template
	out := this.Text
	for i := range this.Params {
		if len(params) <= i {
			out = strings.Replace(out, "%{" + this.Params[i].Name + "}", this.Params[i].Default, -1)
		} else {
			out = strings.Replace(out, "%{" + this.Params[i].Name + "}", params[i], -1)
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
	
	return StageParse(out)
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

