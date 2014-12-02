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

// AXIS VFS, a simple virtual file system API.
// Absurdly eXtremely Incredibly Simple Virtual File System (adjectives are good for making cool acronyms!).
// 
// AXIS uses special OS-independent two part paths, with the first part being a series of "location IDs"
// separated by colons and the second part being a slash separated series of directory/file names.
// 
// The root DataSource is accessed with the path "".
// 
// If you think the name is stupid (it is) you can just call it AXIS and forget what it is supposed to mean,
// after all the "official" name is more of a joke than anything...
// 
package axis

import "io"

// locid1:locid2:dir1/dir2/dir3/file.ext

// DataSource is the heart of AXIS. EVERYTHING in AXIS implements DataSource, and it is the sole interface to all AXIS data.
type DataSource interface {
	// IsDir returns true if the path points to an AXIS directory
	// (this may not be a real directory on the OS file system! AXIS directories can be purely virtual).
	IsDir(path *Path) bool
	
	// Exists returns true if an AXIS path is valid.
	Exists(path *Path) bool
	
	// Delete tries to delete an AXIS resource, this generally removes the backing file/directory from the OS's file system.
	// Read-only resources generally cannot be deleted.
	Delete(path *Path) error
	
	// Directories
	
	// Create tries to create an AXIS resource, this generally creates the backing file/directory on the OS's file system.
	// If the resource already exists then this function should do nothing.
	Create(path *Path) error
	
	// ListDir lists the AXIS directories at a path.
	// Returns a zero length slice if there are no directories or anything went wrong.
	ListDir(path *Path) []string
	
	// ListDir lists the files in a AXIS directory at a path.
	// Returns a zero length slice if there are no files or anything went wrong.
	ListFile(path *Path) []string
	
	// Files
	
	// Read opens an AXIS file for reading and returns the result and any error that may have happened.
	Read(path *Path) (io.ReadCloser, error)
	
	// ReadAll attempts to read an AXIS file into memory and returns the result and any error that may have happened.
	ReadAll(path *Path) ([]byte, error)
	
	// Write opens an AXIS file for writing and returns the result and any error that may have happened.
	// The file is always truncated!
	Write(path *Path) (io.WriteCloser, error)
	
	// WriteAll attempts to replace the contents of an AXIS file with the data given, if an error happened it is returned.
	WriteAll(path *Path, content []byte) error
}

// Collection is used for things like logical directories and the base file system.
type Collection interface {
	DataSource
	
	// Mount a DataSource at the specified location, the name may be used as a location ID or
	// file/directory name, depending on context.
	Mount(name string, ds DataSource)
	
	// GetChild attempts to lookup a child DataSource by path.
	GetChild(path *Path) (DataSource, error)
}

// ================================================================
// Convenience functions
// ================================================================

// IsDir returns true if the path points to an AXIS directory
// (this may not be a real directory on the OS file system! AXIS directories can be purely virtual).
func IsDir(ds DataSource, path string) bool {
	return ds.IsDir(NewPath(path))
}

// Exists returns true if an AXIS path is valid.
func Exists(ds DataSource, path string) bool {
	return ds.Exists(NewPath(path))
}

// Delete tries to delete an AXIS resource, this generally removes the backing file/directory from the OS's file system.
// Read-only resources generally cannot be deleted.
func Delete(ds DataSource, path string) error {
	return ds.Delete(NewPath(path))
}

// Directories

// Create tries to create an AXIS resource, this generally creates the backing file/directory on the OS's file system.
// If the resource already exists then this function should do nothing.
func Create(ds DataSource, path string) error {
	return ds.Create(NewPath(path))
}

// ListDir lists the AXIS directories at a path.
// Returns a zero length slice if there are no directories or anything went wrong.
func ListDir(ds DataSource, path string) []string {
	return ds.ListDir(NewPath(path))
}

// ListDir lists the files in a AXIS directory at a path.
// Returns a zero length slice if there are no files or anything went wrong.
func ListFile(ds DataSource, path string) []string {
	return ds.ListFile(NewPath(path))
}

// Files

// Read opens an AXIS file for reading and returns the result and any error that may have happened.
func Read(ds DataSource, path string) (io.ReadCloser, error) {
	return ds.Read(NewPath(path))
}

// ReadAll attempts to read an AXIS file into memory and returns the result and any error that may have happened.
func ReadAll(ds DataSource, path string) ([]byte, error) {
	return ds.ReadAll(NewPath(path))
}

// Write opens an AXIS file for writing and returns the result and any error that may have happened.
// The file is always truncated!
func Write(ds DataSource, path string) (io.WriteCloser, error) {
	return ds.Write(NewPath(path))
}

// WriteAll attempts to replace the contents of an AXIS file with the data given, if an error happened it is returned.
func WriteAll(ds DataSource, path string, content []byte) error {
	return ds.WriteAll(NewPath(path), content)
}

// Mount a DataSource at the specified location, the name may be used as a location ID or
// file/directory name, depending on context.
func Mount(col Collection, name string, ds DataSource) {
	col.Mount(name, ds)
}

// GetChild attempts to lookup a child DataSource by path.
func GetChild(col Collection, path string) (DataSource, error) {
	return col.GetChild(NewPath(path))
}
