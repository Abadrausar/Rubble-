
(rubble:requireaddon "Libs/Macro/Example" "Libs/Macro")

# A cool fortress entrance, better have a bunch of miners...
# I usually section this up and only dig the core passage to the central stairwell at first,
# else you will spend far too much time digging things that you don't need until much later.

# If you can, see the image "Pricetraded - Annotated Entry Blueprint.xcf", it has many more details
# (including instructions for building the seesaw) all neatly written on overlays for each layer.
# If you have a way to view that image (eg the GIMP or the like) then there is no real reason to read
# most of these comments.

# Generally I place this 15-20 Z-levels below the surface and build my entire fort above the main level
# This requires extending the spiral wagon ramp all the way to the surface (for caravan access).

# It is HIGHLY recommended to start with your cursor aligned to the grid, as everything
# in this design is sized according to the distance the cursor moves when Shift+<direction> is pressed.

# The little tunnel surrounding the wagon ramp is for use in engraving the wall above the down ramps,
# they can also be used for marksdwarves later...

# All invader accessible passages are separated from the walls by a channel, feel free to carve
# fortifications in the walls (Bonus Idea: Carve fortifications in only one wall so you can recover
# any bolts that missed and ended up in the channels after the battle).

# The construct on the far left passage is a "Seesaw of Armok", a kind of powerful minecart grinder.
# If you do not know how to build one then go ahead and redo that side of the fort with your favorite
# kind of trap hall, if you remove the seesaw don't forget the parts that are above the main level.
# (this is getting long enough without trying to explain how to build that too)
# The idea is to make two four barrel seesaws. The lever room has space for two levers per seesaw,
# one for the main operation and one for a special set of anti-flier hatch covers I normally install.
# The reason the walls are dug in a serrated patten on the lower level is so that they may be replaced
# with platinum, because being smashed between a platinum minecart and a rock wall just isn't as cool
# as being smashed between a platinum minecart and a platinum wall... If you are more practical feel
# free to edit the plans and save some work, after all you can always change your mind later :)

# Speaking of the far left side, the four rooms here are for military training, and the area below
# these rooms is for burial. If you tend to lose lots of dwarves you may wish to ditch my fancy
# architecture for something with more floor space, just be careful not to dig into the channels...

# The odd room on the main level near the center stairwell is the lever room. The undesignated spots 
# (shaped like some demented animal) are where the levers go eventually, the 2x2 block on the
# far left side is for the seesaw, if you are not building the seesaw then obviously these spots
# can be used for other things. (the lever room is "position coded", eg every lever spot is placed
# in the room relative to where it's bridge is in the fort)

# If you do not need lever spots for the trap hall (the one to the far left) at all, then move all
# lever positions one tile to the left, this will keep the levers centered in the room.

# It should be obvious where bridges go, all bridges are 3X5 and raise in such a way as create an
# impenetrable barrier (so long as the hatch covers are installed and locked)

# There are many, many down stairs on the main level leading into the channels under where the bridges
# go. These down stairs MUST be covered by locked hatch covers unless you want every dodging goblin to
# have access to your fort! It is recommended to have at least a small military to make cleaning these
# channels possible

# If you wish the gap between the last set of bridges on each path can be filled with traps, this is
# your best anti hidden necromancer defense...

# The dead end on the far right side is for extension to a possible cavern exit (for caravans to use
# to leave the fort during a siege).

# 2-Z above the main level I have a bedroom level, you may or may not want this, if not feel free to
# remove it.

# The small room directly above the entry tunnel is a F.R.O.G. (Fiend Revealing Observation Garrison)
# Cover the holes in the floor with grates and station some animals here for ambush spotters.
# This is only for last ditch detection, as anybody who gets this far can charge into the fort before
# you can seal up unless the center path is already closed!
# Unless you have another FROG farther up stream it is best to keep the center path sealed unless 3/4
# of the fort is parading back and forth to the surface cleaning up a dead caravan or something.
# (in which case a military is advised)

# The big room in the middle of the right path is for the trade depot.

# I use this entry in almost every fort I build, as it allows for a very powerful and flexible defense
# when finished.

# An XCF file is included, it has no use other than providing an easy way to look at the blueprint as one layered image.

var mak = (rubble:macro:new)

# The colors here are less than perfect, I should have stuck to the primaries and secondaries,
# but the basic blueprint predates MacroGen (and by extension the Rubble macro commands) by quite a bit,
# and so I never thought about needing to specify the colors by hand. Oh well, it only had to be done once...
var palette = <smap
	0xffffffff = <sarray DESIGNATE_DIG>
	0xff0000ff = <sarray DESIGNATE_STAIR_DOWN>
	0xff7e00ff = <sarray DESIGNATE_STAIR_UPDOWN>
	0x96ff00ff = <sarray DESIGNATE_STAIR_UP>
	0x9f9f9fff = <sarray DESIGNATE_CHANNEL>
	0xea00ffff = <sarray DESIGNATE_RAMP>
	0xffff00ff = <sarray DESIGNATE_STAIR_DOWN SELECT SELECT DESIGNATE_STAIR_UPDOWN SELECT SELECT>
>

[mak x = 46]
[mak y = 56]

(for 0 -5 -1 block i {
	(rubble:macro:setcursor [mak] [mak x] [mak y] [i])
	(rubble:macro:designate_rle [mak] [palette] (png:load [rubble:raws (str:add "Entry " [i] ".png")]))
	
	(break true)
})

(rubble:macro:setcursor [mak] 46 56 -4)

(rubble:macro:install [mak] "Libs_Macro_Example - Entry")
