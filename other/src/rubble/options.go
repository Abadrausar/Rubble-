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

import "flag"
import "strings"
import "io/ioutil"
import "os"

// Getting usage info and flag parse errors into the log file is a pain in the a**
var Flags = flag.NewFlagSet("rubble", flag.ExitOnError)

var DFDir string
var OutputDir string
var AddonsDir string

var AddonsList string
var ConfigList string

var LexTest bool
var NoRecover bool
var ExitAfterUpdate bool

var ShellMode bool
var ScriptPath string
var Compile string
var BinVersion int
var Validate bool
var NoExit bool
var NoPredefs bool
var RunForcedInit bool

func ParseCommandLine() {
	Flags.SetOutput(logFile)

	// Setting default defaults
	DFDir = ".."
	OutputDir = "../raw/objects"
	AddonsDir = "./addons"

	// Load defaults from config if present
	LogPrintln("Attempting to Read Config File: ./rubble.ini")
	file, err := ioutil.ReadFile("./rubble.ini")
	if err == nil {
		lines := strings.Split(string(file), "\n")
		for i := range lines {
			if strings.HasPrefix(strings.TrimSpace(lines[i]), "[") {
				continue
			}
			if strings.HasPrefix(strings.TrimSpace(lines[i]), "#") {
				continue
			}
			if strings.TrimSpace(lines[i]) == "" {
				continue
			}

			parts := strings.SplitN(lines[i], "=", 2)
			if len(parts) != 2 {
				LogPrintln("  Malformed config line found, skipping.")
				continue
			}

			parts[0] = strings.TrimSpace(parts[0])
			parts[1] = strings.TrimSpace(parts[1])

			switch parts[0] {
			case "dfdir":
				DFDir = parts[1]
			case "outputdir":
				OutputDir = parts[1]
			case "addonsdir":
				AddonsDir = parts[1]
			default:
				LogPrintln("  Invalid config option:", parts[0], ", skipping.")
			}
		}
	} else {
		LogPrintln("Read failed (this is most likely ok)\n  Error:", err)
		LogPrintln("  Using hardcoded defaults.")
	}

	Flags.StringVar(&DFDir, "dfdir", DFDir, "The path to the DF directory. May be relative.")
	Flags.StringVar(&OutputDir, "outputdir", OutputDir, "Where should Rubble write the generated raw files?")
	Flags.StringVar(&AddonsDir, "addonsdir", AddonsDir, "Rubble addons directory.")
	
	Flags.StringVar(&AddonsList, "addons", "", "List of addons to load. This is optional.")
	Flags.StringVar(&ConfigList, "config", "", "List of config overrides. This is optional.")

	Flags.BoolVar(&ExitAfterUpdate, "addonlist", false, "Update the addon list and exit.")

	Flags.BoolVar(&LexTest, "lextest", false, "Run a Rubble lexer test. No files will be written.")
	Flags.BoolVar(&NoRecover, "norecover", false, "Should Rubble not recover errors? Useful for debugging.")
	
	Flags.BoolVar(&ShellMode, "shell", false, "Enter Shell Mode. Shell Mode is a full-up version of the Raptor Shell.")
	Flags.StringVar(&ScriptPath, "script", "", "Shell Mode: Path to the input script, if any. Changes to batch mode. Needed for -compile")
	Flags.StringVar(&Compile, "compile", "", "Shell Mode: Path to the output file. Changes to compile mode. Needs -script to be set.")
	Flags.IntVar(&BinVersion, "binversion", 4, "Shell Mode: Force a specific binary version. Fallback rules still apply.")
	Flags.BoolVar(&Validate, "validate", false, "Shell Mode: Run script through the syntax checker and exit. Use with -script.")
	Flags.BoolVar(&NoExit, "noexit", false, "Shell Mode: If set changes from batch mode to interactive mode. Use with -script.")
	Flags.BoolVar(&NoPredefs, "nopredefs", false, "Shell Mode: If set disables the shell predefs.")
	Flags.BoolVar(&RunForcedInit, "forcedinit", false, "Shell Mode: Run the Rubble forced init scripts before entering shell mode.")
	
	Flags.Parse(os.Args[1:])
}
