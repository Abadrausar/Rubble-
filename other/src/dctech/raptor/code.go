/*
For copyright/license see header in file "doc.go"
*/

package raptor

// Code is a generic storage for script code, used by user commands among other things.
// It is valid for the position information to be nil.
type Code struct {
	// Code Values:
	//	0 = TknCmdBegin
	//	1 = TknCmdEnd
	//	2 = TknDerefBegin
	//	3 = TknDerefEnd
	//	4 = TknObjLitBegin
	//	5 = TknObjLitEnd
	//	6 = TknObjLitSplit
	//	7 = TknCodeBegin
	//	8 = TknCodeEnd
	//	9 = unused
	//	>9 = TknString, value - 10 is the StringTable index
	Code        []uint32
	Positions   []*Position
	StringTable []string
}

// NewCode will generate a new Code object from a CodeSource.
func NewCode(input CodeSource) *Code {
	this := new(Code)
	this.Code = make([]uint32, 0, 20)
	this.Positions = make([]*Position, 0, 20)

	nextstring := 0
	strings := make(map[string]int)

	for {
		input.Advance()

		if input.CurrentTkn().Type == TknINVALID {
			break
		}

		if input.CurrentTkn().Type == TknString {
			if _, ok := strings[input.CurrentTkn().Lexeme]; ok {
				this.Code = append(this.Code, uint32(strings[input.CurrentTkn().Lexeme]+10))
			} else {
				strings[input.CurrentTkn().Lexeme] = nextstring
				this.Code = append(this.Code, uint32(strings[input.CurrentTkn().Lexeme]+10))
				nextstring++
			}
		} else {
			this.Code = append(this.Code, uint32(input.CurrentTkn().Type))
		}

		this.Positions = append(this.Positions, input.CurrentTkn().Pos)
	}

	this.StringTable = make([]string, nextstring)
	for str, i := range strings {
		this.StringTable[i] = str
	}
	return this
}

// NewCodeBlock will generate a new Code object from an existing CodeSource,
// reading from the current token (which MUST be a TknCodeBegin) to the matching TknCodeEnd.
// (Nested blocks are allowed)
func NewCodeBlock(input CodeSource) *Code {
	if input.CurrentTkn().Type != TknCodeBegin {
		panic("NewCodeBlock called on a CodeSource that is not entering a block.")
	}

	// Compiling from a Code, take a short cut.
	// This is important for loop performance.
	if tmp, ok := input.(*CodeReader); ok {
		code := tmp.data
		start := tmp.currentToken + 1

		blockDepth := 1
		i := start
		foundend := false
		for ; i < len(code.Code); i++ {
			if code.Code[i] == TknCodeBegin {
				blockDepth++
			}
			if code.Code[i] == TknCodeEnd {
				blockDepth--
				if blockDepth == 0 {
					foundend = true
					break
				}
			}
		}
		if !foundend {
			panic("NewCodeBlock: Unexpected EOS.")
		}

		this := new(Code)
		this.Code = code.Code[start:i]

		if code.Positions != nil {
			this.Positions = code.Positions[start:i]
		}

		this.StringTable = code.StringTable

		tmp.currentToken = i

		return this
	}

	// Compiling from some other CodeSource
	this := new(Code)
	this.Code = make([]uint32, 0, 20)
	this.Positions = make([]*Position, 0, 20)

	blockDepth := 1

	nextstring := 0
	strings := make(map[string]int)

	for {
		input.Advance()

		if input.CurrentTkn().Type == TknINVALID {
			panic("NewCodeBlock: Unexpected TknINVALID.")
		}

		// These MUST fall through! (Unless the end of the block is found)
		if input.CurrentTkn().Type == TknCodeBegin {
			blockDepth++
		}
		if input.CurrentTkn().Type == TknCodeEnd {
			blockDepth--
			if blockDepth == 0 {
				break
			}
		}

		if input.CurrentTkn().Type == TknString {
			if _, ok := strings[input.CurrentTkn().Lexeme]; ok {
				this.Code = append(this.Code, uint32(strings[input.CurrentTkn().Lexeme]+10))
			} else {
				strings[input.CurrentTkn().Lexeme] = nextstring
				this.Code = append(this.Code, uint32(strings[input.CurrentTkn().Lexeme]+10))
				nextstring++
			}
		} else {
			this.Code = append(this.Code, uint32(input.CurrentTkn().Type))
		}

		this.Positions = append(this.Positions, input.CurrentTkn().Pos)
	}

	this.StringTable = make([]string, nextstring)
	for str, i := range strings {
		this.StringTable[i] = str
	}
	return this
}

// String should convert a Code object back to a string containing valid Raptor code.
func (this *Code) String() string {
	out := "{ "
	for _, code := range this.Code {
		if code < 9 {
			out += TknLexemes[code] + " "
			continue
		}
		if code > 9 {
			out += "\"" + EscapeString(this.StringTable[code-10]) + "\" "
			continue
		}
	}
	return out + "}"
}
