/*
Copyright 2012-2014 by Milo Christiansen

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

// Initializes the raw parser.
// 	raw:init text
// Returns a handle to the raw file.
func CommandRaw_Init(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	script.RetVal = raptor.NewValueObject(newRawLexer(params[0].String()))
}

// Advances the parser to the next raw tag.
// 	raw:advance file
// Returns unchanged.
func CommandRaw_Advance(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}
	if _, ok := params[0].Data.(*rawlexer); !ok {
		panic(script.GeneralCmdError("Parameter 0 is not a raw lexer."))
	}
	lex := params[0].Data.(*rawlexer)

	lex.advance()
}

// Returns true if there is a valid tag.
// 	raw:valid file
// Returns true or false.
func CommandRaw_Valid(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}
	if _, ok := params[0].Data.(*rawlexer); !ok {
		panic(script.GeneralCmdError("Parameter 0 is not a raw lexer."))
	}
	lex := params[0].Data.(*rawlexer)

	script.RetVal = raptor.NewValueBool(lex.tagvalid)
}

// Returns an indexable that always points to the current tag.
// You only need to call this once per file.
// 	raw:current file
// Returns the current raw tag.
func CommandRaw_Current(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}
	if _, ok := params[0].Data.(*rawlexer); !ok {
		panic(script.GeneralCmdError("Parameter 0 is not a raw lexer."))
	}
	lex := params[0].Data.(*rawlexer)

	script.RetVal = raptor.NewValueObject(lex)
}

// Disables the current tag.
// 	raw:disable file
// Returns unchanged.
func CommandRaw_Disable(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}
	if _, ok := params[0].Data.(*rawlexer); !ok {
		panic(script.GeneralCmdError("Parameter 0 is not a raw lexer."))
	}
	lex := params[0].Data.(*rawlexer)

	lex.source[lex.tagbegin] = '-'
	lex.source[lex.tagend] = '-'
}

// Replaces the current tag.
// 	raw:replace file replacement
// Returns unchanged.
func CommandRaw_Replace(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic(script.BadParamCount("2"))
	}
	if _, ok := params[0].Data.(*rawlexer); !ok {
		panic(script.GeneralCmdError("Parameter 0 is not a raw lexer."))
	}
	lex := params[0].Data.(*rawlexer)

	newstring := params[1].String()
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
// 	raw:append file text
// Returns unchanged.
func CommandRaw_Append(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 2 {
		panic(script.BadParamCount("2"))
	}
	if _, ok := params[0].Data.(*rawlexer); !ok {
		panic(script.GeneralCmdError("Parameter 0 is not a raw lexer."))
	}
	lex := params[0].Data.(*rawlexer)

	newstring := params[1].String()
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

// Dumps a raw file as a string.
// Use this to write out your edits.
// 	raw:dump file
// Returns the raw text.
func CommandRaw_Dump(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}
	if _, ok := params[0].Data.(*rawlexer); !ok {
		panic(script.GeneralCmdError("Parameter 0 is not a raw lexer."))
	}
	lex := params[0].Data.(*rawlexer)

	script.RetVal = raptor.NewValueString(string(lex.source))
}
