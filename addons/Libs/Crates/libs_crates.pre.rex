
module rubble:libs_crates
var rubble:libs_crates:data = (sort:new)
var rubble:libs_crates:classes = (sort:new)
var rubble:libs_crates:last = ""

# Define a new crate.
(rubble:template "!CRATE" block id name value product {
	(if (exists [rubble:libs_crates:data] [id]){
		(rubble:abort "Error: Crate already exists.")
	}{
		[rubble:libs_crates:last = [id]]
		[rubble:libs_crates:data [id] = <map
			"name"=[name]
			"value"=[value]
			"product"=[product]
			"bar"=false
		>]
	})
	(ret "")
})

# Define a new crate (containing 10 bars).
(rubble:template "!CRATE_BARS" block id name value product {
	(if (exists [rubble:libs_crates:data] [id]){
		(rubble:abort "Error: Crate already exists.")
	}{
		[rubble:libs_crates:last = [id]]
		[rubble:libs_crates:data [id] = <map
			"name"=(str:add [name] " bars (10)")
			"value"=[value]
			"product"=[product]
			"bar"=true
		>]
	})
	(ret "")
})

# Add last defined crate to class.
# Crates may have more than one class.
# !CRATE_CLASS class | !CRATE_CLASS id class
(rubble:template "!CRATE_CLASS" block params=... {
	var id = [rubble:libs_crates:last]
	var class = ""
	
	(if (int:eq (len [params]) 1){
		# class only
		[class = [params 0]]
	}{
		(if (int:eq (len [params]) 2){
			# id + class
			[id = [params 0]]
			[class = [params 1]]
		}{
			(rubble:abort "Error: Invalid param count to !CRATE_CLASS.")
		})
	})
	
	(if (exists [rubble:libs_crates:classes] [class]){
		[rubble:libs_crates:classes [class] "append" = [id]]
	}{
		[rubble:libs_crates:classes [class] = <array [id]>]
	})
	
	(ret "")
})

# This generates a list of product lines for all crates, use in world gen reactions.
(rubble:template "CRATE_WORLDGEN_REACTION_PRODUCTS" {
	var out = ""
	(foreach [rubble:libs_crates:data] block id _ {
		[out = (str:add [out]
			"\t[PRODUCT:100:1:BAR:NONE:CREATURE_MAT:LIBS_CRATES_CREATURE:" [id] "][PRODUCT_DIMENSION:150]\n"
		)]
		(break true)
	})
	(str:trimspace [out])
})

# This generates a list of product lines for a crate class, use in world gen reactions.
(rubble:template "CRATE_WORLDGEN_REACTION_PRODUCTS_CLASSED" block class {
	(if (exists [rubble:libs_crates:classes] [class]){
		var out = ""
		(foreach [rubble:libs_crates:classes [class]] block _ id {
			[out = (str:add [out] "\n\t"
				"[PRODUCT:100:1:BAR:NONE:CREATURE_MAT:LIBS_CRATES_CREATURE:" [id] "][PRODUCT_DIMENSION:150]"
			)]
			(break true)
		})
		(str:trimspace [out])
	}{
		(ret "")
	})
})

# Generate unpack reactions for all crates.
# Example:
# {CRATE_UNPACK_REACTIONS;CRAFTSMAN;CARPENTRY;ADDON_HOOK_PLAYABLE}
(rubble:template "CRATE_UNPACK_REACTIONS" block building skill techclass auto=true {
	var out = ""
	(foreach [rubble:libs_crates:data] block id data {
		(if [data "bar"] {
			[out = (str:add [out] "\n"
				"{REACTION;UNPACK_" [id] ";" [techclass] "}\n"
				"\t[NAME:unpack " [data "name"] "]\n"
				"\t[BUILDING:" [building] ":NONE]\n"
				"\t[REAGENT:A:150:BAR:NONE:NONE:NONE:" [id] "]\n"
				"\t\t[HAS_MATERIAL_REACTION_PRODUCT:CRATE_" [id] "_MAT]"
				"\t[SKILL:" [skill] "]\n"
				"\t[PRODUCT:100:10:BAR:NONE:GET_MATERIAL_FROM_REAGENT:A:CRATE_" [id] "_MAT]\n"
				(if [auto] {"\t[AUTOMATIC]\n"}{""})
			)]
			(breakloop true)
		})
		
		[out = (str:add [out] "\n"
			"{REACTION;UNPACK_" [id] ";" [techclass] "}\n"
			"\t[NAME:unpack " [data "name"] "]\n"
			"\t[BUILDING:" [building] ":NONE]\n"
			"\t[REAGENT:A:150:BAR:NONE:CREATURE_MAT:LIBS_CRATES_CREATURE:" [id] "]\n"
			"\t[SKILL:" [skill] "]\n"
			"\t" [data "product"] "\n"
			(if [auto] {"\t[AUTOMATIC]\n"}{""})
		)]
		(break true)
	})
	(str:trimspace (rubble:stageparse [out]))
})

# Same as CRATE_UNPACK_REACTIONS, except for only a single crate class.
(rubble:template "CRATE_UNPACK_REACTIONS_CLASSED" block building skill techclass crateclass auto=true {
	(if (exists [rubble:libs_crates:classes] [crateclass]){
		var out = "";
		(foreach [rubble:libs_crates:classes [crateclass]] block _ id {
			(if [rubble:libs_crates:data [id] "bar"] {
				[out = (str:add [out] "\n"
					"{REACTION;UNPACK_" [id] ";" [techclass] "}\n"
					"\t[NAME:unpack " [rubble:libs_crates:data [id] "name"] "]\n"
					"\t[BUILDING:" [building] ":NONE]\n"
					"\t[REAGENT:A:150:BAR:NONE:NONE:NONE:" [id] "]\n"
					"\t\t[HAS_MATERIAL_REACTION_PRODUCT:CRATE_" [id] "_MAT]"
					"\t[SKILL:" [skill] "]\n"
					"\t[PRODUCT:100:10:BAR:NONE:GET_MATERIAL_FROM_REAGENT:A:CRATE_" [id] "_MAT]\n"
					(if [auto] {"\t[AUTOMATIC]\n"}{""})
				)]
				(breakloop true)
			})
			
			[out = (str:add [out] "\n"
				"{REACTION;UNPACK_" [id] ";" [techclass] "}\n"
				"\t[NAME:unpack " [rubble:libs_crates:data [id] "name"] "]\n"
				"\t[BUILDING:" [building] ":NONE]\n"
				"\t[REAGENT:A:150:BAR:NONE:CREATURE_MAT:LIBS_CRATES_CREATURE:" [id] "]\n"
				"\t[SKILL:" [skill] "]\n"
				"\t" [rubble:libs_crates:data [id] "product"] "\n"
				(if [auto] {"\t[AUTOMATIC]\n"}{""})
			)]
			(break true)
		})
		(str:trimspace (rubble:stageparse [out]))
	}{
		(ret "")
	})
})

(rubble:template "#CRATE_MATS" {
	var out = ""
	(foreach [rubble:libs_crates:data] block id data {
		[out = (str:add [out] "\n\n"
			"\t[USE_MATERIAL_TEMPLATE:" [id] ":CRATE_TEMPLATE]\n"
			"\t\t[STATE_NAME_ADJ:ALL_SOLID:" [data "name"] " crate]\n"
			"\t\t[STATE_NAME_ADJ:LIQUID:melted " [data "name"] " crate]\n"
			"\t\t[STATE_NAME_ADJ:GAS:boiling " [data "name"] " crate]\n"
			"\t\t[PREFIX:NONE]\n"
			"\t\t[MATERIAL_VALUE:" [data "value"] "]"
		)]
		
		(if [data "bar"] {
			[out = (str:add [out] "\n\t\t[MATERIAL_REACTION_PRODUCT:CRATE_" [id] "_MAT:" [data "product"] "]")]
		})
		
		(break true)
	})
	(str:trimspace [out])
})
