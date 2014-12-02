
local fluids = rubble.require "fluids"
local powered = rubble.require "powered"
local pitems = rubble.require "powered_items"

function makeExtractMetal(output)
	return function(wshop)
		if not wshop:isUnpowered() then
			if not powered.HasOutput(wshop) then
				return
			end
			
			local apos = powered.Area(wshop)
			if not fluids.eatFromArea(apos.x1, apos.y1, apos.x2, apos.y2, apos.z, true, 2, 4) then
				return
			end
			
			local rand = math.random(100)
			local a = 0
			local b = 1
			if rand > a and rand <= b then
				local mat = dfhack.matinfo.find("INORGANIC:ADAMANTINE")
				item = pitems.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				pitems.Eject(wshop, item)
				return
			end
			
			a = b
			b = b + 3
			if rand > a and rand < b then
				local mat = dfhack.matinfo.find("INORGANIC:GOLD")
				item = pitems.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				pitems.Eject(wshop, item)
				return
			end
		
			a = b
			b = b + 5
			if rand > a and rand < b then
				local mat = dfhack.matinfo.find("INORGANIC:SILVER")
				item = pitems.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				pitems.Eject(wshop, item)
				return
			end
		
			a = b
			b = b + 10
			if rand > a and rand < b then
				local mat = dfhack.matinfo.find("INORGANIC:COPPER")
				item = pitems.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				pitems.Eject(wshop, item)
				return
			end
		
			a = b
			b = b + 2
			if rand > a and rand < b then
				local mat = dfhack.matinfo.find("INORGANIC:PLATINUM")
				item = pitems.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				pitems.Eject(wshop, item)
				return
			end
		
			a = b
			b = b + 10
			if rand > a and rand < b then
				local mat = dfhack.matinfo.find("INORGANIC:TIN")
				item = pitems.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				pitems.Eject(wshop, item)
				return
			end
		
			a = b
			b = b + 2
			if rand > a and rand < b then
				local mat = dfhack.matinfo.find("INORGANIC:ALUMINUM")
				item = pitems.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				pitems.Eject(wshop, item)
				return
			end
		
			a = b
			b = b + 5
			if rand > a and rand < b then
				local mat = dfhack.matinfo.find("INORGANIC:IRON")
				item = pitems.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				pitems.Eject(wshop, item)
				return
			end
		end
	end
end

powered.Register("MAGMA_EXTRACTOR_POWERED", nil, 50, 2000, makeExtractMetal)
