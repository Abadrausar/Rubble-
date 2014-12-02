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

package structure

import "dctech/rex"
import "sync"

type Structure struct {
	data  map[string]*rex.Value
	lock  *sync.RWMutex
}

func NewStructureFromLit(script *rex.Script, keys []string, values []*rex.Value) *rex.Value {
	if keys == nil {
		rex.RaiseError("struct: May not initialize without keys.")
	}
	if len(keys) < 1 {
		rex.RaiseError("struct: Must have at least one key.")
	}
	if keys[0] != "proto" {
		rex.RaiseError("struct: First key must be \"proto\".")
	}
	if values[0].Type != rex.TypIndex {
		rex.RaiseError("struct: \"proto\" key's value must be type index.")
	}
	
	str := new(Structure)
	str.lock = new(sync.RWMutex)
	
	index := values[0].Data.(rex.Indexable)
	str.data = make(map[string]*rex.Value, int(index.Len()))
	for _, i := range index.Keys() {
		str.data[i] = index.Get(i)
	}
	
	for i, k := range keys[1:] {
		str.Set(k, values[i + 1])
	}
	
	return rex.NewValueIndex(str)
}

func (str *Structure) Get(index string) *rex.Value {
	str.lock.RLock()
	defer str.lock.RUnlock()
	
	if _, ok := str.data[index]; ok {
		tmp := str.data[index]
		return tmp
	}
	return rex.NewValue()
}

func (str *Structure) Set(index string, value *rex.Value) bool {
	str.lock.Lock()
	defer str.lock.Unlock()
	
	if _, ok := str.data[index]; ok {
		str.data[index] = value
		return true
	}
	return false
}

func (str *Structure) Exists(index string) bool {
	str.lock.RLock()
	defer str.lock.RUnlock()
	
	_, ok := str.data[index]
	return ok
}

func (str *Structure) Len() int64 {
	return int64(len(str.data))
}

func (str *Structure) Keys() []string {
	str.lock.RLock()
	defer str.lock.RUnlock()
	
	rtn := make([]string, 0, len(str.data))
	for key := range str.data {
		rtn = append(rtn, key)
	}
	return rtn
}

func (str *Structure) String() string {
	return rex.IndexableToString("map", str)
}
