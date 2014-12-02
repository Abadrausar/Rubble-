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

package guts

import "flag"
import "strings"
import "strconv"
import "io/ioutil"
import "os"
import "runtime"
import "dctech/axis"
import "dctech/axis/axiszip"
import "rubble/rblutil"

// Getting usage info and flag parse errors into the log file is a pain in the a**
var Flags *flag.FlagSet

// Only used in very rare special cases related to the "go test" command.
var ContinueOnBadFlag bool

var DFDir string
var OutputDir string
var AddonsDir string

var Installer string

var AddonsList *ArgList
var ConfigList *ArgList

var ZapAddons bool
var ZapConfig bool

var Bench bool

var NoRecover bool
var ExitAfterUpdate bool

var PrepRegion string

var Threaded bool

var Profile string

func ParseCommandLine() {
	if ContinueOnBadFlag {
		Flags = flag.NewFlagSet("rubble", flag.ContinueOnError)
	} else {
		Flags = flag.NewFlagSet("rubble", flag.ExitOnError)
	}

	Flags.SetOutput(logFile)

	// Used for the -zapxxx options.
	addonCount := 0
	configCount := 0

	// Setting hardcoded defaults
	DFDir = ".."
	OutputDir = "df:raw"
	AddonsDir = "rubble:addons"

	AddonsList = new(ArgList)
	ConfigList = new(ArgList)

	NoRecover = false
	ExitAfterUpdate = false

	PrepRegion = ""

	Threaded = false

	// Load defaults from config if present
	LogPrintln("Attempting to Read Config File: ./rubble.ini")
	file, err := ioutil.ReadFile("./rubble.ini")
	if err == nil {
		LogPrintln("  Read OK, loading options from file.")
		lines := strings.Split(string(file), "\n")

		ParseINI(lines, func(key, value string) {
			switch key {
			case "dfdir":
				DFDir = value
			case "outputdir":
				OutputDir = value
			case "addonsdir":
				AddonsDir = value
			case "install":
				Installer = value
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
				LogPrintln("  Invalid config option:", key, ", skipping.")
			}
		})
	} else {
		LogPrintln("  Read failed (this is most likely ok)\n  Error:", err)
		LogPrintln("    Using hardcoded defaults.")
	}

	Flags.StringVar(&DFDir, "dfdir", DFDir, "The path to the DF directory. May be an AXIS path (only the 'rubble' location ID works).")
	Flags.StringVar(&OutputDir, "outputdir", OutputDir, "Where should Rubble write the generated raw files? May be an AXIS path (only the 'rubble' and 'df' location IDs work).")
	Flags.StringVar(&AddonsDir, "addonsdir", AddonsDir, "Rubble addons directory. May be an AXIS path (only the 'rubble', 'df', and 'out' location IDs work).")
	
	Flags.StringVar(&Installer, "install", Installer, "Path to a Rubble installer package to install.")

	Flags.Var(AddonsList, "addons", "List of addons to load. This is optional. If the value is a file path then the file is read as an ini file containing config variables. May be specified more than once.")
	Flags.Var(ConfigList, "config", "List of config variables. This is optional. If the value is a file path then the file is read as an ini file containing config variables. May be specified more than once.")

	Flags.BoolVar(&ZapAddons, "zapaddons", ZapAddons, "Ignore any -addons flags loaded from rubble.ini.")
	Flags.BoolVar(&ZapConfig, "zapconfig", ZapConfig, "Ignore any -config flags loaded from rubble.ini.")

	Flags.BoolVar(&Bench, "bench", Bench, "Display the elapsed time before exiting. Only works if no errors were encountered.")

	Flags.BoolVar(&ExitAfterUpdate, "addonlist", ExitAfterUpdate, "Update the addon list and exit.")

	Flags.BoolVar(&NoRecover, "norecover", NoRecover, "Should Rubble not recover errors? Useful for debugging.")

	Flags.StringVar(&PrepRegion, "prep", PrepRegion, "Name of a world to prepare DF for loading (or \"raw\" for the base raw folder). Use this to make sure tilesets, init changes, and DFHack scripts match the world's requirements.")

	Flags.BoolVar(&Threaded, "threads", Threaded, "Allows Rubble to use more than one processor core, not useful except for running threaded scripts.")

	Flags.StringVar(&Profile, "profile", "", "Output CPU profile information to specified file.")

	Flags.Parse(os.Args[1:])

	// Massage some of the path variables to allow AXIS paths before AXIS is setup.
	DFDir = strings.Replace(DFDir, "rubble:", "./", 1)

	OutputDir = strings.Replace(OutputDir, "rubble:", "./", 1)
	OutputDir = strings.Replace(OutputDir, "df:", DFDir+"/", 1)

	AddonsDir = strings.Replace(AddonsDir, "rubble:", "./", 1)
	AddonsDir = strings.Replace(AddonsDir, "df:", DFDir+"/", 1)
	AddonsDir = strings.Replace(AddonsDir, "out:", OutputDir+"/", 1)

	// Setup the AXIS file system.
	FS = axis.NewFileSystem()
	fs, err := axis.NewOSDir(DFDir)
	if err != nil {
		panic(err)
	}
	FS.Mount("df", fs)

	fs, err = axis.NewOSDir(".")
	if err != nil {
		panic(err)
	}
	FS.Mount("rubble", fs)

	fs, err = axis.NewOSDir(OutputDir)
	if err != nil {
		panic(err)
	}
	FS.Mount("out", fs)

	addonsFS := axis.NewFileSystem()

	fs, err = axis.NewOSDir(AddonsDir)
	if err != nil {
		panic(err)
	}
	addonsFS.Mount("dir", fs)

	zipFS := axis.NewLogicalDir()

	zips, err := ioutil.ReadDir(AddonsDir)
	if err != nil {
		panic(err)
	}

	for _, file := range zips {
		if !file.IsDir() {
			if strings.HasSuffix(file.Name(), ".zip") {
				fs, err := axiszip.NewFile(AddonsDir + "/" + file.Name())
				if err != nil {
					panic(err)
				}
				zipFS.Mount(rblutil.StripExt(file.Name()), fs)
			}
			if strings.HasSuffix(file.Name(), ".zip.b64") {
				fs, err := axiszip.NewFile64(AddonsDir + "/" + file.Name())
				if err != nil {
					panic(err)
				}
				zipFS.Mount(rblutil.StripExt(rblutil.StripExt(file.Name())), fs)
			}
		}
	}
	addonsFS.Mount("zip", zipFS)

	FS.Mount("addons", addonsFS)

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
}

// ArgsList allows an argument to exist multiple times on the command line.
type ArgList []string

func (args *ArgList) String() string {
	if len(*args) == 0 {
		return ""
	}

	rtn := (*args)[0]
	args2 := (*args)[1:]
	for i := range args2 {
		rtn += " " + args2[i]
	}
	return rtn
}

func (args *ArgList) Set(arg string) error {
	*args = append(*args, arg)
	return nil
}

func (args *ArgList) Empty() bool {
	return len(*args) != 0
}
