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

// New creates a new GenII interface value, unconvertible types are converted to Rex nil.
func New(obj interface{}) (sval *rex.Value) {
	rval := reflect.ValueOf(obj)
	return RValueToSValue(rval)
}

// RValueToSValue converts a reflect.Value to a rex.Value, unconvertible types are converted to Rex nil.
// If possible types that satisfy the rex.Indexable interface are returned raw.
func RValueToSValue(rval reflect.Value) *rex.Value {
	if rval.IsValid() && rval.CanInterface() {
		val := rval.Interface()
		if index, ok := val.(rex.Indexable); ok {
			return rex.NewValueIndex(index)
		}
	}
	
	switch rval.Kind() {
	case reflect.String:
		return rex.NewValueString(rval.String())
	
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rex.NewValueInt64(int64(rval.Uint()))
	
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rex.NewValueInt64(rval.Int())
	
	case reflect.Float32, reflect.Float64:
		return rex.NewValueFloat64(rval.Float())
	
	case reflect.Bool:
		return rex.NewValueBool(rval.Bool())
	
	case reflect.Array, reflect.Slice:
		return rex.NewValueIndex(NewArray(rval))
	
	case reflect.Map:
		return rex.NewValueIndex(NewMap(rval))
	
	case reflect.Struct:
		return rex.NewValueIndex(NewStruct(rval))
	
	case reflect.Ptr, reflect.Interface:
		return RValueToSValue(rval.Elem())
	
	default:
		return rex.NewValue()
	}
}

// SValueToRValue converts a rex.Value to a reflect.Value, how it is converted is determined by the
// passed in reflect.Value.
// Panics if the value cannot be stored in the reflect.Value for some reason.
func SValueToRValue(sval *rex.Value, rval reflect.Value) {
	if !rval.CanSet() {
		rex.RaiseError("GenII: Value passed to SValueToRValue is not settable.")
	}
	
	if sval.Type == rex.TypUser {
		if reflect.TypeOf(sval.Data) == rval.Type() {
			rval.Set(reflect.ValueOf(sval.Data))
			return
		}
	}
	
	switch rval.Kind() {
	case reflect.String:
		rval.SetString(sval.String())
	
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		rval.SetUint(uint64(sval.Int64()))
	
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rval.SetInt(sval.Int64())
	
	case reflect.Float32, reflect.Float64:
		rval.SetFloat(sval.Float64())
	
	case reflect.Bool:
		rval.SetBool(sval.Bool())
	
	case reflect.Array, reflect.Slice:
		rex.RaiseError("GenII: Cannot convert to arrays and slices yet.")
	
	case reflect.Map:
		rex.RaiseError("GenII: Cannot convert to maps yet.")
	
	case reflect.Struct:
		rex.RaiseError("GenII: Cannot convert to structs yet.")
	
	case reflect.Ptr, reflect.Interface:
		SValueToRValue(sval, rval.Elem())
	
	default:
		rex.RaiseError("GenII: Value passed to SValueToRValue has a type that Rex cannot convert a value to.")
	}
	
	return
}
