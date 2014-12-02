// Encode and decode files in the Rubble compression formats and base 64.
package main

import (
	"fmt"
	"os"
	"io/ioutil"
	
	"compress/flate"
	
	"bytes"
	"strings"
	
	"encoding/base64"
)

func main() {
	
	if len(os.Args) != 2 {
		fmt.Println("Usage: encoder filename")
	}
	name := os.Args[1]
	content := ReadFile(name)
	
	switch {
	case strings.HasSuffix(name, ".zip"):
		WriteFile(name + ".b64", Split(Encode(content)))
		
	case strings.HasSuffix(name, ".zip.b64"):
		WriteFile(StripExt(name), Decode(content))
		
	case strings.HasSuffix(name, ".b64"):
		WriteFile(StripExt(name), Decompress(Decode(content)))
		
	default:
		WriteFile(name + ".b64", Split(Encode(Compress(content))))
	}
	
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

func WriteFile(name string, file []byte) {
	ioutil.WriteFile(name, file, 0600)
}

func ReadFile(name string) []byte {
	file, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return file
}

func StripExt(name string) string {
	i := len(name) - 1
	for i >= 0 {
		if name[i] == '.' {
			return name[:i]
		}
		i--
	}
	return name
}
