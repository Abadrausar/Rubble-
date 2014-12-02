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

// Value is string+position.
type Value struct {
	Data string
	Pos  *Position
}

func NewValue(val string) *Value {
	this := new(Value)
	this.Data = val
	this.Pos = NewPosition(1, "")
	return this
}

func NewValuePos(val string, pos *Position) *Value {
	this := new(Value)
	this.Data = val
	this.Pos = pos
	return this
}

func NewValueRaptor(val *raptor.Value) *Value {
	this := new(Value)
	this.Data = val.String()
	this.Pos = NewPositionRaptor(val.Pos)
	return this
}

// Raptor copies a Value object into the Raptor equivalent.
func (this *Value) Raptor() *raptor.Value {
	rtn := raptor.NewValueString(this.Data)
	rtn.Pos = this.Pos.Raptor()
	return rtn
}
