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
