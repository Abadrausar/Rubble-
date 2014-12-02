package main

//import "fmt"

func StageParse(input string) string {
	if ParseStage == 0 {
		return PreParse(input)
	} else if ParseStage == 1 {
		return Parse(input)
	} else if ParseStage == 2 {
		return PostParse(input)
	}
	panic("Invalid ParseStage")
}

func PreParse(input string) string {
	out := ""
	lex := NewLexer(input)
	
	for {
		lex.Advance()
		if lex.Current.Type == tknString {
			out += lex.Current.Lexeme
		} else if lex.Current.Type == tknTagBegin {
			lex.GetToken(tknString)
			if lex.Current.Lexeme[0] != '!' {
				// Not a pre tag, copy over until we get a tag end
				out += "{" + lex.Current.Lexeme
				for lex.Advance() {
					if lex.Current.Type == tknTagEnd {
						out += lex.Current.Lexeme
						break
					}
					out += lex.Current.Lexeme
				}
				continue
			}
			name := lex.Current.Lexeme
			params := make([]string, 0, 5)
			for lex.CheckLookAhead(tknDelimiter) {
				lex.GetToken(tknDelimiter)
				lex.GetToken(tknString)
				params = append(params, lex.Current.Lexeme)
			}
			lex.GetToken(tknTagEnd)
			
			if _, ok := Templates[name]; !ok {
				panic("Invalid template: " + name)
			}
			out += Templates[name].Call(params)
		} else if lex.Current.Type == tknINVALID {
			break
		}
	}
	
	return out
}

func Parse(input string) string {
	out := ""
	lex := NewLexer(input)
	
	for {
		lex.Advance()
		if lex.Current.Type == tknString {
			out += lex.Current.Lexeme
		} else if lex.Current.Type == tknTagBegin {
			lex.GetToken(tknString)
			if lex.Current.Lexeme[0] == '#' {
				// Post tag, copy over until we get a tag end
				out += "{" + lex.Current.Lexeme
				for lex.Advance() {
					if lex.Current.Type == tknTagEnd {
						out += lex.Current.Lexeme
						break
					}
					out += lex.Current.Lexeme
				}
				continue
			}
			name := lex.Current.Lexeme
			params := make([]string, 0, 5)
			for lex.CheckLookAhead(tknDelimiter) {
				lex.GetToken(tknDelimiter)
				lex.GetToken(tknString)
				params = append(params, lex.Current.Lexeme)
			}
			lex.GetToken(tknTagEnd)
			
			if _, ok := Templates[name]; !ok {
				panic("Invalid template: " + name)
			}
			out += Templates[name].Call(params)
		} else if lex.Current.Type == tknINVALID {
			break
		}
	}
	
	return out
}

func PostParse(input string) string {
	out := ""
	lex := NewLexer(input)
	
	for {
		lex.Advance()
		if lex.Current.Type == tknString {
			out += lex.Current.Lexeme
		} else if lex.Current.Type == tknTagBegin {
			lex.GetToken(tknString)
			if lex.Current.Lexeme[0] != '#' {
				// Not a post tag, copy over until we get a tag end
				out += "{" + lex.Current.Lexeme
				for lex.Advance() {
					if lex.Current.Type == tknTagEnd {
						out += lex.Current.Lexeme
						break
					}
					out += lex.Current.Lexeme
				}
				continue
			}
			name := lex.Current.Lexeme
			params := make([]string, 0, 5)
			for lex.CheckLookAhead(tknDelimiter) {
				lex.GetToken(tknDelimiter)
				lex.GetToken(tknString)
				params = append(params, lex.Current.Lexeme)
			}
			lex.GetToken(tknTagEnd)
			
			if _, ok := Templates[name]; !ok {
				panic("Invalid template: " + name)
			}
			out += Templates[name].Call(params)
		} else if lex.Current.Type == tknINVALID {
			break
		}
	}
	
	return out
}
