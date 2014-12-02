-- Do stuff with the earth

local eventful = require 'plugins.eventful'
local utils = require 'utils'

function block_events(pos, dir)
	local block = dfhack.maps.ensureTileBlock(pos.x,pos.y,pos.z)
	local found = false
	local uniqueFinds = {}
	if block and #block.block_events > 0 then
		for index,event in ipairs(block.block_events) do
			eventType = event:getType()
			if eventType == 0 then -- mineral
				if event.flags.discovered == false then -- only print if you haven't found it yet
					found = true
					local mineralName = ""
					inorganic = df.inorganic_raw.find(event.inorganic_mat)
					if inorganic.material.gem_name1 ~= "" then
						mineralName = inorganic.material.gem_name1
					else
						mineralName = inorganic.material.state_name.Solid
					end
					line = ""
					if uniqueFinds[mineralName] == nil then
						if dir == "here" then
							if event.flags.cluster == true then
								dfhack.gui.showAnnouncement("This kind of rock often bears " .. mineralName .. "." , COLOR_WHITE, true)
							elseif event.flags.vein == true then
								dfhack.gui.showAnnouncement("There are signs of " .. mineralName .. " nearby." , COLOR_WHITE, true)
							elseif event.flags.cluster_small == true then
								dfhack.gui.showAnnouncement("There are hints of " .. mineralName .. " nearby." , COLOR_WHITE, true)
							elseif event.flags.cluster_one == true then
								dfhack.gui.showAnnouncement("There may be " .. mineralName .. " nearby." , COLOR_WHITE, true)
							else
								dfhack.gui.showAnnouncement("There is " .. mineralName .. " nearby." , COLOR_WHITE, true)
							end
						else
							
							if event.flags.cluster == true then
								dfhack.gui.showAnnouncement("This stone to the "..dir.." is likely to have " .. mineralName .. "." , COLOR_WHITE, true)
							elseif event.flags.vein == true then
								dfhack.gui.showAnnouncement("There are signs of " .. mineralName .. " in the stone to the " .. dir , COLOR_WHITE, true)
							elseif event.flags.cluster_small == true then
								dfhack.gui.showAnnouncement("There are hints of " .. mineralName .. " in the stone to the " .. dir , COLOR_WHITE, true)
							elseif event.flags.cluster_one == true then
								dfhack.gui.showAnnouncement("There may be " .. mineralName .. " in the stone to the " .. dir , COLOR_WHITE, true)
							else
								dfhack.gui.showAnnouncement("There is " .. mineralName .. " nearby." , COLOR_WHITE, true)
							end
						end
						uniqueFinds[mineralName] = 1
					end
				end
			end
		end
	end
	return found
end

function prospect(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value = false
	
	dfhack.gui.showAnnouncement(dfhack.TranslateName(unit.name).." examined the surrounding area.", COLOR_WHITE, true)
    local pos = {}
	pos.x = unit.pos.x
	pos.y = unit.pos.y
	pos.z = unit.pos.z
	found = false
	if block_events(pos, "here") == true then found = true end
	pos.y = pos.y-16
	if block_events(pos, "north") == true then found = true end
	pos.x = pos.x+16
	if block_events(pos, "northeast") == true then found = true end
	pos.y = pos.y+16
	if block_events(pos, "east") == true then found = true end
	pos.y = pos.y+16
	if block_events(pos, "southeast") == true then found = true end
	pos.x = pos.x-16
	if block_events(pos, "south") == true then found = true end
	pos.x = pos.x-16
	if block_events(pos, "southwest") == true then found = true end
	pos.y = pos.y-16
	if block_events(pos, "west") == true then found = true end
	pos.y = pos.y-16
	if block_events(pos, "northwest") == true then found = true end
	if found == false then
		dfhack.gui.showAnnouncement("Nothing unexpected was found.", COLOR_WHITE, true)
	end
end

eventful.registerReaction("LUA_HOOK_PROSPECT_GEOLOGY", prospect)
