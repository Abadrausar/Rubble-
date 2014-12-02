// NCA v5 CSV Commands.
package csv

import "dctech/nca6"
import "encoding/csv"
import "bytes"

// Adds the csv commands to the state.
// The csv commands are:
//	csv:parse
//	csv:getfield
//	csv:recordcount
//	csv:fieldcount
// Note that there is no support for writing csv files
func Setup(state *nca6.State) {
	state.NewNameSpace("csv")
	state.NewNativeCommand("csv:parse", CommandCsv_Parse)
}

// Parse a csv file.
// The returned file is an Indexable of Indexables. Indexing rules are the same as for array.
// 	csv:parse filecontents
// Returns the csv file or an error message. May set the Error flag.
func CommandCsv_Parse(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to csv.parse.")
	}

	csvreader := bytes.NewBuffer([]byte(params[0].String()))
	file, err := csv.NewReader(csvreader).ReadAll()
	if err != nil {
		state.Error = true
		state.RetVal = nca6.NewValueString("error:" + err.Error())
		return
	}

	records := make([]*nca6.Value, len(file))
	for i := range file {
		fields := make([]*nca6.Value, len(file[i]))
		
		for x := range file[i] {
			fields[x] = nca6.NewValueString(file[i][x])
		}
		records[i] = nca6.NewValueObject(CSVFile(fields))
	}
	state.RetVal = nca6.NewValueObject(CSVFile(records))
}
