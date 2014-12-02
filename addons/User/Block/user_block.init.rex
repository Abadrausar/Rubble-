
(if (rubble:groupactive "User/Block"){
	(console:print "    One or more block addons are active, forcing dependency activation.\n")
	(rubble:activate_addon "User/Block")
})
