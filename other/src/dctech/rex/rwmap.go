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

type rwMap struct {
	data map[string]*Value
	lock *sync.RWMutex
}

// NewMap creates a new script map.
func NewMap(data map[string]*Value) EditIndexable {
	return &rwMap{
		data: data,
		lock: new(sync.RWMutex),
	}
}

// NewMapFromLit is an ObjectFactory for script maps.
func NewMapFromLit(script *Script, keys []string, values []*Value) *Value {
	if len(values) != 0 {
		if keys == nil {
			RaiseError("Map may not be initialized without keys.")
		}
	}
	if len(values) != len(keys) {
		RaiseError("Map: Key count does not equal value count.")
	}

	tmp := make(map[string]*Value, len(keys))
	for i := range keys {
		tmp[keys[i]] = values[i]
	}
	return NewValueIndex(&rwMap{
		data: tmp,
		lock: new(sync.RWMutex),
	})
}

func (this *rwMap) Get(index string) *Value {
	this.lock.RLock()
	defer this.lock.RUnlock()
	
	if _, ok := this.data[index]; ok {
		return this.data[index]
	}
	return NewValue()
}

func (this *rwMap) Set(index string, value *Value) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	
	this.data[index] = value
	return true
}

func (this *rwMap) Exists(index string) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	
	_, ok := this.data[index]
	return ok
}

func (this *rwMap) Len() int64 {
	return int64(len(this.data))
}

func (this *rwMap) Keys() []string {
	this.lock.RLock()
	defer this.lock.RUnlock()
	
	rtn := make([]string, 0, len(this.data))
	for key := range this.data {
		rtn = append(rtn, key)
	}
	return rtn
}

func (this *rwMap) String() string {
	this.lock.RLock()
	defer this.lock.RUnlock()
	
	return IndexableToString("map", this)
}

func (this *rwMap) CodeString() string {
	this.lock.RLock()
	defer this.lock.RUnlock()
	
	return IndexableToCodeString("map", this, true)
}
