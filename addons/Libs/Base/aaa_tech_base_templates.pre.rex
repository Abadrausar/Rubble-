
# Tech Templates
var rubble:building_data = <map>
(rubble:template "BUILDING_WORKSHOP" block ... {
	(if (int:lt (len [params]) 2){(rubble:abort "Error: Invalid param count to BUILDING_WORKSHOP.")})
	
	var id = [params 0]
	(for 0 (int:sub (len [params]) 1) 1 block count {
		var class = [params [count]]
		(if (exists [rubble:building_data] [class]){
			[rubble:building_data [class] "append" = [id]]
		}{
			[rubble:building_data [class] = <array [id]>]
		})
		(break true)
	})
	(str:add "[BUILDING_WORKSHOP:" [id] "]")
})
(rubble:template "BUILDING_FURNACE" block ... {
	(if (int:lt (len [params]) 2){(rubble:abort "Error: Invalid param count to BUILDING_FURNACE.")})
	
	var id = [params 0]
	(for 0 (int:sub (len [params]) 1) 1 block count {
		var class = [params [count]]
		(if (exists [rubble:building_data] [class]){
			[rubble:building_data [class] "append" = [id]]
		}{
			[rubble:building_data [class] = <array [id]>]
		})
		(break true)
	})
	(str:add "[BUILDING_FURNACE:" [id] "]")
})
(rubble:template "#USES_BUILDINGS" block ... {
	(if (int:lt (len [params]) 1){(rubble:abort "Error: Invalid param count to #USES_BUILDINGS.")})
	
	var buildings = <map>
	var buildingnames = <array>
	(foreach [params] block key class {
		(if (exists [rubble:building_data] [class]) {
			(foreach [rubble:building_data [class]] block key building {
				(if (exists [buildings] [building]) {
					# nop
				}{
					[buildings [building] = "ok"]
					[buildingnames "append" = [building]]
				})
				(break true)
			})
		})
		(break true)
	})
	
	var out = ""
	(foreach [buildingnames] block key building {
		[out = (str:add [out] "\n\t[PERMITTED_BUILDING:" [building] "]")]
		(break true)
	})
	(ret (str:trimspace [out]))
})

var rubble:reaction_data = <map>
(rubble:template "REACTION" block ... {
	(if (int:lt (len [params]) 2){(rubble:abort "Error: Invalid param count to REACTION.")})
	
	var id = [params 0]
	(for 1 (int:sub (len [params]) 1) 1 block count {
		var class = [params [count]]
		(if (exists [rubble:reaction_data] [class]){
			[rubble:reaction_data [class] "append" = [id]]
		}{
			[rubble:reaction_data [class] = <array [id]>]
		})
		(break true)
	})
	(str:add "[REACTION:" [id] "]")
})
(rubble:template "DFHACK_REACTION" block ... {
	(if (int:lt (len [params]) 3){(rubble:abort "Error: Invalid param count to DFHACK_REACTION.")})
	
	var id = [params 0]
	var action = [params 1]
	(for 2 (int:sub (len [params]) 1) 1 block count {
		var class = [params [count]]
		(if (exists [rubble:reaction_data] [class]){
			[rubble:reaction_data [class] "append" = [id]]
		}{
			[rubble:reaction_data [class] = <array [id]>]
		})
		(break true)
	})
	
	[rubble:dfhack_reactions [id] = [action]]
	(str:add "[REACTION:" [id] "]")
})
(rubble:template "#USES_REACTIONS" block ... {
	(if (int:lt (len [params]) 1){(rubble:abort "Error: Invalid param count to #USES_REACTIONS.")})
	
	var reactions = <map>
	var reactionnames = <array>
	(foreach [params] block key class {
		(if (exists [rubble:reaction_data] [class]) {
			(foreach [rubble:reaction_data [class]] block key reaction {
				(if (exists [reactions] [reaction]) {
					# nop
				}{
					[reactions [reaction] = "ok"]
					[reactionnames "append" = [reaction]]
				})
				(break true)
			})
		})
		(break true)
	})
	
	var out = ""
	(foreach [reactionnames] block key reaction {
		[out = (str:add [out] "\n\t[PERMITTED_REACTION:" [reaction] "]")]
		(break true)
	})
	(ret (str:trimspace [out]))
})

(rubble:template "#USES_TECH" block ... {
	(if (int:lt (len [params]) 1){(rubble:abort "Error: Invalid param count to #USES_TECH.")})
	
	var buildings = <map>
	var buildingnames = <array>
	var reactions = <map>
	var reactionnames = <array>
	(foreach [params] block key class {
		(if (exists [rubble:building_data] [class]) {
			(foreach [rubble:building_data [class]] block key building {
				(if (exists [buildings] [building]) {
					# nop
				}{
					[buildings [building] = "ok"]
					[buildingnames "append" = [building]]
				})
				(break true)
			})
		})
		(if (exists [rubble:reaction_data] [class]) {
			(foreach [rubble:reaction_data [class]] block key reaction {
				(if (exists [reactions] [reaction]) {
					# nop
				}{
					[reactions [reaction] = "ok"]
					[reactionnames "append" = [reaction]]
				})
				(break true)
			})
		})
		(break true)
	})
	
	var out = ""
	(foreach [buildingnames] block key building {
		[out = (str:add [out] "\n\t[PERMITTED_BUILDING:" [building] "]")]
		(break true)
	})
	(foreach [reactionnames] block key reaction {
		[out = (str:add [out] "\n\t[PERMITTED_REACTION:" [reaction] "]")]
		(break true)
	})
	(ret (str:trimspace [out]))
})

# combination of #USES_TECH and #USES_ITEMS, for internal use.
(rubble:template "#_ADDON_HOOK" block class {
	var out = (str:add "# Hook: " [class])
	
	# Items
	(if (exists [rubble:item_data] [class]){
		(foreach [rubble:item_data [class]] block item value {
			var type = [rubble:item_data [class] [item] "type"]
			var rarity = [rubble:item_data [class] [item] "rarity"]
			(if [rubble:item_types [type]] {
				# Has rarity
				[out = (str:add [out] "\n\t[" [type] ":" [item] ":" [rarity] "]")]
			}{
				# Does not have rarity
				[out = (str:add [out] "\n\t[" [type] ":" [item] "]")]
			})
			(break true)
		})
	})
	
	var buildings = <map>
	var buildingnames = <array>
	var reactions = <map>
	var reactionnames = <array>
	
	(if (exists [rubble:building_data] [class]) {
		(foreach [rubble:building_data [class]] block key building {
			(if (exists [buildings] [building]) {
				# nop
			}{
				[buildings [building] = "ok"]
				[buildingnames "append" = [building]]
			})
			(break true)
		})
	})
	(if (exists [rubble:reaction_data] [class]) {
		(foreach [rubble:reaction_data [class]] block key reaction {
			(if (exists [reactions] [reaction]) {
				# nop
			}{
				[reactions [reaction] = "ok"]
				[reactionnames "append" = [reaction]]
			})
			(break true)
		})
	})
	
	(foreach [buildingnames] block key building {
		[out = (str:add [out] "\n\t[PERMITTED_BUILDING:" [building] "]")]
		(break true)
	})
	(foreach [reactionnames] block key reaction {
		[out = (str:add [out] "\n\t[PERMITTED_REACTION:" [reaction] "]")]
		(break true)
	})
	
	(str:trimspace [out])
})

(rubble:template "ADDON_HOOKS" block id playable=false {
	(str:add
		"{#_ADDON_HOOK;ADDON_HOOK_" [id] "}\n"
		"\t{#_ADDON_HOOK;ADDON_HOOK_GENERIC}"
		(if [playable] {"\n\t{#_ADDON_HOOK;ADDON_HOOK_PLAYABLE}"}{""})
	)
})
