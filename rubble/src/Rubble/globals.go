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
