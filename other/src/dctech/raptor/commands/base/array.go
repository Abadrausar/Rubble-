/*
Copyright 2012-2013 by Milo Christiansen

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

package base

import "dctech/raptor"
import "strconv"

type ScriptArray []*raptor.Value

func NewScriptArray() raptor.EditIndexable {
	tmp := make(ScriptArray, 0, 20)
	return &tmp
}

func NewScriptArraySized(size int) raptor.EditIndexable {
	tmp := make(ScriptArray, size)
	return &tmp
}

func NewScriptArrayFromLit(script *raptor.Script, keys []string, values []*raptor.Value) *raptor.Value {
	if keys != nil {
		panic("array may not be initalized with keys.")
	}

	tmp := make(ScriptArray, len(values))
	for i := range values {
		tmp[i] = values[i]
	}
	return raptor.NewValueObject(&tmp)
}

func (this *ScriptArray) Get(index string) *raptor.Value {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		panic("Index not a valid number.")
	}
	if val < 0 || val >= int64(len(*this)) {
		panic("Index out of range.")
	}
	return (*this)[val]
}

func (this *ScriptArray) Set(index string, value *raptor.Value) bool {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		panic("Index not a valid number.")
	}
	if val < 0 {
		panic("Index too small.")
	}
	if val < int64(len(*this)) {
		(*this)[val] = value
		return true
	}
	if val == int64(len(*this)) {
		*this = append(*this, value)
		return true
	}
	panic("Index too large, you may only extend the array by one at a time.")
}

func (this *ScriptArray) Exists(index string) bool {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		return false
	}
	if val < 0 || val >= int64(len(*this)) {
		return false
	}
	return true
}

func (this *ScriptArray) Len() int64 {
	return int64(len(*this))
}

func (this *ScriptArray) Keys() []string {
	rtn := make([]string, 0, len(*this))
	for key := range *this {
		rtn = append(rtn, strconv.FormatInt(int64(key), 10))
	}
	return rtn
}

func (this *ScriptArray) String() string {
	return raptor.GenericIndexableToString("array", this)
}

func (this *ScriptArray) CodeString() string {
	return raptor.GenericIndexableToCodeString("array", this, false)
}
