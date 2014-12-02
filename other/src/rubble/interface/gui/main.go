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

// Rubble GUI Interface
package main

import "reflect"
import "github.com/andlabs/ui"

import "os"
import "fmt"
import "runtime/pprof"
import "io/ioutil"
import "strconv"
import "runtime"

import "rubble"
import "rubble/rblutil"

var DFDir string
var OutputDir string
var AddonsDir *rblutil.ArgList

var AddonsList *rblutil.ArgList
var ConfigList *rblutil.ArgList

var Threaded bool

var Profile string

type Addon struct {
	Active bool
	Name   string
}

func main() {
	// TODO: Write log to a tab, AFAIK there is no multi-line edit control.
	err, log := rblutil.NewLogger()
	if err != nil {
		fmt.Println("Fatal Error:", err)
		os.Exit(1)
	}

	log.Header(rubble.Versions[0])

	DFDir = ".."
	OutputDir = "df:raw"
	AddonsDir = new(rblutil.ArgList)

	AddonsList = new(rblutil.ArgList)
	ConfigList = new(rblutil.ArgList)

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
				AddonsList.Set(value)
			case "config":
				ConfigList.Set(value)
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

	if AddonsDir.Empty() {
		AddonsDir.Set("rubble:addons")
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
		// TODO: Display message box (How?).
		os.Exit(1)
	}

	err = state.Load(*AddonsList, *ConfigList)
	if err != nil {
		if _, ok := err.(rubble.Abort); ok {
			log.Println("Abort:", err)
			// TODO: Display message box (How?).
		} else {
			log.Println("Error:", err)
			// TODO: Display message box (How?).
		}
		os.Exit(1)
	}

	loadAddons = func() []Addon {
		// TODO: handle duplicates.
		out := make([]Addon, len(state.Addons))
		for i := range state.Addons {
			out[i] = Addon{
				state.Addons[i].Active,
				state.Addons[i].Name,
			}
		}
		return out
	}

	onRun = func() {
		addons := []string{}
		d := *addonTbl.Data().(*[]Addon)
		for i := range d {
			if d[i].Active {
				addons = append(addons, d[i].Name)
			}
		}

		err = state.RunPreLoaded(addons, *ConfigList)
		if err != nil {
			if _, ok := err.(rubble.Abort); ok {
				log.Println("Abort:", err)
				// TODO: Display message box (How?).
			} else {
				log.Println("Error:", err)
				// TODO: Display message box (How?).
			}
			winMain.Close()
			ui.Stop()
		}
		log.Println("Done.")
		winMain.Close()
		ui.Stop()
	}

	go ui.Do(initGUI)
	err = ui.Go()
	if err != nil {
		log.Println("Error:", err)
		// TODO: Display message box (How?).
	}
}

var winMain ui.Window

var runBtn ui.Button
var addonTbl ui.Table

func initGUI() {
	addonTbl = ui.NewTable(reflect.TypeOf(Addon{}))
	d := addonTbl.Data().(*[]Addon)
	*d = loadAddons()
	runBtn = ui.NewButton("Run Rubble!")
	runBtn.OnClicked(onRun)
	stack := ui.NewVerticalStack(addonTbl, runBtn)
	stack.SetStretchy(0)

	tab := ui.NewTab()
	tab.Append("Addons", stack)
	tab.Append("Prep", ui.Space())
	tab.Append("Regen", ui.Space())
	tab.Append("Other", ui.Space())
	//tab.Append("Log", ui.Space())

	winMain = ui.NewWindow("Window", 400, 500, tab)
	winMain.OnClosing(func() bool {
		ui.Stop()
		return true
	})
	winMain.Show()
}

var onRun func()

var loadAddons func() []Addon
