
# Announcement Template and Script Install

# Example:
# In a reaction:
# {DFHACK_ANNOUNCE;You have unleashed a dragon!;COLOR_LIGHTRED}
# (the color defaults to white)
#

# Use this template as a replacement for the REACTION template
(rubble:template "DFHACK_REACTION_ANNOUNCE" block id class text color="COLOR_WHITE" {
	[rubble:dfhack_reactions [id] = (str:add 'rubble_announce "' [text] '" ' [color])]
	(rubble:reaction [id] [class])
})

(rubble:template "DFHACK_ANNOUNCE" block text color="COLOR_WHITE" id="" {
	(if (str:cmp [id] "") {
		(if (isnil [rubble:reaction_data [rubble:reaction_last]]) {
			(rubble:abort "Error: rubble:reaction_last is invalid in call to DFHACK_ANNOUNCE.")
		})
		
		[rubble:dfhack_reactions [rubble:reaction_last] = (str:add 'rubble_announce "' [text] '" ' [color])]
	}{
		[rubble:dfhack_reactions [id] = (str:add 'rubble_announce "' [text] '" ' [color])]
	})
	(ret "")
})

# Install the script.
(axis:write [rubble:fs] "out:scripts/rubble_announce.lua" [rubble:raws "libs_dfhack_announcement.lua"])

(ret "")
