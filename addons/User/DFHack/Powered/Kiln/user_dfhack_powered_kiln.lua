
function isClay(mat)
	local rp = mat.material.reaction_product
	for k, v in ipairs(rp.id) do
		if v.value == "FIRED_MAT" then
			return dfhack.matinfo.decode(rp.material.mat_type[k], rp.material.mat_index[k])
		end
	end
	return nil
end

function makeFireBoulder(output)
	return function(wshop)
		if not wshop:isUnpowered() then
			if not rubble.powered.HasOutput(wshop) then
				return
			end
			
			if output == "PEARLASH" then
				local ash = rubble.powered_items.FindItemAtInput(wshop, function(item)
					if df.item_type[item:getType()] == "BAR" then
						local mat = dfhack.matinfo.decode(item)
						local ashmat = dfhack.matinfo.find("ASH:NONE")
						if mat.type == ashmat.type then
							return true
						end
					end
					return false
				end)
				if ash == nil then
					return
				end
				
				magma, bar = rubble.powered_items.FindFuel(wshop)
				if not magma then
					if bar == nil then
						return
					end
					
					dfhack.items.remove(bar)
				end
				
				dfhack.items.remove(ash)
				
				item = rubble.powered_items.CreateItem(dfhack.matinfo.find("PEARLASH:NONE"), 'item_barst', nil, 0)
				item:setDimension(150)
				rubble.powered_items.Eject(wshop, item)
			else
				local boulder = rubble.powered_items.FindItemAtInput(wshop, function(item)
					if df.item_type[item:getType()] == "BOULDER" then
						local mat = isClay(dfhack.matinfo.decode(item))
						if mat ~= nil then
							return true
						end
					end
					return false
				end)
				if boulder == nil then
					return
				end
				
				magma, bar = rubble.powered_items.FindFuel(wshop)
				if not magma then
					if bar == nil then
						return
					end
					
					dfhack.items.remove(bar)
				end
				
				local mat = isClay(dfhack.matinfo.decode(boulder))
				dfhack.items.remove(boulder)
				
				item = rubble.powered_items.CreateItem(mat, 'item_boulderst', nil, 0)
				rubble.powered_items.Eject(wshop, item)
			end
		end
	end
end

local outputs = {
	"CERAMIC",
	"PEARLASH"
}

rubble.powered.Register("KILN_POWERED", outputs, 10, 300, makeFireBoulder)
