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

// Rex Debug Command.
package debug

import "dctech/rex"

import "dctech/iconsole"
import "io"
import "fmt"
import "os"

//import "io"
import "text/tabwriter"

// Adds the debug commands to the state.
// The debug commands are:
//	debug:shell
//	debug:value
//	debug:list
func Setup(state *rex.State) (err error) {
	defer func() {
		if !state.NoRecover {
			if x := recover(); x != nil {
				if y, ok := x.(rex.ScriptError); ok {
					err = y
					return
				}
				panic(x)
			}
		}
	}()
	
	mod := state.RegisterModule("debug")
	mod.RegisterCommand("shell", Command_Shell)
	mod.RegisterCommand("value", Command_Value)
	mod.RegisterCommand("registers", Command_Registers)
	
	return nil
}

// Break into the debugging shell.
// This command will provide an interactive shell until eof is reached,
// on windows this may be simulated by pressing CTRL+Z followed by <ENTER>.
// 	debug:shell
// Returns the return value of the last command to be run.
func Command_Shell(script *rex.Script, params []*rex.Value) {
	script.Println("+++++++++++++++++++++++++++++++++++++++++++++")
	script.Println("Debugging Shell")
	
	var code *rex.Code = nil
	var err error = nil
	var line []byte = nil
	
	console := iconsole.New()
	for {
		line, err = console.Run()
		if err == io.EOF {
			script.Println("Exiting...")
			break
		} else if err != nil {
			script.Println("Read Error:", err, "\nExiting...")
			break
		}

		_, code, err = script.Host.CompileShell(string(line), code)
		if err != nil {
			script.Println("Error:", err)
			break
		}
		
		script.Locals.Add(code)
		err = script.SafeExec(code)
		script.Locals.Remove()
		if err != nil {
			script.Println("Error:", err)
			break
		}
		script.Println("Ret:", script.RetVal)
	}
	
	script.Println("---------------------------------------------")
}

// Print information about a script value.
// 	debug:value value
// Returns unchanged.
func Command_Value(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("debug:value", "1")
	}

	script.Println("+-------------------------------------------+")
	script.Println("|")
	script.Println("| Script Value Inspector")
	if params[0] == nil {
		script.Println("| Value == nil: Aborting")
		script.Println("| (This should never happen!)")
		script.Println("|")
		script.Println("+-------------------------------------------+")
		return
	}
	script.Println("| Data:", params[0].Data)
	script.Println("| Type:", params[0].TypeString())
	script.Println("|")
	script.Println("| Position: ", params[0].Pos)
	if params[0].Type == rex.TypIndex {
		obj := params[0].Data.(rex.Indexable)
		script.Println("|")
		script.Println("| Value Is Indexable")

		len := obj.Len()
		script.Println("|   len: ", len)
		if len > 0 {
			w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
			fmt.Fprintf(w, "|   %v\t%v\n", "Key", "Value")
			for _, key := range obj.Keys() {
				val := obj.Get(key)
				fmt.Fprintf(w, "|   %v\t%v\n", key, val.String())
			}
			w.Flush()
		}
	}
	script.Println("|")
	script.Println("+-------------------------------------------+")
}

// Print the value of all internal registers and flags.
// 	debug:registers
// Returns unchanged.
func Command_Registers(script *rex.Script, params []*rex.Value) {
	script.Println("+-------------------------------------------+")
	script.Println("|")
	script.Println("| Registers")
	script.Println("|   RetVal:   ", script.RetVal)
	script.Println("|")
	script.Println("| Flags")
	script.Println("|   Error:    ", script.Error)
	script.Println("|   NoRecover:", script.Host.NoRecover)
	script.Println("|")
	script.Println("| Exit States")
	script.Println("|   Exit:     ", script.Exit)
	script.Println("|   Return:   ", script.Return)
	script.Println("|   Break:    ", script.Break)
	script.Println("|   BreakLoop:", script.BreakLoop)
	script.Println("|")
	script.Println("+-------------------------------------------+")
}

// TODO: Make command to list all modules
// TODO: Make command to list all commands
// TODO: Make command to print modules by name
// TODO: Make command to print local variables
