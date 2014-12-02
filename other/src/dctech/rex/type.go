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

// ObjectFactory is the function signature for a literal converter.
// keys may be nil, a specific implementation may panic with an error if
// a nil or non-nil keys is expected but not found.
type ObjectFactory func(script *Script, keys []string, values []*Value) *Value

// typeStore is where all the literal converters hang out.
type typeStore struct {
	ntoi map[string]int // Offset from Name
	iton []string       // Name from Offset
	data []ObjectFactory
	lock *sync.RWMutex
}

// newTypeStore returns a new typeStore, and if you couldn't figure that out you have problems.
func newTypeStore() *typeStore {
	return &typeStore{
		ntoi: make(map[string]int, 10),
		iton: make([]string, 0, 10),
		data: make([]ObjectFactory, 0, 50),
		lock: new(sync.RWMutex),
	}
}

// add adds a new type and sets it's ObjectFactory, returns the index.
func (store *typeStore) add(name string, typ ObjectFactory) int {
	store.lock.Lock()
	defer store.lock.Unlock()
	
	slot := len(store.data)
	if _, ok := store.ntoi[name]; ok {
		//RaiseError("Type slot already exists: " + name)
		return store.ntoi[name]
	}
	store.ntoi[name] = slot
	store.iton = append(store.iton, name)
	store.data = append(store.data, typ)
	return slot
}

// get a type's ObjectFactory by index.
func (store *typeStore) get(index int) ObjectFactory {
	store.lock.RLock()
	defer store.lock.RUnlock()
	
	if index < 0 || index >= len(store.data) {
		RaiseError("Type slot does not exist.")
	}
	dat := store.data[index]
	if dat == nil {
		RaiseError("Type slot value invalid.")
	}
	return dat
}

// lookup a type's index by name.
func (store *typeStore) lookup(name string) int {
	store.lock.RLock()
	defer store.lock.RUnlock()
	
	if slot, ok := store.ntoi[name]; ok {
		return slot
	}
	RaiseError("No Type slot named: " + name + " exists.")
	panic("UNREACHABLE")
}
