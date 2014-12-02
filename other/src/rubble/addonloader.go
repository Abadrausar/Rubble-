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

// TODO: Now that more than one addon dir is possible make sure duplicates are handled.
//	Not sure how, maybe add them together, but what about my plans for addon.meta?

func (state *State) LoadAddons() {
	state.loadGlobals("addons:dir:")
	for _, dir := range state.FS.ListDir("addons:dir:") {
		state.loadDir(dir, "addons:dir:"+dir)
	}

	for _, dir := range state.FS.ListDir("addons:zip:") {
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

	for _, dir := range state.FS.ListDir(dirpath) {
		state.loadDir(addonname+"/"+dir, path+dir)
	}
}

func (state *State) loadAddon(addonname, path string) {
	if _, ok := state.AddonsTbl[addonname]; ok {
		panic(Abort("Duplicate addon: " + addonname))
	}

	addon := NewAddon(addonname)

	dirpath := path
	if path != "" && path[len(path)-1] != ':' {
		path += "/"
	}

	// Load Meta File
	content, err := state.FS.ReadAll(path + "addon.meta")
	if err == nil {
		rtn, err := state.ScriptState.CompileAndRun(string(content), path+"addon.meta")
		if err != nil {
			panic(Abort(err.Error()))
		}
		meta, ok := rtn.Data.(*Meta)
		if !ok {
			panic(Abort("Addon meta file for: " + addonname + " did not return a rubble:addonmeta value!"))
		}
		addon.Meta = meta
	}

	// Load Files
	for _, filepath := range state.FS.ListFile(dirpath) {
		content, err := state.FS.ReadAll(path + filepath)
		if err != nil {
			panic(Abort(err.Error()))
		}

		file := NewAddonFile(filepath, dirpath, content)
		classifyFile(file, filepath)
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

	if strings.HasSuffix(filename, ".prep.rex") {
		file.Tags["PrepFile"] = true
		return
	}

	if strings.HasSuffix(filename, ".inst.rex") {
		file.Tags["InstScript"] = true
		return
	}

	if strings.HasSuffix(filename, ".rbl") {
		file.Tags["RawFile"] = true
		file.Tags["NoWrite"] = true
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

	for _, filepath := range state.FS.ListFile(dirpath) {
		if strings.HasSuffix(filepath, ".init.rex") {
			content, err := state.FS.ReadAll(path + filepath)
			if err != nil {
				panic(err)
			}

			state.GlobalFiles.Data[filepath] = NewAddonFile(filepath, dirpath, content)
			state.GlobalFiles.Data[filepath].Tags["InitScript"] = true
			state.GlobalFiles.Order = append(state.GlobalFiles.Order, filepath)
			continue
		}

		if strings.HasSuffix(filepath, ".load.rex") {
			content, err := state.FS.ReadAll(path + filepath)
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
	for _, filename := range source.ListFile(path) {
		if strings.HasSuffix(filename, ".pre.rex") {
			return true
		}
		if strings.HasSuffix(filename, ".post.rex") {
			return true
		}
		if strings.HasSuffix(filename, ".prep.rex") {
			return true
		}
		if strings.HasSuffix(filename, ".rbl") {
			return true
		}
		if strings.HasSuffix(filename, ".txt") {
			return true
		}
	}
	return false
}
