/*
Copyright 2014 by Milo Christiansen

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

package expr

import "strings"
import "dctech/rex"

const (
	tknINVALID = iota
	
	tknAdd
	tknSub
	tknMul
	tknDiv
	tknMod
	
	tknAnd
	tknOr
	tknNot
	
	tknEq
	tknNotEq
	tknGt
	tknGtEq
	tknLt
	tknLtEq
	
	tknVal
	
	tknOParen
	tknCParen
)

type lexer struct {
	look    int
	current int

	source *strings.Reader
	char   byte
	eof    bool // true if there are no more chars to read
}

func newLexer(input string) *lexer {
	lex := new(lexer)

	lex.source = strings.NewReader(input)

	// prime the pump
	lex.nextchar()
	lex.look = tknINVALID
	lex.advance()
	
	return lex
}

// advance retrieves the next token from the stream.
// For most purposes use Gettoken instead.
func (lex *lexer) advance() {
	lex.current = lex.look
	lex.look = tknINVALID
	
	if lex.eof {
		return
	}
	
	lex.eatWS()
	if lex.eof {
		return
	}
	
	// We are at the beginning of a token
	switch lex.char {
	case '+':
		lex.look = tknAdd
		lex.nextchar()
	
	case '-':
		lex.look = tknSub
		lex.nextchar()
	
	case '*':
		lex.look = tknMul
		lex.nextchar()
	
	case '/':
		lex.look = tknDiv
		lex.nextchar()
	
	case '%':
		lex.look = tknMod
		lex.nextchar()
	
	case '&':
		lex.getChar("&")
		lex.look = tknAnd
		lex.nextchar()
	
	case '|':
		lex.getChar("|")
		lex.look = tknOr
		lex.nextchar()
	
	case '!':
		lex.nextchar()
		if lex.match("=") {
			lex.look = tknNotEq
			lex.nextchar()
			break
		}
		lex.look = tknNot
	
	case '=':
		lex.getChar("=")
		lex.look = tknEq
		lex.nextchar()
	
	case '>':
		lex.nextchar()
		if lex.match("=") {
			lex.look = tknGtEq
			lex.nextchar()
			break
		}
		lex.look = tknGt
	
	case '<':
		lex.nextchar()
		if lex.match("=") {
			lex.look = tknLtEq
			lex.nextchar()
			break
		}
		lex.look = tknLt
	
	case '(':
		lex.look = tknOParen
		lex.nextchar()
	
	case ')':
		lex.look = tknCParen
		lex.nextchar()
	
	default:
		lex.look = tknVal
		lex.nextchar()
	}
}

// checkLook checks to see if the look ahead is one of tokenTypes and if so returns true.
func (lex *lexer) checkLook(tokenTypes ...int) bool {
	for _, val := range tokenTypes {
		if lex.look == val {
			return true
		}
	}
	return false
}

func (lex *lexer) getToken(tokenTypes ...int) {
	if !lex.checkLook(tokenTypes...) {
		rex.ErrorGeneralCmd("expr", "Expected token not found.")
	}
	
	lex.advance()
}

// Fetch the next char
func (lex *lexer) nextchar() {
	if lex.eof {return}
	
	// err should only ever be io.EOF
	var err error
	lex.char, err = lex.source.ReadByte()
	if err != nil {
		lex.eof = true
		return
	}
}

func (lex *lexer) match(chars string) bool {
	for _, char := range []byte(chars) {
		if lex.char == char {
			return true
		}
	}
	return false
}
func (lex *lexer) getChar(chars string) {
	lex.nextchar()
	if lex.eof {
		rex.ErrorGeneralCmd("expr", "Unexpected EOS.")
	}
	
	if lex.match(chars) {
		return
	}
	rex.ErrorGeneralCmd("expr", "Expected char not found.")
}

// Eat white space and comments.
func (lex *lexer) eatWS() {
	for {
		if lex.match("\n\r \t") {
			lex.nextchar()
			if lex.eof {return}
			continue
		}
		break
	}
}
