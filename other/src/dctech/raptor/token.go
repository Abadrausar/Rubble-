/*
Copyright 2012-2013 by Milo Christiansen

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

package raptor

import "fmt"

type Token struct {
	Lexeme string
	Type   int
	Pos    *PositionInfo
}

func (this *Token) String() string {
	return tokenTypeToString(this.Type)
}

func tokenTypeToString(tokenType int) string {
	switch tokenType {
	case TknINVALID:
		return "TknINVALID"
	case TknCmdBegin:
		return "TknCmdBegin"
	case TknCmdEnd:
		return "TknCmdEnd"
	case TknDerefBegin:
		return "TknDerefBegin"
	case TknDerefEnd:
		return "TknDerefEnd"
	case TknObjLitBegin:
		return "TknObjLitBegin"
	case TknObjLitEnd:
		return "TknObjLitEnd"
	case TknObjLitSplit:
		return "TknObjLitSplit"
	case TknCodeBegin:
		return "TknCodeBegin"
	case TknCodeEnd:
		return "TknCodeEnd"
	case TknString:
		return "TknString"
	}

	panic("Token type value out of range")
}

// Panics with a message formatted like the following: 
//	Invalid Token: Found: thecurrenttoken. Expected: expected1, expected2, or expected3.
//	Invalid Token: Found: thecurrenttoken. Expected: expected1 or expected2.
//	Invalid Token: Found: thecurrenttoken. Expected: expected.
func ExitOnTokenExpected(token *Token, expected ...int) {
	expectedString := ""
	expectedCount := len(expected) - 1
	for i, val := range expected {
		// Is the only value
		if expectedCount == 0 {
			expectedString = tokenTypeToString(val)
			continue
		}

		// Is last of a list (2 or more)
		if i == expectedCount && expectedCount > 0 {
			expectedString += "or " + tokenTypeToString(val)
			continue
		}

		// Is the first of two
		if expectedCount == 1 {
			expectedString += tokenTypeToString(val) + " "
			continue
		}

		// Is any but the last of a list of 3 or more
		expectedString += tokenTypeToString(val) + ", "
	}

	panic(fmt.Sprintf("Invalid Token: Found: %s. Expected: %s", token.String(), expectedString))
}
