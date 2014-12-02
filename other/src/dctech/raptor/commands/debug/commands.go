/*
Copyright 2012-2014 by Milo Christiansen

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
}

// Break into the debugging shell.
// This command will provide an interactive shell until eof is reached, 
// on windows this may be simulated by pressing CTRL+Z followed by <ENTER>.
// 	debug:shell
// Returns the return value of the last command to be run.
func CommandDebug_Shell(script *raptor.Script, params []*raptor.Value) {
	line := make([]byte, 0, 100)
	curchar := make([]byte, 1, 1)
	escape := false

	script.Println("+-------------------------------------------+")
	script.Println("Raptor Debugging Shell")
	script.Print(">>>")
	for {
		_, err := os.Stdin.Read(curchar)
		if err == io.EOF {
			script.Println("Exiting...")
			break
		} else if err != nil {
			script.Println("Read Error:", err, "\nExiting...")
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
			script.Code.AddString(string(line), raptor.NewPosition(1, 1, ""))
			err := script.SafeExec()
			if err != nil {
				script.Println("Error:", err)
				script.Println("Internal state may be messed up!")
			}
			script.Println("Ret:", script.RetVal)

			line = line[:0]
			script.Print(">>>")
			continue
		}

		if curchar[0] != byte('\n') && escape {
			line = append(line, byte('\\'))
		}

		if curchar[0] == byte('\n') && escape {
			script.Print(">>>")
		}

		escape = false
		line = append(line, curchar...)
	}
	script.Println("+-------------------------------------------+")
}

// Print information about a script value.
// 	debug:value value
// Returns unchanged.
func CommandDebug_Value(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	script.Println("+-------------------------------------------+")
	script.Println("Raptor Script Value Inspector")
	if params[0] == nil {
		script.Println("Value == nil: Aborting")
		script.Println("+-------------------------------------------+")
		return
	}
	script.Println("Data:", params[0].Data)
	switch params[0].Type {
	case raptor.TypString:
		script.Println("Type: TypString")
	case raptor.TypInt:
		script.Println("Type: TypInt")
	case raptor.TypFloat:
		script.Println("Type: TypFloat")
	case raptor.TypBool:
		script.Println("Type: TypBool")
	case raptor.TypObject:
		script.Println("Type: TypObject")
	case raptor.TypCode:
		script.Println("Type: TypCode")
	}
	script.Println(params[0].Pos)
	obj := params[0].Indexable()
	if obj != nil {
		script.Println("Value Is Indexable")

		len := obj.Len()
		script.Println("len: ", len)
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
	script.Println("+-------------------------------------------+")
}

// List all global data. This is mostly for use in interactive shells.
// 	debug:list
// Returns unchanged.
func CommandDebug_List(script *raptor.Script, params []*raptor.Value) {
	script.Println("+-------------------------------------------+")
	script.Println("Raptor Global Data")
	
	script.Println("Namespaces:")
	for i := range script.Host.NameSpaces.BeginLowLevel() {
		script.Println("\t" + i)
	}
	script.Host.NameSpaces.EndLowLevel()
	
	script.Println("Commands:")
	for i := range script.Host.Commands.BeginLowLevel() {
		script.Println("\t" + i)
	}
	script.Host.Commands.EndLowLevel()
	
	script.Println("Variables:")
	envs := []*raptor.Environment(*script.Envs)
	for i := range envs[0].Vars {
		script.Println("\t" + i)
	}
	
	script.Println("Registers:")
	script.Println("\tExit:", script.Exit)
	script.Println("\tReturn:", script.Return)
	script.Println("\tBreak:", script.Break)
	script.Println("\tBreakLoop:", script.BreakLoop)
	script.Println("\tError:", script.Error)
	script.Println("\tNoRecover:", script.Host.NoRecover)
	script.Println("\tRetVal:", script.RetVal)
	script.Println("\tThis:", script.This)
	script.Println("+-------------------------------------------+")
}

// Print information about a namespace.
// 	debug:namespace name
// Returns unchanged.
func CommandDebug_NameSpace(script *raptor.Script, params []*raptor.Value) {
	if len(params) != 1 {
		panic(script.BadParamCount("1"))
	}

	namespace := script.ParseNameSpaceName(params[0].String())

	script.Println("+-------------------------------------------+")
	script.Println("Raptor Namespace Inspector")
	
	script.Println("Namespaces:")
	for i := range namespace.NameSpaces.BeginLowLevel() {
		script.Println("\t" + i)
	}
	namespace.NameSpaces.EndLowLevel()
	
	script.Println("Commands:")
	for i := range namespace.Commands.BeginLowLevel() {
		script.Println("\t" + i)
	}
	namespace.Commands.EndLowLevel()
	
	script.Println("Variables:")
	for i := range namespace.Vars.BeginLowLevel() {
		script.Println("\t" + i)
	}
	namespace.Vars.EndLowLevel()
	
	script.Println("+-------------------------------------------+")
}

// List variables in all environments.
// 	debug:env
// Returns unchanged.
func CommandDebug_Env(script *raptor.Script, params []*raptor.Value) {
	script.Println("+-------------------------------------------+")
	script.Println("Raptor Environment Inspector")
	envs := []*raptor.Environment(*script.Envs)
	for x := range envs {
		script.Println("Environment #", x, ":")
		for i := range envs[x].Vars {
			script.Println("\t" + i)
		}
	}
	script.Println("+-------------------------------------------+")
}
