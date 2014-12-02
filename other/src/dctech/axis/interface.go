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
import "strings"
import "path/filepath"

// locid1:locid2:dir1/dir2/dir3/file.ext

// DataSource is the heart of AXIS. EVERYTHING in AXIS implements DataSource, and it is the sole interface to all AXIS data.
type DataSource interface {
	// IsDir returns true if the path points to an AXIS directory
	// (this may not be a real directory on the OS file system! AXIS directories can be purely virtual).
	IsDir(path string) bool
	
	// Exists returns true if an AXIS path is valid.
	Exists(path string) bool
	
	// Delete tries to delete an AXIS resource, this generally removes the backing file/directory from the OS's file system.
	// Read-only resources generally cannot be deleted.
	Delete(path string) error
	
	// Directories
	
	// Create tries to create an AXIS resource, this generally creates the backing file/directory on the OS's file system.
	// If the resource already exists then this function should do nothing.
	Create(path string) error
	
	// ListDir lists the AXIS directories at a path.
	// Returns a zero length slice if there are no directories or anything went wrong.
	ListDir(path string) []string
	
	// ListDir lists the files in a AXIS directory at a path.
	// Returns a zero length slice if there are no files or anything went wrong.
	ListFile(path string) []string
	
	// Files
	
	// Read opens an AXIS file for reading and returns the result and any error that may have happened.
	Read(path string) (io.ReadCloser, error)
	
	// ReadAll attempts to read an AXIS file into memory and returns the result and any error that may have happened.
	ReadAll(path string) ([]byte, error)
	
	// Write opens an AXIS file for writing and returns the result and any error that may have happened.
	// The file is always truncated!
	Write(path string) (io.WriteCloser, error)
	
	// WriteAll attempts to replace the contents of an AXIS file with the data given, if an error happened it is returned.
	WriteAll(path string, content []byte) error
}

// Collection is used for things like logical directories and the base file system.
type Collection interface {
	DataSource
	
	// Mount a DataSource at the specified location, the name may be used as a location ID or
	// file/directory name, depending on context.
	Mount(name string, ds DataSource)
	
	// GetChild attempts to lookup a child DataSource by path.
	GetChild(path string) (DataSource, error)
}

// The functions below this point are common utilities for working with AXIS paths.
// They are exported so that they may be used by custom DataSource implementations.

// IsAbs returns true if the path is not a relative path (includes no "." or ".." parts).
func IsAbs(path string) bool {
	// AXIS paths may start with an arbitrary number of location IDs that are separated by colons,
	// they have no meaning to this function, so just strip them off.
	locs := strings.Split(path, ":")
	path = locs[len(locs) - 1]
	dirs := strings.Split(path, "/")
	
	for i := range dirs {
		if dirs[i] == ".." || dirs[i] == "." {
			return false
		}
	}
	return true
}

// StripLoc removes the first location ID from the path and returns it along with the remainder.
// If there is no location ID in the path it returns "" for the id and the path as the remainder.
func StripLoc(path string) (loc string, remainder string) {
	parts := strings.SplitN(path, ":", 2)
	if len(parts) != 2 {
		return "", path
	}
	return parts[0], parts[1]
}

// StripDir removes the first directory name from the path and returns it along with the remainder.
// Any location IDs will be stripped.
// Paths consisting of a single directory will have an empty remainder.
func StripDir(path string) (loc string, remainder string) {
	locs := strings.Split(path, ":")
	path = locs[len(locs) - 1]
	
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		return path, ""
	}
	return parts[0], parts[1]
}

// Sanitize turns all OS specific path separators into slashes and then runs IsAbs on the result.
// If IsAbs returns false an empty string and a bad path error are returned.
func Sanitize(path string) (string, error) {
	path = filepath.ToSlash(path)
	
	if !IsAbs(path) {
		return "", NewError(ErrBadPath, path)
	}
	
	return path, nil
}
