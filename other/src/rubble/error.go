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

package rubble

import "dctech/rex"

// Abort is used for errors that should not have attached error position information.
type Abort string

func (err Abort) Error() string {
	return string(err)
}

func RaiseAbort(msg string) {
	panic(Abort(msg))
}

// RblError is used for most common errors caused by Rubble.
type RblError struct {
	Msg string
	Pos *rex.Position
}

func (err RblError) Error() string {
	if err.Pos != nil {
		return err.Msg + " Near: " + err.Pos.String()
	}
	return err.Msg
}

// RaiseError converts a string to a RblError and then panics with it.
// The error's position is filled out when it is caught.
func RaiseError(msg string) {
	panic(RblError{Msg: msg})
}


// InternalError is used for any and every error that is caused by something other than Rubble.
// Examine the Err field to get the original error.
type InternalError struct {
	Err error
	Pos *rex.Position
}

func (err InternalError) Error() string {
	if err.Pos != nil {
		return err.Err.Error() + " Near: " + err.Pos.String()
	}
	return err.Err.Error()
}

// RaiseInternalError converts a string to an InternalError and then panics with it.
// The error's position is filled out when it is caught.
func RaiseInternalError(err error) {
	panic(InternalError{Err: err})
}
