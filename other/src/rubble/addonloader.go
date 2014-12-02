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

package main

import "strings"
import "dctech/dcfs"
import "io/ioutil"

func LoadAddons(path string) []*Addon {
	addonlist := make([]*Addon, 0, 20)
	
	LogPrintln(path)
	source, err := dcfs.NewDirReader(path)
	if err != nil {
		if ShellMode {
			return addonlist
		}
		panic(err)
	}

	for _, dir := range source.ListDirs("") {
		addonlist = loadDir(source, dir, dir, addonlist)
	}

	for _, file := range source.ListFiles("") {
		if strings.HasSuffix(file, ".zip") {
			zip, err := dcfs.NewZipReader(path + "/" + file)
			if err != nil {
				panic(err)
			}

			addonlist = loadDir(zip, StripExt(file), "", addonlist)
		}
	}

	return addonlist
}

func loadDir(source dcfs.DataSource, addonname, path string, addons []*Addon) []*Addon {
	dirpath := path
	if path != "" {
		path += "/"
	}
	
	// Load a init script if any.
	loadInit(source, dirpath)
	
	if containsParseable(source, dirpath) {
		addons = append(addons, loadAddon(source, addonname, dirpath))
	}
	
	dirs := source.ListDirs(dirpath)
	if len(dirs) != 0 {
		for _, dir := range dirs {
			addons = loadDir(source, addonname + "/" + dir, path + dir, addons)
		}
	}
	
	return addons
}

func loadAddon(source dcfs.DataSource, addonname, path string) *Addon {
	addon := NewAddon(addonname)
	LogPrintln(addonname)
	
	dirpath := path
	if path != "" {
		path += "/"
	}
	
	for _, filepath := range source.ListFiles(dirpath) {
		
		LogPrintln("  " + path + filepath)

		file := new(AddonFile)
		file.Path = filepath

		content, err := source.OpenAndRead(path + filepath)
		if err != nil {
			panic(err)
		}
		file.Content = content

		basename := classifyFile(file, filepath)
		addon.Files[basename] = file
	}
	return addon
}

func classifyFile(file *AddonFile, filename string) string {
	if strings.HasSuffix(filename, ".pre.rsf") || strings.HasSuffix(filename, ".pre.rbf") {
		// is pre script
		file.PreScript = true
		return filename
	}
	if strings.HasSuffix(filename, ".post.rsf") || strings.HasSuffix(filename, ".post.rbf") {
		// is post script
		file.PostScript = true
		return filename
	}

	if strings.HasSuffix(filename, ".rbl") {
		// is rubble code (do not write out after parse)
		file.NoWrite = true
		return filename
	}

	if strings.HasSuffix(filename, ".txt") {
		// is raw file
		return filename
	}

	file.UserData = true
	return filename
}

func loadInit(source dcfs.DataSource, path string) {
	dirpath := path
	if path != "" {
		path += "/"
	}
	
	for _, filepath := range source.ListFiles(dirpath) {
		if strings.HasSuffix(filepath, ".init.rsf") {
			
			content, err := source.OpenAndRead(path + filepath)
			if err != nil {
				panic(err)
			}
			ForcedInit[filepath] = content
			ForcedInitOrder = append(ForcedInitOrder, filepath)
			continue
		}
		if strings.HasSuffix(filepath, ".init.rbf") {
			
			content, err := source.OpenAndRead(path + filepath)
			if err != nil {
				panic(err)
			}
			ForcedInit[filepath] = content
			ForcedInitOrder = append(ForcedInitOrder, filepath)
		}
	}
	
	return
}

func containsParseable(source dcfs.DataSource, path string) bool {
	for _, filename := range source.ListFiles(path) {
		if strings.HasSuffix(filename, ".init.rsf") || strings.HasSuffix(filename, ".init.rbf") {
			continue
		}
		if strings.HasSuffix(filename, ".rsf") {
			// is script
			return true
		}
		if strings.HasSuffix(filename, ".rbf") {
			// is binary script
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
	LogPrintln("Updating the Addon List File...")

	out := make([]byte, 0, 2048)
	out = append(out, "[addons]\n"...)
	for i := range addons {
		out = append(out, addons[i].Name+"="...)
		if addons[i].Active {
			out = append(out, "true\n"...)
		} else {
			out = append(out, "false\n"...)
		}
	}

	ioutil.WriteFile(dest+"/addonlist.ini", out, 0600)
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
