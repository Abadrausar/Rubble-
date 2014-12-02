
local buildings = require 'plugins.building-hacks'

function makeLoadCart()
	return function(wshop)
		if not wshop:isUnpowered() then
			local size = 1
			
			local cart = rubble.powered_items.FindItemArea(wshop, function(item)
				return item:isTrackCart()
			end)
			if cart == nil then
				return
			end
			
			-- Get how full the cart is
			local totalvolume = 0
			local cartrefs = cart.general_refs
			for r = 0, #cartrefs - 1, 1 do
				if getmetatable(cartrefs[r])=="general_ref_contains_itemst" then
					totalvolume = totalvolume + df.item.find(cartrefs[r].item_id):getVolume()
				end
			end
			
			local item = rubble.powered_items.FindItemAtInput(wshop, function(item)
				if totalvolume + item:getVolume() <= cart.subtype.container_capacity then
					return true
				end
				return false
			end)
			if item == nil then
				return
			end
			
			dfhack.items.moveToContainer(item, cart)
		end
	end
end

buildings.registerBuilding{
	name="CART_LOADER_POWERED",
	consume=10,
	gears={{x=0,y=0}},
	action={10, makeLoadCart()},
	animate={
		isMechanical=true,
		frames={
			{{x=0,y=0,240,0,7,0}},
			{{x=0,y=0,254,0,7,0}}
		}
	}
}
