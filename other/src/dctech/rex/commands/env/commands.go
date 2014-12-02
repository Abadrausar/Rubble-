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

// Rex OS Environment Commands.
package env

import "dctech/rex"
import "os"
import "os/exec"
import "syscall"
import "path/filepath"

// Adds the OS environment commands to the state.
// The OS environment commands are:
//	env:exec
//	env:getvar
//	env:setvar
//	env:getwd
//	env:setwd
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
	
	mod := state.RegisterModule("env")
	mod.RegisterCommand("exec", Command_Exec)
	mod.RegisterCommand("getvar", Command_GetVar)
	mod.RegisterCommand("setvar", Command_SetVar)
	mod.RegisterCommand("getwd", Command_GetWD)
	mod.RegisterCommand("setwd", Command_SetWD)
	
	return nil
}

// Runs an external program and waits for it to exit.
// 	env:exec path [args]
// Returns unchanged.
func Command_Exec(script *rex.Script, params []*rex.Value) {
	if len(params) < 1 {
		rex.ErrorParamCount("env:exec", ">1")
	}

	strparams := make([]string, len(params)-1)
	for i, val := range params[1:] {
		strparams[i] = val.String()
	}

	path, err := filepath.Abs(params[0].String())
	if err != nil {
		rex.ErrorGeneralCmd("env:exec", err.Error())
	}
	//path = "\"" + path + "\""
	
	cmd := exec.Command(path, strparams...)

	if script.Output == nil {
		if script.Host == nil {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		} else {
			cmd.Stdout = script.Host.Output
			cmd.Stderr = script.Host.Output
		}
	} else {
		cmd.Stdout = script.Output
		cmd.Stderr = script.Output
	}

	err = cmd.Start()
	if err != nil {
		rex.ErrorGeneralCmd("env:exec", err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		rex.ErrorGeneralCmd("env:exec", err.Error())
	}
}

// Gets the value of an environment variable.
// 	env:getvar name
// Returns the variables value.
func Command_GetVar(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("env:getvar", "1")
	}

	script.RetVal = rex.NewValueString(os.Getenv(params[0].String()))
}

// Sets the value of an environment variable.
// 	env:setvar name value
// Returns the old value.
func Command_SetVar(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("env:setvar", "2")
	}
	rtn := rex.NewValueString(os.Getenv(params[0].String()))
	err := os.Setenv(params[0].String(), params[1].String())
	if err != nil {
		rex.ErrorGeneralCmd("env:setvar", err.Error())
	}
	script.RetVal = rtn
}

// Get the working directory.
// 	env:getwd
// Returns the working directory.
func Command_GetWD(script *rex.Script, params []*rex.Value) {
	wd, err := os.Getwd()
	if err != nil {
		rex.ErrorGeneralCmd("env:getwd", err.Error())
	}

	script.RetVal = rex.NewValueString(wd)
}

// Set the working directory.
// 	env:setwd path
// Returns unchanged.
func Command_SetWD(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("env:setwd", "1")
	}

	err := syscall.Chdir(params[0].String())
	if err != nil {
		rex.ErrorGeneralCmd("env:setwd", err.Error())
	}
}
