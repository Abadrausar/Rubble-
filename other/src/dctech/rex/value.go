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

package rex

import "strconv"

// Script value internal type codes.
const (
	TypNil     = iota // Nothing (nil)
	TypString         // string
	TypInt            // int64
	TypFloat          // float64
	TypBool           // bool
	TypCode           // *Code
	TypCommand        // NativeCommand
	TypIndex          // Indexable or EditIndexable
	TypUser           // User data, could be anything
)

// Value is a simple script value.
type Value struct {
	Type int
	Data interface{} // Data type stored here depends on value of Type field, see the TypXxx constants.
	Pos  *Position
}

// NewValue creates a new nil Value.
func NewValue() *Value {
	return &Value{
		Type: TypNil,
		Data: nil,
		Pos: NewPosition(1, 1, ""),
	}
}

// NewValueString creates a new Value from a string.
func NewValueString(val string) *Value {
	return &Value{
		Type: TypString,
		Data: val,
		Pos: NewPosition(1, 1, ""),
	}
}

// NewValueInt64 creates a new Value from a int64.
func NewValueInt64(val int64) *Value {
	return &Value{
		Type: TypInt,
		Data: val,
		Pos: NewPosition(1, 1, ""),
	}
}

// NewValueFloat64 creates a new Value from a float64.
func NewValueFloat64(val float64) *Value {
	return &Value{
		Type: TypFloat,
		Data: val,
		Pos: NewPosition(1, 1, ""),
	}
}

// NewValueBool creates a new Value from a bool.
func NewValueBool(val bool) *Value {
	return &Value{
		Type: TypBool,
		Data: val,
		Pos: NewPosition(1, 1, ""),
	}
}

// NewValueCode creates a new Value from a *Code.
func NewValueCode(val *Code) *Value {
	return &Value{
		Type: TypCode,
		Data: val,
		Pos: NewPosition(1, 1, ""),
	}
}

// NewValueCommand creates a new Value from a NativeCommand.
func NewValueCommand(val NativeCommand) *Value {
	return &Value{
		Type: TypCommand,
		Data: val,
		Pos: NewPosition(1, 1, ""),
	}
}

// NewValueIndex creates a new Value from an Indexable.
func NewValueIndex(val Indexable) *Value {
	return &Value{
		Type: TypIndex,
		Data: val,
		Pos: NewPosition(1, 1, ""),
	}
}

// NewValueUser creates a new user Value.
func NewValueUser(val interface{}) *Value {
	return &Value{
		Type: TypUser,
		Data: val,
		Pos: NewPosition(1, 1, ""),
	}
}

// tokenToValue turns a lexer token into a script value.
// This function may panic if an unconvertible token is input.
func tokenToValue(tok *token) *Value {
	val := new(Value)
	val.Pos = tok.Pos

	// String
	if tok.Type == tknString {
		val.Type = TypString
		val.Data = tok.Lexeme
		return val
	}

	// Bool
	if tok.Type == tknTrue {
		val.Type = TypBool
		val.Data = true
		return val
	}
	if tok.Type == tknFalse {
		val.Type = TypBool
		val.Data = false
		return val
	}

	// Nil
	if tok.Type == tknNil {
		val.Type = TypNil
		val.Data = nil
		return val
	}

	// Number or undelimited string
	if tok.Type == tknRawString {
		intval, err := strconv.ParseInt(tok.Lexeme, 0, 64)
		if err == nil {
			val.Type = TypInt
			val.Data = intval
			return val
		}
		
		fltval, err := strconv.ParseFloat(tok.Lexeme, 64)
		if err == nil {
			val.Type = TypFloat
			val.Data = fltval
			return val
		}
		
		val.Type = TypString
		val.Data = tok.Lexeme
		return val
	}
	
	RaiseError("Invalid token conversion, not convertible to value.")
	panic("UNREACHABLE")
}

// CodeString converts a Value to a parseable string.
// The result of parsing a string returned by CodeString may not be the same as the value it was generated from.
// This function is not very useful as it is impossible to save just anything, only some values can be retrieved.
// The output of this function is not to be displayed to the user, this is strictly for implementing
// possible future support for making object literals from existing Indexables for saving state and the like.
// 
// Warning! In some cases (in particularly code values) this function still does not generate proper parseable code.
func (val *Value) CodeString() string {
	switch val.Type {
	case TypNil:
		return "nil"
	case TypString:
		return "\"" + EscapeString(val.Data.(string)) + "\""
	case TypInt:
		return strconv.FormatInt(val.Data.(int64), 10)
	case TypFloat:
		return strconv.FormatFloat(val.Data.(float64), 'g', -1, 64)
	case TypBool:
		if val.Data.(bool) {
			return "true"
		}
		return "false"
	case TypCode:
		code := val.Data.(*Code)
		out := "{ "
		for _, val := range code.data {
			out += val.toLexeme() + " "
		}
		out += "}"
		return out
	case TypCommand:
		return "\"{NativeCommand}\""
	case TypIndex:
		return val.Data.(Indexable).CodeString()
	case TypUser:
		return "\"<UserData>\""
	}
	RaiseError("Script Value has invalid Type.")
	panic("UNREACHABLE")
}

// String converts a Value to a string.
// The value generated by this function is not always useful for anything but display to a user.
func (val *Value) String() string {
	switch val.Type {
	case TypNil:
		return "nil"
	case TypString:
		return val.Data.(string)
	case TypInt:
		return strconv.FormatInt(val.Data.(int64), 10)
	case TypFloat:
		return strconv.FormatFloat(val.Data.(float64), 'g', -1, 64)
	case TypBool:
		if val.Data.(bool) {
			return "true"
		}
		return "false"
	case TypCode:
		code := val.Data.(*Code)
		out := "{ "
		for _, val := range code.data {
			out += val.toLexeme() + " "
		}
		out += "}"
		return out
	case TypCommand:
		return "{NativeCommand}"
	case TypIndex:
		return val.Data.(Indexable).String()
	case TypUser:
		return "<UserData>"
	}
	RaiseError("Script Value has invalid Type.")
	panic("UNREACHABLE")
}

// Int64 converts a Value to an int64.
// Objects and invalid strings are always 0.
func (val *Value) Int64() int64 {
	switch val.Type {
	case TypNil:
		return 0
	case TypString:
		val, err := strconv.ParseInt(val.Data.(string), 0, 64)
		if err != nil {
			val = 0
		}
		return val
	case TypInt:
		return val.Data.(int64)
	case TypFloat:
		return int64(val.Data.(float64))
	case TypBool:
		if val.Data.(bool) {
			return -1
		}
		return 0
	case TypCode:
		return 0
	case TypCommand:
		return 0
	case TypIndex:
		return 0
	case TypUser:
		return 0
	}
	RaiseError("Script Value has invalid Type.")
	panic("UNREACHABLE")
}

// Float64 converts a Value to an float64.
// Objects and invalid strings are always 0.0.
func (val *Value) Float64() float64 {
	switch val.Type {
	case TypNil:
		return 0.0
	case TypString:
		val, err := strconv.ParseFloat(val.Data.(string), 64)
		if err != nil {
			val = 0.0
		}
		return val
	case TypInt:
		return float64(val.Data.(int64))
	case TypFloat:
		return val.Data.(float64)
	case TypBool:
		if val.Data.(bool) {
			return -1
		}
		return 0.0
	case TypCode:
		return 0.0
	case TypCommand:
		return 0.0
	case TypIndex:
		return 0.0
	case TypUser:
		return 0.0
	}
	RaiseError("Script Value has invalid Type.")
	panic("UNREACHABLE")
}

// Bool converts a Value to a bool.
// Strings return false for "0", "false", or "".
// Ints are true if the value is anything other than 0.
// Floats are converted to ints and then converted to bool by the int rules.
// Code is always true.
// User values and indexables are false if the value's data is nil.
// Nil is always false.
func (val *Value) Bool() bool {
	switch val.Type {
	case TypNil:
		return false
	case TypString:
		return (val.Data.(string) != "" && val.Data.(string) != "0" && val.Data.(string) != "false")
	case TypInt:
		return val.Data.(int64) != 0
	case TypFloat:
		return int64(val.Data.(float64)) != 0
	case TypBool:
		return val.Data.(bool)
	case TypCode:
		return true
	case TypCommand:
		return true
	case TypIndex:
		return val.Data != nil
	case TypUser:
		return val.Data != nil
	}
	RaiseError("Script Value has invalid Type.")
	panic("UNREACHABLE")
}

// TypeString will return the string name of a value's type.
func (val *Value) TypeString() string {
	switch val.Type {
	case TypNil:
		return "nil"
	case TypString:
		return "string"
	case TypInt:
		return "int"
	case TypFloat:
		return "float"
	case TypBool:
		return "bool"
	case TypCode:
		return "code"
	case TypCommand:
		return "command"
	case TypIndex:
		return "index"
	case TypUser:
		return "user"
	}
	RaiseError("Script Value has invalid Type.")
	panic("UNREACHABLE")
}

// NativeCommand is the function signature for a native command handler.
type NativeCommand func(*Script, []*Value)

func (val *Value) call(script *Script, params []*Value) {
	if val.Type != TypCommand && val.Type != TypCode {
		RaiseError("Attempt to call non-runnable value.")
	}
	
	// Native command
	if val.Type == TypCommand {
		val.Data.(NativeCommand)(script, params)
		return
	}

	// User command

	// Parameters are always the first variables in a block.
	block := val.Data.(*Code)
	script.Locals.Add(block)
	
	script.setParams(block, params)

	script.Exec(block)
	script.Locals.Remove()
	script.Return = false
	return
}

// Formatting Functions

// EscapeString will escape (some) special chars in a string to (try to) make it into a valid DQString.
func EscapeString(input string) string {
	in := []byte(input)
	out := make([]byte, 0, len(in))
	for i := range in {
		switch in[i] {
		case '\n':
			out = append(out, "\\n"...)
		case '\r':
			out = append(out, "\\r"...)
		case '"':
			out = append(out, "\\\""...)
		case '\\':
			out = append(out, "\\\\"...)
		default:
			out = append(out, in[i])
		}
	}
	return string(out)
}
