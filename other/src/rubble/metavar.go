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

// MetaVar is used to store meta data about a config var.
type MetaVar struct {
	Name    string // A user friendly name, may be ""
	Val     string
	Choices []string
}

func NewMetaVar(value *rex.Value) *MetaVar {
	if value.Type == rex.TypIndex {
		v, ok := value.Data.(*MetaVar)
		if ok {
			return v
		}
	}
	meta := new(MetaVar)
	meta.Val = value.String()
	return meta
}

func NewMetaVarFromLit(script *rex.Script, keys []string, values []*rex.Value) *rex.Value {
	if keys == nil {
		rex.RaiseError("rubble:addonmetavar may not be initialized without keys.")
	}

	meta := new(MetaVar)
	for i := range keys {
		switch keys[i] {
		case "name":
			meta.Name = values[i].String()
		case "val":
			meta.Val = values[i].String()
		case "choices":
			meta.Choices = append(meta.Choices, values[i].String())
		default:
			rex.RaiseError("rubble:addonmetavar Invalid key name: " + keys[i])
		}
	}
	return rex.NewValueIndex(meta)
}

func (meta *MetaVar) Get(index string) *rex.Value {
	switch index {
	case "name":
		return rex.NewValueString(meta.Name)
	case "val":
		return rex.NewValueString(meta.Val)
	case "choices":
		vals := []*rex.Value{}
		for _, s := range meta.Choices {
			vals = append(vals, rex.NewValueString(s))
		}
		return rex.NewValueIndex(rex.NewStaticArray(vals))
	default:
		return rex.NewValue()
	}
}

func (meta *MetaVar) Exists(index string) bool {
	switch index {
	case "name":
	case "val":
	case "choices":
	default:
		return false
	}
	return true
}

func (meta *MetaVar) Len() int64 {
	return int64(3)
}

func (meta *MetaVar) Keys() []string {
	return []string{
		"name",
		"val",
		"choices",
	}
}

func (meta *MetaVar) String() string {
	return rex.IndexableToString("rubble:addonmetavar", meta)
}

func (meta *MetaVar) CodeString() string {
	return rex.IndexableToCodeString("rubble:addonmetavar", meta, true)
}
