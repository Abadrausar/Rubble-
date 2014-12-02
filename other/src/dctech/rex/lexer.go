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

package rex

import "strconv"
import "strings"

// The lexer is a simple state machine, but rather than implementing one big generic state machine
// with lots of states, I have a little semi-independent one for each type of value that is called
// from a central "dispatcher" switch statement.
// State is entirely implicit.

// All the delimiter chars, these are the characters that are illegal in a raw string.
// Including a character in this string that is not recognized and consumed by the lexer
// will cause an OOM condition!
var delimiters = "\n\r \t#(){}[]<>=:\"'`"

// Lexer is for feeding the compiler only.
type Lexer struct {
	look    *token
	current *token

	source *strings.Reader
	char   rune
	eof    bool // true if there are no more chars to read

	pos *Position

	lexeme []rune

	token    int
	tokenpos *Position

	strdepth int
	objdepth int
}

// Returns a new Lexer.
// About the only time you will need this is when you use CompileExisting.
func NewLexer(input string, pos *Position) *Lexer {
	lex := new(Lexer)

	lex.source = strings.NewReader(input)

	lex.pos = pos.Copy()

	lex.lexeme = make([]rune, 0, 20)

	lex.token = tknINVALID
	lex.tokenpos = pos.Copy()

	lex.strdepth = 0
	lex.objdepth = 0
	
	// prime the pump
	lex.nextchar()
	lex.look = &token{"INVALID", tknINVALID, lex.tokenpos}
	lex.advance()
	
	return lex
}

// advance retrieves the next token from the stream.
// For most purposes use getcurrent instead.
func (lex *Lexer) advance() {
	lex.current = lex.look
	if lex.eof {
		lex.look = &token{"EOF", tknINVALID, lex.tokenpos}
		return
	}
	
	lex.eatWS()
	if lex.eof {
		lex.look = &token{"EOF", tknINVALID, lex.tokenpos}
		return
	}
	
	// We are at the beginning of a token
	switch lex.char {
	case '(':
		lex.look = &token{"", tknCmdBegin, lex.pos}
		lex.nextchar()
	case ')':
		lex.look = &token{"", tknCmdEnd, lex.pos}
		lex.nextchar()
	case '[':
		lex.look = &token{"", tknVarBegin, lex.pos}
		lex.nextchar()
	case ']':
		lex.look = &token{"", tknVarEnd, lex.pos}
		lex.nextchar()
	case '<':
		lex.look = &token{"", tknObjLitBegin, lex.pos}
		lex.nextchar()
	case '>':
		lex.look = &token{"", tknObjLitEnd, lex.pos}
		lex.nextchar()
	case '{':
		lex.look = &token{"", tknCodeBegin, lex.pos}
		lex.nextchar()
	case '}':
		lex.look = &token{"", tknCodeEnd, lex.pos}
		lex.nextchar()
	case ':':
		lex.look = &token{"", tknNameSplit, lex.pos}
		lex.nextchar()
	case '=':
		lex.look = &token{"", tknAssignment, lex.pos}
		lex.nextchar()
	case 'n':
		lex.matchKeyword("nil", tknNil)
	case 't':
		lex.matchKeyword("true", tknTrue)
	case 'f':
		lex.matchKeyword("false", tknFalse)
	case 'm':
		lex.matchKeyword("module", tknDeclModule)
	case 'c':
		lex.matchKeyword("command", tknDeclCommand)
	case 'b':
		lex.matchKeyword("block", tknDeclBlock)
	case 'v':
		lex.matchKeyword("var", tknDeclVar)
	case '"':
		lex.matchString('"')
	case '\'':
		lex.matchString('\'')
	case '`':
		lex.matchString('`')
	default:
		lex.matchRawString()
	}
	
	lex.lexeme = lex.lexeme[0:0]
}

// getcurrent gets the next token, and panics with an error if it's not of type tokenType.
// May cause a panic if the lexer encounters an error.
// Used as a type checked advance.
func (lex *Lexer) getcurrent(tokenTypes ...int) {
	lex.advance()

	for _, val := range tokenTypes {
		if lex.current.Type == val {
			return
		}
	}

	exitOntokenExpected(lex.current, tokenTypes...)
}

// checkLook checks to see if the look ahead is one of tokenTypes and if so returns true.
func (lex *Lexer) checkLook(tokenTypes ...int) bool {
	for _, val := range tokenTypes {
		if lex.look.Type == val {
			return true
		}
	}
	return false
}

// return true if the current char matches one of the chars in the string.
func (lex *Lexer) match(chars string) bool {
	for _, char := range chars {
		if lex.char == char {
			return true
		}
	}
	return false
}

// Fetch the next char (actually a Unicode code point).
// I don't like the way EOF is handled, but there is really no better way that is flexible enough.
func (lex *Lexer) nextchar() {
	if lex.eof {return}
	
	// err should only ever be io.EOF
	var err error
	lex.char, _, err = lex.source.ReadRune()
	if err != nil {
		lex.eof = true
		return
	}
	
	if lex.char == '\n' {
		lex.pos.Line++
		lex.pos.Column = 0
	} else {
		lex.pos.Column++
	}
}

// Add the current char to the lexeme buffer.
func (lex *Lexer) addLexeme() {
	lex.lexeme = append(lex.lexeme, lex.char)
}

// Eat white space and comments.
func (lex *Lexer) eatWS() {
	for {
		if lex.match("#") {
			lex.nextchar()
			if lex.eof {return}
			for {
				if lex.match("\n") {
					lex.nextchar()
					if lex.eof {return}
					break
				}
				lex.nextchar()
				if lex.eof {return}
			}
			continue
		}
		if lex.match("\n\r \t") {
			lex.nextchar()
			if lex.eof {return}
			continue
		}
		break
	}
}

// Try to match a key word, fails to matching a raw string.
func (lex *Lexer) matchKeyword(str string, typ int) {
	pos := lex.pos.Copy()
	
	for _, char := range str {
		if lex.char == char {
			lex.addLexeme()
			lex.nextchar()
			if lex.eof {
				lex.look = &token{string(lex.lexeme), tknRawString, pos}
				return
			}
		} else {
			lex.matchRawString()
			return
		}
	}
	
	if lex.match(delimiters) {
		lex.look = &token{"", typ, pos}
		return
	} else {
		lex.matchRawString()
		return
	}
}

// Match a delimited string, reads to end delimiter or EOF
func (lex *Lexer) matchString(delim rune) {
	pos := lex.pos.Copy()
	
	lex.nextchar()
	if lex.eof || lex.char == delim {
		lex.look = &token{"", tknString, pos}
		lex.nextchar()
		return
	}
	
	for lex.char != delim {
		// Handle escapes
		if lex.char == '\\' {
			lex.nextchar()
			if lex.eof {
				lex.lexeme = append(lex.lexeme, '\\')
				break
			}
			
			switch lex.char {
			case 'n':
				lex.lexeme = append(lex.lexeme, '\n')
			case 'r':
				lex.lexeme = append(lex.lexeme, '\r')
			case 't':
				lex.lexeme = append(lex.lexeme, '\t')
			case '"':
				lex.lexeme = append(lex.lexeme, '"')
			case '\'':
				lex.lexeme = append(lex.lexeme, '\'')
			case '`':
				lex.lexeme = append(lex.lexeme, '`')
			case '\\':
				lex.lexeme = append(lex.lexeme, '\\')
			case 'x':
				lex.nextchar()
				if lex.eof {
					lex.lexeme = append(lex.lexeme, '\\', 'x')
					break
				}
				a := lex.char
				lex.nextchar()
				if lex.eof {
					lex.lexeme = append(lex.lexeme, '\\', 'x', a)
					break
				}
				b := lex.char
				
				rep, err := strconv.ParseInt(string([]rune{a, b}), 16, 8)
				if err != nil {
					lex.lexeme = append(lex.lexeme, '\\', 'x', a, b)
				}
				lex.lexeme = append(lex.lexeme, rune(rep))

			default:
				lex.lexeme = append(lex.lexeme, '\\', lex.char)
			}
			
			lex.nextchar()
			if lex.eof {
				break
			}
			continue
		}
		
		lex.addLexeme()
		lex.nextchar()
		if lex.eof {
			break
		}
	}
	
	if lex.char == delim {
		lex.nextchar()
	}
	lex.look = &token{string(lex.lexeme), tknString, pos}
}

// Match a raw string.
func (lex *Lexer) matchRawString() {
	pos := lex.pos.Copy()
	
	for !lex.match(delimiters) {
		lex.addLexeme()
		lex.nextchar()
		if lex.eof {
			break
		}
	}
	lex.look = &token{string(lex.lexeme), tknRawString, pos}
}
