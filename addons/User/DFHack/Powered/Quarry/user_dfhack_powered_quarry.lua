
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
		if not wshop:isUnpowered() then
			if not rubble.powered.HasOutput(wshop) then
				return
			end
			local apos = rubble.powered.Area(wshop)
			
			if output == "SAND" then
				local mat = findSandArea(apos.x1, apos.y1, apos.x2, apos.y2, apos.z)
				if mat == nil then
					return
				end
				
				local bag = rubble.powered_items.FindItemAtInput(wshop, function(item)
					if item:isBag() and #dfhack.items.getContainedItems(item) == 0 then
						return true
					end
					return false
				end)
				if bag == nil then
					return
				end
				
				local item = rubble.powered_items.CreateItem(mat, 'item_powder_miscst', nil, 0)
				dfhack.items.moveToContainer(item, bag)
				rubble.powered_items.Eject(wshop, bag)
			else
				local mat = findClayArea(apos.x1, apos.y1, apos.x2, apos.y2, apos.z)
				if mat == nil then
					return
				end
				
				item = rubble.powered_items.CreateItem(mat, 'item_boulderst', nil, 0)
				rubble.powered_items.Eject(wshop, item)
			end
		end
	end
end

local outputs = {
	"CLAY",
	"SAND"
}

rubble.powered.Register("QUARRY_POWERED", outputs, 10, 500, makeDigQuarry)
