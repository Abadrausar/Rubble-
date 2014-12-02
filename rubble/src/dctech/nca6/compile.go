package nca6

//import "fmt"
import "dctech/ncalex"

// CompiledScript is a compiled NCA script.
type CompiledScript struct {
	// Code Values:
	//	0 = TknOpenParen
	//	1 = TknCloseParen
	//	2 = TknOpenSqBracket
	//	3 = TknCloseSqBracket
	//	4 = Newline, this is never seen by the parser, it is for error reporting only.
	//	>4 = TknString, value - 5 is the StringTable index
	Code []uint32
	StringTable []string
	Start int
}

func Compile(input string, line int) *CompiledScript {
	this := new(CompiledScript)
	this.Code = make([]uint32, 0)
	this.Start = line
	
	lex := ncalex.NewLexer(input, line, 0)
	
	nextstring := 0
	strings := make(map[string]int)
	lastLine := line
	
	for {
		lex.Advance()
		
		for lex.Current.Line > lastLine {
			this.Code = append(this.Code, 4)
			lastLine++
		}
		
		switch lex.Current.Type {
		case ncalex.TknString:
			if _, ok := strings[lex.Current.Lexeme]; ok {
				this.Code = append(this.Code, uint32(strings[lex.Current.Lexeme]+5))
			} else {
				strings[lex.Current.Lexeme] = nextstring
				this.Code = append(this.Code, uint32(strings[lex.Current.Lexeme]+5))
				nextstring++
			}
		case ncalex.TknOpenParen:
			this.Code = append(this.Code, 0)
		case ncalex.TknCloseParen:
			this.Code = append(this.Code, 1)
		case ncalex.TknOpenSqBracket:
			this.Code = append(this.Code, 2)
		case ncalex.TknCloseSqBracket:
			this.Code = append(this.Code, 3)
		}
		
		if lex.CheckLookAhead(ncalex.TknINVALID) {
			break
		}
	}
	
	this.StringTable = make([]string, nextstring)
	for str, i := range strings {
		this.StringTable[i] = str
	}
	return this
}

// CompiledLexer is a CodeSource wrapper for CompiledScript.
type CompiledLexer struct {
	currentToken int
	line int
	data *CompiledScript
}

func NewCompiledLexer(val *CompiledScript) CodeSource {
	this := new(CompiledLexer)
	this.data = val
	this.line = val.Start
	return this
}

func (this *CompiledLexer) Line() int {
	return this.line
}

func (this *CompiledLexer) CurrentTkn() *ncalex.Token {
	return this.Token(this.currentToken)
}

func (this *CompiledLexer) LookAhead() *ncalex.Token {
	token := this.currentToken
	for {
		token++
		if token >= len(this.data.Code) || token < 0 {
			return this.Token(token)
		}
		code := this.data.Code[token]
		if code == 4 {
			continue
		}
		return this.Token(token)
	}
	panic("UNREACHABLE: *CompiledLexer.LookAhead")
}

// This advances the code one token.
// For most purposes use GetToken instead.
func (this *CompiledLexer) Advance() {
	for {
		this.currentToken++
		if this.currentToken >= len(this.data.Code) || this.currentToken < 0 {
			return
		}
		code := this.data.Code[this.currentToken]
		if code == 4 {
			this.line++
			continue
		}
		return
	}
}

// Gets the next token, and panics with an error if it's not of type tokenType.
// Used as a type checked Advance.
func (this *CompiledLexer) GetToken(tokenTypes ...int) {
	this.Advance()
	current := this.CurrentTkn()
	
	for _, val := range tokenTypes {
		if current.Type == val {
			return
		}
	}

	ncalex.ExitOnTokenExpected(current, tokenTypes...)
}

// Checks to see if the look ahead is one of tokenTypes and if so returns true.
func (this *CompiledLexer) CheckLookAhead(tokenTypes ...int) bool {
	look := this.LookAhead()
	for _, val := range tokenTypes {
		if look.Type == val {
			return true
		}
	}
	return false
}

// Token creates a Token from the data at index.
// Returns TknINVALID if the index is invalid.
// Returns -1 if the index represent a newline.
func (this *CompiledLexer) Token(index int) *ncalex.Token {
	if index >= len(this.data.Code) || index < 0 {
		return &ncalex.Token { "INVALID", ncalex.TknINVALID, this.line, 0}
	}
	code := this.data.Code[index]
	
	if code == 0 {
		return &ncalex.Token { "(", ncalex.TknOpenParen, this.line, 0}
	}
	
	if code == 1 {
		return &ncalex.Token { ")", ncalex.TknCloseParen, this.line, 0}
	}
	
	if code == 2 {
		return &ncalex.Token { "[", ncalex.TknOpenSqBracket, this.line, 0}
	}
	
	if code == 3 {
		return &ncalex.Token { "]", ncalex.TknCloseSqBracket, this.line, 0}
	}
	
	if code == 4 { // This should NEVER, EVER happen
		panic("Trying to turn a control code into a token, Did someone forget to skip past newlines?")
	}
	
	code = code - 5
	if int(code) < len(this.data.StringTable) && code >= 0 {
		return &ncalex.Token { this.data.StringTable[code], ncalex.TknString, this.line, 0}
	}
	
	return &ncalex.Token { "INVALID", ncalex.TknINVALID, this.line, 0}
}
