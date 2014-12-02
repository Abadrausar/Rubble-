
(rubble:template "DFHACK_REACTION_ANNOUNCE" block id class text color="COLOR_WHITE" {
	# *sigh* the things I do to avoid major changes...
	(rubble:stageparse (str:add "{REACTION;" [id] ";" [class] "}"))
})

(rubble:placeholder "DFHACK_ANNOUNCE")
