package main

import "fmt"
import "os"
import "io/ioutil"
import "path/filepath"
import "flag"
import "strings"

type RawFile struct {
	Name string
	Content string
	Namespace string
}

var SourceDirBase = "./source"
var DestDir = "./objects/"
var BaseConfig = "./config.txt"

var RawFiles = make([]*RawFile, 0, 10)
var CurFile *RawFile // used by CONFIG to get the current namespace

var AddonsDir = "./addons"
var Addons = make([]string, 0, 10)
var CurAddonDir = ""

var ParseStage = 0

var compat = flag.Bool("compat", false, "Blast compatability mode.")

func main() {
	flag.Parse()
	
	if *compat {
		SourceDirBase = "./base/objects"
		DestDir = "./raw/objects/"
		BaseConfig = "./base/config.txt"
	}
	
	// Read files
	fmt.Println("Loading Files...")
	
	// Load base files
	filepath.Walk(SourceDirBase, ListBaseFiles)
	ReadConfig(BaseConfig, "base")
	
	// Load files from "./addons" if compat
	if *compat {
		filepath.Walk(AddonsDir, ListAddonDirs)
		for i := range Addons {
			CurAddonDir = Addons[i]
			filepath.Walk(AddonsDir + "/" + CurAddonDir + "/objects", ListAddonFiles)
			ReadConfig(AddonsDir + "/" + CurAddonDir + "/config.txt", CurAddonDir)
		}
	}
	
	// Test lexer
	//for i := range RawFiles {
	//	lex := NewLexer(RawFiles[i].Content)
	//	for lex.Advance() {
	//		fmt.Println(lex.Current, ":", lex.Current.Lexeme)
	//	}
	//}
	
	// Load Templates
	SetupBuiltins()
	
	// preprocess
	fmt.Println("Preprocessing...")
	for i := range RawFiles {
		fmt.Println(RawFiles[i].Name)
		CurFile = RawFiles[i]
		RawFiles[i].Content = PreProcess(RawFiles[i].Content, RawFiles[i].Namespace, StripExt(filepath.Base(RawFiles[i].Name)))
	}
	
	// preparse
	fmt.Println("Preparsing...")
	for i := range RawFiles {
		fmt.Println(RawFiles[i].Name)
		CurFile = RawFiles[i]
		RawFiles[i].Content = PreParse(RawFiles[i].Content)
	}
	ParseStage++
	
	// parse
	fmt.Println("Parsing...")
	for i := range RawFiles {
		fmt.Println(RawFiles[i].Name)
		CurFile = RawFiles[i]
		RawFiles[i].Content = Parse(RawFiles[i].Content)
	}
	ParseStage++
	
	// postparse
	fmt.Println("Postparsing...")
	for i := range RawFiles {
		fmt.Println(RawFiles[i].Name)
		CurFile = RawFiles[i]
		RawFiles[i].Content = PostParse(RawFiles[i].Content)
	}
	ParseStage++
	
	// Write files out
	fmt.Println("Writing files...")
	if *compat {
		for i := range RawFiles {
			name := StripExt(filepath.Base(RawFiles[i].Name))
			file := []byte(name + "\n\n" + RawFiles[i].Content)
			ioutil.WriteFile(DestDir + name + "__" + RawFiles[i].Namespace + ".txt", file, 0600)
		}
	} else {
		for i := range RawFiles {
			file := []byte(StripExt(filepath.Base(RawFiles[i].Name)) + "\n\n" + RawFiles[i].Content)
			ioutil.WriteFile(DestDir + filepath.Base(RawFiles[i].Name), file, 0600)
		}
	}
	fmt.Println("Done.")
}

func ListBaseFiles(path string, info os.FileInfo, err error) error {
	if path == SourceDirBase {
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
	rawfile.Namespace = "base"
	
	RawFiles = append(RawFiles, rawfile)
	
	return nil
}

func ListAddonDirs(path string, info os.FileInfo, err error) error {
	if path == AddonsDir {
		return nil
	}
	if info.IsDir() {
		if !strings.HasPrefix(filepath.Base(path), "__") {
			Addons = append(Addons, filepath.Base(path))
		}
		
		return filepath.SkipDir
	}
	
	return nil
}

func ListAddonFiles(path string, info os.FileInfo, err error) error {
	if path == AddonsDir + "/" + CurAddonDir + "/objects" {
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
	rawfile.Namespace = CurAddonDir
	
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
