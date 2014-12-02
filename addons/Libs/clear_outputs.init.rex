
# This script clears the "out:" directory.

# EXPERTS ONLY!
(if (str:cmp (rubble:configvar "_RUBBLE_NO_CLEAR_") "true") {
	(ret)
})

# Recursively clears a directory.
var cleardir = block path {
	(axis:walkdirs [rubble:fs] [path] block fpath {
		([cleardir] (str:add [path] "/" [fpath]))
		
		(axis:del [rubble:fs] (str:add [path] "/" [fpath]))
		(onerror {
			(console:print "  Error deleting: " [path] "/" [fpath] "\n")
		})
	})
	
	(axis:walkfiles [rubble:fs] [path] block fpath {
		(axis:del [rubble:fs] (str:add [path] "/" [fpath]))
		(onerror {
			(console:print "  Error deleting: " [path] "/" [fpath] "\n")
		})
	})
}

# These are (or at least should be) safe to nuke.
([cleardir] "out:objects")
([cleardir] "out:graphics")
([cleardir] "out:scripts")

# These directories were created by Rubble in the first place, so go ahead and nuke'em.
([cleardir] "out:dfhack")
([cleardir] "out:Docs")

# Completely remove certain obsolete directories.
([cleardir] "out:prep")
(axis:del [rubble:fs] "out:prep")
(error false)

# Now clear out the various junk that has accumulated in the raw directory itself.
(axis:walkfiles [rubble:fs] "out:" block path {
	(axis:del [rubble:fs] (str:add "out:" [path]))
	(onerror {
		(console:print "  Error deleting: " [path] "\n")
	})
})
