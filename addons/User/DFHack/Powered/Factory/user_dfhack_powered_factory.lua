
local eventful = require "plugins.eventful"
local buildings = require 'plugins.building-hacks'

function eatBlocks(wshop, stone)
	local items = {}
	local found = {}
	for i = 0, 7, 1 do
		block = rubble.powered_items.FindItemAtInput(wshop, function(item)
			if df.item_type[item:getType()] == "BLOCKS" then
				if found[item.id] == true then
					return false
				end
				
				local mat = dfhack.matinfo.decode(item)
				if stone and mat.mode == "inorganic" then
					return true
				elseif not stone and mat.mode == "plant" then
					return true
				end
			end
			return false
		end)
		if block == nil then
			return nil
		end
		found[block.id] = true
		items[i] = block
	end
	
	local mat = nil
	for i, block in ipairs(items) do
		mat = dfhack.matinfo.decode(block)
		dfhack.items.remove(block)
	end
	return mat
end

function eatGlass(wshop)
	local items = {}
	for i = 0, 3, 1 do
		glass = rubble.powered_items.FindItemAtInput(wshop, function(item)
			if df.item_type[item:getType()] == "ROUGH" then
				if item.flags.removed == true then
					return false
				end
				
				local mat = dfhack.matinfo.decode(item)
				if mat.mode == "builtin" then
					return true
				end
			end
			return false
		end)
		if glass == nil then
			for i, glass in ipairs(items) do
				glass.flags.removed = false
			end
			return nil
		end
		glass.flags.removed = true
		items[i] = glass
	end
	
	local mat = nil
	for i, glass in ipairs(items) do
		mat = dfhack.matinfo.decode(glass)
		glass.flags.removed = false
		dfhack.items.remove(glass)
	end
	return mat
end

function levelUp(unit, skillId, amount)
	max_skill = 20 
	
	local skill = df.unit_skill:new()
	local foundSkill = false
	for k, soulSkill in ipairs(unit.status.current_soul.skills) do
		if soulSkill.id == skillId then
			skill = soulSkill
			foundSkill = true
			break
		end
	end
 
	if foundSkill then
		-- Let's not train beyond the max skill
		if skill.rating >= max_skill then
			return false
		end
 
		skill.experience = skill.experience + amount
		if skill.experience > 100 * skill.rating + 500 then
			skill.experience = skill.experience - (100 * skill.rating + 500)
			skill.rating = skill.rating + 1
		end
	else
		skill.id = skillId
		skill.experience = amount
		skill.rating = 0
		unit.status.current_soul.skills:insert('#',skill)
	end
	
	return true
end

function factoryStone(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value=false
	local fake_wshop = rubble.powered.MakeFake(unit.pos.x, unit.pos.y, unit.pos.z, 5)
	
	local mat = eatBlocks(fake_wshop, true)
	if mat ~= nil then
		for x = 1, reaction.products[0].count, 1 do
			local item = rubble.powered_items.CreateItemNumeric(mat, reaction.products[0].item_type, reaction.products[0].item_subtype, unit, reaction.skill)
			rubble.powered_items.Eject(fake_wshop, item)
		end
		
		levelUp(unit, reaction.skill, 100)
		
		return
	end
	
	dfhack.gui.showAnnouncement("Your Powered Factory has run out of blocks!", COLOR_LIGHTRED)
end

function factoryWood(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value=false
	local fake_wshop = rubble.powered.MakeFake(unit.pos.x, unit.pos.y, unit.pos.z, 5)
	
	local mat = eatBlocks(fake_wshop, false)
	if mat ~= nil then
		for x = 1, reaction.products[0].count, 1 do
			local item = rubble.powered_items.CreateItemNumeric(mat, reaction.products[0].item_type, reaction.products[0].item_subtype, unit, reaction.skill)
			rubble.powered_items.Eject(fake_wshop, item)
		end
		
		levelUp(unit, reaction.skill, 100)
		
		return
	end
	
	dfhack.gui.showAnnouncement("Your Powered Factory has run out of planks!", COLOR_LIGHTRED)
end

function factoryGlass(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value=false
	local fake_wshop = rubble.powered.MakeFake(unit.pos.x, unit.pos.y, unit.pos.z, 5)
	
	local mat = eatGlass(fake_wshop)
	if mat ~= nil then
		for x = 1, reaction.products[0].count, 1 do
			local item = rubble.powered_items.CreateItemNumeric(mat, reaction.products[0].item_type, reaction.products[0].item_subtype, unit, reaction.skill)
			rubble.powered_items.Eject(fake_wshop, item)
		end
		
		levelUp(unit, reaction.skill, 100)
		
		return
	end
	
	dfhack.gui.showAnnouncement("Your Powered Factory has run out of raw glass!", COLOR_LIGHTRED)
end

for i,reaction in ipairs(df.global.world.raws.reactions) do
	if string.match(reaction.code,'^LUA_HOOK_FACTORY_BLOCK') then
		eventful.registerReaction(reaction.code, factoryStone)
	elseif string.match(reaction.code,'^LUA_HOOK_FACTORY_PLANK') then
		eventful.registerReaction(reaction.code, factoryWood)
	elseif string.match(reaction.code,'^LUA_HOOK_FACTORY_GLASS') then
		eventful.registerReaction(reaction.code, factoryGlass)
	end
end

buildings.registerBuilding{
	name="FACTORY_POWERED",
	consume=50,
	gears={{x=0,y=4},{x=4,y=0},{x=4,y=4},{x=0,y=0}},
	animate={
		isMechanical=true,
		frames={
			{{x=0,y=4,42,0,0,1},{x=4,y=0,42,0,0,1},{x=4,y=4,42,0,0,1},{x=0,y=0,42,0,0,1}},
			{{x=0,y=4,15,0,0,1},{x=4,y=0,15,0,0,1},{x=4,y=4,15,0,0,1},{x=0,y=0,15,0,0,1}}
		}
	}
}
