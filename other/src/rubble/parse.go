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

package rubble

// ParseFile runs the Rubble stage parser on a single file.
func (state *State) ParseFile(input []byte, stage ParseStage, pos *Position) []byte {
	if stage == StgUseCurrent {
		stage = state.Stage
	}
	if stage != StgPreParse && stage != StgParse && stage != StgPostParse {
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
			if _, ok := state.Templates[lex.Look.Lexeme]; !ok {
				state.Log.Println("      Warning: Nonexistent Template Found: " + lex.Look.Lexeme)
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

			out = append(out, state.Templates[name].Call(state, params).Data...)

		case tknINVALID:
			break loop

		default:
			out = append(out, lex.Current.Lexeme...)
		}
	}

	return []byte(out)
}

func stageTemplate(name string, stage ParseStage) bool {
	if len(name) < 1 {
		return false
	}

	if name[0] == '@' {
		return true
	}

	switch stage {
	case StgPreParse:
		return name[0] == '!'

	case StgParse:
		return !(name[0] == '#' || name[0] == '!')

	case StgPostParse:
		return name[0] == '#'

	}
	return false
}

func templateToString(lex *Lexer) string {
	out := "{"
	for {
		lex.Advance()

		if lex.Current.Type == tknTagEnd {
			return out + "}"
		}

		if lex.Current.Type == tknINVALID {
			return out
		}

		out += lex.Current.Lexeme
	}
}

// Expand variables in string using their current values from the state.
// Respects template bodies by not expanding variables inside of one.
// Non-existent variables are skipped unmodified.
func (state *State) ExpandVars(input string) string {
	buf := make([]byte, 0, len(input))

	depth := 0
	x := 0
	for y := 0; y < len(input); y++ {
		if input[y] == '{' {
			depth++
		}
		if input[y] == '}' && depth > 0 {
			depth--
		}

		if input[y] == '$' && y+1 < len(input) && depth == 0 {
			buf = append(buf, input[x:y]...)
			name, w := getVarName(input[y+1:])
			if name == "" {
				buf = append(buf, "${"...)
			} else {
				buf = append(buf, state.VariableData[name]...)
			}
			y += w
			x = y + 1
		}
	}

	return string(buf) + input[x:]
}

func getVarName(input string) (string, int) {
	if input[0] == '{' {
		// Scan to closing brace
		for i := 1; i < len(input); i++ {
			if input[i] == '}' {
				return input[1:i], i + 1
			}
		}
		return "", 2 // Bad syntax
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
