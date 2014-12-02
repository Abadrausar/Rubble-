
(axis:walkfiles [rubble:fs] "out:objects" block path {
	(axis:del [rubble:fs] (str:add "out:objects/" [path]))
	(if (error) {
		(console:print "Error deleting: " [path] "\n")
		(error false)
	})
})

(axis:walkfiles [rubble:fs] "out:graphics" block path {
	(axis:del [rubble:fs] (str:add "out:graphics/" [path]))
	(if (error) {
		(console:print "Error deleting: " [path] "\n")
		(error false)
	})
})

(axis:walkfiles [rubble:fs] "out:prep" block path {
	(axis:del [rubble:fs] (str:add "out:prep/" [path]))
	(if (error) {
		(console:print "Error deleting: " [path] "\n")
		(error false)
	})
})

(axis:walkfiles [rubble:fs] "out:dfhack" block path {
	(axis:del [rubble:fs] (str:add "out:dfhack/" [path]))
	(if (error) {
		(console:print "Error deleting: " [path] "\n")
		(error false)
	})
})

(axis:walkfiles [rubble:fs] "out:Docs" block path {
	(axis:del [rubble:fs] (str:add "out:Docs/" [path]))
	(if (error) {
		(console:print "Error deleting: " [path] "\n")
		(error false)
	})
})

(axis:walkdirs [rubble:fs] "out:Docs" block path {
	(axis:del [rubble:fs] (str:add "out:Docs/" [path]))
	(if (error) {
		(console:print "Error deleting: " [path] "\n")
		(error false)
	})
})

(axis:walkfiles [rubble:fs] "out:" block path {
	(axis:del [rubble:fs] (str:add "out:" [path]))
	(if (error) {
		(console:print "Error deleting: " [path] "\n")
		(error false)
	})
})
