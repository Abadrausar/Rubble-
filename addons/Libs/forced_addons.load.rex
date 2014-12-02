
(console:print "    Forcing activation of critical addons:\n")

(foreach [rubble:addons] block _ addon {
	(if (str:cmp [addon Name] "Libs/Base") {
		(console:print '      "Libs/Base"\n')
		[addon Active = true]
	})
	(if (str:cmp [addon Name] "Libs/Castes") {
		(console:print '      "Libs/Castes"\n')
		[addon Active = true]
	})
	(break true)
})
