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
import "strconv"

type Array reflect.Value

// NewArray creates a new read/write array/slice interface and returns it as an EditIndexable.
// panics if data is not a array or slice.
func NewArray(data reflect.Value) rex.EditIndexable {
	ii := Array(data)
	
	if data.Kind() != reflect.Array && data.Kind() != reflect.Slice {
		rex.RaiseError("GenII: Non-array or slice type passed to NewArray.")
	}
	return &ii
}

func (ii *Array) Get(index string) *rex.Value {
	val, err := strconv.ParseInt(index, 0, 32)
	if err != nil {
		rex.RaiseError("GenII: Index not a valid number.")
	}
	i := int(val)
	
	length := reflect.Value(*ii).Len()
	if i < 0 || i >= length {
		return rex.NewValue()
	}
	
	rval := reflect.Value(*ii).Index(i)
	return RValueToSValue(rval)
}

func (ii *Array) GetI(index int64) *rex.Value {
	i := int(index)
	
	length := reflect.Value(*ii).Len()
	if i < 0 || i >= length {
		return rex.NewValue()
	}
	
	rval := reflect.Value(*ii).Index(i)
	return RValueToSValue(rval)
}

func (ii *Array) Set(index string, sval *rex.Value) bool {
	val, err := strconv.ParseInt(index, 0, 32)
	if err != nil {
		rex.RaiseError("GenII: Index not a valid number.")
	}
	i := int(val)
	
	length := reflect.Value(*ii).Len()
	if i < 0 || i >= length {
		return false
	}
	
	rval := reflect.Value(*ii).Index(i)
	SValueToRValue(sval, rval)
	return true
}

func (ii *Array) SetI(index int64, sval *rex.Value) bool {
	i := int(index)
	
	length := reflect.Value(*ii).Len()
	if i < 0 || i >= length {
		return false
	}
	
	rval := reflect.Value(*ii).Index(i)
	SValueToRValue(sval, rval)
	return true
}

func (ii *Array) Exists(index string) bool {
	val, err := strconv.ParseInt(index, 0, 32)
	if err != nil {
		return false
	}
	i := int(val)
	
	length := reflect.Value(*ii).Len()
	if i < 0 || i >= length {
		return false
	}
	return true
}

func (ii *Array) Len() int64 {
	return int64(reflect.Value(*ii).Len())
}

func (ii *Array) Keys() []string {
	length := reflect.Value(*ii).Len()
	rtn := make([]string, length)
	for i := 0; i < length; i++ {
		rtn[i] = strconv.FormatInt(int64(i), 10)
	}
	return rtn
}

func (ii *Array) String() string {
	//return rex.IndexableToString("GenII:Array", ii)
	return `<GenII:Array>`
}

func (ii *Array) CodeString() string {
	//return rex.IndexableToCodeString("array", ii, false)
	return `"<GenII:Array>"`
}
