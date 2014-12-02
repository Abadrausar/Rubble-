
local buildings = require 'plugins.building-hacks'
local powered = rubble.require "powered"
local pitems = rubble.require "powered_items"

function makeLoadCart()
	return function(wshop)
		if wshop:isUnpowered() or powered.ControllerOff(wshop) then
			return
		end
		
		local cart
		if powered.HasOutput(wshop) then
			cart = pitems.FindItemAt(powered.Outputs(wshop), function(item)
				return item:isTrackCart()
			end)
		else
			cart = pitems.FindItemArea(wshop, function(item)
				return item:isTrackCart()
			end)
		end
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
		
		local item = pitems.FindItemAtInput(wshop, function(item)
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

buildings.registerBuilding{
	name="CART_LOADER_POWERED",
	consume=5,
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
