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

package rex

// token values
// These are order/value sensitive!
const (
	tknINVALID     = iota - 1 // Invalid
	tknCmdBegin               // '('
	tknCmdEnd                 // ')'
	tknVarBegin               // '['
	tknVarEnd                 // ']'
	tknObjLitBegin            // '<'
	tknObjLitEnd              // '>'
	tknNameSplit              // ':'
	tknAssignment             // '='
	
	// Below this point are values unique to tokens
	
	tknCodeBegin              // '{'
	tknCodeEnd                // '}'
	
	tknRawString              // A raw (undelimited) string, normally an identifier or number
	tknString                 // A DQ, SQ, or BQ string
	tknTrue                   // 'true'
	tknFalse                  // 'false'
	tknNil                    // 'nil'
	tknVariadic               // '...'
	
	tknDeclModule             // 'module'
	tknDeclCommand            // 'command'
	tknDeclBlock              // 'block'
	tknDeclVar                // 'var'
)

// token is used to represent the output from the lexer.
// In most cases Lexeme is blank
type token struct {
	Lexeme string
	Type   int
	Pos    *Position
}

// String returns the tokens type as a string.
func (this *token) String() string {
	return tokenTypeToString(this.Type)
}

// OpCode will convert a lexer token to an OpCode ready to insert into the code stream.
// Will produce invalid OpCodes for some token types.
func (tok *token) opCode() *opCode {
	if tok.Type >= tknRawString && tok.Type <= tknVariadic {
		return &opCode{
			Value: tokenToValue(tok),
			Type:  opValue,
			Pos:  tok.Pos.Copy(),
		}
	}

	// this should handle all the types that match up.
	return &opCode{
		Type: tok.Type,
		Pos:  tok.Pos.Copy(),
	}
}

func tokenTypeToString(tokenType int) string {
	switch tokenType {
	case tknINVALID:
		return "tknINVALID"
	case tknCmdBegin:
		return "tknCmdBegin"
	case tknCmdEnd:
		return "tknCmdEnd"
	case tknVarBegin:
		return "tknVarBegin"
	case tknVarEnd:
		return "tknVarEnd"
	case tknObjLitBegin:
		return "tknObjLitBegin"
	case tknObjLitEnd:
		return "tknObjLitEnd"
	case tknCodeBegin:
		return "tknCodeBegin"
	case tknCodeEnd:
		return "tknCodeEnd"
	case tknNameSplit:
		return "tknNameSplit"
	case tknAssignment:
		return "tknAssignment"
	case tknRawString:
		return "tknRawString"
	case tknString:
		return "tknString"
	case tknTrue:
		return "tknTrue"
	case tknFalse:
		return "tknFalse"
	case tknNil:
		return "tknNil"
	case tknVariadic:
		return "tknVariadic"
	case tknDeclModule:
		return "tknDeclModule"
	case tknDeclCommand:
		return "tknDeclCommand"
	case tknDeclBlock:
		return "tknDeclBlock"
	case tknDeclVar:
		return "tknDeclVar"
	}

	RaiseError("token type value out of range")
	return ""
}

// Panics with a message formatted like one of the following:
//	Invalid token: Found: thecurrenttoken. Expected: expected1, expected2, or expected3.
//	Invalid token: Found: thecurrenttoken. Expected: expected1 or expected2.
//	Invalid token: Found: thecurrenttoken. Expected: expected.
//	Invalid token: Found: thecurrenttoken (Lexeme: test). Expected: expected1, expected2, or expected3.
//	Invalid token: Found: thecurrenttoken (Lexeme: test). Expected: expected1 or expected2.
//	Invalid token: Found: thecurrenttoken (Lexeme: test). Expected: expected.
// If the lexeme is long it is truncated.
func exitOnTokenExpected(token *token, expected ...int) {
	expectedString := ""
	expectedCount := len(expected) - 1
	for i, val := range expected {
		// Is the only value
		if expectedCount == 0 {
			expectedString = tokenTypeToString(val)
			continue
		}

		// Is last of a list (2 or more)
		if i == expectedCount && expectedCount > 0 {
			expectedString += "or " + tokenTypeToString(val)
			continue
		}

		// Is the first of two
		if expectedCount == 1 {
			expectedString += tokenTypeToString(val) + " "
			continue
		}

		// Is any but the last of a list of 3 or more
		expectedString += tokenTypeToString(val) + ", "
	}

	found := token.String()
	if token.Lexeme != "" {
		if len(token.Lexeme) < 20 {
			found += " (Lexeme: " + token.Lexeme + ")"
		} else {
			found += " (Lexeme: " + token.Lexeme[:17] + "...)"
		}
	}
	RaiseError("Invalid token: Found: " + found + ". Expected: " + expectedString)
}
