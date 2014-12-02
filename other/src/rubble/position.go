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

import "fmt"
import "dctech/raptor"

// Position stores line and file information.
// Be careful about changing values after creation, these things get passed and stored all over the place.
type Position struct {
	Line   int
	File   string
}

// NewPosition creates and returns a new Position object.
func NewPosition(line int, file string) *Position {
	this := new(Position)
	this.Line = line
	this.File = file
	return this
}

// NewPosition creates and returns a new Position object.
func NewPositionRaptor(pos *raptor.Position) *Position {
	this := new(Position)
	this.Line = pos.Line
	if pos.Column == -1 {
		this.Line = 1
	}
	this.File = pos.File()
	return this
}

// Copy copies a Position object, more useful than it sounds.
func (this *Position) Copy() *Position {
	return NewPosition(this.Line, this.File)
}

// Raptor copies a Position object into the Raptor equivalent.
func (this *Position) Raptor() *raptor.Position {
	return raptor.NewPosition(this.Line, 1, this.File)
}

// String returns strings of one of the following forms: 
//	"somefile|L:x"
//	"L:x"
func (this *Position) String() string {
	out := ""
	
	if this.File != "" {
		out += this.File + "|"
	}
	
	out += fmt.Sprintf("L:%v", this.Line)
	
	return out
}
