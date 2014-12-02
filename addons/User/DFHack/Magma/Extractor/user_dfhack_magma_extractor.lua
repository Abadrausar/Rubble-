
local event = require "plugins.eventful"
local fluids = rubble.require "fluids"

function makeBar(mat)
	local item = df['item_barst']:new()
	item.id = df.global.item_next_id
	df.global.world.items.all:insert('#', item)
	df.global.item_next_id = df.global.item_next_id + 1
	local mat = dfhack.matinfo.find(mat)
	item:setDimension(150)
	item:setMaterial(mat.type)
	item:setMaterialIndex(mat.index)
	item:categorize(true)
	return item
end

function extractMetal(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value=false
	
	if not fluids.eatFromArea(unit.pos.x-1, unit.pos.y-1, unit.pos.x+1, unit.pos.y+1, unit.pos.z, true, 2, 4) then
		dfhack.gui.showAnnouncement("Your Magma Extractor has run out of magma!", COLOR_LIGHTRED)
		return
	end
	
	local rand = math.random(100)
	local a = 0
	local b = 1
	if rand > a and rand <= b then
		out_items:insert('#', makeBar("INORGANIC:ADAMANTINE"))
	end
	
	a = b
	b = b + 3
	if rand > a and rand < b then
		out_items:insert('#', makeBar("INORGANIC:GOLD"))
	end

	a = b
	b = b + 5
	if rand > a and rand < b then
		out_items:insert('#', makeBar("INORGANIC:SILVER"))
	end

	a = b
	b = b + 10
	if rand > a and rand < b then
		out_items:insert('#', makeBar("INORGANIC:COPPER"))
	end

	a = b
	b = b + 2
	if rand > a and rand < b then
		out_items:insert('#', makeBar("INORGANIC:PLATINUM"))
	end

	a = b
	b = b + 10
	if rand > a and rand < b then
		out_items:insert('#', makeBar("INORGANIC:TIN"))
	end

	a = b
	b = b + 2
	if rand > a and rand < b then
		out_items:insert('#', makeBar("INORGANIC:ALUMINUM"))
	end

	a = b
	b = b + 5
	if rand > a and rand < b then
		out_items:insert('#', makeBar("INORGANIC:IRON"))
	end
	
	-- Furnace operator, 1/10 normal exp
	fluids.levelUp(unit, 25, 3)
end

event.registerReaction("LUA_HOOK_MAGMA_EXTRACTOR", extractMetal)
