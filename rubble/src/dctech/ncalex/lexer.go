package ncalex

import "fmt"
import "strconv"

// Lexer states
const (
	stReset = iota
	stEatComment
	stReadDQStr
	stReadCBStr
	stReadStr
)

// Token values
const (
	TknINVALID = iota
	TknString
	TknOpenParen
	TknCloseParen
	TknOpenSqBracket
	TknCloseSqBracket
)

// A semi-generic Lexer framework
type Lexer struct {

	// The next/current tokens
	Look    *Token
	Current *Token

	// The line and column of the last token
	// Used by the parser to generate errors
	PositionValid bool
	LastLine      int
	LastColumn    int

	// The token stream
	stream    <-chan *Token
	errstream <-chan lexError
}

type lexError struct {
	line   int
	column int
	text   string
}

// Returns a new Lexer.
func NewLexer(input string, startline, startcolumn int) *Lexer {
	lexer := new(Lexer)

	if startline >= 0 && startcolumn >= 0 {
		lexer.PositionValid = true
		lexer.LastLine = startline
		lexer.LastColumn = startcolumn
	} else {
		startline = 0
		startcolumn = 0
		lexer.PositionValid = false
	}

	out := make(chan *Token)
	outerr := make(chan lexError)

	go func() {
		line := startline
		column := startcolumn

		token := TknINVALID
		tokenline := startline
		tokencolumn := startcolumn

		state := stReset
		strdepth := 0
		escape := false
		xscape := 2 // extended excape, eg \xff
		xscapeStr := make([]rune, 2)

		lexeme := make([]rune, 0, 20)

		for _, val := range input {
			if val == '\n' {
				line++
				column = 0
			} else {
				column++
			}

			// When reading a delimited token of any type state transitions are not allowed!
			// delimited tokens are: double quote strings, curly bracket strings, and comments.

			if xscape < 2 {
				xscapeStr[xscape] = val
				xscape++
				if xscape == 2 {
					i, err := strconv.ParseInt(string(xscapeStr), 16, 32)
					if err != nil {
						// This should always cause a panic in Advance()
						errstr := fmt.Sprintf("Invalid escape sequence: \\x%v.", string(xscapeStr))
						outerr <- lexError{line, column - 4, errstr}
					}
					lexeme = append(lexeme, rune(i))
				}
				continue
			}

			if escape == true && state == stReadDQStr {
				switch val {
				case 'n':
					lexeme = append(lexeme, '\n')
				case 'r':
					lexeme = append(lexeme, '\r')
				case 't':
					lexeme = append(lexeme, '\t')
				case '\'':
					lexeme = append(lexeme, '\'')
				case '"':
					lexeme = append(lexeme, '"')
				case '\\':
					lexeme = append(lexeme, '\\')
				case 'x':
					xscape = 0
					escape = false
					continue
				default:
					// This should always cause a panic in Advance()
					err := fmt.Sprintf("Invalid escape sequence: \\%v.", string(val))
					outerr <- lexError{line, column - 2, err}
				}
			}

			if escape == true {
				escape = false
				continue
			}

			if val == '\\' && (state == stReadDQStr || state == stEatComment) {
				escape = true
				continue
			}

			if state != stReadCBStr && state != stEatComment {
				if val == '"' {
					if state == stReadDQStr {
						state = stReset
						continue
					}
					out <- &Token{string(lexeme), token, tokenline, tokencolumn}
					state = stReadDQStr
					token = TknString
					tokenline = line
					tokencolumn = column
					lexeme = lexeme[0:0]
					continue
				}
				if state == stReadDQStr {
					lexeme = append(lexeme, val)
					continue
				}
			}

			if state != stReadDQStr && state != stEatComment {
				if val == '{' {
					strdepth++
					if state == stReadCBStr {
						lexeme = append(lexeme, val)
						continue
					}
					out <- &Token{string(lexeme), token, tokenline, tokencolumn}
					state = stReadCBStr
					token = TknString
					tokenline = line
					tokencolumn = column
					lexeme = lexeme[0:0]
					continue
				}
				if val == '}' {
					strdepth--
					if strdepth == 0 {
						state = stReset
						continue
					}
					lexeme = append(lexeme, val)
					continue
				}
				if state == stReadCBStr {
					lexeme = append(lexeme, val)
					continue
				}
			}

			if val == '\'' {
				if state == stEatComment {
					state = stReset
					continue
				}
				state = stEatComment
				continue
			}
			if state == stEatComment {
				continue
			}

			if val == '\n' || val == '\r' || val == ' ' || val == '\t' || val == ',' {
				state = stReset
				continue
			}

			if val == '(' {
				out <- &Token{string(lexeme), token, tokenline, tokencolumn}
				state = stReset
				token = TknOpenParen
				tokenline = line
				tokencolumn = column
				lexeme = lexeme[0:1]
				lexeme[0] = '('
				continue
			}
			if val == ')' {
				out <- &Token{string(lexeme), token, tokenline, tokencolumn}
				state = stReset
				token = TknCloseParen
				tokenline = line
				tokencolumn = column
				lexeme = lexeme[0:1]
				lexeme[0] = ')'
				continue
			}

			if val == '[' {
				out <- &Token{string(lexeme), token, tokenline, tokencolumn}
				state = stReset
				token = TknOpenSqBracket
				tokenline = line
				tokencolumn = column
				lexeme = lexeme[0:1]
				lexeme[0] = '['
				continue
			}
			if val == ']' {
				out <- &Token{string(lexeme), token, tokenline, tokencolumn}
				state = stReset
				token = TknCloseSqBracket
				tokenline = line
				tokencolumn = column
				lexeme = lexeme[0:1]
				lexeme[0] = ']'
				continue
			}

			if state == stReadStr {
				lexeme = append(lexeme, val)
				continue
			}

			out <- &Token{string(lexeme), token, tokenline, tokencolumn}
			state = stReadStr
			token = TknString
			tokenline = line
			tokencolumn = column
			lexeme = lexeme[0:1]
			lexeme[0] = val
		}

		out <- &Token{string(lexeme), token, tokenline, tokencolumn}
		close(out)
	}()

	_ = <-out

	lexer.stream = out
	lexer.errstream = outerr

	lexer.Look = new(Token)
	lexer.Current = new(Token)

	lexer.Advance()
	return lexer
}

// This advances the Lexer one token.
// May cause a panic if the lexer encounters an error
// For most purposes use GetToken instead.
func (this *Lexer) Advance() {
	this.Current = this.Look

	select {
	case this.Look = <-this.stream:
		if this.Look == nil {
			this.Look = &Token{"INVALID", TknINVALID, this.LastLine, this.LastColumn}
		}

		this.LastLine = this.Current.Line
		this.LastColumn = this.Current.Column
	case err := <-this.errstream:
		this.LastLine = err.line
		this.LastColumn = err.column
		panic(err.text)
	}

	// Uncomment this for debugging tokens
	//fmt.Println(this.Current)
}

// Gets the next token, and panics with an error if it's not of type tokenType.
// May cause a panic if the lexer encounters an error
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
