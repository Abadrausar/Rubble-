
# Tech Templates
var rubble:building_last = ""
var rubble:building_data = <map>
var rubble:building_ban_data = <map>
(rubble:template "BUILDING_WORKSHOP" block id class {
	[rubble:building_last = [id]]
	(if (exists [rubble:building_data] [class]){
		[rubble:building_data [class] append = [id]]
	}{
		[rubble:building_data [class] = <array [id]>]
	})
	(str:add "[BUILDING_WORKSHOP:" [id] "]")
})
(rubble:template "BUILDING_FURNACE" block id class {
	[rubble:building_last = [id]]
	(if (exists [rubble:building_data] [class]){
		[rubble:building_data [class] append = [id]]
	}{
		[rubble:building_data [class] = <array [id]>]
	})
	(str:add "[BUILDING_FURNACE:" [id] "]")
})
(rubble:template "BUILDING_ADD_CLASS" block class id=nil {
	(if (isnil [id]) {
		[id = [rubble:building_last]]
	})
	(if (exists [rubble:building_data] [class]){
		[rubble:building_data [class] append = [id]]
	}{
		[rubble:building_data [class] = <array [id]>]
	})
	(ret "")
})
(rubble:template "REMOVE_BUILDING" block id class {
	(if (exists [rubble:building_ban_data] [class]){
		[rubble:building_ban_data [class] [id] = true]
	}{
		[rubble:building_ban_data [class] = <map [id]=true>]
	})
	(ret "")
})
(rubble:template "REMOVE_BUILDING_FROM_PLAYABLES" block id {
	(foreach [rubble:entity_playability] block ent flags {
		(if [flags fort] {}{
			(breakloop true)
		})
		
		var class = (str:add "ADDON_HOOK_" [ent])
		(if (exists [rubble:building_ban_data] [class]){
			[rubble:building_ban_data [class] [id] = true]
		}{
			[rubble:building_ban_data [class] = <map [id]=true>]
		})
		
		(break true)
	})
	
	(if (exists [rubble:building_ban_data] ADDON_HOOK_PLAYABLE){
		[rubble:building_ban_data ADDON_HOOK_PLAYABLE [id] = true]
	}{
		[rubble:building_ban_data ADDON_HOOK_PLAYABLE = <map [id]=true>]
	})
	(ret "")
})
command rubble:uses_buildings params {
	(if (int:lt (len [params]) 1){(rubble:abort "Error: Invalid param count to rubble:uses_buildings.")})
	
	var buildings = <map>
	var buildingnames = <array>
	(foreach [params] block key class {
		(if (exists [rubble:building_data] [class]) {
			(foreach [rubble:building_data [class]] block key building {
				(if (exists [rubble:building_ban_data] [class]) {
					(if [rubble:building_ban_data [class] [building]] {
						(breakloop true)
					})
				})
				
				(if [buildings [building]]  {
					# nop
				}{
					[buildings [building] = true]
					[buildingnames append = [building]]
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
	(str:trimspace [out])
}
(rubble:template "#USES_BUILDINGS" block params=... {(rubble:uses_buildings [params])})

var rubble:reaction_data = <map>
var rubble:reaction_ban_data = <map>
var rubble:reaction_last = ""
command rubble:reaction id class {
	[rubble:reaction_last = [id]]
	(if (exists [rubble:reaction_data] [class]){
		[rubble:reaction_data [class] append = [id]]
	}{
		[rubble:reaction_data [class] = <array [id]>]
	})
	(str:add "[REACTION:" [id] "]")
}
(rubble:template "REACTION" [rubble:reaction])
(rubble:template "DFHACK_REACTION" block id action class {
	[rubble:dfhack_reactions [id] = [action]]
	(rubble:reaction [id] [class])
})
(rubble:template "DFHACK_REACTION_BIND" block action id=nil {
	(if (str:cmp [id] "") {
		(if (isnil [rubble:reaction_data [rubble:reaction_last]]) {
			(rubble:abort "Error: rubble:reaction_last is invalid in call to DFHACK_REACTION_BIND.")
		})
		
		[rubble:dfhack_reactions [rubble:reaction_last] = [action]]
	}{
		[rubble:dfhack_reactions [id] = [action]]
	})
	(ret "")
})
(rubble:template "REACTION_ADD_CLASS" block class id=nil {
	(if (isnil [id]) {
		[id = [rubble:reaction_last]]
	})
	(if (exists [rubble:reaction_data] [class]){
		[rubble:reaction_data [class] append = [id]]
	}{
		[rubble:reaction_data [class] = <array [id]>]
	})
	(ret "")
})
(rubble:template "REMOVE_REACTION" block id class {
	(if (exists [rubble:reaction_ban_data] [class]){
		[rubble:reaction_ban_data [class] [id] = true]
	}{
		[rubble:reaction_ban_data [class] = <map [id]=true>]
	})
	(ret "")
})
(rubble:template "REMOVE_REACTION_FROM_PLAYABLES" block id {
	(foreach [rubble:entity_playability] block ent flags {
		(if [flags fort] {}{
			(breakloop true)
		})
		
		var class = (str:add "ADDON_HOOK_" [ent])
		(if (exists [rubble:reaction_ban_data] [class]){
			[rubble:reaction_ban_data [class] [id] = true]
		}{
			[rubble:reaction_ban_data [class] = <map [id]=true>]
		})
		
		(break true)
	})
	
	(if (exists [rubble:reaction_ban_data] ADDON_HOOK_PLAYABLE){
		[rubble:reaction_ban_data ADDON_HOOK_PLAYABLE [id] = true]
	}{
		[rubble:reaction_ban_data ADDON_HOOK_PLAYABLE = <map [id]=true>]
	})
	(ret "")
})
command rubble:uses_reactions params {
	(if (int:lt (len [params]) 1){(rubble:abort "Error: Invalid param count to rubble:uses_reactions.")})
	
	var reactions = <map>
	var reactionnames = <array>
	(foreach [params] block key class {
		(if (exists [rubble:reaction_data] [class]) {
			(foreach [rubble:reaction_data [class]] block key reaction {
				(if (exists [rubble:reaction_ban_data] [class]) {
					(if [rubble:reaction_ban_data [class] [reaction]] {
						(breakloop true)
					})
				})
				
				(if [reactions [reaction]] {
					# nop
				}{
					[reactions [reaction] = true]
					[reactionnames append = [reaction]]
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
	(str:trimspace [out])
}
(rubble:template "#USES_REACTIONS" block params=... {(rubble:uses_reactions [params])})
(rubble:template "#USES_TECH" block params=... {
	(if (int:lt (len [params]) 1){(rubble:abort "Error: Invalid param count to #USES_TECH.")})
	
	[out = (str:add [out] (rubble:uses_buildings [params]) "\n")]
	[out = (str:add [out] (rubble:uses_reactions [params]))]
	(ret (str:trimspace [out]))
})

# combination of #USES_TECH and #USES_ITEMS, for internal use.
(rubble:template "#_ADDON_HOOK" block class {
	var out = (str:add "# Hook: " [class] "\n")
	
	var items = (rubble:uses_items [class])
	(if (str:cmp [items] "") {}{
		[out = (str:add [out] "\t" [items] "\n")]
	})
	
	var buildings = (rubble:uses_buildings <array [class]>)
	(if (str:cmp [buildings] "") {}{
		[out = (str:add [out] "\t" [buildings] "\n")]
	})
	
	var reactions = (rubble:uses_reactions <array [class]>)
	(if (str:cmp [reactions] "") {}{
		[out = (str:add [out] "\t" [reactions] "\n")]
	})
	
	(str:trimspace [out])
})

(rubble:template "ADDON_HOOKS" block id {
	(str:add
		"{#_ADDON_HOOK;ADDON_HOOK_" [id] "}\n"
		"\t{#_ADDON_HOOK;ADDON_HOOK_GENERIC}\n"
		"\t{#ECHO;{@IF_ENTITY_PLAYABLE;" [id] ";FORT;\n"
		"\t{#_ADDON_HOOK;ADDON_HOOK_PLAYABLE}\n"
		"\t;\n"
		"\t# Entity Not Playable in Fortress Mode.\n"
		"\t}}"
	)
})
