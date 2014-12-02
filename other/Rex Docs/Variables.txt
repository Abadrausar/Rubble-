
Built-in Variables:
    rubble:version      The Rubble version string.
    rubble:versions     A map of all the versions in this series.
    rubble:addons       An indexable version of ALL loaded addons, active or not. 
                            This is a GenII value! Take care.
    rubble:addonstbl    This is supposed to be exactly like rubble:addons, just for lookup by name. 
                            This is a GenII value! Take care.
    rubble:files        An indexable version of all the files in the active addons.
                            Only valid after generation starts (aka pre-scripts or later)!
                            This is a GenII value! Take care.
    rubble:raws         This is an EditIndexable that stores string versions of all active addon files.
                            Only valid after generation starts (aka pre-scripts or later)!
    
    rubble:fs           A reference to the Rubble AXIS VFS file system, use with the axis commands.
    
    rubble:state        A reference to the Rubble State, Rubble uses this internally you should never need to use it.
                            This entry is just so you know it exists.
                            (please don't clobber it, Rubble will crash if you do)
