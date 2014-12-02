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

/*
Shell for the Rex scripting language.

Run "rexsh -h" for usage information.

*/
package main

import "fmt"
import "os"
import "io"
import "io/ioutil"
import "dctech/iconsole"
import "dctech/rex"
import "dctech/rex/commands/base"
import "dctech/rex/commands/boolean"
import "dctech/rex/commands/console"
import "dctech/rex/commands/convert"
import "dctech/rex/commands/debug"
import "dctech/rex/commands/env"
import "dctech/rex/commands/expr"
import "dctech/rex/commands/file"
import "dctech/rex/commands/fileio"
import "dctech/rex/commands/float"
import "dctech/rex/commands/integer"
import "dctech/rex/commands/png"
import "dctech/rex/commands/random"
import "dctech/rex/commands/regex"
import "dctech/rex/commands/sort"
import "dctech/rex/commands/str"
import "dctech/rex/commands/structure"
import "dctech/rex/commands/thread"
import "flag"
import "runtime"
import "runtime/pprof"

// The bit about "Enter Ctrl+Z or Ctrl+Break to Exit" is true in Windows when using the default command prompt.
// Ctrl+Z generates char 26 (0x1A), which is the "DOS EOF char".
// Ctrl+Break works just as good or better, you don't even have to press enter.
// If none of these options work for you just type Alt+26 which is guaranteed to work on all platforms.
var header = ` +----------------------------------------------------------------------------+
 |                                  Rex Shell                                 |
 |                      Enter Ctrl+Z or Ctrl+Break to Exit                    |
 +----------------------------------------------------------------------------+
`

var ScriptPath string
var NoExit bool
var NoRecover bool
var Threaded bool
var Profile string

var preDefPos = rex.NewPosition(76, 1, "main.go")
var preDefs = `
	# nothing for now...
`

func main() {

	flag.StringVar(&ScriptPath, "script", "", "Path to the input script, if any. Changes to batch mode.")
	flag.BoolVar(&NoExit, "noexit", false, "If set changes from batch mode to interactive mode. Use with -script.")
	flag.BoolVar(&NoRecover, "norecover", false, "If set disables error recovery. Use for debugging the runtime.")
	flag.BoolVar(&Threaded, "threads", false, "Allows the shell to use more than one processor core, not useful unless running a threaded script.")
	flag.StringVar(&Profile, "profile", "", "Output CPU profile information to specified file.")

	flag.Parse()

	if Threaded {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	if Profile != "" {
		f, err := os.Create(Profile)
		if err != nil {
			fmt.Println("Profile Error:", err)
			return
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	fmt.Print(header)

	state := rex.NewState()
	script := rex.NewScript()
	state.NoRecover = NoRecover

	rex.SetupArrays(state)
	rex.SetupMaps(state)
	base.Setup(state)
	boolean.Setup(state)
	console.Setup(state)
	convert.Setup(state)
	debug.Setup(state)
	env.Setup(state)
	expr.Setup(state)
	file.Setup(state)
	fileio.Setup(state)
	float.Setup(state)
	integer.Setup(state)
	png.Setup(state)
	random.Setup(state)
	regex.Setup(state)
	sort.Setup(state)
	str.Setup(state)
	structure.Setup(state)
	thread.Setup(state)

	var code *rex.Code = nil
	var val *rex.Value = nil
	var err error = nil

	// Load predefs
	val, code, err = state.CompileShell(preDefs, code)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	ret, err := state.RunShell(script, val)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Ret:", ret)

	// If a batch file is provided run it.
	if ScriptPath != "" {
		fmt.Println("Executing File:", ScriptPath)
		file, err := ioutil.ReadFile(ScriptPath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		val, code, err = state.CompileShell(string(file), code)
		var ret *rex.Value
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			ret, err = state.RunShell(script, val)
			if err != nil {
				fmt.Println("Error:", err)
			}
			fmt.Println("Ret:", ret)
		}

		if !NoExit {
			fmt.Println("Exiting...")
			
			// Non-number values convert to 0, so this is perfect.
			os.Exit(int(ret.Int64()))
		}
	}

	// Interactive Shell
	console := iconsole.New()
	for {
		line, err := console.Run()
		if err == io.EOF {
			fmt.Println("Exiting...")
			break
		} else if err != nil {
			fmt.Println("Read Error:", err, "\nExiting...")
			break
		}

		val, code, err = state.CompileShell(string(line), code)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			ret, err := state.RunShell(script, val)
			if err != nil {
				fmt.Println("Error:", err)
			}
			fmt.Println("Ret:", ret)
		}
	}
}
