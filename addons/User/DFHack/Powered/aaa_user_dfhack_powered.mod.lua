
_ENV = rubble.mkmodule("powered")

local buildings = require 'plugins.building-hacks'

-- The Lua half of this addon takes three parts:
--	Part A: Dealing directly with workshops, finding input and output locations, etc.
--	Part B: Finding, creating, and outputting items.
--	Part C: Handling persistent output type data for workshops with many possible types.
-- This is part A.

-- This is a convenience function for registering powered workshops that use the Rex API
-- provided by this addon.
-- If you do not use the Rex API this function will not be useful.
-- 
-- id should match that used when registering the workshop on the Rubble side of things
-- outputs is the same basic thing as the Rubble equivalent, just that here it is just a
-- table of the IDs, no need for the names.
-- consume is how much power to use
-- ticks is how often to run the function returned by makeaction
-- makeaction is a function that takes one parameter (the output id) and returns
-- a function to handle the action, this function should take a single parameter as well, the workshop.
function Register(id, outputs, consume, ticks, makeaction)
	print("    Registering mechanical workshop: "..id)
	
	local register = function(output)
		if output == nil then
			buildings.registerBuilding{
				name=id,
				consume=consume,
				gears={{x=0,y=2},{x=2,y=0},{x=2,y=2},{x=0,y=0}},
				action={ticks, makeaction("")},
				animate={
					isMechanical=true,
					frames={
						{{x=0,y=2,42,0,0,1},{x=2,y=0,42,0,0,1},{x=2,y=2,42,0,0,1},{x=0,y=0,42,0,0,1}},
						{{x=0,y=2,15,0,0,1},{x=2,y=0,15,0,0,1},{x=2,y=2,15,0,0,1},{x=0,y=0,15,0,0,1}}
					}
				}
			}
		else
			buildings.registerBuilding{
				name=id.."_"..output,
				consume=consume,
				gears={{x=0,y=2},{x=2,y=0},{x=2,y=2},{x=0,y=0}},
				action={ticks, makeaction(output)},
				animate={
					isMechanical=true,
					frames={
						{{x=0,y=2,42,0,0,1},{x=2,y=0,42,0,0,1},{x=2,y=2,42,0,0,1},{x=0,y=0,42,0,0,1}},
						{{x=0,y=2,15,0,0,1},{x=2,y=0,15,0,0,1},{x=2,y=2,15,0,0,1},{x=0,y=0,15,0,0,1}}
					}
				}
			}
		end
	end
	
	if outputs == nil or #outputs == 0 then
		register(nil)
	else
		for _, output in pairs(outputs) do
			register(output)
		end
	end
	
end

-- Returns true if there is an input at the position.
function InputAt(x, y, z)
	building = dfhack.buildings.findAtTile(x, y, z)
	if building ~= nil then
		if getmetatable(building) == "building_workshopst" then
			t = df.building_def.find(building.custom_type)
			if t ~= nil and t ~= -1 then
				if t.code == "DFHACK_INPUT_POWERED" then
					return true
				end
			end
		end
	end
	return false
end

-- Returns true if the workshop has at least one input tile.
function HasInput(wshop)
	local pos = Area(wshop)
	
	for cx = pos.x1, pos.x2, 1 do
		for cy = pos.y1, pos.y2, 1 do
			if InputAt(cx, cy, pos.z) then
				return true
			end
		end
	end
	return false
end

-- Returns a list of input positions.
-- List uses 1-based indexing!
function Inputs(wshop)
	local pos = Area(wshop)
	local rtn = {}
	
	for cx = pos.x1, pos.x2, 1 do
		for cy = pos.y1, pos.y2, 1 do
			if InputAt(cx, cy, pos.z) then
				table.insert(rtn, {x = cx, y = cy, z = pos.z})
			end
		end
	end
	return rtn
end

-- Returns true if there is an output at the position.
function OutputAt(x, y, z)
	building = dfhack.buildings.findAtTile(x, y, z)
	if building ~= nil then
		if getmetatable(building) == "building_workshopst" then
			t = df.building_def.find(building.custom_type)
			if t ~= nil and t ~= -1 then
				if t.code == "DFHACK_OUTPUT_POWERED" then
					return true
				end
			end
		end
	end
	return false
end

-- Returns true if the workshop has at least one output tile.
function HasOutput(wshop)
	local pos = Area(wshop)
	
	for cx = pos.x1, pos.x2, 1 do
		for cy = pos.y1, pos.y2, 1 do
			if OutputAt(cx, cy, pos.z) then
				return true
			end
		end
	end
	return false
end

-- Returns a list of output positions.
-- The locations returned are actually the tiles on the other side
-- of the outputs from the workshop!
-- List uses 1-based indexing!
function Outputs(wshop)
	local pos = Area(wshop)
	local center = Center(wshop)
	local rtn = {}
	
	for cx = pos.x1, pos.x2, 1 do
		for cy = pos.y1, pos.y2, 1 do
			if OutputAt(cx, cy, pos.z) then
				local opos = {x = cx, y = cy, z = pos.z}
				if cx > pos.x2-1 then
					opos.x = opos.x + 1
				elseif cx < pos.x1+1 then
					opos.x = opos.x - 1
				end
				if cy > pos.y2-1 then
					opos.y = opos.y + 1
				elseif cy < pos.y1+1 then
					opos.y = opos.y - 1
				end
				table.insert(rtn, opos)
			end
		end
	end
	return rtn
end

-- Sets an item as forbidden if it is on an input.
-- Keeps dwarves from stealing stuff from the middle of your production lines...
function ForbidIfNeeded(item)
	if InputAt(item.pos.x, item.pos.y, item.pos.z) then
		item.flags.forbid = true
	else
		item.flags.forbid = false
	end
end

-- Returns a {x,y,z} for the workshop's center tile.
function Center(wshop)
	return {x = wshop.centerx, y = wshop.centery, z = wshop.z}
end

-- Returns a {x1,y1,x2,y2,z} that describes the area around the workshop.
function Area(wshop)
	return {x1 = wshop.x1-1, y1 = wshop.y1-1, x2 = wshop.x2+1, y2 = wshop.y2+1, z = wshop.z}
end

-- Returns a fake "workshop" for use with the functions in this API.
-- Use for case where the workshop is not immediately available.
-- Only works for odd sized square workshops (eg 1x1, 3x3, 5x5, etc.).
function MakeFake(x, y, z, size)
	local o = math.floor(size / 2) + 1
	return {centerx = x, centery = y, x1 = x-o, y1 = y-o, x2 = x+o, y2 = y+o, z = z}
end

-- Same as MakeFake, but for any size or shape.
function MakeFakeAdv(cx, cy, x1, y1, x2, y2, z)
	return {centerx = cx, centery = cy, x1 = x1, y1 = y1, x2 = x2, y2 = y2, z = z}
end

return _ENV