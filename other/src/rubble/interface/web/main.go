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

// Rubble Web GUI Interface
package main

import "net/http"

import "os"
import "fmt"
import "bytes"
import "runtime/pprof"
import "time"
import "io/ioutil"
import "strings"
import "strconv"
import "runtime"

import "flag"

import "rubble"
import "rubble/rblutil"

import "dctech/axis"

import "os/exec"
import "path/filepath"

var Addr string

var DFDir string
var OutputDir string
var AddonsDir *rblutil.ArgList

var AddonsList *rblutil.ArgList
var ConfigList *rblutil.ArgList

var Threaded bool

var Profile string

func main() {
	logBuffer := new(bytes.Buffer)
	err, log := rblutil.NewLogger(logBuffer)
	if err != nil {
		fmt.Println("Fatal Error:", err)
		os.Exit(1)
	}

	flags := flag.NewFlagSet("rubble", flag.ExitOnError)
	flags.SetOutput(log)

	log.Header(rubble.Versions[0])

	defer func(){
		err := recover()
		if err != nil {
			log.Println("Unrecovered Error:")
			log.Println("  The following error was not properly recovered, please report this ASAP!")
			log.Printf("  %#v\n", err)
			log.Println("Stack Trace:")
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, true)]
			log.Printf("%s\n", buf)
			os.Exit(1)
		}
	}()

	Addr = "127.0.0.1:2120"

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
			case "addr":
				Addr = value
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

	flags.StringVar(&Addr, "addr", Addr, "The server address.")

	flags.StringVar(&DFDir, "dfdir", DFDir, "The path to the DF directory. May be an AXIS path (only the 'rubble' location ID works).")
	flags.StringVar(&OutputDir, "outputdir", OutputDir, "Where should Rubble write the generated raw files? May be an AXIS path (only the 'rubble' and 'df' location IDs work).")
	flags.Var(AddonsDir, "addonsdir", "Rubble addons directory. May be an AXIS path (only the 'rubble', 'df', and 'out' location IDs work).")

	flags.BoolVar(&Threaded, "threads", Threaded, "Allows Rubble to use more than one processor core, not useful except for running threaded scripts.")

	flags.StringVar(&Profile, "profile", "", "Output CPU profile information to specified file.")

	flags.Parse(os.Args[1:])

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
		os.Exit(1)
	}

	// Note: http.Redirect acts "funny" when you use http.StatusInternalServerError,
	// so I use http.StatusFound even for error-triggered redirects to the log page.
	// (this may just be my browser, but better safe than sorry)

	// Main menu
	http.HandleFunc("/menu", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/menu\"")

		var err error
		err, state = rubble.NewState(DFDir, OutputDir, *AddonsDir, log)
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}

		err = state.Load(*AddonsList, *ConfigList)
		if err != nil {
			if _, ok := err.(rubble.Abort); ok {
				log.Println("Abort:", err)
			} else {
				log.Println("Error:", err)
			}
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}

		fmt.Fprint(w, menuPage)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "./menu", http.StatusFound)
	})

	// Addon List

	http.HandleFunc("/addons", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/addons\"")

		err = addonsPage.Execute(w, state.Addons)
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}
	})

	http.HandleFunc("/addonfile", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/addonfile\"")

		addon, ok := state.AddonsTbl[r.FormValue("addon")]
		if !ok {
			http.Error(w, "Addon does not exist.", http.StatusNotFound)
			return
		}

		file, ok := addon.Files[r.FormValue("file")]
		if !ok {
			http.Error(w, "Addon file does not exist.", http.StatusNotFound)
			return
		}

		err = addonfilePage.Execute(w, file)
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}
	})

	// Normal Generation

	http.HandleFunc("/genaddons", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/genaddons\"")

		err = genaddonsPage.Execute(w, state.Addons)
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}
	})

	http.HandleFunc("/genbranch", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/genbranch\"")

		err := r.ParseForm()
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}

		for i := range r.Form {
			if i == "__EditConfig__" {
				http.Redirect(w, r, "./genvars?"+r.URL.RawQuery, http.StatusFound)
				return
			} else if i == "__Generate__" {
				http.Redirect(w, r, "./genrun?"+r.URL.RawQuery, http.StatusFound)
				return
			}
		}
		http.Redirect(w, r, "./genrun?"+r.URL.RawQuery, http.StatusFound)
	})

	http.HandleFunc("/genvars", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/genvars\"")

		err := r.ParseForm()
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}

		addons := []string{}
		for i := range r.Form {
			if v := r.Form[i]; len(v) > 0 {
				if v[0] == "true" {
					addons = append(addons, i)
				}
			}
		}

		vars := make(map[string]*rubble.MetaVar)
		for _, name := range addons {
			addon := state.AddonsTbl[name]
			for i := range addon.Meta.Vars {
				vars[i] = addon.Meta.Vars[i]
			}
		}

		err = genvarsPage.Execute(w, struct {
			Addons []string
			Vars   map[string]*rubble.MetaVar
		}{addons, vars})
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}
	})

	http.HandleFunc("/genrun", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/genrun\"")

		err := r.ParseForm()
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}

		addons := []string{}
		config := []string{}
		for i := range r.Form {
			if strings.HasPrefix(i, "__CONFIG_VAR_") {
				if v := r.Form[i]; len(v) > 0 {
					config = append(config, strings.TrimPrefix(i, "__CONFIG_VAR_")+"="+v[0])
				}
				continue
			}

			if v := r.Form[i]; len(v) > 0 {
				if v[0] == "true" {
					addons = append(addons, i)
				}
			}
		}
		ConfigList.Set(strings.Join(config, ";"))

		err = state.RunPreLoaded(addons, *ConfigList)
		if err != nil {
			if _, ok := err.(rubble.Abort); ok {
				log.Println("Abort:", err)
			} else {
				log.Println("Error:", err)
			}
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}
		log.Println("Done.")

		http.Redirect(w, r, "./log", http.StatusFound)
	})

	// Prep

	http.HandleFunc("/prep", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/prep\"")

		if !axis.Exists(state.FS, "df:data/save") {
			log.Println("Error: Cannot find save directory.")
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}

		regions := axis.ListDir(state.FS, "df:data/save")
		for i := range regions {
			if regions[i] == "current" {
				regions = append(regions[:i], regions[i+1:]...)
				break
			}
		}

		err = prepPage.Execute(w, regions)
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}
	})

	http.HandleFunc("/preprun", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/preprun\"")

		err := state.PrepModeRun(r.FormValue("region"))
		if err != nil {
			if _, ok := err.(rubble.Abort); ok {
				log.Println("Abort:", err)
			} else {
				log.Println("Error:", err)
			}
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}
		log.Println("Done.")

		http.Redirect(w, r, "./log", http.StatusFound)
	})

	// Regen

	http.HandleFunc("/regen", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/regen\"")

		fmt.Fprint(w, regenPage)
	})

	http.HandleFunc("/regenregion", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/regenregion\"")

		if !axis.Exists(state.FS, "df:data/save") {
			log.Println("Error: Cannot find save directory.")
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}

		regions := axis.ListDir(state.FS, "df:data/save")
		for i := range regions {
			if regions[i] == "current" {
				regions = append(regions[:i], regions[i+1:]...)
				break
			}
		}

		err = regenregionPage.Execute(w, regions)
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}
	})

	http.HandleFunc("/regenaddons", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/regenaddons\"")

		// Activate addons from region addon list
		region := r.FormValue("region")
		state.OutputDir = state.DFDir + "/data/save" + region + "/raw"
		state.Activate([]string{"df:data/save/" + region + "/raw/addonlist.ini"}, []string{"df:data/save/" + region + "/raw/genconfig.ini"})

		err = regenaddonsPage.Execute(w, struct {
			Region string
			Addons []*rubble.Addon
		}{region, state.Addons})
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}
	})

	http.HandleFunc("/regenbranch", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/regenbranch\"")

		err := r.ParseForm()
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}

		for i := range r.Form {
			if i == "__EditConfig__" {
				http.Redirect(w, r, "./regenvars?"+r.URL.RawQuery, http.StatusFound)
				return
			} else if i == "__Generate__" {
				http.Redirect(w, r, "./regenrun?"+r.URL.RawQuery, http.StatusFound)
				return
			}
		}
		http.Redirect(w, r, "./regenrun?"+r.URL.RawQuery, http.StatusFound)
	})

	http.HandleFunc("/regenvars", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/regenvars\"")

		err := r.ParseForm()
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}

		addons := []string{}
		for i := range r.Form {
			if v := r.Form[i]; len(v) > 0 {
				if v[0] == "true" {
					addons = append(addons, i)
				}
			}
		}

		vars := make(map[string]*rubble.MetaVar)
		for _, name := range addons {
			addon := state.AddonsTbl[name]
			for i := range addon.Meta.Vars {
				data, ok := state.VariableData[i]
				if !ok {
					vars[i] = addon.Meta.Vars[i]
					continue
				}
				val := new(rubble.MetaVar)
				val.Name = addon.Meta.Vars[i].Name
				val.Val = data
				val.Choices = addon.Meta.Vars[i].Choices
				vars[i] = val
			}
		}

		err = regenvarsPage.Execute(w, struct {
			Addons []string
			Vars   map[string]*rubble.MetaVar
		}{addons, vars})
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}
	})

	http.HandleFunc("/regenrun", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/regenrun\"")

		err := r.ParseForm()
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}

		addons := []string{}
		config := []string{}
		for i := range r.Form {
			if strings.HasPrefix(i, "__CONFIG_VAR_") {
				if v := r.Form[i]; len(v) > 0 {
					config = append(config, strings.TrimPrefix(i, "__CONFIG_VAR_")+"="+v[0])
				}
				continue
			}

			if v := r.Form[i]; len(v) > 0 {
				if v[0] == "true" {
					addons = append(addons, i)
				}
			}
		}
		ConfigList.Set(strings.Join(config, ";"))

		err = state.RunPreLoaded(addons, *ConfigList)
		if err != nil {
			if _, ok := err.(rubble.Abort); ok {
				log.Println("Abort:", err)
			} else {
				log.Println("Error:", err)
			}
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}
		log.Println("Done.")

		http.Redirect(w, r, "./log", http.StatusFound)
	})

	// Common

	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/log\"")

		err = logPage.Execute(w, logBuffer.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusFound)
		}
	})

	http.HandleFunc("/addondata", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/addondata\"")

		addon, ok := state.AddonsTbl[r.FormValue("addon")]
		if !ok {
			http.Error(w, "Addon does not exist.", http.StatusNotFound)
			return
		}

		err = addondataPage.Execute(w, addon)
		if err != nil {
			log.Println("Error:", err)
			http.Redirect(w, r, "./log", http.StatusFound)
			return
		}
	})

	http.HandleFunc("/css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		fmt.Fprint(w, css)
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/ico")
		fmt.Fprint(w, ico)
	})

	http.HandleFunc("/extras/", func(w http.ResponseWriter, r *http.Request) {
		file := strings.TrimPrefix(r.URL.Path, "/extras/")
		if !isAbs(file) {
			http.Error(w, "Extra file path invalid.", http.StatusNotFound)
			return
		}
		content, err := ioutil.ReadFile("./other/webUI/extras/" + file)
		if err != nil {
			http.Error(w, "Extra file not found.", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "%s", content)
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		log.Println("UI Transition: \"/about\"")

		fmt.Fprint(w, aboutPage)
	})

	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, killPage)
		
		d, _ := time.ParseDuration("1s")
		time.AfterFunc(d, func() { os.Exit(0) })
	})

	log.Separator()
	log.Println("Attempting to Start Your Web Browser (wish me luck)...")
	log.Println("  Attempting to Start: \"./other/webUI/browser\" \"http://" + Addr + "/menu\"")
	path, err := exec.LookPath("./other/webUI/browser")
	if err != nil {
		log.Println("    Browser Startup Failed:\n      Error:", err)
	} else {
		path, err := filepath.Abs(path)
		if err != nil {
			log.Println("    Browser Startup Failed:\n      Error:", err)
		} else {
			cmd := exec.Command(path, "http://" + Addr + "/menu")
			err = cmd.Start()
			if err != nil {
				log.Println("    Browser Startup Failed:\n      Error:", err)
			} else {
				log.Println("  As far as I can tell everything went fine.")
			}
		}
	}

	log.Println("Starting Server...")
	err = http.ListenAndServe(Addr, nil)
	if err != nil {
		log.Println("  Server Startup Failed:\n    Error:", err)
		os.Exit(1)
	}
}

// isAbs returns true if the path is not a relative path (includes no "." or ".." parts).
func isAbs(path string) bool {
	dirs := strings.Split(path, "/")
	
	for i := range dirs {
		if dirs[i] == ".." || dirs[i] == "." {
			return false
		}
	}
	return true
}
