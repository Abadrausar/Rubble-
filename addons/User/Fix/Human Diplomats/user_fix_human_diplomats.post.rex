
var diplomats =
"[POSITION:DIPLOMAT]
		[NAME:Diplomat:Diplomats]
		[NUMBER:1]
		[RESPONSIBILITY:MAKE_INTRODUCTIONS]
		[RESPONSIBILITY:MAKE_PEACE_AGREEMENTS]
		[RESPONSIBILITY:MAKE_TOPIC_AGREEMENTS]
		[PRECEDENCE:70]
		[DO_NOT_CULL]
		[COLOR:7:0:1]
		[MENIAL_WORK_EXEMPTION]
		[SLEEP_PRETENSION]
		[PUNISHMENT_EXEMPTION]
		[ACCOUNT_EXEMPT]
		[DUTY_BOUND]
	[POSITION:TRADE_REP]
		[NAME:Trade Representative:Trade Representatives]
		[NUMBER:1]
		[RESPONSIBILITY:TRADE]
		[PRECEDENCE:40]
		[DO_NOT_CULL]
		[COLOR:7:0:1]
		[MENIAL_WORK_EXEMPTION]
		[SLEEP_PRETENSION]
		[PUNISHMENT_EXEMPTION]
		[ACCOUNT_EXEMPT]
		[DUTY_BOUND]
"

var curEntity = ""

[rubble:raws "entity_default" = (df:raw:walk [rubble:raws "entity_default"] block tag {
	(if (str:cmp [tag id] "ENTITY") {
		(if (int:eq (len [tag]) 1) {
			[curEntity [tag 0]]
		}{
			(rubble:abort "invalid param count to ENTITY raw tag in entity_default")
		})
	})
	
	(if (str:cmp [curEntity] "PLAINS") {
		(if (str:cmp [tag id] "STONE_SHAPE") {
			(if (int:eq (len [tag]) 1) {
				[curEntity = ""]
				[tag replace = (str:add [diplomats] "[STONE_SHAPE:" [tag 0] "]")]
			}{
				(rubble:abort "invalid param count to STONE_SHAPE raw tag in entity_default")
			})
		})
	})
	
	(break true)
})]
