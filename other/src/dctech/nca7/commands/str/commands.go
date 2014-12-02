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

// NCA v7 String Commands.
package str

import "dctech/nca7"
//import "fmt"
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
//	str:left
//	str:right
func Setup(state *nca7.State) {
	state.NewNameSpace("str")
	state.NewNativeCommand("str:add", CommandStr_Add)
	state.NewNativeCommand("str:trimspace", CommandStr_TrimSpace)
	state.NewNativeCommand("str:len", CommandStr_Len)
	state.NewNativeCommand("str:char", CommandStr_Char)
	state.NewNativeCommand("str:fmt", CommandStr_Fmt)
	state.NewNativeCommand("str:cmp", CommandStr_Cmp)
	state.NewNativeCommand("str:left", CommandStr_Left)
	state.NewNativeCommand("str:right", CommandStr_Right)
}

// Appends two or more strings together.
// 	str:add a b [c...]
// Returns the result of appending all parameters together.
func CommandStr_Add(state *nca7.State, params []*nca7.Value) {
	if len(params) <= 1 {
		panic("str:add needs at least two params.")
	}

	result := ""
	for _, val := range params {
		result += val.String()
	}

	state.RetVal = nca7.NewValueString(result)
	return
}

// Trims leading and trailing whitespace from a string.
// 	str:trimspace str
// Returns str with leading and trailing whitespace removed.
func CommandStr_TrimSpace(state *nca7.State, params []*nca7.Value) {
	if len(params) != 1 {
		panic("Wrong param count to str:trimspace.")
	}

	state.RetVal = nca7.NewValueString(strings.TrimSpace(params[0].String()))
	return
}

// Gets the length of a string.
// 	str:len a
// Returns the length.
func CommandStr_Len(state *nca7.State, params []*nca7.Value) {
	if len(params) != 1 {
		panic("Wrong param count to str:len.")
	}

	state.RetVal = nca7.NewValueInt64(int64(len(params[0].String())))
	return
}

// Gets char at pos.
// 	str:char a pos
// Opperand pos is converted to a 64 bit integer. Invalid strings are given the value "0"
// If the position is out of range returns unchanged and sets the error flag.
// Returns the character.
func CommandStr_Char(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong param count to str:char.")
	}

	pos := params[1].Int64()
	if pos >= int64(len(params[0].String())) {
		state.Error = true
		return
	}
	
	state.RetVal = nca7.NewValueString(string(params[0].String()[pos]))
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
// Returns the formated string.
func CommandStr_Fmt(state *nca7.State, params []*nca7.Value) {
	if len(params) < 2 {
		panic("Wrong param count to str:fmt.")
	}
	
	paramcount := len(params)-1
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
				panic("More format codes than params to str:fmt.")
			}
			output = append(output, []rune(params[escapecount].String())...)
			escape = false
			continue
		}
		
		if escape && val == 'd' {
			escapecount++
			if paramcount < escapecount {
				panic("More format codes than params to str:fmt.")
			}
			temp := strconv.FormatInt(params[escapecount].Int64(), 10)
			output = append(output, []rune(temp)...)
			escape = false
			continue
		}
		
		if escape && val == 'x' {
			escapecount++
			if paramcount < escapecount {
				panic("More format codes than params to str:fmt.")
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
				panic("More format codes than params to str:fmt.")
			}
			temp := strconv.FormatInt(params[escapecount].Int64(), 16)
			//output = append(output, '0', 'X')
			output = append(output, []rune(strings.ToUpper(temp))...)
			escape = false
			continue
		}
		
		if escape {
			panic("Invalid format code: %" + string(val) + " to str:fmt.")
		}
		
		output = append(output, val)
	}
	
	state.RetVal = nca7.NewValueString(string(output))
	return
}

// Compair two strings.
// 	str:cmp a b
// Returns true or false.
func CommandStr_Cmp(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong param count to str:cmp.")
	}
	
	result := params[0].String() == params[1].String()
	state.RetVal = nca7.NewValueBool(result)
}

// Return x charecters from the left side of a string.
// 	str:left str x
// Returns the new string or (if the index is out of range) returns the string and sets the error flag.
func CommandStr_Left(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong param count to str:left.")
	}
	
	index := int(params[1].Int64())
	str := params[0].String()
	if len(str) < index {
		state.RetVal = nca7.NewValueString(str)
		state.Error = true
		return
	}
	state.RetVal = nca7.NewValueString(string([]byte(str))[:index])
}

// Return x charecters from the right side of a string.
// 	str:right str x
// Returns the new string or (if the index is out of range) returns the string and sets the error flag.
func CommandStr_Right(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong param count to str:right.")
	}
	
	index := int(params[1].Int64())
	str := params[0].String()
	if len(str) < index {
		state.RetVal = nca7.NewValueString(str)
		state.Error = true
		return
	}
	state.RetVal = nca7.NewValueString(string([]byte(str))[index:])
}
