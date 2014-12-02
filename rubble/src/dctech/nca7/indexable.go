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

package nca7

//import "fmt"

// Indexable represents a ScriptObject that may be used with the 
// dereference opperator ([]) as well as some of the base commands.
// Indexable is readonly.
type Indexable interface {
	Get(index string) *Value
	Exists(index string) bool
	Len() int64
	Keys() []string
}

// EditIndexable represents a ScriptObject that may be used with the 
// dereference opperator ([]) as well as some of the base commands.
// EditIndexable is read/write.
type EditIndexable interface {
	Indexable
	
	// Set returns false if the key could not be created
	// This may be caused by the indexable being semi-readonly
	Set(index string, value *Value) bool 
}

// IndexableFactory is the function signature for a literal converter.
// keys may be nil, a specific implimentation may panic with an error if
// a nil or non-nil keys is expected but not found.
type IndexableFactory func(keys []string, values []*Value) *Value

// ToIndexable will try to convert a ScriptObject to an Indexable, returns nil on failure.
func ToIndexable(val ScriptObject) Indexable {
	new, ok := val.(Indexable)
	if !ok {
		return nil
	}
	return new
}

// ToEditIndexable will try to convert a ScriptObject to an EditIndexable, returns nil on failure.
func ToEditIndexable(val ScriptObject) EditIndexable {
	new, ok := val.(EditIndexable)
	if !ok {
		return nil
	}
	return new
}
