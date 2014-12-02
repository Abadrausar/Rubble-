/*
For copyright/license see header in file "doc.go"
*/

package raptor

// NativeCommand is the function signature for a native command handler.
type NativeCommand func(*Script, []*Value)

// A Command stores a native or user script command. Used by State and NameSpace.
type Command struct {
	// The native command handler, if non-nil all other fields are ignored.
	Handler NativeCommand

	// The user command code.
	Code *Code

	// Parameter names, if this is nil the command has a variable parameter count.
	Params []string
}

// Call runs a command.
func (this *Command) Call(name string, script *Script, params []*Value) {
	// Native command
	if this.Handler != nil {
		script.RunningCmd = name
		this.Handler(script, params)
		script.RunningCmd = ""
		return
	}

	// User command
	if this.Params != nil {
		if len(this.Params) != len(params) {
			panic("Invalid param count to user command:" + name + ".")
		}
	}

	script.Envs.Add(NewEnvironment())

	if this.Params != nil {
		for i := range this.Params {
			script.NewVar(this.Params[i], params[i])
		}
	} else {
		script.AddParamsValue(params...)
	}

	script.Code.Add(NewCodeReader(this.Code))
	script.Exec()
	script.Envs.Remove()
	script.Return = false
	return
}
