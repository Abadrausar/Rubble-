/*
Copyright 2012-2013 by Milo Christiansen

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

package raptor

//import "fmt"

// The Raptor version of this is simpler than the old NCA version
// A much needed refactoring as well as better token type constants did wonders.

// CompiledScript is a pre-lexed Raptor script.
// It is valid for the position information to be nil.
type CompiledScript struct {
	// Code Values:
	//	0 = TknCmdBegin
	//	1 = TknCmdEnd
	//	2 = TknDerefBegin
	//	3 = TknDerefEnd
	//	4 = TknObjLitBegin
	//	5 = TknObjLitEnd
	//	6 = TknObjLitSplit
	//	7 = TknCodeBegin
	//	8 = TknCodeEnd
	//	9 = unused
	//	>9 = TknString, value - 10 is the StringTable index
	Code        []uint32
	Positions   []*PositionInfo
	StringTable []string
}

// Compile will generate a new CompiledScript from a string.
func Compile(input string, pos *PositionInfo) *CompiledScript {
	this := new(CompiledScript)
	this.Code = make([]uint32, 0, 20)
	this.Positions = make([]*PositionInfo, 0, 20)

	lex := NewLexer(input, pos.Line, pos.Column)

	nextstring := 0
	strings := make(map[string]int)

	for {
		lex.Advance()

		if lex.CurrentTkn().Type == TknINVALID {
			break
		}

		if lex.CurrentTkn().Type == TknString {
			if _, ok := strings[lex.CurrentTkn().Lexeme]; ok {
				this.Code = append(this.Code, uint32(strings[lex.CurrentTkn().Lexeme]+10))
			} else {
				strings[lex.CurrentTkn().Lexeme] = nextstring
				this.Code = append(this.Code, uint32(strings[lex.CurrentTkn().Lexeme]+10))
				nextstring++
			}
		} else {
			this.Code = append(this.Code, uint32(lex.CurrentTkn().Type))
		}

		this.Positions = append(this.Positions, lex.CurrentTkn().Pos)
	}

	this.StringTable = make([]string, nextstring)
	for str, i := range strings {
		this.StringTable[i] = str
	}
	return this
}

// String should convert a CompiledScript back to a string containing valid Raptor code.
func (this *CompiledScript) String() string {
	out := "{ "
	for _, code := range this.Code {
		if code < 9 {
			out += tknLexemes[code] + " "
			continue
		}
		if code > 9 {
			out += "\"" + EscapeString(this.StringTable[code-10]) + "\" "
			continue
		}
	}
	return out + "}"
}

// CompileBlock will generate a new CompiledScript from an existing CodeSource, 
// reading from the current token (which MUST be a TknCodeBegin) to the matching TknCodeEnd.
// (Nested blocks are allowed)
func CompileBlock(lex CodeSource) *CompiledScript {
	if lex.CurrentTkn().Type != TknCodeBegin {
		panic("CompileBlock called on a CodeSource that is not entering a block.")
	}

	// Compiling from a CompiledLexer
	if tmp, ok := lex.(*CompiledLexer); ok {
		code := tmp.data
		start := tmp.currentToken + 1

		blockDepth := 1
		i := start
		foundend := false
		for ; i < len(code.Code); i++ {
			if code.Code[i] == TknCodeBegin {
				blockDepth++
			}
			if code.Code[i] == TknCodeEnd {
				blockDepth--
				if blockDepth == 0 {
					foundend = true
					break
				}
			}
		}
		if !foundend {
			panic("CompileBlock: Unexpected EOS.")
		}

		this := new(CompiledScript)
		this.Code = code.Code[start:i]

		if code.Positions != nil {
			this.Positions = code.Positions[start:i]
		}

		this.StringTable = code.StringTable

		tmp.currentToken = i

		return this
	}

	// Compiling from some other CodeSource
	this := new(CompiledScript)
	this.Code = make([]uint32, 0, 20)
	this.Positions = make([]*PositionInfo, 0, 20)

	blockDepth := 1

	nextstring := 0
	strings := make(map[string]int)

	for {
		lex.Advance()

		if lex.CurrentTkn().Type == TknINVALID {
			panic("CompileBlock: Unexpected TknINVALID.")
		}

		// These MUST fall through! (Unless the end of the block is found)
		if lex.CurrentTkn().Type == TknCodeBegin {
			blockDepth++
		}
		if lex.CurrentTkn().Type == TknCodeEnd {
			blockDepth--
			if blockDepth == 0 {
				break
			}
		}

		if lex.CurrentTkn().Type == TknString {
			if _, ok := strings[lex.CurrentTkn().Lexeme]; ok {
				this.Code = append(this.Code, uint32(strings[lex.CurrentTkn().Lexeme]+10))
			} else {
				strings[lex.CurrentTkn().Lexeme] = nextstring
				this.Code = append(this.Code, uint32(strings[lex.CurrentTkn().Lexeme]+10))
				nextstring++
			}
		} else {
			this.Code = append(this.Code, uint32(lex.CurrentTkn().Type))
		}

		this.Positions = append(this.Positions, lex.CurrentTkn().Pos)
	}

	this.StringTable = make([]string, nextstring)
	for str, i := range strings {
		this.StringTable[i] = str
	}
	return this
}

// CompiledLexer is a CodeSource wrapper for CompiledScript.
type CompiledLexer struct {
	currentToken int
	line         int
	data         *CompiledScript
}

// NewCompiledLexer creates a new CodeSource reading from a compiled script.
func NewCompiledLexer(val *CompiledScript) CodeSource {
	this := new(CompiledLexer)
	this.data = val
	this.currentToken = -1
	return this
}

// Position returns a PositionInfo for the last token fetched.
func (this *CompiledLexer) Position() *PositionInfo {
	if !this.ValidatePosition(this.currentToken) || this.data.Positions == nil {
		return NewPositionInfo(this.currentToken, -1)
	}
	return this.data.Positions[this.currentToken]
}

// CurrentTkn returns the current token.
func (this *CompiledLexer) CurrentTkn() *Token {
	return this.Token(this.currentToken)
}

// LookAhead returns the lookahead token.
func (this *CompiledLexer) LookAhead() *Token {
	return this.Token(this.currentToken + 1)
}

// Advance retrieves the next token from the stream.
// For most purposes use GetToken instead.
func (this *CompiledLexer) Advance() {
	this.currentToken++
	if this.currentToken > len(this.data.Code) {
		this.currentToken = len(this.data.Code)
	}
}

// GetToken gets the next token, and panics with an error if it's not of type tokenType.
// Used as a type checked Advance.
func (this *CompiledLexer) GetToken(tokenTypes ...int) {
	this.Advance()
	current := this.CurrentTkn()

	for _, val := range tokenTypes {
		if current.Type == val {
			return
		}
	}

	ExitOnTokenExpected(current, tokenTypes...)
}

// CheckLookAhead checks to see if the look ahead is one of tokenTypes and if so returns true.
func (this *CompiledLexer) CheckLookAhead(tokenTypes ...int) bool {
	look := this.LookAhead()
	for _, val := range tokenTypes {
		if look.Type == val {
			return true
		}
	}
	return false
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
func (this *CompiledLexer) Token(index int) *Token {
	if index >= len(this.data.Code) || index < 0 {
		return &Token{"INVALID", TknINVALID, NewPositionInfo(0, -1)}
	}
	code := this.data.Code[index]

	if code < 9 {
		if this.data.Positions == nil {
			return &Token{tknLexemes[code], int(code), NewPositionInfo(index, -1)}
		}

		return &Token{tknLexemes[code], int(code), this.data.Positions[index]}
	}

	if code == 9 {
		// Unused token code
		return &Token{"INVALID", TknINVALID, NewPositionInfo(0, -1)}
	}

	code = code - 10
	if int(code) < len(this.data.StringTable) && code >= 0 {
		if this.data.Positions == nil {
			return &Token{this.data.StringTable[code], TknString, NewPositionInfo(index, -1)}
		}

		return &Token{this.data.StringTable[code], TknString, this.data.Positions[index]}
	}

	return &Token{"INVALID", TknINVALID, NewPositionInfo(0, -1)}
}

// ValidatePosition will return true if index is a valid token code.
func (this *CompiledLexer) ValidatePosition(index int) bool {
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
