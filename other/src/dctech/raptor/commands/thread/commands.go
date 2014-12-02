/*
Copyright 2012-2014 by Milo Christiansen

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

// Raptor Threading Commands.
package thread

import "dctech/raptor"

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
func Setup(state *raptor.State) {
	state.NewNameSpace("thread")
	state.NewNativeCommand("thread:new", CommandThread_New)
	state.NewNativeCommand("thread:channel", CommandThread_Channel)
	state.NewNativeCommand("thread:send", CommandThread_Send)
	state.NewNativeCommand("thread:receive", CommandThread_Receive)
	state.NewNativeCommand("thread:close", CommandThread_Close)
}

// Creates a new thread.
// 	thread:new code [params...]
// WARNING! The base types map and array are NOT thread safe!
// Returns unchanged.
func CommandThread_New(script *raptor.Script, params []*raptor.Value) {
	if len(params) < 1 {
		panic(script.BadParamCount(">=1"))
	}

	scr := raptor.NewScript()
	scr.Envs.Add(raptor.NewEnvironment())

	scr.AddParamsValue(params[1:]...)

	scr.Code.Add(params[0].CodeSource())
	go func() {
		_, err := script.Host.Run(scr)
		if err != nil {
			panic(err)
		}
	}()

	return
}

// Creates a new communications channel.
// 	thread:channel [buffer]
// Returns the new channel.
func CommandThread_Channel(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 0 && len(params) != 1 {
		panic(script.BadParamCount("0 or 1"))
	}

	if len(params) == 0 {
		script.RetVal = raptor.NewValueObject(make(chan *raptor.Value))
		return
	}

	script.RetVal = raptor.NewValueObject(make(chan *raptor.Value, params[0].Int64()))
	return
}

// Send a value on a communications channel.
// 	thread:send channel value
// Sending to a closed channel is a runtime error.
// Returns unchanged.
func CommandThread_Send(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic(script.BadParamCount("2"))
	}

	if _, ok := params[0].Data.(chan *raptor.Value); !ok {
		panic(script.GeneralCmdError("Parameter 0 is not a channel."))
	}
	channel := params[0].Data.(chan *raptor.Value)
	channel <- params[1]

	return
}

// Receive a value from a communications channel.
// 	thread:receive channel
// If the channel is closed, always returns nil
// Returns the received value.
func CommandThread_Receive(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	if _, ok := params[0].Data.(chan *raptor.Value); !ok {
		panic(script.GeneralCmdError("Parameter 0 is not a channel."))
	}
	channel := params[0].Data.(chan *raptor.Value)
	script.RetVal = <-channel

	return
}

// Close a communications channel.
// 	thread:close channel
// Returns unchanged.
func CommandThread_Close(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	if _, ok := params[0].Data.(chan *raptor.Value); !ok {
		panic(script.GeneralCmdError("Parameter 0 is not a channel."))
	}
	channel := params[0].Data.(chan *raptor.Value)
	close(channel)

	return
}
