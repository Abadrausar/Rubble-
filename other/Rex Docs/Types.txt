
Available Indexable Types:

array
sarray
map
smap
	The standard indexable types.
	See "Rex Docs/Rex Basics"

sort:map
	An alphabetical ordered map.
	Exactly like map, except foreach will always process keys in alphabetical order.
	Can also be created by the commands sort:map or sort:new.

struct
	The struct type is basically a map indexable that has a fixed set of keys.
	The keys are determined based on the value passed to the "proto" key (which
	MUST be first, and MUST be an indexable), any keys which are not given a value
	take their value from the prototype indexable.
	To put it simply a struct is a copy of a map that cannot have new keys added.

rubble:addonmeta
rubble:addonmetavar
	Special struct-like types for various kinds of meta data (in addon.meta files).
	See "Rubble Basics".
