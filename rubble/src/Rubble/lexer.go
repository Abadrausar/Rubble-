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

		for i := 0; i < len(input); i++ {
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
				
				out <- &Token{string(lexeme), token}
				out <- &Token{ ";", tknDelimiter}
				token = tknString
				lexeme = lexeme[0:0]
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
					out <- &Token{string(lexeme), token}
					out <- &Token{ "{", tknTagBegin}
					token = tknString
					lexeme = lexeme[0:0]
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
					out <- &Token{string(lexeme), token}
					out <- &Token{ "}", tknTagEnd}
					token = tknString
					lexeme = lexeme[0:0]
					continue
				}
				commandDepth--
				lexeme = append(lexeme, input[i])
				continue
			}

			lexeme = append(lexeme, input[i])
		}
		
		out <- &Token{string(lexeme), token}
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
