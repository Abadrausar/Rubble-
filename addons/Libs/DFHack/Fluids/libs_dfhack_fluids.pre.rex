
(rubble:requireaddon "Libs/DFHack/Fluids" "Libs/DFHack/Command")

(axis:write [rubble:fs] "df:hack/scripts/rubble_fluids.lua" [rubble:raws "com_libs_dfhack_fluids.lua"])
(rubble:prepfile "com_libs_dfhack_fluids.lua")
(axis:write [rubble:fs] "df:hack/lua/rubble_fluids.lua" [rubble:raws "lib_libs_dfhack_fluids.lua"])
(rubble:prepfile "lib_libs_dfhack_fluids.lua")

# 
# area may be:
#	spot
#	3x3
#	5x5
#	7x7
# 
# fluid should be magma or water
# 

# {DFHACK_FLUIDS_EAT;magma;5x5;4;4}
(rubble:template "DFHACK_FLUIDS_EAT" block fluid area amount minimum {
	(str:add "[PRODUCT:100:1:BOULDER:NONE:" (rubble:libs_dfhack_command:getmat <array 
		"rubble_fluids" "eat" [area] [fluid] [amount] [minimum] "\\LOCATION"
	>) "]")
})

# {DFHACK_FLUIDS_CART;water;spot}
(rubble:template "DFHACK_FLUIDS_CART" block fluid area {
	(str:add "[PRODUCT:100:1:BOULDER:NONE:" (rubble:libs_dfhack_command:getmat <array 
		"rubble_fluids" "cart" [area] [fluid] "\\LOCATION"
	>) "]")
})

# {DFHACK_FLUIDS_SPAWN;water;spot;1}
(rubble:template "DFHACK_FLUIDS_SPAWN" block fluid area amount {
	(str:add "[PRODUCT:100:1:BOULDER:NONE:" (rubble:libs_dfhack_command:getmat <array 
		"rubble_fluids" "spawn" [area] [fluid] [amount] "\\LOCATION"
	>) "]")
})
