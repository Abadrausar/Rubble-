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

// An Enviroment is where all variables are stored.
// Nesting is handled by State via EnvStore.
// Touch only if you know what you are doing!
type Environment struct {
	Vars map[string]*Value
}

func NewEnvironment() *Environment {
	rtn := new(Environment)
	rtn.Vars = make(map[string]*Value)
	return rtn
}

// A NameSpace is a named, hierarchical environment that may
// also contain other namespaces and commands.
// Touch only if you know what you are doing!
type NameSpace struct {
	Vars       map[string]*Value
	Commands   map[string]*Command
	Types      map[string]ObjectFactory
	NameSpaces map[string]*NameSpace
}

func NewNameSpace() *NameSpace {
	rtn := new(NameSpace)
	rtn.Vars = make(map[string]*Value)
	rtn.Commands = make(map[string]*Command)
	rtn.Types = make(map[string]ObjectFactory)
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

// A BlockStore is a stack of CodeSources representing blocks of code.
// Touch only if you know what you are doing!
type BlockStore []CodeSource

func NewBlockStore() *BlockStore {
	rtn := new(BlockStore)
	*rtn = make([]CodeSource, 0, 20)
	return rtn
}

// Add is an exception to the hands-off rule, use this to add more code to a state.
// This creates a new lexer.
func (this *BlockStore) Add(input string) {
	*this = append(*this, NewLexer(input, 1, 0))
}

// AddCodeSource is an exception to the hands-off rule, use this to add more code to a state.
// Use if you already have a Lexer or CompiledLexer.
func (this *BlockStore) AddCodeSource(val CodeSource) {
	*this = append(*this, val)
}

// AddCompiledScript is an exception to the hands-off rule, use this to add more code to a state.
// Use if you already have a CompiledScript.
func (this *BlockStore) AddCompiledScript(val *CompiledScript) {
	*this = append(*this, NewCompiledLexer(val))
}

func (this *BlockStore) Remove() CodeSource {
	rtn := (*this)[len(*this)-1]
	(*this)[len(*this)-1] = nil
	*this = (*this)[:len(*this)-1]
	return rtn
}

func (this *BlockStore) Clear() {
	*this = make([]CodeSource, 0, 20)
}

func (this *BlockStore) Last() CodeSource {
	return (*this)[len(*this)-1]
}
