
local powered = rubble.require "powered"
local pitems = rubble.require "powered_items"

function makeWeaveCloth(output)
	return function(wshop)
		if wshop:isUnpowered() or powered.ControllerOff(wshop) then
			return
		end
		if not powered.HasOutput(wshop) then
			return
		end
		
		local thread = pitems.FindItemAtInput(wshop, function(item)
			if df.item_type[item:getType()] == "THREAD" then
				return true
			end
			return false
		end)
		if thread == nil then
			return
		end
		local mat = dfhack.matinfo.decode(thread)
		dfhack.items.remove(thread)
		
		local item = pitems.CreateItem(mat, 'item_clothst', nil, 0)
		item:setDimension(10000)
		pitems.Eject(wshop, item)
	end
end

powered.Register("LOOM_POWERED", nil, 20, 500, makeWeaveCloth)
