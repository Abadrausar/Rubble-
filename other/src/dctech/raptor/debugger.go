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

// The debugger support is very bare-bones and I think it may not be fully wired into the state's parsing functions.
// Basicly I need to make sure the callback system is fully wired in and then I need to write a wrapper 
// API to make all this useable.

// Debug handler function prototype
type DbgHandler func(*State)

const (
	// Debugger callback ids
	DbgrAdvanceTkn     = iota // Just before a token is retrieved
	DbgrEnterCmd              // Before any related tokens are processed
	DbgrLeaveCmd              // After all related tokens are processed and execution finishes
	DbgrEnterDeref            // Before any related tokens are processed
	DbgrLeaveDeref            // After all related tokens are processed and value is retrived
	DbgrEnterObjLit           // Before any related tokens are processed
	DbgrLeaveObjLit           // After all related tokens are processed and value is created
	DbgrEnterCodeBlock        // Before any related tokens are processed
	DbgrLeaveCodeBlock        // After all related tokens are processed and value is created
	dbgrMaxType               // this one MUST BE LAST
)
