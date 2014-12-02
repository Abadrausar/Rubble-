
(axis:write [rubble:fs] "out:dfhack/user_dfhack_pop_cap.lua" (str:add 	
`
-- "User/DFHack/Pop Cap" DFHack Lua Script

population_cap = population_cap or df.global.d_init.population_cap
strict_population_cap = strict_population_cap or df.global.d_init.strict_population_cap
baby_cap_absolute = baby_cap_absolute or df.global.d_init.baby_cap_absolute
baby_cap_percent = baby_cap_percent or df.global.d_init.baby_cap_percent

population_cap_new = "`(rubble:configvar "population_cap")`"
strict_population_cap_new = "`(rubble:configvar "strict_population_cap")`"
baby_cap_absolute_new = "`(rubble:configvar "baby_cap_absolute")`"
baby_cap_percent_new = "`(rubble:configvar "baby_cap_percent")`"

dfhack.onStateChange.rubble_user_dfhack_pop_cap = function(code)
	if code == SC_WORLD_UNLOADED then
		print("User/DFHack/Pop Cap: Restoring global pop cap data.")
		df.global.d_init.population_cap = population_cap
		df.global.d_init.strict_population_cap = strict_population_cap
		df.global.d_init.baby_cap_absolute = baby_cap_absolute
		df.global.d_init.baby_cap_percent = baby_cap_percent
	end
	if code == SC_WORLD_LOADED then
		print("User/DFHack/Pop Cap: Setting per-world pop cap data.")
		df.global.d_init.population_cap = tonumber(population_cap_new) or population_cap
		df.global.d_init.strict_population_cap = tonumber(strict_population_cap_new) or strict_population_cap
		df.global.d_init.baby_cap_absolute = tonumber(baby_cap_absolute_new) or baby_cap_absolute
		df.global.d_init.baby_cap_percent = tonumber(baby_cap_percent_new) or baby_cap_percent
	end
end
`
))
