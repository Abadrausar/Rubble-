
module rubble:libs_dfhack_command
var rubble:libs_dfhack_command:mats = (sort:new)

# Gets a material to run a specified DFHack command via autosyndrome
# pass in an array of the command name followed by any params, including special autosyndrome commands if you want.
command rubble:libs_dfhack_command:getmat params {
	(if (int:lt (len [params]) 1) {
		(rubble:abort "Error: Wrong param count to rubble:libs_dfhack_command:getmat")
	})
	
	var id = ""
	(for 0 (int:sub (len [params]) 1) 1 block i {
		[id = (str:add [id] "_" [params [i]])]
	})
	
	(if (exists [rubble:libs_dfhack_command:mats] [id]){
	}{
		[rubble:libs_dfhack_command:mats [id] = (str:add
			"[INORGANIC:DFHACK_COMMAND" [id] "]\n"
			"\t[USE_MATERIAL_TEMPLATE:COMMAND_STONE_TEMPLATE]\n"
			"\t[SPECIAL]\n"
			"\t[SYNDROME]\n"
			"\t\t[SYN_CLASS:\\AUTOSYNDROME]\n"
			"\t\t[SYN_CLASS:\\COMMAND]\n"
		)]
		
		(for 0 (int:sub (len [params]) 1) 1 block i {
			[rubble:libs_dfhack_command:mats [id] = (str:add
				[rubble:libs_dfhack_command:mats [id]]
				"\t\t[SYN_CLASS:" [params [i]] "]\n"
			)]
		})
	})
	(str:add "INORGANIC:DFHACK_COMMAND" [id])
}

# Use this template as a replacement for a reaction product line.
(rubble:template "DFHACK_COMMAND" block ... {
	(str:add "[PRODUCT:100:1:BOULDER:NONE:" (rubble:libs_dfhack_command:getmat [params]) "]")
})

(rubble:template "#_DFHACK_COMMAND_MATS" {
	var out = ""
	(foreach [rubble:libs_dfhack_command:mats] block _ mat {
		[out = (str:add [out] "\n" [mat])]
	})
	(str:trimspace [out])
})

(ret "")
