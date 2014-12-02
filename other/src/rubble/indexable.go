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

package main

import "dctech/raptor"

// Indexable Raws

type IndexableRaws int

func NewIndexableRaws() raptor.EditIndexable {
	return new(IndexableRaws)
}

func (this IndexableRaws) Get(index string) *raptor.Value {
	if _, ok := Files.Files[index]; ok {
		if !Files.Files[index].Skip {
			return raptor.NewValueString(string(Files.Files[index].Content))
		}
	}
	return raptor.NewValueString("")
}

func (this IndexableRaws) Set(index string, value *raptor.Value) bool {
	if _, ok := Files.Files[index]; ok {
		if !Files.Files[index].Skip {
			Files.Files[index].Content = []byte(value.String())
			return true
		}
	}
	return false
}

func (this IndexableRaws) Exists(index string) bool {
	if _, ok := Files.Files[index]; ok {
		if !Files.Files[index].Skip {
			return true
		}
	}
	return false
}

func (this IndexableRaws) Len() int64 {
	var length int64
	for _, key := range Files.Order {
		if Files.Files[key].Skip {
			continue
		}
		length++
	}
	return length
}

func (this IndexableRaws) Keys() []string {
	rtn := make([]string, 0, this.Len())
	for _, key := range Files.Order {
		if Files.Files[key].Skip {
			continue
		}
		rtn = append(rtn, key)
	}
	return rtn
}

func (this IndexableRaws) String() string {
	return raptor.GenericIndexableToString("rubble:raws", this)
}

func (this IndexableRaws) CodeString() string {
	return raptor.GenericIndexableToCodeString("map", this, true)
}

// Read-only Map

type ReadOnlyMap map[string]*raptor.Value

func NewReadOnlyMap(keys []string, values []*raptor.Value) raptor.Indexable {
	this := make(ReadOnlyMap, len(keys))
	for i, key := range keys {
		this[key] = values[i]
	}
	return this
}

func (this ReadOnlyMap) Get(index string) *raptor.Value {
	if _, ok := this[index]; ok {
		return this[index]
	}
	return raptor.NewValue()
}

func (this ReadOnlyMap) Exists(index string) bool {
	_, ok := this[index]
	return ok
}

func (this ReadOnlyMap) Len() int64 {
	return int64(len(this))
}

func (this ReadOnlyMap) Keys() []string {
	rtn := make([]string, 0, len(this))
	for key := range this {
		rtn = append(rtn, key)
	}
	return rtn
}

func (this ReadOnlyMap) String() string {
	return raptor.GenericIndexableToString("map", this)
}

func (this ReadOnlyMap) CodeString() string {
	return raptor.GenericIndexableToCodeString("map", this, true)
}
