/*
For copyright/license see header in file "doc.go"
*/

package raptor

import "fmt"

// Position stores line, column, and file information for a Token or Value.
// If the column is -1 then the line number is a token offset instead.
// Be careful about changing values after creation, these things get passed and stored all over the place.
type Position struct {
	Line   int
	Column int
	file   *string
}

// NewPosition creates and returns a new Position object.
// To create a token offset Position pass -1 as the column and the offset as the line.
func NewPosition(line, column int, file string) *Position {
	return newPosition(line, column, &file)
}

// For internal use only.
func newPosition(line, column int, file *string) *Position {
	this := new(Position)
	this.Line = line
	this.Column = column
	this.file = file
	return this
}

// File returns the name of the file this position refers to (or "").
func (this *Position) File() string {
	return *this.file
}

// Copy copies a Position object, more useful than it sounds.
func (this *Position) Copy() *Position {
	return newPosition(this.Line, this.Column, this.file)
}

// String returns strings of one of the following forms:
//	"somefile.rsf|L:x|C:y"
//	"somefile.rbf|T:x"
//	"L:x|C:y"
//	"T:x"
func (this *Position) String() string {
	out := ""

	if *this.file != "" {
		out += *this.file + "|"
	}

	if this.Column == -1 {
		out += fmt.Sprintf("T:%v", this.Line)
	}
	out += fmt.Sprintf("L:%v|C:%v", this.Line, this.Column)

	return out
}
