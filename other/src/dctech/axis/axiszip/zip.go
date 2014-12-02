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

package axiszip

import "dctech/axis"

import "io"
import "io/ioutil"
import "strings"
import "archive/zip"
import "path"
import "path/filepath"

type Data struct {
	zip *zip.ReadCloser
}

// New creates a read-only AXIS DataSource backed by a zip file.
func New(path string) (axis.DataSource, error) {
	this := new(Data)
	var err error
	this.zip, err = zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	return this, nil
}

func (this *Data) IsDir(path string) bool {
	path = strings.TrimRight(path, "/")
	if path == "" {
		return true
	}
	path = filepath.ToSlash(path)
	
	for _, f := range this.zip.File {
		if strings.TrimRight(f.Name, "/") == path {
			return f.FileInfo().IsDir()
		}
	}
	return false
}

func (this *Data) Exists(path string) bool {
	path = strings.TrimRight(path, "/")
	if path == "" {
		return true
	}
	path = filepath.ToSlash(path)
	
	for _, f := range this.zip.File {
		if strings.TrimRight(f.Name, "/") == path {
			return true
		}
	}
	return false
}

func (this *Data) Delete(path string) error {
	return axis.NewError(axis.ErrReadOnly, path)
}

func (this *Data) Create(path string) error {
	return axis.NewError(axis.ErrReadOnly, path)
}

func (this *Data) ListDir(dir string) []string {
	dir = filepath.ToSlash(dir)
	dir = strings.TrimRight(dir, "/")
	if dir == "" {
		dir = "."
	}
	rtn := make([]string, 0, 10)

	for _, f := range this.zip.File {
		if f.FileInfo().IsDir() {
			name := strings.TrimRight(f.Name, "/")
			if path.Dir(name) == dir {
				rtn = append(rtn, path.Base(name))
			}
		}
	}
	return rtn
}

func (this *Data) ListFile(dir string) []string {
	dir = filepath.ToSlash(dir)
	dir = strings.TrimRight(dir, "/")
	if dir == "" {
		dir = "."
	}
	rtn := make([]string, 0, 20)

	for _, f := range this.zip.File {
		if !f.FileInfo().IsDir() {
			if path.Dir(f.Name) == dir {
				rtn = append(rtn, path.Base(f.Name))
			}
		}
	}
	return rtn
}

func (this *Data) Read(path string) (io.ReadCloser, error) {
	path = filepath.ToSlash(path)
	for _, f := range this.zip.File {
		if f.Name == path {
			reader, err := f.Open()
			if err != nil {
				return nil, err
			}
			return reader, nil
		}
	}
	return nil, axis.NewError(axis.ErrNotFound, path)
}

func (this *Data) ReadAll(path string) ([]byte, error) {
	path = filepath.ToSlash(path)
	reader, err := this.Read(path)
	defer reader.Close()
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(reader)
	return content, err
}

func (this *Data) Write(path string) (io.WriteCloser, error) {
	return nil, axis.NewError(axis.ErrReadOnly, path)
}

func (this *Data) WriteAll(path string, content []byte) error {
	return axis.NewError(axis.ErrReadOnly, path)
}