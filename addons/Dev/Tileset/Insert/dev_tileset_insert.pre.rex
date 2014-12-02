
# Rewrite any and all objects that have a specialized shared object template to use
# that template (if they don't already).
# Also inserts calls to #TILE and #COLOR in vermin creatures.

# Use only with untemplated raws!

module rubble:dev_tileset_insert

var export = false
(if (rubble:configvar "DEV_TILESET_INSERT_EXPORT") {
	[export = true]
})

command rubble:dev_tileset_insert:guts name contents prefix tags tagprefix="" {
	var foundfirst = false
	var prev
	var file = (df:raw:parse [contents])
	(foreach [file] block i tag {
		(foreach [tags] block _ type {
			(if (str:cmp [tag id] (str:add [tagprefix] [type])) {
				(if (int:eq (len [tag params]) 1) {
					(if (bool:and (exists [file] (int:sub [i] 1)) [foundfirst]) {
						[file (int:sub [i] 1) append = "\n}"]
					})
					[foundfirst = true]
					[tag replace = (str:add "{" [prefix] [type] ";" [tag params 0] ";")]
				}{
					(rubble:abort (str:add "Error: invalid param count to " [type] " raw tag in last file."))
				})
			})
		})
	})
	
	(if [foundfirst] {
		[file (int:sub (len [file]) 1) append = "\n}"]
	})
	
	var out = (df:raw:dump [file])
	[rubble:raws [name] = [out]]
	
	(if [export] {
		(axis:write [rubble:fs] (str:add "rubble:dev_tileset_insert/" [name]) [out])
	})
}

var itemtypes = <sarray
	"AMMO"
	"ARMOR"
	"DIGGER"
	"GLOVES"
	"HELM"
	"INSTRUMENT"
	"PANTS"
	"SHIELD"
	"SHOES"
	"SIEGEAMMO"
	"TOOL"
	"TOY"
	"TRAPCOMP"
	"WEAPON"
>

(foreach [rubble:raws] block name contents {
	(if (str:cmp "inorganic_" (str:left [name] 10)) {
		(rubble:dev_tileset_insert:guts [name] [contents] "SHARED_" <sarray "INORGANIC">)
	}{
		(if (str:cmp "plant_" (str:left [name] 6)) {
			(rubble:dev_tileset_insert:guts [name] [contents] "SHARED_" <sarray "PLANT">)
		}{
			(if (str:cmp "material_template_" (str:left [name] 18)) {
				(rubble:dev_tileset_insert:guts [name] [contents] "SHARED_" <sarray "MATERIAL_TEMPLATE">)
			}{
				(if (str:cmp "item_" (str:left [name] 5)) {
					(rubble:dev_tileset_insert:guts [name] [contents] "SHARED_ITEM;" [itemtypes] "ITEM_")
				})
			})
		})
	})
	(break true)
})

# Insert templates in creatures.

var curCreature = ""
var changedFile = false
(foreach [rubble:raws] block name contents {
	(if (str:cmp "creature_" (str:left [name] 9)) {
		(console:print "      " [name] "\n")
		
		var newfile = (df:raw:walk [contents] block tag {
			(if (str:cmp [tag id] "CREATURE") {
				[curCreature = [tag 0]]
			})
			
			# I think the following two tags are the only ones that matter for vermin.
			# [COLOR:0:0:0]
			(if (str:cmp [tag id] "COLOR") {
				(if (int:eq (len [tag]) 3) {
					[tag replace = (str:add "[COLOR:{#COLOR;" [curCreature] ";;" [tag 0] ":" [tag 1] ":" [tag 2] "}]")]
					[changedFile = true]
				}{
					(rubble:abort "Error: invalid param count to COLOR raw tag in last file.")
				})
			})
			
			# [CREATURE_TILE:'A']
			(if (str:cmp [tag id] "CREATURE_TILE") {
				(if (int:eq (len [tag]) 1) {
					[tag replace = (str:add "[CREATURE_TILE:{#TILE;" [curCreature] ";;" [tag 0] "}]")]
					[changedFile = true]
				}{
					(rubble:abort "Error: invalid param count to CREATURE_TILE raw tag in last file.")
				})
			})
			
			(break true)
		})
		
		(if [changedFile] {
			[rubble:raws [name] = [newfile]]
			
			(if [export] {
				(axis:write [rubble:fs] (str:add "rubble:dev_tileset_insert/" [name]) [newfile])
			})
		})
		[changedFile = false]
	})
	(break true)
})
