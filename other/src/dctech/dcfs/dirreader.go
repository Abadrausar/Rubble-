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

import "path/filepath"
import "os"
import "io"
import "io/ioutil"
import "errors"

type DirReader struct {
	path string
}

func NewDirReader(path string) (DataSource, error) {
	this := new(DirReader)
	this.path = filepath.ToSlash(path)
	info, err := os.Lstat(this.path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, errors.New("File is not a directory: " + this.path)
	}
	return this, nil
}

func (this *DirReader) Open(path string) (io.ReadCloser, error) {
	path = filepath.ToSlash(path)
	file, err := os.Open(this.path + "/" + path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (this *DirReader) OpenAndRead(path string) ([]byte, error) {
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

func (this *DirReader) ListDirs(dir string) []string {
	dir = filepath.ToSlash(dir)
	if dir == "" {
		dir = "."
	}
	rtn := make([]string, 0, 10)

	files, err := ioutil.ReadDir(this.path + "/" + dir)
	if err != nil {
		return rtn
	}

	for _, file := range files {
		if file.IsDir() {
			rtn = append(rtn, file.Name())
		}
	}
	return rtn
}

func (this *DirReader) ListFiles(dir string) []string {
	dir = filepath.ToSlash(dir)
	if dir == "" {
		dir = "."
	}
	rtn := make([]string, 0, 20)
	
	files, err := ioutil.ReadDir(this.path + "/" + dir)
	if err != nil {
		return rtn
	}

	for _, file := range files {
		if !file.IsDir() {
			rtn = append(rtn, file.Name())
		}
	}
	return rtn
}
