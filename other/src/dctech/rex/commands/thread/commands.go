/*
Copyright 2014 by Milo Christiansen

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

// Rex Threading Commands.
package thread

import "dctech/rex"

// TODO:
//	Come up with a way to simulate the behavior of the select statement
//

// Adds the threading commands to the state.
// The threading commands are:
//	thread:new
//	thread:channel
//	thread:send
//	thread:receive
//	thread:close
func Setup(state *rex.State) (err error) {
	defer func() {
		if !state.NoRecover {
			if x := recover(); x != nil {
				if y, ok := x.(rex.ScriptError); ok {
					err = y
					return
				}
				panic(x)
			}
		}
	}()
	
	mod := state.RegisterModule("thread")
	mod.RegisterCommand("new", Command_New)
	mod.RegisterCommand("channel", Command_Channel)
	mod.RegisterCommand("send", Command_Send)
	mod.RegisterCommand("receive", Command_Receive)
	mod.RegisterCommand("close", Command_Close)
	
	return nil
}

// Creates a new thread.
// 	thread:new [params...] code
// code MUST be a block declaration!
// The params for code must match the params you are passing in!
// Each thread is run in it's own script, so it is not possible to access "up stream" variables.
// Returns unchanged.
func Command_New(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("thread:new", "1")
	}

	if params[1].Type != rex.TypCode {
		rex.ErrorGeneralCmd("thread:new", "Attempt to run non-executable Value.")
	}
	block := params[0].Data.(*rex.Code)
	
	scr := rex.NewScript()
	script.Locals.Add(script.Host, block)
	
	script.SetParams(block, params[1:])
	
	go func() {
		_, err := script.Host.Run(scr, params[0])
		if err != nil {
			panic(err)
		}
	}()
}

// Creates a new communications channel.
// 	thread:channel [buffer]
// Returns the new channel.
func Command_Channel(script *rex.Script, params []*rex.Value) {
	if len(params) != 0 && len(params) != 1 {
		rex.ErrorParamCount("thread:channel", "0 or 1")
	}

	if len(params) == 0 {
		script.RetVal = rex.NewValueUser(make(chan *rex.Value))
		return
	}
	script.RetVal = rex.NewValueUser(make(chan *rex.Value, params[0].Int64()))
}

// Send a value on a communications channel.
// 	thread:send channel value
// Sending to a closed channel is a runtime error.
// Returns unchanged.
func Command_Send(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("thread:send", "2")
	}

	if _, ok := params[0].Data.(chan *rex.Value); !ok {
		rex.ErrorGeneralCmd("thread:send", "Parameter 0 is not a channel.")
	}
	channel := params[0].Data.(chan *rex.Value)
	channel <- params[1]
}

// Receive a value from a communications channel.
// 	thread:receive channel
// If the channel is closed, always returns nil
// Returns the received value.
func Command_Receive(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("thread:receive", "1")
	}

	if _, ok := params[0].Data.(chan *rex.Value); !ok {
		rex.ErrorGeneralCmd("thread:receive", "Parameter 0 is not a channel.")
	}
	channel := params[0].Data.(chan *rex.Value)
	script.RetVal = <-channel
	if script.RetVal == nil {
		script.RetVal = rex.NewValue()
	}
}

// Close a communications channel.
// 	thread:close channel
// Returns unchanged.
func Command_Close(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("thread:close", "1")
	}

	if _, ok := params[0].Data.(chan *rex.Value); !ok {
		rex.ErrorGeneralCmd("thread:close", "Parameter 0 is not a channel.")
	}
	channel := params[0].Data.(chan *rex.Value)
	close(channel)
}
