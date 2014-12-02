
var curEntity = ""

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
		(if (str:cmp [tag "id"] "PERMITTED_REACTION") {
			(if (int:eq (len [tag]) 1) {
				(if (str:cmp [tag 0] "MAKE_CLAY_JUG") {
					[tag disable = true]
				})
				(if (str:cmp [tag 0] "MAKE_CLAY_BRICKS") {
					[tag disable = true]
				})
				(if (str:cmp [tag 0] "MAKE_CLAY_STATUE") {
					[tag disable = true]
				})
				(if (str:cmp [tag 0] "MAKE_LARGE_CLAY_POT") {
					[tag disable = true]
				})
				(if (str:cmp [tag 0] "MAKE_CLAY_CRAFTS") {
					[tag disable = true]
				})
				(if (str:cmp [tag 0] "MAKE_CLAY_HIVE") {
					[tag disable = true]
				})
				(if (str:cmp [tag 0] "GLAZE_JUG") {
					[tag disable = true]
				})
				(if (str:cmp [tag 0] "GLAZE_STATUE") {
					[tag disable = true]
				})
				(if (str:cmp [tag 0] "GLAZE_LARGE_POT") {
					[tag disable = true]
				})
				(if (str:cmp [tag 0] "GLAZE_CRAFT") {
					[tag disable = true]
				})
			}{
				(rubble:abort "Error: invalid param count to PERMITTED_REACTION raw tag in entity_saurian.txt")
			})
		})
	})
	
	(break true)
})]
