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

// Raptor INI Commands.
package ini

import "dctech/raptor"
import "dctech/ini"

// Adds the ini commands to the state.
// The ini commands are:
//	ini:create
//	ini:parse
//	ini:format
//	ini:getvalue
//	ini:setvalue
func Setup(state *raptor.State) {
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
func CommandIni_Create(state *raptor.State, params []*raptor.Value) {
	state.RetVal = raptor.NewValueObject(ini.NewFile())
}

// Parse an ini file.
// 	ini:parse filecontents
// Returns the ini file or an error message. May set the Error flag.
func CommandIni_Parse(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to ini:parse.")
	}

	file, err := ini.Parse(params[0].String())
	if err != nil {
		state.Error = true
		state.RetVal = raptor.NewValueString(err.Error())
		return
	}

	state.RetVal = raptor.NewValueObject(file)
}

// Format an ini file for writing to disk.
// 	ini:format file
// Returns the string representation of the ini file or an error message. May set the Error flag.
func CommandIni_Format(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to ini:format.")
	}

	if params[0].Type != raptor.TypObject {
		panic("Non-Object Param 0 passed to ini:format.")
	}
	if _, ok := params[0].Data.(*ini.File); !ok {
		panic("ini:format's Param 0 is not an *ini.File.")
	}
	file := params[0].Data.(*ini.File)
	state.RetVal = raptor.NewValueString(ini.Format(file))
}

// Gets a specific value from an ini file.
// 	ini:getvalue file sectionname keyname
// Returns the value or an error message. May set the Error flag.
func CommandIni_GetValue(state *raptor.State, params []*raptor.Value) {
	if len(params) != 3 {
		panic("Wrong number of params to ini:getvalue.")
	}

	if params[0].Type != raptor.TypObject {
		panic("Non-Object Param 0 passed to ini:getvalue.")
	}
	if _, ok := params[0].Data.(*ini.File); !ok {
		panic("ini:getvalue's Param 0 is not an *ini.File.")
	}
	file := params[0].Data.(*ini.File)

	val, err := file.Get(params[1].String(), params[2].String())
	if err != nil {
		state.Error = true
		state.RetVal = raptor.NewValueString(err.Error())
		return
	}
	state.RetVal = raptor.NewValueString(val)
}

// Sets a specific value in an ini file.
// 	ini:setvalue file sectionname keyname value
// Returns unchanged or an error message. May set the Error flag.
func CommandIni_SetValue(state *raptor.State, params []*raptor.Value) {
	if len(params) != 4 {
		panic("Wrong number of params to ini:setvalue.")
	}

	if params[0].Type != raptor.TypObject {
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
