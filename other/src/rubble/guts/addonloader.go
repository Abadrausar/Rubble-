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

import "strings"
import "dctech/axis"
import "sort"

func LoadAddons(addons *[]*Addon)  {
	loadGlobals(FS, "addons:dir:")
	for _, dir := range FS.ListDir("addons:dir:") {
		*addons = loadDir(FS, dir, "addons:dir:" + dir, *addons)
	}

	for _, dir := range FS.ListDir("addons:zip:") {
		*addons = loadDir(FS, dir, "addons:zip:" + dir, *addons)
	}
}

func loadDir(source axis.DataSource, addonname, path string, addons []*Addon) []*Addon {
	dirpath := path
	if path != "" {
		path += "/"
	}

	// Load any init or load scripts.
	loadGlobals(source, dirpath)

	if containsParseable(source, dirpath) {
		addons = append(addons, loadAddon(source, addonname, dirpath))
	}

	for _, dir := range source.ListDir(dirpath) {
		addons = loadDir(source, addonname + "/" + dir, path + dir, addons)
	}
		
	return addons
}

func loadAddon(source axis.DataSource, addonname, path string) *Addon {
	addon := NewAddon(addonname)

	dirpath := path
	if path != "" && path[len(path) - 1] != ':' {
		path += "/"
	}

	for _, filepath := range source.ListFile(dirpath) {
		content, err := source.ReadAll(path + filepath)
		if err != nil {
			panic(err)
		}

		file := NewAddonFile(filepath, dirpath, content)
		classifyFile(file, filepath)
		addon.Files[filepath] = file
	}
	return addon
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

func loadGlobals(source axis.DataSource, path string) {
	dirpath := path
	if path != "" {
		path += "/"
	}

	for _, filepath := range source.ListFile(dirpath) {
		if strings.HasSuffix(filepath, ".init.rex") {
			content, err := source.ReadAll(path + filepath)
			if err != nil {
				panic(err)
			}
			
			GlobalFiles.Data[filepath] = NewAddonFile(filepath, dirpath, content)
			GlobalFiles.Data[filepath].Tags["InitScript"] = true
			GlobalFiles.Order = append(GlobalFiles.Order, filepath)
			continue
		}
		
		if strings.HasSuffix(filepath, ".load.rex") {
			content, err := source.ReadAll(path + filepath)
			if err != nil {
				panic(err)
			}
			
			GlobalFiles.Data[filepath] = NewAddonFile(filepath, dirpath, content)
			GlobalFiles.Data[filepath].Tags["LoaderScript"] = true
			GlobalFiles.Order = append(GlobalFiles.Order, filepath)
			continue
		}
	}
	
	sort.Strings(GlobalFiles.Order)
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

func UpdateAddonList(dest string, addons []*Addon) {
	out := make([]byte, 0, 2048)
	out = append(out, "\n# Rubble addon list.\n# Version: " + RubbleVersions[0] + "\n# Automatically generated, do not edit!\n\n[addons]\n"...)
	for i := range addons {
		out = append(out, addons[i].Name + "="...)
		if addons[i].Active {
			out = append(out, "true\n"...)
		} else {
			out = append(out, "false\n"...)
		}
	}

	WriteFile(dest, out)
}

func StripExt(name string) string {
	i := len(name) - 1
	for i >= 0 {
		if name[i] == '.' {
			return name[:i]
		}
		i--
	}
	return name
}
