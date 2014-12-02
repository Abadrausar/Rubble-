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

import "dctech/nca7"

type ScriptMap map[string]*nca7.Value

func NewScriptMap() nca7.EditIndexable {
	return make(ScriptMap, 20)
}

func NewScriptMapFromLit(keys []string, values []*nca7.Value) *nca7.Value {
	if len(values) != 0 {
		if keys == nil {
			panic("Map may not be initalized without keys.")
		}
	}
	
	tmp := make(ScriptMap, len(keys))
	for i := range keys {
		tmp[keys[i]] = values[i]
	}
	return nca7.NewValueObject(&tmp)
}

func (this ScriptMap) Get(index string) *nca7.Value {
	if _, ok := this[index]; ok {
		return this[index]
	}
	return nca7.NewValueInt64(0)
}

func (this ScriptMap) Set(index string, value *nca7.Value) bool {
	this[index] = value
	return true
}

func (this ScriptMap) Exists(index string) bool {
	_, ok := this[index]
	return ok
}

func (this ScriptMap) Len() int64 {
	return int64(len(this))
}

func (this ScriptMap) Keys() []string {
	rtn := make([]string, 0, len(this))
	for key := range this {
		rtn = append(rtn, key)
	}
	return rtn
}
