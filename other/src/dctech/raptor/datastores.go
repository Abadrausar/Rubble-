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

import "sync"

// ObjectFactory is the function signature for a literal converter.
// keys may be nil, a specific implementation may panic with an error if
// a nil or non-nil keys is expected but not found.
type ObjectFactory func(script *Script, keys []string, values []*Value) *Value

// An Environment is where all variables are stored.
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
	Commands   *CommandStore
	NameSpaces *NameSpaceStore
	Types      *TypeStore
	Vars       *ValueStore
}

func NewNameSpace() *NameSpace {
	rtn := new(NameSpace)
	rtn.Vars = NewValueStore()
	rtn.Commands = NewCommandStore()
	rtn.NameSpaces = NewNameSpaceStore()
	rtn.Types = NewTypeStore()
	return rtn
}

// An EnvStore is an order sensitive stack of environments.
// Touch only if you know what you are doing!
type EnvStore []*Environment

func NewEnvStore() *EnvStore {
	rtn := new(EnvStore)
	*rtn = make([]*Environment, 0, 20)
	rtn.Add(NewEnvironment())
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

func (this *EnvStore) Clear() {
	*this = (*this)[0:0]
	this.Add(NewEnvironment())
}

func (this *EnvStore) ClearAllButRoot() {
	*this = (*this)[0:1]
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
	*this = (*this)[0:0]
}

func (this *BlockStore) Last() CodeSource {
	return (*this)[len(*this)-1]
}

// Concurrency support

// CommandStore is for storing script commands.
type CommandStore struct {
	data map[string]*Command
	lock *sync.RWMutex
}

// NewCommandStore creates and initializes a CommandStore.
func NewCommandStore() *CommandStore {
	rtn := new(CommandStore)
	rtn.data = make(map[string]*Command)
	rtn.lock = new(sync.RWMutex)
	return rtn
}

// Fetch retrieves a value from the store. Fetch is reentrant.
func (this *CommandStore) Fetch(name string) *Command {
	this.lock.RLock()
	rtn := this.data[name]
	this.lock.RUnlock()
	return rtn
}

// Store adds/updates a value in the store. Store is reentrant.
func (this *CommandStore) Store(name string, obj *Command) {
	this.lock.Lock()
	this.data[name] = obj
	this.lock.Unlock()
}

// Delete removes a value from the store. Delete is reentrant.
func (this *CommandStore) Delete(name string) {
	this.lock.Lock()
	delete(this.data, name)
	this.lock.Unlock()
}

// Exist checks to see if an item exists in the store. Exist is reentrant.
func (this *CommandStore) Exist(name string) bool {
	this.lock.RLock()
	_, rtn := this.data[name]
	this.lock.RUnlock()
	return rtn
}

// DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
func (this *CommandStore) BeginLowLevel() map[string]*Command {
	this.lock.Lock()
	return this.data
}

// DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
func (this *CommandStore) EndLowLevel() {
	this.lock.Unlock()
}

// NameSpaceStore is for storing script namespaces.
type NameSpaceStore struct {
	data map[string]*NameSpace
	lock *sync.RWMutex
}

// NewNameSpaceStore creates and initializes a NameSpaceStore.
func NewNameSpaceStore() *NameSpaceStore {
	rtn := new(NameSpaceStore)
	rtn.data = make(map[string]*NameSpace)
	rtn.lock = new(sync.RWMutex)
	return rtn
}

// Fetch retrieves a value from the store. Fetch is reentrant.
func (this *NameSpaceStore) Fetch(name string) *NameSpace {
	this.lock.RLock()
	rtn := this.data[name]
	this.lock.RUnlock()
	return rtn
}

// Store adds/updates a value in the store. Store is reentrant.
func (this *NameSpaceStore) Store(name string, obj *NameSpace) {
	this.lock.Lock()
	this.data[name] = obj
	this.lock.Unlock()
}

// Delete removes a value from the store. Delete is reentrant.
func (this *NameSpaceStore) Delete(name string) {
	this.lock.Lock()
	delete(this.data, name)
	this.lock.Unlock()
}

// Exist checks to see if an item exists in the store. Exist is reentrant.
func (this *NameSpaceStore) Exist(name string) bool {
	this.lock.RLock()
	_, rtn := this.data[name]
	this.lock.RUnlock()
	return rtn
}

// DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
func (this *NameSpaceStore) BeginLowLevel() map[string]*NameSpace {
	this.lock.Lock()
	return this.data
}

// DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
func (this *NameSpaceStore) EndLowLevel() {
	this.lock.Unlock()
}

// TypeStore is for storing ObjectFactorys.
type TypeStore struct {
	data map[string]ObjectFactory
	lock *sync.RWMutex
}

// NewTypeStore creates and initializes a TypeStore.
func NewTypeStore() *TypeStore {
	rtn := new(TypeStore)
	rtn.data = make(map[string]ObjectFactory)
	rtn.lock = new(sync.RWMutex)
	return rtn
}

// Fetch retrieves a value from the store. Fetch is reentrant.
func (this *TypeStore) Fetch(name string) ObjectFactory {
	this.lock.RLock()
	rtn := this.data[name]
	this.lock.RUnlock()
	return rtn
}

// Store adds/updates a value in the store. Store is reentrant.
func (this *TypeStore) Store(name string, obj ObjectFactory) {
	this.lock.Lock()
	this.data[name] = obj
	this.lock.Unlock()
}

// Delete removes a value from the store. Delete is reentrant.
func (this *TypeStore) Delete(name string) {
	this.lock.Lock()
	delete(this.data, name)
	this.lock.Unlock()
}

// Exist checks to see if an item exists in the store. Exist is reentrant.
func (this *TypeStore) Exist(name string) bool {
	this.lock.RLock()
	_, rtn := this.data[name]
	this.lock.RUnlock()
	return rtn
}

// DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
func (this *TypeStore) BeginLowLevel() map[string]ObjectFactory {
	this.lock.Lock()
	return this.data
}

// DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
func (this *TypeStore) EndLowLevel() {
	this.lock.Unlock()
}

// TypeStore is for storing ObjectFactorys.
type ValueStore struct {
	data map[string]*Value
	lock *sync.RWMutex
}

// NewTypeStore creates and initializes a TypeStore.
func NewValueStore() *ValueStore {
	rtn := new(ValueStore)
	rtn.data = make(map[string]*Value)
	rtn.lock = new(sync.RWMutex)
	return rtn
}

// Fetch retrieves a value from the store. Fetch is reentrant.
func (this *ValueStore) Fetch(name string) *Value {
	this.lock.RLock()
	rtn := this.data[name]
	this.lock.RUnlock()
	return rtn
}

// Store adds/updates a value in the store. Store is reentrant.
func (this *ValueStore) Store(name string, obj *Value) {
	this.lock.Lock()
	this.data[name] = obj
	this.lock.Unlock()
}

// Delete removes a value from the store. Delete is reentrant.
func (this *ValueStore) Delete(name string) {
	this.lock.Lock()
	delete(this.data, name)
	this.lock.Unlock()
}

// Exist checks to see if an item exists in the store. Exist is reentrant.
func (this *ValueStore) Exist(name string) bool {
	this.lock.RLock()
	_, rtn := this.data[name]
	this.lock.RUnlock()
	return rtn
}

// DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
func (this *ValueStore) BeginLowLevel() map[string]*Value {
	this.lock.Lock()
	return this.data
}

// DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
func (this *ValueStore) EndLowLevel() {
	this.lock.Unlock()
}
