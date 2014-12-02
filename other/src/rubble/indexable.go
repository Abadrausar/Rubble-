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

package rubble

import "dctech/rex"

// Indexable Raws

type IndexableRaws struct {
	files **FileList
}

func NewIndexableRaws(list **FileList) rex.EditIndexable {
	return &IndexableRaws{list}
}

func (raws *IndexableRaws) Get(index string) *rex.Value {
	if _, ok := (*raws.files).Data[index]; ok {
		if !(*raws.files).Data[index].Tags["Skip"] {
			return rex.NewValueString(string((*raws.files).Data[index].Content))
		}
	}
	return rex.NewValue()
}

func (raws *IndexableRaws) Set(index string, value *rex.Value) bool {
	if _, ok := (*raws.files).Data[index]; ok {
		if !(*raws.files).Data[index].Tags["Skip"] {
			(*raws.files).Data[index].Content = []byte(value.String())
			return true
		}
	}
	return false
}

func (raws *IndexableRaws) Exists(index string) bool {
	if _, ok := (*raws.files).Data[index]; ok {
		if !(*raws.files).Data[index].Tags["Skip"] {
			return true
		}
	}
	return false
}

func (raws *IndexableRaws) Len() int64 {
	var length int64
	for _, key := range (*raws.files).Order {
		if (*raws.files).Data[key].Tags["Skip"] {
			continue
		}
		length++
	}
	return length
}

func (raws *IndexableRaws) Keys() []string {
	rtn := make([]string, 0, raws.Len())
	for _, key := range (*raws.files).Order {
		if (*raws.files).Data[key].Tags["Skip"] {
			continue
		}
		rtn = append(rtn, key)
	}
	return rtn
}

func (raws *IndexableRaws) String() string {
	return rex.IndexableToString("rubble:raws", raws)
}
