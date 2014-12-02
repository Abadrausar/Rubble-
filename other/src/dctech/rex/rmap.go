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

type rMap map[string]*Value

// NewStaticMap creates a new read only script map.
func NewStaticMap(data map[string]*Value) Indexable {
	return rMap(data)
}

// NewStaticMapFromLit is an ObjectFactory for read only script maps.
func NewStaticMapFromLit(script *Script, keys []string, values []*Value) *Value {
	if len(values) != 0 {
		if keys == nil {
			RaiseError("Smap may not be initialized without keys.")
		}
	}
	if len(values) != len(keys) {
		RaiseError("Smap: Key count does not equal value count.")
	}

	tmp := make(rMap, len(keys))
	for i := range keys {
		tmp[keys[i]] = values[i]
	}
	return NewValueIndex(&tmp)
}

func (this rMap) Get(index string) *Value {
	if _, ok := this[index]; ok {
		return this[index]
	}
	return NewValue()
}

func (this rMap) Exists(index string) bool {
	_, ok := this[index]
	return ok
}

func (this rMap) Len() int64 {
	return int64(len(this))
}

func (this rMap) Keys() []string {
	rtn := make([]string, 0, len(this))
	for key := range this {
		rtn = append(rtn, key)
	}
	return rtn
}

func (this rMap) String() string {
	return IndexableToString("map", this)
}
