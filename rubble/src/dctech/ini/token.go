package ini

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
