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

// Rex File IO Commands.
package fileio

import "dctech/rex"
import "io/ioutil"

// Adds the file io commands to the state.
// The file io commands are:
//	fileio:read
//	fileio:write
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
	
	mod := state.RegisterModule("fileio")
	mod.RegisterCommand("read", Command_Read)
	mod.RegisterCommand("write", Command_Write)
	
	return nil
}

// Read from file.
// 	fileio:read path
// Returns file contents or an error message. May set the Error flag.
func Command_Read(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("fileio:read", "1")
	}

	file, err := ioutil.ReadFile(params[0].String())
	if err != nil {
		script.Error = true
		script.RetVal = rex.NewValueString("error:" + err.Error())
		return
	}
	script.RetVal = rex.NewValueString(string(file))
}

// Write to file.
// 	fileio:write path contents
// Returns unchanged or an error message. May set the Error flag.
func Command_Write(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("fileio:write", "2")
	}

	err := ioutil.WriteFile(params[0].String(), []byte(params[1].String()), 0666)
	if err != nil {
		script.Error = true
		script.RetVal = rex.NewValueString("error:" + err.Error())
		return
	}
}
