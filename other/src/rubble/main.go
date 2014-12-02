/*
Copyright 2013 by Milo Christiansen

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

package main

import "fmt"
import "os"
import "io/ioutil"
import "path/filepath"
import "strings"
import "sort"
import "dctech/nca7"

func main() {
	
	// Init crash handler
	defer func() {
		if !RblRecover {
			return
		}
		
		if x := recover(); x != nil {
			fmt.Println("    Error:", x)
			fmt.Println("    Near line:", LastLine, "In last file")
		}
	}()
	
	fmt.Println("Rubble v2.0")
	fmt.Println("After Blast comes Rubble.")
	fmt.Println("Initalizing...")
	
	ParseCommandLine()
	
	InitNCA()

	SetupBuiltins()
	
	// Read files
	fmt.Println("=============================================")
	fmt.Println("Loading Files...")
	
	// Load base files
	CurWalkDir = BaseDir
	CurNamespace = "base"
	filepath.Walk(CurWalkDir, ListFiles)
	
	_, err := os.Lstat(BaseDir + "/config.ini")
	if err == nil {
		ReadConfig(BaseDir + "/config.ini")
	}

	// Load addon files
	if AddonsList != "" {
		AddonNames = filepath.SplitList(AddonsList)
	} else {
		filepath.Walk(AddonsDir, ListAddonNames)
	}
	for i := range AddonNames {
		CurWalkDir = AddonsDir + "/" + AddonNames[i]
		CurNamespace = AddonNames[i]
		filepath.Walk(CurWalkDir, ListFiles)
		
		_, err := os.Lstat(CurWalkDir + "/config.ini")
		if err == nil {
			ReadConfig(CurWalkDir + "/config.ini")
		}
	}
	
	// Add config overrides
	fmt.Println("=============================================")
	fmt.Println("Loading Config Overrides...")
	if ConfigList != "" {
		lines := filepath.SplitList(ConfigList)
		for i := range lines {
			if strings.HasPrefix(strings.TrimSpace(lines[i]), "#") {
				continue
			}
			if strings.TrimSpace(lines[i]) == "" {
				continue
			}
			
			parts := strings.SplitN(lines[i], "=", 2)
			if len(parts) != 2 {
				panic("Malformed config line.")
			}
			
			parts[0] = strings.TrimSpace(parts[0])
			VariableData[parts[0]] = strings.TrimSpace(parts[1])
		}
	}
	
	// This is needed to interleave addon and base files
	sort.Strings(RawOrder)
	sort.Strings(PreScriptOrder)
	sort.Strings(PostScriptOrder)
	
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
	
	// Prescripts
	fmt.Println("=============================================")
	fmt.Println("Running Prescripts...")
	for _, i := range PreScriptOrder {
		if PreScripts[i].Skip {
			continue
		}
		
		fmt.Println(PreScripts[i].Path)
		
		GlobalNCAState.Code.Add(PreScripts[i].Content)
		GlobalNCAState.Envs.Add(nca7.NewEnvironment())
		
		rtn, err := GlobalNCAState.Run()
		if err != nil {
			panic("Script Error: " + err.Error())
		}
		
		GlobalNCAState.Envs.Remove()
		
		if rtn == nil {
			PreScripts[i].Content = ""
			continue
		}
		PreScripts[i].Content = rtn.String()
	}
	
	// preparse
	fmt.Println("=============================================")
	fmt.Println("Preparsing...")
	for _, i := range RawOrder {
		if RawFiles[i].Skip {
			continue
		}
		
		fmt.Println(RawFiles[i].Path)
		CurFile = i
		RawFiles[i].Content = PreParse(RawFiles[i].Content)
	}
	ParseStage++
	
	// parse
	fmt.Println("=============================================")
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
	fmt.Println("=============================================")
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
	
	// Expand any remaining vars
	fmt.Println("=============================================")
	fmt.Println("Expanding Vars...")
	for _, i := range RawOrder {
		if RawFiles[i].Skip {
			continue
		}
		
		RawFiles[i].Content = ExpandVars(RawFiles[i].Content)
	}
	
	// Postscripts
	fmt.Println("=============================================")
	fmt.Println("Running Postscripts...")
	for _, i := range PostScriptOrder {
		if PostScripts[i].Skip {
			continue
		}
		
		fmt.Println(PostScripts[i].Path)
		
		GlobalNCAState.Code.Add(PostScripts[i].Content)
		GlobalNCAState.Envs.Add(nca7.NewEnvironment())
		
		rtn, err := GlobalNCAState.Run()
		if err != nil {
			panic("Script Error: " + err.Error())
		}
		
		GlobalNCAState.Envs.Remove()
		
		if rtn == nil {
			PostScripts[i].Content = ""
			continue
		}
		PostScripts[i].Content = rtn.String()
	}
	
	// Write files out
	fmt.Println("=============================================")
	fmt.Println("Writing files...")
	for _, i := range RawOrder {
		if RawFiles[i].Skip || RawFiles[i].NoWrite {
			continue
		}
		
		file := []byte(i + "\n\n" + RawFiles[i].Content)
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
	
	if strings.HasSuffix(filepath.Base(path), ".pre.nca") {
		// is pre script
		name := StripExt(StripExt(filepath.Base(path)))
		scriptfile := new(RawFile)
		scriptfile.Path = path
		scriptfile.Content = string(file)
		
		if _, ok := RawFiles[name]; !ok {
			PreScriptOrder = append(PreScriptOrder, name)
		}
		PreScripts[name] = scriptfile
		return nil
	}
	if strings.HasSuffix(filepath.Base(path), ".post.nca") {
		// is post script
		name := StripExt(StripExt(filepath.Base(path)))
		scriptfile := new(RawFile)
		scriptfile.Path = path
		scriptfile.Content = string(file)
		
		if _, ok := RawFiles[name]; !ok {
			PostScriptOrder = append(PostScriptOrder, name)
		}
		PostScripts[name] = scriptfile
		return nil
	}
	
	if strings.HasSuffix(filepath.Base(path), ".rbl") {
		// is rubble code (do not write out after parse)
		name := StripExt(filepath.Base(path))
		rawfile := new(RawFile)
		rawfile.Path = path
		rawfile.Content = string(file)
		rawfile.NoWrite = true
		
		if _, ok := RawFiles[name]; !ok {
			RawOrder = append(RawOrder, name)
		}
		RawFiles[name] = rawfile
		return nil
	}
	
	if strings.HasSuffix(filepath.Base(path), ".txt") {
		// is raw file
		name := StripExt(filepath.Base(path))
		rawfile := new(RawFile)
		rawfile.Path = path
		rawfile.Content = string(file)
		
		if _, ok := RawFiles[name]; !ok {
			RawOrder = append(RawOrder, name)
		}
		RawFiles[name] = rawfile
		return nil
	}
	
	fmt.Println("    Not a parsable file.")
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
