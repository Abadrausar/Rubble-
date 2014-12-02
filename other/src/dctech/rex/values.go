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

//import "fmt"

// valueStore

// valueStore is a data store for global variables.
type valueStore struct {
	ntoi map[string]int // Offset from Name
	iton []string       // Name from Offset
	data []*Value
	lock *sync.RWMutex
}

// newValueStore creates a new valueStore, ready to use.
func newValueStore() *valueStore {
	return &valueStore{
		ntoi: make(map[string]int, 10),
		iton: make([]string, 0, 10),
		data: make([]*Value, 0, 50),
		lock: new(sync.RWMutex),
	}
}

// add a new variable to the valueStore, returns the index.
func (store *valueStore) add(name string) int {
	store.lock.Lock()
	defer store.lock.Unlock()
	
	slot := len(store.data)
	if _, ok := store.ntoi[name]; ok {
		//RaiseError("Value slot already exists: " + name)
		return store.ntoi[name]
	}
	store.ntoi[name] = slot
	store.iton = append(store.iton, name)
	store.data = append(store.data, nil)
	return slot
}

// set an existing variable by index.
func (store *valueStore) set(index int, val *Value) {
	store.lock.Lock()
	defer store.lock.Unlock()
	
	if index < 0 || index >= len(store.data) {
		RaiseError("Value slot does not exist.")
	}
	store.data[index] = val
}

// addAndSet adds a new variable and sets it's initial value, returns the index.
func (store *valueStore) addAndSet(name string, val *Value) int {
	store.lock.Lock()
	defer store.lock.Unlock()
	
	slot := len(store.data)
	if _, ok := store.ntoi[name]; ok {
		//RaiseError("Value slot already exists.")
		store.data[store.ntoi[name]] = val
		return store.ntoi[name]
	}
	store.ntoi[name] = slot
	store.iton = append(store.iton, name)
	store.data = append(store.data, val)
	return slot
}

// get a variables value by index.
func (store *valueStore) get(index int) *Value {
	store.lock.RLock()
	defer store.lock.RUnlock()
	
	if index < 0 || index >= len(store.data) {
		RaiseError("Value slot does not exist.")
	}
	dat := store.data[index]
	if dat == nil {
		RaiseError("Value slot value invalid.")
	}
	return dat
}

// lookup a variables index by name.
func (store *valueStore) lookup(name string) int {
	store.lock.RLock()
	defer store.lock.RUnlock()
	
	if slot, ok := store.ntoi[name]; ok {
		return slot
	}
	RaiseError("No Value slot named: " + name + " exists.")
	panic("UNREACHABLE")
}

// LocalValueStore

// LocalValueStore is used for storing local variables, the only place this is used is in the Script.
type LocalValueStore struct {
	data []*Value
	envs []*env
}

// newLocalValueStore creates a new (empty) LocalValueStore.
func newLocalValueStore() *LocalValueStore {
	return &LocalValueStore{
		data: make([]*Value, 1024),
		envs: make([]*env, 0, 20),
	}
}

// Env is mostly for internal use by the LocalValueStore.
type env struct {
	top      int // top is last index + 1
	base     int
	block    *Code // for the meta data
	lOffsets []int
	eOffsets []int
}

// lookup gets a global index by name, this searches up-chain and is for use in resolving external variable indexes.
// (only used in lvs.Add)
func (lvs *LocalValueStore) lookup(state *State, name string) int {
	for x := len(lvs.envs)-1; x >= 0; x-- {
		if i, ok := lvs.envs[x].block.meta.ntoi[name]; ok {
			return lvs.envs[x].lOffsets[i]
		}
	}
	
	// Not in local block chain, look for a global with the correct name
	// (and then transform the index into a negative number)
	return 0 - (state.global.vars.lookup(name) + 1)
}

// Add a new environment using the meta-data from code.
// This assumes that the meta-data in code is well formed!
func (lvs *LocalValueStore) Add(state *State, code *Code) {
	// calculate bounds
	base := 0
	if len(lvs.envs) > 0 {
		base = lvs.envs[len(lvs.envs)-1].top
	}
	top := base + len(code.meta.iton)

	// Extend the local store if needed
	for top > len(lvs.data) {
		lvs.data = append(lvs.data, nil)
	}

	// Create the offset tables
	loff := make([]int, len(code.meta.iton))
	for i := range loff {
		loff[i] = i + base
	}

	eoff := make([]int, len(code.meta.eiton))
	for i := range eoff {
		eoff[i] = lvs.lookup(state, code.meta.eiton[i])
	}
	
	// add the env
	lvs.envs = append(lvs.envs, &env{
		top:      top,
		base:     base,
		block:    code,
		lOffsets: loff,
		eOffsets: eoff,
	})
}

// Removes the last environment added.
// It is an error to call this with no environments available!
func (lvs *LocalValueStore) Remove() {
	if len(lvs.envs) == 0 {
		RaiseError("No environments in local store.")
	}

	// clear freed values
	env := lvs.envs[len(lvs.envs)-1]
	for i := env.base; i < env.top; i++ {
		lvs.data[i] = nil
	}

	// pop the env off the stack
	lvs.envs[len(lvs.envs)-1] = nil
	lvs.envs = lvs.envs[:len(lvs.envs)-1]
}

// RemoveNoClear removes the last environment added WITHOUT clearing the freed slots.
// Use with NewCodeShell for implementing interactive shells and the like.
func (lvs *LocalValueStore) RemoveNoClear() {
	if len(lvs.envs) == 0 {
		return // This is called in cases where there may not be an env (rarely)
	}

	// pop the env off the stack
	lvs.envs[len(lvs.envs)-1] = nil
	lvs.envs = lvs.envs[:len(lvs.envs)-1]
}

// Clears the entire LocalValueStore.
func (lvs *LocalValueStore) Clear() {
	if len(lvs.envs) == 0 {
		return
	}

	for i := 0; i < len(lvs.envs); i++ {
		lvs.Remove()
	}
}

// ClearToRoot removes all but the first environment.
// Generally used with RemoveNoclear.
func (lvs *LocalValueStore) ClearToRoot() {
	if len(lvs.envs) < 2 {
		return
	}

	for i := 1; i < len(lvs.envs); i++ {
		lvs.Remove()
	}
}

func (lvs *LocalValueStore) gIndex(index int) int {
	if len(lvs.envs) == 0 {
		RaiseError("No environments in local store.")
	}
	env := lvs.envs[len(lvs.envs)-1]

	offsets := env.lOffsets
	if index < 0 {
		index = 0 - (index + 1)
		offsets = env.eOffsets
	}

	if index < 0 || index >= len(offsets) {
		RaiseError("Local index out of range.")
	}
	gindex := offsets[index]
	if gindex < 0 || gindex >= len(lvs.data) {
		RaiseError("Index out of range.")
	}
	return gindex
}

// Get a value by local index.
func (lvs *LocalValueStore) Get(index int) *Value {
	return lvs.data[lvs.gIndex(index)]
}

// Set a value by local index.
func (lvs *LocalValueStore) Set(index int, val *Value) {
	lvs.data[lvs.gIndex(index)] = val
}

// Returns a positive global index if the variable is actually a global.
// If the variable is local -1 is returned.
func (lvs *LocalValueStore) IsGlobal(index int) int {
	if index > 0 {
		return -1
	}
	index = 0 - (index + 1)
	
	env := lvs.envs[len(lvs.envs)-1]
	
	if index < 0 || index >= len(env.eOffsets) {
		return -1
	}
	if env.eOffsets[index] > 0 {
		return -1
	}
	return 0 - (env.eOffsets[index] + 1)
}
