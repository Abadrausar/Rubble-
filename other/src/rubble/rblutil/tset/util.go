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

var objTypes = map[string]map[string]bool {
	"ITEM_TOOL": map[string]bool {
		"TILE": true,
	},
	"CREATURE": map[string]bool {
		"CASTE_ALTTILE": true,
		"CASTE_COLOR": true,
		"CASTE_GLOWCOLOR": true,
		"CASTE_GLOWTILE": true,
		"CASTE_SOLDIER_ALTTILE": true,
		"CASTE_SOLDIER_TILE": true,
		"CASTE_TILE": true,
		"CREATURE_SOLDIER_TILE": true,
		"CREATURE_TILE": true,
		"COLOR": true,
		"GLOWCOLOR": true,
		"GLOWTILE": true,
		"SOLDIER_ALTTILE": true,
	},
	"PLANT": map[string]bool {
		// Plant growths currently require special handling...
		"GROWTH": true,
		"GROWTH_PRINT": true,
		"PICKED_TILE": true,
		"DEAD_PICKED_TILE": true,
		"SHRUB_TILE": true,
		"DEAD_SHRUB_TILE": true,
		"PICKED_COLOR": true,
		"DEAD_PICKED_COLOR": true,
		"SHRUB_COLOR": true,
		"DEAD_SHRUB_COLOR": true,
		"GRASS_TILES": true,
		"ALT_GRASS_TILES": true,
		"GRASS_COLORS": true,
		"TREE_TILE": true,
		"DEAD_TREE_TILE": true,
		"SAPLING_TILE": true,
		"DEAD_SAPLING_TILE": true,
		"TREE_COLOR": true,
		"DEAD_TREE_COLOR": true,
		"SAPLING_COLOR": true,
		"DEAD_SAPLING_COLOR": true,
	},
	"INORGANIC": map[string]bool {
		"ITEM_SYMBOL": true,
		"TILE": true,
		"DISPLAY_COLOR": true,
		"BUILD_COLOR": true,
		"TILE_COLOR": true,
		"BASIC_COLOR": true,
	},
	"MATERIAL_TEMPLATE": map[string]bool {
		"ITEM_SYMBOL": true,
		"TILE": true,
		"DISPLAY_COLOR": true,
		"BUILD_COLOR": true,
		"TILE_COLOR": true,
		"BASIC_COLOR": true,
	},
}

// Reduce removes all tags that are not useful for tilesets
// and strips the comments from what is left.
func Reduce(tags []*Tag) []*Tag {
	out := make([]*Tag, 0)
	
	keepset := map[string]bool{}
	for _, tag := range tags {
		tag.Comments = ""
		
		if _, ok := objTypes[tag.ID]; ok {
			keepset = objTypes[tag.ID]
			out = append(out, tag)
			continue
		}

		if keepset[tag.ID] {
			out = append(out, tag)
		}
	}
	
	return out
}

// Flatten combines multiple raw files into one large one.
func Flatten(files [][]*Tag) []*Tag {
	l := 0
	for _, tags := range files {
		l += len(tags)
	}
	
	out := make([]*Tag, 0, l)
	for _, tags := range files {
		out = append(out, tags...)
	}
	return out
}

// Tableize makes a set of lookup tables from a raw file (generally the output from Flatten).
// Errors are silently skipped (with any tags that are part of an invalid object being dropped).
func Tableize(tags []*Tag) map[string]map[string]*Tag {
	out := make(map[string]map[string]*Tag)
	
	var curobject map[string]*Tag
	aux := ""
	for _, tag := range tags {
		if _, ok := objTypes[tag.ID]; ok {
			if len(tag.Params) != 1 {
				// Bad param count, definitely an error.
				curobject = nil
				aux = ""
				continue
			} 
			
			if old, ok := out[tag.Params[0]]; ok {
				// Duplicate object! Should be an error?
				curobject = old
				aux = ""
				continue
			}
			
			curobject = make(map[string]*Tag)
			out[tag.Params[0]] = curobject
			aux = ""
			continue
		}
		
		if curobject != nil {
			if tag.ID == "GROWTH" {
				if len(tag.Params) != 1 {
					// Bad param count, definitely an error.
					continue
				}
				
				aux = "GROWTH_" + tag.Params[0] + ":"
				continue
			}
			
			curobject[aux + tag.ID] = tag
		}
	}
	return out
}

// Convenience function
func lookup(table map[string]map[string]*Tag, obj, tag string) *Tag {
	if objT, ok := table[obj]; ok {
		if tagV, ok := objT[tag]; ok {
			return tagV
		}
	}
	return nil
}

// Convenience function
func swapIfExist(tag *Tag, obj, aux string, table map[string]map[string]*Tag) {
	if aux != "" {
		aux += ":"
	}
	
	ntag := lookup(table, obj, aux + tag.ID)
	if ntag == nil {
		return
	}
	
	// We don't want to change the comments and ID is already the same...
	tag.Params = ntag.Params
}

// Normalize does a simple pretty print of a raw file.
// The formatting is VERY simple and stupid, it is only suitable for a tset file.
func Normalize(tags []*Tag) {
	var prev *Tag
	for _, tag := range tags {
		if prev != nil {
			if _, ok := objTypes[tag.ID]; ok {
				prev.Comments = "\n\n"
			} else {
				prev.Comments = "\n\t"
			}
		}
		
		prev = tag
	}
}

// FormatFile takes a set of tags and turns them into a byte slice ready to write.
func FormatFile(tags []*Tag) []byte {
	out := []byte{}

	for _, tag := range tags {
		out = append(out, tag.String()...)
	}
	
	return out
}
