
# Item templates

var rubble:item_current_id = ""
var rubble:item_current_type = ""
var rubble:item_data = <map>
var rubble:item_ban_data = <map>

# Set to true if item type allows rarity
var rubble:item_types = <map
	"AMMO"=false
	"ARMOR"=true
	"DIGGER"=false
	"GLOVES"=true
	"HELM"=true
	"INSTRUMENT"=false
	"PANTS"=true
	"SHIELD"=false
	"SHOES"=true
	"SIEGEAMMO"=false
	"TOOL"=false
	"TOY"=false
	"TRAPCOMP"=false
	"WEAPON"=false
>

var rubble:item_rarities = <map
	"RARE"=1
	"UNCOMMON"=2
	"COMMON"=3
	"FORCED"=4
>

(rubble:template "SHARED_ITEM" block type id def {
	(if (str:cmp [type] "FOOD") {
		# FOOD doesn't need to be registered in an entity, so translate it directly to a SHARED_OBJECT call.
		(ret (rubble:stageparse (str:add "{SHARED_OBJECT;" [id] ";\n[ITEM_" [type] ":" [id] "]\n\t" [def] "\n}")))
	})
	
	[rubble:item_current_id = [id]]
	[rubble:item_current_type = [type]]
	
	(if (exists [rubble:item_types] [type]){
	}{
		(rubble:abort "Error: Invalid item type passed to SHARED_ITEM.")
	})
	
	[def = (regex:replace "\\[TILE:([0-9]+|'.')\\]" [def] (str:add "[TILE:{#TILE;" [id] ";$1}]"))]
	
	(if (str:cmp [type] "DIGGER") {
		[type = "WEAPON"]
	})
	
	(rubble:stageparse (str:add "{SHARED_OBJECT;" [id] ";\n[ITEM_" [type] ":" [id] "]\n\t" [def] "\n}"))
})
(rubble:template "ITEM_CLASS" block ... {
	var type = [rubble:item_current_type]
	var id = [rubble:item_current_id]
	var class = ""
	var rarity = "COMMON"
	
	(if (int:eq (len [params]) 1){
		# class only
		[class = [params 0]]
	}{
		(if (int:eq (len [params]) 2){
			# class + rarity
			[class = [params 0]]
			[rarity = [params 1]]
		}{
			(if (int:eq (len [params]) 3){
				# all but rarity
				[type = [params 0]]
				[id = [params 1]]
				[class = [params 2]]
			}{
				(if (int:eq (len [params]) 4){
					# full call
					[type = [params 0]]
					[id = [params 1]]
					[class = [params 2]]
					[rarity = [params 3]]
				}{
					(rubble:abort "Error: Invalid param count to ITEM_CLASS.")
				})
			})
			
		})
	})
	
	(if (exists [rubble:item_rarities] [rarity]){
	}{
		(rubble:abort "Error: Invalid item rarity passed to ITEM_CLASS.")
	})
	
	(if (exists [rubble:item_data] [class]){
		(if (exists [rubble:item_data [class]] [id]) {
		}{
			[rubble:item_data [class] [id] = <map
				rarity=[rarity]
				type=[type]
			>]
		})
	}{
		[rubble:item_data [class] = <map
			[id]=<map
				rarity=[rarity]
				type=[type]
			>
		>]
	})
	(ret "")
})
(rubble:template "REMOVE_ITEM" block id class {
	(if (exists [rubble:item_ban_data] [class]){
		[rubble:item_ban_data [class] [id] = true]
	}{
		[rubble:item_ban_data [class] = <map [id]=true>]
	})
	(ret "")
})
(rubble:template "REMOVE_ITEM_FROM_PLAYABLES" block id {
	(foreach [rubble:entity_playability] block ent flags {
		(if [flags fort] {}{
			(breakloop true)
		})
		
		var class = (str:add "ADDON_HOOK_" [ent])
		(if (exists [rubble:item_ban_data] [class]){
			[rubble:item_ban_data [class] [id] = true]
		}{
			[rubble:item_ban_data [class] = <map [id]=true>]
		})
		
		(break true)
	})
	
	(if (exists [rubble:item_ban_data] ADDON_HOOK_PLAYABLE){
		[rubble:item_ban_data ADDON_HOOK_PLAYABLE [id] = true]
	}{
		[rubble:item_ban_data ADDON_HOOK_PLAYABLE = <map [id]=true>]
	})
	(ret "")
})
command rubble:uses_items class {
	var out = ""
	(if (exists [rubble:item_data] [class]){
		(foreach [rubble:item_data [class]] block id _ {
			(if (exists [rubble:item_ban_data] [class]) {
				(if [rubble:item_ban_data [class] [id]] {
					(breakloop true)
				})
			})
			
			var type = [rubble:item_data [class] [id] type]
			var rarity = [rubble:item_data [class] [id] rarity]
			(if [rubble:item_types [type]] {
				# Has rarity
				[out = (str:add [out] "\n\t[" [type] ":" [id] ":" [rarity] "]")]
			}{
				# Does not have rarity
				[out = (str:add [out] "\n\t[" [type] ":" [id] "]")]
			})
			(break true)
		})
	})
	(str:trimspace [out])
}
(rubble:template "#USES_ITEMS" [rubble:uses_items])
