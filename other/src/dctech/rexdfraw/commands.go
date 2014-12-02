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

// Rex Dwarf Fortress Raw Commands
package rexdfraw

import "dctech/rex"

// Adds the raw commands to the state.
// The raw commands are:
//	df:raw:walk
//	df:raw:parse
//	df:raw:dump
func Setup(state *rex.State) (err error) {
	defer func() {
		if !state.NoRecover {
			if x := recover(); x != nil {
				if y, ok := x.(rex.ScriptError); ok {
					err = y
					return
				}
				panic(x)
			}
		}
	}()
	
	df := state.FetchModule("df")
	if df == nil {
		df = state.RegisterModule("df")
	}
	
	mod := df.RegisterModule("raw")
	
	mod.RegisterCommand("walk", Command_Walk)
	
	mod.RegisterCommand("parse", Command_Parse)
	mod.RegisterCommand("tag", Command_Tag)
	mod.RegisterCommand("dump", Command_Dump)
	
	return nil
}

// Parses the raws and runs code for every tag.
// Return false to abort early, respects breakloop.
// 	df:raw:walk file code
// Parameters for code:
//	tag
// code MUST be a block created via a block declaration!
// Returns the file (as a string).
func Command_Walk(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("df:raw:walk", "2")
	}

	lex := newRawLexer(params[0].String())

	if params[1].Type != rex.TypCode {
		rex.ErrorGeneralCmd("df:raw:walk", "Attempt to run non-executable Value.")
	}

	block := params[1].Data.(*rex.Code)
	script.Locals.Add(block)
	script.Locals.Set(0, rex.NewValueIndex(lex))
	
	for {
		lex.advance()
		if !lex.tagvalid {
			break
		}
		
		script.Exec(block)
		script.BreakLoop = false
		if !script.RetVal.Bool() || script.ExitFlags() {
			break
		}
	}
	script.Locals.Remove()
	
	script.RetVal = rex.NewValueString(string(lex.source))
}

// Parses a raw file.
// 	df:raw:parse text
// This command is much slower, but also far more flexible, than df:raw:walk.
// Returns an indexable of all the tags in the file (each tag is also indexable).
func Command_Parse(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("df:raw:parse", "1")
	}
	
	script.RetVal = rex.NewValueIndex(LexFile([]byte(params[0].String())))
}

// Creates a new raw tag.
// 	df:raw:tag id [params...]
// Returns a tag of the same format as used by df:raw:parse.
func Command_Tag(script *rex.Script, params []*rex.Value) {
	if len(params) < 1 {
		rex.ErrorParamCount("df:raw:tag", ">0")
	}
	
	tag := NewTag()
	tag.ID = params[0].String()
	for _, val := range params[1:] {
		tag.Params.(rex.EditIndexable).Set("append", val)
	}
	script.RetVal = rex.NewValueIndex(tag)
}

// Dumps a parsed raw file to text.
// 	df:raw:dump file
// This command is for use with the indexable returned by df:raw:parse.
// Returns a string version of the file.
func Command_Dump(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("df:raw:dump", "1")
	}

	if params[0].Type != rex.TypIndex {
		rex.ErrorGeneralCmd("df:raw:dump", "Attempt to range over non-Indexable.")
	}
	index := params[0].Data.(rex.Indexable)
	
	out := ""
	for _, i := range index.Keys() {
		if _, ok := index.Get(i).Data.(*Tag); !ok {
			rex.ErrorGeneralCmd("df:raw:dump", "Non-Tag type in Indexable.")
		}
		
		out += index.Get(i).String()
	}
	script.RetVal = rex.NewValueString(out)
}
