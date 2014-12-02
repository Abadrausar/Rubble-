package main

import "strings"

type NativeTemplate func([]string) string

var Templates = make(map[string]*Template)

type Template struct {
	// Is this a native template?
	Native bool

	// The native template handler
	Handler NativeTemplate

	// The user template text
	Text string

	// Param names for user commands
	Params []string
}

// Call runs a template.
// It is the job of the caller to parse the output for possible embeded template calls.
func (this *Template) Call(params []string) string {
	// Native template
	if this.Native {
		return this.Handler(params)
	}

	// User template
	if len(this.Params) != len(params) {
		panic("Invalid param count to user template.")
	}
	
	out := this.Text
	for i := range this.Params {
		out = strings.Replace(out, "%(" + this.Params[i] + ")", params[i], -1)
	}
	
	return StageParse(out)
}

func NewNativeTemplate(name string, handler NativeTemplate) {
	rtn := new(Template)
	rtn.Native = true
	rtn.Handler = handler
	Templates[name] = rtn
}

func NewUserTemplate(name string, text string, params []string) {
	rtn := new(Template)
	rtn.Text = text
	rtn.Params = params
	Templates[name] = rtn
}

