
# Generic programmable reaction zapper script.

var curEntity = ""

var entities = <smap
	MOUNTAIN = true
	SWAMP = true
>

var strip = <smap
	MAKE_CLAY_JUG=true
	MAKE_CLAY_BRICKS=true
	MAKE_CLAY_STATUE=true
	MAKE_LARGE_CLAY_POT=true
	MAKE_CLAY_CRAFTS=true
	MAKE_CLAY_HIVE=true
	GLAZE_JUG=true
	GLAZE_STATUE=true
	GLAZE_LARGE_POT=true
	GLAZE_CRAFT=true
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
