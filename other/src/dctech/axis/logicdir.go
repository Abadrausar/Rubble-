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

package axis

import "io"
import "sort"

// LogicalDir is a simple logical directory Collection that serves it's mounted DataSources as files or child directories.
type LogicalDir map[string][]DataSource

// NewOSFile creates a new logical AXIS directory.
func NewLogicalDir() Collection {
	return make(LogicalDir)
}

func (dir LogicalDir) Mount(loc string, ds DataSource) {
	dir[loc] = append(dir[loc], ds)
}

func (dir LogicalDir) GetChild(path string) (DataSource, error) {
	path, err := Sanitize(path)
	if err != nil {
		return nil, NewError(ErrBadPath, path)
	}
	
	locID, path := StripDir(path)
	if locID == "" {
		return dir, nil
	}
	
	for x := range dir[locID] {
		ds := dir[locID][x]
		if col, ok := ds.(Collection); ok {
			return col.GetChild(path)
		}
	}
	return nil, NewError(ErrNotFound, path)
}

func (dir LogicalDir) IsDir(path string) bool {
	path, err := Sanitize(path)
	if err != nil {
		return false
	}
	
	if path == "" {
		return true
	}
	
	locID, path := StripDir(path)
	if path == "" {
		return true
	}
	
	for x := range dir[locID] {
		ds := dir[locID][x]
		if ds.Exists(path) {
			return ds.IsDir(path)
		}
	}
	return false
}

func (dir LogicalDir) Exists(path string) bool {
	path, err := Sanitize(path)
	if err != nil {
		return false
	}
	
	if path == "" {
		return true
	}
	
	if _, ok := dir[path]; ok {
		return true
	}
	
	locID, path := StripDir(path)
	if path == "" {
		return false
	}
	
	for x := range dir[locID] {
		ds := dir[locID][x]
		if ds.Exists(path) {
			return true
		}
	}
	return false
}

func (dir LogicalDir) Delete(path string) error {
	path, err := Sanitize(path)
	if err != nil {
		return NewError(ErrBadPath, path)
	}
	
	locID, path := StripDir(path)
	if locID == "" {
		return NewError(ErrBadPath, path)
	}
	
	for x := range dir[locID] {
		ds := dir[locID][x]
		if ds.Exists(path) {
			return ds.Delete(path)
		}
	}
	return NewError(ErrNotFound, path)
}

func (dir LogicalDir) Create(path string) error {
	path, err := Sanitize(path)
	if err != nil {
		return NewError(ErrBadPath, path)
	}
	
	locID, path := StripDir(path)
	if locID == "" {
		return NewError(ErrBadPath, path)
	}
	
	// This should try all mounted locations.
	if _, ok := dir[locID]; ok && len(dir[locID]) > 0 {
		return dir[locID][0].Create(path)
	}
	return NewError(ErrReadOnly, path)
}

func (dir LogicalDir) ListDir(path string) []string {
	path, err := Sanitize(path)
	if err != nil {
		return []string{}
	}
	
	locID, path := StripDir(path)
	if locID == "" {
		rtn := make([]string, 0, 25)
		for item := range dir {
			rtn = append(rtn, item)
		}
		
		sort.Strings(rtn)
		return rtn
	}

	tmp := make(map[string]bool)
	for x := range dir[locID] {
		ds := dir[locID][x]
		if ds.Exists(path) {
			for _, y := range ds.ListDir(path) {
				tmp[y] = true
			}
		}
	}
	
	rtn := make([]string, len(tmp))
	i := 0
	for x := range tmp {
		rtn[i] = x
		i++
	}
	
	sort.Strings(rtn)
	return rtn
}

func (dir LogicalDir) ListFile(path string) []string {
	path, err := Sanitize(path)
	if err != nil {
		return []string{}
	}
	
	locID, path := StripDir(path)
	if locID == "" {
		return []string{}
	}
	
	tmp := make(map[string]bool, 50)
	for x := range dir[locID] {
		ds := dir[locID][x]
		if ds.Exists(path) {
			for _, y := range ds.ListFile(path) {
				tmp[y] = true
			}
		}
	}
	
	rtn := make([]string, len(tmp))
	i := 0
	for x := range tmp {
		rtn[i] = x
		i++
	}
	
	sort.Strings(rtn)
	return rtn
}

func (dir LogicalDir) Read(path string) (io.ReadCloser, error) {
	path, err := Sanitize(path)
	if err != nil {
		return nil, NewError(ErrBadPath, path)
	}
	
	locID, path := StripDir(path)
	if path == "" || locID == "" {
		return nil, NewError(ErrBadPath, path)
	}
	
	for x := range dir[locID] {
		ds := dir[locID][x]
		if ds.Exists(path) {
			return ds.Read(path)
		}
	}
	return nil, NewError(ErrNotFound, path)
}

func (dir LogicalDir) ReadAll(path string) ([]byte, error) {
	path, err := Sanitize(path)
	if err != nil {
		return nil, NewError(ErrBadPath, path)
	}
	
	locID, path := StripDir(path)
	if path == "" || locID == "" {
		return nil, NewError(ErrBadPath, path)
	}
	
	for x := range dir[locID] {
		ds := dir[locID][x]
		if ds.Exists(path) {
			return ds.ReadAll(path)
		}
	}
	return nil, NewError(ErrNotFound, path)
}

func (dir LogicalDir) Write(path string) (io.WriteCloser, error) {
	path, err := Sanitize(path)
	if err != nil {
		return nil, NewError(ErrBadPath, path)
	}
	
	locID, path := StripDir(path)
	if path == "" || locID == "" {
		return nil, NewError(ErrBadPath, path)
	}
	
	for x := range dir[locID] {
		ds := dir[locID][x]
		if ds.Exists(path) {
			return ds.Write(path)
		}
	}
	return nil, NewError(ErrNotFound, path)
}

func (dir LogicalDir) WriteAll(path string, content []byte) error {
	path, err := Sanitize(path)
	if err != nil {
		return NewError(ErrBadPath, path)
	}
	
	locID, path := StripDir(path)
	if path == "" || locID == "" {
		return NewError(ErrBadPath, path)
	}
	
	for x := range dir[locID] {
		ds := dir[locID][x]
		if ds.Exists(path) {
			return ds.WriteAll(path, content)
		}
	}
	return NewError(ErrNotFound, path)
}
