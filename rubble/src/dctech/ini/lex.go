package ini

//Copyright 2013 by Milo Christiansen
//
//DCTech Project License
//	You may not modify the project without prior written permission.
//	You may not redistribute the project without prior written permission.
//	
//	The project is provided 'as-is', without any express or implied warranty. 
//	In no event will the author be held liable for any damages arising from the use of the project.
//	
//	In the event of a dispute about the interpretation of this license you agree that the 
//	interpretation of the author is the correct one.
// 
// Permission is given to include this package in your projects so long as the above license 
// and this exception are reproduced somewhere in your project.
// 

//import "fmt"
//import "strconv"

// Lexer states
const (
	stReset = iota
	stReadComment
	stReadKey
	stReadValue
	stReadSection
)

// Token values
const (
	TknINVALID = iota
	TknKey
	TknValue
	TknSection
	TknComment
)

// A semi-generic Lexer framework
type Lexer struct {

	// The next/current tokens
	Look    Token
	Current Token

	// The line and column of the last token
	// Used by the parser to generate errors
	PositionValid bool
	LastLine      int

	// The token stream
	stream <-chan Token
}

// Returns a new Lexer.
func NewLexer(input string, startline int) *Lexer {
	lexer := new(Lexer)

	if startline >= 0 {
		lexer.PositionValid = true
	} else {
		startline = 0
		lexer.PositionValid = false
	}

	out := make(chan Token)

	go func() {
		line := startline

		token := TknINVALID
		tokenline := 0

		state := stReset
		sectiondepth := 0

		lexeme := make([]rune, 0, 100)

		for _, val := range input {
			if val == '\n' {
				line++
			}

			if state == stReadValue && val != '\n' {
				lexeme = append(lexeme, val)
				continue
			}

			if val == '[' {
				if state == stReadSection {
					lexeme = append(lexeme, val)
					sectiondepth++
					continue
				}
				out <- Token{string(lexeme), token, tokenline}
				state = stReadSection
				token = TknSection
				tokenline = line
				lexeme = lexeme[0:0]
				continue
			}
			if state == stReadSection && val == ']' {
				if sectiondepth > 0 {
					sectiondepth--
					continue
				}
				state = stReset
				continue
			}
			if state == stReadSection {
				lexeme = append(lexeme, val)
				continue
			}

			if state == stReadKey && val == '=' {
				out <- Token{string(lexeme), token, tokenline}
				state = stReadValue
				token = TknValue
				tokenline = line
				lexeme = lexeme[0:0]
				continue
			}
			if state == stReadKey {
				lexeme = append(lexeme, val)
				continue
			}

			if val == '#' {
				out <- Token{string(lexeme), token, tokenline}
				state = stReadComment
				token = TknComment
				tokenline = line
				lexeme = lexeme[0:0]
				continue
			}
			if state == stReadComment && val == '\n' {
				state = stReset
				continue
			}
			if state == stReadComment {
				lexeme = append(lexeme, val)
				continue
			}

			if val == '\n' || val == '\r' || val == ' ' || val == '\t' {
				state = stReset
				continue
			}

			out <- Token{string(lexeme), token, tokenline}
			state = stReadKey
			token = TknKey
			tokenline = line
			lexeme = lexeme[0:1]
			lexeme[0] = val
		}

		out <- Token{string(lexeme), token, tokenline}
		close(out)
	}()

	_ = <-out

	lexer.stream = out

	lexer.Advance()
	return lexer
}

// This advances the Lexer one token.
// May cause a panic if the lexer encounters an error
// For most purposes use GetToken instead.
func (this *Lexer) Advance() {
	this.Current = this.Look

	this.Look = <-this.stream
	this.LastLine = this.Current.Line

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

	exitOnTokenExpected(this.Current, tokenTypes...)
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
