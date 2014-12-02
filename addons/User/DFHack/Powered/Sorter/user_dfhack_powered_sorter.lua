
local eventful = require "plugins.eventful"
local buildings = require 'plugins.building-hacks'
local script = require('gui.script')

alreadyAdjusting = false
function sorterAdjust(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value = false
	
	if alreadyAdjusting == false then
		alreadyAdjusting = true
		
		script.start(function()
			local itemok = false
			local itemtype, itemsubtype
			local matok = false
			local mattype, matindex
			
			local adjust = script.showYesNoPrompt('Sorter Adjust', 'Sort by specific item type?', COLOR_LIGHTGREEN)
			if adjust == true then
				require('gui.materials').ItemTypeDialog{
					text = "Sort what item type?",
					item_filter = function() return true end,
					hide_none = true,
					on_select = script.mkresume(true),
					on_cancel = script.mkresume(false),
					on_close = script.qresume(nil)
				}:show()
				
				itemok, itemtype, itemsubtype = script.wait()
			end
				
			adjust = script.showYesNoPrompt('Sorter Adjust', 'Sort by specific material?', COLOR_LIGHTGREEN)
			if adjust == true then
				matok, mattype, matindex = script.showMaterialPrompt('Sort','Sort what material?')
			end
			
			local ipart = ""
			local mpart = ""
			if itemok then
				ipart = "itype = "..itemtype..", isubtype = "..itemsubtype
			else
				ipart = "itype = nil, isubtype = -1"
			end
			if matok then
				mpart = "mtype = "..mattyp..", mindex = "..matindex
			else
				mpart = "mtype = nil, mindex = -1"
			end
			
			local wshop = rubble.powered.MakeFake(unit.pos.x, unit.pos.y, unit.pos.z, 1)
			rubble.powered_persist.SetOutputType(wshop, "return {"..ipart..", "..mpart.."}")
			
			alreadyAdjusting = false
		end)
	end
end

function makeSortItem()
	return function(wshop)
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
		
		if not rubble.powered.HasOutput(wshop) then
			return
		end
		
		local matchitem = function(item)
			match = true
			if output.itype ~= nil then
				if item:getType() ~= output.itype then
					match = false
				end
				if output.isubtype ~= -1 then
					if item:getSubtype() ~= output.isubtype then
						match = false
					end
				end
			end
			
			if output.mtype ~= nil then
				local mat = dfhack.matinfo.decode(item)
				
				if mat.type ~= output.mtype or mat.index ~= mindex then
					match = false
				end
			end
			
			return match
		end
		
		local item = nil
		if rubble.powered.HasInput(wshop) then
			item = rubble.powered_items.FindItemAtInput(wshop, matchitem)
			if item == nil then
				return
			end
		else
			item = rubble.powered_items.FindItemArea(wshop, matchitem)
			if item == nil then
				return
			end
		end
		
		rubble.powered_items.Eject(wshop, item)
	end
end

buildings.registerBuilding{
	name="DFHACK_RUBBLE_POWERED_SORTER",
	action={10, makeSortItem()},
}
eventful.registerReaction("LUA_HOOK_ADJUST_SORTER", sorterAdjust)
print("    Registered mechanical workshop: DFHACK_RUBBLE_POWERED_SORTER")
