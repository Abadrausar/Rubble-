
# Announcement Template and Script Install

# Example:
# As a reaction product line:
# {DFHACK_ANNOUNCEMENT;You have unleashed a dragon!;COLOR_LIGHTRED}
# (the color defaults to white)
#

(rubble:requireaddon "Libs/DFHack/Announcement" "Libs/DFHack/Command")
module rubble:libs_dfhack_announcement

# Gets a material to show the specified announcement.
command rubble:libs_dfhack_announcement:getmat text color {
	(rubble:libs_dfhack_command:getmat <array "rubble_announce" [text] [color]>)
}

# Use this template as a replacement for a reaction product line.
(rubble:template "DFHACK_ANNOUNCEMENT" block text color="COLOR_WHITE" {
	(str:add "[PRODUCT:100:1:BOULDER:NONE:" (rubble:libs_dfhack_announcement:getmat [text] [color]) "]")
})

# Install the script.
(axis:write [rubble:fs] "df:hack/scripts/rubble_announce.lua" [rubble:raws "libs_dfhack_announcement.lua"])
(rubble:prepfile "libs_dfhack_announcement.lua")

(ret "")
