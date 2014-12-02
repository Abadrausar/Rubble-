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

//import "fmt"

// Code is a block of opCodes with meta-data.
type Code struct {
	data     []*opCode      // The actual code...
	parent   *Code          // This code block's containing block, used for name lookup.
	ntoi     map[string]int // local variables offset from the block start, lookup by name.
	iton     []string       // local variable names, lookup by offset from the block start.
	defaults []*Value       // Default values for variables in this block, normally each slot is nil.
	
	// This is normally 0, but for blocks created with block or command it contains 
	// the parameter count (-1 for variable, eg a single sarray named params).
	params   int

	entoi    map[string]int // external local variables offset from the block start, lookup by name.
	eiton    []string       // external local variable names, lookup by offset from the block start.
}

// NewCode creates an empty Code block, ready to use as a compiler target,
func NewCode(parent *Code) *Code {
	code := new(Code)
	code.data = make([]*opCode, 0, 100)
	code.parent = parent
	code.ntoi = make(map[string]int, 10)
	code.iton = make([]string, 0, 10)
	code.defaults = make([]*Value, 0, 10)
	code.entoi = make(map[string]int, 10)
	code.eiton = make([]string, 0, 10)
	return code
}

// NewCodeShell creates new Code block with the meta-data from source.
// Use for things like interactive shells where variables need to stick around.
func NewCodeShell(source *Code) *Code {
	code := new(Code)
	code.data = make([]*opCode, 0, 100)
	if source == nil {
		code.parent = nil
		code.ntoi = make(map[string]int, 10)
		code.iton = make([]string, 0, 10)
		code.defaults = make([]*Value, 0, 10)
		code.entoi = make(map[string]int, 10)
		code.eiton = make([]string, 0, 10)
	} else {
		code.parent = source.parent
		code.ntoi = make(map[string]int, len(source.ntoi))
		code.iton = make([]string, len(source.iton))
		code.defaults = make([]*Value, len(source.defaults))
		code.entoi = make(map[string]int, len(source.entoi))
		code.eiton = make([]string, len(source.eiton))
		
		for i := range source.ntoi {
			code.ntoi[i] = source.ntoi[i]
		}
		copy(code.iton, source.iton)
		
		copy(code.defaults, source.defaults)
		
		for i := range source.entoi {
			code.entoi[i] = source.entoi[i]
		}
		copy(code.eiton, source.eiton)
	}
	return code
}

// addOp adds an opCode to the block.
func (code *Code) addOp(op *opCode) {
	code.data = append(code.data, op)
}

// String prints all the opcodes in the code block.
func (code *Code) String() string {
	out := ""
	for _, op := range code.data {
		out += " " + op.String()
	}
	return out
}

// lookup returns the local index of a local or external variable by name.
func (code *Code) lookup(name string) int {
	if i, ok := code.ntoi[name]; ok {
		return i
	}

	if i, ok := code.entoi[name]; ok {
		return 0 - (i + 1)
	}

	index := len(code.eiton)
	code.entoi[name] = index
	code.eiton = append(code.eiton, name)
	return 0 - (index + 1)
}

// add is used for adding meta-data for a new local variable.
func (code *Code) add(name string) int {
	if _, ok := code.ntoi[name]; ok {
		RaiseError("Local Value slot already exists.")
	}

	index := len(code.iton)
	code.ntoi[name] = index
	code.iton = append(code.iton, name)
	code.defaults = append(code.defaults, nil)
	return index
}

// addDefault is used for adding meta-data for a new local variable that has a default value.
// (in practice this is used for parameters only)
func (code *Code) addDefault(name string, def *Value) int {
	if _, ok := code.ntoi[name]; ok {
		RaiseError("Local Value slot already exists.")
	}

	index := len(code.iton)
	code.ntoi[name] = index
	code.iton = append(code.iton, name)
	code.defaults = append(code.defaults, def)
	return index
}

// Code Reader

// codeReader is a thin wrapper to make stepping over opCodes in a Code block easy.
type codeReader struct {
	data   *Code
	index  int
	length int
}

// newCodeReader returns a new codeReader object pointing to just before the first opCode in the block.
func newCodeReader(code *Code) *codeReader {
	return &codeReader{
		data:   code,
		index:  -1,
		length: len(code.data),
	}
}

// advance the codeReader one opCode.
func (code *codeReader) advance() {
	code.index++
}

// current returns the current opCode.
func (code *codeReader) current() *opCode {
	if code.index < code.length {
		return code.data.data[code.index]
	}
	return &opCode{Type: opINVALID}
}

// lookAhead returns the next opCode.
func (code *codeReader) lookAhead() *opCode {
	if code.index+1 < code.length {
		return code.data.data[code.index+1]
	}
	return &opCode{Type: opINVALID}
}

// getOpCode advances the codeReader then checks the current type against the types passed in,
// If the current type does not match one of the passed in types it panics with an error.
func (code *codeReader) getOpCode(types ...int) {
	code.advance()

	for _, val := range types {
		if code.current().Type == val {
			return
		}
	}

	exitOnopCodeExpected(code.current(), types...)
}

// checkLookAhead is like getOpCode except it checks the look ahead instead of the current opCode
// and it returns a boolean instead of panicking.
func (code *codeReader) checkLookAhead(types ...int) bool {
	for _, val := range types {
		if code.lookAhead().Type == val {
			return true
		}
	}
	return false
}

// lookAhead returns the next opCode.
func (code *codeReader) beginning(end *opCode) *opCode {
	start := -1 
	switch end.Type {
	case opCmdEnd:
		start = opCmdBegin
	case opVarEnd:
		start = opVarBegin
	case opObjLitEnd:
		start = opObjLitBegin
	default:
		return end
	}
	
	if code.index >= code.length {
		// Error, but we don't want to raise an error here
		// as this is to be called from an error handler.
		return end
	}
	
	depth := 0
	for i := code.index - 1; i >= 0; i-- {
		if code.data.data[i].Type == end.Type {
			depth++
		}
		if code.data.data[i].Type == start {
			if depth > 0 {
				depth--
			} else {
				return code.data.data[i]
			}
		}
	}
	return end
}

// Code Store

// codeStore is a stack of codeReaders.
// The only place this is used is in the Script.
type codeStore []*codeReader

// newCodeStore creates a new empty codeStore, ready to use.
func newCodeStore() *codeStore {
	rtn := new(codeStore)
	*rtn = make([]*codeReader, 0, 20)
	return rtn
}

// add a codeReader to the stack.
func (code *codeStore) add(val *codeReader) {
	*code = append(*code, val)
}

// remove a codeReader from the stack.
func (code *codeStore) remove() *codeReader {
	rtn := (*code)[len(*code)-1]
	(*code)[len(*code)-1] = nil
	*code = (*code)[:len(*code)-1]
	return rtn
}

// clear the codeStore.
func (code *codeStore) clear() {
	*code = (*code)[0:0]
}

// Return the TOS codeReader.
func (code *codeStore) last() *codeReader {
	return (*code)[len(*code)-1]
}

func (code *codeStore) empty() bool {
	return len(*code) == 0
}
