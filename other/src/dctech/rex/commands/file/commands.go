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

// Rex File System Commands.
package file

import "dctech/rex"
import "io/ioutil"
import "os"

// Adds the file system commands to the state.
// The file system commands are:
//	file:deldir
//	file:deltree
//	file:newdir
//	file:del
//	file:exists
//	file:direxists
//	file:walkfiles
//	file:walkdirs
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
	
	mod := state.RegisterModule("file")
	mod.RegisterCommand("deldir", Command_DelDir)
	mod.RegisterCommand("deltree", Command_DelTree)
	mod.RegisterCommand("newdir", Command_NewDir)
	mod.RegisterCommand("del", Command_Del)
	mod.RegisterCommand("exists", Command_Exists)
	mod.RegisterCommand("direxists", Command_DirExists)
	mod.RegisterCommand("walkfiles", Command_WalkFiles)
	mod.RegisterCommand("walkdirs", Command_WalkDirs)
	
	return nil
}

// Delete an empty directory.
// 	file:deldir path
// Returns unchanged or an error message. May set the Error flag.
func Command_DelDir(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("file:deldir", "1")
	}

	info, err := os.Lstat(params[0].String())
	if err != nil {
		script.Error = true
		return
	}
	if info.IsDir() {
		err := os.Remove(params[0].String())
		if err != nil {
			script.Error = true
		}
		return
	}
	script.Error = true
	return
}

// Delete a directory and all sub dirs and files.
// 	file:deltree path
// Returns unchanged or an error message. May set the Error flag.
func Command_DelTree(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("file:deltree", "1")
	}

	info, err := os.Lstat(params[0].String())
	if err != nil {
		script.RetVal = rex.NewValueString(err.Error())
		script.Error = true
		return
	}
	if info.IsDir() {
		err := os.RemoveAll(params[0].String())
		if err != nil {
			script.RetVal = rex.NewValueString(err.Error())
			script.Error = true
		}
		return
	}
	script.RetVal = rex.NewValueString("File not a directory.")
	script.Error = true
	return
}

// Create a new directory.
// 	file:newdir path
// Returns unchanged or an error message. May set the Error flag.
func Command_NewDir(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("file:newdir", "1")
	}

	err := os.Mkdir(params[0].String(), 0600)
	if err != nil {
		script.RetVal = rex.NewValueString(err.Error())
		script.Error = true
	}
	return
}

// Delete a file.
// 	file:del path
// Returns unchanged or an error message. May set the Error flag.
func Command_Del(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("file:del", "1")
	}

	info, err := os.Lstat(params[0].String())
	if err != nil {
		script.RetVal = rex.NewValueString(err.Error())
		script.Error = true
		return
	}
	if !info.IsDir() {
		err := os.Remove(params[0].String())
		if err != nil {
			script.RetVal = rex.NewValueString(err.Error())
			script.Error = true
		}
		return
	}
	script.RetVal = rex.NewValueString("File is a directory.")
	script.Error = true
	return
}

// Checks if a file exists.
// 	file:exists path
// Returns true or false.
func Command_Exists(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("file:exists", "1")
	}

	info, err := os.Lstat(params[0].String())
	if err != nil {
		script.RetVal = rex.NewValueBool(false)
		return
	}
	if info.IsDir() {
		script.RetVal = rex.NewValueBool(false)
		return
	}
	script.RetVal = rex.NewValueBool(true)
}

// Checks if a directory exists.
// 	file:direxists path
// Returns true or false.
func Command_DirExists(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("file:direxists", "1")
	}

	info, err := os.Lstat(params[0].String())
	if err != nil {
		script.RetVal = rex.NewValueBool(false)
		return
	}
	if !info.IsDir() {
		script.RetVal = rex.NewValueBool(false)
		return
	}
	script.RetVal = rex.NewValueBool(true)
	return
}

// Iterate over all the files in a directory.
// 	file:walkfiles path code
// Runs code for every file found, params:
//	path
// code MUST be a block created via a block declaration!
// Returns nil or an error message. May set the Error flag.
func Command_WalkFiles(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("file:walkfiles", "2")
	}

	files, err := ioutil.ReadDir(params[0].String())
	if err != nil {
		script.RetVal = rex.NewValueString(err.Error())
		script.Error = true
		return
	}

	if params[1].Type != rex.TypCode {
		rex.ErrorGeneralCmd("file:walkfiles", "Attempt to run non-executable Value.")
	}
	block := params[1].Data.(*rex.Code)
	
	script.Locals.Add(block)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		script.Locals.Set(0, rex.NewValueString(file.Name()))
		script.Exec(block)
		script.Return = false
	}
	script.Locals.Remove()
	script.RetVal = rex.NewValue()
}

// Iterate over all the directories in a directory.
// 	file:walkdirs path code
// Runs code for every directory found, params:
//	path
// code MUST be a block created via a block declaration!
// Returns nil or an error message. May set the Error flag.
func Command_WalkDirs(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("file:walkdirs", "2")
	}

	files, err := ioutil.ReadDir(params[0].String())
	if err != nil {
		script.RetVal = rex.NewValueString(err.Error())
		script.Error = true
		return
	}
	
	if params[1].Type != rex.TypCode {
		rex.ErrorGeneralCmd("file:walkdirs", "Attempt to run non-executable Value.")
	}
	block := params[1].Data.(*rex.Code)
	
	script.Locals.Add(block)
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		script.Locals.Set(0, rex.NewValueString(file.Name()))
		script.Exec(block)
		script.Return = false
	}
	script.Locals.Remove()
	script.RetVal = rex.NewValue()
}
