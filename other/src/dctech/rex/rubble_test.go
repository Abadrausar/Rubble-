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

package rex_test

import "testing"

//import "dctech/rex"
import "rubble/guts"

// Does a full Rubble run, use for coverage info and the like

// To use:
//	go test -c -covermode=count dctech/rex
//	rex.test.exe -test.coverprofile=rex.cover
// Just running go test causes weird errors (caused, I think, by running Rubble in a temp dir)

// Commented out, due to the above issue (so normal test are not messed up)
//func Test_Rubble(t *testing.T) {
//	guts.ContinueOnBadFlag = true
//	guts.Main()
//}
