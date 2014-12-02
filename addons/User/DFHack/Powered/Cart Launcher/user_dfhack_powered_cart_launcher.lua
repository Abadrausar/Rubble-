
local eventful = require "plugins.eventful"
local buildings = require 'plugins.building-hacks'
local script = require 'gui.script'
local powered = rubble.require "powered"
local pitems = rubble.require "powered_items"
local ppersist = rubble.require "powered_persist"

alreadyAdjusting = false
function launcherAdjust(reaction, unit, in_items, in_reag, out_items, call_native)
	call_native.value = false
	
	if alreadyAdjusting == false then
		alreadyAdjusting = true
		
		script.start(function()
			local thold = "50"
			repeat
				_, thold = script.showInputPrompt('Minecart Launcher Adjust', 'Set how full the minecart needs to be (0-100)', COLOR_LIGHTGREEN, thold)
			until tonumber(thold) and tonumber(thold) >= 0 and tonumber(thold) <= 100
			
			local dirok, dir
			repeat
				dirok, dir = script.showListPrompt('Minecart Launcher Adjust', 'Select launch direction:', COLOR_LIGHTGREEN, {"N", "S", "E", "W"})
			until dirok
			
			if thold ~= "0" then
				thold = tostring(tonumber(thold) / 100)
			end
			
			local launcher_dirs = {
				"vx = 0, vy = -20000", -- N
				"vx = 0, vy = 20000",  -- S
				"vx = 20000, vy = 0",  -- E
				"vx = -20000, vy = 0", -- W
			}
			
			local wshop = powered.MakeFake(unit.pos.x, unit.pos.y, unit.pos.z, 1)
			ppersist.SetOutputType(wshop, "return {"..launcher_dirs[dir]..", threshold = "..thold.."}")
			
			alreadyAdjusting = false
		end)
	end
end

function makeLaunchCart()
	return function(wshop)
		if not wshop:isUnpowered() then
			-- Read settings
			local output = ppersist.GetOutputTypeAsCode(wshop)
			if output == nil then
				return
			end
			
			-- find cart
			local cart = pitems.FindItemArea(wshop, function(item)
				if item:isTrackCart() then
					local cart_capacity = item.subtype.container_capacity
					local totalvolume = 0
					local cartrefs = item.general_refs
					for r = 0, #cartrefs - 1, 1 do
						if getmetatable(cartrefs[r])=="general_ref_contains_itemst" then
							totalvolume = totalvolume + df.item.find(cartrefs[r].item_id):getVolume()
						end
					end
					if totalvolume >= (cart_capacity * output.threshold) then
						return true
					end
				end
				return false
			end)
			if not cart then
				-- No available minecart
				return
			end
			
			-- launch the cart
			local vehicle = df.global.world.vehicles.all[cart.vehicle_id]
			
			-- It may be possible to see if vehicle.time_stopped > 0 instead...
			if vehicle.speed_x == 0 and vehicle.speed_y == 0 and vehicle.speed_z == 0 then
				-- Launch!
				vehicle.speed_x = output.vx
				vehicle.speed_y = output.vy
			end
		end
	end
end

buildings.registerBuilding{
	name="DFHACK_RUBBLE_POWERED_CART_LAUNCHER",
	consume=5,
	gears={{x=0,y=0}},
	action={50, makeLaunchCart()},
	animate={
		isMechanical=true,
		frames={
			{{x=0,y=0,42,7,0,0}},
			{{x=0,y=0,15,7,0,0}}
		}
	}
}
eventful.registerReaction("LUA_HOOK_ADJUST_CART_LAUNCHER", launcherAdjust)
print("    Registered mechanical workshop: DFHACK_RUBBLE_POWERED_CART_LAUNCHER")
