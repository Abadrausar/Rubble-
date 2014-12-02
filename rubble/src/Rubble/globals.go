/*
Copyright 2013 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. You may not alter this software in any way.

3. This notice may not be removed or altered from any source distribution.
*/

package main

type RawFile struct {
	Path string
	Content string
	Namespace string
	Skip bool
}

var RawOrder = make([]string, 0, 100)
var RawFiles = make(map[string]*RawFile, 100)
var CurFile string // the name of the file being parsed
var CurNamespace = "" // used during loading

var CurWalkDir string // the path of the dir being traversed by filepath.Walk

var AddonNames = make([]string, 0, 10)

// The current parse stage
var ParseStage = 0

// This is where template variables and config options are stored
var VariableData = make(map[string]string)

var PrevParams = make([]string, 0)

// Used by the error handler and lexer
var LastLine = 1
