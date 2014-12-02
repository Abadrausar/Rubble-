-- Do fun stuff with water and magma.

-- [[  ]] is a really stupid choice for string delimiters, particularly when you need to do a usage statement.
-- BTW: experience has shown that it is trivial to support multi-line double quote strings in a lexer,
-- (it actually takes more work to disable them) so why do most languages limit them to a single line?
local usage = [[
Usage: 
	fluids [-h|/?]
	]].."rubble_fluids [eat] [spot|3x3|5x5|7x7] [magma|water] [amount [minimum [x y z]]]\n"..
"	rubble_fluids spawn [magma|water] [amount [x y z]]\n"..[[
	rubble_fluids cart [spot|3x3|5x5|7x7] [magma|water] [x y z]

Allows you to:
	Spawn fluids in a single tile.
	Eat fluids, either from a single tile or from any tile below a downward passable tile in an area.
	Fill a minecart with fluids, either in a single tile or in an area.

Examples:
	spawn 4/7 magma at cursor
		rubble_fluids spawn magma 4
	eat 2/7 magma at cursor, but only if there is at least 4/7 there
		rubble_fluids magma 2 4
	fill the minecart at cursor with water
		rubble_fluids cart water
	fill the minecart at coords <x:10, y:20, z:30> with water
		rubble_fluids cart water 10 20 30
	eat 2/7 magma from the first tile below a 3x3 area around the cursor that is 
	downward-passable and has at least 2/7 magma
		rubble_fluids eat 3x3 magma 2

This script is mostly a test driver for the functions in the rubble_fluids module.

If you use Rubble there are some convenience templates avalible in the
"Libs/DFHack/Fluids" addon that make using this script in reactions MUCH easier.
(That addon also installs this script and the required lua module)
]]

local pos=df.global.cursor
local args={...}

if args[1]=="-h" or args[1]=="/?" then
	print(usage)
	return
end

local fluids = require "rubble_fluids"

-- parse options, oh for go's "flag" package...
local spawn = false
local area = false
local offset = 1
local tocart = false
local magma = false
local narg = 1
local amount = 1
local minimum = 1

local x = -1
local y = -1
local z = -1

if pos.x ~= -30000 then
	x = pos.x
	y = pos.y
	z = pos.z
end

if args[narg]=="spawn" then
	spawn = true
	narg = narg + 1
elseif args[narg]=="eat" then
	narg = narg + 1
elseif args[narg]=="cart" then
	tocart = true
	narg = narg + 1
end

if not spawn or tocart then
	if args[narg]=="3x3" then
		area = true
		offset = 1
		narg = narg + 1
	elseif args[narg]=="5x5" then
		area = true
		offset = 2
		narg = narg + 1
	elseif args[narg]=="7x7" then
		area = true
		offset = 3
		narg = narg + 1
	end
end

if args[narg]=="magma" then
	magma = true
	narg = narg + 1
elseif args[narg]=="water" then
	narg = narg + 1
end

if not tocart then
	if args[narg] ~= nil then
		amount = tonumber(args[narg])
		narg = narg + 1
	end
end

if not spawn then
	if args[narg] ~= nil then
		minimum = tonumber(args[narg])
		narg = narg + 1
	end
end

if args[narg] ~= nil then
	x = tonumber(args[narg])
	narg = narg + 1
end

if args[narg] ~= nil then
	y = tonumber(args[narg])
	narg = narg + 1
end

if args[narg] ~= nil then
	z = tonumber(args[narg])
	narg = narg + 1
end

if x == -1 or y == -1 or z == -1 then
	print("Error: invalid position.")
	print("Drop a cursor or specify coords on the command line.")
	print("Type \"rubble_fluids -h\" for help.")
	return
end

if minimum < amount then
	print("Warning: minimum less than amount or not specified, adjusting.")
	minimum = amount
end

-- Do the requested operation.
-- Error handling is a little sparse.
if spawn then
	fluids.spawnFluid(x, y, z, magma, amount)
	return
elseif tocart then
	local cart
	if area then
		cart = fluids.findCartArea(x-offset, y-offset, x+offset, y+offset, z, magma)
	else
		cart = fluids.findCart(x, y, z, magma)
	end
	
	if cart == nil then
		print("Error: No minecart at coords/in area.")
		return
	end
	
	fluids.fillCart(cart, magma)
else
	if area then
		fluids.eatFromArea(x-offset, y-offset, x+offset, y+offset, z, magma, amount, minimum)
		return
	end
	
	fluids.eatFluid(x, y, z, magma, amount, minimum)
	return
end
