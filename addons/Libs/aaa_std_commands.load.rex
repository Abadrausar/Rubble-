
# Load some helpful commands.

# Is the named addon active?
command rubble:addonactive name {
	var addonref = [rubble:addonstbl [name]]
	(if (isnil [addonref]) {
		(ret false)
	}{
		(ret [addonref Active])
	})
}

# Activate an addon (and its dependencies).
# You should generally use addon.meta instead, but just in case...
command rubble:activateaddon me addon {
	(if (rubble:addonactive [me]) {
		var addonref = [rubble:addonstbl [addon]]
		(if (isnil [addonref]) {
			(rubble:abort (str:add "The \"" [me] "\" addon requires the \"" [addon] "\" addon!\n"
			"The required addon is not currently installed, please install the required addon and try again."))
		}{
			[addonref Active = true]
			(foreach [addonref Meta Activates] block _ childname {
				(rubble:activateaddon [addon] [childname])
			})
		})
	})
}

# Check the Rubble version.
command rubble:checkversion addon version {
	(if (isnil [rubble:versions [version]]) {
		(rubble:abort (str:add [addon] " requires Rubble version " [version] " (or a compatible newer version)\n"
			"The current version is older than the requested version (or it is incompatible),\n"
			"Please install the correct Rubble version and try again."))
	}{
		(if (str:cmp [version] [rubble:version]) {
			# All good
		}{
			(console:print (str:add 
				"    " [addon] " requires Rubble version " [version] " (or a compatible newer version)\n"
				"    The current Rubble version is:	" [rubble:version] " Which is newer than requested.\n"
				"    If you encounter issues changing to the requested Rubble version may help.\n"
			))
		})
	})
}
