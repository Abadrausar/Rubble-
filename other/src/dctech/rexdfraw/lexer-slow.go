/*
Copyright 2014 by Milo Christiansen

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

package rexdfraw

import "dctech/rex"

// Lexer states, use the ones from the other lexer
//const (
//	stEatComment = iota
//	stReadTagID
//	stReadParam
//)

func LexFile(input []byte) rex.Indexable {

	out := make([]*rex.Value, 0)

	lexeme := make([]byte, 0, 20)
	rawtag := make([]byte, 0, 40)
	comments := make([]byte, 0, 100)
	state := stEatComment
	line := 1
	var curTag *Tag

	for _, val := range input {
		if val == '\n' {
			line++
		}

		if val == '[' {
			// Start new tag
			if state != stEatComment {
				comments = append(comments, rawtag...)
			}
			rawtag = rawtag[0:1]
			rawtag[0] = '['
			
			if curTag != nil {
				curTag.Comments = string(comments)
				out = append(out, rex.NewValueIndex(curTag))
			}
			comments = comments[0:0]
			curTag = NewTag()
			state = stReadTagID
			continue
		}
		if val == ']' {
			// Close out tag
			if state == stReadTagID {
				curTag.ID = string(lexeme)
				lexeme = lexeme[0:0]
				state = stEatComment
			} else if state == stReadParam {
				curTag.Params.(rex.EditIndexable).Set("append", rex.NewValueString(string(lexeme)))
				lexeme = lexeme[0:0]
				state = stEatComment
			}
			continue
		}
		if val == ':' {
			// Start new param
			if state == stReadTagID {
				curTag.ID = string(lexeme)
				rawtag = append(rawtag, lexeme...)
				lexeme = lexeme[0:0]
				state = stReadParam
			} else if state == stReadParam {
				curTag.Params.(rex.EditIndexable).Set("append", rex.NewValueString(string(lexeme)))
				rawtag = append(rawtag, lexeme...)
				lexeme = lexeme[0:0]
			}
			continue
		}

		if state == stEatComment {
			comments = append(comments, val)
			continue
		}

		// Add char to lexeme
		lexeme = append(lexeme, val)
	}
	if curTag != nil {
		curTag.Comments = string(comments)
		comments = comments[0:0]
		out = append(out, rex.NewValueIndex(curTag))
	}

	return rex.NewArray(out)
}
