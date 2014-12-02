PACKAGE DOCUMENTATION

package sort
    import "dctech/rex/commands/sort"



FUNCTIONS

func Command_Array(script *rex.Script, params []*rex.Value)
    Sorts an array (or any other IntEditIndexable type).

	sort:array array

    This works best if all values in the array are the same type. Returns
    unchanged.

func Command_Map(script *rex.Script, params []*rex.Value)
    Converts any existing indexable into a sort:map.

	sort:map map

    Returns a new sort:map with the same keys/values as the old indexable.

func Command_New(script *rex.Script, params []*rex.Value)
    Creates a new ordered map.

	sort:new

    Returns a new (empty) ordered map.


