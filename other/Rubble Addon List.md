
==============================================
Included Addons:
==============================================

Rubble includes all kinds of useful addons, some for developers and some for end users. These addons are included in an attempt to ensure that Rubble has something for everyone, and to provide a broad base of examples for aspiring addon developers.

If you want to generate the Rubble version of "vanilla DF" you will need to activate the "Base" addon, this addon is already active by default.

All addons will automatically activate any dependencies, so don't worry if the addon list that is generated is a little different from what you selected.

A few of the addons allow you to make gloves via reactions, unfortunately these gloves will be unusable (it's a vanilla bug). If you have DFHack this will be fixed automatically.

This list is automatically generated, for the full description of an addon see it's web UI addon description page.

[spoiler]
Base
	The standard base addon. Contains Rubblized versions of the vanilla raw files.

Dev/Cheat
	A bunch of cheat reactions to allow for fast setup of testing forts.

Dev/Dummy Reactions
	This addon adds "dummy reactions" eg. empty reactions registered to "ADDON_HOOK_PLAYABLE".

Dev/Tileset/Export
	This addon is a helper for extracting tileset information from raws that have the Rubble tileset templates added.

Dev/Tileset/Insert
	This addon adds the SHARED_OBJECT templates to all the tileset sensitive objects.

Libs/Base
	All the base templates and other such stuff, critical to normal operation.

Libs/Castes
	Adds support for dynamically adding castes to a creature.

Libs/Castes/DFHack/Transform
	Adds support for transforming castes.

Libs/Crates
	Adds infrastructure for "crates", eg bars that can be bought from caravans and then "unpacked" into items.

Libs/DFHack/Add Reaction
	Modify the reaction list of hard coded buildings.

Libs/DFHack/Announcement
	This addon provides templates for displaying announcements to the player.

Libs/DFHack/Fix Handedness
	Automatically installs and activates a fix for making gloves in reactions.

Libs/DFHack/Fluids
	Installs the extremely useful rubble.fluids Lua pseudo module and it's companion command rubble_fluids.

Libs/DFHack/Timeout
	Adds DFHack Lua support for persistent time delays (eg they survive a save/load).

Libs/DFHack/Upgrade Building
	Allows you to change the types of certain workshops after they have been built.

Libs/Macro
	A modders resource for generating macros from scripts.

Libs/Macro/Example
	A simple addon that generates a large complicated macro from images, makes a good demo/test of "Libs/Macro".

Libs/Research
	This is the old research addon from Better Dorfs, reworked as a library.

Libs/Wood Additions
	Adds several small tweaks to wood materials.

User/Adv Reactions
	Adds some basic adventure mode crafting reactions, mostly useful to help jump-start use of advfort.

User/Archery
	Adds the archers workshop, a building where you can make all types of ranged weapons and their ammo.

User/Block
	Support for the addons in the "Block" group.

User/Block/Boulder
	Adds the boulder press, which allows you to recolor boulders via the wonders of dwarven masonry.

User/Block/Cutter
	The block cutter allows you to harness the wonderful power of dwarven masonry to cut blocks of any color.

User/Block/Kiln
	Fire blocks of any color from clay.

User/Bone Flux
	Bake bones at the kiln to make calcite, a flux stone.

User/Bonecarver
	Adds a workshop for advanced shaping of bone.

User/Bonfire
	Allows you to start fires.

User/Building Dimensions
	Append all building names with their dimensions, really simple.

User/Cast Anvil
	Gives playable races a reaction to cast an iron anvil from 5 iron bars.

User/Cat Damper
	Removes ADOPTS_OWNER and RETURNS_VERMIN_KILLS_TO_OWNER from the cat, preventing catsplosion and widespread miasma.

User/Cave Color
	Changes color and availability of cave "grasses" to enhance the visual difference between caverns.

User/DFHack/Bulletin Board
	A bulletin board, where your idle dwarves will post thoughts, preferences, etc...

User/DFHack/Cart Filler
	Fill minecarts with water or magma (without big fancy mechanical contraptions).

User/DFHack/Dragon Engine
	A direct port of the dragon engines from the DFHack r5 mini-mod pack.

User/DFHack/Fill Barrel
	Fill barrels from nearby water at the still.

User/DFHack/Magma/Extractor
	Extract minerals from magma, get all kinds of rare and valuable metals! (but only rarely)

User/DFHack/Magma/Melter
	Melt rocks into magma!

User/DFHack/Obsidian Caster
	The obsidian caster makes obsidian boulders directly from adjacent water and magma.

User/DFHack/Pop Cap
	Allows you to set a population cap that will override the default for this world.

User/DFHack/Powered
	Powered workshop support addon.

User/DFHack/Powered/Block Cutter
	A powered version of the block cutter.

User/DFHack/Powered/Carpenter
	A powered carpenters workshop, make items from wood blocks.

User/DFHack/Powered/Cart Filler
	An automatic mechanically powered version of the "User/DFHack/Cart Filler" addon.

User/DFHack/Powered/Cart Launcher
	Launch full minecarts automatically.

User/DFHack/Powered/Cart Loader
	Fill minecarts with items automatically.

User/DFHack/Powered/Factory
	Adds a factory building for making batches of items from blocks, planks and raw glass.

User/DFHack/Powered/Glass Furnace
	A powered glass furnace, turn sand into glass.

User/DFHack/Powered/Kiln
	A powered kiln, turn clay into pottery boulders.

User/DFHack/Powered/Magma Extractor
	An automatic powered version of the "User/DFHack/Magma/Extractor" addon.

User/DFHack/Powered/Mason
	A powered masons workshop, make items from stone blocks.

User/DFHack/Powered/Obsidian Caster
	An automatic powered version of the "User/DFHack/Obsidian Caster" addon.

User/DFHack/Powered/Quarry
	A powered quarry, when built over clay provides a steady supply of boulders.

User/DFHack/Powered/Sawmill
	A powered version of the sawmill.

User/DFHack/Powered/Smelter
	A powered smelter, ore to bars, no dwarves required.

User/DFHack/Powered/Sorter
	Sort items by type and/or material.

User/DFHack/Powered/Wood Furnace
	A powered wood furnace, turn planks into ash or charcoal.

User/DFHack/Prospect
	Allows your miners to prospect the area for likely ores and gems.

User/DFHack/Spatter
	A streamlined version of the DFHack spatter example mod.

User/DFHack/Steam
	The DFHack steam engine example mod.

User/DFHack/Stirling Engine
	Produces a (very) small amount of power for free.

User/Deadly Bones
	Makes creatures take 10x damage from weapons made from their bones.

User/Decorations/Display Case
	Adds a workshop for displaying items.

User/Decorations/Display Case/DFHack
	Gives the display case a pretty custom sidebar menu that displays the contained item.

User/Diggable Slade
	Allows you to dig slade, once you get past it's guardians that is...

User/Dwarf/Castes
	Adds more dwarf castes.

User/Dwarf/Castes/DFHack/Transform
	Allows you to change the caste of your dwarves at the guild hall for a nominal fee.

User/Edible Vermin
	Allows you to cook vermin at the kitchen to get cheap meat, a good emergency food source.

User/Factory
	Adds a large factory building for making batches of cheap items (mostly from wood and stone).

User/Fix/Butcher Inorganic
	Allows you to butcher creatures made from (some) inorganics, FB stone walls anyone?

User/Fix/Fish Pops
	Raises fish vermin populations to the max.

User/Fix/Gays
	Makes humans and dwarves have a strictly straight orientation.

User/Fix/Human Diplomats
	Adds trade representatives and diplomats to the humans.

User/Fix/Undead Melt
	Fixes bad temperature values in material_template_default.

User/Fix/Usable Chitin
	Allows you to use chitin as shell.

User/Fix/Usable Feathers
	Allows you to use feathers as "soft" pearls (you can use them to make crafts).

User/Fix/Usable Scale
	Allows you to use scale like leather.

User/Fix/Vermin Variations
	Makes giant and animal man vermin show up properly.

User/Fix/Whips
	Nerfs whips and scourges so they are not deadly anti-armor weapons.

User/Generic Animal Mats
	Support for the the addons in the "Generic Animal Mats" group.

User/Generic Animal Mats/Blood and Pus
	Generic animal blood and pus.

User/Generic Animal Mats/Bone
	Generic animal bones.

User/Generic Animal Mats/Fat, Tallow, and Soap
	Generic animal fat, tallow, and soap.

User/Generic Animal Mats/Hair and Feathers
	Generic animal hair and feathers.

User/Generic Animal Mats/Horn and Hoof
	Generic animal horn and hooves.

User/Generic Animal Mats/Meat
	Generic animal meat (including organs).

User/Generic Animal Mats/Milk
	Generic animal milk and cheese.

User/Generic Animal Mats/Shell
	Generic animal shell.

User/Generic Animal Mats/Sinew
	Generic animal sinew.

User/Generic Animal Mats/Skin, Scale and Chitin
	Generic animal hide materials.

User/Generic Animal Mats/Tooth and Claw
	Generic animal teeth and claws.

User/Glass Forge
	The glass forge allows you to make armor and weapons from all three kinds of glass.

User/Lamination
	Glue wooden planks together to make laminated logs, which are more valuable than most wood.

User/Metallurgy
	Adds an "Alloy Furnace" that can be used to make alloys.

User/Metallurgy/Aluminum Bronze
	Adds a new alloy to the alloy furnace, aluminum bronze.

User/Metallurgy/Bloomery
	Adds a "Bloomery" that can be used to make iron and steel in a more realistic manner.

User/Metallurgy/DFHack/Volcanic
	Adds a new alloy to the alloy furnace, volcanic.

User/Metallurgy/Mithril
	Adds a new alloy to the alloy furnace, mithril.

User/Metallurgy/Smelter
	Replaces the default ore smelting reactions with a better set.

User/Metallurgy/Smelter/Remove Alloys
	Removes all the vanilla alloy reactions, it is not a good idea to use this without "User/Metallurgy"!

User/Pottery
	Replace the vanilla pottery system with a much expanded version.

User/Pottery/Batch
	Adds a batch kiln that can make 5x the stuff with 4x the material and 1x the fuel.

User/Pottery/Glaze
	Greatly expands the glazing system.

User/Pottery/Porcelain
	Allows you to use orthoclase and microcline to make porcelain (in addition to kaolinite).

User/Quarry
	Dig for some basic raw materials without needing to actually dig.

User/Research
	Build a laboratory and research advanced workshops.

User/Saurian
	Adds a new playable race, the saurians.

User/Saurian/Castes
	Adds some basic castes to the saurians.

User/Saurian/Castes/DFHack/Transform
	Allows you to change the caste of your saurians at the meditation circle.

User/Saurian/Unplayable
	If you just want saurians as an extra race to trade with the dwarves, activate this addon.

User/Sawmill
	Farmable wood plants and advanced wood processing, makes aboveground and/or treeless forts a lot easier.

User/Smoked Meat
	Allows you to smoke meat at the kitchen to preserve it (inedible to vermin and will not rot).

User/Tanning
	Allows creatures to give more than one skin based on size.

User/Tilesets/MLC/DFHack/TWBT
	A simple addon that installs my custom ASCII-like tileset (Text Will Be Text version).

User/Tilesets/MLC/Normal
	A simple addon that installs my custom ASCII-like tileset (normal version).

User/Tilesets/MLC/Pure ASCII
	A simple addon that installs my custom ASCII-like tileset (pure ASCII version).

User/Tilesets/Vanilla
	Restores the vanilla tileset related init options.

User/Tilesets/Vanilla/Graphics
	The vanilla creature graphics, mostly just a demo.

User/Training
	Train your dwarves by reading books.

User/Warcrafter
	A special workshop for all your early game armor and ammo needs.

User/Warcrafter/Adventure
	Adventure mode reactions for armor, weapons, and other useful items.

User/Wicker
	Make all sorts of stuff from woven plants.

User/Zap Aquifers
	Removes all aquifers, simple.
[/spoiler]
