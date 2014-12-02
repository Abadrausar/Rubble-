/*
For copyright/license see header in file "doc.go"
*/

package main

import "fmt"
import "os"
import "io"
import "dctech/raptor"
import "dctech/raptor/commands/base"
import "dctech/raptor/commands/bit"
import "dctech/raptor/commands/boolean"
import "dctech/raptor/commands/console"
import "dctech/raptor/commands/convert"
import "dctech/raptor/commands/csv"
import "dctech/raptor/commands/debug"
import "dctech/raptor/commands/env"
import "dctech/raptor/commands/file"
import "dctech/raptor/commands/fileio"
import "dctech/raptor/commands/float"
import "dctech/raptor/commands/ini"
import "dctech/raptor/commands/integer"
import "dctech/raptor/commands/raw"
import "dctech/raptor/commands/regex"
import "dctech/raptor/commands/sort"
import "dctech/raptor/commands/str"
import "dctech/raptor/commands/thread"
import "io/ioutil"
import "flag"
import "runtime"

// The bit about "Enter Ctrl+Z to exit." is true in windows when using the default command prompt, 
// I have no idea how to generate an EOF in linux/mac.
var header = ` +----------------------------------------------------------------------------+
 |                                Raptor Shell                                |
 |                             Enter Ctrl+Z to Exit                           |
 +----------------------------------------------------------------------------+
`

var ScriptPath string
var Compile string
var BinVersion int
var Validate bool
var ValidateAll bool
var NoExit bool
var NoPredefs bool
var NoRecover bool
var LexTest bool
var Threaded bool

var preDefPos = raptor.NewPosition(53, 1, "main.go")
var preDefs = `
	# increment variable
	(command ++ __name__ {
		(set [__name__] (int:add [[__name__]] 1))
	})
	
	# decrement variable
	(command -- __name__ {
		(set [__name__] (int:sub [[__name__]] 1))
	})
	
	# for loop
	(command for __init__ __test__ __incriment__ __code__ {
		(eval [__init__])
		(loop {
			(if (evalinnew [__test__]) {}{(breakloop false)})
			(evalinnew [__code__])
			(evalinnew [__incriment__])
			(break true)
		})
	})
	
	# while loop
	(command while __test__ __code__ {
		(loop {
			(if (evalinnew [__test__]) {}{(breakloop false)})
			(evalinnew [__code__])
			(break true)
		})
	})
	
	# include a file
	# requires the file io commands to be loaded
	(command include __path__ {
		(evalinparent (fileio:read [__path__]))
	})
`

func main() {

	flag.StringVar(&ScriptPath, "script", "", "Path to the input script, if any. Changes to batch mode. Needed for -compile")
	flag.StringVar(&Compile, "compile", "", "Path to the output file. Changes to compile mode. Needs -script to be set.")
	flag.IntVar(&BinVersion, "binversion", 4, "Force a specific binary version. Fallback rules still apply.")
	flag.BoolVar(&Validate, "validate", false, "Run script through the syntax checker and exit. Use with -script.")
	flag.BoolVar(&ValidateAll, "all", false, "Check commands and object literals, may result in false positives. Use with -validate.")
	flag.BoolVar(&NoExit, "noexit", false, "If set changes from batch mode to interactive mode. Use with -script.")
	flag.BoolVar(&NoPredefs, "nopredefs", false, "If set disables the shell predefs.")
	flag.BoolVar(&NoRecover, "norecover", false, "If set disables error recovery. Use for debuging the runtime.")
	flag.BoolVar(&LexTest, "lextest", false, "Makes the shell run a lexer test and exit.")
	flag.BoolVar(&Threaded, "threads", false, "Allows the shell to use more than one processor core, not useful unless running a threaded script.")

	flag.Parse()

	if Threaded {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	
	fmt.Print(header)

	// Lexer test
	if LexTest {
		lex := raptor.NewLexer(preDefs, preDefPos)
		if ScriptPath != "" {
			file, err := ioutil.ReadFile(ScriptPath)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			lex = raptor.NewLexer(string(file), raptor.NewPosition(1, 1, ScriptPath))
		}

		for {
			lex.Advance()
			fmt.Println(lex.CurrentTkn(), lex.CurrentTkn().Lexeme)
			if raptor.CheckLookAhead(lex, raptor.TknINVALID) {
				return
			}
		}
	}

	state := raptor.NewState()
	script := raptor.NewScript()
	state.NoRecover = NoRecover

	if Validate && !ValidateAll {
		fmt.Println("Validating Script File...")
		if ScriptPath == "" {
			fmt.Println("Validate Error: No script set.")
			return
		}
		
		file, err := ioutil.ReadFile(ScriptPath)
		if err != nil {
			fmt.Println("Validate Error:", err)
			return
		}
		
		err = raptor.LoadFile(ScriptPath, file, script)
		if err != nil {
			fmt.Println("Validate Error:", err)
			return
		}
		
		err = script.Validate()
		if err != nil {
			fmt.Println("Validate Error:", err)
			return
		}
		fmt.Println("Validation Successful.")
		return
	}

	// Load commands
	base.Setup(state)
	bit.Setup(state)
	boolean.Setup(state)
	console.Setup(state)
	convert.Setup(state)
	csv.Setup(state)
	debug.Setup(state)
	env.Setup(state)
	file.Setup(state)
	fileio.Setup(state)
	float.Setup(state)
	ini.Setup(state)
	integer.Setup(state)
	raw.Setup(state)
	regex.Setup(state)
	sort.Setup(state)
	str.Setup(state)
	thread.Setup(state)

	// Add any command line params to the state.
	script.AddParams(flag.Args()...)

	// Load predefs if desired.
	if !NoPredefs {
		fmt.Println("Loading Predefined User Commands...")
		script.Code.Add(raptor.NewLexer(preDefs, preDefPos))
		rtn, err := state.RunShell(script)
		if err != nil {
			fmt.Println("Predefine Error:", err)
			fmt.Println("Predefine Ret:", rtn)
		}
		script.RetVal = raptor.NewValue()
	}

	if Validate && ValidateAll {
		fmt.Println("Validating Script File...")
		if ScriptPath == "" {
			fmt.Println("Validate Error: No script set.")
			return
		}
		
		file, err := ioutil.ReadFile(ScriptPath)
		if err != nil {
			fmt.Println("Validate Error:", err)
			return
		}
		
		err = raptor.LoadFile(ScriptPath, file, script)
		if err != nil {
			fmt.Println("Validate Error:", err)
			return
		}
		
		err = state.Validate(script)
		if err != nil {
			fmt.Println("Validate Error:", err)
			return
		}
		fmt.Println("Validation Successful.")
		return
	}

	// Compile the provided script if Compile is set
	if Compile != "" {
		if ScriptPath == "" {
			fmt.Println("Compile Error: No script set.")
			return
		}
		fmt.Println("Compiling File:", ScriptPath)
		file, err := ioutil.ReadFile(ScriptPath)
		if err != nil {
			fmt.Println("Compile Error:", err)
			return
		}

		in := raptor.NewCode(raptor.NewLexer(string(file), raptor.NewPosition(0, -1, "")))

		// Small binaries can fail to compile because of the string size restrictions,
		// so we fall through and try to compile a normal one in that case.

		var out []byte
		switch BinVersion {
		case 2:
			fmt.Println("Compiling small binary...")
			out, err = raptor.WriteBinaryV2(in)
			if err == nil {
				break
			}
			fmt.Println("Compile Error:", err)
			fmt.Println("Falling back to normal binary.")
			fallthrough
		case 1:
			fmt.Println("Compiling normal binary...")
			out, err = raptor.WriteBinaryV1(in)
			if err != nil {
				fmt.Println("Compile Error:", err)
				fmt.Println("Compile Failed.")
				return
			}
		case 4:
			fmt.Println("Compiling small binary (Compressed string table)...")
			out, err = raptor.WriteBinaryV4(in)
			if err == nil {
				break
			}
			fmt.Println("Compile Error:", err)
			fmt.Println("Falling back to normal binary.")
			fallthrough
		case 3:
			fmt.Println("Compiling normal binary (Compressed string table)...")
			out, err = raptor.WriteBinaryV3(in)
			if err != nil {
				fmt.Println("Compile Error:", err)
				fmt.Println("Compile Failed.")
				return
			}
		default:
			fmt.Println("Invalid binary version. Aborting.")
			return
		}

		err = ioutil.WriteFile(Compile, out, 0600)
		if err != nil {
			fmt.Println("Compile Error:", err)
			return
		}

		fmt.Println("Exiting...")
		return
	}

	// If a batch file is provided run it.
	if ScriptPath != "" && Compile == "" {
		fmt.Println("Executing File:", ScriptPath)
		file, err := ioutil.ReadFile(ScriptPath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		err = raptor.LoadFile(ScriptPath, file, script)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		rtn, err := state.RunShell(script)
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println("Ret:", rtn)

		if !NoExit {
			fmt.Println("Exiting...")
			return
		}
	}

	// Interactive Shell
	escape := false
	line := make([]byte, 0, 100)
	curchar := make([]byte, 1, 1)

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
			script.Code.AddString(string(line), raptor.NewPosition(1, 1, ""))
			rtn, err := state.RunShell(script)
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
}
