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

package axis_test

// This is not a real test!
// Real testing is accomplished by using AXIS in real-world applications (namely Rubble).

import "dctech/axis"
import "fmt"

func Example() {
	FS := axis.NewFileSystem()
	
	fs := axis.NewOSDir(".")
	
	// axis.Mount(FS, "wd", fs) works as well (but is not generally used).
	FS.Mount("wd", fs)

	fmt.Println("Files (according to AXIS OSDir):")
	for _, dir := range axis.ListFile(fs, "") {
		fmt.Println(dir)
	}
	fmt.Println()
	
	fmt.Println("Files (according to AXIS FileSystem):")
	for _, dir := range axis.ListFile(FS, "wd:") {
		fmt.Println(dir)
	}
	
	// Output:
	// Files (according to AXIS OSDir):
	// axis_test.go
	// errors.go
	// filesystem.go
	// interface.go
	// logicdir.go
	// logicfile.go
	// osdir.go
	// osfile.go
	// path.go
	// 
	// Files (according to AXIS FileSystem):
	// axis_test.go
	// errors.go
	// filesystem.go
	// interface.go
	// logicdir.go
	// logicfile.go
	// osdir.go
	// osfile.go
	// path.go
}
