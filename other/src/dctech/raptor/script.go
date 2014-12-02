/*
For copyright/license see header in file "doc.go"
*/

package raptor

import "fmt"
import "io"
import "os"
import "errors"
import "strings"

// Script stores all non-global data and provides an interface to the global data.
// Script.Host is only valid if the script is being run by a State!
type Script struct {
	RetVal    *Value      // The return value of the last command
	Exit      bool        // true when exiting
	Return    bool        // true when a return is active
	Break     bool        // true when a break is active
	BreakLoop bool        // true when a loop break is active
	Error     bool        // set by some commands on error, this is NOT automatically reset!
	This      *Value      // When the value retrieved by the indexing deref operator is a command this is set to the containing Indexable.
	Envs      *EnvStore   // The script environments, you should not need to touch this.
	Code      *BlockStore // This is where script code is stored.
	Output    io.Writer   // Normally set to nil (in which case the global value is used)
	Host      *State      // Set to nil unless this Script is being run, then it is the host State.

	// Set to "" unless a _NATIVE_ command is running, then it is the string command name.
	// Note that if another native command is called by a native command (eg in loop commands) this value can
	// get trashed (for the first command). This could be a stack, but in practice it's not needed,
	// just make sure you don't use it in a command after calling another command.
	// (this means don't use the error generators after you start such a loop)
	RunningCmd string
}

// NewScript creates (and initializes) a new state.
func NewScript() *Script {
	rtn := new(Script)
	rtn.Envs = NewEnvStore()
	rtn.Code = NewBlockStore()
	rtn.Output = nil
	rtn.This = NewValue()
	rtn.RetVal = NewValue()
	return rtn
}

// Variables

// FetchEnv gets a pointer to the env containing var, panics if var does not exist.
// This should only be called after a failed call to ParseName.
func (this *Script) FetchEnv(name string) *Environment {
	for i := len(*this.Envs) - 1; i >= 0; i-- {
		if _, ok := (*this.Envs)[i].Vars[name]; ok {
			return (*this.Envs)[i]
		}
	}
	panic("Undeclared variable: " + name)
}

// NewVar creates a new script variable.
func (this *Script) NewVar(name string, value *Value) {
	space, itemname := this.ParseName(name)
	if space == nil {
		if _, ok := this.Envs.Last().Vars[itemname]; ok {
			panic(fmt.Sprintf("Variable: \"%v\" already declared", name))
		}
		this.Envs.Last().Vars[itemname] = value
		return
	}
	if space.Vars.Exist(itemname) {
		panic(fmt.Sprintf("Variable: \"%v\" already declared", name))
	}
	space.Vars.Store(itemname, value)
}

// DeleteVar deletes a variable, but only if it's in the last environment.
func (this *Script) DeleteVar(name string) *Value {
	space, itemname := this.ParseName(name)
	if space == nil {
		if _, ok := this.Envs.Last().Vars[itemname]; !ok {
			panic(fmt.Sprintf("Variable: \"%v\" not declared in the current environment", name))
		}
		rtn := this.Envs.Last().Vars[itemname]
		delete(this.Envs.Last().Vars, itemname)
		return rtn
	}
	if !space.Vars.Exist(itemname) {
		panic(fmt.Sprintf("Variable: \"%v\" not declared", name))
	}
	rtn := space.Vars.Fetch(itemname)
	space.Vars.Delete(itemname)
	return rtn
}

// GetValue gets the value of a variable.
func (this *Script) GetValue(name string) *Value {
	space, itemname := this.ParseName(name)
	if space == nil {
		return this.FetchEnv(itemname).Vars[itemname]
	}

	if space.Vars.Exist(itemname) {
		return space.Vars.Fetch(itemname)
	}
	panic("Undeclared variable: " + name)
}

// SetValue sets the value of a variable.
func (this *Script) SetValue(name string, value *Value) {
	space, itemname := this.ParseName(name)
	if space == nil {
		this.FetchEnv(itemname).Vars[itemname] = value
		return
	}

	if space.Vars.Exist(itemname) {
		space.Vars.Store(itemname, value)
		return
	}
	panic("Undeclared variable: " + name)
}

// VarExists returns true if variable "name" exists.
func (this *Script) VarExists(name string) bool {
	space, itemname := this.ParseName(name)
	if space == nil {
		for i := len(*this.Envs) - 1; i >= 0; i-- {
			if _, ok := (*this.Envs)[i].Vars[itemname]; ok {
				return true
			}
		}
		return false
	}

	return space.Vars.Exist(itemname)
}

// AddParams creates the special "params" array using strings.
func (this *Script) AddParams(params ...string) {
	array := make([]*Value, len(params))
	for i, val := range params {
		array[i] = NewValueString(val)
	}
	this.NewVar("params", NewValueObject(NewParamsArray(array)))
}

// AddParamsValue creates the special "params" array using values.
func (this *Script) AddParamsValue(params ...*Value) {
	array := make([]*Value, len(params))
	for i, val := range params {
		array[i] = val
	}
	this.NewVar("params", NewValueObject(NewParamsArray(array)))
}

// Namespaces

// NewNameSpace creates a new namespace.
// May only be called with a valid Host.
func (this *Script) NewNameSpace(name string) {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	this.Host.NewNameSpace(name)
}

// DeleteNameSpace deletes a namespace.
// May only be called with a valid Host.
func (this *Script) DeleteNameSpace(name string) {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	this.Host.DeleteNameSpace(name)
}

// Commands

// NewNativeCommand adds a new native command.
// May only be called with a valid Host.
func (this *Script) NewNativeCommand(name string, handler NativeCommand) {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	this.Host.NewNativeCommand(name, handler)
}

// NewUserCommand adds a new user command (what else would it do?).
// May only be called with a valid Host.
func (this *Script) NewUserCommand(name string, code *Value, params []*Value) {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	this.Host.NewUserCommand(name, code, params)
}

// GetCommand fetches a command by it's name.
// May only be called with a valid Host.
func (this *Script) GetCommand(name string) *Command {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	return this.Host.GetCommand(name)
}

// DeleteCommand removes a command.
// May only be called with a valid Host.
func (this *Script) DeleteCommand(name string) {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	this.Host.DeleteCommand(name)
}

// Native Command Parameter Checking

// GeneralCmdError generates an error message for general command errors.
//	"test:cmd": some error message
func (this *Script) GeneralCmdError(msg string) string {
	return "\"" + this.RunningCmd + "\": " + msg
}

// BadParamCount generates an error message for incorrect parameter count errors.
//	"test:cmd": Incorrect parameter count: 1 parameter(s) required.
// If hint is set to "" it leaves that part off.
//	"test:cmd": Incorrect parameter count.
// Hint should be something like "5" or ">2".
func (this *Script) BadParamCount(hint string) string {
	if hint != "" {
		return this.GeneralCmdError("Incorrect parameter count: " + hint + " parameter(s) required.")
	}
	return this.GeneralCmdError("Incorrect parameter count.")
}

// Indexables

// RegisterType registers a new Indexable or the like with a specific name.
// Registering a type allows it to be created with an object literal.
// May only be called with a valid Host.
func (this *Script) RegisterType(name string, handler ObjectFactory) {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	this.Host.RegisterType(name, handler)
}

// GetType retrieves a named types ObjectFactory.
// May only be called with a valid Host.
func (this *Script) GetType(name string) ObjectFactory {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	return this.Host.GetType(name)
}

// Name Parsing

// ParseName parses a name and returns a namespace or nil and the base name.
// Eg. "test:object" will parse to a pointer to the namespace "test", and the string "object".
func (this *Script) ParseName(name string) (*NameSpace, string) {
	if !strings.Contains(name, ":") {
		return nil, name
	}
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	return this.Host.ParseName(name)
}

// ParseNameSpaceName is just like ParseName but for just the name of a namespace.
// May only be called with a valid Host.
func (this *Script) ParseNameSpaceName(name string) *NameSpace {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	return this.Host.ParseNameSpaceName(name)
}

// Output

// Printf writes to the script output if set, else writes to the host output if the host is set, else defaults to Stdout.
func (this *Script) Printf(format string, msg ...interface{}) {
	if this.Output == nil {
		if this.Host == nil {
			fmt.Fprintf(os.Stdout, format, msg...)
			return
		}

		fmt.Fprintf(this.Host.Output, format, msg...)
		return
	}
	fmt.Fprintf(this.Output, format, msg...)
}

// Println writes to the script output if set, else writes to the host output if the host is set, else defaults to Stdout.
func (this *Script) Println(msg ...interface{}) {
	if this.Output == nil {
		if this.Host == nil {
			fmt.Fprintln(os.Stdout, msg...)
			return
		}

		fmt.Fprintln(this.Host.Output, msg...)
		return
	}
	fmt.Fprintln(this.Output, msg...)
}

// Print writes to the script output if set, else writes to the host output if the host is set, else defaults to Stdout.
func (this *Script) Print(msg ...interface{}) {
	if this.Output == nil {
		if this.Host == nil {
			fmt.Fprint(os.Stdout, msg...)
			return
		}

		fmt.Fprint(this.Host.Output, msg...)
		return
	}
	fmt.Fprint(this.Output, msg...)
}

// Other

// ExitFlags returns true if any of the exit flags are set.
// The exit flags are Exit, Return, Break, and BreakLoop.
func (this *Script) ExitFlags() bool {
	return (this.Exit || this.Return || this.Break || this.BreakLoop)
}

// Clears the script.
func (this *Script) Clear() {
	this.ClearMost()
	this.RetVal = NewValue()
}

// Clears the script.
// Does not modify the RetVal register.
func (this *Script) ClearMost() {
	this.Code.Clear()
	this.Envs.Clear()
	this.Host = nil
	this.This = NewValue()
	this.Exit = false
	this.Return = false
	this.Break = false
	this.BreakLoop = false
	this.RunningCmd = ""
}

// Clears the script.
// Does not modify the RetVal register.
// Leaves the "root" environment.
func (this *Script) ClearSome() {
	this.Code.Clear()
	this.Envs.ClearAllButRoot()
	this.Host = nil
	this.This = NewValue()
	this.Exit = false
	this.Return = false
	this.Break = false
	this.BreakLoop = false
	this.RunningCmd = ""
}

// Script Execution

// Exec is exported for the use of commands ONLY!
// Exec is a subset of Run() and so it must be called from a command handler or the like.
func (this *Script) Exec() {
	for !CheckLookAhead(this.Code.Last(), TknINVALID) {
		this.RetVal = this.fetchValue()
		if this.ExitFlags() {
			this.Break = false
			this.Code.Remove()
			return
		}
	}
	this.Code.Remove()
	return
}

// SafeExec is exported for the use of commands ONLY!
// SafeExec traps panics and turns them into errors, use when you want to emulate the behavior
// of State.Run without doing the other stuff Run does.
// Respects this.Host.NoRecover.
// Halts all exit states, but does none of the other cleanup (aside from what Exec already does).
func (this *Script) SafeExec() (err error) {
	err = nil

	defer func() {
		this.Exit = false
		this.Return = false
		this.Break = false
		this.BreakLoop = false

		if this.Host.NoRecover {
			return
		}

		if x := recover(); x != nil {
			switch i := x.(type) {
			case error:
				err = i
			case string:
				err = fmt.Errorf("%v Near %v", i, this.Code.Last().CurrentTkn().Pos)
			default:
				err = errors.New(fmt.Sprint(i))
			}
		}
	}()

	this.Exec()
	return // Required :(
}

func (this *Script) execCommand() {
	GetToken(this.Code.Last(), TknCmdBegin)

	// Get the command's name
	name := this.fetchValue()
	if this.ExitFlags() {
		return
	}

	// Read the commands parameters if any
	params := make([]*Value, 0, 5)
	for !CheckLookAhead(this.Code.Last(), TknCmdEnd) {
		params = append(params, this.fetchValue())
		if this.ExitFlags() {
			return
		}
	}

	GetToken(this.Code.Last(), TknCmdEnd)

	this.GetCommand(name.String()).Call(name.String(), this, params)
}

func (this *Script) derefVar() {
	GetToken(this.Code.Last(), TknDerefBegin)

	// Get the name or value to opperate on
	name := this.fetchValue()
	if this.ExitFlags() {
		return
	}

	// Read any index parameters
	params := make([]*Value, 0, 5)
	for !CheckLookAhead(this.Code.Last(), TknDerefEnd) {
		params = append(params, this.fetchValue())
		if this.ExitFlags() {
			return
		}
	}

	GetToken(this.Code.Last(), TknDerefEnd)

	// simple name deref
	if len(params) == 0 {
		this.RetVal = this.GetValue(name.String())
		return
	}

	val := name
	if name.Type != TypObject {
		val = this.GetValue(name.String())
	}

	last := val
	for i := range params {
		if val.Type != TypObject {
			panic(fmt.Sprint("Non-Object value passed to an indexing deref opperator. Indexing level:", i))
		}
		obj := val.Indexable()
		if obj == nil {
			panic(fmt.Sprint("Non-Indexable object passed to an indexing deref opperator. Indexing level:", i))
		}

		last = val
		val = obj.Get(params[i].String())
	}

	// "this" should work (pardon the pun)
	if val.Type == TypCommand {
		this.This = last
	}
	this.RetVal = val
	return
}

func (this *Script) readObjLit() {
	GetToken(this.Code.Last(), TknObjLitBegin)

	// Get the type name
	name := this.fetchValue().String()
	if this.ExitFlags() {
		return
	}

	// Generate key/value lists
	// you must have keys for all or none
	values := make([]*Value, 0, 10)
	keys := make([]string, 0, 10)
	hasKeys := false

	if !CheckLookAhead(this.Code.Last(), TknObjLitEnd) {
		ret := this.fetchValue()
		if this.ExitFlags() {
			return
		}
		if CheckLookAhead(this.Code.Last(), TknObjLitSplit) {
			hasKeys = true
			keys = append(keys, ret.String())
			GetToken(this.Code.Last(), TknObjLitSplit)
			values = append(values, this.fetchValue())
			if this.ExitFlags() {
				return
			}
		} else {
			values = append(values, ret)
		}
	}
	for !CheckLookAhead(this.Code.Last(), TknObjLitEnd) {
		ret := this.fetchValue()
		if this.ExitFlags() {
			return
		}
		if hasKeys {
			keys = append(keys, ret.String())
			GetToken(this.Code.Last(), TknObjLitSplit)
			values = append(values, this.fetchValue())
			if this.ExitFlags() {
				return
			}
		} else {
			values = append(values, ret)
		}
	}

	GetToken(this.Code.Last(), TknObjLitEnd)

	if hasKeys {
		this.RetVal = this.GetType(name)(this, keys, values)
		return
	}
	this.RetVal = this.GetType(name)(this, nil, values)
	return
}

func (this *Script) readCodeBlock() {
	GetToken(this.Code.Last(), TknCodeBegin)
	this.RetVal = NewValueCode(NewCodeBlock(this.Code.Last()))
}

func (this *Script) fetchValue() *Value {
	switch this.Code.Last().LookAhead().Type {
	case TknString:
		GetToken(this.Code.Last(), TknString)
		return TokenToValue(this.Code.Last().CurrentTkn())

	case TknCmdBegin:
		this.execCommand()
		return this.RetVal

	case TknDerefBegin:
		this.derefVar()
		return this.RetVal

	case TknObjLitBegin:
		this.readObjLit()
		return this.RetVal

	case TknCodeBegin:
		this.readCodeBlock()
		return this.RetVal

	default:
		ExitOnTokenExpected(this.Code.Last().LookAhead(), TknString, TknCmdBegin,
			TknDerefBegin, TknObjLitBegin, TknCodeBegin)
	}
	panic("UNREACHABLE")
}

// Script Validation

// Validate checks a script to make sure all braces are matched and things like derefs are correctly formed.
// Validate CANNOT ensure a script is valid! Validate will only find obvious syntax errors.
// The check is "destructive" eg the code is consumed.
// If validate is called from the state is will run further checks for commands and object literals
// (Note that this may result in false missing command errors).
func (this *Script) Validate() (err error) {
	err = nil

	defer func() {
		if x := recover(); x != nil {
			switch i := x.(type) {
			case error:
				err = i
			case string:
				err = fmt.Errorf("%v Near %v", i, this.Code.Last().CurrentTkn().Pos)
			default:
				err = errors.New(fmt.Sprint(i))
			}
		}

		this.Code.Clear()
	}()

	this.validate()
	return
}

func (this *Script) validate() {
	for !CheckLookAhead(this.Code.Last(), TknINVALID) {
		this.RetVal = this.validateValue()
	}
	this.Code.Remove()
	return
}

func (this *Script) validateCommand() {
	GetToken(this.Code.Last(), TknCmdBegin)

	// Get the command's name
	name := this.validateValue()
	paramcount := -1
	params := 0

	if this.Host != nil {
		// Check params and command existence
		command := this.GetCommand(name.String())
		if command.Params != nil || command.Handler != nil {
			paramcount = len(command.Params)
		}
	}

	// Read the commands parameters if any
	for !CheckLookAhead(this.Code.Last(), TknCmdEnd) {
		this.validateValue()
		params++
	}

	if paramcount != -1 {
		if paramcount != params {
			panic("Invalid param count to " + name.String() + ".")
		}
	}

	GetToken(this.Code.Last(), TknCmdEnd)
}

func (this *Script) validateDeref() {
	GetToken(this.Code.Last(), TknDerefBegin)

	this.validateValue()

	// Read any index parameters
	for !CheckLookAhead(this.Code.Last(), TknDerefEnd) {
		this.validateValue()
	}

	GetToken(this.Code.Last(), TknDerefEnd)
}

func (this *Script) validateObjLit() {
	GetToken(this.Code.Last(), TknObjLitBegin)

	name := this.validateValue()

	if this.Host != nil {
		this.GetType(name.String())
	}

	for !CheckLookAhead(this.Code.Last(), TknObjLitEnd) {
		this.validateValue()
		if CheckLookAhead(this.Code.Last(), TknObjLitSplit) {
			GetToken(this.Code.Last(), TknObjLitSplit)
			this.validateValue()
		}
	}

	GetToken(this.Code.Last(), TknObjLitEnd)
}

func (this *Script) validateCodeBlock() {
	GetToken(this.Code.Last(), TknCodeBegin)
	this.Code.AddCode(NewCodeBlock(this.Code.Last()))
	this.validate()
}

func (this *Script) validateValue() *Value {
	switch this.Code.Last().LookAhead().Type {
	case TknString:
		GetToken(this.Code.Last(), TknString)
		return TokenToValue(this.Code.Last().CurrentTkn())

	case TknCmdBegin:
		this.validateCommand()
		return NewValue()

	case TknDerefBegin:
		this.validateDeref()
		return NewValue()

	case TknObjLitBegin:
		this.validateObjLit()
		return NewValue()

	case TknCodeBegin:
		this.validateCodeBlock()
		return NewValue()

	default:
		ExitOnTokenExpected(this.Code.Last().LookAhead(), TknString, TknCmdBegin,
			TknDerefBegin, TknObjLitBegin, TknCodeBegin)
	}
	panic("UNREACHABLE")
}
