package main

import "fmt"
import "os"
//import "io"
import "io/ioutil"
import "path/filepath"

type RawFile struct {
	Name string
	Content string
}

var SourceDir = "./source"
var DestDir = "./objects/"
var RawFiles = make([]*RawFile, 0, 10)
var ParseStage = 0

func main() {
	// Read files
	fmt.Println("Loading Files...")
	filepath.Walk(SourceDir, List)
	
	// test
	//for i := range RawFiles {
	//	lex := NewLexer(RawFiles[i].Content)
	//	for lex.Advance() {
	//		fmt.Println(lex.Current, ":", lex.Current.Lexeme)
	//	}
	//}
	
	SetupBuiltins()
	
	// preparse
	fmt.Println("Preparsing...")
	for i := range RawFiles {
		fmt.Println(RawFiles[i].Name)
		RawFiles[i].Content = PreParse(RawFiles[i].Content)
	}
	ParseStage++
	
	// parse
	fmt.Println("Parsing...")
	for i := range RawFiles {
		fmt.Println(RawFiles[i].Name)
		RawFiles[i].Content = Parse(RawFiles[i].Content)
	}
	ParseStage++
	
	// postparse
	fmt.Println("Postparsing...")
	for i := range RawFiles {
		fmt.Println(RawFiles[i].Name)
		RawFiles[i].Content = PostParse(RawFiles[i].Content)
	}
	ParseStage++
	
	// Write files out
	fmt.Println("Writing files...")
	for i := range RawFiles {
		file := []byte(StripExt(filepath.Base(RawFiles[i].Name)) + "\n\n" + RawFiles[i].Content)
		ioutil.WriteFile(DestDir + filepath.Base(RawFiles[i].Name), file, 0600)
	}
	fmt.Println("Done.")
}

func List(path string, info os.FileInfo, err error) error {
	if path == SourceDir {
		return nil
	}
	if info.IsDir() {
		return filepath.SkipDir
	}
	
	fmt.Println(path)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	
	rawfile := new(RawFile)
	rawfile.Name = path
	rawfile.Content = string(file)
	
	RawFiles = append(RawFiles, rawfile)
	
	return nil
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
