var rubble:cheat_reaction_count = 0

(rubble:template "DEV_CHEAT_ITEM" block name {
	[rubble:cheat_reaction_count = (int:add [rubble:cheat_reaction_count] 1)]
	(rubble:stageparse (str:add 
		"{REACTION;DEV_CHEAT_" [rubble:cheat_reaction_count] ";ADDON_HOOK_PLAYABLE}\n"
		"\t[NAME:unpack " [name] "]\n\t[BUILDING:WAGON_DEV_CHEAT:NONE]\n\t[SKILL:ANIMALCARE]"
	))
})

(ret "")