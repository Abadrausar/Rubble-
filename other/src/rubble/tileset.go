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

import "dctech/axis"

import "rubble/rblutil"

// TSetModeRun applies a tileset to the specified region.
func TSetModeRun(region, dfdir string, addonsdir []string, addons []string, log *rblutil.Logger) (err error) {
	output := dfdir
	if region == "raw" {
		output += "/raw"
	} else {
		output += "/data/save/" + region + "/raw"
	}
	
	log.Separator()
	log.Println("Entering Tileset Mode for Region:", region)
	
	oops, state := NewState(dfdir, output, addonsdir, log)
	if oops != nil {
		return oops
	}
	defer state.trapError(&err)

	err = state.Load(addons, nil)
	if err != nil {
		return err
	}
	
	state.Files = NewFileList(state.Addons)
	
	// This should keep any files from addons from messing things up.
	for _, file := range state.Files.Data {
		file.Tags["NoWrite"] = true
		file.Tags["RawFile"] = false
	}
	
	rawFiles := []*AddonFile{}
	for _, filepath := range axis.ListFile(state.FS, "out:objects") {
		content, err := axis.ReadAll(state.FS, "out:objects/" + filepath)
		if err != nil {
			return err
		}

		file := NewAddonFile(filepath, "out:objects", content)
		classifyFile(file, filepath)
		
		// All the files should be raw files, but just in case...
		if file.Tags["RawFile"] && !file.Tags["NoWrite"] && !file.Tags["TSetFile"] {
			rawFiles = append(rawFiles, file)
		}
	}
	state.Files.AddFiles(rawFiles...)
	
	err = state.ApplyTileSet()
	if err != nil {
		return err
	}
	
	err = state.Write()
	if err != nil {
		return err
	}
	
	return nil
}
