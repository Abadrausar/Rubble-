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

import "dctech/axis"

func PrepModeRun(region string) {
	LogPrintln("=============================================")
	LogPrintln("Entering Prep Mode for Region:", region)

	path := DFDir
	if region == "raw" {
		// Main raw folder
		path += "/raw/prep"
	} else {
		path += "/data/save/" + region + "/raw/prep"
	}

	fs, err := axis.NewOSDir(path)
	if err != nil {
		panic(err)
	}
	FS.Mount("prep", fs)

	LogPrintln("Loading Prep Files...")
	prep := loadPrep(region)

	// A hack, but a hack that keeps everything working correctly.
	prep.Active = true
	Files = NewFileList([]*Addon{prep})

	LogPrintln("Running Prep Scripts...")
	for _, i := range Files.Order {
		if Files.Files[i].Skip || !Files.Files[i].PrepScript {
			continue
		}

		CurrentFile = i
		LogPrintln("  " + Files.Files[i].Path)

		_, err := GlobalScriptState.CompileAndRun(string(Files.Files[i].Content), i)
		if err != nil {
			panic("Script Error: " + err.Error())
		}
	}
}

func loadPrep(region string) *Addon {
	addon := NewAddon(region + "/Prep")

	for _, filepath := range FS.ListFile("prep:") {

		//LogPrintln("  " + filepath)

		file := new(AddonFile)
		file.Path = filepath

		content, err := FS.ReadAll("prep:" + filepath)
		if err != nil {
			panic(err)
		}
		file.Content = content

		// Most of classifyFile is useless here, but there is no good reason to
		// write a subset for just prep addons
		classifyFile(file, filepath)
		addon.Files[filepath] = file
	}
	return addon
}
