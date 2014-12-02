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

package rubble

import "dctech/rex"
import "strings"

// Meta is used to store meta-data for an Addon.
type Meta struct {
	// Is this addon visible to the user?
	// Set to false for automatically activated libraries and the like.
	Lib bool

	// A one line addon description.
	Header string

	// A longer addon description, may be as long as you like.
	Description string

	// Addon names for addons that are automatically activated when this addon is active.
	Activates []string

	// Configuration variables used by this addon and their default values.
	Vars map[string]*MetaVar
}

func NewMeta() *Meta {
	return &Meta{
		Activates: []string{},
		Vars:      make(map[string]*MetaVar),
	}
}

func NewMetaFromLit(script *rex.Script, keys []string, values []*rex.Value) *rex.Value {
	if keys == nil {
		rex.RaiseError("rubble:addonmeta may not be initialized without keys.")
	}

	meta := NewMeta()
	for i := range keys {
		switch keys[i] {
		case "lib":
			meta.Lib = values[i].Bool()
		case "header":
			meta.Header = strings.TrimSpace(values[i].String())
		case "description":
			meta.Description = strings.TrimSpace(values[i].String())
		case "activates":
			vals := strings.Split(values[i].String(), ";")
			for i := range vals {
				vals[i] = strings.TrimSpace(vals[i])
				if vals[i] == "" {
					vals = append(vals[:i], vals[i+1:]...)
				}
			}
			meta.Activates = vals
		case "vars":
			if values[i].Type == rex.TypIndex {
				index := values[i].Data.(rex.Indexable)
				for _, i := range index.Keys() {
					meta.Vars[i] = NewMetaVar(index.Get(i))
				}
			}
		default:
			rex.RaiseError("rubble:addonmeta Invalid key name: " + keys[i])
		}
	}
	return rex.NewValueIndex(meta)
}

func (meta *Meta) Get(index string) *rex.Value {
	switch index {
	case "lib":
		return rex.NewValueBool(meta.Lib)
	case "header":
		return rex.NewValueString(meta.Header)
	case "description":
		return rex.NewValueString(meta.Description)
	case "activates":
		return rex.NewValueString(strings.Join(meta.Activates, ";"))
	case "vars":
		val, ok := meta.Vars[index]
		if !ok {
			return rex.NewValue()
		}
		return rex.NewValueIndex(val)
	default:
		return rex.NewValue()
	}
}

func (meta *Meta) Set(index string, value *rex.Value) bool {
	switch index {
	case "lib":
		meta.Lib = value.Bool()
	case "header":
		meta.Header = value.String()
	case "description":
		meta.Description = value.String()
	case "activates":
		meta.Activates = strings.Split(value.String(), ";")
	case "vars":
		meta.Vars[index] = NewMetaVar(value)
	default:
		return false
	}
	return true
}

func (meta *Meta) Exists(index string) bool {
	switch index {
	case "lib":
	case "header":
	case "description":
	case "activates":
	case "vars":
	default:
		return false
	}
	return true
}

func (meta *Meta) Len() int64 {
	return int64(5)
}

func (meta *Meta) Keys() []string {
	return []string{
		"lib",
		"header",
		"description",
		"activates",
		"vars",
	}
}

func (meta *Meta) String() string {
	return rex.IndexableToString("rubble:addonmeta", meta)
}

func (meta *Meta) CodeString() string {
	return rex.IndexableToCodeString("rubble:addonmeta", meta, true)
}
