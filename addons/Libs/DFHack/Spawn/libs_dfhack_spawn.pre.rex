
# Spawn Template and Script Install

# Example:
# As a reaction product line:
# {DFHACK_SPAWN_CREATURE;CAT;0;Fuzzy}
#

(rubble:requireaddon "Libs/DFHack/Spawn" "Libs/DFHack/Command")
module rubble:libs_dfhack_spawn

# Gets a material to spawn the specified creature
command rubble:libs_dfhack_spawn:getmat creature caste_number name {
	(rubble:libs_dfhack_command:getmat <array "rubble_spawn" [creature] [caste_number] [name] "\\LOCATION">)
}

# Use this template as a replacement for a reaction product line.
(rubble:template "DFHACK_SPAWN_CREATURE" block creature caste_number name {
	(str:add "[PRODUCT:100:1:BOULDER:NONE:" (rubble:libs_dfhack_spawn:getmat [creature] [caste_number] [name]) "]")
})

# Install the spawn script.
(axis:write [rubble:fs] "df:hack/scripts/rubble_spawn.lua" [rubble:raws "libs_dfhack_spawn.lua"])
(rubble:prepfile "libs_dfhack_spawn.lua")

(ret "")
