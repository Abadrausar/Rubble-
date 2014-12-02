
local buildings = require 'plugins.building-hacks'

function findCart(x, y, z)
	local itemblock = dfhack.maps.ensureTileBlock(x, y, z)
	if itemblock.occupancy[x%16][y%16].item == true then
		for c=#itemblock.items-1,0,-1 do
			cart=df.item.find(itemblock.items[c])
			if cart:isTrackCart() then
				if cart.pos.x == x and cart.pos.y == y and cart.pos.z == z then
					return cart
				end
			end
		end
	end
	return nil
end

function findCartArea(x1, y1, x2, y2, z)
	for cx = x1, x2, 1 do
		for cy = y1, y2, 1 do
			cart = findCart(cx, cy, z)
			if cart ~= nil then
				return cart
			end
		end
	end
	return nil
end

function makeLoadCart(ix, iy)
	return function(wshop)
		if not wshop:isUnpowered() then
			local pos={x = wshop.centerx, y = wshop.centery, z = wshop.z}
			local cart = findCartArea(pos.x-1, pos.y-1, pos.x+1, pos.y+1, pos.z)
			if not cart then
				-- No available minecart
				return
			end
			
			local inputBlock = dfhack.maps.ensureTileBlock(pos.x + ix, pos.y + iy, pos.z)
			if inputBlock.occupancy[(pos.x + ix)%16][(pos.y + iy)%16].item == true then
				for i = #inputBlock.items - 1, 0, -1 do
					local item = df.item.find(inputBlock.items[i])
					local ipos = item.pos
					if ipos.x == (pos.x + ix) and ipos.y == (pos.y + iy) and ipos.z == pos.z and item.flags.on_ground == true then
						local cart_capacity = cart.subtype.container_capacity
						local totalvolume = 0
						local cartrefs = cart.general_refs
						for r = 0, #cartrefs - 1, 1 do
							if getmetatable(cartrefs[r])=="general_ref_contains_itemst" then
								totalvolume = totalvolume + df.item.find(cartrefs[r].item_id):getVolume()
							end
						end
						if totalvolume + item:getVolume() <= cart_capacity then
							dfhack.items.moveToContainer(item, cart)
							return
						end
					end
				end
			end
		end
	end
end

local loader_dirs = {
	N = {ix = 0, iy = -1},
	S = {ix = 0, iy = 1},
	E = {ix = 1, iy = 0},
	W = {ix = -1, iy = 0},
}

function registerLoader(dir, data)
	buildings.registerBuilding{
		name="CART_LOADER_POWERED_"..dir,
		consume=10,
		gears={{x=0,y=0}},
		action={10, makeLoadCart(data.ix, data.iy)},
		animate={
			isMechanical=true,
			frames={
				{{x=0,y=0,240,0,7,0}},
				{{x=0,y=0,254,0,7,0}}
			}
		}
	}
	print("    Registered mechanical workshop: CART_LOADER_POWERED_"..dir)
end

for k,v in pairs(loader_dirs) do
    registerLoader(k, v)
end