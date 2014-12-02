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

import gopath "path"
import "path/filepath"
import "os"
import "io"
import "io/ioutil"

// OSDir provides an interface to an OS directory and its children.
type OSDir struct {
	path string
}

// NewOSDir creates a new OS AXIS directory interface.
func NewOSDir(path string) DataSource {
	dir := new(OSDir)
	dir.path = filepath.ToSlash(path)
	return dir
}

func (dir *OSDir) IsDir(path *Path) bool {
	info, err := os.Lstat(dir.path + "/" + path.Remainder())
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (dir *OSDir) Exists(path *Path) bool {
	_, err := os.Lstat(dir.path + "/" + path.Remainder())
	if err != nil {
		return false
	}
	return true
}

func (dir *OSDir) Delete(path *Path) error {
	return os.Remove(dir.path + "/" + path.Remainder())
}

func (dir *OSDir) Create(path *Path) error {
	err := os.MkdirAll(gopath.Dir(dir.path + "/" + path.Remainder()), 0666)
	if err != nil {
		return err
	}
	f, err := os.Create(dir.path + "/" + path.Remainder())
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func (dir *OSDir) ListDir(path *Path) []string {
	files, err := ioutil.ReadDir(dir.path + "/" + path.Remainder())
	if err != nil {
		return make([]string, 0)
	}
	
	rtn := make([]string, 0, 10)
	for _, file := range files {
		if file.IsDir() {
			rtn = append(rtn, file.Name())
		}
	}
	return rtn
}

func (dir *OSDir) ListFile(path *Path) []string {
	files, err := ioutil.ReadDir(dir.path + "/" + path.Remainder())
	if err != nil {
		return make([]string, 0)
	}

	rtn := make([]string, 0, 20)
	for _, file := range files {
		if !file.IsDir() {
			rtn = append(rtn, file.Name())
		}
	}
	return rtn
}

func (dir *OSDir) Read(path *Path) (io.ReadCloser, error) {
	return os.Open(dir.path + "/" + path.Remainder())
}

func (dir *OSDir) ReadAll(path *Path) ([]byte, error) {
	reader, err := dir.Read(path)
	defer reader.Close()
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(reader)
	return content, err
}

func (dir *OSDir) Write(path *Path) (io.WriteCloser, error) {
	return os.Create(dir.path + "/" + path.Remainder())
}

func (dir *OSDir) WriteAll(path *Path, content []byte) error {
	writer, err := dir.Write(path)
	defer writer.Close()
	if err != nil {
		return err
	}
	_, err = writer.Write(content)
	return err
}
