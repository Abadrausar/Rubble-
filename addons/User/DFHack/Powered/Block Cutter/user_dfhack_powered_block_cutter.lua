
function makeCutBlock(output)
	return function(wshop)
		if not wshop:isUnpowered() then
			if not rubble.powered.HasOutput(wshop) then
				return
			end
			
			local boulder = rubble.powered_items.FindItemAtInput(wshop, function(item)
				if df.item_type[item:getType()] == "BOULDER" then
					local mat = dfhack.matinfo.decode(item)
					if mat.material.flags.ITEMS_HARD and mat.material.flags.IS_STONE then
						return true
					end
					
					-- Does not play nice with fired boulders, so disabled.
					--if not df.global.ui.economic_stone[mat.index] then
					--	return true
					--end
				end
				return false
			end)
			if boulder == nil then
				return
			end
			
			local mat
			if output == "INPUT_MAT" then
				mat = dfhack.matinfo.decode(boulder)
			else
				mat = dfhack.matinfo.find("INORGANIC:BLOCK_"..output)
			end
			
			-- destroy the boulder
			dfhack.items.remove(boulder)
			
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

local outputs = {
	"BLACK",
	"BROWN",
	"LIGHT_BLUE",
	"LIGHT_CYAN",
	"LIGHT_GREEN",
	"LIGHT_GREY",
	"LIGHT_RED",
	"LIGHT_VIOLET",
	"DARK_BLUE",
	"DARK_CYAN",
	"DARK_GREEN",
	"DARK_GREY",
	"DARK_RED",
	"DARK_VIOLET",
	"WHITE",
	"YELLOW",
	"INPUT_MAT"
}

rubble.powered.Register("BLOCK_CUTTER_POWERED", outputs, 20, 300, makeCutBlock)
