/*
Copyright 2012-2013 by Milo Christiansen

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
func CommandFile_DelDir(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:deldir.")
	}

	info, err := os.Lstat(params[0].String())
	if err != nil {
		state.Error = true
		return
	}
	if info.IsDir() {
		err := os.Remove(params[0].String())
		if err != nil {
			state.Error = true
		}
		return
	}
	state.Error = true
	return
}

// Delete a directory and all sub dirs and files.
// 	file:deltree path
// Returns unchanged or an error message. May set the Error flag.
func CommandFile_DelTree(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:deltree.")
	}

	info, err := os.Lstat(params[0].String())
	if err != nil {
		state.RetVal = raptor.NewValueString(err.Error())
		state.Error = true
		return
	}
	if info.IsDir() {
		err := os.RemoveAll(params[0].String())
		if err != nil {
			state.RetVal = raptor.NewValueString(err.Error())
			state.Error = true
		}
		return
	}
	state.RetVal = raptor.NewValueString("File not a directory.")
	state.Error = true
	return
}

// Create a new directory.
// 	file:newdir path
// Returns unchanged or an error message. May set the Error flag.
func CommandFile_NewDir(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:newdir.")
	}

	err := os.Mkdir(params[0].String(), 0600)
	if err != nil {
		state.RetVal = raptor.NewValueString(err.Error())
		state.Error = true
	}
	return
}

// Delete a file.
// 	file:del path
// Returns unchanged or an error message. May set the Error flag.
func CommandFile_Del(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:del.")
	}

	info, err := os.Lstat(params[0].String())
	if err != nil {
		state.RetVal = raptor.NewValueString(err.Error())
		state.Error = true
		return
	}
	if !info.IsDir() {
		err := os.Remove(params[0].String())
		if err != nil {
			state.RetVal = raptor.NewValueString(err.Error())
			state.Error = true
		}
		return
	}
	state.RetVal = raptor.NewValueString("File is a directory.")
	state.Error = true
	return
}

// Checks if a file exists.
// 	file:exists path
// Returns true or false.
func CommandFile_Exists(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:exists.")
	}

	info, err := os.Lstat(params[0].String())
	if err != nil {
		state.RetVal = raptor.NewValueBool(false)
		return
	}
	if info.IsDir() {
		state.RetVal = raptor.NewValueBool(false)
		return
	}
	state.RetVal = raptor.NewValueBool(true)
}

// Checks if a directory exists.
// 	file:direxists path
// Returns true or false.
func CommandFile_DirExists(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:direxists.")
	}

	info, err := os.Lstat(params[0].String())
	if err != nil {
		state.RetVal = raptor.NewValueBool(false)
		return
	}
	if !info.IsDir() {
		state.RetVal = raptor.NewValueBool(false)
		return
	}
	state.RetVal = raptor.NewValueBool(true)
	return
}

// Iterate over all the files in a directory.
// 	file:walkfiles path code
// Calles code (as a command) for every file found:
//	code path
// Returns nil or an error message. May set the Error flag.
func CommandFile_WalkFiles(state *raptor.State, params []*raptor.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to file:walkfiles.")
	}

	files, err := ioutil.ReadDir(params[0].String())
	if err != nil {
		state.RetVal = raptor.NewValueString(err.Error())
		state.Error = true
		return
	}

	code := params[1].CompiledScript()
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		state.Envs.Add(raptor.NewEnvironment())

		state.AddParams(file.Name())

		//lex := raptor.NewCompiledLexer(code)
		//for {
		//	lex.Advance()
		//	state.Println(lex.CurrentTkn(), lex.CurrentTkn().Lexeme)
		//	if lex.CheckLookAhead(raptor.TknINVALID) {
		//		return
		//	}
		//}

		state.Code.AddCodeSource(raptor.NewCompiledLexer(code))
		state.Exec()
		state.Envs.Remove()
		state.Return = false
	}
	state.RetVal = nil
}

// Iterate over all the directories in a directory.
// 	file:walkdirs path code
// Calles code (as a command) for every directory found:
//	code path
// Returns nil or an error message. May set the Error flag.
func CommandFile_WalkDirs(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:walkdirs.")
	}

	files, err := ioutil.ReadDir(params[0].String())
	if err != nil {
		state.RetVal = raptor.NewValueString(err.Error())
		state.Error = true
		return
	}

	code := params[1].CompiledScript()
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		state.Envs.Add(raptor.NewEnvironment())

		state.AddParams(file.Name())

		state.Code.AddCodeSource(raptor.NewCompiledLexer(code))
		state.Exec()
		state.Envs.Remove()
		state.Return = false
	}
	state.RetVal = nil
}
