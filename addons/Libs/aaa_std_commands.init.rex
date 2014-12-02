
# Load some helpful commands.

# Rubble specific commands

# Placeholder template generator.
command rubble:placeholder name {
	(rubble:template [name] block ... {
		"Library template addon disabled."
	})
}

# Dependency checks

command rubble:addonactive name {
	(foreach [rubble:addons] block _ addon {
		(if (str:cmp [addon Name] [name]) {
			(if [addon Active] {
				(ret true)
			}{
				(ret false)
			})
		})
	})
	(ret false)
}

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

command rubble:requireaddon me addon {
	(if (rubble:addonactive [addon]) {
	}{
		(rubble:abort (str:add "The \"" [me] "\" addon requires that the \"" [addon] "\" addon be active!\n"
		"Please activate that addon and try again."))
	})
}

command rubble:incompatibleaddon me addon {
	(if (rubble:addonactive [addon]) {
		(rubble:abort (str:add "The \"" [me] "\" addon is incompatible with the \"" [addon] "\" addon!\n"
		"Please deactivate that addon and try again."))
	})
}

# Addon group checks

command rubble:ingroup addon group {
	[group = (str:add [group] "/")]
	
	(if (str:cmp (str:left [addon] (str:len [group])) [group]) {
		(ret true)
	} {
		(ret false)
	})
}

command rubble:groupactive group {
	var found = false
	(foreach [rubble:addons] block _ addon {
		(if [addon Active] {
			(if (rubble:ingroup [addon Name] [group]) {
				[found = true]
				(break false)
			}{
				(break true)
			})
		}{
			(break true)
		})
	})
	(ret [found])
}

# Example (rubble:grouprequires "Better Dorfs" "Better Dorfs/Base")
command rubble:grouprequires me addon {
	(if (rubble:groupactive [me]) {
		(if (rubble:addonactive [addon]) {
		}{
			(rubble:abort (str:add "The \"" [me] "\" addon group requires that the \"" [addon] "\" addon be active!\n"
			"Please activate that addon and try again."))
		})
	}{
		# No addons in the group are active, do nothing.
	})
}

command rubble:groupincompatible me addon {
	(if (rubble:groupactive [me]) {
		(if (rubble:addonactive [addon]) {
			(rubble:abort (str:add "The \"" [me] "\" addon group is incompatible with the \"" [addon] "\" addon!\n"
			"Please deactivate that addon and try again."))
		})
	}{
		# No addons in the group are active, do nothing.
	})
}

# All hail Rex/GenII!

command rubble:skipfile name {
	(if (exists [rubble:files Files] [name]) {
		[rubble:files Files [name] Skip = true]
	}{
		(rubble:abort (str:add 'rubble:skipfile: "' [name] '" is not the name of a loaded raw file.'))
	})
}

command rubble:nowritefile name {
	(if (exists [rubble:files Files] [name]) {
		[rubble:files Files [name] NoWrite = true]
	}{
		(rubble:abort (str:add 'rubble:nowritefile: "' [name] '" is not the name of a loaded raw file.'))
	})
}

command rubble:graphicsfile name {
	(if (exists [rubble:files Files] [name]) {
		[rubble:files Files [name] Graphics = true]
	}{
		(rubble:abort (str:add 'rubble:graphicsfile: "' [name] '" is not the name of a loaded raw file.'))
	})
}

command rubble:prepfile name {
	(if (exists [rubble:files Files] [name]) {
		[rubble:files Files [name] PrepScript = true]
	}{
		(rubble:abort (str:add 'rubble:prepfile: "' [name] '" is not the name of a loaded raw file.'))
	})
}
