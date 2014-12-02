
local fluids = rubble.require "fluids"

function makeCastObsidian(output)
	return function(wshop)
		if not wshop:isUnpowered() then
			if not rubble.powered.HasOutput(wshop) then
				return
			end
			local apos = rubble.powered.Area(wshop)
			
			if fluids.checkInArea(apos.x1, apos.y1, apos.x2, apos.y2, apos.z, true, 4, 4) then
				if fluids.checkInArea(apos.x1, apos.y1, apos.x2, apos.y2, apos.z, false, 4, 4) then
					fluids.eatFromArea(apos.x1, apos.y1, apos.x2, apos.y2, apos.z, true, 4, 4)
					fluids.eatFromArea(apos.x1, apos.y1, apos.x2, apos.y2, apos.z, false, 4, 4)
					
					local mat = dfhack.matinfo.find("INORGANIC:OBSIDIAN")
					item = rubble.powered_items.CreateItem(mat, 'item_boulderst', nil, 0)
					rubble.powered_items.Eject(wshop, item)
				else
					--dfhack.gui.showAnnouncement("Your Powered Obsidian Caster has run out of water!", COLOR_LIGHTBLUE)
					return
				end
			else
				--dfhack.gui.showAnnouncement("Your Powered Obsidian Caster has run out of magma!", COLOR_LIGHTRED)
				return
			end
		end
	end
end

rubble.powered.Register("OBSIDIAN_CASTER_POWERED", nil, 25, 500, makeCastObsidian)
