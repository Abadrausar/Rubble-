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

package sort

import "dctech/rex"
import "sort"
import "sync"

type OrderedMap struct {
	data  map[string]*rex.Value
	order []string
	lock  *sync.RWMutex
}

func NewOrderedMap() rex.EditIndexable {
	this := new(OrderedMap)
	this.data = make(map[string]*rex.Value, 20)
	this.order = make([]string, 0, 20)
	this.lock = new(sync.RWMutex)
	return this
}

func NewOrderedMapFromIndexable(input rex.Indexable) rex.EditIndexable {
	this := new(OrderedMap)
	keys := input.Keys()
	this.order = make([]string, 0, len(keys))
	this.order = append(this.order, keys...)
	sort.Strings(this.order)

	this.data = make(map[string]*rex.Value, len(keys))
	for _, val := range keys {
		this.data[val] = input.Get(val)
	}
	return this
}

func NewOrderedMapFromLit(script *rex.Script, keys []string, values []*rex.Value) *rex.Value {
	if len(values) != 0 {
		if keys == nil {
			rex.RaiseError("Sort:map may not be initialized without keys.")
		}
	}
	if len(values) != len(keys) {
		rex.RaiseError("Sort:map: Key count does not equal value count.")
	}
	
	this := new(OrderedMap)
	this.order = make([]string, 0, len(keys))
	this.order = append(this.order, keys...)
	sort.Strings(this.order)
	
	this.data = make(map[string]*rex.Value, len(keys))
	for i, val := range keys {
		this.data[val] = values[i]
	}
	return rex.NewValueIndex(this)
}

func (this *OrderedMap) Get(index string) *rex.Value {
	this.lock.RLock()
	if _, ok := this.data[index]; ok {
		tmp := this.data[index]
		this.lock.RUnlock()
		return tmp
	}
	this.lock.RUnlock()
	return rex.NewValueInt64(0)
}

func (this *OrderedMap) Set(index string, value *rex.Value) bool {
	this.lock.Lock()
	if _, ok := this.data[index]; ok {
		this.data[index] = value
		this.lock.Unlock()
		return true
	}
	this.order = append(this.order, index)
	sort.Strings(this.order)
	this.data[index] = value
	this.lock.Unlock()
	return true
}

func (this *OrderedMap) Exists(index string) bool {
	_, ok := this.data[index]
	return ok
}

func (this *OrderedMap) Len() int64 {
	return int64(len(this.data))
}

func (this *OrderedMap) Keys() []string {
	return this.order
}

func (this *OrderedMap) String() string {
	return rex.IndexableToString("map", this)
}

func (this *OrderedMap) CodeString() string {
	return rex.IndexableToCodeString("map", this, true)
}
