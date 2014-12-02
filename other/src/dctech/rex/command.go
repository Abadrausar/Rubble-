/*
Copyright 2014 by Milo Christiansen

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

package rex

import "sync"

// NativeCommand is the function signature for a native command handler.
type NativeCommand func(*Script, []*Value)

// A Command stores a native or user script command.
type Command struct {
	// The native command handler, if non-nil all other fields are ignored.
	handler NativeCommand

	// The user command code.
	// All the parameter info is encoded in the block meta-data.
	block *Value
}

// NewNativeCommand creates a new native command.
func NewNativeCommand(handler NativeCommand) *Command {
	return &Command{
		handler: handler,
	}
}

// NewUserCommand creates a new user command (what else would it do?).
func NewUserCommand(code *Value) *Command {
	return &Command{
		block: code,
	}
}

// call runs a command.
func (cmd *Command) call(script *Script, params []*Value) {
	// Native command
	if cmd.handler != nil {
		cmd.handler(script, params)
		return
	}

	// User command

	// Parameters are always the first variables in a block.
	block := cmd.block.Data.(*Code)
	script.Locals.Add(block)
	
	script.setParams(block, params)

	script.Exec(block)
	script.Locals.Remove()
	script.Return = false
	return
}

// commandStore is a basic data store for commands.
type commandStore struct {
	ntoi map[string]int // Offset from Name
	iton []string       // Name from Offset
	data []*Command
	lock *sync.RWMutex // I may be paranoid, but better safe than sorry...
}

// newCommandStore does exactly what it's name suggests.
func newCommandStore() *commandStore {
	return &commandStore{
		ntoi: make(map[string]int, 10),
		iton: make([]string, 0, 10),
		data: make([]*Command, 0, 50),
		lock: new(sync.RWMutex),
	}
}

// add adds a new slot and sets it's Command object, returns the index.
func (store *commandStore) add(name string, command *Command) int {
	store.lock.Lock()
	defer store.lock.Unlock()
	
	slot := len(store.data)
	if _, ok := store.ntoi[name]; ok {
		RaiseError("Command slot already exists: " + name)
	}
	store.ntoi[name] = slot
	store.iton = append(store.iton, name)
	store.data = append(store.data, command)
	return slot
}

// get a command by index.
func (store *commandStore) get(index int) *Command {
	store.lock.RLock()
	defer store.lock.RUnlock()
	
	if index < 0 || index >= len(store.data) {
		RaiseError("Command slot does not exist.")
	}
	dat := store.data[index]
	if dat == nil {
		RaiseError("Command slot value invalid.")
	}
	return dat
}

// lookup a command index by name.
func (store *commandStore) lookup(name string) int {
	store.lock.RLock()
	defer store.lock.RUnlock()
	
	if slot, ok := store.ntoi[name]; ok {
		return slot
	}
	RaiseError("No Command slot named: " + name + " exists.")
	panic("UNREACHABLE")
}
