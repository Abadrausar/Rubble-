
# Spawn Template and Script Install

module rubble:libs_dfhack_spawn

# Install the spawn script.
(axis:write [rubble:fs] "df:hack/scripts/rubble_spawn.lua" [rubble:raws "libs_dfhack_spawn.lua"])
(rubble:prepfile "libs_dfhack_spawn.lua")


# Use this template as a replacement for the REACTION template
(rubble:template "DFHACK_REACTION_SPAWN" block id class creature caste name {
	[rubble:dfhack_reactions [id] = (str:add 'rubble_spawn "' [creature] '" "' [caste] '" "' [name] '" \\LOCATION')]
	(rubble:reaction [id] [class])
})

(rubble:template "DFHACK_SPAWN" block creature caste name id="" {
	(if (str:cmp [id] "") {
		(if (isnil [rubble:reaction_data [rubble:reaction_last]]) {
			(rubble:abort "Error: rubble:reaction_last is invalid in call to DFHACK_ANNOUNCE.")
		})
		
		[rubble:dfhack_reactions [rubble:reaction_last] = (str:add 'rubble_spawn "' [creature] '" "' [caste] '" "' [name] '" \\LOCATION')]
	}{
		[rubble:dfhack_reactions [id] = (str:add 'rubble_spawn "' [creature] '" "' [caste] '" "' [name] '" \\LOCATION')]
	})
	(ret "")
})

(ret "")
