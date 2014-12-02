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
