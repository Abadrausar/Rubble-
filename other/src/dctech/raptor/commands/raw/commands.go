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

// Raptor DF Raw Commands
package raw

import "dctech/raptor"

// Adds the raw commands to the state.
// The raw commands are:
//	raw:init
//	raw:advance
//	raw:valid
//	raw:current
//	raw:disable
//	raw:replace
//	raw:append
//	raw:dump
func Setup(state *raptor.State) {
	state.NewNameSpace("raw")
	state.NewNativeCommand("raw:init", CommandRaw_Init)
	state.NewNativeCommand("raw:advance", CommandRaw_Advance)
	state.NewNativeCommand("raw:valid", CommandRaw_Valid)
	state.NewNativeCommand("raw:current", CommandRaw_Current)
	state.NewNativeCommand("raw:disable", CommandRaw_Disable)
	state.NewNativeCommand("raw:replace", CommandRaw_Replace)
	state.NewNativeCommand("raw:append", CommandRaw_Append)
	state.NewNativeCommand("raw:dump", CommandRaw_Dump)
}

// Initalizes the raw parser.
// 	raw:init text
// Returns unchanged.
func CommandRaw_Init(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to raw:text.")
	}

	lex = newRawLexer(params[0].String())
}

// Advances the parser to the next raw tag.
// 	raw:advance
// Returns unchanged.
func CommandRaw_Advance(state *raptor.State, params []*raptor.Value) {
	lex.advance()
}

// Returns true if there is a valid tag.
// 	raw:valid
// Returns true or false.
func CommandRaw_Valid(state *raptor.State, params []*raptor.Value) {
	state.RetVal = raptor.NewValueBool(lex.tagvalid)
}

// Returns an indexable that always points to the current tag.
// You only need to call this once.
// 	raw:current
// Returns the current raw tag.
func CommandRaw_Current(state *raptor.State, params []*raptor.Value) {
	state.RetVal = raptor.NewValueObject(lex)
}

// Disables the current tag.
// 	raw:disable
// Returns unchanged.
func CommandRaw_Disable(state *raptor.State, params []*raptor.Value) {
	lex.source[lex.tagbegin] = '-'
	lex.source[lex.tagend] = '-'
}

// Replaces the current tag.
// 	raw:replace replacement
// Returns unchanged.
func CommandRaw_Replace(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to raw:replace.")
	}
	newstring := params[0].String()
	beginstr := lex.source[:lex.tagbegin]
	endstr := lex.source[lex.tagend+1:]

	newlen := len(newstring) + len(beginstr) + len(endstr)

	rtn := make([]byte, 0, newlen)
	rtn = append(rtn, beginstr...)
	rtn = append(rtn, newstring...)
	lex.index = len(rtn) + 1
	rtn = append(rtn, endstr...)
	lex.source = rtn
}

// Adds text after the current tag.
// 	raw:append text
// Returns unchanged.
func CommandRaw_Append(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to raw:append.")
	}
	newstring := params[0].String()
	beginstr := lex.source[:lex.tagend+1]
	endstr := lex.source[lex.tagend+1:]

	newlen := len(newstring) + len(beginstr) + len(endstr)

	rtn := make([]byte, 0, newlen)
	rtn = append(rtn, beginstr...)
	rtn = append(rtn, newstring...)
	lex.index = len(rtn) + 2
	rtn = append(rtn, endstr...)
	lex.source = rtn
}

// Dumps the raw file being opperated on as a string.
// Use this to write out your edits.
// 	raw:dump
// Returns the raw text.
func CommandRaw_Dump(state *raptor.State, params []*raptor.Value) {
	state.RetVal = raptor.NewValueString(string(lex.source))
}
