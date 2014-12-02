/*
For copyright/license see header in file "doc.go"
*/

package raptor

// CodeReader is a CodeSource wrapper for Code.
type CodeReader struct {
	currentToken int
	line         int
	data         *Code
}

// NewCodeReader creates a new CodeSource reading from a Code object.
func NewCodeReader(val *Code) CodeSource {
	this := new(CodeReader)
	this.data = val
	this.currentToken = -1
	return this
}

// CurrentTkn returns the current token.
func (this *CodeReader) CurrentTkn() *Token {
	return this.Token(this.currentToken)
}

// LookAhead returns the lookahead token.
func (this *CodeReader) LookAhead() *Token {
	return this.Token(this.currentToken + 1)
}

// Advance retrieves the next token from the stream.
func (this *CodeReader) Advance() {
	this.currentToken++
	// This little bit here is just to keep the token number from growing too much.
	if this.currentToken > len(this.data.Code) {
		this.currentToken = len(this.data.Code)
	}
}

// Used by Token to generate lexemes.
var tknLexemes = []string{
	"(",
	")",
	"[",
	"]",
	"<",
	">",
	"=",
	"{",
	"}"}

// Token creates a Token from the data at index.
// Returns TknINVALID if the index is invalid.
func (this *CodeReader) Token(index int) *Token {
	if index >= len(this.data.Code) || index < 0 {
		return NewToken("INVALID", TknINVALID, NewPosition(0, -1, ""))
	}
	code := this.data.Code[index]

	if code < 9 {
		if this.data.Positions == nil {
			return NewToken(tknLexemes[code], int(code), NewPosition(index, -1, ""))
		}
		return NewToken(tknLexemes[code], int(code), this.data.Positions[index])
	}

	if code == 9 {
		// Unused token code
		if this.data.Positions == nil {
			return NewToken("INVALID", TknINVALID, NewPosition(0, -1, ""))
		}
		return NewToken("INVALID", TknINVALID, this.data.Positions[index])
	}

	code = code - 10
	if int(code) < len(this.data.StringTable) && code >= 0 {
		if this.data.Positions == nil {
			return NewToken(this.data.StringTable[code], TknString, NewPosition(index, -1, ""))
		}
		return NewToken(this.data.StringTable[code], TknString, this.data.Positions[index])
	}

	// String token out of bounds.
	if this.data.Positions == nil {
		return NewToken("INVALID", TknINVALID, NewPosition(0, -1, ""))
	}
	return NewToken("INVALID", TknINVALID, this.data.Positions[index])
}

// ValidatePosition will return true if index is a valid token code.
func (this *CodeReader) ValidatePosition(index int) bool {
	if index >= len(this.data.Code) || index < 0 {
		return false
	}
	code := this.data.Code[index]

	if code < 9 {
		return true
	}

	if code == 9 {
		// Unused token code
		return false
	}

	code = code - 10
	if int(code) < len(this.data.StringTable) && code >= 0 {
		return true
	}

	return false
}
