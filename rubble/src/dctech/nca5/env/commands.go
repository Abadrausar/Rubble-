// NCA v5 OS Environment Commands.
package env

import "dctech/nca5"
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
func Setup(state *nca5.State) {
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
func CommandEnv_Exec(state *nca5.State, params []*nca5.Value) {
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
func CommandEnv_GetVar(state *nca5.State, params []*nca5.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to getenv.")
	}

	state.RetVal = nca5.NewValue(os.Getenv(params[0].String()))
}

// Sets the value of an environment variable.
// 	env:setvar name value
// Returns the old value.
func CommandEnv_SetVar(state *nca5.State, params []*nca5.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to setenv.")
	}
	rtn := nca5.NewValue(os.Getenv(params[0].String()))
	err := os.Setenv(params[0].String(), params[1].String())
	if err != nil {
		panic(err)
	}
	state.RetVal = rtn
}

// Get the working dir.
// 	env:getwd
// Returns the working directory.
func CommandEnv_GetWD(state *nca5.State, params []*nca5.Value) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	state.RetVal = nca5.NewValue(wd)
}

// Set the working dir.
// 	env:setwd path
// Returns "".
func CommandEnv_SetWD(state *nca5.State, params []*nca5.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to env:setwd.")
	}

	err := syscall.Chdir(params[0].String())
	if err != nil {
		panic(err)
	}
}
