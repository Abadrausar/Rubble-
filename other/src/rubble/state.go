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

import "fmt"

import "net/http"
import "html"

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
import "dctech/rex/commands/random"
import "dctech/rex/commands/regex"
import "dctech/rex/commands/sort"
import "dctech/rex/commands/str"
import "dctech/rex/commands/structure"
import "dctech/rex/commands/thread"

import "dctech/rexdfraw"

import "dctech/axis"
import "dctech/axis/axiszip"
import "dctech/axis/axisrex"

import "rubble/rblutil"

// This is a list of all Rubble versions in this series
// (all the listed versions are assumed to be backwards compatible)
// Index 0 MUST be the current version!
var Versions = []string{
	"5.5",
	"5.4",
	"5.3",
	"5.2",
	"5.1",
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
	PrevParams []*rex.Value

	// The file being parsed/executed right now or "".
	CurrentFile string

	// For debugging, don't touch (unless you know what you are doing).
	NoRecover bool
	
	// The position of the "last token", this may be way off, but should be good enough for error messages.
	// Set by the stage parser.
	ErrPos *rex.Position
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
	state.NewScriptTemplate("!TEMPLATE", rex.NewValueCommand(userTemplateWrap))
	state.PrevParams = make([]*rex.Value, 0)

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
	fs := axis.NewOSDir(state.DFDir)
	state.FS.Mount("df", fs)

	fs = axis.NewOSDir(".")
	state.FS.Mount("rubble", fs)

	fs = axis.NewOSDir(state.OutputDir)
	state.FS.Mount("out", fs)

	addonsFS := axis.NewFileSystem()
	zipFS := axis.NewLogicalDir()
	addonsFS.Mount("zip", zipFS)
	for i := range state.AddonsDir {
		fs = axis.NewOSDir(state.AddonsDir[i])
		addonsFS.Mount("dir", fs)
	}
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
	random.Setup(state.ScriptState)
	regex.Setup(state.ScriptState)
	sort.Setup(state.ScriptState)
	str.Setup(state.ScriptState)
	structure.Setup(state.ScriptState)
	thread.Setup(state.ScriptState)

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
	rbl.RegisterCommand("usertemplate", Command_UserTemplate)
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

func (state *State) trapError(err *error) {
	if !state.NoRecover {
		if x := recover(); x != nil {
			switch y := x.(type) {
			case RblError:
				y.Pos = state.ErrPos.Copy()
				*err = y
			case Abort:
				*err = y
			case error:
				*err = InternalError{Err: y, Pos: state.ErrPos.Copy()}
			default:
				*err = InternalError{Err: fmt.Errorf("%v", x), Pos: state.ErrPos.Copy()}
			}
		}
	}
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

	err = state.WriteReport()
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

	err = state.WriteReport()
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
	state.NewScriptTemplate("!TEMPLATE", rex.NewValueCommand(userTemplateWrap))
	state.PrevParams = make([]*rex.Value, 0)
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
	defer state.trapError(&err)

	state.Stage = StgLoad
	state.Log.Separator()
	state.Log.Println("Loading...")

	state.Log.Println("  Finding and Downloading Web Addons...")
	client := new(http.Client)
	
	for _, filename := range axis.ListFile(state.FS, "addons:dir:") {
		if strings.HasSuffix(filename, ".webload") {
			u, err := axis.ReadAll(state.FS, "addons:dir:" + filename)
			if err != nil {
				return err
			}
			url := strings.TrimSpace(string(u))
			
			r, err := client.Head(url)
			if err != nil {
				return err
			}
			
			content, err := axis.ReadAll(state.FS, "addons:dir:" + rblutil.StripExt(filename) + ".zip")
			if err == nil {
				if r.ContentLength == int64(len(content)) {
					state.Log.Println("    " + rblutil.StripExt(filename) + ": Our copy is up to date.")
					continue
				}
			}
			state.Log.Println("    " + rblutil.StripExt(filename) + ": Downloading.")
			
			r, err = client.Get(url)
			if err != nil {
				state.Log.Println("      Download Error:", err)
				continue
			}
			
			content, err = ioutil.ReadAll(r.Body)
			r.Body.Close()
			if err != nil {
				return err
			}
			
			err = axis.WriteAll(state.FS, "addons:dir:" + rblutil.StripExt(filename) + ".zip", content)
			if err != nil {
				return err
			}
		}
	}
	
	state.Log.Println("  Loading Zipped Addons...")
	aFS, err := axis.GetChild(state.FS, "addons:")
	if err != nil {
		return err
	}
	addonsFS := aFS.(axis.Collection)
	err = axis.Delete(addonsFS, "zip:")
	if err != nil {
		return err
	}
	
	zipFS := axis.NewLogicalDir()
	
	for _, filename := range axis.ListFile(state.FS, "addons:dir:") {
		if strings.HasSuffix(filename, ".zip") {
			state.Log.Println("    " + filename)
			
			content, err := axis.ReadAll(state.FS, "addons:dir:" + filename)
			if err != nil {
				return err
			}
			
			fs, err := axiszip.NewRaw(content)
			if err != nil {
				return err
			}
			zipFS.Mount(rblutil.StripExt(filename), fs)
		}
		if strings.HasSuffix(filename, ".zip.b64") {
			state.Log.Println("    " + filename)
			
			content, err := axis.ReadAll(state.FS, "addons:dir:" + filename)
			if err != nil {
				return err
			}
			
			fs, err := axiszip.NewRaw64(content)
			if err != nil {
				return err
			}
			zipFS.Mount(rblutil.StripExt(rblutil.StripExt(filename)), fs)
		}
	}
	
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
	defer state.trapError(&err)

	state.Stage = StgLoad
	state.Log.Separator()
	state.Log.Println("Activating...")

	state.Log.Println("  Loading Config Variables...")
	if config != nil && len(config) != 0 {
		for _, val := range config {
			file, err := axis.ReadAll(state.FS, val)
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
		file, err := axis.ReadAll(state.FS, val)
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
	var activate func(string)
	activate = func(me string) {
		for j := range state.AddonsTbl[me].Meta.Activates {
			it := state.AddonsTbl[me].Meta.Activates[j]
			if _, ok := state.AddonsTbl[it]; !ok {
				RaiseAbort("The \"" + state.AddonsTbl[me].Name + "\" addon requires the \"" + it + "\" addon!\n" +
					"The required addon is not currently installed, please install the required addon and try again.")
			}
			state.AddonsTbl[it].Active = true
			activate(it)
		}
	}
	for i := range state.AddonsTbl {
		if state.AddonsTbl[i].Active {
			activate(i)
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
	defer state.trapError(&err)

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
		lfiles.Data[i].Content = state.ParseFile(lfiles.Data[i].Content, StgUseCurrent, rex.NewPosition(1, 0, lfiles.Data[i].Name))
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
		lfiles.Data[i].Content = state.ParseFile(lfiles.Data[i].Content, StgUseCurrent, rex.NewPosition(1, 0, lfiles.Data[i].Name))
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
		lfiles.Data[i].Content = state.ParseFile(lfiles.Data[i].Content, StgUseCurrent, rex.NewPosition(1, 0, lfiles.Data[i].Name))
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

// Used by WriteReport
type addonheader struct {
	Name string
	Meta *Meta
}
type varheader struct {
	Name string
	Value string
	Meta *MetaVar
}

// WriteReport adds a "generation report" to the output directory.
// This report includes HTML documentation for the addons used as well as
// information about the configuration variables and their settings.
func (state *State) WriteReport() error {
	state.Log.Separator()
	state.Log.Println("Writing Generation Report...")
	
	state.Log.Println("  Loading Templates...")
	tmpl, err := rblutil.LoadHTMLTemplatesStatic(state.FS)
	if err != nil {
		return err
	}
	
	state.Log.Println("  Generating Addon Pages...")
	headers := make([]addonheader, 0, len(state.Addons))
	for _, addon := range state.Addons {
		if !addon.Active {
			continue
		}
		state.Log.Println("    " + addon.Name)
		
		header := addonheader{addon.Name, addon.Meta}
		headers = append(headers, header)
		
		path := "out:Docs/Addons/" + addon.Name + ".html"
		err := axis.Create(state.FS, path)
		if err != nil {
			return err
		}
		w, err := axis.Write(state.FS, path)
		if err != nil {
			return err
		}
		err = tmpl.ExecuteTemplate(w, "addondata", header)
		w.Close()
		if err != nil {
			return err
		}
	}
	
	state.Log.Println("  Generating Active Addon List Page...")
	err = axis.Create(state.FS, "out:Docs/Addon List.html")
	if err != nil {
		return err
	}
	w, err := axis.Write(state.FS, "out:Docs/Addon List.html")
	if err != nil {
		return err
	}
	err = tmpl.ExecuteTemplate(w, "addonlist", headers)
	w.Close()
	if err != nil {
		return err
	}
	
	state.Log.Println("  Generating Configuration Variable List Page...")
	vars := make([]varheader, 0)
	// We only generate documentation for variables that are documented in the addon.meta file.
	for _, addon := range state.Addons {
		if !addon.Active {
			continue
		}
		
		for i := range addon.Meta.Vars {
			data, ok := state.VariableData[i]
			if !ok {
				// This should never happen!
				continue
			}
			vars = append(vars, varheader{
				i,
				html.EscapeString(data),
				addon.Meta.Vars[i],
			})
		}
	}
	
	err = axis.Create(state.FS, "out:Docs/Configuration.html")
	if err != nil {
		return err
	}
	w, err = axis.Write(state.FS, "out:Docs/Configuration.html")
	if err != nil {
		return err
	}
	err = tmpl.ExecuteTemplate(w, "config", vars)
	w.Close()
	if err != nil {
		return err
	}
	
	state.Log.Println("  Generating Other Pages...")
	err = axis.Create(state.FS, "out:Docs/About.html")
	if err != nil {
		return err
	}
	w, err = axis.Write(state.FS, "out:Docs/About.html")
	if err != nil {
		return err
	}
	err = tmpl.ExecuteTemplate(w, "about", nil)
	w.Close()
	if err != nil {
		return err
	}
	
	state.Log.Println("  Generating Menu...")
	err = axis.Create(state.FS, "out:Menu.html")
	if err != nil {
		return err
	}
	w, err = axis.Write(state.FS, "out:Menu.html")
	if err != nil {
		return err
	}
	err = tmpl.ExecuteTemplate(w, "mainpage", nil)
	w.Close()
	if err != nil {
		return err
	}
	return nil
}

func (state *State) writeFile(path string, content []byte) error {
	if !axis.Exists(state.FS, path) {
		err := axis.Create(state.FS, path)
		if err != nil {
			return err
		}
	}

	return axis.WriteAll(state.FS, path, content)
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
