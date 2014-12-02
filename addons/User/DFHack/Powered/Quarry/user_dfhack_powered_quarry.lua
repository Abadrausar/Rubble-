
local powered = rubble.require "powered"
local pitems = rubble.require "powered_items"

function findClay(x, y, z)
	block = dfhack.maps.ensureTileBlock(x, y, z)
	
	-- I can't figure out how to get a tile's material,
	-- is it even possible?
	-- Let's fake it and hope users won't notice...
	if df.tiletype_material[df.tiletype.attrs[block.tiletype[x%16][y%16]].material] == "SOIL" then
		return dfhack.matinfo.find("INORGANIC:CLAY")
	end
	
	return nil
end

function findClayArea(x1, y1, x2, y2, z)
	for cx = x1, x2, 1 do
		for cy = y1, y2, 1 do
			clay = findClay(cx, cy, z)
			if clay ~= nil then
				return clay
			end
		end
	end
	return nil
end

function findSandArea(x1, y1, x2, y2, z)
	-- Since we can't tell soil type may as well...
	if findClayArea(x1, y1, x2, y2, z) ~= nil then
		return dfhack.matinfo.find("INORGANIC:SAND_TAN")
	end
	return nil
end

function makeDigQuarry(output)
	return function(wshop)
		if wshop:isUnpowered() or powered.ControllerOff(wshop) then
			return
		end
		if not powered.HasOutput(wshop) then
			return
		end
		local apos = powered.Area(wshop)
		
		if output == "SAND" then
			local mat = findSandArea(apos.x1, apos.y1, apos.x2, apos.y2, apos.z)
			if mat == nil then
				return
			end
			
			local bag = pitems.FindItemAtInput(wshop, function(item)
				if item:isBag() and #dfhack.items.getContainedItems(item) == 0 then
					return true
				end
				return false
			end)
			if bag == nil then
				return
			end
			
			local item = pitems.CreateItem(mat, 'item_powder_miscst', nil, 0)
			dfhack.items.moveToContainer(item, bag)
			pitems.Eject(wshop, bag)
		else
			local mat = findClayArea(apos.x1, apos.y1, apos.x2, apos.y2, apos.z)
			if mat == nil then
				return
			end
			
			local item = pitems.CreateItem(mat, 'item_boulderst', nil, 0)
			pitems.Eject(wshop, item)
		end
	end
end

local outputs = {
	"CLAY",
	"SAND"
}

powered.Register("QUARRY_POWERED", outputs, 30, 800, makeDigQuarry)
