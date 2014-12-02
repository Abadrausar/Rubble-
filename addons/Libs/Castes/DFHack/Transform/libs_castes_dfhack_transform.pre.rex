
(rubble:template "DFHACK_REACTION_CASTE_TRANSFORM" block id class race caste delay=0 {
	(if (int:gt [delay] 0) {
		[rubble:dfhack_reactions [id] = (str:add 'rubble_caste_transform -unit \\\\WORKER_ID -race ' [race] ' -caste ' [caste] ' -delay ' [delay])]
	}{
		[rubble:dfhack_reactions [id] = (str:add 'rubble_caste_transform -unit \\\\WORKER_ID -race ' [race] ' -caste ' [caste])]
	})
	(rubble:reaction [id] [class])
})

(rubble:template "DFHACK_CASTE_TRANSFORM" block race caste delay=0 id="" {
	(if (str:cmp [id] "") {
		(if (str:cmp [rubble:reaction_last] "") {
			(rubble:abort "Error: rubble:reaction_last is invalid in call to DFHACK_CASTE_TRANSFORM.")
		})
		
		[id = [rubble:reaction_last]]
	})
	
	(if (int:gt [delay] 0) {
		[rubble:dfhack_reactions [id] = (str:add 'rubble_caste_transform -unit \\\\WORKER_ID -race ' [race] ' -caste ' [caste] ' -delay ' [delay])]
	}{
		[rubble:dfhack_reactions [id] = (str:add 'rubble_caste_transform -unit \\\\WORKER_ID -race ' [race] ' -caste ' [caste])]
	})
	(ret "")
})

(axis:write [rubble:fs] "out:scripts/rubble_caste_transform.lua" [rubble:raws "libs_castes_dfhack_transform.lua"])
