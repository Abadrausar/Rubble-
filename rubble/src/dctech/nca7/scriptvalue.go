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

package nca7

//import "fmt"
import "strconv"

const (
	TypString = iota
	TypInt
	TypFloat
	TypBool
	TypObject
)

// ScriptObject is the interface that a Value's object data must impliment.
type ScriptObject interface {
	
}

// Value is a simple script value.
type Value struct {
	Type int
	Data interface{} // A string, float64, int64, or ScriptObject
	Line   int
	Column int
}

// NewValueString creates a new Value from a string
func NewValueString(val string) *Value {
	this := new(Value)
	this.Type = TypString
	this.Data = val
	this.Line = -1
	this.Column = -1
	return this
}

// NewValueInt64 creates a new Value from a int64
func NewValueInt64(val int64) *Value {
	this := new(Value)
	this.Type = TypInt
	this.Data = val
	this.Line = -1
	this.Column = -1
	return this
}

// NewValueFloat64 creates a new Value from a float64
func NewValueFloat64(val float64) *Value {
	this := new(Value)
	this.Type = TypFloat
	this.Data = val
	this.Line = -1
	this.Column = -1
	return this
}

// NewValueBool creates a new Value from a bool
func NewValueBool(val bool) *Value {
	this := new(Value)
	this.Type = TypBool
	this.Data = val
	this.Line = -1
	this.Column = -1
	return this
}

// NewValueObject creates a new Value from a ScriptObject
func NewValueObject(val ScriptObject) *Value {
	this := new(Value)
	this.Type = TypObject
	this.Data = val
	this.Line = -1
	this.Column = -1
	return this
}

// TokenToValue turns a lexer token into a script value using the following rules.
// 	if lexeme is "true" or "false" type is bool
//	if lexeme can be converted into an int without error type is int
//	if lexeme can be converted into a float without error type is float
//	else type is string
func TokenToValue(tok *Token) *Value {
	this := new(Value)
	this.Line = tok.Line
	this.Column = tok.Column
	
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
	
	intval, err := strconv.ParseInt(tok.Lexeme, 0, 64)
	if err == nil {
		this.Type = TypInt
		this.Data = intval
		return this
	}
	
	fltval, err := strconv.ParseFloat(tok.Lexeme, 64)
	if err == nil {
		this.Type = TypFloat
		this.Data = fltval
		return this
	}
	
	this.Type = TypString
	this.Data = tok.Lexeme
	return this
}

// String converts a Value to a string.
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
	
	case TypObject:
		obj := this.Indexable()
		if obj == nil {
			return "<ScriptObject>"
		}
		
		out := "<Indexable"
		keys := obj.Keys()
		for i := range keys {
			out += ", {" + keys[i] + "}={" + obj.Get(keys[i]).String() + "}"
		}
		return out + ">"
	}
	panic("Script Value has invalid Type.")
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
	
	case TypObject:
		return 0
	}
	panic("Script Value has invalid Type.")
}

// Bool converts a Value to a bool.
// Strings return false for "0", "false", or "".
// Ints are true if the value is anything other than 0.
// Floats are converted to ints and then converted to bool by the int rules.
// Objects are false if the value's data is nil.
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
	
	case TypObject:
		return this.Data != nil
	}
	panic("Script Value has invalid Type.")
}

// AsLexer converts a Value to a *Lexer.
// This will probably be removed.
func (this *Value) AsLexer() *Lexer {
	if this.Type != TypString {
		panic("Only strings may be used as lexers.")
	}
	return NewLexer(this.Data.(string), this.Line, this.Column)
}

// Object handling

// Indexable will try to return the value's data as a Indexable, returns nil on failure.
func (this *Value) Indexable() Indexable {
	if this.Type != TypObject {
		return nil
	}
	return ToIndexable(this.Data)
}

// EditIndexable will try to return the value's data as a EditIndexable, returns nil on failure.
func (this *Value) EditIndexable() EditIndexable {
	if this.Type != TypObject {
		return nil
	}
	return ToEditIndexable(this.Data)
}

