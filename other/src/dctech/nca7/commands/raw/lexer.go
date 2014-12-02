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

package raw

import "fmt"
import "strconv"
import "dctech/nca7"

// Lexer states
const (
	stEatComment = iota
	stReadTagID
	stReadParam
)

var lex *rawlexer

// Lexer is a NCA lexer reading from a string.
type rawlexer struct {
	source []byte
	
	line int
	column int
	
	index int
	state int
	
	tagbegin int
	tagend int
	tagvalid bool
	
	lexeme []byte
	
	curID     string
	curParams []string
}

// Returns a new Lexer.
func newRawLexer(input string) *rawlexer {
	this := new(rawlexer)
	
	this.source = []byte(input)
	
	this.line = 1
	this.column = 1
	
	this.index = 0
	this.state = stEatComment
	
	this.tagbegin = 0
	this.tagend = 0
	this.tagvalid = false
	
	this.lexeme = make([]byte, 0, 20)
	
	//this.advance()
	//this.advance()
	
	return this
}

func (this *rawlexer) advance() {
	
	if this.index >= len(this.source) {
		// EOS
		this.tagvalid = false
		return
	}
	
	for ; this.index < len(this.source); this.index++ {
		//lookok := len(this.source) - this.index
		
		if this.source[this.index] == '\n' && this.line >= 0 {
			this.line++
			this.column = 0
		} else {
			this.column++
		}
		
		// Lexing Begin
		//======================================
		
		if this.source[this.index] == '[' {
			if this.state != stEatComment {
				panic(fmt.Sprint("Nested Raw Tag found on line: ", this.line))
			}
			
			this.tagbegin = this.index
			this.tagend = this.index
			this.curID = ""
			this.curParams = make([]string, 0)
			this.tagvalid = false
			this.state = stReadTagID
			continue
		}
		
		if this.source[this.index] == ']' {
			if this.state == stReadTagID {
				this.tagend = this.index
				this.curID = string(this.lexeme)
				this.lexeme = this.lexeme[0:0]
				this.tagvalid = true
				this.state = stEatComment
				return
			} else if this.state == stReadParam {
				this.tagend = this.index
				this.curParams = append(this.curParams, string(this.lexeme))
				this.lexeme = this.lexeme[0:0]
				this.tagvalid = true
				this.state = stEatComment
				return
			}
			continue
		}
		
		if this.source[this.index] == ':' {
			if this.state == stReadTagID {
				this.curID = string(this.lexeme)
				this.lexeme = this.lexeme[0:0]
				this.state = stReadParam
			} else if this.state == stReadParam {
				this.curParams = append(this.curParams, string(this.lexeme))
				this.lexeme = this.lexeme[0:0]
			}
			continue
		}
		
		if this.state == stEatComment {
			continue
		}
		this.lexeme = append(this.lexeme, this.source[this.index])
	}
	this.tagvalid = false
}

func (this *rawlexer) Get(index string) *nca7.Value {
	if index == "id" {
		return nca7.NewValueString(this.curID)
	}
	
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		panic("Index not a valid number.")
	}
	if val < 0 || val >= int64(len(this.curParams)) {
		panic("Index out of range.")
	}
	return nca7.NewValueString(this.curParams[val])
}

func (this *rawlexer) Exists(index string) bool {
	if index == "id" {
		return true
	}
	
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		return false
	}
	if val < 0 || val >= int64(len(this.curParams)) {
		return false
	}
	return true
}

func (this *rawlexer) Len() int64 {
	return int64(len(this.curParams))
}

func (this *rawlexer) Keys() []string {
	rtn := make([]string, 0, len(this.curParams))
	for key := range this.curParams {
		rtn = append(rtn, strconv.FormatInt(int64(key), 10))
	}
	return rtn
}
