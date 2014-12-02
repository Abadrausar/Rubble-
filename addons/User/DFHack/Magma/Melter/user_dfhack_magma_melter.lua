
local event = require "plugins.eventful"
local fluids = rubble.require "fluids"

function fillCartMagma(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value=false
	
	cart = fluids.findCartArea(unit.pos.x-2, unit.pos.y-2, unit.pos.x+2, unit.pos.y+2, unit.pos.z, true)
	if cart == nil then
		dfhack.gui.showAnnouncement("Your Magma Melter cannot find an empty magma-safe minecart!", COLOR_LIGHTRED)
		return
	end
	
	fluids.fillCart(cart, true)
	
	-- Furnace operator
	fluids.levelUp(unit, 25, 30)
end

function toCisternSeven(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value=false
	
	for count = 1, 7, 1 do
		if not fluids.spawnInArea(unit.pos.x-2, unit.pos.y-2, unit.pos.x+2, unit.pos.y+2, unit.pos.z, true, 1) then
			dfhack.gui.showAnnouncement("Your Magma Melter has run out of room to spawn magma!", COLOR_LIGHTRED)
		end
	end
	
	-- Furnace operator
	fluids.levelUp(unit, 25, 30)
end

function toCisternOne(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value=false
	
	if not fluids.spawnInArea(unit.pos.x-2, unit.pos.y-2, unit.pos.x+2, unit.pos.y+2, unit.pos.z, true, 1) then
		dfhack.gui.showAnnouncement("Your Magma Melter has run out of room to spawn magma!", COLOR_LIGHTRED)
	end
	
	-- Furnace operator
	fluids.levelUp(unit, 25, 30)
end

function allowMagma(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value=false
	
	-- Discover the magma sea.
	for index,value in ipairs(df.global.world.cur_savegame.map_features) do
		local featureType=value:getType()
		if df.feature_type[featureType]=="magma_core_from_layer" then
			value.flags[df.feature_init_flags.Discovered]=true
			dfhack.gui.showAnnouncement("Magma workshops/furnaces are now unlocked!", COLOR_LIGHTRED)
			return
		end
	end
end

event.registerReaction("LUA_HOOK_MAGMA_MELTER_CART", fillCartMagma)
event.registerReaction("LUA_HOOK_MAGMA_MELTER_CISTERN_SEVEN", toCisternSeven)
event.registerReaction("LUA_HOOK_MAGMA_MELTER_CISTERN_ONE", toCisternOne)
event.registerReaction("LUA_HOOK_MAGMA_MELTER_ALLOW", allowMagma)
