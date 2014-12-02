
local buildings = require 'plugins.building-hacks'
local fluids = rubble.fluids

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
	item.flags.removed = true
	return item
end

function makeExtractMetal()
	return function(wshop)
		if not wshop:isUnpowered() then
			local pos={x=wshop.centerx,y=wshop.centery,z=wshop.z}
			
			if not fluids.eatFromArea(pos.x-1, pos.y-1, pos.x+1, pos.y+1, pos.z, true, 2, 4) then
				dfhack.gui.showAnnouncement("Your Powered Magma Extractor has run out of magma!", COLOR_LIGHTRED)
				return
			end
			
			local rand = math.random(100)
			local a = 0
			local b = 1
			if rand > a and rand <= b then
				dfhack.items.moveToGround(makeBar("INORGANIC:ADAMANTINE"), {x=pos.x,y=pos.y,z=pos.z})
			end
			
			a = b
			b = b + 3
			if rand > a and rand < b then
				dfhack.items.moveToGround(makeBar("INORGANIC:GOLD"), {x=pos.x,y=pos.y,z=pos.z})
			end
		
			a = b
			b = b + 5
			if rand > a and rand < b then
				dfhack.items.moveToGround(makeBar("INORGANIC:SILVER"), {x=pos.x,y=pos.y,z=pos.z})
			end
		
			a = b
			b = b + 10
			if rand > a and rand < b then
				dfhack.items.moveToGround(makeBar("INORGANIC:COPPER"), {x=pos.x,y=pos.y,z=pos.z})
			end
		
			a = b
			b = b + 2
			if rand > a and rand < b then
				dfhack.items.moveToGround(makeBar("INORGANIC:PLATINUM"), {x=pos.x,y=pos.y,z=pos.z})
			end
		
			a = b
			b = b + 10
			if rand > a and rand < b then
				dfhack.items.moveToGround(makeBar("INORGANIC:TIN"), {x=pos.x,y=pos.y,z=pos.z})
			end
		
			a = b
			b = b + 2
			if rand > a and rand < b then
				dfhack.items.moveToGround(makeBar("INORGANIC:ALUMINUM"), {x=pos.x,y=pos.y,z=pos.z})
			end
		
			a = b
			b = b + 5
			if rand > a and rand < b then
				dfhack.items.moveToGround(makeBar("INORGANIC:IRON"), {x=pos.x,y=pos.y,z=pos.z})
			end
		end
	end
end

buildings.registerBuilding{
	name="MAGMA_EXTRACTOR_POWERED",
	consume=25,
	gears={{x=0,y=1},{x=1,y=0},{x=2,y=1},{x=1,y=2}},
	action={1000, makeExtractMetal()},
	animate={
		isMechanical=true,
		frames={
			{{x=0,y=1,42,0,0,1},{x=1,y=0,42,0,0,1},{x=2,y=1,42,0,0,1},{x=1,y=2,42,0,0,1}},
			{{x=0,y=1,15,0,0,1},{x=1,y=0,15,0,0,1},{x=2,y=1,15,0,0,1},{x=1,y=2,15,0,0,1}}
		}
	}
}
