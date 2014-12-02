/*
Copyright 2014 by Milo Christiansen

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

package tset

// Apply modifies the tags in file with information from table.
// TODO: This does not handle plant growths correctly!
func Apply(file []*Tag, table map[string]map[string]*Tag) {
	obj := ""
	aux := ""
	
	isPlant := false
	
	for _, tag := range file {
		if _, ok := objTypes[tag.ID]; ok {
			if len(tag.Params) != 1 {
				// Bad param count, definitely an error.
				obj = ""
				aux = ""
				continue
			}
			isPlant = tag.ID == "PLANT"
			
			obj = tag.Params[0]
			aux = ""
			continue
		}
		
		if isPlant {
			if tag.ID == "GROWTH" {
				if len(tag.Params) != 1 {
					// Bad param count, definitely an error.
					continue
				}
				
				aux = "GROWTH_" + tag.Params[0]
				continue
			}
		}
		
		if obj != "" {
			swapIfExist(tag, obj, aux, table)
		}
	}
}

// ApplyToFile parses a file and applies the information from table, then returns the result.
func ApplyToFile(file []byte, table map[string]map[string]*Tag) []byte {
	tags := ParseRaws(file)
	Apply(tags, table)
	return FormatFile(tags)
}
