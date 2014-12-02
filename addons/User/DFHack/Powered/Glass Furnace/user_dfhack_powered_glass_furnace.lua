
local eventful = require "plugins.eventful"
local buildings = require 'plugins.building-hacks'
local script = require('gui.script')
		
function getItemCaption(type, subtype)
    local attrs = df.item_type.attrs[type]
    if subtype == -1 or subtype == nil then
		return attrs.caption
	else
		return df['itemdef_'..string.lower(df.item_type[type])..'st'].find(subtype).name
	end
end

alreadyAdjusting = false
function furnaceAdjust(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value = false
	
	if alreadyAdjusting == false then
		alreadyAdjusting = true
		
		local opt_types = {}
		local opt_names = {}
		
		for i = 0, #reaction.products-1, 1 do
			option = reaction.products[i]
			table.insert(opt_types, "itype = "..option.item_type..", isubtype = "..option.item_subtype)
			table.insert(opt_names, getItemCaption(option.item_type,option.item_subtype))
		end
		
		script.start(function()
			local mats = {"GREEN", "CLEAR"}
			
			local matok, mat = script.showListPrompt('Glass Furnace', 'Select material:', COLOR_LIGHTGREEN, mats)
			local choiceok, choice = script.showListPrompt('Glass Furnace', 'Select item to produce:', COLOR_LIGHTGREEN, opt_names)
			
			local product = ""
			if choiceok and matok then
				product = opt_types[choice]..", mat = "..mats[mat]
			else
				product = "NONE"
			end
			
			local wshop = rubble.powered.MakeFake(unit.pos.x, unit.pos.y, unit.pos.z, 1)
			rubble.powered_persist.SetOutputType(wshop, "return {"..product.."}")
			
			alreadyAdjusting = false
		end)
	end
end

function makeGlass()
	return function(wshop)
		if not wshop:isUnpowered() then
			if not rubble.powered.HasOutput(wshop) then
				return
			end
			
			local outputraw = rubble.powered_persist.GetOutputType(wshop)
			if outputraw == "NONE" then
				return
			end
			
			local f, err = load(outputraw)
			if f == nil then
				error(err)
				return
			end
			local output = f()
			
			local sand = rubble.powered_items.FindItemAtInput(wshop, function(item)
				if item:isBag() then
					local contents = dfhack.items.getContainedItems(item)
					if #contents == 1 and contents[1]:isSand() then
						return true
					end
				end
				return false
			end)
			if sand == nil then
				return
			end
			
			local pearlash = nil
			if output.mat == "CLEAR" then
				pearlash = rubble.powered_items.FindItemAtInput(wshop, function(item)
					if df.item_type[item:getType()] == "BAR" then
						local mat = dfhack.matinfo.decode(item)
						local ashmat = dfhack.matinfo.find("PEARLASH:NONE")
						if mat.type == ashmat.type then
							return true
						end
					end
					return false
				end)
				if pearlash == nil then
					return
				end
			end
			
			magma, bar = rubble.powered_items.FindFuel(wshop)
			if not magma then
				if bar == nil then
					return
				end
				
				dfhack.items.remove(bar)
			end
			
			if pearlash ~= nil then
				dfhack.items.remove(pearlash)
			end
			
			-- empty the sand bag and eject it
			for r = #sand.general_refs - 1, 0, -1 do
				if getmetatable(sand.general_refs[r]) == 'general_ref_contains_itemst' then
					local contained_item = df.item.find(sand.general_refs[r].item_id)
					for r2 = #contained_item.general_refs-1, 0, -1 do
						if getmetatable(contained_item.general_refs[r2]) == 'general_ref_contained_in_itemst' then
							contained_item.general_refs:erase(r2)
							contained_item.flags.in_inventory = false
						end
					end
					sand.general_refs:erase(r)
					dfhack.items.remove(contained_item)
				end
			end
			rubble.powered_items.Eject(wshop, sand)
			
			local glass = dfhack.matinfo.find("GLASS_"..output.mat..":NONE")
			item = rubble.powered_items.CreateItemNumeric(glass, output.itype, output.isubtype, nil, 0)
			rubble.powered_items.SetAutoItemQuality(wshop, item)
			rubble.powered_items.Eject(wshop, item)
		end
	end
end

buildings.registerBuilding{
	name="DFHACK_RUBBLE_POWERED_GLASS_FURNACE",
	consume=10,
	gears={{x=0,y=2},{x=2,y=0},{x=2,y=2},{x=0,y=0}},
	action={300, makeGlass()},
	animate={
		isMechanical=true,
		frames={
			{{x=0,y=2,42,0,0,1},{x=2,y=0,42,0,0,1},{x=2,y=2,42,0,0,1},{x=0,y=0,42,0,0,1}},
			{{x=0,y=2,15,0,0,1},{x=2,y=0,15,0,0,1},{x=2,y=2,15,0,0,1},{x=0,y=0,15,0,0,1}}
		}
	}
}
eventful.registerReaction("LUA_HOOK_ADJUST_GLASS_FURNACE", furnaceAdjust)
print("    Registered mechanical workshop: DFHACK_RUBBLE_POWERED_GLASS_FURNACE")
