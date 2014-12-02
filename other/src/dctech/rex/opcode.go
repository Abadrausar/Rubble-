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

// opcodes type values.
const (
	opINVALID     = iota - 1 // Invalid
	opCmdBegin               // '('
	opCmdEnd                 // ')'
	opVarBegin               // '['
	opVarEnd                 // ']'
	opObjLitBegin            // '<'
	opObjLitEnd              // '>'
	opNameSplit              // ':'
	opAssignment             // '='

	opValue // A literal value (including code blocks)
	opName  // Any name, checked at runtime
)

// opCode is used to represent compiled code.
// All values are encapsulated and all names are reduced to their local indexes.
type opCode struct {
	Type  int
	Index int    // Type == opName, if index is an external local then the true index is 0-(index+1)
	Value *Value // Type == opValue
	Pos   *Position
}

// String returns a string representation of the opCode's type.
// If you want a parseable/more compact result use toLexeme.
func (op *opCode) String() string {
	return opCodeTypeToString(op.Type)
}

func opCodeTypeToString(typ int) string {
	switch typ {
	case opINVALID:
		return "opINVALID"
	case opCmdBegin:
		return "opCmdBegin"
	case opCmdEnd:
		return "opCmdEnd"
	case opVarBegin:
		return "opVarBegin"
	case opVarEnd:
		return "opVarEnd"
	case opObjLitBegin:
		return "opObjLitBegin"
	case opObjLitEnd:
		return "opObjLitEnd"
	case opNameSplit:
		return "opNameSplit"
	case opAssignment:
		return "opAssignment"
	case opValue:
		return "opValue"
	case opName:
		return "opName"
	}

	RaiseError("Type value out of range")
	panic("UNREACHABLE")
}

// toLexeme is a helper function for converting code values to strings.
// Does not return parseable values for type opName and type opValue.
func (op *opCode) toLexeme() string {
	switch op.Type {
	case opINVALID:
		return "\"INVALID_OPCODE\""
	case opCmdBegin:
		return "("
	case opCmdEnd:
		return ")"
	case opVarBegin:
		return "["
	case opVarEnd:
		return "]"
	case opObjLitBegin:
		return "<"
	case opObjLitEnd:
		return ">"
	case opNameSplit:
		return ":"
	case opAssignment:
		return "="
	case opValue:
		return op.Value.String()
	case opName:
		return fmt.Sprint("name(", op.Index, ")")
	}
	RaiseError("Type value out of range")
	panic("UNREACHABLE")
}

// Panics with a message formatted like the following:
//	Invalid opCode: Found: thecurrent. Expected: expected1, expected2, or expected3.
//	Invalid opCode: Found: thecurrent. Expected: expected1 or expected2.
//	Invalid opCode: Found: thecurrent. Expected: expected.
func exitOnopCodeExpected(op *opCode, expected ...int) {
	expectedString := ""
	expectedCount := len(expected) - 1
	for i, val := range expected {
		// Is the only value
		if expectedCount == 0 {
			expectedString = opCodeTypeToString(val)
			continue
		}

		// Is last of a list (2 or more)
		if i == expectedCount && expectedCount > 0 {
			expectedString += "or " + opCodeTypeToString(val)
			continue
		}

		// Is the first of two
		if expectedCount == 1 {
			expectedString += opCodeTypeToString(val) + " "
			continue
		}

		// Is any but the last of a list of 3 or more
		expectedString += opCodeTypeToString(val) + ", "
	}

	RaiseError("Invalid opCode: Found: " + op.String() + ". Expected: " + expectedString)
	panic("UNREACHABLE")
}
