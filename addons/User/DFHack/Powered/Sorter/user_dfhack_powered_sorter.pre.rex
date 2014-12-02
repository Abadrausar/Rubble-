
(if (str:cmp (rubble:configvar "POWERED_SORTER_HARD_MODE") "YES") {
	[rubble:raws "user_dfhack_powered_sorter.lua" =
		(str:replace [rubble:raws "user_dfhack_powered_sorter.lua"] "--SORTER_MODE" "consume = 5," -1)]
	[rubble:raws "user_dfhack_powered_sorter.lua" =
		(str:replace [rubble:raws "user_dfhack_powered_sorter.lua"] "--[[SORTER_POWERED]]" "wshop:isUnpowered() or " -1)]
}{
	[rubble:raws "user_dfhack_powered_sorter.lua" =
		(str:replace [rubble:raws "user_dfhack_powered_sorter.lua"] "--SORTER_MODE" "-- Consumes no power." -1)]
	[rubble:raws "user_dfhack_powered_sorter.lua" =
		(str:replace [rubble:raws "user_dfhack_powered_sorter.lua"] "--[[SORTER_POWERED]]" "" -1)]
})

(rubble:dfhack_loadscript "user_dfhack_powered_sorter.lua")
