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
	meta     *blockMeta     // Block variable information
	params   *blockParams   // Block parameter information (may be nil)
}

// NewCode creates a new empty Code block, ready to use as a compiler target.
func NewCode() *Code {
	block := new(Code)
	block.data = make([]*opCode, 0, 100)
	block.meta = newBlockMeta()
	return block
}

// NewCodeShell creates a new empty Code block with the same initial meta-data as the source.
func NewCodeShell(source *Code) *Code {
	if source == nil {
		return NewCode()
	}
	
	block := new(Code)
	block.data = make([]*opCode, 0, 100)
	block.meta = source.meta.dup()
	if source.params != nil {
		block.params = source.params.dup()
	}
	return block
}

// addOp adds an opCode to the block.
func (block *Code) addOp(op *opCode) {
	block.data = append(block.data, op)
}

// String prints all the opcodes in the block.
func (block *Code) String() string {
	out := ""
	for _, op := range block.data {
		out += " " + op.String()
	}
	return out
}

func (block *Code) addParam(name string, val *Value) {
	if block.params == nil {
		block.params = newBlockParams()
	}
	
	i := block.meta.add(name)
	block.params.add(name, val, i)
}

// ===========================================================================================

type blockMeta struct {
	ntoi     map[string]int // local variables offset from the block start, lookup by name.
	iton     []string       // local variable names, lookup by offset from the block start.
	
	entoi    map[string]int // external local variables offset from the block start, lookup by name.
	eiton    []string       // external local variable names, lookup by offset from the block start.
}

func newBlockMeta() *blockMeta {
	meta := new(blockMeta)
	meta.ntoi = make(map[string]int, 10)
	meta.iton = make([]string, 0, 10)
	meta.entoi = make(map[string]int, 10)
	meta.eiton = make([]string, 0, 10)
	return meta
}

func (source *blockMeta) dup() *blockMeta {
	meta := new(blockMeta)
	meta.ntoi = make(map[string]int, len(source.ntoi))
	meta.iton = make([]string, len(source.iton))
	meta.entoi = make(map[string]int, len(source.entoi))
	meta.eiton = make([]string, len(source.eiton))
	
	for i := range source.ntoi {
		meta.ntoi[i] = source.ntoi[i]
	}
	copy(meta.iton, source.iton)
	
	for i := range source.entoi {
		meta.entoi[i] = source.entoi[i]
	}
	copy(meta.eiton, source.eiton)
	return meta
}

// lookup returns the local index of a local or external variable by name.
func (meta *blockMeta) lookup(name string) int {
	if i, ok := meta.ntoi[name]; ok {
		return i
	}

	if i, ok := meta.entoi[name]; ok {
		return 0 - (i + 1)
	}

	index := len(meta.eiton)
	meta.entoi[name] = index
	meta.eiton = append(meta.eiton, name)
	return 0 - (index + 1)
}

// add is used for adding meta-data for a new local variable.
func (meta *blockMeta) add(name string) int {
	if _, ok := meta.ntoi[name]; ok {
		RaiseError("Local Value slot already exists.")
	}

	index := len(meta.iton)
	meta.ntoi[name] = index
	meta.iton = append(meta.iton, name)
	return index
}

// ===========================================================================================

type blockParams struct {
	ntoi     map[string]int // Name to param index
	iton     []string       // index to param name
	itov     []int          // index to local variable index
	defaults []*Value       // Parameter default values (non-optional params have a nil value here)
}

func newBlockParams() *blockParams {
	meta := new(blockParams)
	meta.ntoi = make(map[string]int, 5)
	meta.iton = make([]string, 0, 5)
	meta.itov = make([]int, 0, 5)
	meta.defaults = make([]*Value, 0, 5)
	return meta
}

func (source *blockParams) dup() *blockParams {
	meta := new(blockParams)
	meta.ntoi = make(map[string]int, len(source.ntoi))
	meta.iton = make([]string, len(source.iton))
	meta.itov = make([]int, len(source.itov))
	
	for i := range source.ntoi {
		meta.ntoi[i] = source.ntoi[i]
	}
	copy(meta.iton, source.iton)
	copy(meta.itov, source.itov)
	return meta
}

func (meta *blockParams) add(name string, val *Value, vIndex int) {
	if _, ok := meta.ntoi[name]; ok {
		RaiseError("Parameter slot already exists.")
	}
	
	meta.ntoi[name] = len(meta.iton)
	meta.iton = append(meta.iton, name)
	meta.itov = append(meta.itov, vIndex)
	meta.defaults = append(meta.defaults, val)
}

// ===========================================================================================

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
	if code.index + 1 < code.length {
		return code.data.data[code.index + 1]
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

// beginning finds the beginning tag of a construct based on it's end tag.
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
