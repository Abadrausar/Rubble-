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

// PrepModeRun runs prep mode for the specified region.
// state.Load does NOT need to be run!
func (state *State) PrepModeRun(region string) (err error) {
	defer state.trapError(&err)

	state.Log.Separator()
	state.Log.Println("Entering Prep Mode for Region:", region)

	path := state.DFDir
	if region == "raw" {
		// Main raw folder
		path += "/raw/prep"
	} else {
		path += "/data/save/" + region + "/raw/prep"
	}

	fs := axis.NewOSDir(path)
	state.FS.Mount("prep", fs)

	state.Log.Println("  Loading Prep Files...")
	state.loadAddon(region+"/Prep", "prep:")
	prep := state.AddonsTbl[region+"/Prep"]

	// A hack, but a hack that keeps everything working correctly.
	prep.Active = true
	files := NewFileList([]*Addon{prep})

	state.Log.Println("  Running Prep Scripts...")
	for _, i := range files.Order {
		if files.Data[i].Tags["Skip"] || !files.Data[i].Tags["PrepFile"] {
			continue
		}

		state.CurrentFile = files.Data[i].Name
		state.Log.Println("    " + files.Data[i].Name)

		_, err := state.ScriptState.CompileAndRun(string(files.Data[i].Content), files.Data[i].Name)
		if err != nil {
			state.CurrentFile = ""
			return err
		}
	}

	state.CurrentFile = ""
	return nil
}
