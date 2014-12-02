
# Adds the templates, commands and variables that make up the automatic DFHack init script support

var rubble:dfhack_commands = <map>
var rubble:dfhack_reactions = <map>

# Use these two to add extra custom stuff
var rubble:dfhack_extras_lua = <array>
var rubble:dfhack_extras_init = <array>

command rubble:dfhack_loadscript scr {
	(if (isnil [rubble:raws [scr]]) {
		(rubble:abort (str:add "    Error: Attempt to install missing DFHack script: " [scr]))
	}{
		(axis:write [rubble:fs] (str:add "out:dfhack/" [scr]) [rubble:raws [scr]])
	})
}

command rubble:dfhack_runcommand com {
	[rubble:dfhack_commands [com] = true]
}

command rubble:dfhack_write {
	var base = "\n-- DFHack init.lua script file\n-- Automatically generated, DO NOT EDIT!\n"
	[base = (str:add [base] "print(\"Loading DFHack scripts from Rubble addons:\")\n")]
	
	[base = (str:add [base] 
		`\n-- Important Globals\n`
		`dfhack.BASE_G.rubble = dfhack.BASE_G.rubble or {}\n`
		`rubble.savedir = SAVE_PATH\n`
		`function rubble.load_script(name)\n`
		`	env = {}\n`
		`	setmetatable(env, { __index = dfhack.BASE_G })\n`
		`	\n`
		`	local f, perr = loadfile(name, 't', env)\n`
		`	if f then\n`
		`		return safecall(f)\n`
		`	end\n`
		`	dfhack.printerr("    Error: "..perr)\n`
		`end\n`
		`function rubble.mkmodule(name)\n`
		`	rubble[name] = rubble[name] or {}\n`
		`	setmetatable(rubble[name], { __index = dfhack.BASE_G })\n`
		`	return rubble[name]\n`
		`end\n`
	)]
	
	[base = (str:add [base] 
		`\n-- Pseudo Modules and Scripts\n`
		`local scrlist = dfhack.internal.getDir(SAVE_PATH.."/raw/dfhack/")\n`
		`if scrlist then\n`
		`	table.sort(scrlist)\n`
		`	for i,name in ipairs(scrlist) do\n`
		`		if string.match(name,'%.mod.lua$') then\n`
		`			print("  Module: "..name)\n`
		`			rubble.load_script(SAVE_PATH.."/raw/dfhack/"..name)\n`
		`		end\n`
		`	end\n`
		`	for i,name in ipairs(scrlist) do\n`
		`		if string.match(name,'%.lua$') and not string.match(name,'%.mod.lua$') then\n`
		`			print("  Script: "..name)\n`
		`			rubble.load_script(SAVE_PATH.."/raw/dfhack/"..name)\n`
		`		end\n`
		`	end\n`
		`end\n`
	)]
	
	[base = (str:add [base] "\n-- Extras:\n")]
	(foreach [rubble:dfhack_extras_lua] block _ txt {
		[base = (str:add [base] [txt])]
		(break true)
	})
	
	[base = (str:add [base] `\n-- Just in case you forgot...\n`)]
	[base = (str:add [base] `dfhack.gui.showAnnouncement("This region's raws were generated with Rubble v` [rubble:version] `!", COLOR_LIGHTGREEN)\n`)]
	
	(axis:write [rubble:fs] "out:init.lua" [base])

	[base = "\n# DFHack onLoad.init file\n# Automatically generated, DO NOT EDIT!\n"]
	
	[base = (str:add [base] "\n# Reactions:\n")]
	(foreach [rubble:dfhack_reactions] block id action {
		[base = (str:add [base] "modtools/reaction-trigger -reactionName \"" [id] "\" -command [ " [action] " ]\n")]
		(break true)
	})
	
	[base = (str:add [base] "\n# Commands:\n")]
	(foreach [rubble:dfhack_commands] block action _ {
		[base = (str:add [base] [action] "\n")]
		(break true)
	})
	
	[base = (str:add [base] "\n# Extras:\n")]
	(foreach [rubble:dfhack_extras_init] block _ txt {
		[base = (str:add [base] [txt])]
		(break true)
	})
	
	(axis:write [rubble:fs] "out:onLoad.init" [base])
}

# Load a script when the world is loaded, the script must be in the raw/dfhack folder
(rubble:template "DFHACK_LOADSCRIPT" block scr {
	(rubble:dfhack_loadscript [scr])
})

# Run a command when the world is loaded.
(rubble:template "DFHACK_RUNCOMMAND" block com {
	(rubble:dfhack_runcommand [com])
})

# The Reaction template is with the tech templates.
