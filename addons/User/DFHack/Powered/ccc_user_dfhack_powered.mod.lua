
_ENV = rubble.mkmodule("powered_persist")

local eventful = require "plugins.eventful"

-- The Lua half of this addon takes several parts:
--	Part A: Dealing directly with workshops, finding input and output locations, etc.
--	Part B: Finding, creating, and outputting items.
--	Part C: Handling persistent output type data for workshops with many possible types.
--	Part D: Finding and changing the state of "switchable" buildings.
-- This is part C.

-- Persistent Output API
-- Use for case where the list of possibilities is WAY to large for the normal way.
-- Use the normal (multiple workshops) way for most cases, as it is easier to provide
-- nice readouts and simpler setting change mechanics.
-- For a workshop to use this API it's ID MUST start with "DFHACK_RUBBLE_POWERED_".

-- Set this workshop's output string.
-- The output string defaults to "NONE" for new buildings.
function SetOutputType(wshop, output)
	local i = "X:"..wshop.centerx.."|Y:"..wshop.centery.."|Z:"..wshop.z
	machines[i] = output
	savepersist()
end

-- Get this workshop's output string.
-- The output string defaults to "NONE" for new buildings.
function GetOutputType(wshop)
	local i = "X:"..wshop.centerx.."|Y:"..wshop.centery.."|Z:"..wshop.z
	if machines[i] == nil then
		machines[i] = "NONE"
	end
	return machines[i]
end

-- Get this workshop's output string and run it as code.
-- Returns nil if the output is "NONE" or there is an error when loading the code.
-- If there is an error it is logged to the DFHack console.
function GetOutputTypeAsCode(wshop)
	local outputraw = GetOutputType(wshop)
	if outputraw == "NONE" then
		return nil
	end
	
	local f, err = load(outputraw)
	if f == nil then
		dfhack.printerr(err)
		return nil
	end
	return f()
end

-- Internal stuff, subject to change without notice!

function savepersist()
	local out = "rubble.powered_persist.machines = {\n"
	for k, v in pairs(machines) do
		if v ~= nil then
			out = out..'\t["'..k..'"] = [['..v..']],\n'
		end
	end
	out = out.."}"
	
	rawmachines.value = out
	rawmachines:save()
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
	savepersist()
end

eventful.onBuildingCreatedDestroyed.rubble_powered=function(building_id)
	refreshMachineList()
end

-- ["X:"..x.."|Y:"..y.."|Z:"..z] = "<output id>"
machines = {}
rawmachines = dfhack.persistent.get("rubble_user_dfhack_powered")
if rawmachines == nil then
	rawmachines, _ = dfhack.persistent.save({key = "rubble_user_dfhack_powered"})
else
	local f, err = load(rawmachines.value)
	if f == nil then
		error(err)
	end
	f()
end
refreshMachineList()

return _ENV
