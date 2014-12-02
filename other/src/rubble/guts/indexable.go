/*
Copyright 2013-2014 by Milo Christiansen

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

package guts

import "dctech/rex"

// Indexable Raws

type IndexableRaws int

func NewIndexableRaws() rex.EditIndexable {
	return new(IndexableRaws)
}

func (this IndexableRaws) Get(index string) *rex.Value {
	if _, ok := Files.Data[index]; ok {
		if !Files.Data[index].Tags["Skip"] {
			return rex.NewValueString(string(Files.Data[index].Content))
		}
	}
	return rex.NewValue()
}

func (this IndexableRaws) Set(index string, value *rex.Value) bool {
	if _, ok := Files.Data[index]; ok {
		if !Files.Data[index].Tags["Skip"] {
			Files.Data[index].Content = []byte(value.String())
			return true
		}
	}
	return false
}

func (this IndexableRaws) Exists(index string) bool {
	if _, ok := Files.Data[index]; ok {
		if !Files.Data[index].Tags["Skip"] {
			return true
		}
	}
	return false
}

func (this IndexableRaws) Len() int64 {
	var length int64
	for _, key := range Files.Order {
		if Files.Data[key].Tags["Skip"] {
			continue
		}
		length++
	}
	return length
}

func (this IndexableRaws) Keys() []string {
	rtn := make([]string, 0, this.Len())
	for _, key := range Files.Order {
		if Files.Data[key].Tags["Skip"] {
			continue
		}
		rtn = append(rtn, key)
	}
	return rtn
}

func (this IndexableRaws) String() string {
	return rex.IndexableToString("rubble:raws", this)
}

func (this IndexableRaws) CodeString() string {
	return rex.IndexableToCodeString("map", this, true)
}
