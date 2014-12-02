// NCA v6 INI Commands.
package ini

import "dctech/nca6"
import "dctech/ini"

// Adds the ini commands to the state.
// The ini commands are:
//	ini:create
//	ini:parse
//	ini:format
//	ini:getvalue
//	ini:setvalue
func Setup(state *nca6.State) {
	state.NewNameSpace("ini")
	state.NewNativeCommand("ini:create", CommandIni_Create)
	state.NewNativeCommand("ini:parse", CommandIni_Parse)
	state.NewNativeCommand("ini:format", CommandIni_Format)
	state.NewNativeCommand("ini:getvalue", CommandIni_GetValue)
	state.NewNativeCommand("ini:setvalue", CommandIni_SetValue)
}

// Create a new (empty) ini file.
// 	ini:create
// Returns the ini file
func CommandIni_Create(state *nca6.State, params []*nca6.Value) {
	state.RetVal = nca6.NewValueObject(ini.NewFile())
}

// Parse an ini file.
// 	ini:parse filecontents
// Returns the ini file or an error message. May set the Error flag.
func CommandIni_Parse(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to ini:parse.")
	}

	file, err := ini.Parse(params[0].String())
	if err != nil {
		state.Error = true
		state.RetVal = nca6.NewValueString(err.Error())
		return
	}

	state.RetVal = nca6.NewValueObject(file)
}

// Format an ini file for writing to disk.
// 	ini:format file
// Returns the string representation of the ini file or an error message. May set the Error flag.
func CommandIni_Format(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to ini:format.")
	}

	if params[0].Type != nca6.TypObject {
		panic("Non-Object Param 0 passed to ini:format.")
	}
	if _, ok := params[0].Data.(*ini.File); !ok {
		panic("ini:format's Param 0 is not an *ini.File.")
	}
	file := params[0].Data.(*ini.File)
	state.RetVal = nca6.NewValueString(ini.Format(file))
}

// Gets a specific value from an ini file.
// 	ini:getvalue file sectionname keyname
// Returns the value or an error message. May set the Error flag.
func CommandIni_GetValue(state *nca6.State, params []*nca6.Value) {
	if len(params) != 3 {
		panic("Wrong number of params to ini:getvalue.")
	}

	if params[0].Type != nca6.TypObject {
		panic("Non-Object Param 0 passed to ini:getvalue.")
	}
	if _, ok := params[0].Data.(*ini.File); !ok {
		panic("ini:getvalue's Param 0 is not an *ini.File.")
	}
	file := params[0].Data.(*ini.File)

	val, err := file.Get(params[1].String(), params[2].String())
	if err != nil {
		state.Error = true
		state.RetVal = nca6.NewValueString(err.Error())
		return
	}
	state.RetVal = nca6.NewValueString(val)
}

// Sets a specific value in an ini file.
// 	ini:setvalue file sectionname keyname value
// Returns unchanged or an error message. May set the Error flag.
func CommandIni_SetValue(state *nca6.State, params []*nca6.Value) {
	if len(params) != 4 {
		panic("Wrong number of params to ini:setvalue.")
	}

	if params[0].Type != nca6.TypObject {
		panic("Non-Object Param 0 passed to ini:setvalue.")
	}
	if _, ok := params[0].Data.(*ini.File); !ok {
		panic("ini:setvalue's Param 0 is not an *ini.File.")
	}
	file := params[0].Data.(*ini.File)

	// Create is a nop if section/key already exists.
	file.Create(params[1].String(), params[2].String())
	file.Set(params[1].String(), params[2].String(), params[3].String())
}
