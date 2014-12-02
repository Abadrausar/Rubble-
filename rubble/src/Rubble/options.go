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
