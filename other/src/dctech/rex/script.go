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

//import "errors"

// Script stores all non-global data and provides an interface to the global data.
// Script.Host is only valid if the script is being run by a State!
type Script struct {
	code   *codeStore
	Locals *LocalValueStore

	RetVal    *Value // The return value of the last command
	Exit      bool   // true when exiting
	Return    bool   // true when a return is active
	Break     bool   // true when a break is active
	BreakLoop bool   // true when a loop break is active
	Error     bool   // set by some commands on error, this is NOT automatically reset!

	Output io.Writer // Normally set to nil (in which case the global value is used)
	Host   *State    // set to nil unless this Script is being run, then it is the host State.
}

// NewScript creates (and initializes) a new script.
func NewScript() *Script {
	return &Script{
		code: newCodeStore(),
		Locals: newLocalValueStore(),
		Output: nil,
		RetVal: NewValue(),
	}
}

// Variables

// SetParams is exactly like calling SetParamsAdv with a nil names.
func (script *Script) SetParams(block *Code, vals []*Value) (err error) {
	return script.SetParamsAdv(block, nil, vals)
}

// SetParamsAdv reads the block meta-data and sets the parameter local variables.
// Values passed in are used first, then any default values defined.
// If there is no default for a slot and all the passed in values are used then an error is generated.
// It is OK to have more values passed in than declared parameters, the extras are ignored.
// names is used for specifying param values by name, setting this to nil means that all params are by position.
func (script *Script) SetParamsAdv(block *Code, names []string, vals []*Value) (err error) {
	err = nil
	defer script.trapErrorRuntime(&err, func(){})
	
	if names != nil && len(names) != len(vals) {
		RaiseError("Parameter name list and parameter value list are not the same length!")
	}
	
	script.setParams(block, names, vals)
	return
}

func (script *Script) setParams(block *Code, names []string, vals []*Value) {
	if block.params == nil {
		// The block has no parameters, so obviously nothing has to be done...
		return
	}
	
	// Handle any named params.
	if names != nil {
		tmp1 := make([]string, 0, len(names))
		tmp2 := make([]*Value, 0, len(vals))
		for i := range names {
			if names[i] != "" {
				k := block.params.itov[block.params.ntoi[names[i]]]
				script.Locals.Set(k, vals[i])
				continue
			}
			tmp1 = append(tmp1, names[i])
			tmp2 = append(tmp2, vals[i])
		}
		names = tmp1
		vals = tmp2
	}
	
	// Now fill everything else by position.
	for i, k := range block.params.itov {
		if len(vals) <= i && block.params.defaults[i] == nil {
			RaiseError("Block declaration and param count do not match.")
		}
		
		// For the last parameter handle a variadic value if needed.
		if i == len(block.params.itov) - 1 {
			if block.params.defaults[i] != nil && block.params.defaults[i].Type == TypVariadic {
				if len(vals) <= i {
					if script.Locals.Get(k) == nil {
						script.Locals.Set(k, NewValueIndex(NewStaticArray([]*Value{})))
					}
					continue
				}
				script.Locals.Set(k, NewValueIndex(NewStaticArray(vals[i:])))
				continue
			}
		}
		
		if len(vals) <= i {
			script.Locals.Set(k, block.params.defaults[i])
			continue
		}
		script.Locals.Set(k, vals[i])
	}
}

// Output

// Printf writes to the script output if set, else writes to the host output if the host is set, else defaults to Stdout.
func (script *Script) Printf(format string, msg ...interface{}) {
	if script.Output == nil {
		if script.Host == nil {
			fmt.Fprintf(os.Stdout, format, msg...)
			return
		}

		fmt.Fprintf(script.Host.Output, format, msg...)
		return
	}
	fmt.Fprintf(script.Output, format, msg...)
}

// Println writes to the script output if set, else writes to the host output if the host is set, else defaults to Stdout.
func (script *Script) Println(msg ...interface{}) {
	if script.Output == nil {
		if script.Host == nil {
			fmt.Fprintln(os.Stdout, msg...)
			return
		}

		fmt.Fprintln(script.Host.Output, msg...)
		return
	}
	fmt.Fprintln(script.Output, msg...)
}

// Print writes to the script output if set, else writes to the host output if the host is set, else defaults to Stdout.
func (script *Script) Print(msg ...interface{}) {
	if script.Output == nil {
		if script.Host == nil {
			fmt.Fprint(os.Stdout, msg...)
			return
		}

		fmt.Fprint(script.Host.Output, msg...)
		return
	}
	fmt.Fprint(script.Output, msg...)
}

// Other

func (script *Script) trapErrorRuntime(err *error, clean func()) {
	if script.Host != nil {
		if script.Host.NoRecover {
			clean()
			return
		}
	}
	
	if x := recover(); x != nil {
		if y, ok := x.(ScriptError); ok {
			if !script.code.empty() {
				y.Pos = script.code.last().current().Pos
			}
			*err = y
		} else if y, ok := x.(error); ok {
			if !script.code.empty() {
				*err = InternalError{Err: y, Pos: script.code.last().current().Pos}
			} else {
				*err = InternalError{Err: y}
			}
		} else {
			if !script.code.empty() {
				*err = InternalError{Err: fmt.Errorf("%v", x), Pos: script.code.last().current().Pos}
			} else {
				*err = InternalError{Err: fmt.Errorf("%v", x)}
			}
		}
	}
	
	clean()
}

// ExitFlags returns true if any of the exit flags are set.
// The exit flags are Exit, Return, Break, and BreakLoop.
func (script *Script) ExitFlags() bool {
	return (script.Exit || script.Return || script.Break || script.BreakLoop)
}

// Clears all registers, removes all locals, and removes all code.
func (script *Script) Clear() {
	script.ClearMost()
	script.RetVal = NewValue()
}

// ClearMost clears all registers, removes all locals, and removes all code.
// Does not modify the RetVal register.
func (script *Script) ClearMost() {
	script.Locals.Clear()
	script.code.clear()
	script.Host = nil // this should already be clear, but just in case...
	script.Exit = false
	script.Return = false
	script.Break = false
	script.BreakLoop = false
}

// ClearSome clears all registers, removes all locals, and removes all code.
// Does not modify the RetVal register.
// Leaves the first block of locals.
func (script *Script) ClearSome() {
	script.Locals.ClearToRoot()
	script.code.clear()
	script.Host = nil // this should already be clear, but just in case...
	script.Exit = false
	script.Return = false
	script.Break = false
	script.BreakLoop = false
}

// Script Execution

// Exec runs script code.
// Exec does NOT do any error recovery and so it should only be called 
// down-stream from one of the Run functions in the State!
func (script *Script) Exec(code *Code) {
	script.code.add(newCodeReader(code))

	// In case of weird failures uncomment next two lines lines.
	//script.Println("opCode Dump:")
	//script.Println(code)
	
	// Run code
	for !script.code.last().checkLookAhead(opINVALID) {
		script.execValue()
		if script.ExitFlags() {
			script.Break = false
			break
		}
	}

	script.code.remove()
}

// SafeExec runs script code.
// SafeExec does error recovery normally reserved for the State Run functions,
// but does none of the other cleanup that the Run functions do.
// Use for sub-shells and the like.
// Should only be called down-stream from one of the Run functions in the State!
func (script *Script) SafeExec(code *Code) (err error) {
	
	err = nil
	defer script.trapErrorRuntime(&err, func(){})
	
	if script.Host == nil {
		RaiseError("SafeExec called with nil Host.")
	}
	
	script.code.add(newCodeReader(code))

	// Run code
	for !script.code.last().checkLookAhead(opINVALID) {
		script.execValue()
		if script.ExitFlags() {
			script.Break = false
			break
		}
	}

	script.code.remove()
	return
}

func (script *Script) execCommand() {
	script.code.last().getOpCode(opCmdBegin)

	// Read the command's name
	module := script.Host.global
	var command *Value = nil
	if script.code.last().checkLookAhead(opName) {
		for {
			script.code.last().getOpCode(opName)
	
			if script.code.last().checkLookAhead(opNameSplit) {
				module = module.modules.get(script.code.last().current().Index)
				script.code.last().getOpCode(opNameSplit)
				continue
			}

			command = module.vars.get(script.code.last().current().Index)
			break
		}
	} else {
		command = script.execValue()
	}

	// Read the command's parameters if any
	params := make([]*Value, 0, 5)
	names := make([]string, 0, 5)
	realname := false
	for !script.code.last().checkLookAhead(opCmdEnd) {
		v := script.execValue()
		if script.ExitFlags() {
			return
		}
		
		if script.code.last().checkLookAhead(opAssignment) {
			script.code.last().getOpCode(opAssignment)
			realname = true
			names = append(names, v.String())
			params = append(params, script.execValue())
			if script.ExitFlags() {
				return
			}
			continue
		}
		names = append(names, "")
		params = append(params, v)
	}

	script.code.last().getOpCode(opCmdEnd)

	if !realname {
		names = nil
	}
	
	command.call(script, names, params)
	return
}

func (script *Script) execVar() {
	script.code.last().getOpCode(opVarBegin)
	
	var val *Value = nil
	module := script.Host.global
	ismod := false
	if script.code.last().checkLookAhead(opName) {
		// Handle the name form
		index := 0
		for {
			script.code.last().getOpCode(opName)
	
			if script.code.last().checkLookAhead(opNameSplit) {
				// It's a module name.
				module = module.modules.get(script.code.last().current().Index)
				ismod = true
	
				script.code.last().getOpCode(opNameSplit)
				continue
			}
	
			index = script.code.last().current().Index
			break
		}
		
		// Non-indexing derefs and sets should only occur with the name form.
		if script.code.last().checkLookAhead(opAssignment) {
			script.code.last().getOpCode(opAssignment)
			script.RetVal = script.execValue()
			if script.ExitFlags() {
				return
			}
			if ismod {
				module.vars.set(index, script.RetVal)
			} else {
				i := script.Locals.IsGlobal(index)
				if i > 0 {
					module.vars.set(i, script.RetVal)
				} else {
					script.Locals.Set(index, script.RetVal)
				}
			}
			script.code.last().getOpCode(opVarEnd)
			return
		}
		if script.code.last().checkLookAhead(opVarEnd) {
			if ismod {
				script.RetVal = module.vars.get(index)
			} else {
				i := script.Locals.IsGlobal(index)
				if i > 0 {
					script.RetVal = module.vars.get(i)
				} else {
					script.RetVal = script.Locals.Get(index)
				}
			}
			script.code.last().getOpCode(opVarEnd)
			return
		}
		
		// It's an indexing form, store the value
		if ismod {
			val = module.vars.get(index)
		} else {
			i := script.Locals.IsGlobal(index)
			if i > 0 {
				val = module.vars.get(i)
			} else {
				val = script.Locals.Get(index)
			}
		}
		
	} else {
		// Handle the value form
		val = script.execValue()
	}
	
	// Read indexing parameters
	params := make([]*Value, 0, 5)
	for {
		params = append(params, script.execValue())
		if script.ExitFlags() {
			return
		}
		if script.code.last().checkLookAhead(opAssignment, opVarEnd) {
			break
		}
	}

	for i := range params[:len(params)-1] {
		if val.Type != TypIndex {
			RaiseError(fmt.Sprint("Non-Indexable object passed to an indexing variable operator. Indexing level:", i))
		}
		if v, ok := val.Data.(IntIndexable); ok && params[i].Type == TypInt {
			val = v.GetI(params[i].Int64())
		} else {
			val = val.Data.(Indexable).Get(params[i].String())
		}
	}

	if val.Type != TypIndex {
		RaiseError("Non-Indexable object passed to an indexing variable operator (last element).")
	}
	
	if script.code.last().checkLookAhead(opAssignment) {
		script.code.last().getOpCode(opAssignment)
		script.RetVal = script.execValue()
		if script.ExitFlags() {
			return
		}

		i := params[len(params)-1]
		if v, ok := val.Data.(IntEditIndexable); ok && i.Type == TypInt {
			r := v.SetI(i.Int64(), script.RetVal)
			if !r {
				RaiseError("Indexing variable set failed (Set returned false).")
			}
			script.RetVal = val
		} else {
			index, ok := val.Data.(EditIndexable)
			if !ok {
				RaiseError("Attempt to use read only value with indexing form variable set operator.")
			}
			r := index.Set(i.String(), script.RetVal)
			if !r {
				RaiseError("Indexing variable set failed (Set returned false).")
			}
			script.RetVal = val
		}
	} else {
		i := params[len(params)-1]
		if v, ok := val.Data.(IntIndexable); ok && i.Type == TypInt {
			script.RetVal = v.GetI(i.Int64())
		} else {
			script.RetVal = val.Data.(Indexable).Get(i.String())
		}
	}

	script.code.last().getOpCode(opVarEnd)
}

func (script *Script) execObjLit() {
	script.code.last().getOpCode(opObjLitBegin)

	// Read the type name
	module := script.Host.global
	var typ ObjectFactory = nil
	for {
		script.code.last().getOpCode(opName)

		if script.code.last().checkLookAhead(opNameSplit) {
			module = module.modules.get(script.code.last().current().Index)
			script.code.last().getOpCode(opNameSplit)
			continue
		}

		typ = module.types.get(script.code.last().current().Index)
		break
	}

	haskeys := false
	keys := make([]string, 0, 20)
	values := make([]*Value, 0, 20)
	if !script.code.last().checkLookAhead(opObjLitEnd) {
		tmp := script.execValue()
		if script.code.last().checkLookAhead(opAssignment) {
			keys = append(keys, tmp.String())
			if script.ExitFlags() {
				return
			}
			script.code.last().getOpCode(opAssignment)
			values = append(values, script.execValue())
			if script.ExitFlags() {
				return
			}
			haskeys = true
		} else {
			values = append(values, tmp)
			if script.ExitFlags() {
				return
			}
		}

		for !script.code.last().checkLookAhead(opObjLitEnd) {
			if haskeys {
				keys = append(keys, script.execValue().String())
				if script.ExitFlags() {
					return
				}
				script.code.last().getOpCode(opAssignment)
				values = append(values, script.execValue())
				if script.ExitFlags() {
					return
				}
			} else {
				values = append(values, script.execValue())
				if script.ExitFlags() {
					return
				}
			}
		}
	}
	if !haskeys {
		keys = nil
	}

	script.RetVal = typ(script, keys, values)
	script.code.last().getOpCode(opObjLitEnd)
}

func (script *Script) execValue() *Value {
	switch script.code.last().lookAhead().Type {
	case opValue:
		script.code.last().getOpCode(opValue)
		script.RetVal = script.code.last().current().Value
		return script.RetVal

	case opCmdBegin:
		script.execCommand()
		return script.RetVal

	case opVarBegin:
		script.execVar()
		return script.RetVal

	case opObjLitBegin:
		script.execObjLit()
		return script.RetVal
	}

	exitOnopCodeExpected(script.code.last().lookAhead(), opValue, opCmdBegin, opVarBegin, opObjLitBegin)
	panic("UNREACHABLE")
}
