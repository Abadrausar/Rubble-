
local powered = rubble.require "powered"
local pitems = rubble.require "powered_items"

function makeMillPlant(output)
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
				if mat.plant.material_defs.type_mill ~= -1 then
					return true
				end
			end
			return false
		end)
		if plant == nil then
			return
		end
		
		local bag = pitems.FindItemAtInput(wshop, function(item)
			if item:isBag() and #dfhack.items.getContainedItems(item) == 0 then
				return true
			end
			return false
		end)
		if bag == nil then
			return
		end
		
		local pmat = dfhack.matinfo.decode(plant)
		local ssize = plant.stack_size
		dfhack.items.remove(plant)
		
		local smat = dfhack.matinfo.decode(pmat.plant.material_defs.type_seed, pmat.plant.material_defs.idx_seed)
		if smat ~= nil then
			for i = 1, ssize, 1 do
				pitems.Eject(wshop, pitems.CreateItem(smat, 'item_seedsst', nil, 0))
			end
		end
		
		local mat = dfhack.matinfo.decode(pmat.plant.material_defs.type_mill, pmat.plant.material_defs.idx_mill)
		local item = pitems.CreateItem(mat, 'item_powder_miscst', nil, 0)
		item:setDimension(150)
		item.stack_size = ssize
		dfhack.items.moveToContainer(item, bag)
		pitems.Eject(wshop, bag)
	end
end

powered.Register("MILL_POWERED", nil, 20, 500, makeMillPlant)
