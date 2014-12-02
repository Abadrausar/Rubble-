package nca6

import "fmt"
import "os"
import "strconv"
import "dctech/ncalex"
import "text/tabwriter"

const (
	TypString = iota
	TypInt
	TypFloat
	TypObject
)

// Indexable represents a ScriptObject that may be used with the 
// dereference opperator ([]) as well as some of the base commands.
type Indexable interface {
	Get(string) *Value
	Set(string, *Value)
	Exists(string) bool
	Len() int64
	Keys() []string
}

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

// NewValueObject creates a new Value from a ScriptObject
func NewValueObject(val ScriptObject) *Value {
	this := new(Value)
	this.Type = TypObject
	this.Data = val
	this.Line = -1
	this.Column = -1
	return this
}

// TokenToValue turns a lexer token into a script value.
// If the lexeme can be converted to int the value is TypInt else it is TypString.
func TokenToValue(tok *ncalex.Token) *Value {
	this := new(Value)
	
	val, err := strconv.ParseInt(tok.Lexeme, 0, 64)
	if err != nil {
		this.Type = TypString
		this.Data = tok.Lexeme
	} else {
		this.Type = TypInt
		this.Data = val
	}
	this.Line = tok.Line
	this.Column = tok.Column
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
	
	case TypObject:
		return "<ScriptObject>"
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
	
	case TypObject:
		return 0
	}
	panic("Script Value has invalid Type.")
}

// Bool converts a Value to a bool.
// Strings return true for: 1, t, T, TRUE, true, True and false for anything else
// Ints are true if the value is anything other than 0.
// Floats are converted to ints and then converted to bool by the int rules.
// Objects are always false.
func (this *Value) Bool() bool {
	switch this.Type {
	case TypString:
		val, err := strconv.ParseBool(this.Data.(string))
		if err != nil {
			val = false
		}
		return val
	
	case TypInt:
		return this.Data.(int64) != 0
	
	case TypFloat:
		return int64(this.Data.(float64)) != 0
	
	case TypObject:
		return false
	}
	panic("Script Value has invalid Type.")
}

// AsLexer converts a Value to a *ncalex.Lexer.
// This will probably be removed.
func (this *Value) AsLexer() *ncalex.Lexer {
	if this.Type != TypString {
		panic("Only strings may be used as lexers.")
	}
	return ncalex.NewLexer(this.Data.(string), this.Line, this.Column)
}

// Object handling

func (this *Value) IsIndexable() bool {
	_, ok := this.Data.(Indexable)
	return ok
}

// CommandValueInspect returns everything you ever wanted to know about a script value.
// Mostly for the use of NCASH and debugging.
//	// Register via: 
//	state.NewNativeCommand("valueinspect", nca6.CommandValueInspect)
func CommandValueInspect(state *State, params []*Value) {
	if len(params) != 1 {
		panic("Wrong number of params to valueinspect.")
	}
	
	fmt.Println("=============================================")
	fmt.Println("NCA Script Value Inspector")
	fmt.Println("Data:", params[0].Data)
	switch params[0].Type {
	case TypString:
		fmt.Println("Type: TypString")
	case TypInt:
		fmt.Println("Type: TypInt")
	case TypFloat:
		fmt.Println("Type: TypFloat")
	case TypObject:
		fmt.Println("Type: TypObject")
	}
	fmt.Println("Line: ", params[0].Line)
	fmt.Println("Column: ", params[0].Column)
	if params[0].IsIndexable() {
		data := params[0].Data.(Indexable)
		fmt.Println("Value Is Indexable")
		fmt.Println("len: ", data.Len())
		
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintf(w, "%v\t%v\n", "Key", "Value")
		for _, key := range data.Keys() {
			val := data.Get(key)
			fmt.Fprintf(w, "%v\t%v\n", key, val.String())
		}
		w.Flush()
	}
	fmt.Println("=============================================")
}





