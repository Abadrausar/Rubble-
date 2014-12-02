
local buildings = require 'plugins.building-hacks'
local fluids = rubble.fluids

function makeFillCart()
	return function(wshop)
		if not wshop:isUnpowered() then
			local pos={x=wshop.centerx,y=wshop.centery,z=wshop.z}
			
			local cart = false
			local magmacart = false
			local magma = false
			local water = false
			local fluid = nil
			
			if fluids.findCartArea(pos.x-2, pos.y-2, pos.x+2, pos.y+2, pos.z, false) then
				cart = true
			end
			if fluids.findCartArea(pos.x-2, pos.y-2, pos.x+2, pos.y+2, pos.z, true) then
				magmacart = true
			end
			if fluids.checkInArea(pos.x-2, pos.y-2, pos.x+2, pos.y+2, pos.z, true, 4, 4) then
				magma = true
			end
			if fluids.checkInArea(pos.x-2, pos.y-2, pos.x+2, pos.y+2, pos.z, false, 4, 4) then
				water = true
			end
			
			if magma and magmacart then
				fluid = true
			elseif water and cart then
				fluid = false
			elseif magma and cart then
				--dfhack.gui.showAnnouncement("Your Powered Cart Filler cannot find an empty magma-safe minecart!", COLOR_LIGHTRED)
				return
			elseif cart then
				--dfhack.gui.showAnnouncement("Your Powered Cart Filler has no nearby fluids!", COLOR_LIGHTRED)
				return
			else
				-- No minecart, just do nothing to stop message spam
				return
			end
			
			cart = fluids.findCartArea(pos.x-2, pos.y-2, pos.x+2, pos.y+2, pos.z, fluid)
			
			-- Handle carts of various sizes
			amount = math.floor(math.floor(cart.subtype.container_capacity/60)/7)
			if amount > 7 then amount = 7 end
			minimum = 4
			if minimum < amount then minimum = amount end
			
			fluids.eatFromArea(pos.x-2, pos.y-2, pos.x+2, pos.y+2, pos.z, fluid, amount, minimum)
			fluids.fillCart(cart, fluid)
		end
	end
end

buildings.registerBuilding{
	name="CART_FILLER_POWERED",
	consume=25,
	gears={{x=0,y=1},{x=1,y=0},{x=2,y=1},{x=1,y=2}},
	action={5, makeFillCart()},
	animate={
		isMechanical=true,
		frames={
			{{x=0,y=1,42,0,0,1},{x=1,y=0,42,0,0,1},{x=2,y=1,42,0,0,1},{x=1,y=2,42,0,0,1}},
			{{x=0,y=1,15,0,0,1},{x=1,y=0,15,0,0,1},{x=2,y=1,15,0,0,1},{x=1,y=2,15,0,0,1}}
		}
	}
}
