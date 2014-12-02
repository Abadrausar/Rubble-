package ncalex

//Copyright 2013 by Milo Christiansen
//
//DCTech Project License
//	You may not modify the project without prior written permission.
//	You may not redistribute the project without prior written permission.
//	
//	The project is provided 'as-is', without any express or implied warranty. 
//	In no event will the author be held liable for any damages arising from the use of the project.
//	
//	In the event of a dispute about the interpretation of this license you agree that the 
//	interpretation of the author is the correct one.
// 
// Permission is given to include this package in your projects so long as the above license 
// and this exception are reproduced somewhere in your project.
// 

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
	case TknOpenParen:
		return "TknOpenParen"
	case TknCloseParen:
		return "TknCloseParen"
	case TknOpenSqBracket:
		return "TknOpenSqBracket"
	case TknCloseSqBracket:
		return "TknCloseSqBracket"
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
