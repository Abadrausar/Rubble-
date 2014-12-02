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

package rblutil

import "dctech/axis"

// Strip the extension from a file name.
// If a file has multiple extensions strip only the last.
func StripExt(name string) string {
	i := len(name) - 1
	for i >= 0 {
		if name[i] == '.' {
			return name[:i]
		}
		i--
	}
	return name
}

// LoadOr is used for loading HTML templates and the like for the web UI or the documentation generator.
// It attempts to load the file from "rubble:other/webUI/", if this fails "file" is written out and then used.
func LoadOr(fs axis.DataSource, name, file string) string {
	content, err := axis.ReadAll(fs, "rubble:other/webUI/" + name)
	if err != nil {
		axis.WriteAll(fs, "rubble:other/webUI/" + name, []byte(file))
		return file
	}
	return string(content)
}
