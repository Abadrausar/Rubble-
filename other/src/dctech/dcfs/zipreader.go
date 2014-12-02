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

package dcfs

import "io"
import "io/ioutil"
import "errors"
import "strings"
import "archive/zip"
import "path"
import "path/filepath"

type ZipReader struct {
	zip *zip.ReadCloser
}

func NewZipReader(path string) (DataSource, error) {
	this := new(ZipReader)
	var err error
	this.zip, err = zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	return this, nil
}

func (this *ZipReader) Open(path string) (io.ReadCloser, error) {
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
	return nil, errors.New("File not found: " + path)
}

func (this *ZipReader) OpenAndRead(path string) ([]byte, error) {
	reader, err := this.Open(path)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	reader.Close()
	return content, nil
}

func (this *ZipReader) ListDirs(dir string) []string {
	rtn := make([]string, 0, 10)
	dir = filepath.ToSlash(dir)
	if dir == "" {
		dir = "."
	}
	
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

func (this *ZipReader) ListFiles(dir string) []string {
	dir = filepath.ToSlash(dir)
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
