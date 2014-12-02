
_ENV = rubble.mkmodule("powered_switchable")

local powered = rubble.require "powered"

-- The Lua half of this addon takes several parts:
--	Part A: Dealing directly with workshops, finding input and output locations, etc.
--	Part B: Finding, creating, and outputting items.
--	Part C: Handling persistent output type data for workshops with many possible types.
--	Part D: Finding and changing the state of "switchable" buildings.
-- This is part D.

-- Is a building switchable?
function IsSwitchable(building)
	if getmetatable(building) == "building_gear_assemblyst"
	or getmetatable(building) == "building_doorst" 
	or getmetatable(building) == "building_hatchst" 
	or getmetatable(building) == "building_grate_wallst" 
	or getmetatable(building) == "building_grate_floorst" 
	or getmetatable(building) == "building_bars_verticalst" 
	or getmetatable(building) == "building_bars_floorst" 
	or getmetatable(building) == "building_floodgatest"
	or getmetatable(building) == "building_weaponst"
	or getmetatable(building) == "building_bridgest" then
		return true
	elseif getmetatable(building) == "building_trapst" and building.trap_type == df.emums.trap_type["Lever"] then
		return true
	end
	return false
end

-- Returns a list of all switchable buildings adjacent to the workshop.
function Switchables(wshop)
	local outputs = {}
	local apos = powered.Area(wshop)
	for cx = apos.x1, apos.x2, 1 do
		for cy = apos.y1, apos.y2, 1 do
			cb = dfhack.buildings.findAtTile(cx, cy, apos.z)
			if cb ~= nil then
				if IsSwitchable(cb) then
					table.insert(outputs, cb)
				end
			end
		end
	end
	return outputs
end

-- Change the state of a switchable building.
function SwitchBuilding(building, state)
	if getmetatable(building) == "building_gear_assemblyst" then
		if state == false then
			if building.gear_flags.disengaged == false then
				building:setTriggerState(1)
				building.gear_flags.disengaged = true
			end
		elseif state == true then
			if building.gear_flags.disengaged == true then
				building:setTriggerState(1)
				building.gear_flags.disengaged = false
			end
		end
	elseif getmetatable(building) == "building_doorst" 
	or getmetatable(building) == "building_hatchst" then
		building.door_flags.operated_by_mechanisms=true
		if state == false then
			building.close_timer = 0
		elseif state == true then
			building.close_timer = 1
		end
	elseif getmetatable(building) == "building_grate_wallst" 
	or getmetatable(building) == "building_grate_floorst" 
	or getmetatable(building) == "building_bars_verticalst" 
	or getmetatable(building) == "building_bars_floorst" 
	or getmetatable(building) == "building_floodgatest"
	or getmetatable(building) == "building_weaponst"
	or getmetatable(building) == "building_bridgest" then
		if state == false then
			if building.gate_flags.closed == false then
				building.gate_flags.closing = true
				building.timer = 1
			end
		elseif state == true then
			if building.gate_flags.closed == true then
				building.gate_flags.opening = true
				building.timer = 1
			end
		end
	elseif getmetatable(building) == "building_trapst" and building.trap_type == df.emums.trap_type["Lever"] then
		-- This is the same basic code used in "lever.rb", but translated to Lua
		-- This badly needs testing...
		if state == false then
			if building.state == 1 then
				for _, mech in pairs(building.linked_mechanisms) do
					for _, ref in pairs(mech.general_refs) do
						if getmetatable(ref) == "general_ref_building_holderst" then
							--ref.building_id
							-- Is this the right state?
							ref.building_tg:setTriggerState(1)
						end
					end
				end
				building.state = 0
			end
		elseif state == true then
			if building.state == 0 then
				for _, mech in pairs(building.linked_mechanisms) do
					for _, ref in pairs(mech.general_refs) do
						if getmetatable(ref) == "general_ref_building_holderst" then
							--ref.building_id
							-- Is this the right state?
							ref.building_tg:setTriggerState(0)
						end
					end
				end
				building.state = 1
			end
		end
		
	end
end

-- Switch all adjacent buildings
function SwitchBuildings(wshop, state)
	local outputs = Switchables(wshop)
	for _, output in pairs(outputs) do
		SwitchBuilding(output, state)
	end
end

return _ENV
