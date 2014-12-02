
local eventful = require "plugins.eventful"
local buildings = require 'plugins.building-hacks'
local script = require 'gui.script'
local powered = rubble.require "powered"
local pitems = rubble.require "powered_items"
local ppersist = rubble.require "powered_persist"

function isClay(mat)
	local rp = mat.material.reaction_product
	for k, v in ipairs(rp.id) do
		if v.value == "FIRED_MAT" then
			return dfhack.matinfo.decode(rp.material.mat_type[k], rp.material.mat_index[k])
		end
	end
	return nil
end

function getItemCaption(type, subtype)
    local attrs = df.item_type.attrs[type]
    if subtype == -1 or subtype == nil then
		return attrs.caption
	else
		return df['itemdef_'..string.lower(df.item_type[type])..'st'].find(subtype).name
	end
end

alreadyAdjusting = false
function kilnAdjust(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value = false
	
	if alreadyAdjusting == false then
		alreadyAdjusting = true
		
		local opt_types = {"special = 'PEARLASH'"}
		local opt_names = {"pearlash"}
		
		for i = 0, #reaction.products-1, 1 do
			option = reaction.products[i]
			table.insert(opt_types, "itype = "..option.item_type..", isubtype = "..option.item_subtype)
			table.insert(opt_names, getItemCaption(option.item_type, option.item_subtype))
		end
		
		script.start(function()
			local choiceok, choice = script.showListPrompt('Glass Furnace', 'Select item to produce:', COLOR_LIGHTGREEN, opt_names)
			
			local product = ""
			if choiceok then
				product = opt_types[choice]
			else
				product = "NONE"
			end
			
			local wshop = powered.MakeFake(unit.pos.x, unit.pos.y, unit.pos.z, 3)
			ppersist.SetOutputType(wshop, "return {"..product.."}")
			
			alreadyAdjusting = false
		end)
	end
end

function makeFireClay()
	return function(wshop)
		if not wshop:isUnpowered() then
			if not powered.HasOutput(wshop) then
				return
			end
			
			local output = ppersist.GetOutputTypeAsCode(wshop)
			if output == nil then
				return
			end
			
			if output.special == "PEARLASH" then
				local ash = pitems.FindItemAtInput(wshop, function(item)
					if df.item_type[item:getType()] == "BAR" then
						local mat = dfhack.matinfo.decode(item)
						local ashmat = dfhack.matinfo.find("ASH:NONE")
						if mat.type == ashmat.type then
							return true
						end
					end
					return false
				end)
				if ash == nil then
					return
				end
				
				local magma, bar = pitems.FindFuel(wshop)
				if not magma then
					if bar == nil then
						return
					end
					
					dfhack.items.remove(bar)
				end
				
				dfhack.items.remove(ash)
				
				local item = pitems.CreateItem(dfhack.matinfo.find("PEARLASH:NONE"), 'item_barst', nil, 0)
				item:setDimension(150)
				pitems.Eject(wshop, item)
				return
			end
			
			local boulder = pitems.FindItemAtInput(wshop, function(item)
				if df.item_type[item:getType()] == "BOULDER" then
					local mat = isClay(dfhack.matinfo.decode(item))
					if mat ~= nil then
						return true
					end
				end
				return false
			end)
			if boulder == nil then
				return
			end
			
			local magma, bar = pitems.FindFuel(wshop)
			if not magma then
				if bar == nil then
					return
				end
				
				dfhack.items.remove(bar)
			end
			
			local mat = isClay(dfhack.matinfo.decode(boulder))
			dfhack.items.remove(boulder)
			
			local item = pitems.CreateItemNumeric(mat, output.itype, output.isubtype, nil, 0)
			pitems.SetAutoItemQuality(wshop, item)
			pitems.Eject(wshop, item)
		end
	end
end

buildings.registerBuilding{
	name="DFHACK_RUBBLE_POWERED_KILN",
	consume=20,
	gears={{x=0,y=2},{x=2,y=0},{x=2,y=2},{x=0,y=0}},
	action={500, makeFireClay()},
	animate={
		isMechanical=true,
		frames={
			{{x=0,y=2,42,0,0,1},{x=2,y=0,42,0,0,1},{x=2,y=2,42,0,0,1},{x=0,y=0,42,0,0,1}},
			{{x=0,y=2,15,0,0,1},{x=2,y=0,15,0,0,1},{x=2,y=2,15,0,0,1},{x=0,y=0,15,0,0,1}}
		}
	}
}
eventful.registerReaction("LUA_HOOK_ADJUST_KILN", kilnAdjust)
print("    Registered mechanical workshop: DFHACK_RUBBLE_POWERED_GLASS_FURNACE")
