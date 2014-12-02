
# Generic programmable reaction zapper script.

var curEntity = ""

var entities = <smap
	MOUNTAIN = true
	SWAMP = true
>

var strip = <smap
	BRASS_MAKING=true
	BRASS_MAKING2=true
	BRONZE_MAKING=true
	BRONZE_MAKING2=true
	ELECTRUM_MAKING=true
	ELECTRUM_MAKING2=true
	BILLON_MAKING=true
	BILLON_MAKING2=true
	PEWTER_FINE_MAKING=true
	PEWTER_FINE_MAKING2=true
	PEWTER_TRIFLE_MAKING=true
	PEWTER_TRIFLE_MAKING2=true
	PEWTER_LAY_MAKING=true
	PIG_IRON_MAKING=true
	STEEL_MAKING=true
	NICKEL_SILVER_MAKING=true
	BLACK_BRONZE_MAKING=true
	STERLING_SILVER_MAKING=true
	ROSE_GOLD_MAKING=true
	BISMUTH_BRONZE_MAKING=true
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
