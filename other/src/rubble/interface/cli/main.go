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

// Rubble CLI Interface
package main

import "os"
import "fmt"
import "runtime/pprof"
import "time"
import "io/ioutil"
import "strconv"
import "runtime"

import "flag"

import "rubble"
import "rubble/rblutil"

var DFDir string
var OutputDir string
var AddonsDir *rblutil.ArgList

var AddonsList *rblutil.ArgList
var ConfigList *rblutil.ArgList

var ZapAddons bool
var ZapConfig bool

var Bench bool

var NoRecover bool
var ExitAfterUpdate bool

var PrepRegion string

var Threaded bool

var Profile string

func main() {
	timeStart := time.Now()
	err, log := rblutil.NewLogger()
	if err != nil {
		fmt.Println("Fatal Error:", err)
		os.Exit(1)
	}

	flags := flag.NewFlagSet("rubble", flag.ExitOnError)
	flags.SetOutput(log)

	log.Header(rubble.Versions[0])

	// Used for the -zapxxx options.
	addonCount := 0
	configCount := 0

	// Setting hardcoded defaults
	DFDir = ".."
	OutputDir = "df:raw"
	AddonsDir = new(rblutil.ArgList)

	AddonsList = new(rblutil.ArgList)
	ConfigList = new(rblutil.ArgList)

	PrepRegion = ""

	// Load defaults from config if present
	log.Println("Attempting to Read Config File: ./rubble.ini")
	file, err := ioutil.ReadFile("./rubble.ini")
	if err == nil {
		log.Println("  Read OK, loading options from file.")
		rblutil.ParseINI(string(file), "\n", func(key, value string) {
			switch key {
			case "dfdir":
				DFDir = value
			case "outputdir":
				OutputDir = value
			case "addonsdir":
				AddonsDir.Set(value)
			case "addons":
				addonCount++
				AddonsList.Set(value)
			case "config":
				configCount++
				ConfigList.Set(value)
			case "zapaddons":
				ZapAddons, _ = strconv.ParseBool(value)
			case "zapconfig":
				ZapConfig, _ = strconv.ParseBool(value)
			case "bench":
				Bench, _ = strconv.ParseBool(value)
			case "addonlist":
				ExitAfterUpdate, _ = strconv.ParseBool(value)
			case "norecover":
				NoRecover, _ = strconv.ParseBool(value)
			case "prep":
				PrepRegion = value
			case "threads":
				Threaded, _ = strconv.ParseBool(value)
			case "profile":
				Profile = value
			default:
				log.Println("  Invalid config option:", key, ", skipping.")
			}
		})
	} else {
		log.Println("  Read failed (this is most likely ok)\n  Error:", err)
		log.Println("    Using hardcoded defaults.")
	}

	flags.StringVar(&DFDir, "dfdir", DFDir, "The path to the DF directory. May be an AXIS path (only the 'rubble' location ID works).")
	flags.StringVar(&OutputDir, "outputdir", OutputDir, "Where should Rubble write the generated raw files? May be an AXIS path (only the 'rubble' and 'df' location IDs work).")
	flags.Var(AddonsDir, "addonsdir", "Rubble addons directory. May be an AXIS path (only the 'rubble', 'df', and 'out' location IDs work).")

	flags.Var(AddonsList, "addons", "List of addons to load. This is optional. If the value is a file path then the file is read as an ini file containing config variables. May be specified more than once.")
	flags.Var(ConfigList, "config", "List of config variables. This is optional. If the value is a file path then the file is read as an ini file containing config variables. May be specified more than once.")

	flags.BoolVar(&ZapAddons, "zapaddons", ZapAddons, "Ignore any -addons flags loaded from rubble.ini.")
	flags.BoolVar(&ZapConfig, "zapconfig", ZapConfig, "Ignore any -config flags loaded from rubble.ini.")

	flags.BoolVar(&Bench, "bench", Bench, "Display the elapsed time before exiting. Only works if no errors were encountered.")

	flags.BoolVar(&ExitAfterUpdate, "addonlist", ExitAfterUpdate, "Update the addon list and exit.")

	flags.BoolVar(&NoRecover, "norecover", NoRecover, "Should Rubble not recover errors? Useful for debugging.")

	flags.StringVar(&PrepRegion, "prep", PrepRegion, "Name of a world to prepare DF for loading (or \"raw\" for the base raw folder). Use this to make sure tilesets, init changes, and DFHack scripts match the world's requirements.")

	flags.BoolVar(&Threaded, "threads", Threaded, "Allows Rubble to use more than one processor core, not useful except for running threaded scripts.")

	flags.StringVar(&Profile, "profile", "", "Output CPU profile information to specified file.")

	flags.Parse(os.Args[1:])

	if AddonsDir.Empty() {
		AddonsDir.Set("rubble:addons")
	}

	if ZapAddons {
		if addonCount > 0 && addonCount < len(*AddonsList) {
			*AddonsList = (*AddonsList)[addonCount:]
		}
	}

	if ZapConfig {
		if configCount > 0 && configCount < len(*ConfigList) {
			*ConfigList = (*ConfigList)[configCount:]
		}
	}

	if Threaded {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	if Profile != "" {
		log.Println("Writing CPU profiling data to file:", Profile)
		f, err := os.Create(Profile)
		if err != nil {
			log.Println("  Profile Error:", err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	err, state := rubble.NewState(DFDir, OutputDir, *AddonsDir, log)
	if err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}
	state.NoRecover = NoRecover
	state.ScriptState.NoRecover = NoRecover

	if ExitAfterUpdate {
		err := state.Load(*AddonsList, *ConfigList)
		if Bench {
			log.Println("Run time: ", time.Since(timeStart))
		}
		if err != nil {
			if _, ok := err.(rubble.Abort); ok {
				log.Println("Abort:", err)
			} else {
				log.Println("Error:", err)
			}
			os.Exit(1)
		}
		log.Println("Done.")
		return
	}

	if PrepRegion != "" {
		err := state.PrepModeRun(PrepRegion)
		if Bench {
			log.Println("Run time: ", time.Since(timeStart))
		}
		if err != nil {
			if _, ok := err.(rubble.Abort); ok {
				log.Println("Abort:", err)
			} else {
				log.Println("Error:", err)
			}
			os.Exit(1)
		}
		log.Println("Done.")
		return
	}

	err = state.Run(*AddonsList, *ConfigList)
	if Bench {
		log.Println("Run time: ", time.Since(timeStart))
	}
	if err != nil {
		if _, ok := err.(rubble.Abort); ok {
			log.Println("Abort:", err)
		} else {
			log.Println("Error:", err)
		}
		os.Exit(1)
	}
	log.Println("Done.")
}
