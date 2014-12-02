
var curEntity = ""
var curPosition = ""

(if (exists [rubble:raws] "entity_saurian.txt"){}{(ret "")})
[rubble:raws "entity_saurian.txt" = (df:raw:walk [rubble:raws "entity_saurian.txt"] block tag {
	(if (str:cmp [tag "id"] "ENTITY") {
		(if (int:eq (len [tag]) 1) {
			[curEntity = [tag 0]]
		}{
			(rubble:abort "Error: invalid param count to ENTITY raw tag in entity_saurian.txt")
		})
	})
	
	(if (str:cmp [curEntity] "SWAMP") {
		(if (str:cmp [tag "id"] "CIV_CONTROLLABLE") {
			[tag disable = true]
		})
		
		(if (str:cmp [tag "id"] "POSITION") {
			(if (int:eq (len [tag]) 1) {
				[curPosition = [tag 0]]
			}{
				(rubble:abort "Error: invalid param count to POSITION raw tag in entity entity_saurian.txt")
			})
		})
		
		(if (str:cmp [curPosition] "TRADE_REP") {
			(if (str:cmp [tag "id"] "RESPONSIBILITY") {
				(if (int:eq (len [tag]) 1) {
					(if (str:cmp [tag 0] "ESTABLISH_COLONY_TRADE_AGREEMENTS") {
						[tag replace = "[RESPONSIBILITY:TRADE]"]
					})
				}{
					(rubble:abort "Error: invalid param count to RESPONSIBILITY raw tag in entity entity_saurian.txt")
				})
			})
		})
	})
	
	(break true)
})]
