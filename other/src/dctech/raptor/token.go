/*
For copyright/license see header in file "doc.go"
*/

package raptor

import "fmt"

// Token values
// THESE ARE ORDER/VALUE SENSITIVE!!!!!!
const (
	TknINVALID     = iota - 1 // Invalid
	TknCmdBegin               // (
	TknCmdEnd                 // )
	TknDerefBegin             // [
	TknDerefEnd               // ]
	TknObjLitBegin            // <
	TknObjLitEnd              // >
	TknObjLitSplit            // =
	TknCodeBegin              // {
	TknCodeEnd                // }
	TknString                 // A string, this needs to come last.
)

// Used to generate lexemes from token numbers.
// There is no entry for TknINVALID or TknString.
var TknLexemes = []string{
	"(",
	")",
	"[",
	"]",
	"<",
	">",
	"=",
	"{",
	"}"}

type Token struct {
	Lexeme string
	Type   int
	Pos    *Position
}

// NewToken is for use by CodeSource implementations.
// To help streamline things pos is copied, not used directly.
func NewToken(lexeme string, typ int, pos *Position) *Token {
	this := new(Token)
	this.Lexeme = lexeme
	this.Type = typ
	this.Pos = pos.Copy()
	return this
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
