package nca5

import "fmt"
import "regexp"
import "strings"
import "strconv"
import "errors"
import "dctech/ncalex"

// Panicking with an ExitScript exits the current script.
//	// example:
//	panic(ExitScript("return value"))
// Needless to say this is for commands only.
type ExitScript string

// State handles EVERYTHING to do with nca scripts, including running them.
// The majority of the fields are exported for the use of commands only.
type State struct {
	NameSpaces map[string]*NameSpace
	Commands   map[string]*Command
	RetVal     *Value // The return value of the last command
	Return     bool // true when a return is active
	Break      bool // true when a break is active
	Error      bool // set by some commands on error, this is NOT automaticly reset!
	NoRecover  bool // Do not recover errors, this makes it easier to debug the parser.
	Envs       *EnvStore // The script environments, you should not need to touch this
	Code       *BlockStore // This is where script code is stored
}

// NewState creates (and initalizes) a new state.
func NewState() *State {
	rtn := new(State)
	rtn.NameSpaces = make(map[string]*NameSpace)
	rtn.Commands = make(map[string]*Command)
	rtn.Envs = NewEnvStore()
	rtn.Envs.Add(NewEnvironment())
	rtn.Code = NewBlockStore()
	return rtn
}

// Variables and maps

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

// NewMap creates a new map.
func (this *State) NewMap(name string) {
	this.NewVar(name, NewValueFromObject(NewScriptMap()))
}

// NewArray creates a new array.
func (this *State) NewArray(name string, size int) {
	this.NewVar(name, NewValueFromObject(NewScriptArraySized(size)))
}

// GetValueObject gets a value from a map or array.
func (this *State) GetValueObject(name, index string) *Value {
	val := this.GetValue(name)
	if val.HasObject() {
		return val.Get(index)
	}
	panic(fmt.Sprintf("Not an object: \"%v\"", name))
}

// SetValueObject sets a value in a map or array.
func (this *State) SetValueObject(name, index string, value *Value) {
	val := this.GetValue(name)
	if val.HasObject() {
		val.Set(index, value)
		return
	}
	panic(fmt.Sprintf("Not an object: \"%v\"", name))
}

// GetObjectKeys gets a list of an array or maps keys for advanced command usage.
func (this *State) GetObjectKeys(name string) []string {
	val := this.GetValue(name)
	if val.HasObject() {
		return val.Keys()
	}
	panic(fmt.Sprintf("Not an object: \"%v\"", name))
}

// ObjectKeyExists returns true if key "key" exists in map or array "name".
func (this *State) ObjectKeyExists(name, key string) bool {
	val := this.GetValue(name)
	if val.HasObject() {
		return val.Exists(key)
	}
	panic(fmt.Sprintf("Not an object: \"%v\"", name))
}

// AddParams creates the special "params" array using strings.
func (this *State) AddParams(params ...string) {
	this.NewArray("params", len(params))
	for i, val := range params {
		this.SetValueObject("params", strconv.FormatInt(int64(i), 10), NewValue(val))
	}
}

// AddParamsValue creates the special "params" array using values.
func (this *State) AddParamsValue(params ...*Value) {
	this.NewArray("params", len(params))
	for i, val := range params {
		this.SetValueObject("params", strconv.FormatInt(int64(i), 10), val)
	}
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
		rtn.Code = code
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
	rtn.Code = code
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

// Exec

// Run executes a nca script.
func (this *State) Run() (ret *Value, err error) {
	err = nil

	defer func() {
		ret = this.RetVal // this line is the only one not related to error handling
		
		if this.NoRecover {
			return
		}

		if x := recover(); x != nil {
			switch i := x.(type) {
			case ExitScript:
				ret = NewValue(string(i))
			case error:
				err = i
			case string:
				if this.Code.Last().PositionValid {
					err = fmt.Errorf("%v On Line: %v, Column: %v.", i, this.Code.Last().LastLine, this.Code.Last().LastColumn)
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
	this.Return = false
	return
}

// Exec is exported for the use of commands ONLY!
// A subset of Run() and so it must be called from a command handler or the like.
func (this *State) Exec() {
	for this.Code.Last().CheckLookAhead(ncalex.TknOpenParen) {
		this.execCommand()

		if this.Break {
			this.Break = false
			this.Code.Remove()
			return
		}

		if this.Return {
			this.Code.Remove()
			return
		}
	}
	if this.Code.Last().CheckLookAhead(ncalex.TknINVALID) {
		this.Code.Remove()
		return
	}
	panic("Malformed input.")
}

func (this *State) execCommand() {
	this.Code.Last().GetToken(ncalex.TknOpenParen)

	// Get the commands name
	name := new(Value)
	switch this.Code.Last().Look.Type {
	case ncalex.TknString:
		this.Code.Last().GetToken(ncalex.TknString)
		name = TokenToValue(this.Code.Last().Current)
	case ncalex.TknOpenParen:
		this.execCommand()
		name = this.RetVal
		if this.Return {
			return
		}
	case ncalex.TknOpenSqBracket:
		this.derefVar()
		name = this.RetVal
		if this.Return {
			return
		}
	default:
		ncalex.ExitOnTokenExpected(this.Code.Last().Look, ncalex.TknString, ncalex.TknOpenParen, ncalex.TknOpenSqBracket)
	}

	// Read the commands parameters if any
	params := make([]*Value, 0, 5)
	for !this.Code.Last().CheckLookAhead(ncalex.TknCloseParen) {
		switch this.Code.Last().Look.Type {
		case ncalex.TknString:
			this.Code.Last().GetToken(ncalex.TknString)
			params = append(params, TokenToValue(this.Code.Last().Current))
		case ncalex.TknOpenParen:
			this.execCommand()
			if this.Return {
				return
			}
			params = append(params, this.RetVal)
		case ncalex.TknOpenSqBracket:
			this.derefVar()
			if this.Return {
				return
			}
			params = append(params, this.RetVal)
		default:
			ncalex.ExitOnTokenExpected(this.Code.Last().Look, ncalex.TknString, ncalex.TknOpenParen, ncalex.TknOpenSqBracket, ncalex.TknCloseParen)
		}
	}

	this.Code.Last().GetToken(ncalex.TknCloseParen)

	this.GetCommand(name.String()).Call(this, params)
}

func (this *State) derefVar() {
	this.Code.Last().GetToken(ncalex.TknOpenSqBracket)

	// Get the variables name
	name := new(Value)
	switch this.Code.Last().Look.Type {
	case ncalex.TknString:
		this.Code.Last().GetToken(ncalex.TknString)
		name = TokenToValue(this.Code.Last().Current)
	case ncalex.TknOpenParen:
		this.execCommand()
		name = this.RetVal
		if this.Return {
			return
		}
	case ncalex.TknOpenSqBracket:
		this.derefVar()
		name = this.RetVal
		if this.Return {
			return
		}
	default:
		ncalex.ExitOnTokenExpected(this.Code.Last().Look, ncalex.TknString, ncalex.TknOpenParen, ncalex.TknOpenSqBracket)
	}

	// Get the index value if any
	index := new(Value)
	switch this.Code.Last().Look.Type {
	case ncalex.TknString:
		this.Code.Last().GetToken(ncalex.TknString)
		index = TokenToValue(this.Code.Last().Current)
	case ncalex.TknOpenParen:
		this.execCommand()
		index = this.RetVal
		if this.Return {
			return
		}
	case ncalex.TknOpenSqBracket:
		this.derefVar()
		index = this.RetVal
		if this.Return {
			return
		}
	}

	this.Code.Last().GetToken(ncalex.TknCloseSqBracket)

	// Get and return the value
	if index.String() == "" {
		this.RetVal = this.GetValue(name.String())
		return
	}
	this.RetVal = this.GetValueObject(name.String(), index.String())
	return
}
