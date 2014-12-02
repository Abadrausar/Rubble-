// NCA v4 File System Commands.
package file

import "dctech/nca4"
import "io/ioutil"
import "os"

// Adds the file io commands to the state.
// The file io commands are:
//	fileio:read
//	fileio:write
func Setup(state *nca4.State) {
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
// Returns unchanged. Sets the Error flag if dir is not empty.
func CommandFile_DelDir(state *nca4.State, params []*nca4.Value) {
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
// Returns unchanged. Sets the Error flag if dir is not empty.
func CommandFile_DelTree(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:deltree.")
	}
	
	
	err := os.RemoveAll(params[0].String())
	if err != nil {
		state.Error = true
	}
	return
}

// Create a new directory.
// 	file:newdir path
// Returns unchanged. 
func CommandFile_NewDir(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:newdir.")
	}
	
	err := os.Mkdir(params[0].String(), 0600)
	if err != nil {
		state.Error = true
	}
	return
}

// Delete a file.
// 	file:del path
// Returns unchanged. May set the Error flag.
func CommandFile_Del(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:del.")
	}
	
	info, err := os.Lstat(params[0].String())
	if err != nil {
		state.Error = true
		return
	}
	if !info.IsDir() {
		err := os.Remove(params[0].String())
		if err != nil {
			state.Error = true
		}
		return
	}
	state.Error = true
	return
}

// Checks if a fileexists.
// 	file:exists path
// Returns 0 or -1.
func CommandFile_Exists(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:exists.")
	}
	
	info, err := os.Lstat(params[0].String())
	if err != nil {
		state.RetVal = nca4.NewValue("0")
		return
	}
	if info.IsDir() {
		state.RetVal = nca4.NewValue("0")
		return
	}
	state.RetVal = nca4.NewValue("-1")
	return
}

// Checks if a directory exists.
// 	file:direxists path
// Returns 0 or -1.
func CommandFile_DirExists(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:direxists.")
	}
	
	info, err := os.Lstat(params[0].String())
	if err != nil {
		state.RetVal = nca4.NewValue("0")
		return
	}
	if !info.IsDir() {
		state.RetVal = nca4.NewValue("0")
		return
	}
	state.RetVal = nca4.NewValue("-1")
	return
}

// Iterate over all the files in a directory.
// 	file:walkfiles path code
// Calles code (as a command) for every file found:
//	code path
// Returns nil or an error message. May set the Error flag.
func CommandFile_WalkFiles(state *nca4.State, params []*nca4.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to file:walkfiles.")
	}
	
	files, err := ioutil.ReadDir(params[0].String())
	if err != nil {
		state.RetVal = nca4.NewValue(err.Error())
		state.Error = true
		return
	}
	
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		state.Envs.Add(nca4.NewEnvironment())
		
		state.AddParams(file.Name())
	
		state.Code.AddLexer(params[1].AsLexer())
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
// Returns nil.
func CommandFile_WalkDirs(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to file:walkdirs.")
	}
	
	files, err := ioutil.ReadDir(params[0].String())
	if err != nil {
		state.RetVal = nca4.NewValue(err.Error())
		state.Error = true
		return
	}
	
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		
		state.Envs.Add(nca4.NewEnvironment())
		
		state.AddParams(file.Name())
	
		state.Code.AddLexer(params[1].AsLexer())
		state.Exec()
		state.Envs.Remove()
		state.Return = false
	}
	state.RetVal = nil
}
