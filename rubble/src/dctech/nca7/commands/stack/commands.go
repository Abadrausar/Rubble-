/*
Copyright 2012-2013 by Milo Christiansen

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

// NCA v7 Stack Commands.
package stack

import "dctech/nca7"

type stack struct {
	ptr  int
	data []*nca7.Value
}

func newstack() *stack {
	rtn := new(stack)
	rtn.data = make([]*nca7.Value, 25, 25)
	return rtn
}

// Adds the stack commands to the state.
// The stack commands are:
//	stack:new
//	stack:grow
//	stack:push
//	stack:pop
//	stack:top
func Setup(state *nca7.State) {
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
func CommandStack_New(state *nca7.State, params []*nca7.Value) {
	state.RetVal = nca7.NewValueObject(newstack())
}

// Grows the stack.
// 	stack:grow stack howmuch
// Returns the new size or an error message. May set the Error flag.
func CommandStack_Grow(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to stack:grow.")
	}

	if params[0].Type != nca7.TypObject {
		panic("Non-Object Param 0 passed to stack:grow.")
	}
	if _, ok := params[0].Data.(*stack); !ok {
		panic("stack:grow's Param 0 is not a *stack.")
	}
	stack := params[0].Data.(*stack)

	tmp := make([]*nca7.Value, int(params[1].Int64()))
	stack.data = append(stack.data, tmp...)
	state.RetVal = nca7.NewValueInt64(int64(len(stack.data)))
}

// Pushes a value onto a stack.
// 	stack:push stack value
// Returns unchanged or an error message. May set the Error flag.
func CommandStack_Push(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to stack:push.")
	}

	if params[0].Type != nca7.TypObject {
		panic("Non-Object Param 0 passed to stack:push.")
	}
	if _, ok := params[0].Data.(*stack); !ok {
		panic("stack:push's Param 0 is not a *stack.")
	}
	stack := params[0].Data.(*stack)

	if stack.ptr >= len(stack.data) {
		state.Error = true
		state.RetVal = nca7.NewValueString("Stack Overflow.")
		return
	}

	stack.data[stack.ptr] = params[1]
	stack.ptr++
}

// Pops a value off the stack.
// 	stack:pop stack
// Returns the value or an error message. May set the Error flag.
func CommandStack_Pop(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to stack:pop.")
	}

	if params[0].Type != nca7.TypObject {
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
func CommandStack_Top(state *nca7.State, params []*nca7.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to stack:top.")
	}

	if params[0].Type != nca7.TypObject {
		panic("Non-Object Param 0 passed to stack:top.")
	}
	if _, ok := params[0].Data.(*stack); !ok {
		panic("stack:top's Param 0 is not a *stack.")
	}
	stack := params[0].Data.(*stack)

	state.RetVal = stack.data[stack.ptr-1]
}
