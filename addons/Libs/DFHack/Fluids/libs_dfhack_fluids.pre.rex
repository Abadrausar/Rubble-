
(rubble:requireaddon "Libs/DFHack/Fluids" "Libs/DFHack/Command")

(axis:write [rubble:fs] "df:hack/scripts/rubble_fluids.lua" [rubble:raws "com_libs_dfhack_fluids.lua"])
(rubble:prepfile "com_libs_dfhack_fluids.lua")
(axis:write [rubble:fs] "df:hack/lua/rubble_fluids.lua" [rubble:raws "lib_libs_dfhack_fluids.lua"])
(rubble:prepfile "lib_libs_dfhack_fluids.lua")
