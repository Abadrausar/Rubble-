
# Item templates

var rubble:item_current_id = ""
var rubble:item_current_type = ""
var rubble:item_data = <map>

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
		(if (exists [rubble:item_data [class]] [rubble:item_current_id]) {
		}{
			[rubble:item_data [class] [rubble:item_current_id] = <map
				"rarity"=[rarity]
				"type"=[rubble:item_current_type]
			>]
		})
	}{
		[rubble:item_data [class] = <map
			[rubble:item_current_id]=<map
				"rarity"=[rarity]
				"type"=[rubble:item_current_type]
			>
		>]
	})
	(ret "")
})
(rubble:template "#USES_ITEMS" block class {
	var out = ""
	(if (exists [rubble:item_data] [class]){
		(foreach [rubble:item_data [class]] block key value {
			var type = [rubble:item_data [class] [key] "type"]
			var rarity = [rubble:item_data [class] [key] "rarity"]
			(if [rubble:item_types [type]] {
				# Has rarity
				[out = (str:add [out] "\n\t[" [type] ":" [key] ":" [rarity] "]")]
			}{
				# Does not have rarity
				[out = (str:add [out] "\n\t[" [type] ":" [key] "]")]
			})
			(break true)
		})
	})
	(str:trimspace [out])
})
