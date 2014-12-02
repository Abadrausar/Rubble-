package main

import "fmt"
import "regexp"
import "math/rand"
import "strings"
import "io/ioutil"

func ReadConfig(path, namespace string) {
	fmt.Println("Reading Config File:", path)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	
	lines := strings.Split(string(file), "\n")
	for i := range lines {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "#") {
			continue
		}
		if strings.TrimSpace(lines[i]) == "" {
			continue
		}
		
		parts := strings.SplitN(lines[i], "=", 2)
		if len(parts) != 2 {
			panic("Malformed config line.")
		}
		
		name := namespace + "__CONFIG_" + strings.TrimSpace(parts[0])
		varData[name] = strings.TrimSpace(parts[1])
	}
}

var atAtAtRegEx = regexp.MustCompile("@@@@*")
var nameAtRegEx = regexp.MustCompile("[a-zA-Z_]+@")

func PreProcess(input, namespace, filename string) string {
	// Replace @@@ (or @@@@, @@@@@...) with a unique identifier. "GEN_00000000_"
	input = atAtAtRegEx.ReplaceAllStringFunc(input, func(in string) string {
		return fmt.Sprintf("GEN_%X_", int32(rand.Int31()))
	})
	
	// Replace @@ with a file-unique identifier. "namespace__filename__"
	input = strings.Replace(input, "@@", fmt.Sprintf("%v__%v__", namespace, filename), -1)
	
	// Replace namespace@ with a unique identifier of namespace. "namespace__"
	input = nameAtRegEx.ReplaceAllStringFunc(input, func(in string) string {
		in = strings.TrimLeft(in, "@")
		return fmt.Sprintf("%v__", in)
	})
	
	// Replace @ with a unique identifier of the base or addon this file belongs to. "namespace__"
	input = strings.Replace(input, "@", fmt.Sprintf("%v__", namespace), -1)
	
	return input
}

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
