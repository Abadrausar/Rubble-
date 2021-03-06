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

import "os"
import "io"
import "io/ioutil"
import "bytes"
import "strings"
import "archive/zip"
import "path"
import "encoding/base64"

type Data struct {
	zip *zip.Reader
}

// NewFile creates a read-only AXIS DataSource backed by a zip file.
func NewFile(path string) (axis.DataSource, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	file := bytes.NewReader(content)
	this := new(Data)
	this.zip, err = zip.NewReader(file, int64(file.Len()))
	if err != nil {
		return nil, err
	}
	return this, nil
}

// NewFile64 creates a read-only AXIS DataSource backed by a zip file encoded in base 64.
func NewFile64(path string) (axis.DataSource, error) {
	content, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer content.Close()
	
	bc := base64.NewDecoder(base64.StdEncoding, content)
	data, err := ioutil.ReadAll(bc)
	if err != nil {
		return nil, err
	}
	
	file := bytes.NewReader(data)
	this := new(Data)
	this.zip, err = zip.NewReader(file, int64(file.Len()))
	if err != nil {
		return nil, err
	}
	return this, nil
}

// NewRaw creates a read-only AXIS DataSource backed by a zip file read into a byte slice.
func NewRaw(content []byte) (axis.DataSource, error) {
	file := bytes.NewReader(content)
	this := new(Data)
	var err error
	this.zip, err = zip.NewReader(file, int64(file.Len()))
	if err != nil {
		return nil, err
	}
	return this, nil
}

// NewRaw64 creates a read-only AXIS DataSource backed by a zip file read into a byte slice and encoded in base 64.
func NewRaw64(content []byte) (axis.DataSource, error) {
	b := bytes.NewReader(content)
	bc := base64.NewDecoder(base64.StdEncoding, b)
	data, err := ioutil.ReadAll(bc)
	if err != nil {
		return nil, err
	}
	
	file := bytes.NewReader(data)
	this := new(Data)
	this.zip, err = zip.NewReader(file, int64(file.Len()))
	if err != nil {
		return nil, err
	}
	return this, nil
}

func (this *Data) IsDir(path *axis.Path) bool {
	if path.Done() {
		return true
	}
	
	for _, f := range this.zip.File {
		if strings.TrimRight(f.Name, "/") == path.Remainder() {
			return f.FileInfo().IsDir()
		}
	}
	return false
}

func (this *Data) Exists(path *axis.Path) bool {
	if path.Done() {
		return true
	}
	
	for _, f := range this.zip.File {
		if strings.TrimRight(f.Name, "/") == path.Remainder() {
			return true
		}
	}
	return false
}

func (this *Data) Delete(path *axis.Path) error {
	return axis.NewError(axis.ErrReadOnly, path)
}

func (this *Data) Create(path *axis.Path) error {
	return axis.NewError(axis.ErrReadOnly, path)
}

func (this *Data) ListDir(dir *axis.Path) []string {
	rtn := make([]string, 0, 10)

	for _, f := range this.zip.File {
		if f.FileInfo().IsDir() {
			name := strings.TrimRight(f.Name, "/")
			if path.Dir(name) == dir.Remainder() {
				rtn = append(rtn, path.Base(name))
			}
		}
	}
	return rtn
}

func (this *Data) ListFile(dir *axis.Path) []string {
	rtn := make([]string, 0, 20)

	for _, f := range this.zip.File {
		if !f.FileInfo().IsDir() {
			if path.Dir(f.Name) == dir.Remainder() {
				rtn = append(rtn, path.Base(f.Name))
			}
		}
	}
	return rtn
}

func (this *Data) Read(path *axis.Path) (io.ReadCloser, error) {
	for _, f := range this.zip.File {
		if f.Name == path.Remainder() {
			reader, err := f.Open()
			if err != nil {
				return nil, err
			}
			return reader, nil
		}
	}
	return nil, axis.NewError(axis.ErrNotFound, path)
}

func (this *Data) ReadAll(path *axis.Path) ([]byte, error) {
	reader, err := this.Read(path)
	defer reader.Close()
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(reader)
	return content, err
}

func (this *Data) Write(path *axis.Path) (io.WriteCloser, error) {
	return nil, axis.NewError(axis.ErrReadOnly, path)
}

func (this *Data) WriteAll(path *axis.Path, content []byte) error {
	return axis.NewError(axis.ErrReadOnly, path)
}
