
var ores = <map
	IRON=true
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
			
			(break true)
		})]
	})
	(break true)
})
