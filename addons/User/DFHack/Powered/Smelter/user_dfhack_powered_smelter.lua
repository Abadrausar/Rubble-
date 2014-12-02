
function table_contains(t, el)
	for _, v in pairs(t) do
		if v == el then
			return true
		end
	end
	return false
end

function isOre(mat)
	local products = {}
	
	if mat.material.flags.IS_STONE == true then
		if #mat.inorganic.metal_ore.mat_index > 0 then
			for ore=0, #mat.inorganic.metal_ore.mat_index-1, 1 do
				ore_index = mat.inorganic.metal_ore.mat_index[ore]
				ore_prob = mat.inorganic.metal_ore.probability[ore]
				for p = 1, 4, 1 do
					if math.random(100) <= ore_prob then
						table.insert(products, dfhack.matinfo.decode(0, ore_index))
					end
				end
			end
		end
		
		-- Handle ores for "User/Metallurgy/Smelter"
		-- (ignore the "poor" quality silver ores)
		local ores = {
			IRON_ORE = "IRON",
			NICKEL_ORE = "NICKEL",
			GOLD_ORE = "GOLD",
			SILVER_ORE = "SILVER",
			COPPER_ORE = "COPPER",
			LEAD_ORE = "LEAD",
			ZINC_ORE = "ZINC",
			TIN_ORE = "TIN",
			PLATINUM_ORE = "PLATINUM",
			BISMUTH_ORE = "BISMUTH",
			ALUMINUM_ORE = "ALUMINUM"
		}
		
		local rc = mat.material.reaction_class
		for k, v in ipairs(rc) do
			metal = ores[v.value]
			if metal ~= nil then
				local metalmat = dfhack.matinfo.find("INORGANIC:"..metal)
				table.insert(products, metalmat)
				table.insert(products, metalmat)
				table.insert(products, metalmat)
				table.insert(products, metalmat)
			end
		end
	end
	
	if #products ~= 0 then
		return products
	end
	return nil
end

function makeSmeltBoulder(output)
	return function(wshop)
		if not wshop:isUnpowered() then
			if not rubble.powered.HasOutput(wshop) then
				return
			end
			
			local boulder = rubble.powered_items.FindItemAtInput(wshop, function(item)
				if df.item_type[item:getType()] == "BOULDER" then
					local a = isOre(dfhack.matinfo.decode(item))
					if a ~= nil then
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
			
			local products = isOre(dfhack.matinfo.decode(boulder))
			dfhack.items.remove(boulder)
			
			for _, metal in ipairs(products) do
				item = rubble.powered_items.CreateItem(metal, 'item_barst', nil, 0)
				item:setDimension(150)
				rubble.powered_items.Eject(wshop, item)
			end
			
		end
	end
end

rubble.powered.Register("SMELTER_POWERED", nil, 10, 300, makeSmeltBoulder)
