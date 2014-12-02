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

package base_test

import "testing"
import "bytes"

import "dctech/rex"
import "dctech/rex/commands/base"

// Run code and return the result + the console output.
func testCommand(name, test string, t *testing.T, code string) (*rex.Value, string) {
	// New state each time so tests won't mess each other up.
	state := rex.NewState()
	base.Setup(state)
	
	// Redirect all output to a buffer.
	output := new(bytes.Buffer)
	state.Output = output
	
	script := rex.NewScript()
	val, err := state.CompileToValue(code, rex.NewPosition(1, 1, ""))
	if err != nil {
		t.Error("Command: " + name + " Test: " + test + " Compile Error: ", err)
	}
	rtn, err := state.Run(script, val)
	if err != nil {
		t.Error("Command: " + name + " Test: " + test + " Runtime Error: ", err)
	}
	
	return rtn, output.String()
}

func returnError(name, test, err string) string {
	return "Command: " + name + " Test: " + test + " Return Error: " + err
}

func outputError(name, test, err string) string {
	return "Command: " + name + " Test: " + test + " Output Error: " + err
}

// Now for the tests...

// Really nothing to test, but here for consistency.
func Test_Nop(t *testing.T) {
	name := "nop"
	
	rtn, _ := testCommand(name, "1", t, "true (nop)")
	if rtn.Type != rex.TypBool {
		t.Error(returnError(name, "1", "Type not bool."))
	}
	if !rtn.Bool() {
		t.Error(returnError(name, "1", "Value not true."))
	}
}

func Test_Ret(t *testing.T) {
	name := "ret"
	
	rtn, _ := testCommand(name, "1", t, "(ret true)")
	if rtn.Type != rex.TypBool {
		t.Error(returnError(name, "1", "Type not bool."))
	}
	if !rtn.Bool() {
		t.Error(returnError(name, "1", "Value not true."))
	}
	
	rtn, _ = testCommand(name, "2", t, "(ret 'abc')")
	if rtn.Type != rex.TypString {
		t.Error(returnError(name, "2", "Type not string."))
	}
	if rtn.String() != "abc" {
		t.Error(returnError(name, "2", "Value not 'abc'."))
	}
	
	rtn, _ = testCommand(name, "3", t, `
	command b {(ret true) false}
	command a {false (b)}
	(a)`)
	if rtn.Type != rex.TypBool {
		t.Error(returnError(name, "3", "Type not bool."))
	}
	if !rtn.Bool() {
		t.Error(returnError(name, "3", "Value not true."))
	}
}

func Test_Break(t *testing.T) {
	name := "break"
	
	rtn, _ := testCommand(name, "1", t, "(break true)")
	if rtn.Type != rex.TypBool {
		t.Error(returnError(name, "1", "Type not bool."))
	}
	if !rtn.Bool() {
		t.Error(returnError(name, "1", "Value not true."))
	}
	
	rtn, _ = testCommand(name, "2", t, "(break 'abc')")
	if rtn.Type != rex.TypString {
		t.Error(returnError(name, "2", "Type not string."))
	}
	if rtn.String() != "abc" {
		t.Error(returnError(name, "2", "Value not 'abc'."))
	}
	
	rtn, _ = testCommand(name, "3", t, `
	command b {(break true) false}
	command a {false (b)}
	(a)`)
	if rtn.Type != rex.TypBool {
		t.Error(returnError(name, "3", "Type not bool."))
	}
	if !rtn.Bool() {
		t.Error(returnError(name, "3", "Value not true."))
	}
}

func Test_BreakLoop(t *testing.T) {
	name := "breakloop"
	
	rtn, _ := testCommand(name, "1", t, "(breakloop true)")
	if rtn.Type != rex.TypBool {
		t.Error(returnError(name, "1", "Type not bool."))
	}
	if !rtn.Bool() {
		t.Error(returnError(name, "1", "Value not true."))
	}
	
	rtn, _ = testCommand(name, "2", t, "(breakloop 'abc')")
	if rtn.Type != rex.TypString {
		t.Error(returnError(name, "2", "Type not string."))
	}
	if rtn.String() != "abc" {
		t.Error(returnError(name, "2", "Value not 'abc'."))
	}
	
	rtn, _ = testCommand(name, "3", t, `
	command b {(breakloop true) false}
	command a {false (b)}
	(a)`)
	if rtn.Type != rex.TypBool {
		t.Error(returnError(name, "3", "Type not bool."))
	}
	if !rtn.Bool() {
		t.Error(returnError(name, "3", "Value not true."))
	}
	
	rtn, _ = testCommand(name, "4", t, `
	(loop {
		(breakloop false)
		0
	})`)
	if rtn.Type != rex.TypBool {
		t.Error(returnError(name, "4", "Type not bool."))
	}
	if rtn.Bool() {
		t.Error(returnError(name, "4", "Value not false."))
	}
}

func Test_Eval(t *testing.T) {
	name := "eval"
	
	rtn, _ := testCommand(name, "1", t, "(eval {true})")
	if rtn.Type != rex.TypBool {
		t.Error(returnError(name, "1", "Type not bool."))
	}
	if !rtn.Bool() {
		t.Error(returnError(name, "1", "Value not true."))
	}
}

func Test_Error(t *testing.T) {
	name := "error"
	
	rtn, _ := testCommand(name, "1", t, "(error -1) (error)")
	if rtn.Type != rex.TypBool {
		t.Error(returnError(name, "1", "Type not bool."))
	}
	if !rtn.Bool() {
		t.Error(returnError(name, "1", "Value not true."))
	}
	
	rtn, _ = testCommand(name, "2", t, "(error)")
	if rtn.Type != rex.TypBool {
		t.Error(returnError(name, "2", "Type not bool."))
	}
	if rtn.Bool() {
		t.Error(returnError(name, "2", "Value not false."))
	}
}

// TODO: Tests for all other commands
