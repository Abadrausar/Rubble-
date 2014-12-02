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

// A Module is a named, hierarchical global data store that may
// contain other Modules, variables, types, and commands.
type Module struct {
	modules  *moduleStore
	types    *typeStore
	vars     *valueStore
}

// newModule will create a new, empty Module.
func newModule() *Module {
	return &Module{
		modules: newModuleStore(),
		types: newTypeStore(),
		vars: newValueStore(),
	}
}

// RegisterVar will create a new Module variable.
func (mod *Module) RegisterVar(name string, val *Value) {
	mod.vars.addAndSet(name, val)
}

// RegisterCommand will create a new Module native command.
func (mod *Module) RegisterCommand(name string, handler NativeCommand) {
	mod.vars.addAndSet(name, NewValueCommand(handler))
}

// RegisterType will create a new Module indexable type.
func (mod *Module) RegisterType(name string, typ ObjectFactory) {
	mod.types.add(name, typ)
}

// RegisterModule will create a new subModule.
// The Module is returned so that it can be immediately used to register global data.
func (mod *Module) RegisterModule(name string) *Module {
	index := mod.modules.add(name, newModule())
	return mod.modules.get(index)
}

// FetchModule will retrieve a module by name.
// Returns nil if the module does not exist.
func (mod *Module) FetchModule(name string) *Module {
	if !mod.modules.exists(name) {
		return nil
	}
	return mod.modules.get(mod.modules.lookup(name))
}

// Module store is exactly like all the other data stores.
// (this is one of those rare cases where I miss generics...)
type moduleStore struct {
	ntoi map[string]int // Offset from Name
	iton []string       // Name from Offset
	data []*Module
	lock *sync.RWMutex
}

// newModuleStore does exactly what it's name suggests.
func newModuleStore() *moduleStore {
	return &moduleStore{
		ntoi: make(map[string]int, 10),
		iton: make([]string, 0, 10),
		data: make([]*Module, 0, 50),
		lock: new(sync.RWMutex),
	}
}

// add adds a new Module slot and sets it's value, returns the index.
func (store *moduleStore) exists(name string) bool {
	store.lock.RLock()
	defer store.lock.RUnlock()
	
	_, ok := store.ntoi[name]
	return ok
}

// add adds a new Module slot and sets it's value, returns the index.
func (store *moduleStore) add(name string, mod *Module) int {
	store.lock.Lock()
	defer store.lock.Unlock()
	
	slot := len(store.data)
	if _, ok := store.ntoi[name]; ok {
		//RaiseError("Module slot already exists: " + name)
		return store.ntoi[name]
	}
	store.ntoi[name] = slot
	store.iton = append(store.iton, name)
	store.data = append(store.data, mod)
	return slot
}

// get a Module by index.
func (store *moduleStore) get(index int) *Module {
	store.lock.RLock()
	defer store.lock.RUnlock()

	if index < 0 || index >= len(store.data) {
		RaiseError("Module slot does not exist.")
	}
	dat := store.data[index]
	if dat == nil {
		RaiseError("Module slot value invalid.")
	}
	return dat
}

// lookup a Module index by name.
func (store *moduleStore) lookup(name string) int {
	store.lock.RLock()
	defer store.lock.RUnlock()
	
	if slot, ok := store.ntoi[name]; ok {
		return slot
	}
	RaiseError("No Module slot named: " + name + " exists.")
	panic("UNREACHABLE")
}
