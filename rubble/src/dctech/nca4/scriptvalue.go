package nca4

//import "fmt"
import "strconv"
import "dctech/ncalex"

// Value is a simple script value. Basically a string with (if possible) line and column info from the lexer.
// The position info is used if the value is ever used as code.
type Value struct {
	data   string
	line   int
	column int
}

// NewValue creates a new value from a string
func NewValue(str string) *Value {
	val := new(Value)
	val.data = str
	val.line = -1
	val.column = -1
	return val
}

// NewValueFromI64 creates a new value from an int64.
func NewValueFromI64(input int64) *Value {
	val := new(Value)
	val.data = strconv.FormatInt(input, 10)
	val.line = -1
	val.column = -1
	return val
}

// TokenToValue turns a lexer token into a script value.
func TokenToValue(tok *ncalex.Token) *Value {
	val := new(Value)
	val.data = tok.Lexeme
	val.line = tok.Line
	val.column = tok.Column
	return val
}

// String converts a Value to a string.
func (this *Value) String() string {
	return this.data
}

// Int64 converts a Value to an int64.
func (this *Value) Int64() int64 {
	val, err := strconv.ParseInt(this.data, 0, 64)
	if err != nil {
		val = 0
	}
	return val
}

// UInt8 converts a Value to an uint8.
// Only used in lshift and rshift.
func (this *Value) UInt8() uint8 {
	val, err := strconv.ParseInt(this.data, 0, 8)
	if err != nil {
		val = 0
	}
	return uint8(val)
}

// Bool converts a Value to a bool.
// "0" or "" evaluate as false everything else is true.
func (this *Value) Bool() bool {
	if this.data == "" || this.data == "0" {
		return false
	}
	return true
}

// AsLexer converts a Value to a *ncalex.Lexer.
func (this *Value) AsLexer() *ncalex.Lexer {
	return ncalex.NewLexer(this.data, this.line, this.column)
}
