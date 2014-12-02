
_ENV = rubble.mkmodule("powered_items")

-- The Lua half of this addon takes three parts:
--	Part A: Dealing directly with workshops, finding input and output locations, etc.
--	Part B: Finding, creating, and outputting items.
--	Part C: Handling persistent output type data for workshops with many possible types.
-- This is part B.

-- Find either adjacent magma or a bar of fuel on an input tile.
-- Returns magma, bar.
-- If magma is true bar will be nil, if magma is false bar will be either a bar of fuel or nil.
function FindFuel(wshop)
	local fuelmat = dfhack.matinfo.find("COAL:COKE")
	local check = function(item)
		if df.item_type[item:getType()] == "BAR" then
			local mat = dfhack.matinfo.decode(item)
			if mat.type == fuelmat.type then
				return true
			end
		end
		return false
	end
	
	local pos = rubble.powered.Area(wshop)
	if rubble.fluids.checkInArea(pos.x1, pos.y1, pos.x2, pos.y2, pos.z, true, 4, 4) then
		return true, nil
	end
	
	return false, FindItemAtInput(wshop, check)
end

-- Returns an item or nil.
-- Searches adjacent to and on top of the workshop.
-- check should be a function that take an item and returns true if it
-- is like the one you are looking for.
function FindItemArea(wshop, check)
	local apos = rubble.powered.Area(wshop)
	local find = function(x, y, z)
		local itemblock = dfhack.maps.ensureTileBlock(x, y, z)
		if itemblock.occupancy[x%16][y%16].item == true then
			for c=#itemblock.items-1,0,-1 do
				local item = df.item.find(itemblock.items[c])
				if item.pos.x == x and item.pos.y == y and item.pos.z == z then
					if check(item) then
						return item
					end
				end
			end
		end
		return nil
	end
	
	for cx = apos.x1, apos.x2, 1 do
		for cy = apos.y1, apos.y2, 1 do
			item = find(cx, cy, apos.z)
			if item ~= nil then
				return item
			end
		end
	end
	return nil
end

-- Find an item on an adjacent input tile.
-- Checks all input tiles, not just the first one it finds.
-- Returns an item or nil
-- check should be a function that take an item and returns true if it
-- is like the one you are looking for.
function FindItemAtInput(wshop, check)
	local inputs = rubble.powered.Inputs(wshop)
	local find = function(x, y, z)
		local itemblock = dfhack.maps.ensureTileBlock(x, y, z)
		if itemblock.occupancy[x%16][y%16].item == true then
			for c=#itemblock.items-1,0,-1 do
				local item = df.item.find(itemblock.items[c])
				if item.pos.x == x and item.pos.y == y and item.pos.z == z then
					if check(item) then
						return item
					end
				end
			end
		end
		return nil
	end
	
	for _, pos in pairs(inputs) do
		local item = find(pos.x, pos.y, pos.z)
		if item ~= nil then
			return item
		end
	end
	return nil
end

-- Create a basic item, you will have to set dimensions, subtype or stack size if needed.
-- The removed flag is set (as needed by moveToGround), so remember to clear this if you need to!
-- If unit is not nil then the item quality is based on it's skill.
-- id should be an item id of the form "item_barst" or "item_boulderst".
-- The item is returned.
function CreateItem(mat, id, unit, skill)
	local item = df[id]:new()
	
	item.id = df.global.item_next_id
	df.global.world.items.all:insert('#',item)
	df.global.item_next_id = df.global.item_next_id+1
	
	item:setMaterial(mat.type)
	item:setMaterialIndex(mat.index)
	
	item:setMakerRace(df.global.ui.race_id)
	if unit ~= nil then
		item:assignQuality(unit, skill)
		item:setMaker(unit.id)
	end
	
	item:categorize(true)
	item.flags.removed=true
	
	return item
end

-- Like CreateItem, but with numeric type and subtype.
function CreateItemNumeric(mat, typ, subtyp, unit, skill)
	local item = CreateItem(mat, 'item_'..string.lower(df.item_type[typ])..'st', unit, skill)
	
	if subtyp ~= -1 then
		item:setSubtype(subtyp)
	end
	
	return item
end

-- Like CreateItem, but just use "BAR" or "bar" instead of "item_barst".
function CreateItemBasic(mat, id, unit, skill)
	return CreateItem(mat, 'item_'..string.lower(id)..'st', unit, skill)
end

-- Put the item on one of the workshop's output tiles.
-- If there are no output tiles the item is placed in the workshop center.
-- If the item is to be placed on an input tile it will be forbidden.
-- If there is more than one output tile one will be chosen at random.
function Eject(wshop, item)
	local outputs = rubble.powered.Outputs(wshop)
	if #outputs == 0 then
		opos = rubble.powered.Center(wshop)
		dfhack.items.moveToGround(item, opos)
		return
	end
	
	dfhack.items.moveToGround(item, outputs[math.random(#outputs)])
	rubble.powered.ForbidIfNeeded(item)
end

-- Set an item's quality based on a workshops "skill", aka the quality of it's components.
-- All masterworks gives a skill of 10, all base quality a skill of 0.
-- A skill of 10 gives mostly well-crafted items with smattering of fine
-- and exceptional items (and the rare masterwork and even rarer base quality item).
-- A skill of 0 has trouble producing anything but base quality.
function SetAutoItemQuality(wshop, item)
	local totalQuality = 0
	local partNumber = 0
	for i = 0, #wshop.contained_items - 1, 1 do
		ic = wshop.contained_items[i].item
		-- only take mechanisms and trap components into account
		if ic:getType() == 66 or ic:getType() == 67 then
			partNumber = partNumber + 1
			totalQuality = totalQuality + ic.quality
		end
	end
	local skill
	if partNumber > 0 then
		return math.floor((totalQuality / partNumber) * 2)
	else
		return 0
	end
	
	-- Almost certainly not the proper algorithm, but it works.
	local quality = 0
	if math.random(5) < skill then quality = quality + 1 end
	if math.random(10) < skill then quality = quality + 1 end
	if math.random(15) < skill then quality = quality + 1 end
	if math.random(20) < skill then quality = quality + 1 end
	if math.random(25) < skill and math.random(3) == 1 then quality = quality + 1 end
	if math.random(30) < skill and math.random(3) == 1 then quality = quality + 1 end
	if math.random(35) < skill and math.random(3) == 1 then quality = quality + 1 end
	if math.random(40) < skill and math.random(3) == 1 then quality = quality + 1 end
	if math.random(45) < skill and math.random(3) == 1 then quality = quality + 1 end
	if math.random(50) < skill and math.random(3) == 1 then quality = quality + 1 end
	if quality > 5 then
		quality = 5
	end
	item:setQuality(quality)
end

return _ENV