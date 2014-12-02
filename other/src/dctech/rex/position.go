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

import "fmt"

// Position stores line, column, and file information.
// Be careful about changing values after creation, these things get passed and stored all over the place.
type Position struct {
	Line   int
	Column int
	file   *string
}

// NewPosition creates and returns a new Position object.
func NewPosition(line, column int, file string) *Position {
	return &Position{
		Line: line,
		Column: column,
		file: &file,
	}
}

// File returns the name of the file this position refers to (or "").
func (pos *Position) File() string {
	return *pos.file
}

// Copy copies a Position object.
// More useful than it sounds.
func (pos *Position) Copy() *Position {
	return &Position{
		Line: pos.Line,
		Column: pos.Column,
		file: pos.file,
	}
}

// String returns strings of one of the following forms:
//	"somefile.rex|L:x|C:y"
//	"L:x|C:y"
func (pos *Position) String() string {
	if *pos.file == "" {
		return fmt.Sprintf("L:%v|C:%v", pos.Line, pos.Column)
	}
	return fmt.Sprintf("%v|L:%v|C:%v", *pos.file, pos.Line, pos.Column)
}
