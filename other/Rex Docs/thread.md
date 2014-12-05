PACKAGE DOCUMENTATION

package thread
    import "dctech/rex/commands/thread"



FUNCTIONS

func Command_Channel(script *rex.Script, params []*rex.Value)
    Creates a new communications channel.

	thread:channel [buffer]

    Returns the new channel.

func Command_Close(script *rex.Script, params []*rex.Value)
    Close a communications channel.

	thread:close channel

    Returns unchanged.

func Command_New(script *rex.Script, params []*rex.Value)
    Creates a new thread.

	thread:new [params...] code

    code MUST be a block declaration! The params for code must match the
    params you are passing in! Each thread is run in it's own script, so it
    is not possible to access "up stream" variables. Returns unchanged.

func Command_Receive(script *rex.Script, params []*rex.Value)
    Receive a value from a communications channel.

	thread:receive channel

    If the channel is closed, always returns nil Returns the received value.

func Command_Send(script *rex.Script, params []*rex.Value)
    Send a value on a communications channel.

	thread:send channel value

    Sending to a closed channel is a runtime error. Returns unchanged.


