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

import "flag"

var DFDir string
var OutputDir string
var ConfigDir string
var BaseDir string
var AddonsDir string

var LexTest bool
var Recover bool

func ParseCommandLine() {
	flag.StringVar(&DFDir ,"dfdir", "..", "What is the base DF directory?")
	flag.StringVar(&OutputDir ,"outputdir", "./objects", "Where should Rubble write the generated raw files?")
	flag.StringVar(&ConfigDir ,"configdir", "./source", "Rubble config dir.")
	flag.StringVar(&BaseDir ,"basedir", "./source/base", "Rubble base dir.")
	flag.StringVar(&AddonsDir ,"addonsdir", "./source/addons", "Rubble addons dir. (may be a path list)")
	
	flag.BoolVar(&LexTest, "lextest", false, "Run a lexer test. No files will be written.")
	flag.BoolVar(&Recover, "recover", true, "Should NCA4 recover errors? Useful for debugging.")
	
	flag.Parse()
}
