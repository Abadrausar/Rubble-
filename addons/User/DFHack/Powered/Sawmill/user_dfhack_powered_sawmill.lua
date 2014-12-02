
local powered = rubble.require "powered"
local pitems = rubble.require "powered_items"

function makeCutWood(output)
	return function(wshop)
		if not wshop:isUnpowered() then
			if not powered.HasOutput(wshop) then
				return
			end
			
			-- find a log
			local wood = pitems.FindItemAtInput(wshop, function(item)
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
			pitems.Eject(wshop, pitems.CreateItem(mat, 'item_blocksst', nil, 0))
			pitems.Eject(wshop, pitems.CreateItem(mat, 'item_blocksst', nil, 0))
			pitems.Eject(wshop, pitems.CreateItem(mat, 'item_blocksst', nil, 0))
			pitems.Eject(wshop, pitems.CreateItem(mat, 'item_blocksst', nil, 0))
		end
	end
end

powered.Register("SAWMILL_POWERED", nil, 30, 1000, makeCutWood)
