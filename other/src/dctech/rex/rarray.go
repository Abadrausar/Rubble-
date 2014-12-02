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

import "strconv"

type rArray []*Value

// NewStaticArray creates a new read only script array.
func NewStaticArray(data []*Value) Indexable {
	tmp := rArray(data)
	return &tmp
}

// NewStaticArray is an ObjectFactory for read only script arrays.
func NewStaticArrayFromLit(script *Script, keys []string, values []*Value) *Value {
	if keys != nil {
		RaiseError("Sarray may not be initialized with keys.")
	}

	tmp := make(rArray, len(values))
	for i := range values {
		tmp[i] = values[i]
	}
	return NewValueIndex(&tmp)
}

func (this *rArray) Get(index string) *Value {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		RaiseError("Index not a valid number.")
	}
	if val < 0 {
		RaiseError("Index out of range: Too small.")
	}
	if val >= int64(len(*this)) {
		RaiseError("Index out of range: Too large.")
	}
	return (*this)[val]
}

func (this *rArray) Exists(index string) bool {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		return false
	}
	if val < 0 || val >= int64(len(*this)) {
		return false
	}
	return true
}

func (this *rArray) Len() int64 {
	return int64(len(*this))
}

func (this *rArray) Keys() []string {
	rtn := make([]string, 0, len(*this))
	for key := range *this {
		rtn = append(rtn, strconv.FormatInt(int64(key), 10))
	}
	return rtn
}

func (this *rArray) String() string {
	return IndexableToString("sarray", this)
}

func (this *rArray) CodeString() string {
	return IndexableToCodeString("sarray", this, false)
}
