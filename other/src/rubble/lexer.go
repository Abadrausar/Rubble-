/*
Copyright 2013 by Milo Christiansen

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

package main

// A semi-generic Lexer framework
type Lexer struct {

	// The next/current tokens
	Look    *Token
	Current *Token

	source []byte

	line      int
	tokenline int

	index int

	depth int

	commandDepth int

	lexeme []byte

	token int
}

// Returns a new Lexer.
func NewLexer(dat []byte) *Lexer {
	this := new(Lexer)

	this.Look = &Token{"INVALID", tknINVALID, -1}
	this.Current = &Token{"INVALID", tknINVALID, -1}

	this.source = dat

	this.line = 1
	this.tokenline = 1

	this.index = 0

	this.depth = 0

	this.commandDepth = 0

	this.lexeme = make([]byte, 0, 20)

	this.token = tknString

	this.Advance()

	return this
}

// This advances the Lexer one token.
// For most purposes use GetToken instead.
func (this *Lexer) Advance() {
	if this.index > len(this.source) {
		this.Current = this.Look
		LastLine = this.Current.Line
		this.Look = &Token{"INVALID", tknINVALID, this.tokenline}
		return
	}

	for ; this.index < len(this.source); this.index++ {
		dat := this.source
		i := this.index
		lookok := len(this.source) - this.index

		if dat[i] == '\n' {
			this.line++
		}

		if this.depth < 0 {
			panic("Lexer template depth less than 0 (Unmatched curly brackets)")
		}

		if dat[i] == ';' || dat[i] == '{' || dat[i] == '}' {
			if 0 < i && lookok > 0 {
				if dat[i-1] == '\'' && dat[i+1] == '\'' {
					this.lexeme = append(this.lexeme, dat[i])
					continue
				}
			}

			if len(this.lexeme) > 0 && this.depth == 0 {
				this.Current = this.Look
				LastLine = this.Current.Line
				this.Look = &Token{string(this.lexeme), tknString, this.tokenline}

				this.tokenline = -1
				this.lexeme = this.lexeme[0:0]
				return
			}
		}

		if dat[i] == ';' || dat[i] == '}' {
			if len(this.lexeme) > 0 && this.depth == 1 {
				this.Current = this.Look
				LastLine = this.Current.Line
				this.Look = &Token{string(this.lexeme), tknString, this.tokenline}

				this.tokenline = -1
				this.lexeme = this.lexeme[0:0]
				return
			}
		}

		if dat[i] == ';' {
			if this.depth != 1 {
				this.lexeme = append(this.lexeme, dat[i])
				continue
			}

			this.Current = this.Look
			LastLine = this.Current.Line
			this.Look = &Token{";", tknDelimiter, this.line}
			this.tokenline = this.line
			this.lexeme = this.lexeme[0:0]
			this.index++
			return
		}

		if dat[i] == '{' {
			this.depth++
			if this.depth > 1 {
				this.lexeme = append(this.lexeme, dat[i])
				continue
			}

			this.Current = this.Look
			LastLine = this.Current.Line
			this.Look = &Token{"{", tknTagBegin, this.line}
			this.tokenline = this.line
			this.lexeme = this.lexeme[0:0]
			this.index++
			return
		}

		if dat[i] == '}' {
			this.depth--
			if this.depth != 0 {
				this.lexeme = append(this.lexeme, dat[i])
				continue
			}

			this.Current = this.Look
			LastLine = this.Current.Line
			this.Look = &Token{"}", tknTagEnd, this.line}
			this.tokenline = this.line
			this.lexeme = this.lexeme[0:0]
			this.index++
			return
		}

		this.lexeme = append(this.lexeme, dat[i])
	}

	if this.index == len(this.source) {
		this.Current = this.Look
		LastLine = this.Current.Line
		this.Look = &Token{string(this.lexeme), tknString, this.tokenline}
		this.index++
		return
	}
}

// Gets the next token, and panics with an error if it's not of type tokenType.
// Used as a type checked Advance
func (this *Lexer) GetToken(tokenTypes ...int) {
	this.Advance()

	for _, val := range tokenTypes {
		if this.Current.Type == val {
			return
		}
	}

	ExitOnTokenExpected(this.Current, tokenTypes...)
}

// Checks to see if the look ahead is one of tokenTypes and if so returns true
func (this *Lexer) CheckLookAhead(tokenTypes ...int) bool {
	for _, val := range tokenTypes {
		if this.Look.Type == val {
			return true
		}
	}
	return false
}
