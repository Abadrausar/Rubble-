// NCA v6 File System Commands.
package file

import "dctech/nca6"
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
func Setup(state *nca6.State) {
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
func CommandFile_DelDir(state *nca6.State, params []*nca6.Value) {
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
func CommandFile_DelTree(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:deltree.")
	}
	
	info, err := os.Lstat(params[0].String())
	if err != nil {
		state.RetVal = nca6.NewValueString(err.Error())
		state.Error = true
		return
	}
	if info.IsDir() {
		err := os.RemoveAll(params[0].String())
		if err != nil {
			state.RetVal = nca6.NewValueString(err.Error())
			state.Error = true
		}
		return
	}
	state.RetVal = nca6.NewValueString("File not a directory.")
	state.Error = true
	return
}

// Create a new directory.
// 	file:newdir path
// Returns unchanged or an error message. May set the Error flag.
func CommandFile_NewDir(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:newdir.")
	}
	
	err := os.Mkdir(params[0].String(), 0600)
	if err != nil {
		state.RetVal = nca6.NewValueString(err.Error())
		state.Error = true
	}
	return
}

// Delete a file.
// 	file:del path
// Returns unchanged or an error message. May set the Error flag.
func CommandFile_Del(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:del.")
	}
	
	info, err := os.Lstat(params[0].String())
	if err != nil {
		state.RetVal = nca6.NewValueString(err.Error())
		state.Error = true
		return
	}
	if !info.IsDir() {
		err := os.Remove(params[0].String())
		if err != nil {
			state.RetVal = nca6.NewValueString(err.Error())
			state.Error = true
		}
		return
	}
	state.RetVal = nca6.NewValueString("File is a directory.")
	state.Error = true
	return
}

// Checks if a file exists.
// 	file:exists path
// Returns 0 or -1.
func CommandFile_Exists(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:exists.")
	}
	
	info, err := os.Lstat(params[0].String())
	if err != nil {
		state.RetVal = nca6.NewValueInt64(0)
		return
	}
	if info.IsDir() {
		state.RetVal = nca6.NewValueInt64(0)
		return
	}
	state.RetVal = nca6.NewValueInt64(-1)
	return
}

// Checks if a directory exists.
// 	file:direxists path
// Returns 0 or -1.
func CommandFile_DirExists(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:direxists.")
	}
	
	info, err := os.Lstat(params[0].String())
	if err != nil {
		state.RetVal = nca6.NewValueInt64(0)
		return
	}
	if !info.IsDir() {
		state.RetVal = nca6.NewValueInt64(0)
		return
	}
	state.RetVal = nca6.NewValueInt64(-1)
	return
}

// Iterate over all the files in a directory.
// 	file:walkfiles path code
// Calles code (as a command) for every file found:
//	code path
// Returns nil or an error message. May set the Error flag.
func CommandFile_WalkFiles(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to file:walkfiles.")
	}
	
	files, err := ioutil.ReadDir(params[0].String())
	if err != nil {
		state.RetVal = nca6.NewValueString(err.Error())
		state.Error = true
		return
	}
	
	code := nca6.Compile(params[1].String(), params[1].Line)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		state.Envs.Add(nca6.NewEnvironment())
		
		state.AddParams(file.Name())
	
		state.Code.AddCodeSource(nca6.NewCompiledLexer(code))
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
func CommandFile_WalkDirs(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:walkdirs.")
	}
	
	files, err := ioutil.ReadDir(params[0].String())
	if err != nil {
		state.RetVal = nca6.NewValueString(err.Error())
		state.Error = true
		return
	}
	
	code := nca6.Compile(params[1].String(), params[1].Line)
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		
		state.Envs.Add(nca6.NewEnvironment())
		
		state.AddParams(file.Name())
	
		state.Code.AddCodeSource(nca6.NewCompiledLexer(code))
		state.Exec()
		state.Envs.Remove()
		state.Return = false
	}
	state.RetVal = nil
}
