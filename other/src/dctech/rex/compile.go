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

// All of these functions are part of the State, I just put them here to help keep clutter under control.

// CompileToValue compiles code to a TypCode value, value will not be valid if there is an error.
// This is the most commonly used of the compile functions.
func (state *State) CompileToValue(val string, pos *Position) (*Value, error) {
	lex := NewLexer(val, pos)
	code, err := state.Compile(lex)
	return NewValueCode(code), err
}

// CompileShell compiles code with the meta-data from an existing code block and returns the 
// result as both a Value and a Code object, ready to cycle again.
// Use for interactive shells and the like.
func (state *State) CompileShell(val string, code *Code) (value *Value, block *Code, err error) {
	block = NewCodeShell(code)
	lex := NewLexer(val, NewPosition(1, 1, ""))

	err = state.CompileExisting(lex, block)
	
	value = NewValueCode(block)
	return
}

// Compile is as basic as it gets as far as the compile functions are concerned.
// It just compiles code and returns the result directly.
func (state *State) Compile(lex *Lexer) (code *Code, err error) {
	code = NewCode(nil)
	err = state.CompileExisting(lex, code)
	return
}

// CompileExisting is for compiling code when you need to insert variables into the block meta-data first.
// All other compile functions call this one at some level.
func (state *State) CompileExisting(lex *Lexer, code *Code) (err error) {
	err = nil
	defer state.trapErrorCompile(lex, &err)

	for !lex.checkLook(tknINVALID) {
		state.compileValue(lex, code)
	}
	return
}

func (state *State) compileCommand(lex *Lexer, code *Code) {
	lex.getcurrent(tknCmdBegin)
	code.addOp(lex.current.opCode())

	// Compile the command's name
	state.compileName(lex, code, tknCmdBegin)

	// Compile the commands parameters if any
	for !lex.checkLook(tknCmdEnd) {
		state.compileValue(lex, code)
	}

	lex.getcurrent(tknCmdEnd)
	code.addOp(lex.current.opCode())
	return
}

func (state *State) compileVar(lex *Lexer, code *Code) {
	lex.getcurrent(tknVarBegin)
	code.addOp(lex.current.opCode())

	if lex.checkLook(tknRawString) {
		state.compileName(lex, code, tknVarBegin)
	} else {
		state.compileValue(lex, code)
	}
	
	for !lex.checkLook(tknVarEnd, tknAssignment) {
		state.compileValue(lex, code)
	}

	if lex.checkLook(tknAssignment) {
		lex.getcurrent(tknAssignment)
		code.addOp(lex.current.opCode())

		state.compileValue(lex, code)
	}

	lex.getcurrent(tknVarEnd)
	code.addOp(lex.current.opCode())
	return
}

func (state *State) compileObjLit(lex *Lexer, code *Code) {
	lex.getcurrent(tknObjLitBegin)
	code.addOp(lex.current.opCode())

	state.compileName(lex, code, tknObjLitBegin)

	haskeys := false
	if !lex.checkLook(tknObjLitEnd) {
		state.compileValue(lex, code)
		if lex.checkLook(tknAssignment) {
			lex.getcurrent(tknAssignment)
			code.addOp(lex.current.opCode())
			state.compileValue(lex, code)
			haskeys = true
		}

		for !lex.checkLook(tknObjLitEnd) {
			state.compileValue(lex, code)
			if lex.checkLook(tknAssignment) {
				if !haskeys {
					RaiseError("Bad Object Literal: Inconsistent key state.")
				}

				lex.getcurrent(tknAssignment)
				code.addOp(lex.current.opCode())
				state.compileValue(lex, code)
			}
		}
	}

	lex.getcurrent(tknObjLitEnd)
	code.addOp(lex.current.opCode())
	return
}

func (state *State) compileCodeBlock(lex *Lexer, code *Code) {
	block := NewCode(code)
	lex.getcurrent(tknCodeBegin)
	pos := lex.current.Pos.Copy()

	// Nesting is handled automatically, three cheers for recursive decent!
	for !lex.checkLook(tknCodeEnd) {
		state.compileValue(lex, block)
	}
	lex.getcurrent(tknCodeEnd)

	val := NewValueCode(block)
	val.Pos = pos
	code.addOp(&opCode{
		Value: val,
		Type:  opValue,
		Pos:  lex.current.Pos.Copy(),
	})
	return
}

func (state *State) compileCommandBody(lex *Lexer, code *Code) *Value {
	lex.getcurrent(tknCodeBegin)
	pos := lex.current.Pos.Copy()

	// Nesting is handled automatically, three cheers for recursive decent!
	for !lex.checkLook(tknCodeEnd) {
		state.compileValue(lex, code)
	}
	lex.getcurrent(tknCodeEnd)

	val := NewValueCode(code)
	val.Pos = pos
	return val
}

func (state *State) compileBlockDeclare(lex *Lexer, code *Code) *Code {
	params := make([]string, 0, 10)
	defaults := make([]*Value, 0, 10)
	for !lex.checkLook(tknCodeBegin) {
		lex.getcurrent(tknRawString)
		if lex.current.Lexeme == "..." {
			if len(params) != 0 {
				RaiseError("Unexpected variable param marker.")
			}
			params = nil
			break
		}
		params = append(params, lex.current.Lexeme)
		
		// Do we have a default value?
		if lex.checkLook(tknAssignment) {
			lex.getcurrent(tknAssignment)
			lex.getcurrent(tknString, tknTrue, tknFalse, tknNil, tknRawString)
			val := tokenToValue(lex.current)
			defaults = append(defaults, val)
		} else {
			defaults = append(defaults, nil)
		}
	}

	block := NewCode(code)
	if params != nil {
		block.params = len(params)
		for i := range params {
			block.addDefault(params[i], defaults[i])
		}
	} else {
		block.params = -1
		block.add("params")
	}
	return block
}

func (state *State) compileName(lex *Lexer, code *Code, typ int) {
	var module *Module = nil
	
	for {
		lex.getcurrent(tknRawString)

		if lex.checkLook(tknNameSplit) {
			// It's a Module name.
			index := -1
			if module == nil {
				index = state.modules.lookup(lex.current.Lexeme)
				module = state.modules.get(index)
			} else {
				index = module.modules.lookup(lex.current.Lexeme)
				module = module.modules.get(index)
			}

			code.addOp(&opCode{
				Index: index,
				Type:  opName,
				Pos:   lex.current.Pos.Copy(),
			})
			lex.getcurrent(tknNameSplit)
			code.addOp(&opCode{
				Type: opNameSplit,
				Pos:  lex.current.Pos.Copy(),
			})
			continue
		}

		// It's a variable, type or command name.
		index := 0

		switch typ {
		case tknCmdBegin:
			// Command name
			if module == nil {
				// Global
				index = state.global.vars.lookup(lex.current.Lexeme)
				break
			}
			// Module
			index = module.vars.lookup(lex.current.Lexeme)
			break

		case tknVarBegin:
			// Variable name
			if module == nil {
				// Local
				index = code.lookup(lex.current.Lexeme)
				break
			}
			// Module
			index = module.vars.lookup(lex.current.Lexeme)
			break

		case tknObjLitBegin:
			// Type name
			if module == nil {
				// Global
				index = state.types.lookup(lex.current.Lexeme)
				break
			}
			// Module
			index = module.types.lookup(lex.current.Lexeme)
			break
		}

		code.addOp(&opCode{
			Index: index,
			Type:  opName,
			Pos:   lex.current.Pos.Copy(),
		})
		return
	}
}

func (state *State) compileDeclModule(lex *Lexer, code *Code) {
	lex.getcurrent(tknDeclModule)

	var module *Module = nil
	for {
		lex.getcurrent(tknRawString)

		if lex.checkLook(tknNameSplit) {
			// It's a Module name.
			index := 0
			if module == nil {
				index = state.modules.lookup(lex.current.Lexeme)
				module = state.modules.get(index)
			} else {
				index = module.modules.lookup(lex.current.Lexeme)
				module = module.modules.get(index)
			}
			lex.getcurrent(tknNameSplit)
			continue
		}

		// We have found the name we are declaring
		if module == nil {
			if !state.modules.exists(lex.current.Lexeme) {
				state.modules.add(lex.current.Lexeme, newModule())
			}
		} else {
			if !module.modules.exists(lex.current.Lexeme) {
				module.modules.add(lex.current.Lexeme, newModule())
			}
		}
		return
	}
}

func (state *State) compileDeclCommand(lex *Lexer, code *Code) {
	lex.getcurrent(tknDeclCommand)

	var module *Module = nil
	for {
		lex.getcurrent(tknRawString)

		if lex.checkLook(tknNameSplit) {
			// It's a Module name.
			index := 0
			if module == nil {
				index = state.modules.lookup(lex.current.Lexeme)
				module = state.modules.get(index)
			} else {
				index = module.modules.lookup(lex.current.Lexeme)
				module = module.modules.get(index)
			}
			lex.getcurrent(tknNameSplit)
			continue
		}

		// We have found the name we are declaring
		name := lex.current.Lexeme

		// compile the block
		block := state.compileBlockDeclare(lex, code)
		body := state.compileCommandBody(lex, block)
		
		if module == nil {
			state.global.vars.addAndSet(name, body)
			return
		}
		module.vars.addAndSet(name, body)
		return
	}
}

func (state *State) compileDeclBlock(lex *Lexer, code *Code) {
	lex.getcurrent(tknDeclBlock)
	
	block := state.compileBlockDeclare(lex, code)
	body := state.compileCommandBody(lex, block)

	code.addOp(&opCode{
		Value: body,
		Type:  opValue,
		Pos:   lex.current.Pos.Copy(),
	})
}

func (state *State) compileDeclVar(lex *Lexer, code *Code) {
	lex.getcurrent(tknDeclVar)

	// This generates a set expression for the variable.
	var module *Module = nil
	code.addOp(&opCode{
		Type: opVarBegin,
		Pos: lex.current.Pos.Copy(),
	})
	for {
		lex.getcurrent(tknRawString)

		if lex.checkLook(tknNameSplit) {
			// It's a Module name.
			index := 0
			if module == nil {
				index = state.modules.lookup(lex.current.Lexeme)
				module = state.modules.get(index)
			} else {
				index = module.modules.lookup(lex.current.Lexeme)
				module = module.modules.get(index)
			}
			code.addOp(&opCode{
				Type:  opName,
				Index: index,
				Pos:   lex.current.Pos.Copy(),
			})
			lex.getcurrent(tknNameSplit)
			code.addOp(lex.current.opCode())
			continue
		}

		// We have found the name we are declaring
		index := 0
		if module == nil {
			// it's a local!
			index = code.add(lex.current.Lexeme)
		} else {
			// it's a Module variable
			index = module.vars.add(lex.current.Lexeme)
		}

		code.addOp(&opCode{
			Type: opName,
			Index: index,
			Pos: lex.current.Pos.Copy(),
		})
		if lex.checkLook(tknAssignment) {
			lex.getcurrent(tknAssignment)
			code.addOp(lex.current.opCode())
			state.compileValue(lex, code)
		} else {
			code.addOp(&opCode{
				Type: opAssignment,
				Pos:  lex.current.Pos.Copy(),
			})
			code.addOp(&opCode{
				Type:  opValue,
				Value: NewValue(),
				Pos:  lex.current.Pos.Copy(),
			})
		}
		code.addOp(&opCode{
			Type: opVarEnd,
			Pos:  lex.current.Pos.Copy(),
		})
		return
	}
}

func (state *State) compileValue(lex *Lexer, code *Code) {
	switch lex.look.Type {
	case tknString, tknTrue, tknFalse, tknNil, tknRawString:
		lex.getcurrent(tknString, tknTrue, tknFalse, tknNil, tknRawString)
		code.addOp(lex.current.opCode())
		return
	
	case tknCmdBegin:
		state.compileCommand(lex, code)
		return

	case tknVarBegin:
		state.compileVar(lex, code)
		return

	case tknObjLitBegin:
		state.compileObjLit(lex, code)
		return

	case tknCodeBegin:
		state.compileCodeBlock(lex, code)
		return

	case tknDeclModule:
		state.compileDeclModule(lex, code)
		return
	
	case tknDeclCommand:
		state.compileDeclCommand(lex, code)
		return
	
	case tknDeclBlock:
		state.compileDeclBlock(lex, code)
		return
	
	case tknDeclVar:
		state.compileDeclVar(lex, code)
		return
	}

	// This one line handles almost ALL compile time errors, wow.
	exitOntokenExpected(lex.look, tknString, tknTrue, tknFalse, tknNil, tknRawString, tknCmdBegin,
		tknVarBegin, tknObjLitBegin, tknCodeBegin, tknDeclModule, tknDeclCommand, tknDeclBlock, tknDeclVar)
}
