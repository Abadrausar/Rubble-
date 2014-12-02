
# Adds the templates, commands and variables that make up the automatic DFHack init script support

var rubble:dfhack_scripts = <map>
var rubble:dfhack_commands = <map>
var rubble:dfhack_reactions = <map>

command rubble:dfhack_loadscript scr {
	(if (isnil [rubble:raws [scr]]) {
		(rubble:abort (str:add "    Error: Attempt to install missing DFHack script: " [scr]))
	}{
		(if (isnil [rubble:dfhack_scripts [scr]]) {
			(axis:write [rubble:fs] (str:add "out:dfhack/" [scr]) [rubble:raws [scr]])
			[rubble:dfhack_scripts [scr] = true]
		})
	})
}

command rubble:dfhack_runcommand com {
	var comstring = ""
	(foreach [com] block _ param {
		[comstring = (str:add [comstring] [param] ' ')]
		(break true)
	})
	[rubble:dfhack_commands [comstring] = [com]]
}

command rubble:dfhack_write {
	var base = "\n-- DFHack init.lua script file\n-- Automatically generated, DO NOT EDIT!\n"
	[base = (str:add [base] "print(\"Initializing Rubble DFHack run time addon support:\")\n")]
	
	[base = (str:add [base] "\n-- Scripts:\nlocal scrdir = SAVE_PATH..\"/raw/dfhack/\"\n")]
	(foreach [rubble:dfhack_scripts] block scr _ {
		[base = (str:add [base] "\nprint(\"  Script: " [scr] "\")\n")]
		[base = (str:add [base] "dofile(scrdir..\"" [scr] "\")\n")]
		(break true)
	})
	
	[base = (str:add [base] "\n-- Commands:\n")]
	(foreach [rubble:dfhack_commands] block comstring com {
		[base = (str:add [base] "\nprint(\"  Command: " [comstring] "\")\n")]
		[base = (str:add [base] "dfhack.run_command(")]
		(foreach [com] block _ param {
			[base = (str:add [base] '"' [param] '",')]
			(break true)
		})
		(str:trimright [base] 1)
		[base = (str:add [base] "\")\n")]
		(break true)
	})
	
	[base = (str:add [base] "\n-- Reactions:\n")]
	[base = (str:add [base] "-- Not Yet Implemented\n")]
	(foreach [rubble:dfhack_reactions] block id action {
		[base = (str:add [base] "-- " [id] " -> " [action] "\n")]
		(break true)
	})
	
	(axis:write [rubble:fs] "out:init.lua" [base])
}

# Load a script when the world is loaded, the script must be in the raw/dfhack folder
(rubble:template "DFHACK_LOADSCRIPT" block scr {
	(rubble:dfhack_loadscript [scr])
})

# Run a command when the world is loaded.
(rubble:template "DFHACK_RUNCOMMAND" block ... {
	(rubble:dfhack_runcommand [params])
})

# The Reaction template is with the tech templates.
