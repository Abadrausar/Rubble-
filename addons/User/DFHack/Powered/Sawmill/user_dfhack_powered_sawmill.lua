
function makeCutWood(output)
	return function(wshop)
		if not wshop:isUnpowered() then
			if not rubble.powered.HasOutput(wshop) then
				return
			end
			
			-- find a log
			local wood = rubble.powered_items.FindItemAtInput(wshop, function(item)
				if df.item_type[item:getType()] == "WOOD" then
					return true
				end
				return false
			end)
			if not wood then
				return
			end
			
			local mat = dfhack.matinfo.decode(wood)
			
			-- destroy the log
			dfhack.items.remove(wood)
			
			-- create 4 blocks
			item = rubble.powered_items.CreateItem(mat, 'item_blocksst', nil, 0)
			rubble.powered_items.Eject(wshop, item)
			
			item = rubble.powered_items.CreateItem(mat, 'item_blocksst', nil, 0)
			rubble.powered_items.Eject(wshop, item)
			
			item = rubble.powered_items.CreateItem(mat, 'item_blocksst', nil, 0)
			rubble.powered_items.Eject(wshop, item)
			
			item = rubble.powered_items.CreateItem(mat, 'item_blocksst', nil, 0)
			rubble.powered_items.Eject(wshop, item)
		end
	end
end

rubble.powered.Register("SAWMILL_POWERED", nil, 20, 300, makeCutWood)
