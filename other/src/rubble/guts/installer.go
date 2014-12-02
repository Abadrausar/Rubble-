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
	install := loadAddon(FS, pack + "/Install", "install:files:")

	// A hack, but a hack that keeps everything working correctly.
	install.Active = true
	Files = NewFileList([]*Addon{install})

	LogPrintln("Running Install Scripts...")
	for _, i := range Files.Order {
		if Files.Data[i].Tags["Skip"] || !Files.Data[i].Tags["InstScript"] {
			continue
		}

		CurrentFile = Files.Data[i].Name
		LogPrintln("  " + Files.Data[i].Name)

		_, err := GlobalScriptState.CompileAndRun(string(Files.Data[i].Content), Files.Data[i].Name)
		if err != nil {
			panic("Script Error: " + err.Error())
		}
	}
}
