package nca4

import "dctech/ncalex"

// An Enviroment is where all maps and variables are stored.
// Nesting is handled by State via EnvStore.
// Touch only if you know what you are doing!
type Environment struct {
	Vars map[string]*Value
	Maps map[string]map[string]*Value
}

func NewEnvironment() *Environment {
	rtn := new(Environment)
	rtn.Vars = make(map[string]*Value)
	rtn.Maps = make(map[string]map[string]*Value)
	return rtn
}

// A NameSpace is a named, hierarchical environment that may
// also contain other namespaces and commands.
// Touch only if you know what you are doing!
type NameSpace struct {
	Vars       map[string]*Value
	Maps       map[string]map[string]*Value
	Commands   map[string]*Command
	NameSpaces map[string]*NameSpace
}

func NewNameSpace() *NameSpace {
	rtn := new(NameSpace)
	rtn.Vars = make(map[string]*Value)
	rtn.Maps = make(map[string]map[string]*Value)
	rtn.Commands = make(map[string]*Command)
	rtn.NameSpaces = make(map[string]*NameSpace)
	return rtn
}

// An EnvStore is an order sensitive stack of environments.
// Touch only if you know what you are doing!
type EnvStore []*Environment

func NewEnvStore() *EnvStore {
	rtn := new(EnvStore)
	*rtn = make([]*Environment, 0, 20)
	return rtn
}

func (this *EnvStore) Add(val *Environment) {
	*this = append(*this, val)
}

func (this *EnvStore) Remove() *Environment {
	rtn := (*this)[len(*this)-1]
	(*this)[len(*this)-1] = nil
	*this = (*this)[:len(*this)-1]
	// There must always be at least one valid env
	if len(*this) == 0 {
		this.Add(NewEnvironment())
	}
	return rtn
}

func (this *EnvStore) Last() *Environment {
	return (*this)[len(*this)-1]
}

// A BlockStore is a stack of *ncalex.Lexer representing blocks of code.
// Touch only if you know what you are doing!
type BlockStore []*ncalex.Lexer

func NewBlockStore() *BlockStore {
	rtn := new(BlockStore)
	*rtn = make([]*ncalex.Lexer, 0, 20)
	return rtn
}

// Add is an exception to the hands-off rule, use this to add more code to a state.
func (this *BlockStore) Add(input string) {
	*this = append(*this, ncalex.NewLexer(input, 1, 1))
}

// AddLexer is an exception to the hands-off rule, use this to add more code to a state.
func (this *BlockStore) AddLexer(val *ncalex.Lexer) {
	*this = append(*this, val)
}

func (this *BlockStore) Remove() *ncalex.Lexer {
	rtn := (*this)[len(*this)-1]
	(*this)[len(*this)-1] = nil
	*this = (*this)[:len(*this)-1]
	return rtn
}

func (this *BlockStore) Clear() {
	*this = make([]*ncalex.Lexer, 0, 20)
}

func (this *BlockStore) Last() *ncalex.Lexer {
	return (*this)[len(*this)-1]
}
