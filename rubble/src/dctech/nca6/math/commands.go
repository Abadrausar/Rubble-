// NCA v6 Math Commands.
package math

import "dctech/nca6"

// Adds the math commands to the state.
// The math commands are:
//	add
//	sub
//	div
//	mod
//	mul
func Setup(state *nca6.State) {
	state.NewNativeCommand("add", CommandAdd)
	state.NewNativeCommand("sub", CommandSub)
	state.NewNativeCommand("div", CommandDiv)
	state.NewNativeCommand("mod", CommandMod)
	state.NewNativeCommand("mul", CommandMul)
}

// Adds two nca6.Values.
// Opperands are converted to 64 bit integers. Invalid strings are given the nca6.Value "0"
// 	add a b
// Returns a + b
func CommandAdd(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to add.")
	}

	state.RetVal = nca6.NewValueInt64(params[0].Int64() + params[1].Int64())
	return
}

// Subtracts two nca6.Values.
// Opperands are converted to 64 bit integers. Invalid strings are given the nca6.Value "0"
// 	sub a b
// Returns a - b
func CommandSub(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to sub.")
	}

	state.RetVal = nca6.NewValueInt64(params[0].Int64() - params[1].Int64())
	return
}

// Divides two nca6.Values.
// Opperands are converted to 64 bit integers. Invalid strings are given the nca6.Value "0"
// 	div a b
// Returns a / b
func CommandDiv(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to div.")
	}

	state.RetVal = nca6.NewValueInt64(params[0].Int64() / params[1].Int64())
	return
}

// Gives the remainder of dividing two nca6.Values.
// Opperands are converted to 64 bit integers. Invalid strings are given the nca6.Value "0"
// 	mod a b
// Returns a % b
func CommandMod(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to mod.")
	}

	state.RetVal = nca6.NewValueInt64(params[0].Int64() % params[1].Int64())
	return
}

// Multiplies two nca6.Values.
// Opperands are converted to 64 bit integers. Invalid strings are given the nca6.Value "0"
// 	mul a b
// Returns a * b
func CommandMul(state *nca6.State, params []*nca6.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to mul.")
	}

	state.RetVal = nca6.NewValueInt64(params[0].Int64() * params[1].Int64())
	return
}
