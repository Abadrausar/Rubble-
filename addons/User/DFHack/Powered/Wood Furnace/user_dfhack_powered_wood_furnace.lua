
function makeBurnWood(output)
	return function(wshop)
		if not wshop:isUnpowered() then
			if not rubble.powered.HasOutput(wshop) then
				return
			end
			
			plank = rubble.powered_items.FindItemAtInput(wshop, function(block)
				if df.item_type[block:getType()] == "BLOCKS" then
					local mat = dfhack.matinfo.decode(block)
					if mat.mode == "plant" then
						return true
					end
				end
				return false
			end)
			if plank == nil then
				return
			end
			dfhack.items.remove(plank)
			
			local mat
			if output == "ASH" then
				mat = dfhack.matinfo.find("ASH:NONE")
			else
				mat = dfhack.matinfo.find("COAL:CHARCOAL")
			end
			item = rubble.powered_items.CreateItem(mat, 'item_barst', nil, 0)
			item:setDimension(150)
			rubble.powered_items.Eject(wshop, item)
		end
	end
end

local outputs = {
	"ASH",
	"CHARCOAL"
}

rubble.powered.Register("WOOD_FURNACE_POWERED", outputs, 10, 200, makeBurnWood)
