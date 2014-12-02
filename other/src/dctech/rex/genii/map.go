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

package genii

import "dctech/rex"
import "reflect"

type Map reflect.Value

// NewMap creates a new read/write map interface and returns it as an EditIndexable.
// panics if data is not a map with a string key type.
func NewMap(data reflect.Value) rex.EditIndexable {
	ii := Map(data)
	
	if data.Kind() != reflect.Map {
		rex.RaiseError("GenII: Non-map type passed to NewMap.")
	}
	
	if data.Type().Key().Kind() != reflect.String {
		rex.RaiseError("GenII: Map type does not have string keys.")
	}
	return &ii
}

func (ii *Map) Get(index string) *rex.Value {
	key := reflect.ValueOf(index)
	
	rval := reflect.Value(*ii).MapIndex(key)
	zero := reflect.Value{}
	if rval == zero {
		rex.RaiseError("GenII: Index does not exist.")
	}
	return RValueToSValue(rval)
}

func (ii *Map) Set(index string, sval *rex.Value) bool {
	key := reflect.ValueOf(index)
	
	rval := reflect.Value(*ii).MapIndex(key)
	zero := reflect.Value{}
	if rval == zero {
		rex.RaiseError("GenII: Index does not exist.")
	}
	SValueToRValue(sval, rval)
	return true
}

func (ii *Map) Exists(index string) bool {
	key := reflect.ValueOf(index)
	rval := reflect.Value(*ii).MapIndex(key)
	zero := reflect.Value{}
	if rval == zero {
		return false
	}
	return true
}

func (ii *Map) Len() int64 {
	return int64(reflect.Value(*ii).Len())
}

func (ii *Map) Keys() []string {
	keysR := reflect.Value(*ii).MapKeys()
	keys := make([]string, len(keysR))
	for i := range keysR {
		keys[i] = keysR[i].String()
	}
	return keys
}

func (ii *Map) String() string {
	//return rex.IndexableToString("GenII:Map", ii)
	return `<GenII:Map>`
}

func (ii *Map) CodeString() string {
	//return rex.IndexableToCodeString("map", ii, true)
	return `"<GenII:Map>"`
}

