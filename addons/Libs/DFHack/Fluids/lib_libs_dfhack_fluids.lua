
local _ENV = mkmodule("rubble_fluids")

-- There must be a function like this already, but I couldn't find it.
-- Returns true if the specified tile is downward passable.
function passableDown(x, y, z)
	local ttype = dfhack.maps.getTileType(x, y, z)
	if ttype == 1   or   -- downslope (any kind), why can't the others be like this *sigh*
	   ttype == 32  or   -- open space
	   ttype == 36  or   -- obsidian ud stair
	   ttype == 37  or   -- obsidian d stair
	   ttype == 39  or   -- soil ud stair
	   ttype == 40  or   -- soil d stair
	   ttype == 55  or   -- stone ud stair
	   ttype == 56  or   -- stone d stair
	   ttype == 515 or   -- constructed ud stair
	   ttype == 516 then -- constructed d stair
		return true
	end
	return false
end

-- Eat fluid from the specified tile, returns true if it succeeds.
function eatFluid(x, y, z, magma, amount, minimum)
	local block = dfhack.maps.ensureTileBlock(x,y,z)
	
	if block.designation[x%16][y%16].flow_size >= minimum then
		if block.designation[x%16][y%16].liquid_type == magma then
			block.designation[x%16][y%16].flow_size = block.designation[x%16][y%16].flow_size - amount
		else
			return false
		end
	else
		return false
	end
	
	dfhack.maps.enableBlockUpdates(block,true,true)
	return true
end

-- Eat from below a specified area, there needs to be access to the fluid via a downward passable tile
function eatFromArea(x1, y1, x2, y2, z, magma, amount, minimum)
	for cx = x1, x2, 1 do
		for cy = y1, y2, 1 do
			if passableDown(cx, cy, z) then
				if eatFluid(cx, cy, z-1, magma, amount, minimum) then
					return true
				end
			end
		end
	end
	return false
end

-- spawn fluid
-- returns the amount of fluid spawned
function spawnFluid(x, y, z, magma, amount)
	local block = dfhack.maps.ensureTileBlock(x,y,z)
	local spawned = 0
	
	if block.designation[x%16][y%16].flow_size < 7 then
		if block.designation[x%16][y%16].liquid_type == magma then
			spawn = block.designation[x%16][y%16].flow_size + amount
			spawned = spawn - block.designation[x%16][y%16].flow_size
			if spawn > 7 then spawn = 7 end
			block.designation[x%16][y%16].flow_size = spawn
		else
			return spawned
		end
	else
		return spawned
	end
	
	dfhack.maps.enableBlockUpdates(block,true,true)
	return spawned
end

-- Returns true if item is magma safe.
function magmaSafe(item)
	local mat = dfhack.matinfo.decode(item)
	if mat.material.heat.heatdam_point > 12000 and
	mat.material.heat.melting_point > 12000 and
	mat.material.heat.ignite_point > 12000 and
	mat.material.heat.boiling_point > 12000 then
		return true
	end
	return false
end

-- Find an empty (possibly magma safe) cart in the specified tile
function findCart(x, y, z, magmasafe)
	local itemblock = dfhack.maps.ensureTileBlock(x, y, z)
	if itemblock.occupancy[x%16][y%16].item == true then
		for c=#itemblock.items-1,0,-1 do
			cart=df.item.find(itemblock.items[c])
			if cart:isTrackCart() then
				if cart.pos.x == x and cart.pos.y == y and cart.pos.z == z then
					if #dfhack.items.getContainedItems(cart) == 0 then
						if magmasafe then
							if magmaSafe(cart) then
								return cart
							end
						else
							return cart
						end
					end
				end
			end
		end
	end
	return nil
end

-- Find a empty (possibly magma safe) cart in the specified area.
function findCartArea(x1, y1, x2, y2, z, magmasafe)
	for cx = x1, x2, 1 do
		for cy = y1, y2, 1 do
			cart = findCart(cx, cy, z, magmasafe)
			if cart ~= nil then
				return cart
			end
		end
	end
	return nil
end

-- Fill a minecart with magma or water.
function fillCart(cart, magma)
	capacity = math.floor(cart.subtype.container_capacity/60)
	
	local item=df['item_liquid_miscst']:new()
	item.id=df.global.item_next_id
	df.global.world.items.all:insert('#',item)
	df.global.item_next_id=df.global.item_next_id+1
	
	local mat
	if magma then
		mat = dfhack.matinfo.find("INORGANIC:NONE")
		-- This keeps the magma from becoming solid after a few seconds
		-- apparently items start out at room temperature.
		item.temperature.whole = 12000
	else
		mat = dfhack.matinfo.find("WATER:NONE")
	end
	
	item:setMaterial(mat.type)
	item:setMaterialIndex(mat.index)
	item.stack_size = capacity
	item:categorize(true)
	
	dfhack.items.moveToContainer(item, cart)
end

return _ENV
