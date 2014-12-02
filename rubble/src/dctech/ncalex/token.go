package ncalex

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
