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
import "strings"

// isAbs returns true if the path is not a relative path (includes no "." or ".." parts).
func isAbs(path string) bool {
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

// Path represents a parsed AXIS path.
type Path struct {
	locations []string
	usedlocs int
	
	dirs []string
	useddirs int
}

// NewPath parses a string into an AXIS path.
func NewPath(source string) *Path {
	source = strings.TrimSpace(filepath.ToSlash(source))
	if !isAbs(source) {
		return nil
	}
	
	path := new(Path)
	if len(source) == 0 {
		path.locations = make([]string, 0)
		path.dirs = make([]string, 0)
	} else if source[len(source) - 1:] == ":" {
		path.locations = strings.Split(source, ":")
		path.locations = path.locations[:len(path.locations) - 1]
		path.dirs = make([]string, 0)
	} else {
		locs := strings.Split(source, ":")
		source = locs[len(locs) - 1]
		path.locations = locs[:len(locs) - 1]
		path.dirs = strings.Split(source, "/")
		if path.dirs[len(path.dirs) - 1] == "" {
			path.dirs = path.dirs[:len(path.dirs) - 1]
		}
	}
	return path
}

func (path *Path) String() string {
	out := ""
	for i := range path.locations {
		out += path.locations[i] + ":"
	}
	for i := range path.dirs {
		out += path.dirs[i] + "/"
	}
	
	if len(out) != 0 {
		if out[len(out) - 1] == '/' {
			return out[:len(out) - 1]
		}
	}
	return out
}

// Done returns true once all path elements are marked used.
func (path *Path) Done() bool {
	if len(path.locations) > path.usedlocs {
		return false
	}
	if len(path.dirs) > path.useddirs {
		return false
	}
	return true
}

// Remainder returns a string constructed from the remaining path elements.
// Location IDs are always ignored (as this is meant to be used when talking to the OS).
// If the returned string would be empty returns ".".
func (path *Path) Remainder() string {
	if len(path.dirs) > path.useddirs {
		out := ""
		for _, v := range path.dirs[path.useddirs:] {
			out += v + "/"
		}
		if len(out) != 0 {
			if out[len(out) - 1] == '/' {
				return out[:len(out) - 1]
			}
		}
		return out
	}
	return "."
}

// NextLoc returns the next location ID in the path.
// Returns "" if no location IDs remain unused.
func (path *Path) NextLoc() string {
	if len(path.locations) <= path.usedlocs {
		return ""
	}
	
	rtn := path.locations[path.usedlocs]
	path.usedlocs++
	return rtn
}

// RevertLoc resets the last location ID as being unused, if no used location IDs remain does nothing.
// Be careful! if NextLoc returns a blank string no location IDs were marked used!
func (path *Path) RevertLoc() {
	path.usedlocs--
	if path.usedlocs < 0 {
		path.usedlocs = 0
	}
}

// NextDir returns the next directory name in the path (the last element may be a file name).
// Returns "" if locations remain unused or there are no directory names left.
func (path *Path) NextDir() string {
	if len(path.locations) > path.usedlocs {
		return ""
	}
	
	if len(path.dirs) <= path.useddirs {
		return ""
	}
	
	rtn := path.dirs[path.useddirs]
	path.useddirs++
	return rtn
}

// RevertDir resets the last directory as being unused, if no used directories remain does nothing.
// Be careful! if NextDir returns a blank string no directories were marked used!
func (path *Path) RevertDir() {
	path.useddirs--
	if path.useddirs < 0 {
		path.useddirs = 0
	}
}
