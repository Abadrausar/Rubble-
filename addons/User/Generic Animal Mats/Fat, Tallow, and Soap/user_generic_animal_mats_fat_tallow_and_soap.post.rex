
(console:print "    Fixing Tissues...\n")

var mats = <map
	FAT_TEMPLATE = true
	TALLOW_TEMPLATE = true
>

var found = false
(foreach [rubble:raws] block name content {
	(if (str:cmp "tissue_template_" (str:left [name] 16)) {
		(console:print "      " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "TISSUE_TEMPLATE") {
				(if (int:eq (len [tag]) 1) {
					(if [mats [tag 0]] {
						[found = true]
					}{
						[found = false]
					})
				}{
					(rubble:abort "Error: invalid param count to TISSUE_TEMPLATE raw tag in last file.")
				})
			})
			
			(if (bool:and (str:cmp [tag id] "TISSUE_MATERIAL") [found]) {
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

(console:print "    Fixing Material Template...\n")

[found = false]
[rubble:raws "material_template_default.txt" = (df:raw:walk [rubble:raws "material_template_default.txt"] block tag {
	(if (str:cmp [tag id] "MATERIAL_TEMPLATE") {
		(if (int:eq (len [tag]) 1) {
			(if [mats [tag 0]] {
				[found = true]
			}{
				[found = false]
			})
		}{
			(rubble:abort "Error: invalid param count to TISSUE_TEMPLATE raw tag in last file.")
		})
	})
	
	(if (bool:and (str:cmp [tag id] "MATERIAL_REACTION_PRODUCT") [found]) {
		(if (int:eq (len [tag]) 3) {
			(if (str:cmp [tag 0] "RENDER_MAT") {
				[tag replace = (str:add "[MATERIAL_REACTION_PRODUCT:RENDER_MAT:CREATURE_MAT:ANIMAL:" [tag 1] "]")]
			})
			
			(if (str:cmp [tag 0] "SOAP_MAT") {
				[tag replace = (str:add "[MATERIAL_REACTION_PRODUCT:SOAP_MAT:CREATURE_MAT:ANIMAL:" [tag 1] "]")]
			})
		})
	})
	(break true)
})]
