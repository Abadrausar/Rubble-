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

type Struct reflect.Value

// NewStruct creates a new read/write struct interface and returns it as an EditIndexable.
// panics if data is not a struct.
func NewStruct(data reflect.Value) rex.EditIndexable {
	ii := Struct(data)
	
	if data.Kind() != reflect.Struct {
		rex.RaiseError("GenII: Non-struct type passed to NewStruct.")
	}
	return &ii
}

func (ii *Struct) Get(index string) *rex.Value {
	if _, ok := reflect.Value(*ii).Type().FieldByName(index); !ok {
		return rex.NewValue()
	}
	
	rval := reflect.Value(*ii).FieldByName(index)
	return RValueToSValue(rval)
}

func (ii *Struct) Set(index string, sval *rex.Value) bool {
	if _, ok := reflect.Value(*ii).Type().FieldByName(index); !ok {
		rex.RaiseError("GenII: Field: " + index + " Does not exist in struct.")
	}
	
	rval := reflect.Value(*ii).FieldByName(index)
	SValueToRValue(sval, rval)
	return true
}

func (ii *Struct) Exists(index string) bool {
	_, ok := reflect.Value(*ii).Type().FieldByName(index)
	return ok
}

func (ii *Struct) Len() int64 {
	return int64(reflect.Value(*ii).NumField())
}

func (ii *Struct) Keys() []string {
	fields := reflect.Value(*ii).NumField()
	rtn := make([]string, fields)
	typ := reflect.Value(*ii).Type()
	
	for i := 0; i < fields; i++ {
		rtn[i] = typ.Field(i).Name
	}
	
	return rtn
}

func (ii *Struct) String() string {
	//return rex.IndexableToString("GenII:Struct", ii)
	return `<GenII:Struct>`
}

func (ii *Struct) CodeString() string {
	//return rex.IndexableToCodeString("map", ii, true)
	return `"<GenII:Struct>"`
}

