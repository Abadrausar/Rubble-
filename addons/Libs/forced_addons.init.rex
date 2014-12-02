
# init scripts are the only place you can use rubble:activate_addon safely,
# anywhere else runs the risk of not parsing all files

(console:print "    Forcing activation of critical addons:\n")

(console:print '      "Libs/Base"\n')
(rubble:activate_addon "Libs/Base")

(console:print '      "Libs/Castes"\n')
(rubble:activate_addon "Libs/Castes")
