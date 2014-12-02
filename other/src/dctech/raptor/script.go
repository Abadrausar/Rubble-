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

package raptor

import "fmt"
import "io"
import "errors"
import "strings"

// Script stores all non-global data and provides an interface to the global data.
// Script.Host is only valid if the script is being run by a State!
type Script struct {
	RetVal     *Value       // The return value of the last command
	Exit       bool         // true when exiting
	Return     bool         // true when a return is active
	Break      bool         // true when a break is active
	BreakLoop  bool         // true when a loop break is active
	Error      bool         // set by some commands on error, this is NOT automatically reset!
	This       *Value       // When the value retrieved by the indexing deref operator is a command this is set to the containing Indexable.
	Envs       *EnvStore    // The script environments, you should not need to touch this.
	Code       *BlockStore  // This is where script code is stored.
	Output     io.Writer    // Normally set to nil (in which case the global value is used)
	Host       *State
}

// NewScript creates (and initializes) a new state.
func NewScript() *Script {
	rtn := new(Script)
	rtn.Envs = NewEnvStore()
	rtn.Code = NewBlockStore()
	rtn.Output = nil
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
func (this *Script) NewNameSpace(name string) {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	this.Host.NewNameSpace(name)
}

// DeleteNameSpace deletes a namespace.
func (this *Script) DeleteNameSpace(name string) {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	this.Host.DeleteNameSpace(name)
}

// Commands

// NewNativeCommand adds a new native command.
func (this *Script) NewNativeCommand(name string, handler NativeCommand) {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	this.Host.NewNativeCommand(name, handler)
}

// NewUserCommand adds a new user command (what else would it do?).
func (this *Script) NewUserCommand(name string, code *Value, params []*Value) {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	this.Host.NewUserCommand(name, code, params)
}

// GetCommand fetches a command by it's name.
func (this *Script) GetCommand(name string) *Command {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	return this.Host.GetCommand(name)
}

// DeleteCommand removes a command.
func (this *Script) DeleteCommand(name string) {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	this.Host.DeleteCommand(name)
}

// Indexables

// RegisterType registers a new Indexable or the like with a specific name.
// Registering a type allows it to be created with an object literal.
func (this *Script) RegisterType(name string, handler ObjectFactory) {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	this.Host.RegisterType(name, handler)
}

// GetType retrieves a named types ObjectFactory.
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
func (this *Script) ParseNameSpaceName(name string) *NameSpace {
	if this.Host == nil {
		panic("Script.Host == nil: This function should not be called here!")
	}
	return this.Host.ParseNameSpaceName(name)
}

// Output

func (this *Script) Printf(format string, msg ...interface{}) {
	if this.Output == nil {
		fmt.Fprintf(this.Host.Output, format, msg...)
		return
	}
	fmt.Fprintf(this.Output, format, msg...)
}

func (this *Script) Println(msg ...interface{}) {
	if this.Output == nil {
		fmt.Fprintln(this.Host.Output, msg...)
		return
	}
	fmt.Fprintln(this.Output, msg...)
}

func (this *Script) Print(msg ...interface{}) {
	if this.Output == nil {
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

// Script Execution

// Exec is exported for the use of commands ONLY!
// Exec is a subset of Run() and so it must be called from a command handler or the like.
func (this *Script) Exec() {
	for !this.Code.Last().CheckLookAhead(TknINVALID) {
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
				err = fmt.Errorf("%v Near %v", i, this.Code.Last().Position())
			default:
				err = errors.New(fmt.Sprint(i))
			}
		}
	}()

	this.Exec()
	return // Required :(
}

func (this *Script) execCommand() {
	this.Host.DbgCallback(DbgrAdvanceTkn)
	this.Code.Last().GetToken(TknCmdBegin)

	// Get the command's name
	name := this.fetchValue()
	if this.ExitFlags() {
		return
	}

	// Read the commands parameters if any
	params := make([]*Value, 0, 5)
	for !this.Code.Last().CheckLookAhead(TknCmdEnd) {
		params = append(params, this.fetchValue())
		if this.ExitFlags() {
			return
		}
	}

	this.Host.DbgCallback(DbgrAdvanceTkn)
	this.Code.Last().GetToken(TknCmdEnd)

	this.GetCommand(name.String()).Call(this, params)
}

func (this *Script) derefVar() {
	this.Host.DbgCallback(DbgrAdvanceTkn)
	this.Code.Last().GetToken(TknDerefBegin)

	// Get the name or value to opperate on
	name := this.fetchValue()
	if this.ExitFlags() {
		return
	}

	// Read any index parameters
	params := make([]*Value, 0, 5)
	for !this.Code.Last().CheckLookAhead(TknDerefEnd) {
		params = append(params, this.fetchValue())
		if this.ExitFlags() {
			return
		}
	}

	this.Host.DbgCallback(DbgrAdvanceTkn)
	this.Code.Last().GetToken(TknDerefEnd)

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
	this.Host.DbgCallback(DbgrAdvanceTkn)
	this.Code.Last().GetToken(TknObjLitBegin)

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

	if !this.Code.Last().CheckLookAhead(TknObjLitEnd) {
		ret := this.fetchValue()
		if this.ExitFlags() {
			return
		}
		if this.Code.Last().CheckLookAhead(TknObjLitSplit) {
			hasKeys = true
			keys = append(keys, ret.String())
			this.Host.DbgCallback(DbgrAdvanceTkn)
			this.Code.Last().GetToken(TknObjLitSplit)
			values = append(values, this.fetchValue())
			if this.ExitFlags() {
				return
			}
		} else {
			values = append(values, ret)
		}
	}
	for !this.Code.Last().CheckLookAhead(TknObjLitEnd) {
		ret := this.fetchValue()
		if this.ExitFlags() {
			return
		}
		if hasKeys {
			keys = append(keys, ret.String())
			this.Host.DbgCallback(DbgrAdvanceTkn)
			this.Code.Last().GetToken(TknObjLitSplit)
			values = append(values, this.fetchValue())
			if this.ExitFlags() {
				return
			}
		} else {
			values = append(values, ret)
		}
	}

	this.Host.DbgCallback(DbgrAdvanceTkn)
	this.Code.Last().GetToken(TknObjLitEnd)

	if hasKeys {
		this.RetVal = this.GetType(name)(this, keys, values)
		return
	}
	this.RetVal = this.GetType(name)(this, nil, values)
	return
}

func (this *Script) readCodeBlock() {
	this.Host.DbgCallback(DbgrAdvanceTkn)
	this.Code.Last().GetToken(TknCodeBegin)

	this.RetVal = NewValueCode(CompileBlock(this.Code.Last()))
}

func (this *Script) fetchValue() *Value {
	switch this.Code.Last().LookAhead().Type {
	case TknString:
		this.Host.DbgCallback(DbgrAdvanceTkn)
		this.Code.Last().GetToken(TknString)
		return TokenToValue(this.Code.Last().CurrentTkn())

	case TknCmdBegin:
		this.Host.DbgCallback(DbgrEnterCmd)
		this.execCommand()
		this.Host.DbgCallback(DbgrLeaveCmd)
		return this.RetVal

	case TknDerefBegin:
		this.Host.DbgCallback(DbgrEnterDeref)
		this.derefVar()
		this.Host.DbgCallback(DbgrLeaveDeref)
		return this.RetVal

	case TknObjLitBegin:
		this.Host.DbgCallback(DbgrEnterObjLit)
		this.readObjLit()
		this.Host.DbgCallback(DbgrLeaveObjLit)
		return this.RetVal

	case TknCodeBegin:
		this.Host.DbgCallback(DbgrEnterCodeBlock)
		this.readCodeBlock()
		this.Host.DbgCallback(DbgrLeaveCodeBlock)
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
func (this *Script) Validate() (err error) {
	err = nil

	defer func() {
		if x := recover(); x != nil {
			switch i := x.(type) {
			case error:
				err = i
			case string:
				err = fmt.Errorf("%v Near %v", i, this.Code.Last().Position())
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
	for !this.Code.Last().CheckLookAhead(TknINVALID) {
		this.RetVal = this.validateValue()
	}
	this.Code.Remove()
	return
}

func (this *Script) validateCommand() {
	this.Code.Last().GetToken(TknCmdBegin)

	// Get the command's name
	this.validateValue()

	// Read the commands parameters if any
	for !this.Code.Last().CheckLookAhead(TknCmdEnd) {
		this.validateValue()
	}
	
	this.Code.Last().GetToken(TknCmdEnd)
}

func (this *Script) validateDeref() {
	this.Code.Last().GetToken(TknDerefBegin)

	this.validateValue()

	// Read any index parameters
	for !this.Code.Last().CheckLookAhead(TknDerefEnd) {
		this.validateValue()
	}

	this.Code.Last().GetToken(TknDerefEnd)
}

func (this *Script) validateObjLit() {
	this.Code.Last().GetToken(TknObjLitBegin)

	this.validateValue()

	for !this.Code.Last().CheckLookAhead(TknObjLitEnd) {
		this.validateValue()
		if this.Code.Last().CheckLookAhead(TknObjLitSplit) {
			this.Code.Last().GetToken(TknObjLitSplit)
			this.validateValue()
		}
	}
	
	this.Code.Last().GetToken(TknObjLitEnd)
}

func (this *Script) validateCodeBlock() {
	this.Code.Last().GetToken(TknCodeBegin)
	this.Code.AddCompiledScript(CompileBlock(this.Code.Last()))
	this.validate()
}

func (this *Script) validateValue() *Value {
	switch this.Code.Last().LookAhead().Type {
	case TknString:
		this.Code.Last().GetToken(TknString)
		return TokenToValue(this.Code.Last().CurrentTkn())

	case TknCmdBegin:
		this.validateCommand()
		return this.RetVal

	case TknDerefBegin:
		this.validateDeref()
		return this.RetVal

	case TknObjLitBegin:
		this.validateObjLit()
		return this.RetVal

	case TknCodeBegin:
		this.validateCodeBlock()
		return this.RetVal

	default:
		ExitOnTokenExpected(this.Code.Last().LookAhead(), TknString, TknCmdBegin,
			TknDerefBegin, TknObjLitBegin, TknCodeBegin)
	}
	panic("UNREACHABLE")
}




