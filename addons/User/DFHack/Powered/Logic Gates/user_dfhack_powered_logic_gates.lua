
local powered = rubble.require "powered"
local pswitchable = rubble.require "powered_switchable"

function isInput(x, y, z, cb, wshop)
	local pos = powered.Center(wshop)
	if cb.machine.machine_id ~= -1 then
		if getmetatable(cb) == "building_axle_horizontalst" then
			if (cb.is_vertical == false and (x == pos.x)) or (cb.is_vertical == true and (y == pos.y))
				return true
			end
		elseif getmetatable(cb) == "building_axle_verticalst" then
			if cb.z ~= pos.z and x == pos.x and y == pos.y then
				return true
			end
		end
	end
	return false
end

-- only works with single tile workshops
function findInputs(wshop)
	local offsets = {
		{x = 1, y = 0, z = 0},
		{x = -1, y = 0, z = 0},
		{x = 0, y = 1, z = 0},
		{x = 0, y = -1, z = 0},
		{x = 0, y = 0, z = 1}
	}
	local pos = powered.Center(wshop)
	local inputs = {}
	
	for _, off in pairs(offsets) do
		local cb = dfhack.buildings.findAtTile(pos.x + off.x, pos.y + off.y, pos.z + off.z)
		if isInput(pos.x + off.x, pos.y + off.y, pos.z + off.z, cb, wshop) then
			table.insert(inputs, cb)
		end
	end
	return inputs
end

function makeLogicGate(gatetyp)
	return function(wshop)
		if not wshop:isUnpowered() then
			local inputs = findInputs(wshop)
			local outputs = pswitchable.Switchables(wshop)
			
			if #inputs == 0 or #outputs == 0 then
				return
			end
			
			local pI = false
			local uI = false
			for _, input in ipairs(inputs) do
				if input.machine_id ~= -1 then
					-- Is this assumption true? (id -> index)
					--machine = df.global.world.machines.all[input.machine_id]
					machine = input:getMachineInfo()
					if machine.flags.active then
						pI = true
					else
						uI = true
					end
					if pI == true and uI == true then
						break
					end
				end
			end
			
			if pI == false and uI == false then
				-- No valid inputs
				-- This is probably impossible
				return
			end
			
			local command = nil
			if gatetyp == "AND" then
				if pI == true and uI == false then
					command = true
				else
					command = false
				end
			elseif gatetyp == "OR" then
				if pI == true then
					command = true
				else
					command = false
				end
			elseif gatetyp == "NOT" then
				if pI == true and uI == false then
					command = false 
				elseif pI == false and uI == true then
					command = true
				end
			elseif gatetyp == "NAND" then
				if uI == true then
					command = true
				else
					command = false
				end
			elseif gatetyp == "NOR" then
				if uI == true and pI == false then
					command = true
				else
					command = false
				end
			elseif gatetyp == "XOR" then
				if pI == true and uI == true then
					command = true
				else
					command = false
				end
			elseif gatetyp == "XNOR" then
				if (pI == true and uI == false) or (pI == false and uI == true) then
					command = false
				else
					command = true
				end
			end
			
			-- Can happen when you have a multiple input NOT with inconsistent values.
			if command == nil then
				return
			end
			
			for _, output in ipairs(outputs) do
				pswitchable.SwitchBuilding(output, command)
			end
		end	
	end		
end

local outputs = {
	"AND",
	"OR"
	"NOT"
	--"NAND"
	--"NOR"
	--"XOR"
	--"XNOR"
}

for _, output in pairs(outputs) do
	buildings.registerBuilding{
		name="LOGIC_GATE_"..output,
		gears={{x=0,y=0}},
		action={10, makeLogicGate(output)},
		animate={
			isMechanical=true,
			frames={
				{{x=0,y=0,42,0,7,0}},
				{{x=0,y=0,15,0,7,0}}
			}
		}
	}
end
