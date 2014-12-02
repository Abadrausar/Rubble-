
local event = require "plugins.eventful"
local fluids = rubble.require "fluids"

function fillCartMagma(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value=false
	
	cart = fluids.findCartArea(unit.pos.x-2, unit.pos.y-2, unit.pos.x+2, unit.pos.y+2, unit.pos.z, true)
	if cart == nil then
		dfhack.gui.showAnnouncement("Your Cart Filler cannot find an empty magma-safe minecart!", COLOR_LIGHTRED)
		return
	end
	
	-- Handle carts of various sizes
	amount = math.floor(math.floor(cart.subtype.container_capacity/60)/7)
	if amount > 7 then amount = 7 end
	minimum = 4
	if minimum < amount then minimum = amount end
	
	if fluids.eatFromArea(unit.pos.x-2, unit.pos.y-2, unit.pos.x+2, unit.pos.y+2, unit.pos.z, true, amount, minimum) then
		fluids.fillCart(cart, true)
	else
		dfhack.gui.showAnnouncement("Your Cart Filler has run out of magma!", COLOR_LIGHTRED)
		return
	end
	
	-- Pump operator
	fluids.levelUp(unit, 70, 30)
end

function fillCartWater(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value=false
	
	cart = fluids.findCartArea(unit.pos.x-2, unit.pos.y-2, unit.pos.x+2, unit.pos.y+2, unit.pos.z, false)
	if cart == nil then
		dfhack.gui.showAnnouncement("Your Cart Filler cannot find an empty minecart!", COLOR_LIGHTBLUE)
		return
	end
	
	-- Handle carts of various sizes
	amount = math.floor(math.floor(cart.subtype.container_capacity/60)/7)
	if amount > 7 then amount = 7 end
	minimum = 4
	if minimum < amount then minimum = amount end
	
	if fluids.eatFromArea(unit.pos.x-2, unit.pos.y-2, unit.pos.x+2, unit.pos.y+2, unit.pos.z, false, amount, minimum) then
		fluids.fillCart(cart, false)
	else
		dfhack.gui.showAnnouncement("Your Cart Filler has run out of water!", COLOR_LIGHTBLUE)
		return
	end
	
	-- Pump operator
	fluids.levelUp(unit, 70, 30)
end

event.registerReaction("LUA_HOOK_FILL_CART_MAGMA", fillCartMagma)
event.registerReaction("LUA_HOOK_FILL_CART_WATER", fillCartWater)
