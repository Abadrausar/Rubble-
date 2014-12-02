
var ores = <map
	IRON=true
	NICKEL=true
	GOLD=true
	#SILVER=true # silver is a special case
	COPPER=true
	LEAD=true
	ZINC=true
	TIN=true
	PLATINUM=true
	BISMUTH=true
	ALUMINUM=true
>

(foreach [rubble:raws] block name content {
	(if (str:cmp "inorganic_" (str:left [name] 10)) {
		(console:print "    " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "METAL_ORE") {
				(if [ores [tag 0]] {
					[tag replace = (str:add "[REACTION_CLASS:" [tag 0] "_ORE]")]
				})
			})
			
			(if (str:cmp [tag id] "METAL_ORE") {
				(if (str:cmp [tag 0] SILVER)  {
					(if (int:lt [tag 1] 100) {
						[tag replace = (str:add "[REACTION_CLASS:" [tag 0] "_ORE_POOR]")]
					}{
						[tag replace = (str:add "[REACTION_CLASS:" [tag 0] "_ORE]")]
					})
					
					
				})
			})
			
			(break true)
		})]
	})
	(break true)
})
