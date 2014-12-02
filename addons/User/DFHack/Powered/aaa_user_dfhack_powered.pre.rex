
(rubble:dfhack_loadscript "aaa_user_dfhack_powered.mod.lua")
(rubble:dfhack_loadscript "bbb_user_dfhack_powered.mod.lua")
(rubble:dfhack_loadscript "ccc_user_dfhack_powered.mod.lua")
(rubble:dfhack_loadscript "ddd_user_dfhack_powered.mod.lua")

module rubble:dfhack_powered

# for use with struct
# <struct proto=[rubble:dfhack_powered:output]>
var rubble:dfhack_powered:output = <map
	id=""
	name=""
>

# for use with struct
# <struct proto=[rubble:dfhack_powered:wshop]>
var rubble:dfhack_powered:wshop = <map
	# The player visible building name, this MUST be specified.
	name=nil
	
	# Array of rubble:dfhack_powered:output structs.
	# Leave unset for no output morphs.
	outputs=<sarray>
	
	# The addon hook to use for reactions and workshops.
	hook="ADDON_HOOK_PLAYABLE"
	
	# Tile and color of the center workshop tile.
	centertile="42"
	centercolor="7:0:0"
	
	# Any build items you may want to add beyond the standard.
	extraitems=""
	
	# If true makes a 1x1 workshop rather than a 3x3.
	# If this is true centertile and centercolor are ignored (set from the Lua side).
	small=false
>

var rubble:dfhack_powered:wshops = <map>

# To register a new workshop:
# [rubble:dfhack_powered:wshops [id] = <struct proto=[rubble:dfhack_powered:wshop] name="Powered Something">]

(rubble:template "!POWERED_WSHOP_REACTIONS" block id {
	var out = ""
	
	var wshop = [rubble:dfhack_powered:wshops [id]]
	(if (isnil [wshop]) {
		(rubble:abort (str:add "Attempt to generate reactions for non-existent powered workshop: " [id]))
	})
	
	(if (int:gt (len [wshop outputs]) 0) {
		(foreach [wshop outputs] block _ meoutput {
				var me = (str:add [id] "_" [meoutput id])
				[out = (str:add [out] `\n{REACTION;` [me] `_OUTPUT;` [wshop hook] `}\n\t[BUILDING:` [me] `:NONE]\n\t[NAME:producing ` [meoutput name]`]\n\t[SKILL:MECHANICS]\n`)]
				[out = (str:add [out] `\n{REACTION;` [me] `_SEPERATOR;` [wshop hook] `}\n\t[BUILDING:` [me] `:NONE]\n\t[NAME:-----------------]\n\t[SKILL:MECHANICS]\n`)]
				
				(foreach [wshop outputs] block _ output {
					(breakloopif (str:cmp [meoutput id] [output id]) true)
					
					var it = (str:add [id] "_" [output id])
					[out = (str:add [out] `\n{DFHACK_REACTION_UPGRADE_BUILDING;` [me] `;` [it] `;output ` [output name]`;` [wshop hook] `}\n\t[SKILL:MECHANICS]\n`)]
					
					(break true)
				})
			(break true)
		})
	}{
		(rubble:abort (str:add "Attempt to generate reactions for workshop without changeable outputs: " [id]))
	})
	
	(str:trimspace [out])
})

(rubble:template "!POWERED_WSHOP_BUILDINGS" block id {
	var out = ""
	
	var wshop = [rubble:dfhack_powered:wshops [id]]
	(if (isnil [wshop]) {
		(rubble:abort (str:add "Attempt to generate buildings for non-existent powered workshop: " [id]))
	})
	(if (isnil [wshop name]) {
		(rubble:abort (str:add "Attempt to generate buildings for powered workshop with no name: " [id]))
	})
	
	var generate_shop = block output {
		# Let's see you do this in another language :)
		var hook = (if (isnil [output]) {
			[wshop hook]
		}{
			(if (str:cmp [output id] [wshop outputs 0 id]) {
				[wshop hook]
			}{
				"NIL"
			})
		})
		
		(if [wshop small] {
			[out = (str:add [out] `
{BUILDING_WORKSHOP;` [id] (if (isnil [output]) {``}{(str:add `_` [output id])}) `;` [hook] `}
	[NAME:` [wshop name] `]
	[NAME_COLOR:7:0:1]
	[BUILD_LABOR:MECHANIC]
	[DIM:1:1]
	[WORK_LOCATION:1:1]
	[BLOCK:1:0]
	[TILE:0:1:42]
	[COLOR:0:1:0:7:0]
	[TILE:1:1:128]
	[COLOR:1:1:0:7:0]
	[TILE:2:1:15]
	[COLOR:2:1:0:7:0]
	[TILE:3:1:42]
	[COLOR:3:1:0:7:0]
	[BUILD_ITEM:1:BLOCKS:NONE:NONE:NONE]
		[BUILDMAT]
	[BUILD_ITEM:1:TRAPPARTS:NONE:NONE:NONE]
		[CAN_USE_ARTIFACT]
` (if (isnil [wshop extraitems]) {``}{[wshop extraitems]})
			)]
		}{
			[out = (str:add [out] `
{BUILDING_WORKSHOP;` [id] (if (isnil [output]) {``}{(str:add `_` [output id])}) `;` [hook] `}
	[NAME:` [wshop name] `]
	[NAME_COLOR:7:0:1]
	[BUILD_LABOR:MECHANIC]
	[DIM:3:3]
	[WORK_LOCATION:2:2]
	[BLOCK:1:0:0:0]
	[BLOCK:2:0:0:0]
	[BLOCK:3:0:0:0]
	[TILE:0:1:32:254:32]
	[TILE:0:2:128:32:128]
	[TILE:0:3:42:32:128]
	[COLOR:0:1:0:0:0:0:0:1:0:0:0]
	[COLOR:0:2:0:0:1:0:0:0:0:0:1]
	[COLOR:0:3:7:0:0:0:0:0:0:0:1]
	[TILE:1:1:32:32:32]
	[TILE:1:2:32:32:32]
	[TILE:1:3:32:32:32]
	[COLOR:1:1:0:0:0:0:0:0:0:0:0]
	[COLOR:1:2:0:0:0:0:0:0:0:0:0]
	[COLOR:1:3:0:0:0:0:0:0:0:0:0]
	[TILE:2:1:254:32:42]
	[TILE:2:2:32:42:17]
	[TILE:2:3:128:32:42]
	[COLOR:2:1:0:0:1:0:0:0:0:0:1]
	[COLOR:2:2:0:0:0:7:0:0:0:0:1]
	[COLOR:2:3:0:0:1:0:0:0:0:0:1]
	[TILE:3:1:42:31:42]
	[TILE:3:2:16:` [wshop centertile] `:17]
	[TILE:3:3:42:30:42]
	[COLOR:3:1:0:0:1:0:0:1:0:0:1]
	[COLOR:3:2:0:0:1:` [wshop centercolor] `:0:0:1]
	[COLOR:3:3:0:0:1:0:0:1:0:0:1]
	
	[BUILD_ITEM:4:BLOCKS:NONE:NONE:NONE]
		[BUILDMAT]
	[BUILD_ITEM:2:TRAPPARTS:NONE:NONE:NONE]
		[CAN_USE_ARTIFACT]
` (if (isnil [wshop extraitems]) {``}{[wshop extraitems]})
			)]
		})
	}
	
	(if (int:gt (len [wshop outputs]) 0) {
		# outputs only
		(foreach [wshop outputs] block _ output {
			([generate_shop] [output])
			(break true)
		})
	}{
		([generate_shop] nil)
	})
	
	(str:trimspace [out])
})
