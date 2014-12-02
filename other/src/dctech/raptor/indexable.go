/*
Copyright 2012-2013 by Milo Christiansen

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

package raptor

//import "fmt"

// Indexable represents an object that may be used with the 
// dereference opperator ([]) as well as some of the base commands.
// Indexable is readonly.
type Indexable interface {
	Get(index string) *Value
	Exists(index string) bool
	Len() int64
	Keys() []string

	// String should return the contents in the best format for humans to read.
	// The output should resemble Raptor code but it does not need to use things like proper types.
	String() string

	// Code should return the contents as valid Raptor code, if this is impossible return a 
	// string that is descriptive of the type eg "<SomeUnconvertableType>".
	// If the Indexable does not have a type name it is ok to use map or array if the data is written out in a form
	// the map or array ObjectFactory can understand.
	CodeString() string
}

// EditIndexable represents an object that may be used with the 
// dereference opperator ([]) as well as some of the base commands.
// EditIndexable is read/write.
type EditIndexable interface {
	Indexable

	// Set returns false if the key could not be created.
	// This may be caused by the indexable being semi-readonly.
	Set(index string, value *Value) bool
}

// ToIndexable will try to convert an object to an Indexable, returns nil on failure.
func ToIndexable(val EmptyInterface) Indexable {
	new, ok := val.(Indexable)
	if !ok {
		return nil
	}
	return new
}

// ToEditIndexable will try to convert an object to an EditIndexable, returns nil on failure.
func ToEditIndexable(val EmptyInterface) EditIndexable {
	new, ok := val.(EditIndexable)
	if !ok {
		return nil
	}
	return new
}

// GenericIndexableToString is a helper function for implimenting an Indexable's String method.
// If typ is empty the generic type Indexable is used.
func GenericIndexableToString(typ string, data Indexable) string {
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

// GenericIndexableToCodeString is a helper function for implimenting an Indexable's CodeString method.
// If typ is empty the generic type Indexable is used and the result is a string not an object literal.
// If haskeys is false the result does not have key data.
func GenericIndexableToCodeString(typ string, data Indexable, haskeys bool) string {
	out := ""
	if typ != "" {
		out += "<" + typ
	} else {
		out += "\"<Indexable"
	}
	keys := data.Keys()
	for i := range keys {
		if haskeys {
			out += ", " + keys[i] + "=" + data.Get(keys[i]).CodeString()
		} else {
			out += ", " + data.Get(keys[i]).CodeString()
		}
	}
	if typ != "" {
		out += ">"
	} else {
		out += ">\""
	}
	return out
}
