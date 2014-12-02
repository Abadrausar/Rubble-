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

import "fmt"
import "regexp"
import "strings"
import "errors"

// State handles EVERYTHING to do with nca scripts, including running them.
// The majority of the fields are exported for the use of commands only.
type State struct {
	NameSpaces map[string]*NameSpace
	Commands   map[string]*Command
	Types      map[string]IndexableFactory
	RetVal     *Value // The return value of the last command
	Return     bool // true when a return is active
	Break      bool // true when a break is active
	Exit       bool // true when exiting
	Error      bool // set by some commands on error, this is NOT automaticly reset!
	NoRecover  bool // Do not recover errors, this makes it easier to debug the parser.
	Envs       *EnvStore // The script environments, you should not need to touch this.
	Code       *BlockStore // This is where script code is stored.
	Debug      []DbgHandler // This is the debug handlers, normaly nil unless a debugger is installed.
}

// NewState creates (and initalizes) a new state.
func NewState() *State {
	rtn := new(State)
	rtn.NameSpaces = make(map[string]*NameSpace)
	rtn.Commands = make(map[string]*Command)
	rtn.Types = make(map[string]IndexableFactory)
	rtn.Envs = NewEnvStore()
	rtn.Envs.Add(NewEnvironment())
	rtn.Code = NewBlockStore()
	rtn.Debug = nil
	return rtn
}

// Variables

// NewVar creates a new script variable.
func (this *State) NewVar(name string, value *Value) {
	space, itemname := this.ParseName(name)
	if space == nil {
		if _, ok := this.Envs.Last().Vars[itemname]; ok {
			panic(fmt.Sprintf("Variable: \"%v\" already declared", name))
		}
		this.Envs.Last().Vars[itemname] = value
		return
	}
	if _, ok := space.Vars[itemname]; ok {
		panic(fmt.Sprintf("Variable: \"%v\" already declared", name))
	}
	space.Vars[itemname] = value
}

// DeleteVar deletes a variable, but only if it's in the last environment.
func (this *State) DeleteVar(name string) *Value {
	space, itemname := this.ParseName(name)
	if space == nil {
		if _, ok := this.Envs.Last().Vars[itemname]; !ok {
			panic(fmt.Sprintf("Variable: \"%v\" not declared in the current environment", name))
		}
		rtn := this.Envs.Last().Vars[itemname]
		delete(this.Envs.Last().Vars, itemname)
		return rtn
	}
	if _, ok := space.Vars[itemname]; !ok {
		panic(fmt.Sprintf("Variable: \"%v\" not declared", name))
	}
	rtn := space.Vars[itemname]
	delete(space.Vars, itemname)
	return rtn
}

// GetValue gets the value of a variable.
func (this *State) GetValue(name string) *Value {
	space, itemname := this.ParseName(name)
	if space == nil {
		for i := len(*this.Envs) - 1; i >= 0; i-- {
			if val, ok := (*this.Envs)[i].Vars[itemname]; ok {
				return val
			}
		}
		panic(fmt.Sprintf("Undeclared variable: \"%v\"", name))
	}

	if val, ok := space.Vars[itemname]; ok {
		return val
	}
	panic(fmt.Sprintf("Undeclared variable: \"%v\"", name))
}

// SetValue sets the value of a variable.
func (this *State) SetValue(name string, value *Value) {
	space, itemname := this.ParseName(name)
	if space == nil {
		for i := len(*this.Envs) - 1; i >= 0; i-- {
			if _, ok := (*this.Envs)[i].Vars[itemname]; ok {
				(*this.Envs)[i].Vars[itemname] = value
				return
			}
		}
		panic(fmt.Sprintf("Undeclared variable: \"%v\"", name))
	}

	if _, ok := space.Vars[itemname]; ok {
		space.Vars[itemname] = value
		return
	}
	panic(fmt.Sprintf("Undeclared variable: \"%v\"", name))
}

// VarExists returns true if variable "name" exists.
func (this *State) VarExists(name string) bool {
	space, itemname := this.ParseName(name)
	if space == nil {
		for i := len(*this.Envs) - 1; i >= 0; i-- {
			if _, ok := (*this.Envs)[i].Vars[itemname]; ok {
				return true
			}
		}
		return false
	}

	if _, ok := space.Vars[itemname]; ok {
		return true
	}
	return false
}

// AddParams creates the special "params" array using strings.
func (this *State) AddParams(params ...string) {
	array := make([]*Value, len(params))
	for i, val := range params {
		array[i] = NewValueString(val)
	}
	this.NewVar("params", NewValueObject(NewParamsArray(array)))
}

// AddParamsValue creates the special "params" array using values.
func (this *State) AddParamsValue(params ...*Value) {
	array := make([]*Value, len(params))
	for i, val := range params {
		array[i] = val
	}
	this.NewVar("params", NewValueObject(NewParamsArray(array)))
}

var nameSplitRegex = regexp.MustCompile("^(([^:]*:)*)([^:]*)$")

// ParseName parses a name and returns a namespace or nil and the base name.
// Eg. "test:code" will parse to a pointer to the namespace "test" and "code".
func (this *State) ParseName(name string) (*NameSpace, string) {
	if !strings.Contains(name, ":") {
		return nil, name
	}

	namesplit := nameSplitRegex.FindStringSubmatch(name)
	namespacelist := strings.TrimRight(namesplit[1], ":")
	basename := namesplit[3]
	namespace := this.ParseNameSpaceName(namespacelist)
	return namespace, basename
}

// Namespaces

// NewNameSpace creates a new namespace.
func (this *State) NewNameSpace(name string) {
	if strings.Contains(name, ":") {
		namesplit := nameSplitRegex.FindStringSubmatch(name)
		namespacelist := strings.TrimRight(namesplit[1], ":")
		basename := namesplit[3]
		namespace := this.ParseNameSpaceName(namespacelist)
		if _, ok := namespace.NameSpaces[basename]; ok {
			panic(fmt.Sprintf("Namespace: \"%v\" already declared", name))
		}
		namespace.NameSpaces[basename] = NewNameSpace()
		return
	}

	if _, ok := this.NameSpaces[name]; ok {
		panic(fmt.Sprintf("Namespace: \"%v\" already declared", name))
	}
	this.NameSpaces[name] = NewNameSpace()
}

// DeleteNameSpace deletes a namespace.
func (this *State) DeleteNameSpace(name string) {
	if strings.Contains(name, ":") {
		namesplit := nameSplitRegex.FindStringSubmatch(name)
		namespacelist := strings.TrimRight(namesplit[1], ":")
		basename := namesplit[3]
		namespace := this.ParseNameSpaceName(namespacelist)
		if _, ok := namespace.NameSpaces[basename]; ok {
			delete(namespace.NameSpaces, basename)
		}
		return
	}

	if _, ok := this.NameSpaces[name]; ok {
		delete(this.NameSpaces, name)
	}
}

// ParseNameSpaceName is just like ParseName but for the name of a namespace 
// instead of a variable, map, or command.
func (this *State) ParseNameSpaceName(name string) *NameSpace {
	if !strings.Contains(name, ":") {
		if this.NameSpaces[name] == nil {
			panic(fmt.Sprintf("Undeclared Namespace: \"%v\"", name))
		}
		return this.NameSpaces[name]
	}

	names := strings.Split(name, ":")

	return this.fetchNameSpaceName(names[1:], this.NameSpaces[names[0]])
}

func (this *State) fetchNameSpaceName(names []string, namespace *NameSpace) *NameSpace {
	if len(names) == 1 {
		if this.NameSpaces[names[0]] == nil {
			panic(fmt.Sprintf("Undeclared Namespace: \"%v\"", names[0]))
		}
		return namespace.NameSpaces[names[0]]
	}
	return this.fetchNameSpaceName(names[1:], namespace.NameSpaces[names[0]])
}

// Commands

// NewNativeCommand adds a new native command.
func (this *State) NewNativeCommand(name string, handler NativeCommand) {
	space, itemname := this.ParseName(name)
	if space != nil {
		rtn := new(Command)
		rtn.Native = true
		rtn.Handler = handler
		space.Commands[itemname] = rtn
		return
	}
	rtn := new(Command)
	rtn.Native = true
	rtn.Handler = handler
	this.Commands[itemname] = rtn
	return
}

// NewUserCommand adds a new user command (what else would it do?).
func (this *State) NewUserCommand(name string, code *Value, params []*Value) {
	space, itemname := this.ParseName(name)
	if space != nil {
		rtn := new(Command)
		rtn.Code = Compile(code.String(), code.Line)
		if params == nil {
			rtn.VarParams = true
		} else {
			for _, val := range params {
				rtn.Params = append(rtn.Params, val.String())
			}
		}

		space.Commands[itemname] = rtn
		return
	}

	rtn := new(Command)
	rtn.Code = Compile(code.String(), code.Line)
	if params == nil {
		rtn.VarParams = true
	} else {
		for _, val := range params {
			rtn.Params = append(rtn.Params, val.String())
		}
	}

	this.Commands[itemname] = rtn
	return
}

// GetCommand fetches a command by it's name.
func (this *State) GetCommand(name string) *Command {
	space, itemname := this.ParseName(name)
	if space != nil {
		if val, ok := space.Commands[itemname]; ok {
			return val
		}
		panic(fmt.Sprintf("Undeclared command: \"%v\"", name))
	}

	if val, ok := this.Commands[itemname]; ok {
		return val
	}
	panic(fmt.Sprintf("Undeclared command: \"%v\"", name))
}

// DeleteCommand removes a command.
func (this *State) DeleteCommand(name string) {
	space, itemname := this.ParseName(name)
	if space != nil {
		delete(space.Commands, itemname)
		return
	}
	delete(this.Commands, itemname)
	return
}

// Types

// NewNativeCommand adds a new native command.
func (this *State) RegisterIndexable(name string, handler IndexableFactory) {
	space, itemname := this.ParseName(name)
	if space != nil {
		
		space.Types[itemname] = handler
		return
	}
	this.Types[itemname] = handler
	return
}

// Debugger

func (this *State) RegisterDbgCallback(typ int, handler DbgHandler) {
	if this.Debug == nil {
		this.Debug = make([]DbgHandler, dbgrMaxType)
	}
	if typ >= dbgrMaxType || typ < 0 {
		panic("Callback type out of range.")
	}
	this.Debug[typ] = handler
}

func (this *State) DbgCallback(typ int) {
	if this.Debug == nil {
		return
	}
	if typ >= dbgrMaxType || typ < 0 {
		panic("Callback type out of range.")
	}
	
	if this.Debug[typ] == nil {
		return
	}
	this.Debug[typ](this)
}

// Exec

// Run executes a NCA script.
func (this *State) Run() (ret *Value, err error) {
	err = nil

	defer func() {
		ret = this.RetVal // this line is the only one not related to error handling
		
		if this.NoRecover {
			return
		}

		if x := recover(); x != nil {
			switch i := x.(type) {
			case error:
				err = i
			case string:
				if this.Code.Last().Line() != -1 {
					err = fmt.Errorf("%v On Line: %v Column: %v.", i, this.Code.Last().Line(), this.Code.Last().Column())
				} else {
					err = fmt.Errorf("%v", i)
				}
			default:
				err = errors.New(fmt.Sprint(i))
			}
			
			this.Code.Clear() // Remove any junk hanging around
		}
	}()

	this.Exec()
	this.Break = false // Just in case...
	this.Return = false
	this.Exit = false
	return
}

// Exec is exported for the use of commands ONLY!
// Exec is a subset of Run() and so it must be called from a command handler or the like.
func (this *State) Exec() {
	for !this.Code.Last().CheckLookAhead(TknINVALID) {
		this.RetVal = this.fetchValue()
		if this.Return || this.Exit || this.Break {
			this.Break = false
			this.Code.Remove()
			return
		}
	}
	this.Code.Remove()
	return
}

func (this *State) execCommand() {
	this.DbgCallback(DbgrAdvanceTkn)
	this.Code.Last().GetToken(TknCmdBegin)

	// Get the command's name
	name := this.fetchValue()
	if this.Break || this.Return || this.Exit {
		return
	}

	// Read the commands parameters if any
	params := make([]*Value, 0, 5)
	for !this.Code.Last().CheckLookAhead(TknCmdEnd) {
		params = append(params, this.fetchValue())
		if this.Break || this.Return || this.Exit {
			return
		}
	}

	this.DbgCallback(DbgrAdvanceTkn)
	this.Code.Last().GetToken(TknCmdEnd)

	this.GetCommand(name.String()).Call(this, params)
}

func (this *State) derefVar() {
	this.DbgCallback(DbgrAdvanceTkn)
	this.Code.Last().GetToken(TknDerefBegin)

	// Get the name or value to opperate on
	name := this.fetchValue()
	if this.Break || this.Return || this.Exit {
		return
	}

	// Read any index parameters
	params := make([]*Value, 0, 5)
	for !this.Code.Last().CheckLookAhead(TknDerefEnd) {
		params = append(params, this.fetchValue())
		if this.Break || this.Return || this.Exit {
			return
		}
	}
	
	this.DbgCallback(DbgrAdvanceTkn)
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
	
	for i := range params {
		if val.Type != TypObject {
			panic(fmt.Sprint("Non-Object value passed to an indexing deref opperator. Indexing level:", i))
		}
		obj := val.Indexable()
		if obj == nil {
			panic(fmt.Sprint("Non-Indexable object passed to an indexing deref opperator. Indexing level:", i))
		}
		
		val = obj.Get(params[i].String())
	}
	
	this.RetVal = val
	return
}

func (this *State) readObjLit() {
	this.DbgCallback(DbgrAdvanceTkn)
	this.Code.Last().GetToken(TknObjLitBegin)
	
	// Get the type name
	name := this.fetchValue().String()
	if this.Break || this.Return || this.Exit {
		return
	}
	
	// Generate key/value lists
	// you must have keys for all or none
	values := make([]*Value, 0, 10)
	keys := make([]string, 0, 10)
	hasKeys := false
	
	if !this.Code.Last().CheckLookAhead(TknObjLitEnd) {
		ret := this.fetchValue()
		if this.Break || this.Return || this.Exit {
			return
		}
		if this.Code.Last().CheckLookAhead(TknObjLitSplit) {
			hasKeys = true
			keys = append(keys, ret.String())
			this.DbgCallback(DbgrAdvanceTkn)
			this.Code.Last().GetToken(TknObjLitSplit)
			values = append(values, this.fetchValue())
			if this.Break || this.Return || this.Exit {
				return
			}
		} else {
			values = append(values, ret)
		}
	}
	for !this.Code.Last().CheckLookAhead(TknObjLitEnd) {
		ret := this.fetchValue()
		if this.Break || this.Return || this.Exit {
			return
		}
		if hasKeys {
			keys = append(keys, ret.String())
			this.DbgCallback(DbgrAdvanceTkn)
			this.Code.Last().GetToken(TknObjLitSplit)
			values = append(values, this.fetchValue())
			if this.Break || this.Return || this.Exit {
				return
			}
		} else {
			values = append(values, ret)
		}
	}
	
	this.DbgCallback(DbgrAdvanceTkn)
	this.Code.Last().GetToken(TknObjLitEnd)
	
	if hasKeys {
		this.RetVal = this.Types[name](keys, values)
		return
	}
	this.RetVal = this.Types[name](nil, values)
	return
}

func (this *State) fetchValue() *Value {
	switch this.Code.Last().LookAhead().Type {
	case TknString:
		this.DbgCallback(DbgrAdvanceTkn)
		this.Code.Last().GetToken(TknString)
		return TokenToValue(this.Code.Last().CurrentTkn())
	
	case TknCmdBegin:
		this.DbgCallback(DbgrEnterCmd)
		this.execCommand()
		this.DbgCallback(DbgrLeaveCmd)
		return this.RetVal
	
	case TknDerefBegin:
		this.DbgCallback(DbgrEnterDeref)
		this.derefVar()
		this.DbgCallback(DbgrLeaveDeref)
		return this.RetVal
	
	case TknObjLitBegin:
		this.DbgCallback(DbgrEnterObjLit)
		this.readObjLit()
		this.DbgCallback(DbgrLeaveObjLit)
		return this.RetVal
	
	default:
		ExitOnTokenExpected(this.Code.Last().LookAhead(), TknString, TknCmdBegin, TknDerefBegin, TknObjLitBegin)
	}
	panic("UNREACHABLE")
}
