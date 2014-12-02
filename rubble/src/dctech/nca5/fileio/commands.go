// NCA v5 File IO Commands.
package fileio

import "dctech/nca5"
import "io/ioutil"

// Adds the file io commands to the state.
// The file io commands are:
//	fileio:read
//	fileio:write
func Setup(state *nca5.State) {
	state.NewNameSpace("fileio")
	state.NewNativeCommand("fileio:read", CommandFileIO_Read)
	state.NewNativeCommand("fileio:write", CommandFileIO_Write)
}

// Read from file.
// 	fileio:read path
// Returns file contents or an error message. May set the Error flag.
func CommandFileIO_Read(state *nca5.State, params []*nca5.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to fileio:read.")
	}

	file, err := ioutil.ReadFile(params[0].String())
	if err != nil {
		state.Error = true
		state.RetVal = nca5.NewValue("error:" + err.Error())
		return
	}
	state.RetVal = nca5.NewValue(string(file))
}

// Write to file.
// 	fileio:write path contents
// Returns unchanged or an error message. May set the Error flag.
func CommandFileIO_Write(state *nca5.State, params []*nca5.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to fileio:write.")
	}

	// I have no idea what "0600" means but I saw it in an example.
	// well I do know that it is a file permission.
	err := ioutil.WriteFile(params[0].String(), []byte(params[1].String()), 0600)
	if err != nil {
		state.Error = true
		state.RetVal = nca5.NewValue("error:" + err.Error())
		return
	}
}
