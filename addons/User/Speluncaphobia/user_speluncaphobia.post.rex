
var creatures = (eval (rubble:configvar "USER_SPELUNCAPHOBIA_CREATURES"))

(if (type [creatures] index) {}{
	[creatures = <map 
		SAURIAN=true
		HUMAN=true
	>]
})

(console:print "    Giving Specified Creatures Speluncaphobia...\n")
(foreach [rubble:raws] block name content {
	(if (str:cmp "creature_" (str:left [name] 9)) {
		(console:print "      " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "CREATURE") {
				(if (int:eq (len [tag]) 1) {
					(if [creatures [tag 0]]{
						[tag append = "
	[CAN_DO_INTERACTION:SPELUNCAPHOBIA]
		[CDI:TARGET:A:SELF_ONLY]
		[CDI:WAIT_PERIOD:33600]
		[CDI:FREE_ACTION]
	
	[CAN_DO_INTERACTION:SPELUNCAPHOBIA_SLOW]
		[CDI:TARGET:A:SELF_ONLY]
		[CDI:WAIT_PERIOD:3050]
		[CDI:FREE_ACTION]
"]
					})
				}{
					(rubble:abort "Error: invalid param count to CREATURE raw tag in last file.")
				})
			})
			
			(break true)
		})]
	})
	(break true)
})
