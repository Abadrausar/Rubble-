
(console:print "    Making bones more powerful vs. their source...\n")
(console:print "    (for creatures that use BODY_DETAIL_PLAN:STANDARD_MATERIALS)\n")

var found = false
(foreach [rubble:raws] block name content {
	(if (str:cmp "creature_" (str:left [name] 9)) {
		(console:print "      " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "BODY_DETAIL_PLAN") {
				(if (str:cmp [tag 0] "STANDARD_MATERIALS") {
					[tag append = "\n\t[MATERIAL_FORCE_MULTIPLIER:LOCAL_CREATURE_MAT:BONE:10:1]\n"]
				})
			})
			
			(break true)
		})]
	})
	(break true)
})
