// NCA v4 CSV Commands.
package csv

import "dctech/nca4"
import "encoding/csv"
import "bytes"

var csvFiles [][][]string

func init() {
	csvFiles = make([][][]string, 0, 10)
}

// Adds the csv commands to the state.
// The csv commands are:
//	csv:parse
//	csv:getfield
//	csv:recordcount
//	csv:fieldcount
// Note that there is no support for writing csv files
func Setup(state *nca4.State) {
	state.NewNameSpace("csv")
	state.NewNativeCommand("csv:parse", CommandCsv_Parse)
	state.NewNativeCommand("csv:getfield", CommandCsv_GetField)
	state.NewNativeCommand("csv:recordcount", CommandCsv_RecordCount)
	state.NewNativeCommand("csv:fieldcount", CommandCsv_FieldCount)
}

// Parse a csv file.
// 	csv:parse filecontents
// Returns the csv file handle or an error message. May set the Error flag.
func CommandCsv_Parse(state *nca4.State, params []*nca4.Value) {
	handle := len(csvFiles)
	if len(params) != 1 {
		panic("Wrong number of params to csv.parse.")
	}

	csvreader := bytes.NewBuffer([]byte(params[0].String()))
	file, err := csv.NewReader(csvreader).ReadAll()
	if err != nil {
		state.Error = true
		state.RetVal = nca4.NewValue("error:" + err.Error())
		return
	}

	csvFiles = append(csvFiles, file)

	state.RetVal = nca4.NewValueFromI64(int64(handle))
}

// Gets a specific field from a csv file.
// 	csv:getfield handle record field
// Returns the value or an error message. May set the Error flag.
func CommandCsv_GetField(state *nca4.State, params []*nca4.Value) {
	if len(params) != 3 {
		panic("Wrong number of params to csv:getfield.")
	}

	handle := int(params[0].Int64())
	if handle >= len(csvFiles) {
		state.Error = true
		state.RetVal = nca4.NewValue("error:Invalid Handle.")
		return
	}

	record := int(params[1].Int64())
	if record >= len(csvFiles[handle]) {
		state.Error = true
		state.RetVal = nca4.NewValue("error:Invalid Record.")
		return
	}

	field := int(params[2].Int64())
	if field >= len(csvFiles[handle][record]) {
		state.Error = true
		state.RetVal = nca4.NewValue("error:Invalid Field.")
		return
	}

	state.RetVal = nca4.NewValue(csvFiles[handle][record][field])
}

// Gets the number of records in a csv file.
// 	csv:recordcount handle
// Returns the count or an error message. May set the Error flag.
func CommandCsv_RecordCount(state *nca4.State, params []*nca4.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to csv:recordcount.")
	}

	handle := int(params[0].Int64())
	if handle >= len(csvFiles) {
		state.Error = true
		state.RetVal = nca4.NewValue("error:Invalid Handle.")
		return
	}

	state.RetVal = nca4.NewValueFromI64(int64(len(csvFiles[handle])))
}

// Gets a specific field from a csv file.
// 	csv:fieldcount handle record
// Returns the count or an error message. May set the Error flag.
func CommandCsv_FieldCount(state *nca4.State, params []*nca4.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to csv:fieldcount.")
	}

	handle := int(params[0].Int64())
	if handle >= len(csvFiles) {
		state.Error = true
		state.RetVal = nca4.NewValue("error:Invalid Handle.")
		return
	}

	record := int(params[1].Int64())
	if record >= len(csvFiles[handle]) {
		state.Error = true
		state.RetVal = nca4.NewValue("error:Invalid Record.")
		return
	}

	state.RetVal = nca4.NewValueFromI64(int64(len(csvFiles[handle][record])))
}
