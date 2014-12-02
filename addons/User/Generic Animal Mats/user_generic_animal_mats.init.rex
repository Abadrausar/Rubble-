
(if (rubble:groupactive "User/Generic Animal Mats"){
	(console:print "    One or more generic animal material addons are active, forcing dependency activation.\n")
	(rubble:activate_addon "User/Generic Animal Mats")
})
