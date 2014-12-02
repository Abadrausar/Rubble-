/*
Copyright 2013-2014 by Milo Christiansen

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

package main

// The parser
func Parse(input []byte, stage int, pos *Position) []byte {
	if stage == stgUseCurrent {
		stage = ParseStage
	}
	if stage != stgPreParse && stage != stgParse && stage != stgPostParse {
		panic("Parse called with invalid parse stage")
	}

	// 100k sounds like a lot, but there are vanilla raw files that are almost 400k
	// Most seem to be around 50k-70k so 100k is not too high
	out := make([]byte, 0, 102400)
	lex := NewLexer(input, pos)

loop:
	for {
		lex.Advance()
		switch lex.Current.Type {
		case tknTagBegin:
			if !stageTemplate(lex.Look.Lexeme, stage) {
				out = append(out, templateToString(lex)...)
				continue
			}

			lex.GetToken(tknString)
			name := lex.Current.Lexeme
			params := make([]*Value, 0, 5)
			for lex.CheckLookAhead(tknDelimiter) {
				lex.GetToken(tknDelimiter)
				if lex.CheckLookAhead(tknString) {
					lex.GetToken(tknString)
					params = append(params, lex.Current.Value())
				} else {
					params = append(params, NewValue(""))
				}
			}
			lex.GetToken(tknTagEnd)

			if _, ok := Templates[name]; !ok {
				panic("Invalid template: " + name)
			}
			out = append(out, Templates[name].Call(params).Data...)

		case tknINVALID:
			break loop

		default:
			out = append(out, lex.Current.Lexeme...)
		}
	}

	return []byte(out)
}

func stageTemplate(name string, stage int) bool {
	if len(name) < 1 {
		panic("Invalid template name.")
	}

	if name[0] == '@' {
		return true
	}

	switch stage {
	case stgPreParse:
		return name[0] == '!'

	case stgParse:
		return !(name[0] == '#' || name[0] == '!')

	case stgPostParse:
		return name[0] == '#'

	}
	return false
}

// Takes a lexer where the current token is a tknTagBegin and the lookahead is the template name.
// Returns a string form of the template (and any nested templates).
// The lexer is advanced so the current token is the template's tknTagEnd.
func templateToString(lex *Lexer) string {
	if lex.Current.Type != tknTagBegin {
		panic("templateToString: Bad beginning.")
	}

	out := "{"
	for {
		lex.Advance()

		if lex.Current.Type == tknTagEnd {
			return out + "}"
		}

		if lex.Current.Type == tknINVALID {
			panic("templateToString: tknINVALID.")
		}

		out += lex.Current.Lexeme
	}
	panic("templateToString: Unexpected EOS.")
}

// This is a modified version of os.Expand
func ExpandVars(input string) string {
	buf := make([]byte, 0, len(input))

	depth := 0
	i := 0
	for j := 0; j < len(input); j++ {
		if input[j] == '{' {
			depth++
		}
		if input[j] == '}' && depth > 0 {
			depth--
		}

		if input[j] == '$' && j+1 < len(input) && depth == 0 {
			buf = append(buf, input[i:j]...)
			name, w := getVarName(input[j+1:])
			if name == "" {
				buf = append(buf, '{')
			} else {
				buf = append(buf, VariableData[name]...)
			}
			j += w
			i = j + 1
		}
	}

	return string(buf) + input[i:]
}

func getVarName(input string) (string, int) {
	if input[0] == '{' {
		// Scan to closing brace
		for i := 1; i < len(input); i++ {
			if input[i] == '}' {
				return input[1:i], i + 1
			}
		}
		return "", 1 // Bad syntax
	}
	// Scan alphanumerics.
	var i int
	for i = 0; i < len(input) && isAlphaNum(input[i]); i++ {
	}
	return input[:i], i
}

func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}
