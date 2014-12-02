
(foreach [rubble:raws] block name content {
	(if (str:cmp "entity_" (str:left [name] 7)) {
		
		(console:print "    " [name] "\n")
		
		[rubble:raws [name] =
			(str:replace [content] "[PERMITTED_REACTION:TAN_A_HIDE]" "[PERMITTED_REACTION:TAN_A_HIDE_GLOB]" -1)]
		
	})
	(break true)
})