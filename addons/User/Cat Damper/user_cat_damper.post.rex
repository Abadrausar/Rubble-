
var foundcat = false

(console:print "    Stripping tags ADOPTS_OWNER and RETURNS_VERMIN_KILLS_TO_OWNER from CAT\n")
(foreach [rubble:raws] block name content {
	(if (str:cmp "creature_" (str:left [name] 9)) {
		(console:print "      " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "CREATURE") {
				(if (int:eq (len [tag]) 1) {
					(if (str:cmp [tag 0] "CAT"){
						[foundcat = true]
					}{
						[foundcat = false]
					})
				}{
					(rubble:abort "Error: invalid param count to CREATURE raw tag in last file.")
				})
			})
			
			(if (bool:and (str:cmp [tag id] "ADOPTS_OWNER") [foundcat]) {
				[tag disable = true]
			})
			
			(if (bool:and (str:cmp [tag id] "RETURNS_VERMIN_KILLS_TO_OWNER") [foundcat]) {
				[tag disable = true]
			})
			
			(break true)
		})]
	})
	(break true)
})
