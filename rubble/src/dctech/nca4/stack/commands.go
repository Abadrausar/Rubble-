// NCA v4 Stack Commands.
package stack

import "dctech/nca4"

type stack struct {
	ptr  int
	data []*nca4.Value
}

func newstack() *stack {
	rtn := new(stack)
	rtn.data = make([]*nca4.Value, 25, 25)
	return rtn
}

var stacks []*stack

func init() {
	stacks = make([]*stack, 0, 5)
}

// Adds the stack commands to the state.
// The stack commands are:
//	stack:new
//	stack:grow
//	stack:push
//	stack:pop
//	stack:top
func Setup(state *nca4.State) {
	state.NewNameSpace("stack")
	state.NewNativeCommand("stack:new", CommandStack_New)
	state.NewNativeCommand("stack:grow", CommandStack_Grow)
	state.NewNativeCommand("stack:push", CommandStack_Push)
	state.NewNativeCommand("stack:pop", CommandStack_Pop)
	state.NewNativeCommand("stack:top", CommandStack_Top)
}

// Create a new stack of size 25.
// 	stack:new
// Returns the stack handle
func CommandStack_New(state *nca4.State, params []*nca4.Value) {
	handle := len(stacks)

	stacks = append(stacks, newstack())

	state.RetVal = nca4.NewValueFromI64(int64(handle))
}

// Grows the stack.
// 	stack:grow handle howmuch
// Returns the new size or an error message. May set the Error flag.
func CommandStack_Grow(state *nca4.State, params []*nca4.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to stack:grow.")
	}

	handle := int(params[0].Int64())
	if handle >= len(stacks) {
		state.Error = true
		state.RetVal = nca4.NewValue("error:Invalid Handle.")
		return
	}

	tmp := make([]*nca4.Value, int(params[1].Int64()))
	stacks[handle].data = append(stacks[handle].data, tmp...)
	state.RetVal = nca4.NewValueFromI64(int64(len(stacks[handle].data)))
}

// Pushes a value onto a stack.
// 	stack:push handle value
// Returns unchanged or an error message. May set the Error flag.
func CommandStack_Push(state *nca4.State, params []*nca4.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to stack:push.")
	}

	handle := int(params[0].Int64())
	if handle >= len(stacks) {
		state.Error = true
		state.RetVal = nca4.NewValue("error:Invalid Handle.")
		return
	}

	if stacks[handle].ptr >= len(stacks[handle].data) {
		state.Error = true
		state.RetVal = nca4.NewValue("error:Stack Overflow.")
		return
	}

	stacks[handle].data[stacks[handle].ptr] = params[1]
	stacks[handle].ptr++
}

// Pops a value off the stack.
// 	stack:pop handle
// Returns the value or an error message. May set the Error flag.
func CommandStack_Pop(state *nca4.State, params []*nca4.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to stack:push.")
	}

	handle := int(params[0].Int64())
	if handle >= len(stacks) {
		state.Error = true
		state.RetVal = nca4.NewValue("error:Invalid Handle.")
		return
	}

	stacks[handle].ptr--
	state.RetVal = stacks[handle].data[stacks[handle].ptr]
}

// Get TOS.
// 	stack:top handle
// Returns the value or an error message. May set the Error flag.
func CommandStack_Top(state *nca4.State, params []*nca4.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to stack:push.")
	}

	handle := int(params[0].Int64())
	if handle >= len(stacks) {
		state.Error = true
		state.RetVal = nca4.NewValue("error:Invalid Handle.")
		return
	}

	state.RetVal = stacks[handle].data[stacks[handle].ptr-1]
}
