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

package rex

// ScriptError is used for any and every error that is caused by Rex.
type ScriptError struct {
	Msg string
	Pos *Position
}

func (err ScriptError) Error() string {
	if err.Pos != nil {
		return err.Msg + " Near: " + err.Pos.String()
	}
	return err.Msg
}

// RaiseError converts a string to a ScriptError and then panics with it.
// The error's position is filled out when it is caught by the state.
func RaiseError(msg string) {
	panic(ScriptError{Msg: msg})
}

// ErrorGeneralCmd generates an error message for general command errors.
//	"test:cmd": some error message
// The generated message is passed directly to RaiseError.
func ErrorGeneralCmd(name, msg string) {
	RaiseError("\"" + name + "\": " + msg)
}

// ErrorParamCount generates an error message for incorrect parameter count errors.
//	"test:cmd": Incorrect parameter count: 1 parameter(s) required.
// If hint is set to "" it leaves that part off.
//	"test:cmd": Incorrect parameter count.
// Hint should be something like "5" or ">2".
// The generated message is passed directly to RaiseError.
func ErrorParamCount(name, hint string) {
	if hint != "" {
		ErrorGeneralCmd(name, "Incorrect parameter count: " + hint + " parameter(s) required.")
	}
	ErrorGeneralCmd(name, "Incorrect parameter count.")
}


// InternalError is used for any and every error that is caused by something other than Rex.
// Examine the Err field to get the original error.
type InternalError struct {
	Err error
	Pos *Position
}

func (err InternalError) Error() string {
	if err.Pos != nil {
		return err.Err.Error() + " Near: " + err.Pos.String()
	}
	return err.Err.Error()
}

// RaiseInternalError converts a string to a ScriptError and then panics with it.
// The error's position is filled out when it is caught by the state.
func RaiseInternalError(err error) {
	panic(InternalError{Err: err})
}
