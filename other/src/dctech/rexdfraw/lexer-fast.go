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

package rexdfraw

//import "fmt"
import "strconv"
import "dctech/rex"

// Lexer states
const (
	stEatComment = iota
	stReadTagID
	stReadParam
)

// rawlexer is a Dwarf Fortress raw file lexer reading from a string.
// The lexer interfaces with Rex by acting like an Indexable.
type rawlexer struct {
	source []byte

	line   int
	column int

	index int
	state int

	tagbegin int
	tagend   int
	tagvalid bool

	lexeme []byte

	curID     string
	curParams []string
}

// Returns a new Lexer.
func newRawLexer(input string) *rawlexer {
	lex := new(rawlexer)

	lex.source = []byte(input)

	lex.line = 1
	lex.column = 1

	lex.index = 0
	lex.state = stEatComment

	lex.tagbegin = 0
	lex.tagend = 0
	lex.tagvalid = false

	lex.lexeme = make([]byte, 0, 20)

	return lex
}

func (lex *rawlexer) advance() {

	if lex.index >= len(lex.source) {
		// EOS
		lex.tagvalid = false
		return
	}

	for ; lex.index < len(lex.source); lex.index++ {
		//lookok := len(lex.source) - lex.index

		if lex.source[lex.index] == '\n' && lex.line >= 0 {
			lex.line++
			lex.column = 0
		} else {
			lex.column++
		}

		// Lexing Begin
		//======================================

		if lex.source[lex.index] == '[' {
			if lex.state != stEatComment {
				// Just get back on track...
				//rex.RaiseError(fmt.Sprint("Nested Raw Tag found on line: ", lex.line))
				lex.lexeme = lex.lexeme[0:0]
			}

			lex.tagbegin = lex.index
			lex.tagend = lex.index
			lex.curID = ""
			lex.curParams = make([]string, 0)
			lex.tagvalid = false
			lex.state = stReadTagID
			continue
		}

		if lex.source[lex.index] == ']' {
			if lex.state == stReadTagID {
				lex.tagend = lex.index
				lex.curID = string(lex.lexeme)
				lex.lexeme = lex.lexeme[0:0]
				lex.tagvalid = true
				lex.state = stEatComment
				return
			} else if lex.state == stReadParam {
				lex.tagend = lex.index
				lex.curParams = append(lex.curParams, string(lex.lexeme))
				lex.lexeme = lex.lexeme[0:0]
				lex.tagvalid = true
				lex.state = stEatComment
				return
			}
			continue
		}

		if lex.source[lex.index] == ':' {
			if lex.state == stReadTagID {
				lex.curID = string(lex.lexeme)
				lex.lexeme = lex.lexeme[0:0]
				lex.state = stReadParam
			} else if lex.state == stReadParam {
				lex.curParams = append(lex.curParams, string(lex.lexeme))
				lex.lexeme = lex.lexeme[0:0]
			}
			continue
		}

		if lex.state == stEatComment {
			continue
		}
		lex.lexeme = append(lex.lexeme, lex.source[lex.index])
	}
	lex.tagvalid = false
}

func (lex *rawlexer) Get(index string) *rex.Value {
	if index == "id" {
		return rex.NewValueString(lex.curID)
	}

	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		rex.RaiseError("Index not a valid number.")
	}
	if val < 0 || val >= int64(len(lex.curParams)) {
		rex.RaiseError("Index out of range.")
	}
	return rex.NewValueString(lex.curParams[val])
}

func (lex *rawlexer) Set(index string, value *rex.Value) bool {
	if index == "replace" {
		newstring := value.String()
		beginstr := lex.source[:lex.tagbegin]
		endstr := lex.source[lex.tagend+1:]
	
		newlen := len(newstring) + len(beginstr) + len(endstr)
	
		rtn := make([]byte, 0, newlen)
		rtn = append(rtn, beginstr...)
		rtn = append(rtn, newstring...)
		lex.index = len(rtn)
		rtn = append(rtn, endstr...)
		lex.source = rtn
		return true
	}

	if index == "append" {
		newstring := value.String()
		beginstr := lex.source[:lex.tagend+1]
		endstr := lex.source[lex.tagend+1:]
	
		newlen := len(newstring) + len(beginstr) + len(endstr)
	
		rtn := make([]byte, 0, newlen)
		rtn = append(rtn, beginstr...)
		rtn = append(rtn, newstring...)
		lex.index = len(rtn)
		rtn = append(rtn, endstr...)
		lex.source = rtn
		return true
	}

	if index == "prepend" {
		newstring := value.String()
		beginstr := lex.source[:lex.tagbegin]
		endstr := lex.source[lex.tagbegin:]
	
		newlen := len(newstring) + len(beginstr) + len(endstr)
	
		rtn := make([]byte, 0, newlen)
		rtn = append(rtn, beginstr...)
		rtn = append(rtn, newstring...)
		lex.index = len(rtn)
		rtn = append(rtn, endstr...)
		lex.source = rtn
		return true
	}

	if index == "disable" && value.Bool() {
		if value.Bool() {
			lex.source[lex.tagbegin] = '-'
			lex.source[lex.tagend] = '-'
		} else {
			lex.source[lex.tagbegin] = '['
			lex.source[lex.tagend] = ']'
		}
		return true
	}
	
	return false
}

func (lex *rawlexer) Exists(index string) bool {
	if index == "id" {
		return true
	}

	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		return false
	}
	if val < 0 || val >= int64(len(lex.curParams)) {
		return false
	}
	return true
}

func (lex *rawlexer) Len() int64 {
	return int64(len(lex.curParams))
}

func (lex *rawlexer) Keys() []string {
	rtn := make([]string, 0, len(lex.curParams))
	for key := range lex.curParams {
		rtn = append(rtn, strconv.FormatInt(int64(key), 10))
	}
	return rtn
}

func (lex *rawlexer) String() string {
	return rex.IndexableToString("raw:Lexer", lex)
}

func (lex *rawlexer) CodeString() string {
	return rex.IndexableToCodeString("map", lex, true)
}
