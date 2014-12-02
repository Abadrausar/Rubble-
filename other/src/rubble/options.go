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
import "fmt"
import "strings"
import "io/ioutil"

var DFDir string
var OutputDir string
var BaseDir string
var AddonsDir string

var AddonsList string
var ConfigList string

var LexTest bool
var NCARecover bool
var RblRecover bool

func ParseCommandLine() {
	// Setting default defaults
	DFDir = ".."
	OutputDir = "../raw/objects"
	BaseDir = "./base"
	AddonsDir = "./addons"
	
	// Load defaults from config if present
	fmt.Println("Attempting Read of Config File: ./rubble.cfg")
	file, err := ioutil.ReadFile("./rubble.cfg")
	if err == nil {
		lines := strings.Split(string(file), "\n")
		for i := range lines {
			if strings.HasPrefix(strings.TrimSpace(lines[i]), "#") {
				continue
			}
			if strings.TrimSpace(lines[i]) == "" {
				continue
			}
			
			parts := strings.SplitN(lines[i], "=", 2)
			if len(parts) != 2 {
				fmt.Println("    Malformed config line found, skipping.")
				continue
			}
			
			parts[0] = strings.TrimSpace(parts[0])
			
			switch parts[0] {
			case "dfdir":
				DFDir = parts[1]
			case "outputdir":
				OutputDir = parts[1]
			case "basedir":
				BaseDir = parts[1]
			case "addonsdir":
				AddonsDir = parts[1]
			default:
				fmt.Println("    Invalid config option:", parts[0], ", skipping.")
			}
		}
	} else {
		fmt.Println("Read failed (this is most likely ok)\n    Error:", err)
		fmt.Println("    Using hardcoded defaults.")
	}
	
	flag.StringVar(&DFDir ,"dfdir", DFDir, "The path to the DF directory. May be relative.")
	flag.StringVar(&OutputDir ,"outputdir", OutputDir, "Where should Rubble write the generated raw files?")
	flag.StringVar(&BaseDir ,"basedir", BaseDir, "Rubble base directory.")
	flag.StringVar(&AddonsDir ,"addonsdir", AddonsDir, "Rubble addons directory.)")
	
	flag.StringVar(&AddonsList, "addons", "", "List of addons to load. This is optional.")
	flag.StringVar(&ConfigList, "config", "", "List of config overrides. This is optional.")
	
	flag.BoolVar(&LexTest, "lextest", false, "Run a Rubble lexer test. No files will be written.")
	flag.BoolVar(&NCARecover, "ncarecover", true, "Should NCA recover errors? Useful for debugging.")
	flag.BoolVar(&RblRecover, "rblrecover", true, "Should Rubble recover errors? Useful for debugging.")
	
	flag.Parse()
}
