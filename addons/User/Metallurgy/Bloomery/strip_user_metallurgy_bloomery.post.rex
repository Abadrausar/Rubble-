
# Generic programmable reaction zapper script.

var curEntity = ""

var entities = <smap
	MOUNTAIN = true
	SWAMP = true
>

var strip = <smap
	ALLOY_PIG_IRON=true
	ALLOY_STEEL=true
	PIG_IRON_MAKING=true
	STEEL_MAKING=true
	SMELT_ORE_IRON=true
>

(foreach [rubble:raws] block name content {
	(if (str:cmp "entity_" (str:left [name] 7)) {
		(console:print "    " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag "id"] "ENTITY") {
				[curEntity = [tag 0]]
			})
			
			(if [entities [curEntity]] {
				(if (str:cmp [tag id] "PERMITTED_REACTION") {
					(if [strip [tag 0]] {
						[tag disable = true]
					})
				})
			})
			
			(break true)
		})]
	})
})
