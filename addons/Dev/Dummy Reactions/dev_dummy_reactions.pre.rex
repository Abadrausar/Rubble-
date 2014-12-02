
var howmany = (rubble:configvar "DEV_DUMMY_REACTION_COUNT")
(if (str:cmp [howmany] "") {
	[howmany = 15]
})

var out = '[OBJECT:REACTION]\n\nThe following are "dummy" reactions used to add content after worldgen.\n'
var count = 0
(loop {
	(if (int:lt [count] [howmany]) {
		[out = (str:add [out] 
			"\n{REACTION;DUMMY_REACTION_" [count] ";ADDON_HOOK_PLAYABLE}\n"
			"\tAdd your reaction body here\n"
		)]
		
		[count = (int:add [count] 1)]
		(break true)
	}{
		(break false)
	})
})

[rubble:raws "reaction_dev_dummy_reactions.txt" = [out]]
