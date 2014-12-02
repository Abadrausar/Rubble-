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

package rex

import "fmt"
import "io"
import "os"

// State handles all global script data and settings.
// The majority of the fields are exported for the use of commands only.
type State struct {
	modules   *moduleStore
	global    *Module
	types     *typeStore
	NoRecover bool      // Do not recover errors, this makes it easier to debug internal errors.
	Output    io.Writer // Normally set to os.Stdout, this can be changed to redirect to a log file or the like.
}

// NewState creates (and initializes) a new state.
func NewState() *State {
	mod := newModuleStore()
	globals := newModule()
	mod.add("global", globals)
	
	return &State{
		modules: mod,
		global: globals,
		types: newTypeStore(),
		Output: os.Stdout,
	}
}

// Host API

// RegisterCommand registers a new native command.
// WARNING! If you try to redeclare a command this function will case a panic!
func (state *State) RegisterCommand(name string, handler NativeCommand) {
	state.global.RegisterCommand(name, handler)
}

// RegisterType registers a new global indexable type.
// WARNING! If you try to redeclare a type this function will case a panic!
func (state *State) RegisterType(name string, typ ObjectFactory) {
	state.types.add(name, typ)
}

// RegisterModule creates a new global module.
// The module is returned so that it can be immediately used to register global data.
// WARNING! If you try to redeclare a module this function will case a panic!
func (state *State) RegisterModule(name string) *Module {
	index := state.modules.add(name, newModule())
	return state.modules.get(index)
}

// FetchModule will retrieve a module by name.
// Returns nil if the module does not exist.
func (state *State) FetchModule(name string) *Module {
	if !state.modules.exists(name) {
		return nil
	}
	return state.modules.get(state.modules.lookup(name))
}

// Output

// Printf is exactly like it's brother from the fmt package, the only difference is 
// that it sends it's output to the State's designated console output writer.
func (state *State) Printf(format string, msg ...interface{}) {
	fmt.Fprintf(state.Output, format, msg...)
}

// Println is exactly like it's brother from the fmt package, the only difference is 
// that it sends it's output to the State's designated console output writer.
func (state *State) Println(msg ...interface{}) {
	fmt.Fprintln(state.Output, msg...)
}

// Print is exactly like it's brother from the fmt package, the only difference is 
// that it sends it's output to the State's designated console output writer.
func (state *State) Print(msg ...interface{}) {
	fmt.Fprint(state.Output, msg...)
}

// Error handling

// Do not try refactoring these to call a single error recovery function!
// Causes "index out of range" errors, as if there is no error then there may be no valid position
// (I think that's what was wrong)

func (state *State) trapErrorRuntime(script *Script, err *error, clean func()) {
	if !state.NoRecover {
		if x := recover(); x != nil {
			if y, ok := x.(ScriptError); ok {
				if !script.code.empty() {
					y.Pos = script.code.last().beginning(script.code.last().current()).Pos.Copy()
				}
				*err = y
			} else if y, ok := x.(error); ok {
				if !script.code.empty() {
					*err = InternalError{Err: y, Pos: script.code.last().beginning(script.code.last().current()).Pos.Copy()}
				} else {
					*err = InternalError{Err: y}
				}
			} else {
				if !script.code.empty() {
					*err = InternalError{Err: fmt.Errorf("%v", x), Pos: script.code.last().beginning(script.code.last().current()).Pos.Copy()}
				} else {
					*err = InternalError{Err: fmt.Errorf("%v", x)}
				}
			}
		}
	}
	
	clean()
}

func (state *State) trapErrorCompile(lex *Lexer, err *error) {
	if !state.NoRecover {
		if x := recover(); x != nil {
			if y, ok := x.(ScriptError); ok {
				y.Pos = lex.current.Pos.Copy()
				*err = y
			} else if y, ok := x.(error); ok {
				*err = InternalError{Err: y, Pos: lex.current.Pos.Copy()}
			} else {
				*err = InternalError{Err: fmt.Errorf("%v", x), Pos: lex.current.Pos.Copy()}
			}
		}
	}
}

// Script Execution

// CompileAndRun is for the lazy, compile and run a script in one step.
// If there is an error ret may not be valid!
func (state *State) CompileAndRun(code, filename string) (ret *Value, err error) {
	script := NewScript()
	
	val, err := state.CompileToValue(code, NewPosition(1, 1, filename))
	if err != nil {
		return
	}
	
	ret, err = state.Run(script, val)
	return
}

// Run is the generic "run script" function, it handles all the overhead of running a block,
// including local variables, error handling, and cleanup.
func (state *State) Run(script *Script, code *Value) (ret *Value, err error) {
	script.Host = state

	err = nil
	defer state.trapErrorRuntime(script, &err, script.ClearMost)
	
	if code.Type != TypCode {
		RaiseError("Attempt to run non-executable Value.")
	}
	block := code.Data.(*Code)
	
	script.Locals.Add(block)
	script.Exec(block)
	ret = script.RetVal
	script.Locals.Remove()
	return
}

// RunShell is mostly like Run, except that it makes an effort to preserve top level local variables.
// The passed in script is expected to still have the root block of locals defined (it is OK if that is not so)
func (state *State) RunShell(script *Script, code *Value) (ret *Value, err error) {
	script.Host = state

	err = nil
	defer state.trapErrorRuntime(script, &err, script.ClearSome)

	if code.Type != TypCode {
		RaiseError("Attempt to run non-executable Value.")
	}
	block := code.Data.(*Code)
	
	script.Locals.RemoveNoClear()
	script.Locals.Add(block)
	script.Exec(block)
	ret = script.RetVal
	return
}

// RunPreped is mostly like Run, except that is does not handle local variables.
// Use when you need to play around with said variables before running.
// Note that unlike RunShell RunPreped does a full cleanup cycle!
func (state *State) RunPreped(script *Script, code *Code) (ret *Value, err error) {
	script.Host = state

	err = nil
	defer state.trapErrorRuntime(script, &err, script.ClearMost)
	
	script.Exec(code)
	ret = script.RetVal
	return
}

// RunCommand runs a script value as a command.
// Does all the normal cleanup, just like Run.
func (state *State) RunCommand(script *Script, code *Value, params []*Value) (ret *Value, err error) {
	script.Host = state

	err = nil
	defer state.trapErrorRuntime(script, &err, script.ClearMost)
	
	if code.Type != TypCode && code.Type != TypCommand {
		RaiseError("Attempt to run non-executable Value.")
	}
	code.call(script, params)
	ret = script.RetVal
	return
}

// RunMinimal is mostly like Run, except that the only thing handled is error recovery.
// No cleanup is done and local variables must be setup before calling.
func (state *State) RunMinimal(script *Script, code *Code) (ret *Value, err error) {
	script.Host = state

	err = nil
	defer state.trapErrorRuntime(script, &err, func(){})
	
	script.Exec(code)
	ret = script.RetVal
	return
}
