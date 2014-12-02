/*
Copyright 2012-2013 by Milo Christiansen

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

package nca7

// NativeCommand is the function signature for a native command handler.
type NativeCommand func(*State, []*Value)

// A Command stores a native or user script command. Used by State and NameSpace.
type Command struct {
	// Is this a native command?
	Native bool

	// The native command handler
	Handler NativeCommand

	// The user command code,
	Code *CompiledScript

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

	state.Code.AddCodeSource(NewCompiledLexer(this.Code))
	state.Exec()
	state.Envs.Remove()
	state.Return = false
	return
}
