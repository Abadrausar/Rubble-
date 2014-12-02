/*
Copyright 2013 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. You may not alter this software in any way.

3. This notice may not be removed or altered from any source distribution.
*/

package main

import "fmt"

// Token Types
const (
	tknINVALID = iota
	tknString
	tknTagBegin
	tknTagEnd
	tknDelimiter
)

type Token struct {
	Lexeme string
	Type   int
	Line   int
}

func (this *Token) String() string {
	return tokenTypeToString(this.Type)
}

func tokenTypeToString(tokenType int) string {
	switch tokenType {
	case tknINVALID:
		return "TknINVALID"
	case tknString:
		return "TknString"
	case tknTagBegin:
		return "tknTagBegin"
	case tknTagEnd:
		return "tknTagEnd"
	case tknDelimiter:
		return "tknDelimiter"
	}

	panic("Token type value out of range")
}

// Panics with the message: 
//	Invalid Token: Found: thecurrenttoken. Expected: expected.
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
