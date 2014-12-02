/*
Copyright 2013-2014 by Milo Christiansen

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

package rubble

import "dctech/rex"

type lexer struct {
	Look    *token
	Current *token

	source []byte

	pos      *rex.Position
	tokenpos *rex.Position

	index int

	depth int

	commandDepth int

	lexeme []byte

	token int
}

func newLexer(dat []byte, pos *rex.Position) *lexer {
	lex := new(lexer)

	lex.Look = &token{"INVALID", tknINVALID, pos}
	lex.Current = &token{"INVALID", tknINVALID, pos}

	lex.source = dat

	lex.pos = pos.Copy()
	lex.tokenpos = pos.Copy()

	lex.index = 0

	lex.depth = 0

	lex.commandDepth = 0

	lex.lexeme = make([]byte, 0, 20)

	lex.token = tknString

	lex.Advance()

	return lex
}

// This advances the Lexer one token.
// For most purposes use GetToken instead.
func (lex *lexer) Advance() {
	if lex.index > len(lex.source) {
		lex.Current = lex.Look
		lex.Look = &token{"INVALID", tknINVALID, lex.tokenpos.Copy()}
		return
	}

	for ; lex.index < len(lex.source); lex.index++ {
		dat := lex.source
		i := lex.index
		lookok := len(lex.source) - lex.index

		if dat[i] == '\n' {
			lex.pos.Line++
			lex.pos.Column = 0
		} else {
			lex.pos.Column++
		}

		if dat[i] == ';' || dat[i] == '{' || dat[i] == '}' {
			if 0 < i && lookok > 0 {
				if dat[i-1] == '\'' && dat[i+1] == '\'' {
					lex.lexeme = append(lex.lexeme, dat[i])
					continue
				}
			}

			if len(lex.lexeme) > 0 && lex.depth == 0 {
				lex.Current = lex.Look
				lex.Look = &token{string(lex.lexeme), tknString, lex.tokenpos.Copy()}

				lex.tokenpos = lex.pos.Copy()
				lex.lexeme = lex.lexeme[0:0]
				return
			}
		}

		if dat[i] == ';' || dat[i] == '}' {
			if len(lex.lexeme) > 0 && lex.depth == 1 {
				lex.Current = lex.Look
				lex.Look = &token{string(lex.lexeme), tknString, lex.tokenpos.Copy()}

				lex.tokenpos = lex.pos.Copy()
				lex.lexeme = lex.lexeme[0:0]
				return
			}
		}

		if dat[i] == ';' {
			if lex.depth != 1 {
				lex.lexeme = append(lex.lexeme, dat[i])
				continue
			}

			lex.Current = lex.Look
			lex.Look = &token{";", tknDelimiter, lex.tokenpos.Copy()}
			lex.tokenpos = lex.pos.Copy()
			lex.lexeme = lex.lexeme[0:0]
			lex.index++
			return
		}

		if dat[i] == '{' {
			lex.depth++
			if lex.depth > 1 {
				lex.lexeme = append(lex.lexeme, dat[i])
				continue
			}

			lex.Current = lex.Look
			lex.Look = &token{"{", tknTagBegin, lex.tokenpos.Copy()}
			lex.tokenpos = lex.pos.Copy()
			lex.lexeme = lex.lexeme[0:0]
			lex.index++
			return
		}

		if dat[i] == '}' {
			lex.depth--
			if lex.depth != 0 {
				if lex.depth < 0 {
					lex.depth = 0
				}
				lex.lexeme = append(lex.lexeme, dat[i])
				continue
			}

			lex.Current = lex.Look
			lex.Look = &token{"}", tknTagEnd, lex.tokenpos.Copy()}
			lex.tokenpos = lex.pos.Copy()
			lex.lexeme = lex.lexeme[0:0]
			lex.index++
			return
		}

		lex.lexeme = append(lex.lexeme, dat[i])
	}

	if lex.index == len(lex.source) {
		lex.Current = lex.Look
		lex.Look = &token{string(lex.lexeme), tknString, lex.tokenpos.Copy()}
		lex.index++
		return
	}
}

// Gets the next token, and panics with an error if it's not of type tokenType.
// Used as a type checked Advance
func (lex *lexer) GetToken(tokenTypes ...int) {
	lex.Advance()

	for _, val := range tokenTypes {
		if lex.Current.Type == val {
			return
		}
	}

	exitOnTokenExpected(lex.Current, tokenTypes...)
}

// Checks to see if the look ahead is one of tokenTypes and if so returns true
func (lex *lexer) CheckLookAhead(tokenTypes ...int) bool {
	for _, val := range tokenTypes {
		if lex.Look.Type == val {
			return true
		}
	}
	return false
}
