
(console:print "    Fixing Creatures...\n")

(foreach [rubble:raws] block name content {
	(if (str:cmp "creature_" (str:left [name] 9)) {
		(console:print "      " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "TENDONS") {
				(if (int:eq (len [tag]) 3) {
					[tag replace = "[TENDONS:CREATURE_MAT:ANIMAL:SINEW]"]
				}{
					(console:print "      TENDONS tag in creature has non-standard arg count\n")
					(console:print "        This is probably OK so ignoring.\n")
				})
			})
			
			(if (str:cmp [tag id] "LIGAMENTS") {
				(if (int:eq (len [tag]) 3) {
					[tag replace = "[LIGAMENTS:CREATURE_MAT:ANIMAL:SINEW]"]
				}{
					(console:print "      LIGAMENTS tag in creature has non-standard arg count\n")
					(console:print "        This is probably OK so ignoring.\n")
				})
			})
			
			(break true)
		})]
	})
	(break true)
})
