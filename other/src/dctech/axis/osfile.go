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

import "path/filepath"
import "os"
import "io"
import "io/ioutil"

// OSFile allows you to provide an interface to a single OS file.
type OSFile struct {
	path string
	rw bool
}

// NewOSFile creates a new OS AXIS file interface.
// rw controls if Delete, Write, and WriteAll work.
func NewOSFile(path string, rw bool) DataSource {
	file := new(OSFile)
	file.rw = rw
	file.path = filepath.ToSlash(path)
	return file
}

func (file *OSFile) IsDir(path *Path) bool {
	return false
}

func (file *OSFile) Exists(path *Path) bool {
	if !path.Done() {
		return false
	}
	
	info, err := os.Lstat(file.path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}

func (file *OSFile) Delete(path *Path) error {
	if !path.Done() {
		return NewError(ErrNotFound, path)
	}
	
	if !file.rw {
		return NewError(ErrReadOnly, path)
	}
	return os.Remove(file.path)
}

func (file *OSFile) Create(path *Path) error {
	if !path.Done() {
		return NewError(ErrNotFound, path)
	}
	
	if !file.Exists(path) {
		f, err := os.Create(file.path)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

func (file *OSFile) ListDir(path *Path) []string {
	return []string{}
}

func (file *OSFile) ListFile(path *Path) []string {
	return []string{}
}

func (file *OSFile) Read(path *Path) (io.ReadCloser, error) {
	if !path.Done() {
		return nil, NewError(ErrNotFound, path)
	}
	
	return os.Open(file.path)
}

func (file *OSFile) ReadAll(path *Path) ([]byte, error) {
	if !path.Done() {
		return nil, NewError(ErrNotFound, path)
	}
	
	reader, err := file.Read(path)
	defer reader.Close()
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(reader)
	return content, err
}

func (file *OSFile) Write(path *Path) (io.WriteCloser, error) {
	if !path.Done() {
		return nil, NewError(ErrNotFound, path)
	}
	if !file.rw {
		return nil, NewError(ErrReadOnly, path)
	}
	
	return os.Create(file.path)
}

func (file *OSFile) WriteAll(path *Path, content []byte) error {
	if !path.Done() {
		return NewError(ErrNotFound, path)
	}
	if !file.rw {
		return NewError(ErrReadOnly, path)
	}
	
	writer, err := file.Write(path)
	defer writer.Close()
	if err != nil {
		return err
	}
	_, err = writer.Write(content)
	return err
}
