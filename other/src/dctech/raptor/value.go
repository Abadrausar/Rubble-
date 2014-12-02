/*
For copyright/license see header in file "doc.go"
*/

package raptor

//import "fmt"
import "strconv"

const (
	TypNil     = iota // Nothing.
	TypObject         // Generic data, possibly an Indexable.
	TypString         // A string
	TypInt            // An int64
	TypFloat          // A float64
	TypBool           // A boolean
	TypCode           // A Code object
	TypCommand        // A command reference, basically a string with a special type value.
)

// Why do I have this type? It makes refactoring easier.
type EmptyInterface interface{}

// Value is a simple script value.
type Value struct {
	Type int
	Data EmptyInterface
	Pos  *Position
}

// NewValue creates a new nil Value.
func NewValue() *Value {
	this := new(Value)
	this.Type = TypNil
	this.Data = nil
	this.Pos = NewPosition(0, -1, "")
	return this
}

// NewValueString creates a new Value from a string.
func NewValueString(val string) *Value {
	this := new(Value)
	this.Type = TypString
	this.Data = val
	this.Pos = NewPosition(0, -1, "")
	return this
}

// NewValueInt64 creates a new Value from a int64.
func NewValueInt64(val int64) *Value {
	this := new(Value)
	this.Type = TypInt
	this.Data = val
	this.Pos = NewPosition(0, -1, "")
	return this
}

// NewValueFloat64 creates a new Value from a float64.
func NewValueFloat64(val float64) *Value {
	this := new(Value)
	this.Type = TypFloat
	this.Data = val
	this.Pos = NewPosition(0, -1, "")
	return this
}

// NewValueBool creates a new Value from a bool.
func NewValueBool(val bool) *Value {
	this := new(Value)
	this.Type = TypBool
	this.Data = val
	this.Pos = NewPosition(0, -1, "")
	return this
}

// NewValueCode creates a new Value from a *Code.
func NewValueCode(val *Code) *Value {
	this := new(Value)
	this.Type = TypCode
	this.Data = val
	this.Pos = NewPosition(0, -1, "")
	return this
}

// NewValueObject creates a new Value with the type TypObject.
func NewValueObject(val EmptyInterface) *Value {
	this := new(Value)
	this.Type = TypObject
	this.Data = val
	this.Pos = NewPosition(0, -1, "")
	return this
}

// NewValueCommand creates a new Value with the type TypCommand.
func NewValueCommand(val string) *Value {
	this := new(Value)
	this.Type = TypCommand
	this.Data = val
	this.Pos = NewPosition(0, -1, "")
	return this
}

// TokenToValue turns a lexer token into a script value using the following rules.
// 	if lexeme is "true" or "false" type is bool
//	if lexeme can be converted into an int without error type is int
//	else type is string
func TokenToValue(tok *Token) *Value {
	this := new(Value)
	this.Pos = tok.Pos

	if tok.Lexeme == "true" {
		this.Type = TypBool
		this.Data = true
		return this
	}

	if tok.Lexeme == "false" {
		this.Type = TypBool
		this.Data = false
		return this
	}

	if tok.Lexeme == "nil" {
		this.Type = TypNil
		this.Data = nil
		return this
	}

	intval, err := strconv.ParseInt(tok.Lexeme, 10, 64)
	if err == nil {
		this.Type = TypInt
		this.Data = intval
		return this
	}

	this.Type = TypString
	this.Data = tok.Lexeme
	return this
}

// CodeString converts a Value to a parseable string.
// The result of parsing a string returned by CodeString may not be the same as the value it was generated from.
// This function is not very useful as it is imposible to save just anything, only some values can be retrived.
// Also the code generated by this function may need the command "getcommand" to be defined.
// The output of this function is not to be displayed to the user, this is strictly for implimenting
// possible future support for making object literals from existing Indexables for saving state and the like.
func (this *Value) CodeString() string {
	switch this.Type {
	case TypString:
		return "\"" + EscapeString(this.Data.(string)) + "\""

	case TypInt:
		return strconv.FormatInt(this.Data.(int64), 10)

	case TypFloat:
		return strconv.FormatFloat(this.Data.(float64), 'g', -1, 64)

	case TypBool:
		if this.Data.(bool) {
			return "true"
		}
		return "false"

	case TypCode:
		return this.Code().String()

	case TypCommand:
		return "(getcommand \"" + this.Data.(string) + "\")"

	case TypObject:
		index := this.Indexable()
		if index == nil {
			return "\"<UserData>\""
		}

		return index.CodeString()

	case TypNil:
		return "nil"
	}
	panic("CodeString: Script Value has invalid Type.")
}

// String converts a Value to a string.
// The value generated by this function is not always useful for anything but display to a user.
func (this *Value) String() string {
	switch this.Type {
	case TypString:
		return this.Data.(string)

	case TypInt:
		return strconv.FormatInt(this.Data.(int64), 10)

	case TypFloat:
		return strconv.FormatFloat(this.Data.(float64), 'g', -1, 64)

	case TypBool:
		if this.Data.(bool) {
			return "true"
		}
		return "false"

	case TypCode:
		return this.Code().String()

	case TypCommand:
		return this.Data.(string)

	case TypObject:
		index := this.Indexable()
		if index == nil {
			return "<UserData>"
		}

		return index.String()

	case TypNil:
		return "nil"
	}
	panic("String: Script Value has invalid Type.")
}

// Int64 converts a Value to an int64.
// Objects and invalid strings are always 0.
func (this *Value) Int64() int64 {
	switch this.Type {
	case TypString:
		val, err := strconv.ParseInt(this.Data.(string), 0, 64)
		if err != nil {
			val = 0
		}
		return val

	case TypInt:
		return this.Data.(int64)

	case TypFloat:
		return int64(this.Data.(float64))

	case TypBool:
		if this.Data.(bool) {
			return -1
		}
		return 0

	case TypCode:
		return 0

	case TypCommand:
		return 0

	case TypObject:
		return 0

	case TypNil:
		return 0
	}
	panic("Int64: Script Value has invalid Type.")
}

// Float64 converts a Value to an float64.
// Objects and invalid strings are always 0.0.
func (this *Value) Float64() float64 {
	switch this.Type {
	case TypString:
		val, err := strconv.ParseFloat(this.Data.(string), 64)
		if err != nil {
			val = 0.0
		}
		return val

	case TypInt:
		return float64(this.Data.(int64))

	case TypFloat:
		return this.Data.(float64)

	case TypBool:
		if this.Data.(bool) {
			return 1.0
		}
		return 0.0

	case TypCode:
		return 0.0

	case TypCommand:
		return 0.0

	case TypObject:
		return 0.0

	case TypNil:
		return 0.0
	}
	panic("Float64: Script Value has invalid Type.")
}

// Bool converts a Value to a bool.
// Strings return false for "0", "false", or "".
// Ints are true if the value is anything other than 0.
// Floats are converted to ints and then converted to bool by the int rules.
// Code and commands are always true
// Objects are false if the value's data is nil.
// Nil is always false
func (this *Value) Bool() bool {
	switch this.Type {
	case TypString:
		return (this.Data.(string) != "" && this.Data.(string) != "0" && this.Data.(string) != "false")

	case TypInt:
		return this.Data.(int64) != 0

	case TypFloat:
		return int64(this.Data.(float64)) != 0

	case TypBool:
		return this.Data.(bool)

	case TypCode:
		return true

	case TypCommand:
		return true

	case TypObject:
		return this.Data != nil

	case TypNil:
		return false
	}
	panic("Bool: Script Value has invalid Type.")
}

// Code converts a Value to a Code object.
func (this *Value) Code() *Code {
	switch this.Type {
	case TypString:
		return NewCode(NewLexer(this.Data.(string), this.Pos))
	case TypCode:
		return this.Data.(*Code)
	default:
		return NewCode(NewLexer(this.String(), this.Pos))
	}
	panic("UNREACHABLE")
}

// CodeSource converts a Value to a CodeSource.
func (this *Value) CodeSource() CodeSource {
	switch this.Type {
	case TypString:
		return NewLexer(this.Data.(string), this.Pos)
	case TypCode:
		return NewCodeReader(this.Data.(*Code))
	default:
		return NewLexer(this.String(), this.Pos)
	}
	panic("UNREACHABLE")
}

// Indexable will try to return the value's data as an Indexable, returns nil on failure.
func (this *Value) Indexable() Indexable {
	if this.Type != TypObject {
		return nil
	}
	return ToIndexable(this.Data)
}

// EditIndexable will try to return the value's data as an EditIndexable, returns nil on failure.
func (this *Value) EditIndexable() EditIndexable {
	if this.Type != TypObject {
		return nil
	}
	return ToEditIndexable(this.Data)
}

// TypeString will return the string name of a value's type.
func (this *Value) TypeString() string {
	switch this.Type {
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

	case TypObject:
		return "object"

	case TypNil:
		return "nil"
	}
	panic("Type: Script Value has invalid Type.")
}

// Formatting Functions

// EscapeString will escape (some) special chars in a string to (try to) make it into a valid Raptor DQString.
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
