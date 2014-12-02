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

import "sort"

// AddonFile represents any file in an addon.
type AddonFile struct {
	Name    string // File name.
	Source  string // Addon path (AXIS syntax, including loc ids).
	Content []byte
	Tags    map[string]bool
}

// NewAddonFile creates a new addon file with the specified name.
func NewAddonFile(name, source string, content []byte) *AddonFile {
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
	Meta   *Meta
	Files  map[string]*AddonFile
}

// NewAddon creates a new addon with the specified name.
func NewAddon(name string) *Addon {
	this := new(Addon)
	this.Name = name
	this.Meta = NewMeta()
	this.Files = make(map[string]*AddonFile)
	return this
}

// FileList stores an ordered list of AddonFiles.
type FileList struct {
	Order []string
	Data  map[string]*AddonFile
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

// NewFileListFunc creates a FileList from the files of all the active addons passed in.
// Inactive addons are ignored. If filter returns false for a file it is not included.
// If the addon list is nil an empty (but valid) FileList is returned.
func NewFileListFunc(data []*Addon, filter func(*AddonFile) bool) *FileList {
	list := new(FileList)
	list.Order = make([]string, 0, 100)
	list.Data = make(map[string]*AddonFile, 100)

	if data == nil {
		return list
	}

	for _, addon := range data {
		if addon.Active {
			for name, file := range addon.Files {
				if filter(file) {
					if _, ok := list.Data[name]; !ok {
						list.Order = append(list.Order, name)
					}
					list.Data[name] = file
				}
			}
		}
	}

	sort.Strings(list.Order)

	return list
}

func (list *FileList) AddFiles(files ...*AddonFile) {
	for _, file := range files {
		if _, ok := list.Data[file.Name]; !ok {
			list.Order = append(list.Order, file.Name)
		}
		list.Data[file.Name] = file
	}
	sort.Strings(list.Order)
}
