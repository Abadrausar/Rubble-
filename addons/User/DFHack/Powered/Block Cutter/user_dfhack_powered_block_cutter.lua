
local powered = rubble.require "powered"
local pitems = rubble.require "powered_items"

function makeCutBlock(output)
	return function(wshop)
		if wshop:isUnpowered() or powered.ControllerOff(wshop) then
			return
		end
		if not powered.HasOutput(wshop) then
			return
		end
		
		local boulder = pitems.FindItemAtInput(wshop, function(item)
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
		pitems.Eject(wshop, pitems.CreateItem(mat, 'item_blocksst', nil, 0))
		pitems.Eject(wshop, pitems.CreateItem(mat, 'item_blocksst', nil, 0))
		pitems.Eject(wshop, pitems.CreateItem(mat, 'item_blocksst', nil, 0))
		pitems.Eject(wshop, pitems.CreateItem(mat, 'item_blocksst', nil, 0))
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

powered.Register("BLOCK_CUTTER_POWERED", outputs, 30, 1000, makeCutBlock)
