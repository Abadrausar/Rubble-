
(axis:write [rubble:fs] "out:scripts/rubble_change-building.lua" [rubble:raws "libs_dfhack_upgrade_building.lua"])

(rubble:template "DFHACK_REACTION_UPGRADE_BUILDING" block source dest name class {
	var id = (str:add "UPGRADE_BUILDING_" [source] "_TO_" [dest])
	[rubble:dfhack_reactions [id] = (str:add "rubble_change-building -type \"" [dest] "\" -loc [ \\\\LOCATION ]")]
	
	(str:trimspace (str:add (rubble:reaction [id] [class])
		"\n\t[NAME:" [name] "]\n"
		"\t[BUILDING:" [source] ":NONE]"
	))
})
