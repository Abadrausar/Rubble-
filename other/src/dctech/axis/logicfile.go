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

import "bytes"
import "io"

// LogicalFile is a simple logical file, basically a simple interface to a byte buffer.
type LogicalFile struct {
	rw bool
	data *byteBuffer
}

// NewLogicalFile creates a new logical AXIS file.
// rw controls if Delete, Write, and WriteAll work, it is still possible to modify a file
// by abusing Read and ReadAll (but don't do that!).
func NewLogicalFile(data []byte, rw bool) DataSource {
	return &LogicalFile{
		rw: rw,
		data: newByteBuffer(data),
	}
}

func (file *LogicalFile) IsDir(path *Path) bool {
	return false
}

func (file *LogicalFile) Exists(path *Path) bool {
	if path.Done() {
		return true
	}
	return false
}

// Delete works a little differently from normal here, it blanks the files contents rather than deleting it entire
// (which would be impossible).
func (file *LogicalFile) Delete(path *Path) error {
	if !path.Done() {
		return NewError(ErrNotFound, path)
	}
	
	if !file.rw {
		return NewError(ErrReadOnly, path)
	}
	file.data.Reset()
	return nil
}

func (file *LogicalFile) Create(path *Path) error {
	if !path.Done() {
		return NewError(ErrNotFound, path)
	}
	return nil
}

func (file *LogicalFile) ListDir(path *Path) []string {
	return []string{}
}

func (file *LogicalFile) ListFile(path *Path) []string {
	return []string{}
}

func (file *LogicalFile) Read(path *Path) (io.ReadCloser, error) {
	if !path.Done() {
		return nil, NewError(ErrNotFound, path)
	}
	
	return file.data, nil
}

func (file *LogicalFile) ReadAll(path *Path) ([]byte, error) {
	if !path.Done() {
		return nil, NewError(ErrNotFound, path)
	}
	
	return file.data.Bytes(), nil
}

func (file *LogicalFile) Write(path *Path) (io.WriteCloser, error) {
	if !path.Done() {
		return nil, NewError(ErrNotFound, path)
	}
	if !file.rw {
		return nil, NewError(ErrReadOnly, path)
	}
	
	file.data.Reset()
	return file.data, nil
}

func (file *LogicalFile) WriteAll(path *Path, content []byte) error {
	if !path.Done() {
		return NewError(ErrNotFound, path)
	}
	if !file.rw {
		return NewError(ErrReadOnly, path)
	}
	
	file.data = newByteBuffer(content)
	return nil
}

// Internals

type byteBuffer struct {
	*bytes.Buffer
}

func newByteBuffer(data []byte) *byteBuffer {
	return &byteBuffer{
		bytes.NewBuffer(data),
	}
}

func (buff *byteBuffer) Close() error {
	return nil
}
