package main

import "fmt"
import "os"
import "io/ioutil"
import "path/filepath"
import "strings"
import "sort"

func main() {
	fmt.Println("Rubble v1.0")
	fmt.Println("After Blast comes Rubble.")
	fmt.Println("Initalizing...")
	
	ParseCommandLine()
	
	InitNCA()

	SetupBuiltins()
	
	// Read files
	fmt.Println("Loading Files...")
	
	// Load base files
	CurWalkDir = BaseDir
	CurNamespace = "base"
	filepath.Walk(CurWalkDir, ListFiles)
	
	_, err := os.Lstat(ConfigDir + "/base.ini")
	if err == nil {
		ReadConfig(ConfigDir + "/base.ini")
	}

	// Load addon files
	filepath.Walk(AddonsDir, ListAddonNames)
	for i := range AddonNames {
		CurWalkDir = AddonsDir + "/" + AddonNames[i]
		CurNamespace = AddonNames[i]
		filepath.Walk(CurWalkDir, ListFiles)
		
		_, err := os.Lstat(ConfigDir + "/" + CurNamespace + ".ini")
		if err == nil {
			ReadConfig(ConfigDir + "/" + CurNamespace + ".ini")
		}
	}
	
	// This is needed to interleave addon and base files
	sort.Strings(RawOrder)
	
	// Test lexer
	if LexTest {
		for _, i := range RawOrder {
			lex := NewLexer(RawFiles[i].Content)
			for lex.Advance() {
				fmt.Println(lex.Current, ":", lex.Current.Lexeme)
			}
		}
		return
	}
	
	// preparse
	fmt.Println("Preparsing...")
	for _, i := range RawOrder {
		if RawFiles[i].Skip {
			fmt.Println("Skipping File:", RawFiles[i].Path)
			continue
		}
		
		fmt.Println(RawFiles[i].Path)
		CurFile = i
		RawFiles[i].Content = PreParse(RawFiles[i].Content)
	}
	ParseStage++
	
	// parse
	fmt.Println("Parsing...")
	for _, i := range RawOrder {
		if RawFiles[i].Skip {
			continue
		}
		
		fmt.Println(RawFiles[i].Path)
		CurFile = i
		RawFiles[i].Content = Parse(RawFiles[i].Content)
	}
	ParseStage++
	
	// postparse
	fmt.Println("Postparsing...")
	for _, i := range RawOrder {
		if RawFiles[i].Skip {
			continue
		}
		
		fmt.Println(RawFiles[i].Path)
		CurFile = i
		RawFiles[i].Content = PostParse(RawFiles[i].Content)
	}
	ParseStage++
	
	// Write files out
	fmt.Println("Writing files...")
	for _, i := range RawOrder {
		if RawFiles[i].Skip {
			continue
		}
		
		file := []byte(i + "\n\n" + ExpandVars(RawFiles[i].Content))
		ioutil.WriteFile(OutputDir + "/" + i + ".txt", file, 0600)
	}
	fmt.Println("Done.")
}

func ListFiles(path string, info os.FileInfo, err error) error {
	if path == CurWalkDir {
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
	
	name := StripExt(filepath.Base(path))
	rawfile := new(RawFile)
	rawfile.Path = path
	rawfile.Content = string(file)
	rawfile.Namespace = CurNamespace
	
	if _, ok := RawFiles[name]; !ok {
		RawOrder = append(RawOrder, name)
	}
	
	RawFiles[name] = rawfile
	
	return nil
}

func ListAddonNames(path string, info os.FileInfo, err error) error {
	if path == AddonsDir {
		return nil
	}
	if info.IsDir() {
		if !strings.HasPrefix(filepath.Base(path), "__") {
			AddonNames = append(AddonNames, filepath.Base(path))
		}
		
		return filepath.SkipDir
	}
	
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
