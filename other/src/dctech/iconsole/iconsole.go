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
	char []byte
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

func (this *Info) Run() ([]byte, error) {
	escape := false
	
	out := make([]byte, 0, 100)
	
	fmt.Fprint(this.Out, this.Prompt)
	for {
		_, err := os.Stdin.Read(this.char)
		if err != nil {
			return out, err
		}
		
		if this.char[0] == byte('\r') {
			continue
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
