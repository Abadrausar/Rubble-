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

// NCA v7 Debug Command.
package debug

import "fmt"
import "dctech/nca7"
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
func Setup(state *nca7.State) {
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
func CommandDebug_Shell(state *nca7.State, params []*nca7.Value) {
	line := make([]byte, 0, 100)
	curchar := make([]byte, 1, 1)
	escape := false

	fmt.Println("=============================================")
	fmt.Println("NCA Debugging Shell")
	fmt.Print(">>>")
	for {
		_, err := os.Stdin.Read(curchar)
		if err == io.EOF {
			fmt.Println("Exiting...")
			break
		} else if err != nil {
			fmt.Println("Read Error:", err, "\nExiting...")
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
				fmt.Println("Error:", err)
			}
			fmt.Println("Ret:", rtn)
			
			line = line[:0]
			fmt.Print(">>>")
			continue
		}

		if curchar[0] != byte('\n') && escape {
			line = append(line, byte('\\'))
		}
		
		if curchar[0] == byte('\n') && escape {
			fmt.Print(">>>")
		}

		escape = false
		line = append(line, curchar...)
	}
	fmt.Println("=============================================")
}

// Print information about a script value.
// 	debug:value value
// Returns unchanged.
func CommandDebug_Value(state *nca7.State, params []*nca7.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to debug:value.")
	}
	
	fmt.Println("=============================================")
	fmt.Println("NCA Script Value Inspector")
	fmt.Println("Data:", params[0].Data)
	switch params[0].Type {
	case nca7.TypString:
		fmt.Println("Type: TypString")
	case nca7.TypInt:
		fmt.Println("Type: TypInt")
	case nca7.TypFloat:
		fmt.Println("Type: TypFloat")
	case nca7.TypBool:
		fmt.Println("Type: TypBool")
	case nca7.TypObject:
		fmt.Println("Type: TypObject")
	}
	fmt.Println("Line: ", params[0].Line)
	fmt.Println("Column: ", params[0].Column)
	obj := params[0].Indexable()
	if obj != nil {
		fmt.Println("Value Is Indexable")
		fmt.Println("len: ", obj.Len())
		
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintf(w, "%v\t%v\n", "Key", "Value")
		for _, key := range obj.Keys() {
			val := obj.Get(key)
			fmt.Fprintf(w, "%v\t%v\n", key, val.String())
		}
		w.Flush()
	}
	fmt.Println("=============================================")
}

// List all global data. This is mostly for use in interactive shells.
// 	debug:list
// Returns unchanged.
func CommandDebug_List(state *nca7.State, params []*nca7.Value) {
	fmt.Println("=============================================")
	fmt.Println("NCA Global Data")
	fmt.Println("Namespaces:")
	for i := range state.NameSpaces {
		fmt.Println("\t" + i)
	}
	fmt.Println("Commands:")
	for i := range state.Commands {
		fmt.Println("\t" + i)
	}
	fmt.Println("Variables:")
	envs := []*nca7.Environment(*state.Envs)
	for i := range envs[0].Vars {
		fmt.Println("\t" + i)
	}
	fmt.Println("Registers:")
	fmt.Println("\tReturn=", state.Return)
	fmt.Println("\tBreak=", state.Break)
	fmt.Println("\tExit=", state.Exit)
	fmt.Println("\tError=", state.Error)
	fmt.Println("\tNoRecover=", state.NoRecover)
	fmt.Println("=============================================")
}

// Print information about a namespace.
// 	debug:namespace name
// Returns unchanged.
func CommandDebug_NameSpace(state *nca7.State, params []*nca7.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to debug:namespace.")
	}
	
	namespace := state.ParseNameSpaceName(params[0].String())
	
	fmt.Println("=============================================")
	fmt.Println("NCA Namespace Inspector")
	fmt.Println("Commands:")
	for i := range namespace.Commands {
		fmt.Println("\t" + i)
	}
	fmt.Println("Namespaces:")
	for i := range namespace.NameSpaces {
		fmt.Println("\t" + i)
	}
	fmt.Println("Variables:")
	for i := range namespace.Vars {
		fmt.Println("\t" + i)
	}
	fmt.Println("=============================================")
}

// List variables in all environments.
// 	debug:env
// Returns unchanged.
func CommandDebug_Env(state *nca7.State, params []*nca7.Value) {
	fmt.Println("=============================================")
	fmt.Println("NCA Environment Inspector")
	envs := []*nca7.Environment(*state.Envs)
	for x := range envs {
		fmt.Println("Environment #", x, ":")
		for i := range envs[x].Vars {
			fmt.Println("\t" + i)
		}
	}
	fmt.Println("=============================================")
}

// Set the return value to nil.
// 	debug:clrret
// Returns nil.
func CommandDebug_ClrRet(state *nca7.State, params []*nca7.Value) {
	state.RetVal = nil
}
