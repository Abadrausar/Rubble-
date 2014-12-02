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

package rubble

import "strings"
import "dctech/axis"
import "sort"

// ForEachAddon provides a generic API for running an action for each addon loaded.
// Runs even for addons that have no parseable files (or even no files at all)!
func (state *State) ForEachAddon(action func(addonname, path string)) {
	state.forEachDir("", "addons:dir:", action)
	
	for _, dir := range axis.ListDir(state.FS, "addons:zip:") {
		state.forEachDir(dir, "addons:zip:"+dir, action)
	}
}

func (state *State) forEachDir(addonname, path string, action func(addonname, path string)) {
	action(addonname, path)
	
	dirpath := path
	if path != "" {
		path += "/"
	}
	for _, dir := range axis.ListDir(state.FS, dirpath) {
		state.forEachDir(addonname+"/"+dir, path+dir, action)
	}
}

// TODO: Now that more than one addon dir is possible make sure duplicates are handled.
//	Not sure how, maybe add them together, but what about my plans for addon.meta?

func (state *State) LoadAddons() {
	state.loadGlobals("addons:dir:")
	for _, dir := range axis.ListDir(state.FS, "addons:dir:") {
		state.loadDir(dir, "addons:dir:"+dir)
	}

	for _, dir := range axis.ListDir(state.FS, "addons:zip:") {
		state.loadDir(dir, "addons:zip:"+dir)
	}
}

func (state *State) loadDir(addonname, path string) {
	dirpath := path
	if path != "" {
		path += "/"
	}

	// Load any init or load scripts.
	state.loadGlobals(dirpath)

	if containsParseable(state.FS, dirpath) {
		state.loadAddon(addonname, dirpath)
	}

	for _, dir := range axis.ListDir(state.FS, dirpath) {
		state.loadDir(addonname+"/"+dir, path+dir)
	}
}

func (state *State) loadAddon(addonname, path string) {
	if _, ok := state.AddonsTbl[addonname]; ok {
		RaiseError("Duplicate addon: " + addonname)
	}

	addon := NewAddon(addonname)

	dirpath := path
	if path != "" && path[len(path)-1] != ':' {
		path += "/"
	}

	// Load Meta File
	content, err := axis.ReadAll(state.FS, path + "addon.meta")
	if err == nil {
		rtn, err := state.ScriptState.CompileAndRun(string(content), path+"addon.meta")
		if err != nil {
			panic(err)
		}
		meta, ok := rtn.Data.(*Meta)
		if !ok {
			RaiseError("Addon meta file for: " + addonname + " did not return a rubble:addonmeta value!")
		}
		addon.Meta = meta
	} else {
		addon.Meta = NewMeta()
	}
	
	// Load Files
	for _, filepath := range axis.ListFile(state.FS, dirpath) {
		content, err := axis.ReadAll(state.FS, path + filepath)
		if err != nil {
			panic(err)
		}

		file := NewAddonFile(filepath, dirpath, content)
		classifyFile(file, filepath)
		if file.Tags["TSetFile"] || file.Tags["TSetScript"] {
			addon.Meta.TSet = true
		}
		addon.Files[filepath] = file
	}
	state.Addons = append(state.Addons, addon)
	state.AddonsTbl[addon.Name] = addon
}

func classifyFile(file *AddonFile, filename string) {
	if strings.HasSuffix(filename, ".pre.rex") {
		file.Tags["PreScript"] = true
		return
	}

	if strings.HasSuffix(filename, ".post.rex") {
		file.Tags["PostScript"] = true
		return
	}

	if strings.HasSuffix(filename, ".tset.rex") {
		file.Tags["TSetScript"] = true
		return
	}

	if strings.HasSuffix(filename, ".rbl") {
		file.Tags["RawFile"] = true
		file.Tags["NoWrite"] = true
		return
	}

	if strings.HasSuffix(filename, ".tset") {
		file.Tags["TSetFile"] = true
		return
	}

	if strings.HasSuffix(filename, ".txt") {
		file.Tags["RawFile"] = true
		return
	}
}

func (state State) loadGlobals(path string) {
	dirpath := path
	if path != "" {
		path += "/"
	}

	for _, filepath := range axis.ListFile(state.FS, dirpath) {
		if strings.HasSuffix(filepath, ".init.rex") {
			content, err := axis.ReadAll(state.FS, path + filepath)
			if err != nil {
				panic(err)
			}

			state.GlobalFiles.Data[filepath] = NewAddonFile(filepath, dirpath, content)
			state.GlobalFiles.Data[filepath].Tags["InitScript"] = true
			state.GlobalFiles.Order = append(state.GlobalFiles.Order, filepath)
			continue
		}

		if strings.HasSuffix(filepath, ".load.rex") {
			content, err := axis.ReadAll(state.FS, path + filepath)
			if err != nil {
				panic(err)
			}

			state.GlobalFiles.Data[filepath] = NewAddonFile(filepath, dirpath, content)
			state.GlobalFiles.Data[filepath].Tags["LoaderScript"] = true
			state.GlobalFiles.Order = append(state.GlobalFiles.Order, filepath)
			continue
		}
	}

	sort.Strings(state.GlobalFiles.Order)
	return
}

func containsParseable(source axis.DataSource, path string) bool {
	for _, filename := range axis.ListFile(source, path) {
		if strings.HasSuffix(filename, ".pre.rex") {
			return true
		}
		if strings.HasSuffix(filename, ".post.rex") {
			return true
		}
		if strings.HasSuffix(filename, ".tset.rex") {
			return true
		}
		if strings.HasSuffix(filename, ".rbl") {
			return true
		}
		if strings.HasSuffix(filename, ".tset") {
			return true
		}
		if strings.HasSuffix(filename, ".txt") {
			return true
		}
		if filename == "addon.meta" {
			return true
		}
	}
	return false
}
