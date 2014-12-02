
local powered = rubble.require "powered"
local pitems = rubble.require "powered_items"

function makeBurnWood(output)
	return function(wshop)
		if wshop:isUnpowered() or powered.ControllerOff(wshop) then
			return
		end
		if not powered.HasOutput(wshop) then
			return
		end
		
		local plank = pitems.FindItemAtInput(wshop, function(block)
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
		local item = pitems.CreateItem(mat, 'item_barst', nil, 0)
		item:setDimension(150)
		pitems.Eject(wshop, item)
	end
end

local outputs = {
	"ASH",
	"CHARCOAL"
}

powered.Register("WOOD_FURNACE_POWERED", outputs, 20, 500, makeBurnWood)
