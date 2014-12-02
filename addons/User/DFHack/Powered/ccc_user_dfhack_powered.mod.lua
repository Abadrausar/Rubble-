
_ENV = rubble.mkmodule("powered_persist")

local eventful = require "plugins.eventful"

-- The Lua half of this addon takes three parts:
--	Part A: Dealing directly with workshops, finding input and output locations, etc.
--	Part B: Finding, creating, and outputting items.
--	Part C: Handling persistent output type data for workshops with many possible types.
-- This is part C.

-- Persistent Output API
-- Use for case where the list of possibilities is WAY to large for the normal way.
-- Use the normal (multiple workshops) way for most cases, as it is more reliable.
-- For a workshop to use this API it's ID MUST start with "DFHACK_RUBBLE_POWERED_".

-- Set this workshops output string.
-- The output string defaults to "NONE" for new buildings.
function SetOutputType(wshop, output)
	local i = "X:"..wshop.centerx.."|Y:"..wshop.centery.."|Z:"..wshop.z
	machines[i] = output
	
	savePersist()
end

-- Get this workshops output string.
-- The output string defaults to "NONE" for new buildings.
function GetOutputType(wshop)
	local i = "X:"..wshop.centerx.."|Y:"..wshop.centery.."|Z:"..wshop.z
	if machines[i] == nil then
		machines[i] = "NONE"
	end
	return machines[i]
end

-- Internal stuff, subject to change without notice!

-- ["X:"..x.."|Y:"..y.."|Z:"..z] = "<output id>"
machines = {}

eventful.onBuildingCreatedDestroyed.rubble_powered=function(building_id)
	refreshMachineList()
end

function savePersist()
	print("rubble.powered_persist: Saving persistence data.")
	local dat = io.open(rubble.savedir..'/raw/rubble_user_dfhack_powered.persist.lua', 'w')
	dat:write("\n")
	dat:write("rubble.powered_persist.machines = {\n")
	for k, v in pairs(machines) do
		if v ~= nil then
			dat:write('\t["'..k..'"] = "'..v..'",\n')
		end
	end
	dat:write("}\n")
	dat:close()
end

function refreshMachineList()
	local newmachines = {}
	for i=0, #df.global.world.buildings.all-1, 1 do
		building = df.global.world.buildings.all[i]
		if getmetatable(building) == "building_workshopst" then
			t = df.building_def.find(building.custom_type)
			if t ~= nil and t ~= -1 then
				if string.match(t.code, "^DFHACK_RUBBLE_POWERED_") then
					local i = "X:"..building.centerx.."|Y:"..building.centery.."|Z:"..building.z
					if machines[i] ~= nil then
						newmachines[i] = machines[i]
					else
						newmachines[i] = "NONE"
					end
				end
			end
		end
	end
	machines = newmachines
	
	savePersist()
end

print("rubble.powered: Loading persistence data.")
rubble.load_script(rubble.savedir.."/raw/rubble_user_dfhack_powered.persist.lua")

refreshMachineList()

return _ENV
