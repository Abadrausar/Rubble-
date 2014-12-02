
local event = require "plugins.eventful"
local fluids = rubble.fluids

function fillBarrel(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value=false
	
	if fluids.eatFromArea(unit.pos.x-2, unit.pos.y-2, unit.pos.x+2, unit.pos.y+2, unit.pos.z, false, 1, 2) then
		call_native.value=true
	else
		dfhack.gui.showAnnouncement("Your Still has no nearby water!", COLOR_LIGHTBLUE)
		return
	end
	
	-- Brewer
	fluids.levelUp(unit, 14, 30)
end

event.registerReaction("LUA_HOOK_FILL_BARREL_STILL", fillBarrel)
event.registerReaction("LUA_HOOK_FILL_BUCKET_STILL", fillBarrel)
