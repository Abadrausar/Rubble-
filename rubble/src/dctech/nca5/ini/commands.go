// NCA v5 INI Commands.
package ini

import "dctech/nca5"
import "dctech/ini"

var iniFiles []*ini.File

func init() {
	iniFiles = make([]*ini.File, 0, 10)
}

// Adds the ini commands to the state.
// The ini commands are:
//	ini:create
//	ini:parse
//	ini:format
//	ini:getvalue
//	ini:setvalue
func Setup(state *nca5.State) {
	state.NewNameSpace("ini")
	state.NewNativeCommand("ini:create", CommandIni_Create)
	state.NewNativeCommand("ini:parse", CommandIni_Parse)
	state.NewNativeCommand("ini:format", CommandIni_Format)
	state.NewNativeCommand("ini:getvalue", CommandIni_GetValue)
	state.NewNativeCommand("ini:setvalue", CommandIni_SetValue)
}

// Create a new (empty) ini file.
// 	ini:create
// Returns the ini file handle
func CommandIni_Create(state *nca5.State, params []*nca5.Value) {
	filehandle := len(iniFiles)

	iniFiles = append(iniFiles, ini.NewFile())

	state.RetVal = nca5.NewValueFromI64(int64(filehandle))
}

// Parse an ini file.
// 	ini:parse filecontents
// Returns the ini file handle or an error message. May set the Error flag.
func CommandIni_Parse(state *nca5.State, params []*nca5.Value) {
	filehandle := len(iniFiles)
	if len(params) != 1 {
		panic("Wrong number of params to ini:parse.")
	}

	file, err := ini.Parse(params[0].String())
	if err != nil {
		state.Error = true
		state.RetVal = nca5.NewValue("error:" + err.Error())
		return
	}

	iniFiles = append(iniFiles, file)

	state.RetVal = nca5.NewValueFromI64(int64(filehandle))
}

// Format an ini file for writing to disk.
// 	ini:format filehandle
// Returns the string representation of the ini file or an error message. May set the Error flag.
func CommandIni_Format(state *nca5.State, params []*nca5.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to ini:format.")
	}

	filehandle := int(params[0].Int64())
	if filehandle >= len(iniFiles) {
		state.Error = true
		state.RetVal = nca5.NewValue("error:Invalid Handle.")
		return
	}

	state.RetVal = nca5.NewValue(ini.Format(iniFiles[filehandle]))
}

// Gets a specific value from an ini file.
// 	ini:getvalue inihandle sectionname keyname
// Returns the value or an error message. May set the Error flag.
func CommandIni_GetValue(state *nca5.State, params []*nca5.Value) {
	if len(params) != 3 {
		panic("Wrong number of params to ini:getvalue.")
	}

	filehandle := int(params[0].Int64())
	if filehandle >= len(iniFiles) {
		state.Error = true
		state.RetVal = nca5.NewValue("error:Invalid Handle.")
		return
	}

	val, err := iniFiles[filehandle].Get(params[1].String(), params[2].String())
	if err != nil {
		state.Error = true
		state.RetVal = nca5.NewValue("error:" + err.Error())
		return
	}
	state.RetVal = nca5.NewValue(val)
}

// Sets a specific value in an ini file.
// 	ini:setvalue inihandle sectionname keyname value
// Returns unchanged or an error message. May set the Error flag.
func CommandIni_SetValue(state *nca5.State, params []*nca5.Value) {
	if len(params) != 4 {
		panic("Wrong number of params to ini:setvalue.")
	}

	filehandle := int(params[0].Int64())
	if filehandle >= len(iniFiles) {
		state.Error = true
		state.RetVal = nca5.NewValue("error:Invalid Handle.")
		return
	}

	// Create is a nop if section/key already exists.
	iniFiles[filehandle].Create(params[1].String(), params[2].String())
	iniFiles[filehandle].Set(params[1].String(), params[2].String(), params[3].String())
}
