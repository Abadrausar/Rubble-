
local powered = rubble.require "powered"
local pitems = rubble.require "powered_items"

function makeSpinThread(output)
	return function(wshop)
		if wshop:isUnpowered() or powered.ControllerOff(wshop) then
			return
		end
		if not powered.HasOutput(wshop) then
			return
		end
		
		local plant = pitems.FindItemAtInput(wshop, function(item)
			if df.item_type[item:getType()] == "PLANT" then
				local mat = dfhack.matinfo.decode(item)
				if mat.plant.material_defs.type_thread ~= -1 then
					return true
				end
			end
			return false
		end)
		if plant == nil then
			return
		end
		local pmat = dfhack.matinfo.decode(plant)
		if plant.stack_size > 1 then
			plant.stack_size = plant.stack_size - 1
		else
			dfhack.items.remove(plant)
		end
		
		local smat = dfhack.matinfo.decode(pmat.plant.material_defs.type_seed, pmat.plant.material_defs.idx_seed)
		if smat ~= nil then
			pitems.Eject(wshop, pitems.CreateItem(smat, 'item_seedsst', nil, 0))
		end
		
		local mat = dfhack.matinfo.decode(pmat.plant.material_defs.type_thread, pmat.plant.material_defs.idx_thread)
		local item = pitems.CreateItem(mat, 'item_threadst', nil, 0)
		item:setDimension(15000)
		pitems.Eject(wshop, item)
	end
end

powered.Register("SPINNER_POWERED", nil, 20, 500, makeSpinThread)
