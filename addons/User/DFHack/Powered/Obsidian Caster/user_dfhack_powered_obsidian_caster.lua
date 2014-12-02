
local fluids = rubble.require "fluids"
local powered = rubble.require "powered"
local pitems = rubble.require "powered_items"

function makeCastObsidian(output)
	return function(wshop)
		if not wshop:isUnpowered() then
			if not powered.HasOutput(wshop) then
				return
			end
			local apos = powered.Area(wshop)
			
			if fluids.checkInArea(apos.x1, apos.y1, apos.x2, apos.y2, apos.z, true, 4, 4) then
				if fluids.checkInArea(apos.x1, apos.y1, apos.x2, apos.y2, apos.z, false, 4, 4) then
					fluids.eatFromArea(apos.x1, apos.y1, apos.x2, apos.y2, apos.z, true, 4, 4)
					fluids.eatFromArea(apos.x1, apos.y1, apos.x2, apos.y2, apos.z, false, 4, 4)
					
					local mat = dfhack.matinfo.find("INORGANIC:OBSIDIAN")
					local item = pitems.CreateItem(mat, 'item_boulderst', nil, 0)
					pitems.Eject(wshop, item)
				end
			end
		end
	end
end

powered.Register("OBSIDIAN_CASTER_POWERED", nil, 30, 800, makeCastObsidian)
