
local fluids = rubble.require "fluids"

function makeExtractMetal(output)
	return function(wshop)
		if not wshop:isUnpowered() then
			if not rubble.powered.HasOutput(wshop) then
				return
			end
			
			local apos = rubble.powered.Area(wshop)
			if not fluids.eatFromArea(apos.x1, apos.y1, apos.x2, apos.y2, apos.z, true, 2, 4) then
				--dfhack.gui.showAnnouncement("Your Powered Magma Extractor has run out of magma!", COLOR_LIGHTRED)
				return
			end
			
			local rand = math.random(100)
			local a = 0
			local b = 1
			if rand > a and rand <= b then
				local mat = dfhack.matinfo.find("INORGANIC:ADAMANTINE")
				item = rubble.powered_items.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				rubble.powered_items.Eject(wshop, item)
				return
			end
			
			a = b
			b = b + 3
			if rand > a and rand < b then
				local mat = dfhack.matinfo.find("INORGANIC:GOLD")
				item = rubble.powered_items.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				rubble.powered_items.Eject(wshop, item)
				return
			end
		
			a = b
			b = b + 5
			if rand > a and rand < b then
				local mat = dfhack.matinfo.find("INORGANIC:SILVER")
				item = rubble.powered_items.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				rubble.powered_items.Eject(wshop, item)
				return
			end
		
			a = b
			b = b + 10
			if rand > a and rand < b then
				local mat = dfhack.matinfo.find("INORGANIC:COPPER")
				item = rubble.powered_items.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				rubble.powered_items.Eject(wshop, item)
				return
			end
		
			a = b
			b = b + 2
			if rand > a and rand < b then
				local mat = dfhack.matinfo.find("INORGANIC:PLATINUM")
				item = rubble.powered_items.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				rubble.powered_items.Eject(wshop, item)
				return
			end
		
			a = b
			b = b + 10
			if rand > a and rand < b then
				local mat = dfhack.matinfo.find("INORGANIC:TIN")
				item = rubble.powered_items.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				rubble.powered_items.Eject(wshop, item)
				return
			end
		
			a = b
			b = b + 2
			if rand > a and rand < b then
				local mat = dfhack.matinfo.find("INORGANIC:ALUMINUM")
				item = rubble.powered_items.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				rubble.powered_items.Eject(wshop, item)
				return
			end
		
			a = b
			b = b + 5
			if rand > a and rand < b then
				local mat = dfhack.matinfo.find("INORGANIC:IRON")
				item = rubble.powered_items.CreateItem(mat, 'item_barst', nil, 0)
				item:setDimension(150)
				rubble.powered_items.Eject(wshop, item)
				return
			end
		end
	end
end

rubble.powered.Register("MAGMA_EXTRACTOR_POWERED", nil, 30, 1000, makeExtractMetal)
