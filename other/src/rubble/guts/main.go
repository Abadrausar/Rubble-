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

// Rubble Guts package, needed for use with cover when used as a Rex test.
package guts

import "strings"
import "os"
import "runtime/pprof"
import "time"
import "sort"

func Main() {
	timeStart := time.Now()

	// init logging
	InitLogging()

	// Init crash handler
	defer func() {
		if NoRecover {
			return
		}

		if x := recover(); x != nil {
			LogPrintln("Error:", x)
			LogPrintln("  Near", LastLine)
			os.Exit(1)
		}
	}()

	LogPrintln("Rubble " + RubbleVersions[0])
	LogPrintln("After Blast comes Rubble.")
	LogPrintln("=============================================")

	ParseCommandLine()

	if Profile != "" {
		LogPrintln("Writing CPU profiling data to file:", Profile)
		f, err := os.Create(Profile)
		if err != nil {
			LogPrintln("  Profile Error:", err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// If we are in prep or install mode do not load any addon (or do any other addon type stuff)
	if PrepRegion == "" && Installer == "" {
		LogPrintln("=============================================")
		LogPrintln("Loading Addons...")

		Addons = LoadAddons()

		LogPrintln("=============================================")
		LogPrintln("Activating Addons...")

		addonNames := make([]string, 0, 10)
		if len(*AddonsList) == 0 {
			LogPrintln("No Addons Specified via Command Line.")
			LogPrintln("Using default list file.")

			// Lucky me, nonexistent addons are ignored :)
			// So if the addonlist.ini file is not there this entry is a NOP.
			AddonsList.Set("addons:dir:addonlist.ini")
		}

		for _, val := range *AddonsList {
			//file, err := ioutil.ReadFile(val)
			file, err := FS.ReadAll(val)
			if err == nil {
				LogPrintln("  Loading List File: " + val)

				lines := strings.Split(string(file), "\n")
				ParseINI(lines, func(key, value string) {
					value = strings.ToLower(value)
					if value == "true" {
						addonNames = append(addonNames, key)
					}
				})
			} else {
				addonNames = append(addonNames, strings.Split(val, ";")...)
			}
		}

		if len(addonNames) == 0 {
			LogPrintln("No Addons Marked Active")
			LogPrintln("Updating the Addon List File...")
			UpdateAddonList("addons:dir:addonlist.ini", Addons)
			LogPrintln("Rubble has no files to parse, aborting.")
			os.Exit(1)
		}

		addonNameLookupTbl := make(map[string]bool)
		for _, name := range addonNames {
			addonNameLookupTbl[name] = true
		}

		LogPrintln("Active Addons:")
		for i := range Addons {
			if addonNameLookupTbl[Addons[i].Name] {
				Addons[i].Active = true
				LogPrintln("  " + Addons[i].Name)
			}
		}

		LogPrintln("Updating the Addon List File...")
		UpdateAddonList("addons:dir:addonlist.ini", Addons)

		if ExitAfterUpdate {
			if Bench {
				timeEnd := time.Now()
				LogPrintln("Run time: ", timeEnd.Sub(timeStart))
			}
			LogPrintln("Done.")
			return
		}

		// Convert Addons to File List
		Files = NewFileList(Addons)

		if len(Files.Order) == 0 {
			LogPrintln("Files list is empty. (no active addons have parseable files)")
			LogPrintln("Rubble has no files to parse, aborting.")
			if Bench {
				LogPrintln("Run time: ", time.Since(timeStart))
			}
			os.Exit(1)
		}
	}

	LogPrintln("=============================================")
	LogPrintln("Initializing...")

	LogPrintln("Loading Config Variables...")
	if len(*ConfigList) != 0 {
		for _, val := range *ConfigList {
			var lines []string

			file, err := FS.ReadAll(val)
			if err == nil {
				LogPrintln("  Loading Config File: " + val)
				// Wow, I'm lazy...
				lines = strings.Split(string(file), "\n")
			} else {
				lines = strings.Split(val, ";")
			}

			ParseINI(lines, func(key, value string) {
				VariableData[key] = value
			})
		}
	} else {
		LogPrintln("  No value(s) specified.")
	}

	LogPrintln("Initializing Scripting Subsystem.")
	InitScripting()

	if Installer != "" {
		InstallModeRun(Installer)
		if Bench {
			LogPrintln("Run time: ", time.Since(timeStart))
		}
		LogPrintln("Done.")
		return
	}

	if PrepRegion != "" {
		PrepModeRun(PrepRegion)
		if Bench {
			LogPrintln("Run time: ", time.Since(timeStart))
		}
		LogPrintln("Done.")
		return
	}

	LogPrintln("Adding Builtins.")
	NewNativeTemplate("!TEMPLATE", tempTemplate)

	LogPrintln("=============================================")
	LogPrintln("Running Init Scripts...")
	sort.Strings(ForcedInitOrder)
	for _, i := range ForcedInitOrder {
		CurrentFile = i
		LogPrintln("  " + i)

		_, err := GlobalScriptState.CompileAndRun(string(ForcedInit[i]), i)
		if err != nil {
			panic("Script Error: " + err.Error())
		}
	}

	LogPrintln("=============================================")
	LogPrintln("Running Prescripts...")
	ParseStage = stgPreScripts
	for _, i := range Files.Order {
		if Files.Files[i].Skip || !Files.Files[i].PreScript {
			continue
		}

		CurrentFile = i
		LogPrintln("  " + Files.Files[i].Path)

		_, err := GlobalScriptState.CompileAndRun(string(Files.Files[i].Content), i)
		if err != nil {
			panic("Script Error: " + err.Error())
		}
	}

	LogPrintln("=============================================")
	LogPrintln("Preparsing...")
	ParseStage = stgPreParse
	for _, i := range Files.Order {
		if Files.Files[i].Skip || (!Files.Files[i].IsRaw() && !Files.Files[i].IsGraphic()) {
			continue
		}

		CurrentFile = i
		LogPrintln("  " + Files.Files[i].Path)
		Files.Files[i].Content = Parse(Files.Files[i].Content, stgUseCurrent, NewPosition(1, Files.Files[i].Path))
	}

	LogPrintln("=============================================")
	LogPrintln("Parsing...")
	ParseStage = stgParse
	for _, i := range Files.Order {
		if Files.Files[i].Skip || (!Files.Files[i].IsRaw() && !Files.Files[i].IsGraphic()) {
			continue
		}

		CurrentFile = i
		LogPrintln("  " + Files.Files[i].Path)
		Files.Files[i].Content = Parse(Files.Files[i].Content, stgUseCurrent, NewPosition(1, Files.Files[i].Path))
	}

	LogPrintln("=============================================")
	LogPrintln("Postparsing...")
	ParseStage = stgPostParse
	for _, i := range Files.Order {
		if Files.Files[i].Skip || (!Files.Files[i].IsRaw() && !Files.Files[i].IsGraphic()) {
			continue
		}

		CurrentFile = i
		LogPrintln("  " + Files.Files[i].Path)
		Files.Files[i].Content = Parse(Files.Files[i].Content, stgUseCurrent, NewPosition(1, Files.Files[i].Path))
	}

	LogPrintln("=============================================")
	LogPrintln("Expanding Variables...")
	ParseStage = stgGlobalExpand
	for _, i := range Files.Order {
		if Files.Files[i].Skip || (!Files.Files[i].IsRaw() && !Files.Files[i].IsGraphic()) {
			continue
		}

		CurrentFile = i
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

		CurrentFile = i
		LogPrintln("  " + Files.Files[i].Path)

		_, err := GlobalScriptState.CompileAndRun(string(Files.Files[i].Content), i)
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

		file := []byte(StripExt(i) + "\n\n# Automatically generated, do not edit!\n# Source: " +
			Files.Files[i].Source + "/" + Files.Files[i].Path + "\n\n" + string(Files.Files[i].Content))
		WriteFile("out:objects/"+i, file)
	}

	LogPrintln("Writing graphics files...")
	for _, i := range Files.Order {
		if Files.Files[i].Skip || Files.Files[i].NoWrite || !Files.Files[i].IsGraphic() {
			continue
		}

		LogPrintln("  " + Files.Files[i].Path)

		file := []byte(StripExt(i) + "\n\n# Automatically generated, do not edit!\n# Source: " +
			Files.Files[i].Source + "/" + Files.Files[i].Path + "\n\n" + string(Files.Files[i].Content))
		WriteFile("out:graphics/"+i, file)
	}

	LogPrintln("Writing prep files...")
	os.Mkdir(OutputDir+"/prep", 0600)
	for _, i := range Files.Order {
		if Files.Files[i].Skip || Files.Files[i].NoWrite || !Files.Files[i].PrepScript {
			continue
		}

		LogPrintln("  " + Files.Files[i].Path)

		WriteFile("out:prep/"+i, Files.Files[i].Content)
	}

	LogPrintln("Writing addon list to raw directory...")
	LogPrintln("  addonlist.ini")
	UpdateAddonList("out:addonlist.ini", Addons)

	LogPrintln("Writing config variables to raw directory...")
	LogPrintln("  genconfig.ini")
	DumpConfig("out:genconfig.ini")

	if Bench {
		LogPrintln("Run time: ", time.Since(timeStart))
	}
	LogPrintln("Done.")
}
