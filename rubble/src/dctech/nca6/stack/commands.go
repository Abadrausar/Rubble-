// NCA v6 Stack Commands.
package stack

import "dctech/nca6"

type stack struct {
	ptr  int
	data []*nca6.Value
}

func newstack() *stack {
	rtn := new(stack)
	rtn.data = make([]*nca6.Value, 25, 25)
	return rtn
}

// Adds the stack commands to the state.
// The stack commands are:
//	stack:new
//	stack:grow
//	stack:push
//	stack:pop
//	stack:top
func Setup(state *nca6.State) {
	state.NewNameSpace("stack")
	state.NewNativeCommand("stack:new", CommandStack_New)
	state.NewNativeCommand("stack:grow", CommandStack_Grow)
	state.NewNativeCommand("stack:push", CommandStack_Push)
	state.NewNativeCommand("stack:pop", CommandStack_Pop)
	state.NewNativeCommand("stack:top", CommandStack_Top)
}

// Create a new stack of size 25.
// 	stack:new
// Returns the new stack
func CommandStack_New(state *nca6.State, params []*nca6.Value) {
	state.RetVal = nca6.NewValueObject(newstack())
}

// Grows the stack.
// 	stack:grow stack howmuch
// Returns the new size or an error message. May set the Error flag.
func CommandStack_Grow(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to stack:grow.")
	}

	if params[0].Type != nca6.TypObject {
		panic("Non-Object Param 0 passed to stack:grow.")
	}
	if _, ok := params[0].Data.(*stack); !ok {
		panic("stack:grow's Param 0 is not a *stack.")
	}
	stack := params[0].Data.(*stack)

	tmp := make([]*nca6.Value, int(params[1].Int64()))
	stack.data = append(stack.data, tmp...)
	state.RetVal = nca6.NewValueInt64(int64(len(stack.data)))
}

// Pushes a value onto a stack.
// 	stack:push stack value
// Returns unchanged or an error message. May set the Error flag.
func CommandStack_Push(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to stack:push.")
	}

	if params[0].Type != nca6.TypObject {
		panic("Non-Object Param 0 passed to stack:push.")
	}
	if _, ok := params[0].Data.(*stack); !ok {
		panic("stack:push's Param 0 is not a *stack.")
	}
	stack := params[0].Data.(*stack)

	if stack.ptr >= len(stack.data) {
		state.Error = true
		state.RetVal = nca6.NewValueString("Stack Overflow.")
		return
	}

	stack.data[stack.ptr] = params[1]
	stack.ptr++
}

// Pops a value off the stack.
// 	stack:pop stack
// Returns the value or an error message. May set the Error flag.
func CommandStack_Pop(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to stack:pop.")
	}

	if params[0].Type != nca6.TypObject {
		panic("Non-Object Param 0 passed to stack:pop.")
	}
	if _, ok := params[0].Data.(*stack); !ok {
		panic("stack:pop's Param 0 is not a *stack.")
	}
	stack := params[0].Data.(*stack)

	stack.ptr--
	state.RetVal = stack.data[stack.ptr]
}

// Get TOS.
// 	stack:top stack
// Returns the value or an error message. May set the Error flag.
func CommandStack_Top(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to stack:top.")
	}

	if params[0].Type != nca6.TypObject {
		panic("Non-Object Param 0 passed to stack:top.")
	}
	if _, ok := params[0].Data.(*stack); !ok {
		panic("stack:top's Param 0 is not a *stack.")
	}
	stack := params[0].Data.(*stack)

	state.RetVal = stack.data[stack.ptr-1]
}
