package nca4

// NativeCommand is the function signature for a native command handler.
type NativeCommand func(*State, []*Value)

// A Command stores a native or user script command. Used by State and NameSpace.
type Command struct {
	// Is this a native command?
	Native bool

	// The native command handler
	Handler NativeCommand

	// The user command code,
	Code *Value

	// Does this user command have a variable number of params?
	VarParams bool

	// Param names for user commands
	Params []string
}

// Call runs a command.
func (this *Command) Call(state *State, params []*Value) {
	// Native command
	if this.Native {
		this.Handler(state, params)
		return
	}

	// User command
	if !this.VarParams {
		if len(this.Params) != len(params) {
			panic("Invalid param count to user command.")
		}
	}
	
	state.Envs.Add(NewEnvironment())

	if !this.VarParams {
		for i := range this.Params {
			state.NewVar(this.Params[i], params[i])
		}
	} else {
		state.AddParamsValue(params...)
	}

	state.Code.AddLexer(this.Code.AsLexer())
	state.Exec()
	state.Envs.Remove()
	state.Return = false
	return
}
