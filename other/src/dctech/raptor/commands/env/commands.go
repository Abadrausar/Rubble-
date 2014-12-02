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

// Raptor OS Environment Commands.
package env

import "dctech/raptor"
import "os"
import "os/exec"
import "syscall"

// Adds the OS environment commands to the state.
// The OS environment commands are:
//	env:exec
//	env:getvar
//	env:setvar
//	env:getwd
//	env:setwd
func Setup(state *raptor.State) {
	state.NewNameSpace("env")
	state.NewNativeCommand("env:exec", CommandEnv_Exec)
	state.NewNativeCommand("env:getvar", CommandEnv_GetVar)
	state.NewNativeCommand("env:setvar", CommandEnv_SetVar)
	state.NewNativeCommand("env:getwd", CommandEnv_GetWD)
	state.NewNativeCommand("env:setwd", CommandEnv_SetWD)
}

// Runs an external program and waits for it to exit.
// 	env:exec path [args]
// Returns unchanged.
func CommandEnv_Exec(script *raptor.Script, params []*raptor.Value) {
	if len(params) < 1 {
		panic("Wrong number of params to env:exec.")
	}

	strparams := make([]string, len(params)-1)
	for i, val := range params[1:] {
		strparams[i] = val.String()
	}

	cmd := exec.Command(params[0].String(), strparams...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

// Gets the value of an environment variable.
// 	env:getvar name
// Returns the variables value.
func CommandEnv_GetVar(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to getenv.")
	}

	script.RetVal = raptor.NewValueString(os.Getenv(params[0].String()))
}

// Sets the value of an environment variable.
// 	env:setvar name value
// Returns the old value.
func CommandEnv_SetVar(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to setenv.")
	}
	rtn := raptor.NewValueString(os.Getenv(params[0].String()))
	err := os.Setenv(params[0].String(), params[1].String())
	if err != nil {
		panic(err)
	}
	script.RetVal = rtn
}

// Get the working dir.
// 	env:getwd
// Returns the working directory.
func CommandEnv_GetWD(script *raptor.Script, params []*raptor.Value) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	script.RetVal = raptor.NewValueString(wd)
}

// Set the working dir.
// 	env:setwd path
// Returns unchanged.
func CommandEnv_SetWD(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to env:setwd.")
	}

	err := syscall.Chdir(params[0].String())
	if err != nil {
		panic(err)
	}
}
