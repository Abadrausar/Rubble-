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

// Raptor Debug Command.
package debug

import "dctech/raptor"
import "fmt"
import "os"
import "io"
import "text/tabwriter"

// Adds the debug commands to the state. 
// The debug commands are:
//	debug:shell
//	debug:value
//	debug:list
//	debug:namespace
//	debug:env
//	debug:clrret
func Setup(state *raptor.State) {
	state.NewNameSpace("debug")
	state.NewNativeCommand("debug:shell", CommandDebug_Shell)
	state.NewNativeCommand("debug:value", CommandDebug_Value)
	state.NewNativeCommand("debug:list", CommandDebug_List)
	state.NewNativeCommand("debug:namespace", CommandDebug_NameSpace)
	state.NewNativeCommand("debug:env", CommandDebug_Env)
	state.NewNativeCommand("debug:clrret", CommandDebug_ClrRet)
}

// Break into the debugging shell.
// This command will provide an interactive shell until eof is reached, 
// on windows this may be simulated by pressing CTRL+Z followed by <ENTER>.
// 	debug:shell
// Returns the return value of the last command to be run.
func CommandDebug_Shell(state *raptor.State, params []*raptor.Value) {
	line := make([]byte, 0, 100)
	curchar := make([]byte, 1, 1)
	escape := false

	state.Println("+-------------------------------------------+")
	state.Println("Raptor Debugging Shell")
	state.Print(">>>")
	for {
		_, err := os.Stdin.Read(curchar)
		if err == io.EOF {
			state.Println("Exiting...")
			break
		} else if err != nil {
			state.Println("Read Error:", err, "\nExiting...")
			break
		}

		if curchar[0] == byte('\r') {
			continue
		}

		if curchar[0] == byte('\\') && !escape {
			escape = true
			continue
		}

		if curchar[0] == byte('\n') && !escape {
			state.Code.Add(string(line))
			rtn, err := state.Run()
			if err != nil {
				state.Println("Error:", err)
			}
			state.Println("Ret:", rtn)

			line = line[:0]
			state.Print(">>>")
			continue
		}

		if curchar[0] != byte('\n') && escape {
			line = append(line, byte('\\'))
		}

		if curchar[0] == byte('\n') && escape {
			state.Print(">>>")
		}

		escape = false
		line = append(line, curchar...)
	}
	state.Println("+-------------------------------------------+")
}

// Print information about a script value.
// 	debug:value value
// Returns unchanged.
func CommandDebug_Value(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to debug:value.")
	}

	state.Println("+-------------------------------------------+")
	state.Println("Raptor Script Value Inspector")
	if params[0] == nil {
		state.Println("Value == nil: Aborting")
		state.Println("+-------------------------------------------+")
		return
	}
	state.Println("Data:", params[0].Data)
	switch params[0].Type {
	case raptor.TypString:
		state.Println("Type: TypString")
	case raptor.TypInt:
		state.Println("Type: TypInt")
	case raptor.TypFloat:
		state.Println("Type: TypFloat")
	case raptor.TypBool:
		state.Println("Type: TypBool")
	case raptor.TypObject:
		state.Println("Type: TypObject")
	case raptor.TypCode:
		state.Println("Type: TypCode")
	}
	state.Println(params[0].Pos)
	obj := params[0].Indexable()
	if obj != nil {
		state.Println("Value Is Indexable")

		len := obj.Len()
		state.Println("len: ", len)
		if len > 0 {
			w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
			fmt.Fprintf(w, "%v\t%v\n", "Key", "Value")
			for _, key := range obj.Keys() {
				val := obj.Get(key)
				fmt.Fprintf(w, "%v\t%v\n", key, val.String())
			}
			w.Flush()
		}
	}
	state.Println("+-------------------------------------------+")
}

// List all global data. This is mostly for use in interactive shells.
// 	debug:list
// Returns unchanged.
func CommandDebug_List(state *raptor.State, params []*raptor.Value) {
	state.Println("+-------------------------------------------+")
	state.Println("Raptor Global Data")
	state.Println("Namespaces:")
	for i := range state.NameSpaces {
		state.Println("\t" + i)
	}
	state.Println("Commands:")
	for i := range state.Commands {
		state.Println("\t" + i)
	}
	state.Println("Variables:")
	envs := []*raptor.Environment(*state.Envs)
	for i := range envs[0].Vars {
		state.Println("\t" + i)
	}
	state.Println("Registers:")
	state.Println("\tExit:", state.Exit)
	state.Println("\tReturn:", state.Return)
	state.Println("\tBreak:", state.Break)
	state.Println("\tBreakLoop:", state.BreakLoop)
	state.Println("\tError:", state.Error)
	state.Println("\tNoRecover:", state.NoRecover)
	state.Println("\tRetVal:", state.RetVal)
	state.Println("\tThis:", state.This)
	state.Println("+-------------------------------------------+")
}

// Print information about a namespace.
// 	debug:namespace name
// Returns unchanged.
func CommandDebug_NameSpace(state *raptor.State, params []*raptor.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to debug:namespace.")
	}

	namespace := state.ParseNameSpaceName(params[0].String())

	state.Println("+-------------------------------------------+")
	state.Println("Raptor Namespace Inspector")
	state.Println("Commands:")
	for i := range namespace.Commands {
		state.Println("\t" + i)
	}
	state.Println("Namespaces:")
	for i := range namespace.NameSpaces {
		state.Println("\t" + i)
	}
	state.Println("Variables:")
	for i := range namespace.Vars {
		state.Println("\t" + i)
	}
	state.Println("+-------------------------------------------+")
}

// List variables in all environments.
// 	debug:env
// Returns unchanged.
func CommandDebug_Env(state *raptor.State, params []*raptor.Value) {
	state.Println("+-------------------------------------------+")
	state.Println("Raptor Environment Inspector")
	envs := []*raptor.Environment(*state.Envs)
	for x := range envs {
		state.Println("Environment #", x, ":")
		for i := range envs[x].Vars {
			state.Println("\t" + i)
		}
	}
	state.Println("+-------------------------------------------+")
}

// Set the return value to nil.
// 	debug:clrret
// Returns nil.
func CommandDebug_ClrRet(state *raptor.State, params []*raptor.Value) {
	state.RetVal = nil
}
