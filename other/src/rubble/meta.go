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
// Exposed to scripts as the "rubble:addonmeta" indexable type.
type Meta struct {
	// Is this addon an automatically managed library?
	Lib bool

	// A one line addon description.
	// For use by user interfaces (not used directly by Rubble).
	Header string

	// A longer addon description, may be as long as you like.
	// If Header is an adequate description leave this empty.
	// For use by user interfaces (not used directly by Rubble).
	Description string

	// How is the value of Description and Header formatted?
	// Valid values: "html" (formatted HTML code)
	// Anything else is assumed to mean the documentation is pre-formatted plain text.
	Format string
	
	// Addon names for addons that are automatically activated when this addon is active.
	Activates []string

	// Addon names for addons that are incompatible with this addon.
	Incompatible []string

	// Configuration variables used by this addon and their default values.
	// For use by user interfaces (not used directly by Rubble).
	Vars map[string]*MetaVar
}

func NewMeta() *Meta {
	return &Meta{
		Activates: []string{},
		Incompatible: []string{},
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
		case "format":
			meta.Format = strings.TrimSpace(values[i].String())
		case "activates":
			vals := strings.Split(values[i].String(), ";")
			for i := range vals {
				vals[i] = strings.TrimSpace(vals[i])
				if vals[i] == "" {
					vals = append(vals[:i], vals[i+1:]...)
				}
			}
			meta.Activates = vals
		case "incompatible":
			vals := strings.Split(values[i].String(), ";")
			for i := range vals {
				vals[i] = strings.TrimSpace(vals[i])
				if vals[i] == "" {
					vals = append(vals[:i], vals[i+1:]...)
				}
			}
			meta.Incompatible = vals
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
	case "format":
		return rex.NewValueString(meta.Format)
	case "activates":
		return rex.NewValueString(strings.Join(meta.Activates, ";"))
	case "incompatible":
		return rex.NewValueString(strings.Join(meta.Incompatible, ";"))
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
	case "format":
		meta.Format = value.String()
	case "activates":
		meta.Activates = strings.Split(value.String(), ";")
	case "incompatible":
		meta.Incompatible = strings.Split(value.String(), ";")
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
	case "format":
	case "activates":
	case "incompatible":
	case "vars":
	default:
		return false
	}
	return true
}

func (meta *Meta) Len() int64 {
	return int64(7)
}

func (meta *Meta) Keys() []string {
	return []string{
		"lib",
		"header",
		"description",
		"format",
		"activates",
		"incompatible",
		"vars",
	}
}

func (meta *Meta) String() string {
	return rex.IndexableToString("rubble:addonmeta", meta)
}
