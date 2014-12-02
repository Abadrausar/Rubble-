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

// LogicalDir is a simple logical directory Collection that serves it's mounted DataSources as files or child directories.
type LogicalDir map[string]DataSource

// NewOSFile creates a new logical AXIS directory.
func NewLogicalDir() Collection {
	return make(LogicalDir)
}

func (dir LogicalDir) Mount(loc string, ds DataSource) {
	dir[loc] = ds
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
		return false
	}
	
	if dir[locID].Exists(path) {
		return dir[locID].Exists(path)
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
	
	if dir[locID].Exists(path) {
		return true
	}
	return false
}

func (dir LogicalDir) Delete(path string) error {
	locID, path := StripDir(path)
	if locID == "" {
		return NewError(ErrBadPath, path)
	}
	
	if dir[locID].Exists(path) {
		return dir[locID].Delete(path)
	}
	return NewError(ErrNotFound, path)
}

func (dir LogicalDir) Create(path string) error {
	locID, path := StripDir(path)
	if locID == "" {
		return NewError(ErrBadPath, path)
	}
	
	if dir[locID].Exists(path) {
		return dir[locID].Create(path)
	}
	return NewError(ErrReadOnly, path)
}

func (dir LogicalDir) ListDir(path string) []string {
	locID, path := StripDir(path)
	if locID == "" {
		rtn := make([]string, 0, 25)
		for item := range dir {
			rtn = append(rtn, item)
		}
		return rtn
	}

	rtn := make([]string, 0, 25)
	if dir[locID].Exists(path) {
		for _, item := range dir[locID].ListDir(path) {
			rtn = append(rtn, item)
		}
	}
	return rtn
}

func (dir LogicalDir) ListFile(path string) []string {
	locID, path := StripDir(path)
	if locID == "" {
		return make([]string, 0)
	}
	
	rtn := make([]string, 0, 25)
	if dir[locID].Exists(path) {
		for _, item := range dir[locID].ListFile(path) {
			rtn = append(rtn, item)
		}
	}
	return rtn
}

func (dir LogicalDir) Read(path string) (io.ReadCloser, error) {
	locID, path := StripDir(path)
	if path == "" || locID == "" {
		return nil, NewError(ErrBadPath, path)
	}
	
	if dir[locID].Exists(path) {
		return dir[locID].Read(path)
	}
	return nil, NewError(ErrNotFound, path)
}

func (dir LogicalDir) ReadAll(path string) ([]byte, error) {
	locID, path := StripDir(path)
	if path == "" || locID == "" {
		return nil, NewError(ErrBadPath, path)
	}
	
	if dir[locID].Exists(path) {
		return dir[locID].ReadAll(path)
	}
	return nil, NewError(ErrNotFound, path)
}

func (dir LogicalDir) Write(path string) (io.WriteCloser, error) {
	locID, path := StripDir(path)
	if path == "" || locID == "" {
		return nil, NewError(ErrBadPath, path)
	}
	
	return dir[locID].Write(path)
}

func (dir LogicalDir) WriteAll(path string, content []byte) error {
	locID, path := StripDir(path)
	if path == "" || locID == "" {
		return NewError(ErrBadPath, path)
	}
	
	return dir[locID].WriteAll(path, content)
}
