
# Adds the templates, commands and variables that make up the automatic DFHack init script support

var rubble:dfhack_commands = <map>
var rubble:dfhack_reactions = <map>

# Use these two to add extra custom stuff
var rubble:dfhack_extras_lua = <array>
var rubble:dfhack_extras_init = <array>

command rubble:dfhack_loadscript scr {
	(if (isnil [rubble:raws [scr]]) {
		(rubble:abort (str:add "Attempt to install missing DFHack script: " [scr]))
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
`
-- Important Globals and Loader Functions
dfhack.BASE_G.rubble = dfhack.BASE_G.rubble or {}
rubble.savedir = SAVE_PATH
function rubble.load_script(name)
	env = {}
	setmetatable(env, { __index = dfhack.BASE_G })
	
	local f, perr = loadfile(name, 't', env)
	if f then
		return safecall(f)
	end
	dfhack.printerr("    Error: "..perr)
end
function rubble.load_module(name)
	env = {}
	setmetatable(env, { __index = dfhack.BASE_G })
	
	local f, perr = loadfile(name, 't', env)
	if f then
		local a, module = safecall(f)
		if not module then
			dfhack.printerr("    Error: Nil module returned.")
			return nil
		end
		rubble[module.module] = module
		return module
	end
	dfhack.printerr("    Error: "..perr)
end
function rubble.mkmodule(name)
	rubble[name] = rubble[name] or {}
	setmetatable(rubble[name], { __index = dfhack.BASE_G })
	rubble[name].module = name
	return rubble[name]
end
function rubble.require(name)
	return rubble[name]
end
`
	)]
	
	[base = (str:add [base] 
`
-- The Pseudo Module and Script Loader
function rubble.reload_scripts()
	local scrlist = dfhack.internal.getDir(rubble.savedir.."/raw/dfhack/")
	if scrlist then
		table.sort(scrlist)
		for i,name in ipairs(scrlist) do
			if string.match(name,'%.mod.lua$') then
				print("  Module: "..name)
				rubble.load_module(rubble.savedir.."/raw/dfhack/"..name)
			end
		end
		for i,name in ipairs(scrlist) do
			if string.match(name,'%.lua$') and not string.match(name,'%.mod.lua$') then
				print("  Script: "..name)
				rubble.load_script(rubble.savedir.."/raw/dfhack/"..name)
			end
		end
	else
		print("  No scripts installed.")
	end
end
rubble.reload_scripts()
`
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
		[base = (str:add [base] [txt] "\n")]
		(break true)
	})
	
	(axis:write [rubble:fs] "out:onLoad.init" [base])
}

# Install a startup script or pseudo module by name.
(rubble:template "DFHACK_LOADSCRIPT" block scr {
	(rubble:dfhack_loadscript [scr])
})

# Run a command when the world is loaded.
(rubble:template "DFHACK_RUNCOMMAND" block com {
	(rubble:dfhack_runcommand [com])
})

# The DFHACK_REACTION template is with the tech templates.
