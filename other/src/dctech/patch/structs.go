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

package patch

import "strconv"
import "bytes"

type Operation int8

const (
	DiffDelete Operation = -1
	DiffInsert Operation = 1
	DiffEqual  Operation = 0
)

// Diff represents one diff operation
type Diff struct {
	Type Operation
	Text string
}

// Patch represents one patch operation.
type Patch struct {
	diffs   []*Diff
	start1  int
	start2  int
	length1 int
	length2 int
}

// String converts a single patch to a string.
func (p *Patch) String() string {
	coords1 := strconv.Itoa(p.start1 + 1) + "," + strconv.Itoa(p.length1)
	coords2 := strconv.Itoa(p.start2 + 1) + "," + strconv.Itoa(p.length2)
	
	var text bytes.Buffer
	text.WriteString("@@ -" + coords1 + " +" + coords2 + " @@\n")

	// Escape the body of the patch with %xx notation.
	for _, aDiff := range p.diffs {
		switch aDiff.Type {
		case DiffInsert:
			text.WriteString("+")
		case DiffDelete:
			text.WriteString("-")
		case DiffEqual:
			text.WriteString(" ")
		}

		text.WriteString(aDiff.Text)
		text.WriteString("\n")
	}

	return text.String()
}
