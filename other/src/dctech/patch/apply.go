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
//import "fmt"

// Apply applies a series of ordered patches to a string.
// The patches MUST be in descending order!
func Apply(patches []*Patch, text string) (string, []bool) {
	if len(patches) == 0 {
		return text, []bool{}
	}
	applied := make([]bool, len(patches))
	
	text = strings.Replace(text, "\r", "", -1)
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return text, applied
	}
	
apply:
	for x, patch := range patches {
		if len(patch.diffs) == 0 {
			continue
		}
		
		i := findStart(lines, patch)
		if i == -1 {
			continue
		}
		
		// Apply the patch
		
		// First we check to make sure all is good
		pStart := i
		pEnd := i
		mid := make([]string, 0)
		for _, op := range patch.diffs {
			switch op.Type {
			case DiffInsert:
				mid = append(mid, op.Text)
				continue
			case DiffDelete:
				if lines[i] != op.Text {
					continue apply
				}
				pEnd++
			case DiffEqual:
				if lines[i] != op.Text {
					continue apply
				}
				mid = append(mid, op.Text)
				pEnd++
			}
			i++
		}
		
		// Then we apply the patch
		start := lines[:pStart]
		end := lines[pEnd:]
		lines = make([]string, 0, len(start) + len(mid) + len(end))
		lines = append(lines, start...)
		lines = append(lines, mid...)
		lines = append(lines, end...)
		applied[x] = true
	}
	return strings.Join(lines, "\n"), applied
}

func findStart(lines []string, patch *Patch) int {
	for i := 0; i < 10; i++ {
		if patch.start2 + i < len(lines) {
			if lines[patch.start2 + i] == patch.diffs[0].Text {
				return patch.start2 + i
			}
		}
		
		if patch.start2 - i > 0 {
			if lines[patch.start2 - i] == patch.diffs[0].Text {
				return patch.start2 - i
			}
		}
	}
	return -1
}
