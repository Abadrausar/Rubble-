// NCA v6 Debug Command.
package debug

import "fmt"
import "dctech/nca6"
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
func Setup(state *nca6.State) {
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
func CommandDebug_Shell(state *nca6.State, params []*nca6.Value) {
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
func CommandDebug_Value(state *nca6.State, params []*nca6.Value) {
	if len(params) != 1 {
		panic("Wrong number of params to debug:value.")
	}
	
	fmt.Println("=============================================")
	fmt.Println("NCA Script Value Inspector")
	fmt.Println("Data:", params[0].Data)
	switch params[0].Type {
	case nca6.TypString:
		fmt.Println("Type: TypString")
	case nca6.TypInt:
		fmt.Println("Type: TypInt")
	case nca6.TypFloat:
		fmt.Println("Type: TypFloat")
	case nca6.TypObject:
		fmt.Println("Type: TypObject")
	}
	fmt.Println("Line: ", params[0].Line)
	fmt.Println("Column: ", params[0].Column)
	if params[0].IsIndexable() {
		data := params[0].Data.(nca6.Indexable)
		fmt.Println("Value Is Indexable")
		fmt.Println("len: ", data.Len())
		
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintf(w, "%v\t%v\n", "Key", "Value")
		for _, key := range data.Keys() {
			val := data.Get(key)
			fmt.Fprintf(w, "%v\t%v\n", key, val.String())
		}
		w.Flush()
	}
	fmt.Println("=============================================")
}

// List all global data. This is mostly for use in interactive shells.
// 	debug:list
// Returns unchanged.
func CommandDebug_List(state *nca6.State, params []*nca6.Value) {
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
	envs := []*nca6.Environment(*state.Envs)
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
func CommandDebug_NameSpace(state *nca6.State, params []*nca6.Value) {
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
func CommandDebug_Env(state *nca6.State, params []*nca6.Value) {
	fmt.Println("=============================================")
	fmt.Println("NCA Environment Inspector")
	envs := []*nca6.Environment(*state.Envs)
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
func CommandDebug_ClrRet(state *nca6.State, params []*nca6.Value) {
	state.RetVal = nil
}
