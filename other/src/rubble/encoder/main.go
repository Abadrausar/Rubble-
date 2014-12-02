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

// Rubble File Encoder.
// Encode and decode files in base 64 encoded DEFLATE compression.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"rubble/rblutil"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: encoder filename")
	}
	name := os.Args[1]
	content := ReadFile(name)

	switch {
	case strings.HasSuffix(name, ".zip"):
		WriteFile(name+".b64", rblutil.Split(rblutil.Encode(content)))

	case strings.HasSuffix(name, ".zip.b64"):
		WriteFile(rblutil.StripExt(name), rblutil.Decode(rblutil.StripWS(content)))

	case strings.HasSuffix(name, ".b64"):
		WriteFile(rblutil.StripExt(name), rblutil.Decompress(rblutil.Decode(rblutil.StripWS(content))))

	default:
		WriteFile(name+".b64", rblutil.Split(rblutil.Encode(rblutil.Compress(content))))
	}
}

func WriteFile(name string, file []byte) {
	ioutil.WriteFile(name, file, 04755)
}

func ReadFile(name string) []byte {
	file, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return file
}
