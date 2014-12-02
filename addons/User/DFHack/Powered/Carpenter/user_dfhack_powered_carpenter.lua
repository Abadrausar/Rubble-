
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
function carpenterAdjust(reaction, unit, in_items, in_reag, out_items, call_native)
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
			local choiceok, choice = script.showListPrompt('Carpenter', 'Select item to produce:', COLOR_LIGHTGREEN, opt_names)
			
			local product = ""
			if choiceok then
				product = opt_types[choice]
			else
				product = "NONE"
			end
			
			local wshop = rubble.powered.MakeFake(unit.pos.x, unit.pos.y, unit.pos.z, 1)
			rubble.powered_persist.SetOutputType(wshop, "return {"..product.."}")
			
			alreadyAdjusting = false
		end)
	end
end

function makeWood()
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
			
			local items = {}
			local found = {}
			for i = 1, 3, 1 do
				block = rubble.powered_items.FindItemAtInput(wshop, function(item)
					if df.item_type[item:getType()] == "BLOCKS" then
						if found[item.id] == true then
							return false
						end
						
						local mat = dfhack.matinfo.decode(item)
						if mat.mode == "plant" then
							return true
						end
					end
					return false
				end)
				if block == nil then
					return
				end
				found[block.id] = true
				items[i] = block
			end
			
			local mat = nil
			for _, block in ipairs(items) do
				mat = dfhack.matinfo.decode(block)
				dfhack.items.remove(block)
			end
			
			item = rubble.powered_items.CreateItemNumeric(mat, output.itype, output.isubtype, nil, 0)
			rubble.powered_items.SetAutoItemQuality(wshop, item)
			rubble.powered_items.Eject(wshop, item)
		end
	end
end

buildings.registerBuilding{
	name="DFHACK_RUBBLE_POWERED_CARPENTER",
	consume=30,
	gears={{x=0,y=2},{x=2,y=0},{x=2,y=2},{x=0,y=0}},
	action={300, makeWood()},
	animate={
		isMechanical=true,
		frames={
			{{x=0,y=2,42,0,0,1},{x=2,y=0,42,0,0,1},{x=2,y=2,42,0,0,1},{x=0,y=0,42,0,0,1}},
			{{x=0,y=2,15,0,0,1},{x=2,y=0,15,0,0,1},{x=2,y=2,15,0,0,1},{x=0,y=0,15,0,0,1}}
		}
	}
}
eventful.registerReaction("LUA_HOOK_ADJUST_CARPENTER", carpenterAdjust)
print("    Registered mechanical workshop: DFHACK_RUBBLE_POWERED_CARPENTER")
