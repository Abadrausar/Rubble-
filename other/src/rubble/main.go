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

import "io/ioutil"
import "path/filepath"
import "strings"
import "dctech/raptor"
import "os"

func main() {
	// init logging
	InitLogging()
	
	// Init crash handler
	defer func() {
		if NoRecover {
			return
		}

		if x := recover(); x != nil {
			LogPrintln("Error:", x)
			LogPrintln("  Near line:", LastLine, "In last file.")
			os.Exit(1)
		}
	}()

	LogPrintln("Rubble v3.6")
	LogPrintln("After Blast comes Rubble.")
	LogPrintln("=============================================")

	ParseCommandLine()

	LogPrintln("=============================================")
	LogPrintln("Loading Addons...")

	Addons = LoadAddons(AddonsDir)

	LogPrintln("=============================================")
	LogPrintln("Activating Addons...")

	addonNames := make([]string, 0, 10)
	if AddonsList != "" {
		LogPrintln("Addons Specified via Command Line.")
		addonNames = filepath.SplitList(AddonsList)
	} else {
		LogPrintln("Addons Not Specified via Command line.")
		LogPrintln("Loading " + AddonsDir + "/addonlist.ini...")
		file, err := ioutil.ReadFile(AddonsDir + "/addonlist.ini")
		if err != nil && !ShellMode {
			LogPrintln("Read failed (this is bad if unexpected)\n  Error:", err)
			UpdateAddonList(AddonsDir, Addons)
			LogPrintln("Rubble has no files to parse, aborting.")
			return
		} else if ShellMode {
			LogPrintln("Read failed (this is bad if unexpected)\n  Error:", err)
			UpdateAddonList(AddonsDir, Addons)
			LogPrintln("Rubble has no files to parse but shell mode is active, continuing.")
		}

		lines := strings.Split(string(file), "\n")
		for i := range lines {
			if strings.HasPrefix(strings.TrimSpace(lines[i]), "#") {
				continue
			}
			if strings.HasPrefix(strings.TrimSpace(lines[i]), "[") {
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
			parts[1] = strings.ToLower(strings.TrimSpace(parts[1]))
			if parts[1] == "true" {
				addonNames = append(addonNames, parts[0])
			}
		}
	}

	if len(addonNames) == 0 && !ShellMode {
		LogPrintln("No Addons Marked Active")
		LogPrintln("Rubble has no files to parse, aborting.")
		return
	}
	
	addonNameLookupTbl := make(map[string]bool)
	for _, name := range addonNames {
		addonNameLookupTbl[name] = true
	}

	for i := range Addons {
		if addonNameLookupTbl[Addons[i].Name] {
			Addons[i].Active = true
			LogPrintln("  " + Addons[i].Name)
		}
	}

	UpdateAddonList(AddonsDir, Addons)

	if ExitAfterUpdate {
		LogPrintln("Done.")
		return
	}

	// Convert Addons to File List
	Files = NewFileList(Addons)

	if len(Files.Order) == 0 && !ShellMode {
		LogPrintln("Files list is empty. (no active addons have parseable files)")
		LogPrintln("Rubble has no files to parse, aborting.")
		return
	}

	// Test lexer
	if LexTest && !ShellMode {
		for _, i := range Files.Order {
			lex := NewLexer(Files.Files[i].Content)
			for {
				lex.Advance()
				if lex.Current.Type == tknINVALID {
					break
				}
				LogPrintln(lex.Current, ":", lex.Current.Lexeme)
			}
		}
		return
	}

	LogPrintln("=============================================")
	LogPrintln("Initializing...")

	LogPrintln("Loading Variables from Command Line.")
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

	LogPrintln("Initializing Scripting Subsystem.")
	InitScripting()
	
	if !ShellMode || RunForcedInit {
		LogPrintln("Running Forced Init Scripts.")
		for i := range ForcedInit {
			script := raptor.NewScript()
			err := raptor.LoadFile(ForcedInit[i], script)
			if err != nil {
				panic("Script Error: " + err.Error())
			}
	
			_, err = GlobalRaptorState.Run(script)
			if err != nil {
				panic("Script Error: " + err.Error())
			}
		}
	} else {
		LogPrintln("Skipping Forced Init Scripts.")
	}
	
	LogPrintln("Adding Builtins.")
	SetupBuiltins()
	
	// Embedded Raptor Shell
	if ShellMode {
		ShellModeRun()
		os.Exit(0)
	}
	
	LogPrintln("=============================================")
	LogPrintln("Running Prescripts...")
	ParseStage = stgPreScripts
	for _, i := range Files.Order {
		if Files.Files[i].Skip || !Files.Files[i].PreScript {
			continue
		}

		LogPrintln("  " + Files.Files[i].Path)

		script := raptor.NewScript()
		err := raptor.LoadFile(Files.Files[i].Content, script)
		if err != nil {
			panic("Script Error: " + err.Error())
		}

		_, err = GlobalRaptorState.Run(script)
		if err != nil {
			panic("Script Error: " + err.Error())
		}
	}

	LogPrintln("=============================================")
	LogPrintln("Preparsing...")
	ParseStage = stgPreParse
	for _, i := range Files.Order {
		if Files.Files[i].Skip || !Files.Files[i].IsRaw() {
			continue
		}

		LogPrintln("  " + Files.Files[i].Path)
		Files.Files[i].Content = Parse(Files.Files[i].Content, stgUseCurrent)
	}

	LogPrintln("=============================================")
	LogPrintln("Parsing...")
	ParseStage = stgParse
	for _, i := range Files.Order {
		if Files.Files[i].Skip || !Files.Files[i].IsRaw() {
			continue
		}

		LogPrintln("  " + Files.Files[i].Path)
		Files.Files[i].Content = Parse(Files.Files[i].Content, stgUseCurrent)
	}

	LogPrintln("=============================================")
	LogPrintln("Postparsing...")
	ParseStage = stgPostParse
	for _, i := range Files.Order {
		if Files.Files[i].Skip || !Files.Files[i].IsRaw() {
			continue
		}

		LogPrintln("  " + Files.Files[i].Path)
		Files.Files[i].Content = Parse(Files.Files[i].Content, stgUseCurrent)
	}

	LogPrintln("=============================================")
	LogPrintln("Expanding Variables...")
	ParseStage = stgGlobalExpand
	for _, i := range Files.Order {
		if Files.Files[i].Skip || !Files.Files[i].IsRaw() {
			continue
		}

		LogPrintln("  " + Files.Files[i].Path)
		Files.Files[i].Content = []byte(ExpandVars(string(Files.Files[i].Content)))
	}

	LogPrintln("=============================================")
	LogPrintln("Running Postscripts...")
	ParseStage = stgPostScripts
	for _, i := range Files.Order {
		if Files.Files[i].Skip || !Files.Files[i].PostScript {
			continue
		}

		LogPrintln("  " + Files.Files[i].Path)

		script := raptor.NewScript()
		err := raptor.LoadFile(Files.Files[i].Content, script)
		if err != nil {
			panic("Script Error: " + err.Error())
		}
		
		_, err = GlobalRaptorState.Run(script)
		if err != nil {
			panic("Script Error: " + err.Error())
		}
	}

	LogPrintln("=============================================")
	LogPrintln("Writing files...")
	for _, i := range Files.Order {
		if Files.Files[i].Skip || Files.Files[i].NoWrite || !Files.Files[i].IsRaw() {
			continue
		}

		LogPrintln("  " + Files.Files[i].Path)

		// HACK: Redo this
		file := []byte(StripExt(i) + "\n\n" + string(Files.Files[i].Content))
		err := ioutil.WriteFile(OutputDir+"/"+i, file, 0600)
		if err != nil {
			panic("Write Error: " + err.Error())
		}
	}
	LogPrintln("Done.")
}
