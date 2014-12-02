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

// Rubble main package, this contains everything that makes Rubble actually work.
package rubble

import "strings"
import "strconv"
import "io/ioutil"

import "dctech/rex"
import "dctech/rex/genii"
import "dctech/rex/commands/base"
import "dctech/rex/commands/boolean"
import "dctech/rex/commands/console"
import "dctech/rex/commands/convert"
import "dctech/rex/commands/debug"

//import "dctech/rex/commands/env"
import "dctech/rex/commands/expr"

//import "dctech/rex/commands/file"
//import "dctech/rex/commands/fileio"
import "dctech/rex/commands/float"
import "dctech/rex/commands/integer"
import "dctech/rex/commands/png"
import "dctech/rex/commands/regex"
import "dctech/rex/commands/sort"
import "dctech/rex/commands/str"

import "dctech/rexdfraw"

import "dctech/axis"
import "dctech/axis/axiszip"
import "dctech/axis/axisrex"

import "rubble/rblutil"

// This is a list of all Rubble versions in this series
// (all the listed versions are assumed to be backwards compatible)
// Index 0 MUST be the current version!
var Versions = []string{
	"5.0",
	"4.7",
	"4.6",
	"4.5",
	"4.4",
	"4.3",
	"4.2",
	"4.1",
	"4.0",
	"pre4",
}

// Abort is a "marker type", eg it is used to discriminate an "error" panic from a "get me out of here" panic.
// In practice the only real difference is that aborts are normally logged with an "Abort:" prefix while
// everything else gets an "Error:" prefix
type Abort string

func (a Abort) Error() string { return string(a) }

// Parse stage constants
type ParseStage int

const (
	StgUseCurrent ParseStage = iota - 1
	StgLoad
	StgInit
	StgPreScripts
	StgPreParse
	StgParse
	StgPostParse
	StgGlobalExpand
	StgPostScripts
	StgWrite
)

// This MUST be deferred in EVERY function that runs script code!
func trapAbort(err *error) {
	if x := recover(); x != nil {
		if y, ok := x.(Abort); ok {
			*err = y
			return
		}
		panic(x)
	}
}

// State is the core of Rubble, everything connects to the state at some level.
type State struct {
	Log *rblutil.Logger

	// The global AXIS filesystem.
	FS axis.Collection

	// Important paths (OS paths, not AXIS!).
	DFDir     string
	OutputDir string
	AddonsDir []string

	// The files of the active loaded addons.
	Files *FileList

	// Global scripts (load and init scripts).
	GlobalFiles *FileList

	// All the loaded addons.
	Addons    []*Addon
	AddonsTbl map[string]*Addon

	// Like the Rubble State, Rex's State is it's core as well.
	// This is where all the persistent script stuff is stored.
	ScriptState *rex.State

	// The current parse stage, used by the stage parser.
	// Most of the stages are useless, but they don't hurt anything so I left them in.
	Stage ParseStage

	// This is where config variables are stored.
	VariableData map[string]string

	// A map of all the templates.
	Templates map[string]*Template

	// The parameters of the last template call, for use with the ... parameter.
	PrevParams []*Value

	// The file being parsed/executed right now or "".
	CurrentFile string

	// For debugging, don't touch (unless you know what you are doing).
	NoRecover bool
}

// NewState creates a new Rubble State with the paths provided.
func NewState(dfdir, output string, addonsdir []string, log *rblutil.Logger) (error, *State) {
	// First create the basic state.
	state := new(State)

	state.Log = log

	state.Log.Separator()
	state.Log.Println("Initializing...")
	state.Log.Println("  Creating New State...")

	state.FS = axis.NewFileSystem()
	state.Files = NewFileList(nil)
	state.GlobalFiles = NewFileList(nil)
	state.Addons = make([]*Addon, 0, 64)
	state.AddonsTbl = make(map[string]*Addon)
	state.ScriptState = rex.NewState()
	state.VariableData = make(map[string]string)
	state.Templates = make(map[string]*Template)
	state.NewNativeTemplate("!TEMPLATE", tempTemplate)
	state.PrevParams = make([]*Value, 0)

	// Massage some of the path variables to allow AXIS paths before AXIS is setup.
	state.DFDir = dfdir
	state.DFDir = strings.Replace(state.DFDir, "rubble:", "./", 1)

	state.OutputDir = output
	state.OutputDir = strings.Replace(state.OutputDir, "rubble:", "./", 1)
	state.OutputDir = strings.Replace(state.OutputDir, "df:", state.DFDir+"/", 1)

	state.AddonsDir = addonsdir
	for i := range state.AddonsDir {
		state.AddonsDir[i] = strings.Replace(state.AddonsDir[i], "rubble:", "./", 1)
		state.AddonsDir[i] = strings.Replace(state.AddonsDir[i], "df:", state.DFDir+"/", 1)
		state.AddonsDir[i] = strings.Replace(state.AddonsDir[i], "out:", state.OutputDir+"/", 1)
	}

	// Now setup the global stuff like Rex and AXIS.

	state.Log.Println("  Initializing AXIS VFS...")
	fs, err := axis.NewOSDir(state.DFDir)
	if err != nil {
		return err, nil
	}
	state.FS.Mount("df", fs)

	fs, err = axis.NewOSDir(".")
	if err != nil {
		return err, nil
	}
	state.FS.Mount("rubble", fs)

	fs, err = axis.NewOSDir(state.OutputDir)
	if err != nil {
		return err, nil
	}
	state.FS.Mount("out", fs)

	addonsFS := axis.NewFileSystem()
	zipFS := axis.NewLogicalDir()
	for i := range state.AddonsDir {
		fs, err = axis.NewOSDir(state.AddonsDir[i])
		if err != nil {
			return err, nil
		}
		addonsFS.Mount("dir", fs)

		zips, err := ioutil.ReadDir(state.AddonsDir[i])
		if err != nil {
			return err, nil
		}

		for _, file := range zips {
			if !file.IsDir() {
				if strings.HasSuffix(file.Name(), ".zip") {
					fs, err := axiszip.NewFile(state.AddonsDir[i] + "/" + file.Name())
					if err != nil {
						return err, nil
					}
					zipFS.Mount(rblutil.StripExt(file.Name()), fs)
				}
				if strings.HasSuffix(file.Name(), ".zip.b64") {
					fs, err := axiszip.NewFile64(state.AddonsDir[i] + "/" + file.Name())
					if err != nil {
						return err, nil
					}
					zipFS.Mount(rblutil.StripExt(rblutil.StripExt(file.Name())), fs)
				}
			}
		}
	}
	addonsFS.Mount("zip", zipFS)

	state.FS.Mount("addons", addonsFS)

	state.Log.Println("  Initializing Rex Scripting...")
	// Load most commands
	rex.SetupArrays(state.ScriptState)
	rex.SetupMaps(state.ScriptState)
	base.Setup(state.ScriptState)
	boolean.Setup(state.ScriptState)
	console.Setup(state.ScriptState)
	convert.Setup(state.ScriptState)
	debug.Setup(state.ScriptState)
	//env.Setup(state.ScriptState)
	expr.Setup(state.ScriptState)
	//file.Setup(state.ScriptState)
	//fileio.Setup(state.ScriptState)
	float.Setup(state.ScriptState)
	integer.Setup(state.ScriptState)
	png.Setup(state.ScriptState)
	regex.Setup(state.ScriptState)
	sort.Setup(state.ScriptState)
	str.Setup(state.ScriptState)

	genii.Setup(state.ScriptState)

	rexdfraw.Setup(state.ScriptState)

	axisrex.Setup(state.ScriptState)

	rbl := state.ScriptState.RegisterModule("rubble")

	rbl.RegisterVar("version", rex.NewValueString(Versions[0]))
	versions := make(map[string]*rex.Value, len(Versions))
	for i := range Versions {
		versions[Versions[i]] = rex.NewValueBool(true)
	}
	rbl.RegisterVar("versions", rex.NewValueIndex(rex.NewStaticMap(versions)))

	rbl.RegisterVar("fs", rex.NewValueUser(state.FS))
	rbl.RegisterVar("state", rex.NewValueUser(state))

	rbl.RegisterVar("addons", genii.New(&state.Addons))
	rbl.RegisterVar("addonstbl", genii.New(&state.AddonsTbl))

	rbl.RegisterVar("files", genii.New(&state.Files))

	rbl.RegisterVar("raws", rex.NewValueIndex(NewIndexableRaws(&state.Files)))

	rbl.RegisterType("addonmeta", NewMetaFromLit)
	rbl.RegisterType("addonmetavar", NewMetaVarFromLit)

	rbl.RegisterCommand("currentfile", Command_CurrentFile)

	rbl.RegisterCommand("template", Command_Template)
	rbl.RegisterCommand("stageparse", Command_Parse)
	rbl.RegisterCommand("calltemplate", Command_CallTemplate)
	rbl.RegisterCommand("expandvars", Command_ExpandVars)

	rbl.RegisterCommand("patch", Command_Patch)
	rbl.RegisterCommand("decompress", Command_Decompress)
	rbl.RegisterCommand("compress", Command_Compress)

	rbl.RegisterCommand("abort", Command_Abort)

	rbl.RegisterCommand("configvar", Command_ConfigVar)

	rbl.RegisterCommand("filetag", Command_FileTag)
	rbl.RegisterCommand("gfiletag", Command_GFileTag)

	rbl.RegisterCommand("newaddon", Command_NewAddon)
	rbl.RegisterCommand("newfile", Command_NewFile)

	// Redirect output to logger
	state.ScriptState.Output = state.Log

	return nil, state
}

// Run runs a full Rubble parse cycle.
// See *State.Load for parameter descriptions.
func (state *State) Run(addons, config []string) (err error) {
	err = state.Load(addons, config)
	if err != nil {
		return err
	}

	err = state.Parse()
	if err != nil {
		return err
	}

	err = state.Write()
	if err != nil {
		return err
	}

	return nil
}

// RunPreLoaded runs a full Rubble parse cycle, minus loading addons.
// See *State.Activate for parameter descriptions.
func (state *State) RunPreLoaded(addons, config []string) (err error) {
	err = state.Activate(addons, config)
	if err != nil {
		return err
	}

	err = state.Parse()
	if err != nil {
		return err
	}

	err = state.Write()
	if err != nil {
		return err
	}

	return nil
}

// Clear clears everything loaded by Load and resets most other fields as well.
// Things not cleared:
//	Script state
//	Anything related to AXIS
//	Anything related to the Logger
func (state *State) Clear() {
	state.Files = NewFileList(nil)
	state.GlobalFiles = NewFileList(nil)
	state.Addons = make([]*Addon, 0, 64)
	state.AddonsTbl = make(map[string]*Addon)
	state.VariableData = make(map[string]string)
	state.Templates = make(map[string]*Template)
	state.NewNativeTemplate("!TEMPLATE", tempTemplate)
	state.PrevParams = make([]*Value, 0)
}

// Load loads all addons and determines which ones should be
// active, then writes the default addon list file.
// This is where most of the configuration magic happens...
// addons contains addon activation information.
// Each entry must be either:
//	A list of addon names delimited by semicolons (each addon given is activated)
//	The name of an INI file that contains addon names and their activation state, paths must use the AXIS syntax.
// If addons is nil the default addons file is used (addons:dir:addonlist.ini).
// config is exactly the same as addons, just for configuration variables.
func (state *State) Load(addons, config []string) (err error) {
	defer trapAbort(&err)

	state.Stage = StgLoad
	state.Log.Separator()
	state.Log.Println("Loading...")

	state.Log.Println("  Loading Addons...")
	state.LoadAddons()

	return state.Activate(addons, config)
}

// Activate reruns some of the steps Load does.
// addons contains addon activation information.
// Each entry must be either:
//	A list of addon names delimited by semicolons (each addon given is activated)
//	The name of an INI file that contains addon names and their activation state, paths must use the AXIS syntax.
// If addons is empty the default addons file is used (addons:dir:addonlist.ini).
// config is exactly the same as addons, just for configuration variables.
func (state *State) Activate(addons, config []string) (err error) {
	defer trapAbort(&err)

	state.Stage = StgLoad
	state.Log.Separator()
	state.Log.Println("Activating...")

	state.Log.Println("  Loading Config Variables...")
	if config != nil && len(config) != 0 {
		for _, val := range config {
			file, err := state.FS.ReadAll(val)
			if err == nil {
				state.Log.Println("    Loading Config File: " + val)
				rblutil.ParseINI(string(file), "\n", func(key, value string) {
					state.VariableData[key] = value
				})
				continue
			}
			rblutil.ParseINI(val, ";", func(key, value string) {
				state.VariableData[key] = value
			})
		}
	} else {
		state.Log.Println("    No variables specified.")
	}

	state.Log.Println("  Generating Active Addon List...")
	if addons == nil {
		state.Log.Println("    No addons specified, using default addon list file.")
		addons = []string{"addons:dir:addonlist.ini"}
	}

	addonNames := make([]string, 0, 10)
	for _, val := range addons {
		file, err := state.FS.ReadAll(val)
		if err == nil {
			state.Log.Println("    Loading List File: " + val)
			rblutil.ParseINI(string(file), "\n", func(key, value string) {
				value = strings.ToLower(value)
				if ok, _ := strconv.ParseBool(value); ok {
					addonNames = append(addonNames, key)
				}
			})
		} else {
			addonNames = append(addonNames, strings.Split(val, ";")...)
		}
	}

	state.Log.Println("  Activating Addons...")
	for _, name := range addonNames {
		if _, ok := state.AddonsTbl[name]; ok {
			state.AddonsTbl[name].Active = true
		}
	}

	state.Log.Println("  Pruning Library Addons...")
	for i := range state.Addons {
		if state.Addons[i].Meta.Lib {
			state.Addons[i].Active = false
		}
	}

	state.Log.Println("  Activating Required Addons from Meta Data...")
	for i := range state.Addons {
		if state.Addons[i].Active {
			for j := range state.Addons[i].Meta.Activates {
				name := state.Addons[i].Meta.Activates[j]
				if _, ok := state.AddonsTbl[name]; !ok {
					panic(Abort("The \"" + state.Addons[i].Name + "\" addon requires the \"" + name + "\" addon!\n" +
						"The required addon is not currently installed, please install the required addon and try again."))
				}
				state.AddonsTbl[name].Active = true
			}
		}
	}

	state.Log.Println("  Running Loader Scripts...")
	gfiles := state.GlobalFiles
	for _, i := range gfiles.Order {
		if gfiles.Data[i].Tags["Skip"] || !gfiles.Data[i].Tags["LoaderScript"] {
			continue
		}

		state.CurrentFile = gfiles.Data[i].Name
		state.Log.Println("    " + gfiles.Data[i].Name)

		_, err := state.ScriptState.CompileAndRun(string(gfiles.Data[i].Content), gfiles.Data[i].Name)
		if err != nil {
			return err
		}
	}

	state.Log.Println("  Active Addons:")
	for i := range state.Addons {
		if state.Addons[i].Active {
			state.Log.Println("    " + state.Addons[i].Name)
		}
	}

	state.Log.Println("  Updating the Default Addon List File...")
	return state.UpdateAddonList("addons:dir:addonlist.ini")
}

// Parse handles everything from running init scripts to running post scripts.
// If an error is returned state.CurrentFile and state.Stage will still be set to their last values.
func (state *State) Parse() (err error) {
	defer trapAbort(&err)

	state.Log.Separator()
	state.Log.Println("Generating Sorted Active File List...")
	state.Files = NewFileList(state.Addons)

	state.Stage = StgInit
	state.Log.Separator()
	state.Log.Println("Running Init Scripts...")
	gfiles := state.GlobalFiles
	for _, i := range gfiles.Order {
		if gfiles.Data[i].Tags["Skip"] || !gfiles.Data[i].Tags["InitScript"] {
			continue
		}

		state.CurrentFile = gfiles.Data[i].Name
		state.Log.Println("  " + gfiles.Data[i].Name)

		_, err := state.ScriptState.CompileAndRun(string(gfiles.Data[i].Content), gfiles.Data[i].Name)
		if err != nil {
			return err
		}
	}

	state.Stage = StgPreScripts
	state.Log.Separator()
	state.Log.Println("Running Prescripts...")
	lfiles := state.Files
	for _, i := range lfiles.Order {
		if lfiles.Data[i].Tags["Skip"] || !lfiles.Data[i].Tags["PreScript"] {
			continue
		}

		state.CurrentFile = lfiles.Data[i].Name
		state.Log.Println("  " + state.Files.Data[i].Name)

		_, err := state.ScriptState.CompileAndRun(string(lfiles.Data[i].Content), lfiles.Data[i].Name)
		if err != nil {
			return err
		}
	}

	state.Stage = StgPreParse
	state.Log.Separator()
	state.Log.Println("Preparsing...")
	for _, i := range state.Files.Order {
		if lfiles.Data[i].Tags["Skip"] || !lfiles.Data[i].Tags["RawFile"] {
			continue
		}

		state.CurrentFile = lfiles.Data[i].Name
		state.Log.Println("  " + lfiles.Data[i].Name)
		lfiles.Data[i].Content = state.ParseFile(lfiles.Data[i].Content, StgUseCurrent, NewPosition(1, lfiles.Data[i].Name))
	}

	state.Stage = StgParse
	state.Log.Separator()
	state.Log.Println("Parsing...")
	for _, i := range state.Files.Order {
		if lfiles.Data[i].Tags["Skip"] || !lfiles.Data[i].Tags["RawFile"] {
			continue
		}

		state.CurrentFile = lfiles.Data[i].Name
		state.Log.Println("  " + lfiles.Data[i].Name)
		lfiles.Data[i].Content = state.ParseFile(lfiles.Data[i].Content, StgUseCurrent, NewPosition(1, lfiles.Data[i].Name))
	}

	state.Stage = StgPostParse
	state.Log.Separator()
	state.Log.Println("Postparsing...")
	for _, i := range state.Files.Order {
		if lfiles.Data[i].Tags["Skip"] || !lfiles.Data[i].Tags["RawFile"] {
			continue
		}

		state.CurrentFile = lfiles.Data[i].Name
		state.Log.Println("  " + lfiles.Data[i].Name)
		lfiles.Data[i].Content = state.ParseFile(lfiles.Data[i].Content, StgUseCurrent, NewPosition(1, lfiles.Data[i].Name))
	}

	state.Stage = StgGlobalExpand
	state.Log.Separator()
	state.Log.Println("Expanding Variables...")
	for _, i := range lfiles.Order {
		if lfiles.Data[i].Tags["Skip"] || !lfiles.Data[i].Tags["RawFile"] {
			continue
		}

		state.CurrentFile = lfiles.Data[i].Name
		state.Log.Println("  " + lfiles.Data[i].Name)
		lfiles.Data[i].Content = []byte(state.ExpandVars(string(lfiles.Data[i].Content)))
	}

	state.Stage = StgPostScripts
	state.Log.Separator()
	state.Log.Println("Running Postscripts...")
	for _, i := range lfiles.Order {
		if lfiles.Data[i].Tags["Skip"] || !lfiles.Data[i].Tags["PostScript"] {
			continue
		}

		state.CurrentFile = lfiles.Data[i].Name
		state.Log.Println("  " + state.Files.Data[i].Name)

		_, err := state.ScriptState.CompileAndRun(string(lfiles.Data[i].Content), lfiles.Data[i].Name)
		if err != nil {
			return err
		}
	}

	state.CurrentFile = ""
	return nil
}

// Write handles writing the files to their output directories.
// TODO: Maybe scripts that write custom outputs should queue them up instead and have this function write EVERYTHING out?
func (state *State) Write() error {
	state.Stage = StgWrite
	state.Log.Separator()
	state.Log.Println("Writing Files...")
	state.Log.Println("  Writing Raw Files...")
	lfiles := state.Files
	for _, i := range lfiles.Order {
		if lfiles.Data[i].Tags["Skip"] || lfiles.Data[i].Tags["NoWrite"] || lfiles.Data[i].Tags["GraphicsFile"] || !lfiles.Data[i].Tags["RawFile"] {
			continue
		}

		state.Log.Println("    " + lfiles.Data[i].Name)

		file := []byte(rblutil.StripExt(lfiles.Data[i].Name) + "\n\n# Automatically generated, do not edit!\n# Source: " +
			lfiles.Data[i].Source + "/" + lfiles.Data[i].Name + "\n\n" + string(lfiles.Data[i].Content))
		err := state.writeFile("out:objects/"+lfiles.Data[i].Name, file)
		if err != nil {
			return err
		}
	}

	state.Log.Println("  Writing Graphics Files...")
	for _, i := range lfiles.Order {
		if lfiles.Data[i].Tags["Skip"] || lfiles.Data[i].Tags["NoWrite"] || !lfiles.Data[i].Tags["GraphicsFile"] {
			continue
		}

		state.Log.Println("    " + lfiles.Data[i].Name)

		file := []byte(rblutil.StripExt(lfiles.Data[i].Name) + "\n\n# Automatically generated, do not edit!\n# Source: " +
			lfiles.Data[i].Source + "/" + lfiles.Data[i].Name + "\n\n" + string(lfiles.Data[i].Content))
		err := state.writeFile("out:graphics/"+lfiles.Data[i].Name, file)
		if err != nil {
			return err
		}
	}

	state.Log.Println("  Writing Prep Files...")
	for _, i := range lfiles.Order {
		if lfiles.Data[i].Tags["Skip"] || lfiles.Data[i].Tags["NoWrite"] || !lfiles.Data[i].Tags["PrepFile"] {
			continue
		}

		state.Log.Println("    " + lfiles.Data[i].Name)

		err := state.writeFile("out:prep/"+lfiles.Data[i].Name, lfiles.Data[i].Content)
		if err != nil {
			return err
		}
	}

	state.Log.Println("  Writing Addon List to Raw Directory...")
	state.Log.Println("    addonlist.ini")
	err := state.UpdateAddonList("out:addonlist.ini")
	if err != nil {
		return err
	}

	state.Log.Println("  Writing Config Variables to Raw Directory...")
	state.Log.Println("    genconfig.ini")
	return state.DumpConfig("out:genconfig.ini")
}

func (state *State) writeFile(path string, content []byte) error {
	if !state.FS.Exists(path) {
		err := state.FS.Create(path)
		if err != nil {
			return err
		}
	}

	return state.FS.WriteAll(path, content)
}

// Data Dumpers

// DumpConfig writes all configuration variables (in INI format) to the indicated file.
func (state *State) DumpConfig(path string) error {
	out := "\n# Rubble config variable dump.\n# Automatically generated, do not edit!\n\n[config]\n"

	for i := range state.VariableData {
		out += i + " = " + state.VariableData[i] + "\n"
	}

	return state.writeFile(path, []byte(out))
}

// UpdateAddonList writes a list of all addons and their activation status (in INI format)- to the indicated file.
func (state *State) UpdateAddonList(dest string) error {
	out := make([]byte, 0, 2048)
	out = append(out, "\n# Rubble addon list.\n# Version: "+Versions[0]+"\n# Automatically generated, do not edit!\n\n[addons]\n"...)
	for i := range state.Addons {
		if !state.Addons[i].Meta.Lib {
			out = append(out, state.Addons[i].Name+"="...)
			if state.Addons[i].Active {
				out = append(out, "true\n"...)
			} else {
				out = append(out, "false\n"...)
			}
		}
	}

	return state.writeFile(dest, out)
}
