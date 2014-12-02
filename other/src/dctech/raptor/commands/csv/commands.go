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

// Raptor CSV Commands.
package csv

import "dctech/raptor"
import "encoding/csv"
import "bytes"

// Adds the csv commands to the state.
// The csv commands are:
//	csv:parse
// Note that there is no support for writing csv files
func Setup(state *raptor.State) {
	state.NewNameSpace("csv")
	state.NewNativeCommand("csv:parse", CommandCsv_Parse)
}

// Parse a csv file.
// The returned file is an Indexable of Indexables. Indexing rules are the same as for array.
// 	csv:parse filecontents
// Returns the csv file or an error message. May set the Error flag.
func CommandCsv_Parse(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to csv.parse.")
	}

	csvreader := bytes.NewBuffer([]byte(params[0].String()))
	file, err := csv.NewReader(csvreader).ReadAll()
	if err != nil {
		state.Error = true
		state.RetVal = raptor.NewValueString("error:" + err.Error())
		return
	}

	records := make([]*raptor.Value, len(file))
	for i := range file {
		fields := make([]*raptor.Value, len(file[i]))

		for x := range file[i] {
			fields[x] = raptor.NewValueString(file[i][x])
		}
		tmp := CSVFile(fields)
		records[i] = raptor.NewValueObject(&tmp)
	}
	tmp := CSVFile(records)
	state.RetVal = raptor.NewValueObject(&tmp)
}
