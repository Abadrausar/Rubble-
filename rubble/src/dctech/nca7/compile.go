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

package nca7

//import "fmt"

// WARNING WARNING WARNING
// I, in a stupid moment, used "magic numbers" all over in here
// any changes will probibly break all sorts of things.

// CompiledScript is a compiled NCA script.
type CompiledScript struct {
	// Code Values:
	//	0 = TknCmdBegin
	//	1 = TknCmdEnd
	//	2 = TknDerefBegin
	//	3 = TknDerefEnd
	//	4 = TknObjLitBegin
	//	5 = TknObjLitEnd
	//	6 = TknObjLitSplit
	//	7 = unused
	//	8 = unused
	//	9 = Newline, this is never seen by the parser, it is for error reporting only.
	//	>9 = TknString, value - 10 is the StringTable index
	Code []uint32
	StringTable []string
	Start int
}

// Compile will generate a new CompiledScript from a string.
func Compile(input string, line int) *CompiledScript {
	this := new(CompiledScript)
	this.Code = make([]uint32, 0)
	this.Start = line
	
	lex := NewLexer(input, line, 0)
	
	nextstring := 0
	strings := make(map[string]int)
	lastLine := line
	
	for {
		lex.Advance()
		
		if lex.CurrentTkn().Type == TknINVALID {
			break
		}
		
		for lex.CurrentTkn().Line > lastLine {
			this.Code = append(this.Code, 9)
			lastLine++
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
			this.Code = append(this.Code, uint32(lex.CurrentTkn().Type-2))
		}
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
	line int
	data *CompiledScript
}

// NewCompiledLexer creates a new CodeSource reading from a compiled script.
func NewCompiledLexer(val *CompiledScript) CodeSource {
	this := new(CompiledLexer)
	this.data = val
	this.line = val.Start
	return this
}

// Line returns the line of the last token fetched.
func (this *CompiledLexer) Line() int {
	return this.line
}

// Column always returns 0.
func (this *CompiledLexer) Column() int {
	return 0
}

// CurrentTkn returns the current token.
func (this *CompiledLexer) CurrentTkn() *Token {
	return this.Token(this.currentToken)
}

// LookAhead returns the lookahead token.
func (this *CompiledLexer) LookAhead() *Token {
	token := this.currentToken
	for {
		token++
		if token >= len(this.data.Code) || token < 0 {
			return this.Token(token)
		}
		code := this.data.Code[token]
		if code == 9 {
			continue
		}
		return this.Token(token)
	}
	panic("UNREACHABLE: *CompiledLexer.LookAhead")
}

// This advances the code one token.
// For most purposes use GetToken instead.
func (this *CompiledLexer) Advance() {
	for {
		this.currentToken++
		if this.currentToken >= len(this.data.Code) || this.currentToken < 0 {
			return
		}
		code := this.data.Code[this.currentToken]
		if code == 9 {
			this.line++
			continue
		}
		return
	}
}

// Gets the next token, and panics with an error if it's not of type tokenType.
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

// Checks to see if the look ahead is one of tokenTypes and if so returns true.
func (this *CompiledLexer) CheckLookAhead(tokenTypes ...int) bool {
	look := this.LookAhead()
	for _, val := range tokenTypes {
		if look.Type == val {
			return true
		}
	}
	return false
}

// Token creates a Token from the data at index.
// Returns TknINVALID if the index is invalid.
func (this *CompiledLexer) Token(index int) *Token {
	if index >= len(this.data.Code) || index < 0 {
		return &Token { "INVALID", TknINVALID, this.line, 0}
	}
	code := this.data.Code[index]
	
	if code == 0 {
		return &Token { "(", TknCmdBegin, this.line, 0}
	}
	
	if code == 1 {
		return &Token { ")", TknCmdEnd, this.line, 0}
	}
	
	if code == 2 {
		return &Token { "[", TknDerefBegin, this.line, 0}
	}
	
	if code == 3 {
		return &Token { "]", TknDerefEnd, this.line, 0}
	}
	
	if code == 4 {
		return &Token { "<", TknObjLitBegin, this.line, 0}
	}
	
	if code == 5 {
		return &Token { ">", TknObjLitEnd, this.line, 0}
	}
	
	if code == 6 {
		return &Token { "=", TknObjLitSplit, this.line, 0}
	}
	
	if code == 7 {
		panic("Invalid token code: 7 unused.")
	}
	
	if code == 8 {
		panic("Invalid token code: 8 unused.")
	}
	
	if code == 9 { // This should NEVER, EVER happen
		panic("Invalid token code: 9 control code.")
	}
	
	code = code - 10
	if int(code) < len(this.data.StringTable) && code >= 0 {
		return &Token { this.data.StringTable[code], TknString, this.line, 0}
	}
	
	return &Token { "INVALID", TknINVALID, this.line, 0}
}
