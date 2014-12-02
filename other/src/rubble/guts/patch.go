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

package guts

import "dctech/rex"
import "dctech/patch"

import "fmt"
import "strings"
import "bytes"
import "io/ioutil"

import "compress/flate"
import "encoding/base64"

// Script commands related to patching and loading addons from scripts.

func InitScriptingPatch() {
	rbl := GlobalScriptState.FetchModule("rubble")

	rbl.RegisterCommand("patch", Command_Patch)
	rbl.RegisterCommand("decompress", Command_Decompress)
	rbl.RegisterCommand("compress", Command_Compress)
}

// Applies a patch to a string.
// 	rubble:patch string patch
// Returns the patched text.
func Command_Patch(script *rex.Script, params []*rex.Value) {
	if len(params) != 2 {
		rex.ErrorParamCount("rubble:patch", "2")
	}

	patches := patch.FromText(params[1].String())
	if patches == nil {
		rex.ErrorGeneralCmd("rubble:patch", "Invalid patch string.")
	}
	
	text, applied := patch.Apply(patches, params[0].String())
	
	count := 0
	for i := range applied {
		if applied[i] {
			count++
		}
	}
	if len(patches) != count {
		rex.ErrorGeneralCmd("rubble:patch", fmt.Sprintf("Not all patches applied: %v out of %v", count, len(patches)))
	}
	
	script.RetVal = rex.NewValueString(text)
}

// Decompresses a Rubble encoded string.
// 	rubble:decompress string
// Returns the decoded text.
func Command_Decompress(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("rubble:decompress", "1")
	}

	a := strings.NewReader(params[0].String())
	ac := base64.NewDecoder(base64.StdEncoding, a)
	datac, err := ioutil.ReadAll(ac)
	if err != nil {
		rex.ErrorGeneralCmd("rubble:decompress", err.Error())
	}
	
	b := bytes.NewReader(datac)
	bc := flate.NewReader(b)
	data, err := ioutil.ReadAll(bc)
	if err != nil {
		rex.ErrorGeneralCmd("rubble:decompress", err.Error())
	}
	bc.Close()
	
	script.RetVal = rex.NewValueString(string(data))
}

// Compresses a string.
// 	rubble:compress string
// Returns the encoded text.
func Command_Compress(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("rubble:compress", "1")
	}

	a := new(bytes.Buffer)
	ac, _ := flate.NewWriter(a, 9)
	_, err := ac.Write([]byte(params[0].String()))
	if err != nil {
		rex.ErrorGeneralCmd("rubble:compress", err.Error())
	}
	ac.Close()
	
	b := new(bytes.Buffer)
	bc := base64.NewEncoder(base64.StdEncoding, b)
	_, err = bc.Write(a.Bytes())
	if err != nil {
		rex.ErrorGeneralCmd("rubble:compress", err.Error())
	}
	bc.Close()
	
	script.RetVal = rex.NewValueString(b.String())
}
