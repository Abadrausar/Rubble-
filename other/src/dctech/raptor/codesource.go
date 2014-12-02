/*
For copyright/license see header in file "doc.go"
*/

package raptor

//import "fmt"

// A CodeSource represents a token stream.
type CodeSource interface {
	CurrentTkn() *Token
	LookAhead() *Token
	Advance()
}

// GetToken gets the next token, and panics with an error if it's not of type tokenType.
// May cause a panic if the lexer encounters an error
// Used as a type checked Advance
func GetToken(code CodeSource, tokenTypes ...int) {
	code.Advance()

	for _, val := range tokenTypes {
		if code.CurrentTkn().Type == val {
			return
		}
	}

	ExitOnTokenExpected(code.CurrentTkn(), tokenTypes...)
}

// CheckLookAhead checks to see if the look ahead is one of tokenTypes and if so returns true
func CheckLookAhead(code CodeSource, tokenTypes ...int) bool {
	for _, val := range tokenTypes {
		if code.LookAhead().Type == val {
			return true
		}
	}
	return false
}
