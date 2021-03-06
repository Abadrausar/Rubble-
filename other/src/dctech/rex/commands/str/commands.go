/*
Copyright 2014 by Milo Christiansen

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

// Rex String Commands.
package str

import "dctech/rex"
import "strconv"
import "strings"

// Adds the string commands to the state.
// The string commands are:
//	str:add
//	str:trimspace
//	str:len
//	str:char
//	str:fmt
//	str:cmp
//	str:find
//	str:left
//	str:trimleft
//	str:right
//	str:trimright
//	str:mid
//	str:replace
//	str:tolower
//	str:toupper
func Setup(state *rex.State) (err error) {
	defer func() {
		if !state.NoRecover {
			if x := recover(); x != nil {
				if y, ok := x.(rex.ScriptError); ok {
					err = y
					return
				}
				panic(x)
			}
		}
	}()
	
	mod := state.RegisterModule("str")
	mod.RegisterCommand("add", Command_Add)
	mod.RegisterCommand("trimspace", Command_TrimSpace)
	mod.RegisterCommand("len", Command_Len)
	mod.RegisterCommand("char", Command_Char)
	mod.RegisterCommand("fmt", Command_Fmt)
	mod.RegisterCommand("cmp", Command_Cmp)
	mod.RegisterCommand("find", Command_Find)
	mod.RegisterCommand("left", Command_Left)
	mod.RegisterCommand("trimleft", Command_TrimLeft)
	mod.RegisterCommand("right", Command_Right)
	mod.RegisterCommand("trimright", Command_TrimRight)
	mod.RegisterCommand("mid", Command_Mid)
	mod.RegisterCommand("replace", Command_Replace)
	mod.RegisterCommand("tolower", Command_ToLower)
	mod.RegisterCommand("toupper", Command_ToUpper)
	
	return nil
}

// Appends two or more strings together.
// 	str:add a b [c...]
// Returns the result of appending all parameters together.
func Command_Add(script *rex.Script, params []*rex.Value) {
	if len(params) < 2 {
		rex.ErrorParamCount("str:add", ">=2")
	}

	result := ""
	for _, val := range params {
		result += val.String()
	}

	script.RetVal = rex.NewValueString(result)
	return
}

// Trims leading and trailing whitespace from a string.
// 	str:trimspace str
// Returns str with leading and trailing whitespace removed.
func Command_TrimSpace(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("str:trimspace", "1")
	}

	script.RetVal = rex.NewValueString(strings.TrimSpace(params[0].String()))
	return
}

// Gets the length of a string.
// 	str:len a
// Returns the length.
func Command_Len(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("str:len", "1")
	}

	script.RetVal = rex.NewValueInt64(int64(len(params[0].String())))
	return
}

// Gets char at pos.
// 	str:char a pos
// Operand pos is converted to a 64 bit integer. Invalid strings are given the value 0
// If the position is out of range returns unchanged and sets the error flag.
// Returns the character.
func Command_Char(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("str:char", "2")
	}

	pos := params[1].Int64()
	if pos >= int64(len(params[0].String())) {
		script.Error = true
		return
	}

	script.RetVal = rex.NewValueString(string(params[0].String()[pos]))
	return
}

// Formats a string.
// 	str:fmt format values...
// valid format "verbs" for fmtstr:
//	%% literal percent
//	%s the raw string
//	%d as a decimal number
//	%x as a lowercase hexadecimal number
//	%X as an upercase hexadecimal number
// Returns the formatted string.
func Command_Fmt(script *rex.Script, params []*rex.Value) {
	if len(params) < 2 {
		rex.ErrorParamCount("str:fmt", "2")
	}

	paramcount := len(params) - 1
	escapecount := 0
	escape := false
	output := make([]rune, 0, 100)
	for _, val := range params[0].String() {
		if val == '%' && escape {
			output = append(output, val)
			escape = false
			continue
		}
		if val == '%' && !escape {
			escape = true
			continue
		}

		if escape && val == 's' {
			escapecount++
			if paramcount < escapecount {
				rex.ErrorGeneralCmd("str:fmt", "More format codes than parameters.")
			}
			output = append(output, []rune(params[escapecount].String())...)
			escape = false
			continue
		}

		if escape && val == 'd' {
			escapecount++
			if paramcount < escapecount {
				rex.ErrorGeneralCmd("str:fmt", "More format codes than parameters.")
			}
			temp := strconv.FormatInt(params[escapecount].Int64(), 10)
			output = append(output, []rune(temp)...)
			escape = false
			continue
		}

		if escape && val == 'x' {
			escapecount++
			if paramcount < escapecount {
				rex.ErrorGeneralCmd("str:fmt", "More format codes than parameters.")
			}
			temp := strconv.FormatInt(params[escapecount].Int64(), 16)
			//output = append(output, '0', 'x')
			output = append(output, []rune(temp)...)
			escape = false
			continue
		}

		if escape && val == 'X' {
			escapecount++
			if paramcount < escapecount {
				rex.ErrorGeneralCmd("str:fmt", "More format codes than parameters.")
			}
			temp := strconv.FormatInt(params[escapecount].Int64(), 16)
			//output = append(output, '0', 'X')
			output = append(output, []rune(strings.ToUpper(temp))...)
			escape = false
			continue
		}

		if escape {
			rex.ErrorGeneralCmd("str:fmt", "Invalid format code: \"%"+string(val)+"\".")
		}

		output = append(output, val)
	}

	script.RetVal = rex.NewValueString(string(output))
	return
}

// Compare two strings.
// 	str:cmp a b
// Returns true or false.
func Command_Cmp(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("str:cmp", "2")
	}

	result := params[0].String() == params[1].String()
	script.RetVal = rex.NewValueBool(result)
}

// Search for a substring in a string.
// 	str:find string substring
// Returns the position of the substring or -1 if the substring is not present.
func Command_Find(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("str:find", "2")
	}

	script.RetVal = rex.NewValueInt64(int64(strings.Index(params[0].String(), params[1].String())))
}

// Return x characters from the left side of a string.
// 	str:left str x
// Returns the new string or (if the index is out of range) returns the string and sets the error flag.
func Command_Left(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("str:left", "2")
	}

	index := int(params[1].Int64())
	str := params[0].String()
	if len(str) < index {
		script.RetVal = rex.NewValueString(str)
		script.Error = true
		return
	}
	script.RetVal = rex.NewValueString(string([]byte(str))[:index])
}

// Remove x characters from the left side of a string.
// 	str:trimleft str x
// Returns the new string or (if the index is out of range) returns the string and sets the error flag.
func Command_TrimLeft(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("str:trimleft", "2")
	}

	index := int(params[1].Int64())
	str := params[0].String()
	if len(str) < index {
		script.RetVal = rex.NewValueString(str)
		script.Error = true
		return
	}
	script.RetVal = rex.NewValueString(string([]byte(str))[index:])
}

// Return x characters from the right side of a string.
// 	str:right str x
// Returns the new string or (if the index is out of range) returns the string and sets the error flag.
func Command_Right(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("str:right", "2")
	}

	index := int(params[1].Int64())
	str := params[0].String()
	if len(str) < index {
		script.RetVal = rex.NewValueString(str)
		script.Error = true
		return
	}
	script.RetVal = rex.NewValueString(string([]byte(str))[len(str)-index:])
}

// Remove x characters from the right side of a string.
// 	str:trimright str x
// Returns the new string or (if the index is out of range) returns the string and sets the error flag.
func Command_TrimRight(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("str:trimright", "2")
	}

	index := int(params[1].Int64())
	str := params[0].String()
	if len(str) < index {
		script.RetVal = rex.NewValueString(str)
		script.Error = true
		return
	}
	script.RetVal = rex.NewValueString(string([]byte(str))[:len(str)-index])
}

// Returns count characters from a string.
// 	str:mid str start count
// Returns the new string or (if start or count is out of range) returns as close a
// result as possible and sets the error flag.
func Command_Mid(script *rex.Script, params []*rex.Value) {
	if len(params) != 3 {
		rex.ErrorParamCount("str:mid", "3")
	}

	start := int(params[1].Int64())
	count := int(params[2].Int64())
	str := params[0].String()
	if len(str) < start {
		script.RetVal = rex.NewValueString("")
		script.Error = true
		return
	}
	if len(str) < start+count {
		script.RetVal = rex.NewValueString(string([]byte(str))[start:])
		script.Error = true
		return
	}
	script.RetVal = rex.NewValueString(string([]byte(str))[start : start+count])
}

// Replaces search with replace in source.
// 	str:replace source search replace occurrence
// Occurrence gives a number of times to carry out the replacement,
// use -1 to replace all.
// Returns the new string.
func Command_Replace(script *rex.Script, params []*rex.Value) {
	if len(params) != 4 {
		rex.ErrorParamCount("str:replace", "4")
	}

	tmp := strings.Replace(params[0].String(), params[1].String(), params[2].String(), int(params[3].Int64()))
	script.RetVal = rex.NewValueString(tmp)
}

// Converts str to lower case.
// 	str:tolower str
// Returns the string with all letters converted to lower case.
func Command_ToLower(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("str:tolower", "1")
	}

	script.RetVal = rex.NewValueString(strings.ToLower(params[0].String()))
}

// Converts str to upper case.
// 	str:toupper str
// Returns the string with all letters converted to upper case.
func Command_ToUpper(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("str:toupper", "1")
	}

	script.RetVal = rex.NewValueString(strings.ToUpper(params[0].String()))
}
