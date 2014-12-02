// Encode and decode files in the Rubble compression formats and base 64.
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
