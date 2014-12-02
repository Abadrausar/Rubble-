
local eventful = require "plugins.eventful"
local buildings = require 'plugins.building-hacks'
local script = require 'gui.script'
local powered = rubble.require "powered"
local pitems = rubble.require "powered_items"
local ppersist = rubble.require "powered_persist"

function getItemCaption(type, subtype)
    local attrs = df.item_type.attrs[type]
    if subtype == -1 or subtype == nil then
		return attrs.caption
	else
		return df['itemdef_'..string.lower(df.item_type[type])..'st'].find(subtype).name
	end
end

alreadyAdjusting = false
function masonAdjust(reaction, unit, in_items, in_reag, out_items, call_native)
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
			local choiceok, choice = script.showListPrompt('Powered Mason', 'Select item to produce:', COLOR_LIGHTGREEN, opt_names)
			
			local product = ""
			if choiceok then
				product = opt_types[choice]
			else
				product = "NONE"
			end
			
			local wshop = powered.MakeFake(unit.pos.x, unit.pos.y, unit.pos.z, 1)
			ppersist.SetOutputType(wshop, "return {"..product.."}")
			
			alreadyAdjusting = false
		end)
	end
end

function makeStone()
	return function(wshop)
		if not wshop:isUnpowered() then
			if not powered.HasOutput(wshop) then
				return
			end
			
			local output = ppersist.GetOutputTypeAsCode(wshop)
			if output == nil then
				return
			end
			
			local blocks = pitems.FindXItemsAtInput(wshop, 3, function(item)
				if df.item_type[item:getType()] == "BLOCKS" then
					if found[item.id] == true then
						return false
					end
					
					local mat = dfhack.matinfo.decode(item)
					if mat.mode == "inorganic" then
						return true
					end
				end
				return false
			end)
			if blocks == nil then
				return
			end
			
			local mat = nil
			for _, block in ipairs(blocks) do
				mat = dfhack.matinfo.decode(block)
				dfhack.items.remove(block)
			end
			
			local item = pitems.CreateItemNumeric(mat, output.itype, output.isubtype, nil, 0)
			pitems.SetAutoItemQuality(wshop, item)
			pitems.Eject(wshop, item)
		end
	end
end

buildings.registerBuilding{
	name="DFHACK_RUBBLE_POWERED_MASON",
	consume=30,
	gears={{x=0,y=2},{x=2,y=0},{x=2,y=2},{x=0,y=0}},
	action={500, makeStone()},
	animate={
		isMechanical=true,
		frames={
			{{x=0,y=2,42,0,0,1},{x=2,y=0,42,0,0,1},{x=2,y=2,42,0,0,1},{x=0,y=0,42,0,0,1}},
			{{x=0,y=2,15,0,0,1},{x=2,y=0,15,0,0,1},{x=2,y=2,15,0,0,1},{x=0,y=0,15,0,0,1}}
		}
	}
}
eventful.registerReaction("LUA_HOOK_ADJUST_MASON", masonAdjust)
print("    Registered mechanical workshop: DFHACK_RUBBLE_POWERED_MASON")
