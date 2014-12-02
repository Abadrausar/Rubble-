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

// Interactive Console helper library.
package iconsole

import "fmt"
import "os"
import "io"

// Info contains information about a specific console.
type Info struct {
	Prompt string
	In     io.Reader
	Out    io.Writer
	Escape byte
	char   []byte
}

// New creates a new interactive console information object.
// Default Values:
//	Prompt		>>>
//	In			stdin
//	Out			stdout
//	Escape		\
func New() *Info {
	this := new(Info)
	this.Prompt = ">>>"
	this.In = os.Stdin
	this.Out = os.Stdout
	this.Escape = '\\'
	this.char = make([]byte, 1, 1)
	return this
}

// Run prints a prompt and then reads input one byte at a time until EOF or newline.
// If EOF is is found an error is returned, in any case the line read up to the end condition
// is returned, if Run exits with an error there may be data returned as well.
// 
// To workaround a go bug Run will simulate EOF if it finds char 26 (the Windows/DOS EOF char).
// 
// I have no idea how well (or if) this works with Unicode.
// 
// If newline is prefixed by the defined escape char then it is added to the output buffer, and
// Run continues to look for input, in all other cases the escape char is not given special treatment
// (eg you do not need to double it up)
//
// Ways to generate EOF (on windows):
//	Ctrl+Z			Sends char 26
//	Ctrl+Break		Seems to work like a combo of Alt+26 followed by Enter
//	Ctrl+C			Same as Ctrl+Break?
//	Alt+26			Generates an EOF condition
//	Alt+026			Same as Ctrl+Z
//	Alt+27			Same as Alt+26?
// I am not sure, but there may be more...
func (this *Info) Run() ([]byte, error) {
	escape := false

	out := make([]byte, 0, 100)
	eatEOF := false // state var for eating the newline after an EOF char.

	fmt.Fprint(this.Out, this.Prompt)
	for {
		_, err := os.Stdin.Read(this.char)
		if err != nil {
			return out, err
		}

		// The DOS EOF char.
		// This causes Ctrl+Z and Alt+026 to work for later go versions (1.2).
		if this.char[0] == 26 {
			eatEOF = true
			continue
		}
		
		if this.char[0] == byte('\r') {
			continue
		}
		
		if eatEOF && (this.char[0] != byte('\n')) {
			eatEOF = false
		}

		if eatEOF && (this.char[0] == byte('\n')) {
			return out, io.EOF
		}

		if this.char[0] == this.Escape && !escape {
			escape = true
			continue
		}

		if this.char[0] == byte('\n') && !escape {
			return out, nil
		}

		if this.char[0] != byte('\n') && escape {
			out = append(out, this.Escape)
		}

		if this.char[0] == byte('\n') && escape {
			fmt.Fprint(this.Out, this.Prompt)
		}

		escape = false
		out = append(out, this.char...)
	}

	panic("UNREACHABLE")
}
