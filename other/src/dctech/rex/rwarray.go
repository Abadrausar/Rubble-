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

import "strconv"
import "sync"

type rwArray struct {
	data []*Value
	lock *sync.RWMutex
}

// NewArray creates a new script array.
func NewArray(data []*Value) EditIndexable {
	return &rwArray{
		data: data,
		lock: new(sync.RWMutex),
	}
}

// NewArrayFromLit is an ObjectFactory for script arrays.
func NewArrayFromLit(script *Script, keys []string, values []*Value) *Value {
	if keys != nil {
		RaiseError("Array may not be initialized with keys.")
	}

	return NewValueIndex(&rwArray{
		data: values,
		lock: new(sync.RWMutex),
	})
}

func (array *rwArray) Get(index string) *Value {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		RaiseError("Index not a valid number.")
	}
	
	array.lock.RLock()
	defer array.lock.RUnlock()
	
	if val < 0 || val >= int64(len(array.data)) {
		RaiseError("Index out of range.")
	}
	return array.data[val]
}

func (array *rwArray) Set(index string, value *Value) bool {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		if index == "append" {
			array.lock.Lock()
			array.data = append(array.data, value)
			array.lock.Unlock()
			return true
		}
		RaiseError("Index not a valid number.")
	}
	if val < 0 {
		RaiseError("Index out of bounds: Too small.")
	}
	
	array.lock.Lock()
	defer array.lock.Unlock()
	
	if val < int64(len(array.data)) {
		array.data[val] = value
		return true
	}
	RaiseError("Index out of bounds: Too large.")
	panic("UNREACHABLE")
}

func (array *rwArray) Exists(index string) bool {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		return false
	}
	if val < 0 || val >= int64(len(array.data)) {
		return false
	}
	return true
}

func (array *rwArray) Len() int64 {
	return int64(len(array.data))
}

func (array *rwArray) Keys() []string {
	array.lock.RLock()
	defer array.lock.RUnlock()
	
	rtn := make([]string, 0, len(array.data))
	for key := range array.data {
		rtn = append(rtn, strconv.FormatInt(int64(key), 10))
	}
	return rtn
}

func (array *rwArray) String() string {
	array.lock.RLock()
	defer array.lock.RUnlock()
	
	return IndexableToString("array", array)
}
