/*
Copyright 2013 by Milo Christiansen

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

package ini

import "fmt"

type Token struct {
	Lexeme string
	Type   int
	Line   int
}

// Returns the type of Token as a string value.
// Used for debugging and error reporting
func (this *Token) TypeToString() string {
	return tokenTypeToString(this.Type)
}

func (this *Token) String() string {
	return fmt.Sprintf("Lexeme: %s Type: %s Line: %d", this.Lexeme, this.TypeToString(), this.Line)
}

func tokenTypeToString(tokenType int) string {
	switch tokenType {
	case TknINVALID:
		return "TknINVALID"
	case TknKey:
		return "TknKey"
	case TknValue:
		return "TknValue"
	case TknSection:
		return "TknSection"
	case TknComment:
		return "TknComment"
	}

	panic("Token type value out of range")
}

// Panics with the message: 
//	Error: Invalid Token: Found: thecurrenttoken. Expected: expected.
func exitOnTokenExpected(token Token, expected ...int) {
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

	panic(newErr(TokenExpected, token.TypeToString(), expectedString, token.Line))
}
