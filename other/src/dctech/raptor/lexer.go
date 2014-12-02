/*
For copyright/license see header in file "doc.go"
*/

package raptor

import "strconv"

// Lexer states
const (
	stReset = iota
	stEatLineComment
	stReadDQStr
	stReadStr
)

// Lexer is a Raptor CodeSource reading from a string.
type Lexer struct {
	look    *Token
	current *Token

	source string

	pos *Position

	index int
	state int

	lexeme []byte

	token    int
	tokenpos *Position

	strdepth int
	objdepth int
}

// Returns a new Lexer.
// pos is forced to be a line, column position.
func NewLexer(input string, pos *Position) *Lexer {
	this := new(Lexer)

	this.source = input

	this.pos = pos.Copy()
	if this.pos.Column == -1 {
		this.pos.Line = 1
		this.pos.Column = 1
	}

	this.index = 0
	this.state = stReset

	this.lexeme = make([]byte, 0, 20)

	this.token = TknINVALID
	this.tokenpos = pos.Copy()

	this.strdepth = 0
	this.objdepth = 0

	this.Advance()
	this.Advance()

	return this
}

// Advance retrieves the next token from the stream.
// For most purposes use GetToken instead.
func (this *Lexer) Advance() {
	if this.index > len(this.source) {
		this.current = this.look
		this.look = NewToken("INVALID", TknINVALID, this.tokenpos)
		return
	}

	for ; this.index < len(this.source); this.index++ {
		dat := this.source
		i := this.index
		lookok := len(this.source) - this.index

		if dat[i] == '\n' && this.pos.Line >= 0 {
			this.pos.Line++
			this.pos.Column = 0
		} else {
			this.pos.Column++
		}

		// Lexing Begin
		//======================================

		// Comments
		if this.state != stReadDQStr {
			if dat[i] == '#' {
				this.state = stEatLineComment
				continue
			}
			if dat[i] == '\n' && this.state == stEatLineComment {
				this.state = stReset
				continue
			}
			if this.state == stEatLineComment {
				continue
			}
		}

		// Escape
		if dat[i] == '\\' && this.state == stReadDQStr {
			if lookok < 1 {
				panic("Unexpected end of stream.")
			}

			switch dat[i+1] {
			case 'n':
				this.lexeme = append(this.lexeme, '\n')
			case 'r':
				this.lexeme = append(this.lexeme, '\r')
			case 't':
				this.lexeme = append(this.lexeme, '\t')
			case '"':
				this.lexeme = append(this.lexeme, '"')
			case '\\':
				this.lexeme = append(this.lexeme, '\\')
			case 'x':
				if lookok < 3 {
					panic("Unexpected end of stream.")
				}
				rep, err := strconv.ParseInt(string([]byte{dat[i+2], dat[i+3]}), 16, 8)
				if err != nil {
					panic("Invalid escape sequence: \\x" + string(dat[i+2]) + string(dat[i+3]) + ".")
				}
				this.lexeme = append(this.lexeme, byte(rep))
				this.index += 2

			default:
				panic("Invalid escape sequence: \\" + string(dat[i+1]) + ".")
			}
			this.index++
			continue
		}

		// Double Quote Strings
		if dat[i] == '"' {
			if this.state == stReadDQStr {
				this.state = stReset
				continue
			}
			this.current = this.look
			this.look = NewToken(string(this.lexeme), this.token, this.tokenpos)
			this.state = stReadDQStr
			this.token = TknString
			this.tokenpos = this.pos.Copy()
			this.lexeme = this.lexeme[0:0]
			this.index++
			return
		}
		if this.state == stReadDQStr {
			this.lexeme = append(this.lexeme, dat[i])
			continue
		}

		// Delimiters
		if dat[i] == '\n' || dat[i] == '\r' || dat[i] == ' ' || dat[i] == '\t' || dat[i] == ',' {
			this.state = stReset
			continue
		}

		// Parentheses
		if dat[i] == '(' {
			this.current = this.look
			this.look = NewToken(string(this.lexeme), this.token, this.tokenpos)
			this.state = stReset
			this.token = TknCmdBegin
			this.tokenpos = this.pos.Copy()
			this.lexeme = this.lexeme[0:1]
			this.lexeme[0] = '('
			this.index++
			return
		}
		if dat[i] == ')' {
			this.current = this.look
			this.look = NewToken(string(this.lexeme), this.token, this.tokenpos)
			this.state = stReset
			this.token = TknCmdEnd
			this.tokenpos = this.pos.Copy()
			this.lexeme = this.lexeme[0:1]
			this.lexeme[0] = ')'
			this.index++
			return
		}

		// Square Brackets
		if dat[i] == '[' {
			this.current = this.look
			this.look = NewToken(string(this.lexeme), this.token, this.tokenpos)
			this.state = stReset
			this.token = TknDerefBegin
			this.tokenpos = this.pos.Copy()
			this.lexeme = this.lexeme[0:1]
			this.lexeme[0] = '['
			this.index++
			return
		}
		if dat[i] == ']' {
			this.current = this.look
			this.look = NewToken(string(this.lexeme), this.token, this.tokenpos)
			this.state = stReset
			this.token = TknDerefEnd
			this.tokenpos = this.pos.Copy()
			this.lexeme = this.lexeme[0:1]
			this.lexeme[0] = ']'
			this.index++
			return
		}

		// Angle Brackets
		if dat[i] == '<' {
			this.objdepth++
			this.current = this.look
			this.look = NewToken(string(this.lexeme), this.token, this.tokenpos)
			this.state = stReset
			this.token = TknObjLitBegin
			this.tokenpos = this.pos.Copy()
			this.lexeme = this.lexeme[0:1]
			this.lexeme[0] = '<'
			this.index++
			return
		}
		if dat[i] == '>' {
			this.objdepth--
			this.current = this.look
			this.look = NewToken(string(this.lexeme), this.token, this.tokenpos)
			this.state = stReset
			this.token = TknObjLitEnd
			this.tokenpos = this.pos.Copy()
			this.lexeme = this.lexeme[0:1]
			this.lexeme[0] = '>'
			this.index++
			return
		}
		if dat[i] == '=' && this.objdepth > 0 {
			this.current = this.look
			this.look = NewToken(string(this.lexeme), this.token, this.tokenpos)
			this.state = stReset
			this.token = TknObjLitSplit
			this.tokenpos = this.pos.Copy()
			this.lexeme = this.lexeme[0:1]
			this.lexeme[0] = '='
			this.index++
			return
		}

		if dat[i] == '{' {
			this.current = this.look
			this.look = NewToken(string(this.lexeme), this.token, this.tokenpos)
			this.state = stReset
			this.token = TknCodeBegin
			this.tokenpos = this.pos.Copy()
			this.lexeme = this.lexeme[0:1]
			this.lexeme[0] = '{'
			this.index++
			return
		}
		if dat[i] == '}' {
			this.current = this.look
			this.look = NewToken(string(this.lexeme), this.token, this.tokenpos)
			this.state = stReset
			this.token = TknCodeEnd
			this.tokenpos = this.pos.Copy()
			this.lexeme = this.lexeme[0:1]
			this.lexeme[0] = '}'
			this.index++
			return
		}

		// Raw Strings
		if this.state == stReadStr {
			this.lexeme = append(this.lexeme, dat[i])
			continue
		}

		this.current = this.look
		this.look = NewToken(string(this.lexeme), this.token, this.tokenpos)
		this.state = stReadStr
		this.token = TknString
		this.tokenpos = this.pos.Copy()
		this.lexeme = this.lexeme[0:1]
		this.lexeme[0] = dat[i]
		this.index++
		return
	}

	if this.index == len(this.source) {
		this.current = this.look
		this.look = NewToken(string(this.lexeme), this.token, this.tokenpos)
		this.index++
		return
	}
}

// CurrentTkn returns the current token.
func (this *Lexer) CurrentTkn() *Token {
	return this.current
}

// LookAhead returns the lookahead token.
func (this *Lexer) LookAhead() *Token {
	return this.look
}
