(console:print "    Fixing Creatures...\n")

(foreach [rubble:raws] block name content {
	(if (str:cmp "creature_" (str:left [name] 9)) {
		(console:print "      " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "MILKABLE") {
				(if (int:eq (len [tag]) 3) {
					[tag replace = (str:add "[MILKABLE:CREATURE_MAT:ANIMAL:MILK:" [tag 2] "]")]
				}{
					(console:print "      MILKABLE tag in creature has non-standard arg count\n")
					(console:print "        This is probably OK so ignoring.\n")
				})
			})
			
			(break true)
		})]
	})
	(break true)
})

(console:print "    Fixing Material Template...\n")

var found = false
[rubble:raws "material_template_default.txt" = (df:raw:walk [rubble:raws "material_template_default.txt"] block tag {
	(if (str:cmp [tag id] "MATERIAL_TEMPLATE") {
		(if (int:eq (len [tag]) 1) {
			(if (str:cmp [tag 0] "MILK_TEMPLATE"){
				[found = true]
			}{
				[found = false]
			})
		}{
			(rubble:abort "Error: invalid param count to TISSUE_TEMPLATE raw tag in last file.")
		})
	})
	
	(if (bool:and (str:cmp [tag id] "MATERIAL_REACTION_PRODUCT") [found]) {
		(if (int:gt (len [tag]) 1) {
			(if (str:cmp [tag 0] "CHEESE_MAT") {
				[tag replace = "[MATERIAL_REACTION_PRODUCT:CHEESE_MAT:CREATURE_MAT:ANIMAL:CHEESE]"]
			})
		})
	})
	(break true)
})]
