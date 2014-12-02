
(console:print "    Fixing Creatures...\n")

(foreach [rubble:raws] block name content {
	(if (str:cmp "creature_" (str:left [name] 9)) {
		(console:print "      " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "BLOOD") {
				(if (int:eq (len [tag]) 3) {
					[tag replace = (str:add "[BLOOD:CREATURE_MAT:ANIMAL:" [tag 1] ":LIQUID]")]
				}{
					(console:print "      BLOOD tag in creature has non-standard arg count\n")
					(console:print "        This is probably OK so ignoring.\n")
				})
			})
			
			(if (str:cmp [tag id] "PUS") {
				(if (int:eq (len [tag]) 3) {
					[tag replace = (str:add "[PUS:CREATURE_MAT:ANIMAL:" [tag 1] ":LIQUID]")]
				}{
					(console:print "      PUS tag in creature has non-standard arg count\n")
					(console:print "        This is probably OK so ignoring.\n")
				})
			})
			
			(if (str:cmp [tag id] "TENDONS") {
				(if (int:eq (len [tag]) 3) {
					[tag replace = "[TENDONS:CREATURE_MAT:ANIMAL:SINEW]"]
				}{
					(console:print "      TENDONS tag in creature has non-standard arg count\n")
					(console:print "        This is probably OK so ignoring.\n")
				})
			})
			
			(if (str:cmp [tag id] "LIGAMENTS") {
				(if (int:eq (len [tag]) 3) {
					[tag replace = "[LIGAMENTS:CREATURE_MAT:ANIMAL:SINEW]"]
				}{
					(console:print "      LIGAMENTS tag in creature has non-standard arg count\n")
					(console:print "        This is probably OK so ignoring.\n")
				})
			})
			
			(if (str:cmp [tag id] "MILKABLE") {
				(if (int:eq (len [tag]) 3) {
					[tag replace = "[MILKABLE:CREATURE_MAT:ANIMAL:MILK]"]
				}{
					(console:print "      MILKABLE tag in creature has non-standard arg count\n")
					(console:print "        This is probably OK so ignoring.\n")
				})
			})
			
			(if (str:cmp [tag id] "EGG_MATERIAL") {
				(if (int:eq (len [tag]) 3) {
					[tag replace = (str:add "[EGG_MATERIAL:CREATURE_MAT:ANIMAL:" [tag 1] ":" [tag 2] "]")]
				}{
					(console:print "      EGG_MATERIAL tag in creature has non-standard arg count\n")
					(console:print "        This is probably OK so ignoring.\n")
				})
			})
			
			(break true)
		})]
	})
	(break true)
})

(console:print "    Fixing Tissues...\n")

(foreach [rubble:raws] block name content {
	(if (str:cmp "tissue_template_" (str:left [name] 16)) {
		(console:print "      " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "TISSUE_MATERIAL") {
				(if (int:gt (len [tag]) 1) {
					(if (str:cmp [tag 0] "LOCAL_CREATURE_MAT") {
						(if(int:eq (len [tag]) 3){
							[tag replace = (str:add "[TISSUE_MATERIAL:CREATURE_MAT:ANIMAL:" [tag 1] ":" [tag 2] "]")]
						})
						(if (int:eq (len [tag]) 2) {
							[tag replace = (str:add "[TISSUE_MATERIAL:CREATURE_MAT:ANIMAL:" [tag 1] "]")]
						})
					})
				})
			})
			(break true)
		})]
	})
	(break true)
})

(console:print "    Fixing Material Templates...\n")

[rubble:raws "material_template_default.txt" = (str:replace [rubble:raws "material_template_default.txt"] "LOCAL_CREATURE_MAT" "CREATURE_MAT:ANIMAL" -1)]
