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

func LoadAddons() []*Addon {
	addonlist := make([]*Addon, 0, 20)

	loadInit(FS, "addons:dir:")
	for _, dir := range FS.ListDir("addons:dir:") {
		addonlist = loadDir(FS, dir, "addons:dir:" + dir, addonlist)
	}

	for _, dir := range FS.ListDir("addons:zip:") {
		addonlist = loadDir(FS, dir, "addons:zip:" + dir, addonlist)
	}

	return addonlist
}

func loadDir(source axis.DataSource, addonname, path string, addons []*Addon) []*Addon {
	dirpath := path
	if path != "" {
		path += "/"
	}

	// Load a init script if any.
	loadInit(source, dirpath)

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

	// This just bloats the log for no great gain.
	//LogPrintln(addonname)

	dirpath := path
	if path != "" {
		path += "/"
	}

	for _, filepath := range source.ListFile(dirpath) {

		//LogPrintln("  " + path + filepath)

		file := new(AddonFile)
		file.Path = filepath
		file.Source = dirpath

		content, err := source.ReadAll(path + filepath)
		if err != nil {
			panic(err)
		}
		file.Content = content

		classifyFile(file, filepath)
		addon.Files[filepath] = file
	}
	return addon
}

func classifyFile(file *AddonFile, filename string) {
	if strings.HasSuffix(filename, ".pre.rex") {
		// is pre script
		file.PreScript = true
		return
	}
	if strings.HasSuffix(filename, ".post.rex") {
		// is post script
		file.PostScript = true
		return
	}
	if strings.HasSuffix(filename, ".prep.rex") {
		// is prep script
		file.PrepScript = true
		return
	}

	if strings.HasSuffix(filename, ".inst.rex") {
		// is install script
		file.InstScript = true
		return
	}

	if strings.HasSuffix(filename, ".rbl") {
		// is rubble code (do not write out after parse)
		file.NoWrite = true
		return
	}

	if strings.HasSuffix(filename, ".txt") {
		// is raw file
		return
	}

	file.UserData = true
	return
}

func loadInit(source axis.DataSource, path string) {
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
			ForcedInit[filepath] = content
			ForcedInitOrder = append(ForcedInitOrder, filepath)
			continue
		}
	}

	return
}

func containsParseable(source axis.DataSource, path string) bool {
	for _, filename := range source.ListFile(path) {
		if strings.HasSuffix(filename, ".init.rex") || strings.HasSuffix(filename, ".inst.rex") {
			continue
		}
		if strings.HasSuffix(filename, ".rex") {
			// is script
			return true
		}
		if strings.HasSuffix(filename, ".rbl") {
			// is rubble code
			return true
		}
		if strings.HasSuffix(filename, ".txt") {
			// is raw file
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
