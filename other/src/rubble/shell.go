/*
Copyright 2013 by Milo Christiansen

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

import "os"
import "io"
import "dctech/raptor"
import "io/ioutil"

func ShellModeRun() {
	LogPrintln("=============================================")
	LogPrintln("Entering Shell Mode.")

	// Lexer test
	if LexTest {
		lex := raptor.NewLexer(preDefs, 64, 1)
		if ScriptPath != "" {
			file, err := ioutil.ReadFile(ScriptPath)
			if err != nil {
				LogPrintln("Error:", err)
				return
			}
			lex = raptor.NewLexer(string(file), 1, 1)
		}

		for {
			lex.Advance()
			LogPrintln(lex.CurrentTkn(), lex.CurrentTkn().Lexeme)
			if lex.CheckLookAhead(raptor.TknINVALID) {
				return
			}
		}
	}
	
	script := raptor.NewScript()

	if Validate && !ValidateAll {
		LogPrintln("Validating Script File...")
		if ScriptPath == "" {
			LogPrintln("Validate Error: No script set.")
			return
		}
		
		file, err := ioutil.ReadFile(ScriptPath)
		if err != nil {
			LogPrintln("Validate Error:", err)
			return
		}
		
		err = raptor.LoadFile(file, script)
		if err != nil {
			LogPrintln("Validate Error:", err)
			return
		}
		
		err = script.Validate()
		if err != nil {
			LogPrintln("Validate Error:", err)
			return
		}
		LogPrintln("Validation Successful.")
		return
	}

	// Add any command line params to the GlobalRaptorState.
	script.AddParams(Flags.Args()...)

	// Load predefs if desired.
	if !NoPredefs {
		LogPrintln("Loading Predefined User Commands...")
		script.Code.AddCodeSource(raptor.NewLexer(preDefs, 64, 1))
		rtn, err := GlobalRaptorState.RunShell(script)
		if err != nil {
			LogPrintln("Predefine Error:", err)
			LogPrintln("Predefine Ret:", rtn)
		}
		script.RetVal = nil
	}

	if Validate && ValidateAll {
		LogPrintln("Validating Script File...")
		if ScriptPath == "" {
			LogPrintln("Validate Error: No script set.")
			return
		}
		
		file, err := ioutil.ReadFile(ScriptPath)
		if err != nil {
			LogPrintln("Validate Error:", err)
			return
		}
		
		err = raptor.LoadFile(file, script)
		if err != nil {
			LogPrintln("Validate Error:", err)
			return
		}
		
		err = GlobalRaptorState.Validate(script)
		if err != nil {
			LogPrintln("Validate Error:", err)
			return
		}
		LogPrintln("Validation Successful.")
		return
	}

	// Compile the provided script if Compile is set
	if Compile != "" {
		if ScriptPath == "" {
			LogPrintln("Compile Error: No script set.")
			return
		}
		LogPrintln("Compiling File:", ScriptPath)
		file, err := ioutil.ReadFile(ScriptPath)
		if err != nil {
			LogPrintln("Compile Error:", err)
			return
		}

		in := raptor.Compile(string(file), raptor.NewPositionInfo(0, -1))

		// Small binaries can fail to compile because of the string size restrictions,
		// so we fall through and try to compile a normal one in that case.

		var out []byte
		switch BinVersion {
		case 2:
			LogPrintln("Compiling small binary...")
			out, err = raptor.WriteBinaryV2(in)
			if err == nil {
				break
			}
			LogPrintln("Compile Error:", err)
			LogPrintln("Falling back to normal binary.")
			fallthrough
		case 1:
			LogPrintln("Compiling normal binary...")
			out, err = raptor.WriteBinaryV1(in)
			if err != nil {
				LogPrintln("Compile Error:", err)
				LogPrintln("Compile Failed.")
				return
			}
		case 4:
			LogPrintln("Compiling small binary (Compressed string table)...")
			out, err = raptor.WriteBinaryV4(in)
			if err == nil {
				break
			}
			LogPrintln("Compile Error:", err)
			LogPrintln("Falling back to normal binary.")
			fallthrough
		case 3:
			LogPrintln("Compiling normal binary (Compressed string table)...")
			out, err = raptor.WriteBinaryV3(in)
			if err != nil {
				LogPrintln("Compile Error:", err)
				LogPrintln("Compile Failed.")
				return
			}
		default:
			LogPrintln("Invalid binary version. Aborting.")
			return
		}

		err = ioutil.WriteFile(Compile, out, 0600)
		if err != nil {
			LogPrintln("Compile Error:", err)
			return
		}

		LogPrintln("Exiting...")
		return
	}

	// If a batch file is provided run it.
	if ScriptPath != "" && Compile == "" {
		LogPrintln("Executing File:", ScriptPath)
		file, err := ioutil.ReadFile(ScriptPath)
		if err != nil {
			LogPrintln("Error:", err)
			return
		}

		err = raptor.LoadFile(file, script)
		if err != nil {
			LogPrintln("Error:", err)
			return
		}

		rtn, err := GlobalRaptorState.RunShell(script)
		if err != nil {
			LogPrintln("Error:", err)
		}
		LogPrintln("Ret:", rtn)

		if !NoExit {
			LogPrintln("Exiting...")
			return
		}
	}

	// Interactive Shell
	escape := false
	line := make([]byte, 0, 100)
	curchar := make([]byte, 1, 1)

	LogPrint(">>>")
	for {
		_, err := os.Stdin.Read(curchar)
		if err == io.EOF {
			LogPrintln("Exiting...")
			break
		} else if err != nil {
			LogPrintln("Read Error:", err, "\nExiting...")
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
			script.Code.Add(string(line))
			rtn, err := GlobalRaptorState.RunShell(script)
			if err != nil {
				LogPrintln("Error:", err)
			}
			LogPrintln("Ret:", rtn)

			line = line[:0]
			LogPrint(">>>")
			continue
		}

		if curchar[0] != byte('\n') && escape {
			line = append(line, byte('\\'))
		}

		if curchar[0] == byte('\n') && escape {
			LogPrint(">>>")
		}

		escape = false
		line = append(line, curchar...)
	}
}

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
