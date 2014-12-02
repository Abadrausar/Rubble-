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

// Rubble Tileset Package.
// This package has functions geared towards reading raw files and
// applying tileset information to them. Also has functions used to
// generate the tileset files.
package tset

// Lexer states
const (
	stEatComment = iota
	stReadTagID
	stReadParam
)

func ParseRaws(input []byte) []*Tag {
	out := make([]*Tag, 0)

	lexeme := make([]byte, 0, 20)
	rawtag := make([]byte, 0, 40)
	comments := make([]byte, 0, 100)
	state := stEatComment
	var curTag *Tag
	
	line := 1
	column := 1

	for _, val := range input {
		// Unused, but better to have it and not use it than not to have it at all...
		if val == '\n' {
			line++
			column = 1
		} else {
			column++
		}

		if val == '[' {
			// Start new tag
			if state != stEatComment {
				comments = append(comments, rawtag...)
				comments = append(comments, lexeme...)
			}
			rawtag = rawtag[0:1]
			rawtag[0] = '['
			
			if curTag != nil {
				curTag.Comments = string(comments)
				out = append(out, curTag)
			}
			comments = comments[0:0]
			curTag = new(Tag)
			state = stReadTagID
			continue
		}
		if val == ']' {
			// Close out tag
			if state == stReadTagID {
				curTag.ID = string(lexeme)
				lexeme = lexeme[0:0]
				state = stEatComment
				continue
			} else if state == stReadParam {
				curTag.Params = append(curTag.Params, string(lexeme))
				lexeme = lexeme[0:0]
				state = stEatComment
				continue
			}
		}
		if val == ':' {
			// Start new param
			if state == stReadTagID {
				curTag.ID = string(lexeme)
				rawtag = append(rawtag, lexeme...)
				rawtag = append(rawtag, ':')
				lexeme = lexeme[0:0]
				state = stReadParam
				continue
			} else if state == stReadParam {
				curTag.Params = append(curTag.Params, string(lexeme))
				rawtag = append(rawtag, lexeme...)
				rawtag = append(rawtag, ':')
				lexeme = lexeme[0:0]
				continue
			}
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
		out = append(out, curTag)
	}

	return out
}
