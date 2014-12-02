
Rubble: After Blast comes Rubble

==============================================
Overview:
==============================================

In Rubble many little (and not so little) mods (called addons) are combined and run through a special script powered parser to generate standard Dwarf Fortress raw files.

These addons are, by construction, able to be assembled in many configurations with minimal direct internal dependencies. Most of the time an addon is completely standalone or dependent on only a few common items/materials present in the vanilla raws.

The beauty of this system is two fold: on one hand users have many choices and can construct their own private version of DF with minimal effort, on the other hand modders can make use of the power of Rubble to make installation of their mod automatic and use the template/scripting system to automate most if not all of the more repetitive parts of modding.

Of course not all addons are made for users, whole groups of addons are made specially for modders to help with testing and automating certain common tasks (not to mention library addons full of ready to use templates and scripts).

Rubble has been in continuous development since ~6/13, and many changes and improvements have been made since the first version (which kinda sucked :p). All of my mods have been made to use it, so I have extensive experience with Rubble modding and along the way fixed most of the bugs and streamlined things as much as possible both for modders and users.

I hope Rubble fills your needs for a general content installer and modding tool (and if not post your suggestions!)

==============================================
Where to Find Help:
==============================================

If you do not plan on doing any modding with Rubble the only thing you need to read is this readme.
(particularly the "Install" section)

Required Reading for Modders (in most-important-first order):
	This readme
	HowTo Rubble - Tutorials
	Rubble Basics - Lots of stuff about addons and Rubble in general
	Rubble Base Templates - The template documentation for the "Base/Templates" addon
	Rubble Native Templates - The template documentation for the built-in templates
	Rubble Libs Templates - The template documentation for the library templates
	Everything in "Rex Docs" - the Rex documentation, NEEDED if you plan to do any scripting
	The included addons - A little short on comments but full of great examples

==============================================
Install:
==============================================

If you have anything in your raw folder that you want to keep back it up! Rubble will delete all your existing raw files, including creature graphics and the DFHack onLoad.init and init.lua files!

Delete or otherwise remove any old Rubble version you may have.
Extract the new Rubble to "<DF Directory>/rubble".
Install any custom addons you may have to "<DF Directory>/rubble/addons"
Run "rubble -addonlist" OR run "Rubble GUI.exe"

Now you are good to go! Documentation is in the "rubble/other" folder as is source code and OSX/Linux binaries.

If you use OSX or Linux, 32 bit binaries for these OSes can be found in the "rubble/other" directory. If you want 64 bit binaries you can compile them yourself, source code is in "rubble/other/src" (along with basic build instructions).

==============================================
Running Rubble (Manually):
==============================================

Windows users can skip this section, just run the GUI, it does all this (mostly) automatically.

Rubble is a command line application, therefore at least some knowledge of how to use a command prompt will be very helpful.

If you can (eg. you are running on windows) it is HIGHLY recommended to use the GUI. The GUI automates the process of updating and editing the addon list and is generally much faster then doing everything by hand (plus you don't have to mess around with the command prompt, if you dislike that kind of thing)

For basic documentation on each command line option run "rubble -h".

Common Tasks:

To activate or deactivate a Rubble addon manually you may set it's entry in addons/addonlist.ini to "true" or "false". If you just installed an addon it will not have an entry until Rubble has run at least once (after the addon was installed). 
If you want to run Rubble without generating anything so as to update the addon list file just run 'rubble -addonlist' (the GUI does this as part of it's startup process).

Some addons may allow additional configuration via "config variables", these are generally for advanced users and may be specified with the -config command line option.

If you want to run multiple worlds with radically different addon sets it is a good idea to run "rubble -prep=<region name>" before switching worlds (this is mostly only for tilesets now days, other addons should be good to go unless they do something weird and non-standard).

All the Rubble default addons are deterministic, meaning that it is possible to generate the raws twice and (as long as you use the same addons) each time the raws will be exactly the same. This makes it possible to regenerate the raws for a world to allow things like switching tilesets. Unfortunately this is a semi-complicated task to handle by hand, as there is a lot of overhead that the GUI would normally have handled that command line users will have to do manually.
The first step is to change the "addonlist.ini" file that is in the worlds raw directory, make sure not to change it too much or you can mess up your world!
If you wish you can change "genconfig.ini" as well, but that is usually a bad idea.
Now for regenerating the raws, what follows is an example command line to do that for the save "region1":
	rubble -zapaddons -zapconfig -addons="df:data/save/region1/raw/addonlist.ini" -config="df:data/save/region1/raw/genconfig.ini" -outputdir="df:data/save/region1/raw"

==============================================
Included Addons:
==============================================

Rubble 4 includes all kinds of useful addons, some for developers and some for end users. These addons are included in an attempt to ensure that Rubble has something for everyone, and to provide a broad base of examples for aspiring addon developers.

If you want to generate the Rubble version of "vanilla DF" you will need to activate the "Base" addon, this addon is already active in the default addonlist.ini ("Libs/Base" and "Libs/Castes" are also required, but those addons are forced active so you do not need to worry about them).

All of the DFHack stuff is made under the assumption that what worked with 34.11 and DFHack r5e1 will work with the 40.x DFHack when it is available.
If this is not true some or all of the DFHack related addons will not work.

A few of the addons allow you to make gloves via reactions, unfortunately these gloves will be unusable (it's a vanilla bug). To make these gloves usable you need a DFHack script to set the glove handedness.

Base:
	The standard base addon. Contains Rubblized versions of the vanilla raw files. 
	Do not disable unless you have a replacement.
	The following changes have been made from vanilla:
		Some of the files have been re-indented (tag order is not effected).
		The dwarf creature has templates from the "Libs/Castes" addon installed.
		Specialized SHARED_OBJECT templates have been inserted where appropriate.
		REACTION and BUILDING_WORKSHOP templates have been inserted.
		@BUILD_KEY templates have been inserted.
		#ADDON_HOOKS have been added to all entities.
		#ADVENTURE_TIER tags have been added to all adventure mode playable entities.
		Vermin creatures have #TILE and #COLOR templates inserted.
		All file IDs have been stripped.
	Requires "Libs/Base" and "Libs/Castes" (both of which are automatically activated.)

Dev:
	The addons in this group are for addon developers, most people may safely ignore them.

	Dev/Cheat:
		A bunch of cheat reactions to allow for fast setup of testing forts.
		Build a "Heavy Supply Wagon" workshop, both the building and all reactions are free.
		Uses animal caretaker skill (the most useless skill I could think of :p)
	
	Dev/Dummy Reactions:
		This addon adds "dummy reactions" eg. empty reactions registered to "ADDON_HOOK_PLAYABLE".
		By default 15 reactions are created, if you want a different number of reactions to be generated pass '-config="DEV_DUMMY_REACTION_COUNT=count"' to Rubble when generating.
		These reactions are empty by default, use if you want to add new content after worldgen.
	
	Dev/Tileset/Export:
		This addon is a helper for extracting tileset information from raws that have the Rubble tileset templates added. See "Dev/Tileset/Insert".
		The extracted tile numbers will be written to "./dev_export_tileset.tile_info.rbl".
		Basically an automatic replacement for "#WRITE_TILESET".
	
	Dev/Tileset/Insert:
		This addon adds the SHARED_OBJECT templates to all the tileset sensitive objects.
		Advanced users may want to set the config var "DEV_TILESET_INSERT_EXPORT=true" in order to get the resulting files output to the current working directory, this is a great help when porting a mod to Rubble.
		Use with "Dev/Tileset/Export" to export tile/color info from non-Rubble raws.
		This addon also handles vermin creature tile numbers and colors.

Libs:
	The addons in this group are libraries, eg they provide functionality used by other addons.
	Any templates added by these addons are fully documented, see "Rubble Libs Templates.txt" and "Rubble Base Templates.txt".

	Libs/Base:
		All the base templates and other such stuff, critical to normal operation.
		This addon is a "critical addon" meaning that it is needed by the base for normal operation.
		You do not need to activate this addon as it is automatically forced active.
	
	Libs/Castes:
		Adds support for dynamically adding castes to a creature.
		This addon is a "critical addon" meaning that it is needed by the base for normal operation.
		You do not need to activate this addon as it is automatically forced active.
	
	Libs/Castes/DFHack/Transform:
		Adds support for transforming castes via autoSyndrome and syndromeTrigger.
		Requires the "Libs/Castes" addon.
	
	Libs/Crates:
		Adds infrastructure for "crates", eg bars that can be bought from caravans and then "unpacked" via special reactions into items.
		Most of the templates added by this addon are always available (even when the addon is not active).
	
	Libs/DFHack/Add Reaction:
		Allows you to add reactions to buildings that are normally not possible to add reactions to, also allows you to remove all hardcoded reactions from a building.
		Requires the "Base/Templates" addon.
		WARNING! Largely untested!
	
	Libs/DFHack/Announcement:
		This addon provides templates for displaying announcements to the player.
		The announcements work via specialized variants of the DFHACK_REACTION and DFHACK_REACTION_BIND templates.
		The templates added by this addon is always available (even when the addon is not active).
	
	Libs/DFHack/Fluids:
		Adds the extremely useful rubble_fluids DFHack Lua library and command.
		(adds no templates)
	
	Libs/Macro:
		A modders resource for generating macros from scripts.
		Adds no templates, only script commands.
	
	Libs/Macro/Example:
		A simple addon that generates a large complicated macro from images, makes a good demo/test of "Libs/Macro".
		This addon takes ~7-8 seconds to generate (the cursor movement generation command is kinda slow, it's a work in progress), in any case activating this addon more than triples normal run time.
	
	Libs/Research:
		This is the old research addon from Better Dorfs, reworked as a library.
		All templates now take a "class" (basically a group id) and the reaction generator takes a building name, skill id, and reaction class. These changes allow modders to use the same framework for multiple different research buildings.
		The templates for adding new research topics are always active, so optional support for addons that use this library is easy to add to your industry mods.
		Many default addons register topics to the class "RESEARCH_GENERIC", so handle that class to if you want to support these addons topics (nothing bad will happen if you don't).

User:
	All the addons in this group are things that are of interest only to end users.
	
	User/Add CHILD Tags:
		Adds a CHILD tag to every creature that does not already have one, thereby allowing them to breed in fortress mode.
		(Note that sexless creatures will still not breed)
		DOES_NOT_EXIST, NOT_LIVING, and EQUIPMENT_WAGON count as CHILD tags for the purpose of this script, vermin are also skipped.
	
	User/Adv Reactions:
		Adds some basic adventure mode crafting reactions, mostly useful to help jump-start use of advfort.
		Allows you to make crude picks and axes from knapped stone, anvils from boulders, backpacks from bags (that you can find in any human hamlet), and "emergency" crutches from bone.
		All adventure mode reactions use the knapping skill.
	
	User/Block/Bone:
		Good for really intimidating statues over the entry.
		For those who just want "bone blocks" instead of "goblin bone blocks" or "elf bone blocks" (that BTW, cannot be stockpiled) there is an extra reaction that allows you to make blocks from a "bone" stone instead of the actual bone material.
		You can also make "treated bone blocks" that are fire safe, but you only get half as many per reaction.
	
	User/Block/Boulder:
		For the OCD among us, I present the boulder press. The boulder press allows you to recolor boulders via the wonders of dwarven masonry. 
		How does it work? I have no idea, but dwarven masons are renowned for their ability to do strange things with stone :)
		Like the other shops of it's kind the boulder press allows you to spend twice as much to get a better quality product.
		Fine boulders are not only worth more, but magma safe as well.
	
	User/Block/Cutter:
		The block cutter allows you to harness the wonderful power of dwarven masonry to cut 4 blocks of any color from 1 hard non-economic boulder of any kind. If you want to make higher value constructions you can also cut a boulder into 2 better quality blocks instead.
		Fine blocks are not only worth more, but magma safe as well.
	
	User/Block/Kiln:
		The Block Kiln works exactly like the block cutter except that you fire blocks from clay instead of cutting them from boulders. 5 clay gets you 10 normal blocks or 5 fine blocks.
		Fine blocks are not only worth more, but magma safe as well.
	
	User/Bone Flux:
		Bake bones at the kiln to make calcite, a flux stone.
		Very useful for steel making in flux-poor locations.
	
	User/Bonfire:
		Start a fire and burn them goblins! Either start a short wood fueled fire or start a longer coke fueled one.
		Good for starting brush fires or burning trash.
	
	User/Building Dimensions:
		Append all building names with their dimensions, really simple.
	
	User/Cast Anvil:
		Gives playable races a reaction to cast an iron anvil from 5 iron bars.
	
	User/Cat Damper:
		Removes ADOPTS_OWNER and RETURNS_VERMIN_KILLS_TO_OWNER from the cat, preventing catsplosion and widespread miasma.
	
	User/Cave Color:
		Changes color and availability of cave "grasses". Level 1 is green, yellow and brown, level 2 is mostly blue and cyan, and level 3 is red and gray. 
		Changes are cosmetic only. Does not change non-grass plants or trees.
		Makes it easy to tell which cavern level you are looking at.
	
	User/Decorations/Display Case:
		Adds a useless workshop that can be built with 1 block and any 1 other item, build with an artifact to make room value skyrocket or just use to keep things out of reach of those pesky thieves.
	
	User/Decorations/Display Case/DFHack:
		Gives the display case a pretty custom sidebar menu that displays the contained item.
		The script used is a direct, unchanged, copy of the one from the DFHack r5 mini-mod pack.
	
	User/DFHack/Cart Filler:
		Fill minecarts with water or magma (without big fancy mechanical contraptions).
		To access fluids the filler needs to be adjacent to a downward passable tile leading to a cistern!
		To fill a minecart it needs to be adjacent to or on top of the filler, the minecart is NOT used as a reagent in the reaction. The best way to use the filler is to place it next to a minecart stop (or stops) and a cistern, this allows you to fill minecarts and send them on their way with a minimum of fuss.
		The amount of liquid used is determined by the minecart capacity, but at least 4/7 fluid is needed for the filler to function at all.
		Uses the pump operator skill (to fill carts) and mechanics (to build the filler).
		Requires the "Libs/DFHack/Fluids" default Rubble addon.
	
	User/DFHack/Cart Filler/Powered:
		An automatic mechanically powered version of the "DFHack/Cart Filler" addon.
		This addon works exactly the same with these additional rules (which are evaluated in order):
			If both magma and a magma safe minecart are present magma is loaded
			If any minecart and water is present water is loaded
			If magma is present and the only minecart(s) are not magma safe an announcement is printed
			If no minecarts are present nothing happens
			If no fluids are present an announcement is printed
		(If these rules cause announcement spam they may need to be changed)
		Requires the "Libs/DFHack/Fluids" default Rubble addon.
	
	User/DFHack/Dragon Engine:
		A direct port of the dragon engines from the DFHack r5 mini-mod pack.
		Dragon engines are basically machines that take magma and spit it out as dragon fire, very deadly short range siege engines.
	
	User/DFHack/Magma/Extractor:
		Extract minerals from magma, get all kinds of rare and valuable metals! (but only rarely)
		A magma extractor can be very profitable, provided you can feed it enough magma.
		Each cycle uses 2/7 magma and at least 4/7 magma needs to be available.
		Uses the furnace operator skill.
		Requires the "Libs/DFHack/Fluids" default Rubble addon.
	
	User/DFHack/Magma/Extractor/Powered:
		An automatic powered version of the "DFHack/Magma/Extractor" addon.
		This addon works exactly the same as "DFHack/Magma/Extractor", just without dwarf labor.
		Requires the "Libs/DFHack/Fluids" default Rubble addon.
	
	User/DFHack/Magma/Melter:
		Melt rocks into magma! 
		The magma is put into an adjacent magma-safe minecart or an adjacent cistern.
		When using minecarts this works exactly like the magma cart filler, when using a cistern it will attempt to spawn magma into any open space below an adjacent downward passable tile.
		DO NOT BUILD THIS BUILDING NEXT TO YOUR MAIN STAIRWELL!
		It is possible to enable magma workshops/furnaces here as well, so if you don't have a volcano there is no need to dig to the magma sea if you don't wish to.
		Uses the furnace operator skill.
		Requires the "Libs/DFHack/Fluids" default Rubble addon.
	
	User/DFHack/Obsidian Caster:
		The obsidian caster makes obsidian boulders directly from adjacent water and magma, which needs to be on the level below and accessible via downward passable tiles.
		4/7 water and magma are used for every boulder (so clever engineers can still be more efficient).
		Uses the pump operator skill (to operate) and mechanics (to build the caster).
		Requires the "Libs/DFHack/Fluids" default Rubble addon.
	
	User/DFHack/Obsidian Caster/Powered:
		An automatic powered version of the "DFHack/Obsidian Caster" addon.
		This addon works exactly the same as "DFHack/Obsidian Caster", just without dwarf labor.
		Requires the "Libs/DFHack/Fluids" default Rubble addon.
	
	User/Diggable Slade:
		Allows you to dig slade, once you get past it's guardians that is...
		Digging slade may lead to bugs, so be careful.
	
	User/Edible Vermin:
		Allows you to cook vermin at the kitchen to get cheap meat, a good emergency food source.
		(urist likes vermin for their great taste.)
	
	User/Factory:
		Adds a large factory building for making batches of cheap items (mostly from wood and stone).
		Allows you to make all kinds of furniture as well as sets of clothing, batches of weapons, and sets of armor.
		All factory reactions are just a little bit cheaper than the vanilla version.
		All factory items that can be made from stone and wood also have a version that can be made from blocks or planks, save even more resources!
		Uses various skills.
	
	User/Fix/Butcher Inorganic:
		Allows you to butcher creatures made from (some) inorganics, FB stone walls anyone?
		(Not 100% sure this works, copied from Modest Mod)
	
	User/Fix/Fish Pops:
		Raises fish vermin populations to the max.
		This should stop you from fishing out a lake or ocean ever again! (unless you catch ~30,000 fish...)
	
	User/Fix/Human Diplomats:
		Adds trade representatives and diplomats to the humans.
	
	User/Fix/Undead Melt:
		Fixes bad temperature values in material_template_default.
		Replaces material_template_default with a new copy, but the only changes are the temperature values.
		In any case this should make it possible to melt undead with magma.
	
	User/Fix/Usable Chitin:
		Allows you to use chitin as shell.
	
	User/Fix/Usable Feathers:
		Allows you to use feathers as "soft" pearls (you can use them to make crafts).
	
	User/Fix/Usable Scale:
		Allows you to use scale like leather.
	
	User/Fix/Vermin Variations:
		This adds two fixes:
			Forces all giant animals to have a life span of at least 10 years and
			Replaces all pool biome tags in giant and animal-man creatures with the correct lake biome.
		This is useful to make variations on vermin appear in game and to help them live long enough to be interesting.
	
	User/Fix/Whips:
		Nerfs whips and scourges so they are not deadly anti-armor weapons.
		
	User/Generic Animal Mats:
		This addon attempts to make all animal materials generic.
		It works fine with the standard addons, but it may not work with all third party addons!
		AFAIK all the standard addons are compatible with this one.
	
	User/Glass Forge:
		The glass forge allows you to make armor and weapons from all three kinds of glass.
		Green glass is a little under iron grade, clear glass is equivalent to iron, and crystal glass is a little under steel.
		If you melt down glass items you will get "bars" of glass, you can melt these bars back to raw glass if you wish, or you can use them at the forge to make other glass items. If you wish you can also make bars directly from raw glass.
	
	User/Lamination:
		Glue wooden planks together to make laminated logs, which are more valuable than most generic woods.
		Uses the alchemy skill.
	
	User/Pottery:
		Replace the vanilla pottery system with a MUCH expanded version. 
		Glazing is done via powders ground at a millstone or quern from ash. For more glaze types see the "User/Pottery/Glaze" addon.
	
	User/Pottery/Batch:
		Adds a batch kiln that can make 5x the stuff with 4x the material and 1x the fuel.
		Like the normal kiln addon you can make many, many, more things from clay than you can make in vanilla.
		Glazing and collecting clay can only be done at the normal kiln.
	
	User/Pottery/Glaze:
		Adds the glaze compounder, a workshop for mixing glazes from other color glazes.
		Also allows you to grind some minerals into glaze powder.
		Requires the "Pottery" addon.
	
	User/Quarry:
		Dig for some basic raw materials without needing to actually dig.
		There are two types of quarry:
			Quarry: 8 blocks, 4 mechanisms, 2 giant corkscrews, 10x10
				Dig for clay; 100% chance of a clay boulder, 10% chance of additional fire clay boulder
				Dig for sand; 100% chance for a bag of sand
				Dig for stone; 2% chance for 12 different layer stones
			Makeshift Quarry, free, 5x5
				Dig for clay; 20% chance of a clay boulder
				Dig for sand; 20% chance for a bag of sand
		You can use a quarry to dig for resources that do not exist in your current fort, so if you use quarries you never need to worry if you proposed embark has sand or clay.
		Uses the pump operator skill.
	
	User/Research:
		Build a laboratory and research advanced workshops.
		Handles the research class "RESEARCH_GENERIC", which is used by many of the default addons.
		You may also copy any research discovery here (for free). Uses the alchemy skill.
		Requires "Libs/Research".
	
	User/Sawmill:
		Farmable wood plants and advanced wood processing, makes aboveground and/or treeless forts a lot easier.
		If you wish you can scrap most wooden items (at a small loss), good for getting rid of all that elf crap you took from the last caravan :)
		The sawmill uses the carpentry and wood burning skills for its reactions.
		Don't forget a serrated disk!
		There are two kinds of farmable wood plants:
			Name			Biome			Growth Time		Log Color	Notes
			lumber bamboo	non freezing	1/2 year		Brown		Edible seeds
			stalkshroom		caves lvl 1-2	1 year			Dark Gray	Heavy and slightly more valuable
	
	User/Saurian:
		A new race to play as, see the Saurian guide in "Rubble Addon Extras.txt" for more details.
		This addon only includes the saurian creature and entity as well as saurian sized weapons and some basic reactions. Nothing else is included.
		It is HIGHLY recommended to activate at least the "User/Quarry", "User/Glass Forge", "User/Sawmill", "User/Warcrafter" and "User/Pottery" addons.
		To dig as saurians you need to buy a pick and resize it at the craftsaurian.
		If this addon is activated dwarves are unplayable and the dwarven outpost liaison becomes a trade rep (unless you activate "Saurian/Unplayable").
		Requires the "Glass Forge" addon.
		Requires the "Libs/Castes" default Rubble addon.
	
	User/Saurian/Basic Castes:
		Adds some basic castes to the saurians.
		See the saurian guide for details.
		Requires the "User/Saurian" addon.
		Requires the "Libs/Castes" default Rubble addon.
	
	User/Saurian/Basic Castes/DFHack/Transform:
		Allows you to change the cast of your saurians at the new meditation circle workshop.
		Meditating takes from 1/2 to 3 seasons (depending on what the result caste is).
		Requires the "User/Saurian" and "User/Saurian/Basic Castes" addons.
		Requires the "Libs/Castes" and "Libs/Castes/DFHack/Transform" default Rubble addons.
	
	User/Saurian/Unplayable:
		If you just want saurians as an extra race to trade with the dwarves, activate this addon.
		Requires the "User/Saurian" addon.
	
	User/Tanning:
		Makes creatures give skin "globs" instead of "skins" when butchered, this allows large animals to produce more than one piece of leather once tanned. 
		If "User/Fix/Usable Scale" is activated then scale is changed as well.
		Based on a mod/modders resource by sackhead.
	
	User/Tileset/MLC/DFHack/TWBT:
		A simple addon that installs my custom ASCII-like tileset.
		This version uses the "Text Will Be Text" DFHack plugin (you will need to install that plugin separately).
		Unlike most ASCII tilesets this one has special tiles for most buildings and items, making it (almost) fully graphical.
		Comes only in the 16x16 size.
		This addon makes no raw changes, so it should be usable with all ASCII compatible mods.
	
	User/Tileset/MLC/Normal:
		A simple addon that installs my custom ASCII-like tileset.
		This version is very close to default ASCII (a few tiles are graphical).
		Comes in two sizes, a 16x16 for fullscreen mode and a 10x10 for windowed mode.
		This addon makes no raw changes, so it should be usable with all ASCII compatible mods.
	
	User/Tileset/MLC/Pure ASCII:
		A simple addon that installs my custom ASCII-like tileset.
		This version has absolutely no graphical tiles.
		Comes in two sizes, a 16x16 for fullscreen mode and a 10x10 for windowed mode.
		This addon makes no raw changes, so it should be usable with all ASCII compatible mods.
	
	User/Tileset/MLC/Tracks:
		A simple addon that installs my custom ASCII-like tileset.
		This version has really nice track graphics. There are a bunch of things that will look a little funny, but tracks look great (tile layout was sub-optimal even for 34.11, now with multi-tile trees it is rather bad).
		Comes in two sizes, a 16x16 for fullscreen mode and a 10x10 for windowed mode.
		This addon makes no raw changes, so it should be usable with all ASCII compatible mods.
	
	User/Tileset/Vanilla:
		Restores the vanilla tileset related init options.
		It is a very good idea to use this addon whenever you are using the vanilla tileset, as otherwise prepping worlds generated from those raws may not work quite right.
		This addon makes no raw changes, so it should be usable with all ASCII compatible mods.
	
	User/Tileset/Vanilla/Graphics:
		The vanilla creature graphics, mostly just a demo.
	
	User/Training:
		Make ink from ash, and paper from wood, then make a book from paper and ink and read it for exp in your chosen skill.
		Note that EVERY skill in the game can be trained this way, even the useless ones ;). Making books uses the Alchemy skill.
		Note that skills that do not have a labor appear to allow ANYONE to do the reaction, use workshop profiles.
	
	User/Warcrafter:
		A special workshop for all your early game armor and ammo needs, the warcrafter allows you to make:
			Scalemail armor from leather, bone, and shell
			Bolts and arrows from bone and shell
			Heavy cloaks from leather
			Green glass weapons and armor (if the "User/Glass Forge" addon is active)
		Making bone and shell items requires you to cut them into "scales" first.
	
	User/Warcrafter/Adventure:
		Adds a bunch of adventure mode reactions, mostly for making use of bones and skin (via bone and leather items, mostly armor but also a few other useful things including a few bone weapons)
		This addon is mostly for those who do not wish to use the advfort DFHack script, but still want to do some basic crafting.
		Supports "User/Tanning" and "User/Fix/Usable Scale", but a DF bug keeps the support for "User/Tanning" from working correctly, it works but it requires lots of hunting to get much leather (the tanning reaction always takes two stacks of skin, and produces one tanned hide no matter how large the stacks are, so "User/Tanning" actually REDUCES the hide output in adventure mode).
		Requires "User/Warcrafter" (for the armor items and materials).
	
	User/Wicker:
		Make all sorts of stuff from woven plants. Uses the process plants skill (for bundling) and weaving (for making stuff).
	
	User/Zap Aquifers:
		Disables all AQUIFER tags, very simple.

==============================================
BUGS (and other issues):
==============================================

Known Issues (non-bug):
	
	You should not use the extension .txt for readme files, as this will cause Rubble to parse those files.
		See "Rubble Addons.txt" for more details.
	
	It may be possible to see FBs made from "caste transformation", there is nothing I can do about this as autoSyndrome requires the use of inorganic materials.
	
	The scalemail armor from the "User/Warcrafter" addon is always marked "Foreign" in the equipment screen, you will have to explicitly assign it.
		IMHO this is better than sometimes not having all the armor pieces (which is what happened before).

Bugs:
	
	Several vanilla bugs make themselves known:
		Adventure mode reactions do not always work the same as fortress mode reactions:
			"User/Warcrafter/Adventure" does not work quite right with "User/Tanning" because the adventure mode tanning reaction uses too much skin. This cannot be fixed by me.
		
		Reactions produce unusable gloves:
			The glove reactions are left in only because DFHack can be used to fix this bug, users that lack DFHack (or who have not run the proper script) should not make any gloves at custom workshops.

	There are no known (non-vanilla) bugs.

When making an error report please post the FULL log file! Posting just a few lines tells me almost nothing about the problem, after all I made Rubble write all that stuff for a reason :)

If any of the documentation is not 100% clear just ask. 
I know everything there is to know about how Rubble works and so I tend to forget to put 
"obvious" stuff in the docs. A reminder that I forgot something is more a help than a hindrance.

In the event that I cannot be contacted on the Bay12 forums (user name "milo christiansen"), my email address is:
	milo.christiansen (at) gmail (dot) com
Please wait 1-4 weeks before giving up hope, as my Internet access is VERY irregular (and I check my email even less often).

==============================================
Changelog:
==============================================
v4.2 (for DF 40.5)
	Updated the "Base" addon to 40.5.
	Removed "User/Nerf/Ranged" (it's changes are now part of vanilla!).
	Moved "User/Nerf/Whips" to "User/Fix/Whips".
	Added "User/Generic Animal Mats", generic animal materials for all!
	Tweaked the DFHack support a little.
	Removed "Libs/DFHack/Spawn", as 40.x DFHack will come with a unit spawn script there is no need for it.
	Added "User/Cat Damper", slightly less annoying cats.
	Added "User/Fix/Human Diplomats", adds trade reps and diplomats to the humans.

v4.1 (for DF 40.3)
	Added "Libs/Research", library support for a research reaction system.
	Moved all "user" addons to the "User" group.
	Added a bunch of addons from Better Dorfs (all were put into the "User" group).
	Added "User/Warcrafter/Adventure", non-DFHack adventure mode crafting.
	Reworked the DFHack support templates in "Libs/Base", they should be all ready to go now.
	Removed "Libs/DFHack/Command".
		Use DFHACK_REACTION or DFHACK_REACTION_BIND (from "Libs/Base") instead, they have radically
		different syntax and usage, but they use the new reaction binding script (which is much
		better than the old autoSyndrome method).
	The REACTION template did not work exactly as documented, now it does.
		This caused a minor loss of functionality (which wasn't used anywhere anyway) for a great
		gain in simplicity.
	Updated the tileset templates to handle the new tree related init settings.
	Removed support for "Libs/DFHack/Command" from "Libs/DFHack/Fluids".
	In the "Base" addon the dwarf creature now has the caste templates for "Libs/Castes" support.
	"Libs/Castes" is now forced active during normal generation cycles.
	Updated "Libs/DFHack/Announcement" to use the new DFHack system.
	Updated "Libs/DFHack/Spawn" to use the new DFHack system.
	Added "User/Tilesets/MLC/Pure ASCII", a version of my tileset that has no graphical tiles.
	Removed the rubble:activeaddons script variable, it was only used once, and that was easy to replace.
	Added script command rubble:add_file, add a file after generation starts (experts only!).
		This is mostly to make patches possible, most users will have no use for this command.

v4.0 (for DF 40.3)
	First public release for the 4.x series.
	Replaced "Base" with a base made from the raws for DF 40.3.
	While the DFHack library addons were not removed do not assume that they will work with DFHack for DF2014.
		In particular anything that needs "Libs/DFHack/Command" probably won't work.
		"Libs/Caste/DFHack/Transform" will probably work just fine as will "Libs/DFHack/Add Reaction".

v4pre (for DF 34.11, changes from Rubble 3.13)
	Changed scripting subsystems from Raptor to Rex.
		Rex is a little less flexible but ultimately faster and more powerful scripting language loosely based on Raptor.
		Many, many things have changed, but only a few of them really make a big difference:
			Some commands like set, var, command, and namespace are replaced by new syntax constructs.
			Code is "compiled" to a much lower form, allowing for faster execution.
			A lot of the more "dynamic" stuff from Raptor is gone:
				It is not possible to use a string as code (directly, eval still works).
				Some commands require code blocks to have variables "predeclared" via a special block declaration syntax.
					This replaces many uses of the old params array with a much clearer, less error-prone syntax.
				The various advanced eval commands are gone.
			You may now use single quotes (') and back quotes (`) as string delimiters.
		See the files in the "Rex Docs" folder for more info
	As a direct result of the change in scripting subsystems Rubble is now MUCH faster.
	More internal data is exported to scripts (via GenII, so be careful!).
	Some native commands have been ported to init scripts.
	Fixed lots of little scripting bugs that were discovered when porting.
	Removed the builtin templates !SCRIPT, SCRIPT, and #SCRIPT.
	Added a new base template: @SCRIPT.
	The !SCRIPT_TEMPLATE builtin template was removed, use the rubble:template script command.
	The syntax of the rubble:template command has changed significantly (it's simpler in some ways).
	Removed shell mode, it was useful, but ultimately added much complexity for small gain.
	Removed the lexer test. This was an internal debugging utility that has outlived it's usefulness.
	Removed the "Shared Animal Mats" addon, as it is too fragile for inclusion with Rubble.
	Added "Fix/Fish Pops", raises all fish population to the max, never fish out an area again.
	Fixed rare case where the addon list would not update.
	The script raw parser has moved from the "raw" module to the "df:raw" module.
		The parser has been simplified and is joined by a slow, but more flexible cousin.
	Added "Dev/Tileset/Export", automatically extract tile info from raws that have the Rubble tileset templates.
	Added "Dev/Tileset/Insert", automatically inserts Rubble tileset templates into untemplated raws.
	Removed the #WRITE_TILESET template, it is now obsolete.
	Removed support for anonymous template parameters (eg %1) as this was never useful, and has become less so over time.
	The "outputdir" setting has changed, it is now the path to the "raw" directory, not "raw/objects".
	Renamed "Tilesets/MLC" to "Tilesets/MLC/Normal".
	Added "Tilesets/MLC/Tracks", A version of the MLC tileset that has graphical tracks.
	Added "Tilesets/MLC/DFHack/TWBT", A version of the MLC tileset for use with Text Will Be Text.
	The addon loader no longer writes every addon and file name to the log, all this did was bloat the log for no good reason.
	It is now possible to specify -config or -addons more than once, this applies to rubble.ini as well.
		Each occurrence is parsed separately, but the results are added together.
	Added the ability to pass a file name to -config, in that case the config variables will be read from the file.
	Added the ability to pass a file name to -addons, in that case the file is treated as an addon list file.
	Removed live profiling, it is normally really nice to have, but Rubble completes too fast for it to be usable.
	Made sure Rubble always exits with code 1 when it does not complete properly.
	Made sure Rubble prints "Done." as the last thing when it succeeds.
	Finally got it so that ALL native script command documentation is automatically generated.
		The script documents should never go out of sync now (user command docs are still done by hand though).
	Rubble now prefixes all parsed files with a short message stating that the file was automatically generated.
		This message also includes the source path (uses the AXIS syntax).
	Added new addon, "Building Dimensions": Append all building names with their dimensions.
	Lists are now always split at ';' characters instead of using the OS-specific path list separator char.
		This change should make things more consistent across OSes.
	rubble.ini now uses the same rules as the command line flag handler for parsing boolean values.
	Added an "installer" mode, it simply loads files from a single zip and then runs any ".inst.rex" scripts,
		use to install non-raw mods and addons automatically.
	Created a new GUI with support for prepping or regenerating a region (still windows only).
	The GUI now launches Rubble in installer mode when it is started with a command line argument.
		(this makes it support drag and drop installing)
	The GUI now supports a "minimal=true" option in it's INI that makes it not create the "Prep" and "Regen" tabs.
	Made all addons deterministic, this should make region regeneration possible.
		(in this case deterministic means that generated raw objects will always be in the same order for a given input)
		Changed addons:
			"Libs/Castes"
			"Libs/Castes/DFHack/Transform"
			"Libs/Crates"
			"Libs/DFHack/Command"
	Added two new templates: SHARED_OBJECT_KILL_TAG and SHARED_OBJECT_REPLACE_TAG, modify shared objects with minimum fuss.
	Rewrote "Fix/Undead Melt" to use SHARED_OBJECT_REPLACE_TAG, it is now much less brittle.
	All file access is now negotiated through AXIS VFS, this makes Rubble's directory structure MUCH more flexible.
		It should no longer be possible to write to arbitrary locations on the file system, so security should be much improved.
	The "file" and "fileio" script modules are no longer loaded, use the commands in the "axis" module instead.
	The various script directory variables have been removed as they are not needed with AXIS VFS.
	Added new base template: @IF_SKIP, conditionally skip a file with minimal fuss.
	Added addon "Tileset/Vanilla/Graphics": mostly just a demo of how to make creature graphic addons.
	Added four new tileset templates to "Base/Templates", these templates are for finer grain control of font settings.
	Added new addons, "Libs/Macro" and "Libs/Macro/Example": Generate macros from scripts.
	Added native script command rubble:activate_addon, force an addon active at any point.
	Renamed "Base/Templates" to "Libs/Base".
	"Libs/Base" is now forced active by init script during normal generation cycles.
	Removed "Base/Clear", replaced with an init script.
	Replaced the "Libs/Base" DFHack support with a much more flexible and less error prone system, no more extra install steps!
	Updated "Libs/DFHack/Add Reaction" to use the new DFHack support.
		There is no longer a need to prep between regions when using this addon.
	Moved "Base/Files" to "Base"
	Updated "Libs/Caste/DFHack/Transform" to use the newest DFHack stuff.
	Updated a LOT of documentation, and still more needs to be written...
