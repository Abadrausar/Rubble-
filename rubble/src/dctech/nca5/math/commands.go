// NCA v5 Math Commands.
package math

import "dctech/nca5"

// Adds the math commands to the state.
// The math commands are:
//	add
//	sub
//	div
//	mod
//	mul
func Setup(state *nca5.State) {
	state.NewNativeCommand("add", CommandAdd)
	state.NewNativeCommand("sub", CommandSub)
	state.NewNativeCommand("div", CommandDiv)
	state.NewNativeCommand("mod", CommandMod)
	state.NewNativeCommand("mul", CommandMul)
}

// Adds two nca5.Values.
// Opperands are converted to 64 bit integers. Invalid strings are given the nca5.Value "0"
// 	add a b
// Returns a + b
func CommandAdd(state *nca5.State, params []*nca5.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to add.")
	}

	state.RetVal = nca5.NewValueFromI64(params[0].Int64() + params[1].Int64())
	return
}

// Subtracts two nca5.Values.
// Opperands are converted to 64 bit integers. Invalid strings are given the nca5.Value "0"
// 	sub a b
// Returns a - b
func CommandSub(state *nca5.State, params []*nca5.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to sub.")
	}

	state.RetVal = nca5.NewValueFromI64(params[0].Int64() - params[1].Int64())
	return
}

// Divides two nca5.Values.
// Opperands are converted to 64 bit integers. Invalid strings are given the nca5.Value "0"
// 	div a b
// Returns a / b
func CommandDiv(state *nca5.State, params []*nca5.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to div.")
	}

	state.RetVal = nca5.NewValueFromI64(params[0].Int64() / params[1].Int64())
	return
}

// Gives the remainder of dividing two nca5.Values.
// Opperands are converted to 64 bit integers. Invalid strings are given the nca5.Value "0"
// 	mod a b
// Returns a % b
func CommandMod(state *nca5.State, params []*nca5.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to mod.")
	}

	state.RetVal = nca5.NewValueFromI64(params[0].Int64() % params[1].Int64())
	return
}

// Multiplies two nca5.Values.
// Opperands are converted to 64 bit integers. Invalid strings are given the nca5.Value "0"
// 	mul a b
// Returns a * b
func CommandMul(state *nca5.State, params []*nca5.Value) {
	if len(params) != 2 {
		panic("Wrong number of params to mul.")
	}

	state.RetVal = nca5.NewValueFromI64(params[0].Int64() * params[1].Int64())
	return
}
