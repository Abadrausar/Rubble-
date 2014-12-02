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

package rex

// Indexable represents an object that may be used with the
// variable operator ([]) as well as some of the base commands.
// Indexable is readonly.
type Indexable interface {
	// Get should return a type nil value for non-existent keys.
	Get(index string) *Value
	Exists(index string) bool
	Len() int64
	Keys() []string

	// String should return the contents in the best format for humans to read.
	// The output should resemble Rex code but it does not need to use things like proper types.
	String() string
}

// EditIndexable represents an object that may be used with the
// variable operator ([]) as well as some of the base commands.
// EditIndexable is read/write.
// It is STRONGLY encouraged that all EditIndexables be thread safe!
type EditIndexable interface {
	Indexable

	// Set returns false if the key could not be created.
	// This may be caused by the indexable being semi-readonly.
	Set(index string, value *Value) bool
}

// IntIndexable is an indexable that may be indexed by integer.
// Indexes MUST range from 0 to Len() - 1, otherwise all sorts of things will break!
type IntIndexable interface {
	Indexable

	// GetI should return a type nil value for non-existent keys.
	GetI(index int64) *Value
}

// IntEditIndexable is an edit indexable that may be indexed by integer.
// Indexes MUST range from 0 to Len() - 1, otherwise all sorts of things will break!
type IntEditIndexable interface {
	EditIndexable

	// SetI returns false if the key could not be created.
	// This may be caused by the indexable being semi-readonly.
	// Conventionally -1 is used as a shorthand for append (but this is not required).
	SetI(index int64, value *Value) bool
}

// IndexableToString is a helper function for implementing an Indexable's String method.
// If typ is empty the generic type Indexable is used.
func IndexableToString(typ string, data Indexable) string {
	out := "<"
	if typ != "" {
		out += typ
	} else {
		out += "Indexable"
	}
	keys := data.Keys()
	for i := range keys {
		out += ", " + keys[i] + "=" + data.Get(keys[i]).String()
	}
	return out + ">"
}
