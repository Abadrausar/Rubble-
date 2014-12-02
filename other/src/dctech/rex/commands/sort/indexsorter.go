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

package sort

import "dctech/rex"

type indexSorter struct {
	data rex.IntEditIndexable
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (sorter *indexSorter) Less(i, j int) bool {
	a := sorter.data.GetI(int64(i))
	b := sorter.data.GetI(int64(j))
	
	switch a.Type {
	case rex.TypNil:
		// Nils go to the end of the line
		return b.Type != rex.TypNil
	case rex.TypString:
		return a.String() < b.String()
	case rex.TypInt:
		return a.Int64() < b.Int64()
	case rex.TypFloat:
		return a.Float64() < b.Float64()
	default:
		// Non-comparable types always come after comparable types
		return b.Type == rex.TypNil || b.Type == rex.TypString || b.Type == rex.TypInt || b.Type == rex.TypFloat
	}
	rex.RaiseError("Script Value has invalid Type.")
	panic("UNREACHABLE")
}

func (sorter *indexSorter) Swap(i, j int) {
	tmp := sorter.data.GetI(int64(i))
	sorter.data.SetI(int64(i), sorter.data.GetI(int64(j)))
	sorter.data.SetI(int64(j), tmp)
}

func (sorter *indexSorter) Len() int {
	return int(sorter.data.Len())
}
