
local event = require "plugins.eventful"
local fluids = rubble.require "fluids"

function castObsidian(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value=false
	
	if fluids.checkInArea(pos.x-2, pos.y-2, pos.x+2, pos.y+2, pos.z, true, 4, 4) then
		if fluids.checkInArea(pos.x-2, pos.y-2, pos.x+2, pos.y+2, pos.z, false, 4, 4) then
			fluids.eatFromArea(pos.x-2, pos.y-2, pos.x+2, pos.y+2, pos.z, true, 4, 4)
			fluids.eatFromArea(pos.x-2, pos.y-2, pos.x+2, pos.y+2, pos.z, false, 4, 4)
			
			local item=df['item_boulderst']:new()
			item.id=df.global.item_next_id
			df.global.world.items.all:insert('#',item)
			df.global.item_next_id=df.global.item_next_id+1
			local mat = dfhack.matinfo.find("INORGANIC:OBSIDIAN")
			item.stack_size = 1
			item:setMaterial(mat.type)
			item:setMaterialIndex(mat.index)
			item:categorize(true)
			
			out_items:insert('#',item)
			
		else
			dfhack.gui.showAnnouncement("Your Obsidian Caster has run out of water!", COLOR_LIGHTBLUE)
			return
		end
	else
		dfhack.gui.showAnnouncement("Your Obsidian Caster has run out of magma!", COLOR_LIGHTRED)
		return
	end
	
	-- Pump operator
	fluids.levelUp(unit, 70, 30)
end

event.registerReaction("LUA_HOOK_CAST_OBSIDIAN", castObsidian)
