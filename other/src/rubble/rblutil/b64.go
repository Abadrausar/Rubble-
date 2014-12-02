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

package rblutil

import (
	"compress/flate"
	"io/ioutil"
	"bytes"
	"encoding/base64"
)

func StripWS(content []byte) []byte {
	return bytes.Map(func(in rune) rune {
		if in == ' ' || in == '\t' || in == '\n' || in == '\r' {
			return -1
		}
		return in
	}, content)
}

func Split(content []byte) []byte {
	out := make([]byte, 0, len(content) + len(content) / 80)
	
	x := 0
	for i := range content {
		if x >= 80 {
			out = append(out, '\n')
			x = 0
		}
		out = append(out, content[i])
		x++
	}
	
	return out
}

func Encode(content []byte) []byte {
	b := new(bytes.Buffer)
	bc := base64.NewEncoder(base64.StdEncoding, b)
	_, err := bc.Write(content)
	if err != nil {
		panic(err)
	}
	bc.Close()
	return b.Bytes()
}

func Decode(content []byte) []byte {
	a := bytes.NewReader(content)
	ac := base64.NewDecoder(base64.StdEncoding, a)
	data, err := ioutil.ReadAll(ac)
	if err != nil {
		panic(err)
	}
	return data
}

func Compress(content []byte) []byte {
	b := new(bytes.Buffer)
	bc, _ := flate.NewWriter(b, 9)
	_, err := bc.Write(content)
	if err != nil {
		panic(err)
	}
	bc.Close()
	return b.Bytes()
}

func Decompress(content []byte) []byte {
	a := bytes.NewReader(content)
	ac := flate.NewReader(a)
	data, err := ioutil.ReadAll(ac)
	if err != nil {
		panic(err)
	}
	ac.Close()
	return data
}
