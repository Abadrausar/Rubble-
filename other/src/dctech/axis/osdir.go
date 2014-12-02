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
func NewOSDir(path string) (DataSource, error) {
	dir := new(OSDir)
	dir.path = filepath.ToSlash(path)
	
	info, err := os.Lstat(dir.path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, NewError(ErrNotDir, path)
	}
	return dir, nil
}

func (dir *OSDir) IsDir(path string) bool {
	path, err := Sanitize(path)
	if err != nil {
		return false
	}
	
	info, err := os.Lstat(dir.path + "/" + path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (dir *OSDir) Exists(path string) bool {
	path, err := Sanitize(path)
	if err != nil {
		return false
	}
	
	_, err = os.Lstat(dir.path + "/" + path)
	if err != nil {
		return false
	}
	return true
}

func (dir *OSDir) Delete(path string) error {
	path, err := Sanitize(path)
	if err != nil {
		return err
	}
	
	return os.Remove(dir.path + "/" + path)
}

func (dir *OSDir) Create(path string) error {
	path, err := Sanitize(path)
	if err != nil {
		return err
	}
	
	err = os.MkdirAll(gopath.Dir(dir.path + "/" + path), 0666)
	if err != nil {
		return err
	}
	f, err := os.Create(dir.path + "/" + path)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func (dir *OSDir) ListDir(path string) []string {
	path, err := Sanitize(path)
	if err != nil {
		return make([]string, 0)
	}
	if path == "" {
		path = "."
	}

	files, err := ioutil.ReadDir(dir.path + "/" + path)
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

func (dir *OSDir) ListFile(path string) []string {
	path, err := Sanitize(path)
	if err != nil {
		return make([]string, 0)
	}
	if path == "" {
		path = "."
	}

	files, err := ioutil.ReadDir(dir.path + "/" + path)
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

func (dir *OSDir) Read(path string) (io.ReadCloser, error) {
	path, err := Sanitize(path)
	if err != nil {
		return nil, err
	}
	
	return os.Open(dir.path + "/" + path)
}

func (dir *OSDir) ReadAll(path string) ([]byte, error) {
	path, err := Sanitize(path)
	if err != nil {
		return nil, err
	}
	
	reader, err := dir.Read(path)
	defer reader.Close()
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(reader)
	return content, err
}

func (dir *OSDir) Write(path string) (io.WriteCloser, error) {
	path, err := Sanitize(path)
	if err != nil {
		return nil, err
	}
	
	return os.Create(dir.path + "/" + path)
}

func (dir *OSDir) WriteAll(path string, content []byte) error {
	path, err := Sanitize(path)
	if err != nil {
		return err
	}
	
	writer, err := dir.Write(path)
	defer writer.Close()
	if err != nil {
		return err
	}
	_, err = writer.Write(content)
	return err
}
