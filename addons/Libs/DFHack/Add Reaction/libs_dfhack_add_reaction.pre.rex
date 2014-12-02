
module rubble:libs_dfhack_add_reaction
var rubble:libs_dfhack_add_reaction:reactions = <array>
var rubble:libs_dfhack_add_reaction:buildings = <array>

command rubble:libs_dfhack_add_reaction:add reaction building {
	[rubble:libs_dfhack_add_reaction:reactions "append" = (str:add '"' [reaction] '", "' [building] '"')]
}

command rubble:libs_dfhack_add_reaction:clear building {
	[rubble:libs_dfhack_add_reaction:buildings "append" = [building]]
}

command rubble:libs_dfhack_add_reaction:write {
	var base = "\n-- Rubble \"Libs/DFHack/Add Reaction\" script file\n-- Automatically generated, DO NOT EDIT!\n\n"
	
	[base = (str:add [base] "local event = require \"plugins.eventful\"\n\n")]
	
	[base = (str:add [base] "-- Remove reactions from shops\n")]
	(foreach [rubble:libs_dfhack_add_reaction:buildings] block _ building {
		[base = (str:add [base] "event.removeNative(\"" [building] "\")\n")]
		(break true)
	})
	
	[base = (str:add [base] "\n-- Add reactions to shops\n")]
	(foreach [rubble:libs_dfhack_add_reaction:reactions] block _ reaction {
		[base = (str:add [base] "event.addReactionToShop(" [reaction] ")\n")]
		(break true)
	})
	
	[rubble:raws "libs_dfhack_add_reaction.lua" = [base]]
	(rubble:dfhack_loadscript "libs_dfhack_add_reaction.lua")
}

(rubble:template "REACTION_ADD" block reaction building {
	(rubble:libs_dfhack_add_reaction:add [reaction] [building])
})

(rubble:template "REACTION_CLEAR_LIST" block building {
	(rubble:libs_dfhack_add_reaction:clear [building])
})

(ret "")
