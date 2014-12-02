
--[[
Rubble Persistent Timeout DFHack Lua Pseudo Module

Copyright 2014 Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
]]

local _ENV = rubble.mkmodule("timeout")

data = {}

function savePersist()
	print("rubble.timeout: Saving persistence data.")
	local dat = io.open(rubble.savedir..'/raw/rubble_libs_dfhack_timeout.persist.lua', 'w')
	dat:write("\n")
	dat:write("rubble.timeout.data = {\n")
	for k, v in pairs(data) do
		if v ~= nil then
			dat:write('\t["'..k..'"] = { delay = '..v.delay..', command = [['..v.command..']] },\n')
		end
	end
	dat:write("}\n")
	dat:close()
end

print("rubble.timeout: Loading persistence data.")
rubble.load_script(rubble.savedir.."/raw/rubble_libs_dfhack_timeout.persist.lua")

function add(id, delay, command)
	data[id] = {delay = delay, command = command}
	savePersist()
end

function del(id)
	data[id] = nil
	savePersist()
end

function tick()
	for k, v in pairs(data) do
		if v ~= nil then
			v.delay = v.delay - 1
			if v.delay <= 0 then
				dfhack.run_command(v.command)
				del(k)
			end
		end
	end
	
	dfhack.timeout(1, 'ticks', tick)
end
tick()

return _ENV
