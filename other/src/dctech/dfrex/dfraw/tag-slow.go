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

package dfraw

import "dctech/rex"

type Tag struct {
	ID string
	Params rex.Indexable
	
	Replaced string
	
	Disabled bool
	Prepended string
	Appended string
	
	// Filled with all the trailing "stuff" (including white space)
	Comments string
}

func NewTag() *Tag {
	return &Tag {
		Params: rex.NewArray([]*rex.Value{}),
	}
}

func (t *Tag) Get(index string) *rex.Value {
	if index == "replace" {
		return rex.NewValueString(t.Replaced)
	}
	
	if index == "append" {
		return rex.NewValueString(t.Appended)
	}

	if index == "prepend" {
		return rex.NewValueString(t.Prepended)
	}

	if index == "disable" {
		return rex.NewValueBool(t.Disabled)
	}
	
	if index == "comments" {
		return rex.NewValueString(t.Comments)
	}

	if index == "id" {
		return rex.NewValueString(t.ID)
	}

	if index == "params" {
		return rex.NewValueIndex(t.Params)
	}
	
	rex.RaiseError("Index not valid.")
	panic("UNREACHABLE")
}

func (t *Tag) Set(index string, value *rex.Value) bool {
	if index == "replace" {
		t.Replaced = value.String()
		return true
	}
	
	if index == "append" {
		t.Appended = value.String()
		return true
	}

	if index == "prepend" {
		t.Prepended = value.String()
		return true
	}

	if index == "disable" {
		t.Disabled = value.Bool()
		return true
	}
	
	if index == "comments" {
		t.Comments = value.String()
		return true
	}
	
	if index == "id" {
		t.ID = value.String()
		return true
	}
	
	if index == "params" {
		if value.Type != rex.TypIndex {
			rex.RaiseError("Value passed as new params value not indexable.")
		}
		
		t.Params = value.Data.(rex.Indexable)
		return true
	}
	
	rex.RaiseError("Index not valid.")
	panic("UNREACHABLE")
}

func (t *Tag) Exists(index string) bool {
	if index == "replace" || index == "append" || index == "prepend" || index == "disable" || index == "comments" {
		return true
	}

	if index == "id" || index == "params" {
		return true
	}
	return false
}

func (t *Tag) Len() int64 {
	return 7
}

func (t *Tag) Keys() []string {
	return []string {
		"prepend",
		"id",
		"params",
		"append",
		"comments",
		"replace",
		"disable",
	}
}

func (t *Tag) String() string {
	if t.Replaced != "" {
		return t.Prepended + t.Replaced + t.Appended + t.Comments
	}
	
	b := "["
	e := "]"
	if t.Disabled {
		b = "-"
		e = "-"
	}
	
	out := t.Prepended + b + t.ID
	for _, val := range t.Params.Keys() {
		out += ":" + t.Params.Get(val).String()
	}
	return out + e + t.Appended + t.Comments
}

func (t *Tag) CodeString() string {
	return "\"" + t.String() + "\""
}
