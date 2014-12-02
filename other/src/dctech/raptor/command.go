/*
For copyright/license see header in file "doc.go"
*/

package raptor

// NativeCommand is the function signature for a native command handler.
type NativeCommand func(*Script, []*Value)

// TODO: Add support for default params.

// A Command stores a native or user script command. Used by State and NameSpace.
type Command struct {
	// Is this a native command?
	Native bool

	// The native command handler
	Handler NativeCommand

	// The user command code,
	Code *Code

	// Does this user command have a variable number of params?
	VarParams bool

	// Param names for user commands
	Params []string
}

// Call runs a command.
func (this *Command) Call(script *Script, params []*Value) {
	// Native command
	if this.Native {
		this.Handler(script, params)
		return
	}

	// User command
	if !this.VarParams {
		if len(this.Params) != len(params) {
			panic("Invalid param count to user command.")
		}
	}

	script.Envs.Add(NewEnvironment())

	if !this.VarParams {
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
