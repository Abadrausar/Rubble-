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

import "sort"

// All the addons are loaded to here.
var Addons []*Addon

// Files is filled with all the active addon files at the end of the addon activation stage.
// Files is mapped to the Raptor indexable rubble:raws.
var Files *FileList

// Parse stage constants
const (
	stgUseCurrent = iota // Order the parser to use the current stage
	stgLoadAndInit
	stgPreScripts
	stgPreParse
	stgParse
	stgPostParse
	stgGlobalExpand
	stgPostScripts
)

// The current parse stage.
var ParseStage = stgLoadAndInit

// This is where template variables and config options are stored.
var VariableData = make(map[string]string)

// The parameters of the last called template.
var PrevParams = make([]string, 0)

// Used by the error handler and lexer.
var LastLine = 1

// AddonFile represents any file in an addon.
type AddonFile struct {
	Path       string
	Content    []byte
	PreScript  bool
	PostScript bool
	UserData   bool // any unparsable file, eg. images for tileset addons.
	Skip       bool
	NoWrite    bool
}

// IsRaw returns true is the addon file is a raw file or rbl file.
func (this *AddonFile) IsRaw() bool {
	return !this.PreScript && !this.PostScript && !this.UserData
}

// Addon represents an addon.
type Addon struct {
	Name   string
	Active bool
	Files  map[string]*AddonFile
}

// NewAddon creates a new addon with the specified name.
func NewAddon(name string) *Addon {
	this := new(Addon)
	this.Name = name
	this.Files = make(map[string]*AddonFile)
	return this
}

// FileList stores an ordered list of AddonFiles.
type FileList struct {
	Order []string
	Files map[string]*AddonFile
}

// NewFileList creates a FileList from the files of all the active addons passed in.
// Inactive addons are ignored.
func NewFileList(data []*Addon) *FileList {
	this := new(FileList)
	this.Order = make([]string, 0, 100)
	this.Files = make(map[string]*AddonFile, 100)

	for _, addon := range data {
		if addon.Active {
			for name, file := range addon.Files {
				this.Order = append(this.Order, name)
				this.Files[name] = file
			}
		}
	}

	sort.Strings(this.Order)

	return this
}