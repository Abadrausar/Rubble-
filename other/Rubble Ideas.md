
If I don't write it down I forget it, hence this file.
If you see something nice in here that you want to use with your mod feel free to steal it :)

You know, the problem with formal TODO lists like this is that they rapidly fill up with things that I should do, but am too lazy to finish for the foreseeable future...

The fastest way to ensure something never gets done is to add it to this list, *sigh*

================================================================

Make the effect of "User/Speluncaphobia" wear off when maximum cave adaption hits.

"User/Speluncaphobia/DFHack" is bugged (that or DFHack is). Determine where the problem is and fix it (or get it fixed)!
	It's an item-trigger issue, wait for it to be fixed and things should start working.

Playable everyone:
	Remember: Most of the special stuff needs to be in generic addons, if you want to play dwarves with elven wood working that should be possible. The race addon should just set a general "flavor".
	
	Humans
		Done. (sort-of)
		I should do a crates addon to encourage trading.
		Make an embassy addon to allow forcing caravans, diplomats, and sieges.
	
	Elves
		Need a way to make grown items and special weapons grade "wood" logs.
		No real way to disable stone or metal working, so if you want an "elven heretic" fort...
	
	Goblins
		No eat and drink will make it easy to survive if you turtle up
		This may be a good time to make a raid system.
	
	Kobolds
		Will require a MAJOR rework of the entity!

Other playable races:
	Gnomes:
		Mostly just smaller dwarves for flavor, nothing special
		Excellent mechanics, stone workers, etc, poor anything surface related.
		Should they be good smiths or should gnomes focus on stone? Probably stone I think.
		Obviously using the powered workshops is encouraged

Finish the powered sensors:
	Sensors:
		Unit (variants for pets and hostiles?)
		Cart (settable cargo thresholds)
		(others?)
	Other:
		Repeater (wait y ticks, turn on, wait x ticks, turn off, repeat)

Crystal Thorn
	A slow growing (but not too slow!) farmable plant that has two uses:
		It may be processed to a low value raw gem
		It may be milled to sand
	~4 plants for each gem or 1 plant per unit of sand
		Processing to gems will require a custom reaction, milling to sand can probably use the vanilla system
		By using vanilla milling nothing special needs to be done to make them work with the powered mill addon

Port the old "BD/Farming/Basics", but with faster plant growth.
	The old addon was too hard (extra slow plant growth + complicated processing = starve).
		If plants grow too slowly they will not be able to compete with other ways to get food.
	Should I support the new "garden" plants or should they be removed?
		I am inclined to remove most of them, but expand the list of plants from the old version.
	Add more "special" plants (like the above idea), make sure there are several plants for each niche, etc.
		Some plants should be similar, but with different processing and/or byproducts
			For example several oil plants, some fast growing, but producing a cheap toxic oil (good only for soap) and some slow growing, but capable of producing an edible oil that can be cooked OR made into soap.
			Another example would be berries, they can be eaten or brewed, while some other plants may be ground then pressed, to provide both food AND drink from the same plant (although press cake would be less valuable than the whole plant).
	Now that I have a good system for config vars variable difficulty is possible...

A powered screw press.
	May be difficult, I may have to make this one less flexible than I would like.
	[HAS_MATERIAL_REACTION_PRODUCT:HONEYCOMB_PRESS_MAT] (unpressed honeycomb tools only)
	[HAS_MATERIAL_REACTION_PRODUCT:PRESS_LIQUID_MAT] (unpressed globs and plant growths only)
	Press to jugs:
		[EMPTY]
		[DOES_NOT_ABSORB]
		[HAS_TOOL_USE:LIQUID_CONTAINER]

Support milling items to paste at the powered mill
	Vanilla uses [HAS_MATERIAL_REACTION_PRODUCT:PRESS_LIQUID_MAT], but only for seed items
		This will need to be extensible by other addons.

A powered forge?
	Maybe, maybe not...

Port the "BD/Hives" addons.
	All of them require a good bit of reworking.
		The spider addon is designed for use with Essential DF (uses creatures that do not exist in vanilla)
		The rainbow beetle dyes desperately need to be changed (color names do not match dye descriptions)

Make the unpowered quarry require soil (if DFHack is available)

Gem forge
	Create armor, weapons, furniture, tools, etc from gems
	Probably requires a new material template for all gems (to make them sharp and hard)
	Mostly gems should make good edged weapons, but poor armor
	mostly decorative, good for nobles and fancy stuff like that (but gem swords may make a good early weapon)
	Not stockpilable? Not sure, but if so (and if there is no easy fix) this won't happen.
