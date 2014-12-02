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

package main

import "fmt"
import "os"
import "io"
import "dctech/nca7"
import "dctech/nca7/commands/base"
import "dctech/nca7/commands/bit"
import "dctech/nca7/commands/console"
import "dctech/nca7/commands/csv"
import "dctech/nca7/commands/debug"
import "dctech/nca7/commands/env"
import "dctech/nca7/commands/file"
import "dctech/nca7/commands/fileio"
import "dctech/nca7/commands/ini"
import "dctech/nca7/commands/integer"
import "dctech/nca7/commands/regex"
import "dctech/nca7/commands/stack"
import "dctech/nca7/commands/str"
import "io/ioutil"
import "flag"
import "runtime/pprof"

// cmd.exe can't display line drawing chars with the default font.
//var header = ` ╔═══════════════════════════════════════════════════════════════════════════╗
// ║                                  NCASH 7                                  ║
// ║                          No Clever Acronym SHell                          ║
// ║                           Enter Ctrl+Z to exit.                           ║
// ╚═══════════════════════════════════════════════════════════════════════════╝
//`

var header = ` +---------------------------------------------------------------------------+
 |                                  NCASH 7                                  |
 |                          No Clever Acronym SHell                          |
 |                           Enter Ctrl+Z to exit.                           |
 +---------------------------------------------------------------------------+
`

var filename = flag.String("script", "", "The file to run in batch mode. If this is omited NCASH will run in interactive mode.")
var preload = flag.String("load", "", "Path to a file containing preload code. This file is evaluated before anything else. Does not cause a switch to batch mode")
var loadPreDefs = flag.Bool("predef", true, "Should ncash load a few extra predefined user commands? These commands are for, while, include, ++, and --.")
var recover = flag.Bool("recover", true, "Should nca recover errors? If not NCASH will CRASH on script errors.")

var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to file.")

var lextest = flag.Bool("lextest", false, "Run a lexer test.")

var preDefs = `
	'increment variable'
	(command ++ __name__ {
		(set [__name__] (int:add [[__name__]] 1))
	})
	
	'decrement variable'
	(command -- name {
		(set [name] (int:sub [[name]] 1))
	})
	
	'for loop'
	(command for __init__ __test__ __incriment__ __code__ {
		(eval [__init__])
		(loop {
			(if (evalinnew [__test__]) {
				(evalinnew [__code__])
				(evalinnew [__incriment__])
				(break -1)
			}{
				(break 0)
			})
		})
	})
	
	'simple while loop'
	(command while __test__ __code__ {
		(loop {
			(if (evalinnew [__test__]) {
				(evalinnew [__code__])
				(break -1)
			}{
				(break 0)
			})
		})
	})
	
	'include a file'
	'requires the file io commands to be loaded'
	(command include __path__ {
		(evalinparent (fileio:read [__path__]))
	})
`

func main() {
	
	flag.Parse()
	
	// Profiling - I have no idea if this works, I can't get pprof working.
	if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            panic(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
	
	// Lexer test
	if *lextest {
        lex := nca7.NewLexer(preDefs, 37, 0)
		
		for {
			lex.Advance()
			fmt.Println(lex.CurrentTkn(), lex.CurrentTkn().Lexeme)
			if lex.CheckLookAhead(nca7.TknINVALID) {
				return
			}
		}
    }
	
	fmt.Print(header)

	line := make([]byte, 0, 100)
	curchar := make([]byte, 1, 1)

	state := nca7.NewState()
	state.NoRecover = !*recover

	// Load commands
	base.Setup(state)
	bit.Setup(state)
	console.Setup(state)
	csv.Setup(state)
	debug.Setup(state)
	env.Setup(state)
	file.Setup(state)
	fileio.Setup(state)
	ini.Setup(state)
	integer.Setup(state)
	regex.Setup(state)
	stack.Setup(state)
	str.Setup(state)
	
	// Add any command line params to the state
	state.AddParams(flag.Args()...)
	
	// Run preload file if any
	if *preload != "" {
		fmt.Println("Executing Preload File:", *preload)
		file, err := ioutil.ReadFile(*preload)
		if err != nil {
			fmt.Println("Preload Error:", err)
			return
		}

		state.Code.Add(string(file))
		rtn, err := state.Run()
		if err != nil {
			fmt.Println("Preload Error:", err)
		}
		fmt.Println("Preload Ret:", rtn)
	}
	
	// Load predefs if desired
	if *loadPreDefs {
		fmt.Println("Loading Predefined User Commands...")
		state.Code.Add(preDefs)
		rtn, err := state.Run()
		if err != nil {
			fmt.Println("Predefine Error:", err)
			fmt.Println("Predefine Ret:", rtn)
		}
	}
	
	// If a batch file is provided run it
	if *filename != "" {
		fmt.Println("Executing File:", *filename)
		file, err := ioutil.ReadFile(*filename)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		state.Code.Add(string(file))
		rtn, err := state.Run()
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println("Ret:", rtn)
		//fmt.Println("Exiting...")
		//return
	}

	escape := false

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
}
