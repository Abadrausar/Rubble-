
local buildings = require 'plugins.building-hacks'
local fluids = rubble.fluids

function makeCastObsidian()
	return function(wshop)
		if not wshop:isUnpowered() then
			local pos={x=wshop.centerx,y=wshop.centery,z=wshop.z}
			
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
					
					-- If the item is not marked as "removed" detachItem will fail, which will cause
					-- dfhack.items.moveToGround to fail as well. May be a bug?
					item.flags.removed=true
					
					dfhack.items.moveToGround(item, {x=pos.x,y=pos.y,z=pos.z})
				else
					dfhack.gui.showAnnouncement("Your Powered Obsidian Caster has run out of water!", COLOR_LIGHTBLUE)
					return
				end
			else
				dfhack.gui.showAnnouncement("Your Powered Obsidian Caster has run out of magma!", COLOR_LIGHTRED)
				return
			end
		end
	end
end

buildings.registerBuilding{
	name="OBSIDIAN_CASTER_POWERED",
	consume=25,
	gears={{x=0,y=1},{x=1,y=0},{x=2,y=1},{x=1,y=2}},
	action={500, makeCastObsidian()},
	animate={
		isMechanical=true,
		frames={
			{{x=0,y=1,42,0,0,1},{x=1,y=0,42,0,0,1},{x=2,y=1,42,0,0,1},{x=1,y=2,42,0,0,1}},
			{{x=0,y=1,15,0,0,1},{x=1,y=0,15,0,0,1},{x=2,y=1,15,0,0,1},{x=1,y=2,15,0,0,1}}
		}
	}
}
