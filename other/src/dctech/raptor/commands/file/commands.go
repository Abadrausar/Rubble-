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

// Raptor File System Commands.
package file

import "dctech/raptor"
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
func Setup(state *raptor.State) {
	state.NewNameSpace("file")
	state.NewNativeCommand("file:deldir", CommandFile_DelDir)
	state.NewNativeCommand("file:deltree", CommandFile_DelTree)
	state.NewNativeCommand("file:newdir", CommandFile_NewDir)
	state.NewNativeCommand("file:del", CommandFile_Del)
	state.NewNativeCommand("file:exists", CommandFile_Exists)
	state.NewNativeCommand("file:direxists", CommandFile_DirExists)
	state.NewNativeCommand("file:walkfiles", CommandFile_WalkFiles)
	state.NewNativeCommand("file:walkdirs", CommandFile_WalkDirs)
}

// Delete an empty directory.
// 	file:deldir path
// Returns unchanged or an error message. May set the Error flag.
func CommandFile_DelDir(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
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
func CommandFile_DelTree(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	info, err := os.Lstat(params[0].String())
	if err != nil {
		script.RetVal = raptor.NewValueString(err.Error())
		script.Error = true
		return
	}
	if info.IsDir() {
		err := os.RemoveAll(params[0].String())
		if err != nil {
			script.RetVal = raptor.NewValueString(err.Error())
			script.Error = true
		}
		return
	}
	script.RetVal = raptor.NewValueString("File not a directory.")
	script.Error = true
	return
}

// Create a new directory.
// 	file:newdir path
// Returns unchanged or an error message. May set the Error flag.
func CommandFile_NewDir(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	err := os.Mkdir(params[0].String(), 0600)
	if err != nil {
		script.RetVal = raptor.NewValueString(err.Error())
		script.Error = true
	}
	return
}

// Delete a file.
// 	file:del path
// Returns unchanged or an error message. May set the Error flag.
func CommandFile_Del(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	info, err := os.Lstat(params[0].String())
	if err != nil {
		script.RetVal = raptor.NewValueString(err.Error())
		script.Error = true
		return
	}
	if !info.IsDir() {
		err := os.Remove(params[0].String())
		if err != nil {
			script.RetVal = raptor.NewValueString(err.Error())
			script.Error = true
		}
		return
	}
	script.RetVal = raptor.NewValueString("File is a directory.")
	script.Error = true
	return
}

// Checks if a file exists.
// 	file:exists path
// Returns true or false.
func CommandFile_Exists(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	info, err := os.Lstat(params[0].String())
	if err != nil {
		script.RetVal = raptor.NewValueBool(false)
		return
	}
	if info.IsDir() {
		script.RetVal = raptor.NewValueBool(false)
		return
	}
	script.RetVal = raptor.NewValueBool(true)
}

// Checks if a directory exists.
// 	file:direxists path
// Returns true or false.
func CommandFile_DirExists(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	info, err := os.Lstat(params[0].String())
	if err != nil {
		script.RetVal = raptor.NewValueBool(false)
		return
	}
	if !info.IsDir() {
		script.RetVal = raptor.NewValueBool(false)
		return
	}
	script.RetVal = raptor.NewValueBool(true)
	return
}

// Iterate over all the files in a directory.
// 	file:walkfiles path code
// Calles code (as a command) for every file found:
//	code path
// Returns nil or an error message. May set the Error flag.
func CommandFile_WalkFiles(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic(script.BadParamCount("2"))
	}

	files, err := ioutil.ReadDir(params[0].String())
	if err != nil {
		script.RetVal = raptor.NewValueString(err.Error())
		script.Error = true
		return
	}

	code := params[1].Code()
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		script.Envs.Add(raptor.NewEnvironment())

		script.AddParams(file.Name())

		script.Code.Add(raptor.NewCodeReader(code))
		script.Exec()
		script.Envs.Remove()
		script.Return = false
	}
	script.RetVal = nil
}

// Iterate over all the directories in a directory.
// 	file:walkdirs path code
// Calles code (as a command) for every directory found:
//	code path
// Returns nil or an error message. May set the Error flag.
func CommandFile_WalkDirs(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic(script.BadParamCount("2"))
	}

	files, err := ioutil.ReadDir(params[0].String())
	if err != nil {
		script.RetVal = raptor.NewValueString(err.Error())
		script.Error = true
		return
	}

	code := params[1].Code()
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		script.Envs.Add(raptor.NewEnvironment())

		script.AddParams(file.Name())

		script.Code.Add(raptor.NewCodeReader(code))
		script.Exec()
		script.Envs.Remove()
		script.Return = false
	}
	script.RetVal = nil
}
