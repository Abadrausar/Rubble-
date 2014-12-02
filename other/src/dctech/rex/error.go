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

type ErrType int

// Error types, most errors use type ErrTypUndefined so far.
const (
	ErrTypUndefined ErrType = iota // Most errors use this type for now.
	ErrTypGenCommand               // Used for any error created by ErrorGeneralCmd.
	ErrTypParamCount               // Used for any error created by ErrorParamCount.
	ErrTypGenCompile               // Unused (for now...)
	ErrTypGenRuntime               // Unused (for now...)
)

// ScriptError is used for any and every error that is caused by Rex.
type ScriptError struct {
	Msg string
	Pos *Position
	
	// Error type, eventually this will be usable for advanced error checking.
	Type ErrType
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

// RaiseError converts a string to a ScriptError and then panics with it.
// The error's position is filled out when it is caught by the state.
func RaiseErrorType(msg string, typ ErrType) {
	panic(ScriptError{Msg: msg, Type: typ})
}

// ErrorGeneralCmd generates an error message for general command errors.
//	"test:cmd": some error message
// The generated message is passed directly to RaiseErrorType.
func ErrorGeneralCmd(name, msg string) {
	RaiseErrorType("\"" + name + "\": " + msg, ErrTypGenCommand)
}

// ErrorParamCount generates an error message for incorrect parameter count errors.
//	"test:cmd": Incorrect parameter count: 1 parameter(s) required.
// If hint is set to "" it leaves that part off.
//	"test:cmd": Incorrect parameter count.
// Hint should be something like "5" or ">2".
// The generated message is passed directly to RaiseErrorType.
func ErrorParamCount(name, hint string) {
	if hint != "" {
		RaiseErrorType("\"" + name + "\": " + "Incorrect parameter count: " + hint + " parameter(s) required.", ErrTypParamCount)
	}
	RaiseErrorType("\"" + name + "\": " + "Incorrect parameter count.", ErrTypParamCount)
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

// RaiseInternalError converts an error to a InternalError and then panics with it.
// The error's position is filled out when it is caught by the state.
func RaiseInternalError(err error) {
	panic(InternalError{Err: err})
}
