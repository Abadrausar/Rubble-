
local buildings = require 'plugins.building-hacks'

-- Returns carts that are >= 75% full
function findCart(x, y, z)
	local itemblock = dfhack.maps.ensureTileBlock(x, y, z)
	if itemblock.occupancy[x%16][y%16].item == true then
		for c=#itemblock.items-1,0,-1 do
			cart=df.item.find(itemblock.items[c])
			if cart:isTrackCart() then
				if cart.pos.x == x and cart.pos.y == y and cart.pos.z == z then
					local cart_capacity = cart.subtype.container_capacity
					local totalvolume = 0
					local cartrefs = cart.general_refs
					for r = 0, #cartrefs - 1, 1 do
						if getmetatable(cartrefs[r])=="general_ref_contains_itemst" then
							totalvolume = totalvolume + df.item.find(cartrefs[r].item_id):getVolume()
						end
					end
					if (cart_capacity - totalvolume) <= (cart_capacity / 4) then
						return cart
					end
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

function makeLaunchCart(vx, vy)
	return function(wshop)
		if not wshop:isUnpowered() then
			local pos = {x = wshop.centerx, y = wshop.centery, z = wshop.z}
			local cart = findCartArea(pos.x-1, pos.y-1, pos.x+1, pos.y+1, pos.z)
			if not cart then
				-- No available minecart
				return
			end
			
			local vehicle = df.global.world.vehicles.all[cart.vehicle_id]
			
			-- It may be possible to see if vehicle.time_stopped > 0 instead...
			if vehicle.speed_x == 0 and vehicle.speed_y == 0 and vehicle.speed_z == 0 then
				-- Launch!
				vehicle.speed_x = vx
				vehicle.speed_y = vy
			end
		end
	end
end

local launcher_dirs = {
	N = {vx = 0, vy = -20000},
	S = {vx = 0, vy = 20000},
	E = {vx = 20000, vy = 0},
	W = {vx = -20000, vy = 0},
}

function registerLauncher(dir, data)
	buildings.registerBuilding{
		name="CART_LAUNCHER_POWERED_"..dir,
		consume=5,
		gears={{x=0,y=0}},
		action={50, makeLaunchCart(data.vx, data.vy)},
		animate={
			isMechanical=true,
			frames={
				{{x=0,y=0,42,7,0,0}},
				{{x=0,y=0,15,7,0,0}}
			}
		}
	}
	print("    Registered mechanical workshop: CART_LAUNCHER_POWERED_"..dir)
end

for k,v in pairs(launcher_dirs) do
    registerLauncher(k, v)
end