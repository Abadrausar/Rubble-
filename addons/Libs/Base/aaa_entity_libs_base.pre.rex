
var rubble:adventure_tier_data = 0
var rubble:entity_playability = <map>
var rubble:entity_nobles = <map>
var rubble:entity_nobles_add = <map>

command rubble:entity_playable id key {
	(if (exists [rubble:entity_playability] [id]) {}{
		(ret false)
	})
	(if (exists [rubble:entity_playability [id]] [key]) {
		(if [rubble:entity_playability [id] [key]] {
			(ret true)
		}{
			(ret false)
		})
		(ret "")
	})
	(ret false)
}

# {@ENTITY_PLAYABLE;MOUNTAIN;true;true;false}
(rubble:template "!ENTITY_PLAYABLE" block id fort adv indiv {
	[rubble:entity_playability [id] = <map
		fort = [fort]
		adv = [adv]
		indiv = [indiv]
	>]
	(str:add "{#_ENTITY_PLAYABLE;" [id] "}")
})

# For internal use only!
(rubble:template "#_ENTITY_PLAYABLE" block id {
	(if (expr "!a && !b && !c" 
		[rubble:entity_playability [id] fort]
		[rubble:entity_playability [id] adv]
		[rubble:entity_playability [id] indiv])
	{
		(ret "# Entity is Not Playable.")
	})
	
	var out = ""
	(if [rubble:entity_playability [id] fort] {
		[out = (str:add [out] "[CIV_CONTROLLABLE]")]
	})
	(if [rubble:entity_playability [id] adv] {
		[rubble:adventure_tier_data = (int:add [rubble:adventure_tier_data] 1)]
		[out = (str:add [out] "[ADVENTURE_TIER:" [rubble:adventure_tier_data] "]")]
	})
	(if [rubble:entity_playability [id] indiv] {
		[out = (str:add [out] "[INDIV_CONTROLLABLE]")]
	})
	(ret [out])
})

(rubble:template "ENTITY_PLAYABLE_EDIT" block id key value {
	[key = (str:tolower [key])]
	
	(if (exists [rubble:entity_playability] [id]) {}{
		(rubble:abort "Attempt to change non-existent entity in call to ENTITY_PLAYABLE_EDIT.")
	})
	(if (exists [rubble:entity_playability [id]] [key]) {
		[rubble:entity_playability [id] [key] = [value]]
		(ret "")
	})
	(rubble:abort "Invalid playability key in call to ENTITY_PLAYABLE_EDIT.")
})

(rubble:template "@IF_ENTITY_PLAYABLE" block id key then else="" {
	[key = (str:tolower [key])]
	
	(if (exists [rubble:entity_playability] [id]) {}{
		(ret (rubble:stageparse [else]))
	})
	(if (exists [rubble:entity_playability [id]] [key]) {
		(if [rubble:entity_playability [id] [key]] {
			(ret (rubble:stageparse [then]))
		}{
			(ret (rubble:stageparse [else]))
		})
		(ret "")
	})
	(ret (rubble:stageparse [else]))
})

(rubble:template "#ENTITY_NOBLES" block id default {
	(if (isnil [rubble:entity_nobles [id]]) {
		(rubble:stageparse (str:add [default] (if (isnil [rubble:entity_nobles_add [id]]) {""}{
			(str:add "\n\t" [rubble:entity_nobles_add [id]])
		})))
	}{
		(rubble:stageparse [rubble:entity_nobles [id]])
	})
})

(rubble:template "ENTITY_ADD_NOBLE" block id noble {
	(if (isnil [rubble:entity_nobles_add [id]]) {
		[rubble:entity_nobles_add [id] = [noble]]
	}{
		[rubble:entity_nobles_add [id] = (str:add [rubble:entity_nobles_add [id]] "\n\t" [noble])]
	})
})

(rubble:template "ENTITY_REPLACE_NOBLES" block id nobles {
	[rubble:entity_nobles [id] = [nobles]]
})
