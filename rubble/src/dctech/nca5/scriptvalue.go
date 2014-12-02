package nca5

//import "fmt"
import "strconv"
import "dctech/ncalex"

type Indexable interface {
	Get(string) *Value
	Set(string, *Value)
	Exists(string) bool
	Len() int64
	Keys() []string
}

// Value is a simple script value.
type Value struct {
	data   string
	object Indexable
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

// NewValueFromObject creates a new value from am Indexable
func NewValueFromObject(val Indexable) *Value {
	rtn := new(Value)
	rtn.data = strconv.FormatInt(val.Len(), 10)
	rtn.object = val
	rtn.line = -1
	rtn.column = -1
	return rtn
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
	if this.object != nil {
		panic("Attempt to use object *Value as lexer")
	}
	return ncalex.NewLexer(this.data, this.line, this.column)
}

// Object handling

// ClearObject clears the value's object data
func (this *Value) ClearObject() {
	this.object = nil
}

func (this *Value) SetObject(val Indexable) {
	this.object = val
	this.data = strconv.FormatInt(val.Len(), 10)
}

func (this *Value) HasObject() bool {
	return this.object != nil
}

func (this *Value) Get(index string) *Value {
	if this.object == nil {
		panic("Attempt to use object function on a *Value that has no object data.")
	}
	
	return this.object.Get(index)
}

func (this *Value) Set(index string, value *Value) {
	if this.object == nil {
		panic("Attempt to use object function on a *Value that has no object data.")
	}
	this.object.Set(index, value)
	this.data = strconv.FormatInt(this.Len(), 10)
}

func (this *Value) Exists(index string) bool {
	if this.object == nil {
		panic("Attempt to use object function on a *Value that has no object data.")
	}
	
	return this.object.Exists(index)
}

func (this *Value) Len() int64 {
	if this.object == nil {
		panic("Attempt to use object function on a *Value that has no object data.")
	}
	
	return this.object.Len()
}

func (this *Value) Keys() []string {
	if this.object == nil {
		panic("Attempt to use object function on a *Value that has no object data.")
	}
	
	return this.object.Keys()
}
