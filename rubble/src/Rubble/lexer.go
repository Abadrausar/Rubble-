/*
Copyright 2013 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. You may not alter this software in any way.

3. This notice may not be removed or altered from any source distribution.
*/

package main

//import "fmt"

// Lexer states
const (
	stReadString = iota
	stReadCommand
)

// A semi-generic Lexer framework
type Lexer struct {

	// The next/current tokens
	Look    *Token
	Current *Token
	
	// The token stream
	stream    <-chan *Token
}

// Returns a new Lexer.
func NewLexer(input string) *Lexer {
	out := make(chan *Token)
	
	go func() {
		
		token := tknString

		state := stReadString
		
		lexeme := make([]byte, 0, 20)
		commandDepth := 0
		
		line := 1
		tknline := 1

		for i := 0; i < len(input); i++ {
			if input[i] == '\n' {
				line++
			}
			
			if input[i] == ';' {
				if state == stReadString || commandDepth > 0 {
					lexeme = append(lexeme, input[i])
					continue
				}
				if 0 < i && len(input) > i+1 {
					if input[i-1] == '\'' && input[i+1] == '\'' {
						lexeme = append(lexeme, input[i])
						continue
					}
				}
				
				out <- &Token{string(lexeme), token, tknline}
				out <- &Token{ ";", tknDelimiter, tknline}
				token = tknString
				lexeme = lexeme[0:0]
				tknline = line
				continue
			}
			
			if input[i] == '{' {
				if 0 < i && len(input) > i+1 {
					if input[i-1] == '\'' && input[i+1] == '\'' {
						lexeme = append(lexeme, input[i])
						continue
					}
				}
				if state == stReadString {
					state = stReadCommand
					out <- &Token{string(lexeme), token, tknline}
					out <- &Token{ "{", tknTagBegin, tknline}
					token = tknString
					lexeme = lexeme[0:0]
					tknline = line
					continue
				}
				commandDepth++
				lexeme = append(lexeme, input[i])
				continue
			}
			if input[i] == '}' {
				if 0 < i && len(input) > i+1 {
					if input[i-1] == '\'' && input[i+1] == '\'' {
						lexeme = append(lexeme, input[i])
						continue
					}
				}
				if state == stReadCommand && commandDepth == 0 {
					state = stReadString
					out <- &Token{string(lexeme), token, tknline}
					out <- &Token{ "}", tknTagEnd, tknline}
					token = tknString
					lexeme = lexeme[0:0]
					tknline = line
					continue
				}
				commandDepth--
				lexeme = append(lexeme, input[i])
				continue
			}

			lexeme = append(lexeme, input[i])
		}
		
		out <- &Token{string(lexeme), token, tknline}
		close(out)
	}()
	
	lexer := new(Lexer)
	lexer.stream = out
	lexer.Look = new(Token)
	lexer.Current = new(Token)
	lexer.Advance()
	return lexer
}

// This advances the Lexer one token.
// For most purposes use GetToken instead.
func (this *Lexer) Advance() bool {
	this.Current = this.Look
	this.Look = <-this.stream
	
	LastLine = this.Current.Line
	
	if this.Look == nil {
		this.Look = new(Token)
	}
	
	return this.Current.Type != tknINVALID
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
