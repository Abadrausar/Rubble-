/*
Copyright 2013 by Milo Christiansen

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

	source, err := dcfs.NewDirReader(path)
	if err != nil {
		panic(err)
	}

	for _, dir := range source.ListDirs(".") {
		addonlist = append(addonlist, loadAddon(source, dir, dir))
	}

	for _, file := range source.ListFiles(".") {
		if strings.HasSuffix(file, ".pack.zip") {
			zip, err := dcfs.NewZipReader(path + "/" + file)
			if err != nil {
				panic(err)
			}

			for _, dir := range zip.ListDirs(".") {
				addonlist = append(addonlist, loadAddon(zip, StripExt(StripExt(file))+"/"+dir, dir))
			}
			continue
		}

		if strings.HasSuffix(file, ".zip") {
			zip, err := dcfs.NewZipReader(path + "/" + file)
			if err != nil {
				panic(err)
			}
			addonlist = append(addonlist, loadAddon(zip, StripExt(file), "."))
			continue
		}
	}

	return addonlist
}

func loadAddon(source dcfs.DataSource, addonname, path string) *Addon {
	addon := NewAddon(addonname)
	LogPrintln(addonname)
	
	dirpath := path
	if path == "." {
		path = ""
	} else {
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
		return StripExt(StripExt(filename))
	}
	if strings.HasSuffix(filename, ".post.rsf") || strings.HasSuffix(filename, ".post.rbf") {
		// is post script
		file.PostScript = true
		return StripExt(StripExt(filename))
	}

	if strings.HasSuffix(filename, ".rbl") {
		// is rubble code (do not write out after parse)
		file.NoWrite = true
		return StripExt(filename)
	}

	if strings.HasSuffix(filename, ".txt") {
		// is raw file
		return StripExt(filename)
	}

	file.UserData = true
	return filename
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
