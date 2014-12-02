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

/*
Rex Generic Indexable Interface.

This package implements some basic indexable types that provides an interface to
arbitrary go data. Unfortunately there is no real way to make a generic object literal
converter, so all object creation needs to take place from the host or via custom
native commands.

Struct fields are accessible by using their name as the key (case sensitive), only exported
fields are editable, but all fields are readable.

Some types cannot be converted to a script value, in such cases script value nil will be returned.
Script values cannot be converted to some go types, in such cases an error is generated.

Due to limitations in the Indexable interface it is not possible to convert any map type
that does not have string keys.
It would be possible (albeit very difficult) to support some of the other basic types as keys,
but the effort is not worth the reward.

Slice types do not support the append key yet.

Map types do not support adding keys yet.

Be warned that this interface will expose an object and all of it's children to scripts,
things can be very easily messed up if precautions are not taken!
It is highly recommended that only objects that are designed for script access be exported!

No special case code is inserted to deal with rex.Values, they are treated like any other struct!

Any value that implements the Indexable interface will use that implementation rather then using GenII,
for this to work the value has to be settable (this is a limitation of go's reflection library).

While concurrent access of a data store is handled by the state there is no provision for
safe access of individual values. This is not a problem most of the time, but indexables that are 
used by more than one script at once can cause issues.
If you need concurrent access to a GenII object you will need to provide a synchronization interface to scripts.
Most of the time it is a better idea to write a custom indexable, as otherwise there is no way to enforce
synchronization.
The lazy route is to have a lock that scripts agree to grab and free when accessing the interface,
but as mentioned there is no way to enforce this (and if an error happens the lock can get "stuck").

This package is like the bad side of town, if you get mugged don't come crying to me :)

*/
package genii

// TODO: Should rex.Values receive special handling?

// TODO: Allow adding keys to map types.

// TODO: How hard would it be to add the "append" key to slice types?
