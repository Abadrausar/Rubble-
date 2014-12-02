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
import "dctech/axis/axiszip"
import "path/filepath"

func InstallModeRun(pack string) {
	path := pack
	pack = StripExt(filepath.Base(pack))
	
	LogPrintln("=============================================")
	LogPrintln("Installing Package: ", pack)

	instFS := axis.NewFileSystem()
	
	fs, err := axis.NewOSFile(path, false)
	if err != nil {
		panic(err)
	}
	instFS.Mount("raw", fs)
	
	fs, err = axiszip.New(path)
	if err != nil {
		panic(err)
	}
	instFS.Mount("files", fs)
	FS.Mount("install", instFS)

	LogPrintln("Loading Install Files...")
	install := loadInst(pack)

	// A hack, but a hack that keeps everything working correctly.
	install.Active = true
	Files = NewFileList([]*Addon{install})

	LogPrintln("Running Install Scripts...")
	for _, i := range Files.Order {
		if Files.Files[i].Skip || !Files.Files[i].InstScript {
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

func loadInst(pack string) *Addon {
	addon := NewAddon(pack + "/Install")

	for _, filepath := range FS.ListFile("install:files:") {

		file := new(AddonFile)
		file.Path = filepath

		content, err := FS.ReadAll("install:files:" + filepath)
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