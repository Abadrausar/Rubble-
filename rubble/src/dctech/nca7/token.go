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

package nca7

import "fmt"

type Token struct {
	Lexeme string
	Type   int
	Line   int
	Column int
}

func (this *Token) String() string {
	return tokenTypeToString(this.Type)
}

func tokenTypeToString(tokenType int) string {
	switch tokenType {
	case TknINVALID:
		return "TknINVALID"
	case TknString:
		return "TknString"
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
	}

	panic("Token type value out of range")
}

// Panics with the message: 
//	Error: Invalid Token: Found: thecurrenttoken. Expected: expected. On Line: x Column: y
func ExitOnTokenExpected(token *Token, expected ...int) {
	expectedString := ""
	doonce := true
	for _, val := range expected {
		if doonce {
			doonce = false
			expectedString = tokenTypeToString(val)
			continue
		}
		expectedString += " or " + tokenTypeToString(val)
	}

	panic(fmt.Sprintf("Invalid Token: Found: %s. Expected: %s.", token.String(), expectedString))
}
