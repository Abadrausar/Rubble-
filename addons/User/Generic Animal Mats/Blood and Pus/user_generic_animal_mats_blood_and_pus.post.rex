
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
			
			(break true)
		})]
	})
	(break true)
})
