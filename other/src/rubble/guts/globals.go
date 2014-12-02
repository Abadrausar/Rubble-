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

import "sort"
import "dctech/axis"

// This is a list of all Rubble versions in this series
// (all the listed versions are assumed to be backwards compatible)
// Index 0 MUST be the current version!
var RubbleVersions = []string{
	"4.7",
	"4.6",
	"4.5",
	"4.4",
	"4.3",
	"4.2",
	"4.1",
	"4.0",
	"pre4",
}

var FS axis.Collection
var ExAddonsFS axis.Collection

// All the addons are loaded to here.
var Addons = make([]*Addon, 0, 64)

// Files is filled with all the active addon files at the end of the addon activation stage.
// Files is mapped to the Raptor indexable rubble:raws.
var Files *FileList

// Filled by the addonloader with all init and load scripts.
var GlobalFiles = NewFileList(nil)

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
var PrevParams = make([]*Value, 0)

// Used by the error handler and lexer.
var LastLine = NewPosition(-1, "bad_position")

// Used by rubble:currentfile
var CurrentFile string

// AddonFile represents any file in an addon.
type AddonFile struct {
	Name       string // File name.
	Source     string // Addon name (including loc ids).
	Content    []byte
	Tags       map[string]bool
}

// NewAddonFile creates a new addon with the specified name.
func NewAddonFile(name string, source string, content []byte) *AddonFile {
	file := new(AddonFile)
	file.Name = name
	file.Source = source
	file.Content = content
	file.Tags = make(map[string]bool)
	return file
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
	Data map[string]*AddonFile
}

// NewFileList creates a FileList from the files of all the active addons passed in.
// Inactive addons are ignored.
// If the addon list is nil an empty (but valid) FileList is returned.
func NewFileList(data []*Addon) *FileList {
	list := new(FileList)
	list.Order = make([]string, 0, 100)
	list.Data = make(map[string]*AddonFile, 100)

	if data == nil {
		return list
	}
	
	for _, addon := range data {
		if addon.Active {
			for name, file := range addon.Files {
				if _, ok := list.Data[name]; !ok {
					list.Order = append(list.Order, name)
				}
				list.Data[name] = file
			}
		}
	}

	sort.Strings(list.Order)

	return list
}

func DumpConfig(path string) {
	out := "\n# Rubble config variable dump.\n# Automatically generated, do not edit!\n\n[config]\n"

	for i := range VariableData {
		out += i + " = " + VariableData[i] + "\n"
	}

	WriteFile(path, []byte(out))
}

func WriteFile(path string, content []byte) {
	if !FS.Exists(path) {
		err := FS.Create(path)
		if err != nil {
			panic("Write Error: " + err.Error())
		}
	}
	
	err := FS.WriteAll(path, content)
	if err != nil {
		panic("Write Error: " + err.Error())
	}
}
