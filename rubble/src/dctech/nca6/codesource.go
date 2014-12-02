package nca6

//import "fmt"
import "dctech/ncalex"

// A CodeSource represents a token stream, ether directly from the lexer or from a compiled script.
type CodeSource interface {
	Line() int
	CurrentTkn() *ncalex.Token
	LookAhead() *ncalex.Token
	Advance()
	GetToken(...int)
	CheckLookAhead(...int) bool
}
