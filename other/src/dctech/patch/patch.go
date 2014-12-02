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

// A simple patch library.
// Currently only supports applying patches using a very simple algorithm.
package patch

import "strings"
import "strconv"
import "bytes"
import "regexp"

// ToText takes a list of patches and returns a textual representation.
func ToText(patches []*Patch) string {
	var text bytes.Buffer
	for _, aPatch := range patches {
		text.WriteString(aPatch.String())
	}
	return text.String()
}

var headerRegex = regexp.MustCompile("^@@ -(\\d+),(\\d+) \\+(\\d+),(\\d+) @@$")

// FromText parses a textual representation of patches and returns a List of Patch objects.
func FromText(text string) []*Patch {
	patches := make([]*Patch, 0)
	text = strings.Replace(text, "\r", "", -1)
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return nil
	}

	i := 0
	for i < len(lines) {
		if !headerRegex.MatchString(lines[i]) {
			return nil
		}
		
		header := headerRegex.FindStringSubmatch(lines[i])
		
		patch := new(Patch)
		patch.start1, _ = strconv.Atoi(header[1])
		patch.start1--
		patch.length1, _ = strconv.Atoi(header[2])
		patch.start2, _ = strconv.Atoi(header[3])
		patch.start2--
		patch.length2, _ = strconv.Atoi(header[4])
		i++
		
		var sign uint8
		for i < len(lines) {
			if len(lines[i]) > 0 {
				sign = lines[i][0]
			} else {
				i++
				continue
			}

			line := lines[i][1:]
			if sign == '-' {
				// Deletion.
				patch.diffs = append(patch.diffs, &Diff{DiffDelete, line})
			} else if sign == '+' {
				// Insertion.
				patch.diffs = append(patch.diffs, &Diff{DiffInsert, line})
			} else if sign == ' ' {
				// Minor equality.
				patch.diffs = append(patch.diffs, &Diff{DiffEqual, line})
			} else if sign == '@' {
				// Start of next patch.
				break
			} else {
				// Bad mode sign, treat like a comment (should probably be an error)
				continue
			}
			
			i++
		}
		
		patches = append(patches, patch)
	}
	return patches
}
