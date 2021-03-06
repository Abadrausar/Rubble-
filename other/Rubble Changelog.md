
==============================================
Changelog:
==============================================
v5.7 (for DF 40.18)
	Updated to 40.18 (no actual changes needed).
	Completely rewrote the tileset system, the new system is MUCH better (and MUCH easier to use).
		The old tileset templates are gone, replaced with a new parse step and some special scripts, see "Rubble Basics".
		Note that even templates that would still be relevant had to be replaced with script commands for technical reasons.
	Rewrote some of the standard templates to remove obsolete tileset support (the specialized SHARED_OBJECT templates).
	Remade the "User/Tilesets/Phoebus" addon to use the new tileset system.
	Added "User/Tileset/ASCII", use to switch back to ASCII mappings from another tileset.
	Rubble now writes tileset information every time you generate, look for "out:current.tset".
	Removed all addons in the "Dev/Tileset" group (as they are obsolete).
	Stripped all calls to #TILE and #COLOR from the "Base" addon.
	Added a new operation mode to the interfaces to allow changing tilesets on worlds in progress.
	Removed prep mode and all support for prep scripts (this mode is obsolete with the new tileset system).
	Loads of little internal changes to the Rubble engine to support the above changes.
	Fixed possible bug with the @GRAPHICS_FILE template.
	Renamed the #INSTALL_GRAPHICS_FILE template to @INSTALL_GRAPHICS_FILE.
	Fixed the "User/Tileset/Vanilla/Graphics" addon to fit the changed graphics template.
	Added new base template: @INSTALL_IMAGES_AS_GRAPHICS, runs @INSTALL_GRAPHICS_FILE for all PNG images in an addon.
	The "Skip" file tag is now cleared for all files when addons are activated.
		This should not make any difference, but just in case.
	Versions prior to this one are marked incompatible (what with no more prep support and the tileset templates being gone).
		This will only effect addons that run a Rubble version check (which is none of them AFAIK).
	The web UI now times generation, these times only cover the non-interactive portions.
	Stopped various things from trashing the default addon list file, I hope I didn't break anything...
	It is now possible to set "_RUBBLE_NO_CLEAR_" to "true" to suppress clearing the output directory (experts only!)
	Rubble now handles the files in "out:objects/text", see the "@TEXT_FILE" template.
	Updated the "Base" addon to include the text files.
	The clear script now does a more complete job of nuking the contents of "out:".

v5.6 (for DF 40.15)
	"Dev/Tileset/Insert" now adds tileset templates to all creatures.
	Modified "Base" to have tileset templates in all creatures.
	The "User/Tileset/MLC/DFHack/TWBT" addon now has better graphical tracks (including individual track ramps!).
	Slightly revised all the tilesheets in the "User/Tileset/MLC" addon group.
	Added "User/Tilesets/Phoebus", support for the Phoebus tileset (and any other tileset that uses the same mapping).
		This addon does NOT include the tileset images, only the init and raw mappings.
		FYI this addon took less than 10 minutes to make (I love the "Dev/Tileset" addons).
	Any old addon activation state information is now cleared before activating addons.
		It should now be possible to *deactivate* addons in the web UI (some weird issues make sense now...).
	Fixed the many, many, bugs in "User/DFHack/Powered/Logic Gates", it works now (and it's really cool!).
	Added "User/DFHack/Powered/Item Sensor", detect items and take action based on their count.
		Finally, the ability to automatically stop powered workshops from over-producing (when used with the new controller).
	The "User/DFHack/Powered" addon now includes a new basic building, a "powered workshop controller".
		The controller allows you to "shut down" a powered workshop without cutting it's power.
	Modified all powered workshops to use the new controller building.
	Items output from powered workshops will now drop if they come out over open space or a down slope.
	Fixed the require power option in "User/DFHack/Powered/Sorter", I only had half of what I needed.
	It is now possible to "invert" powered sorter settings (take everything EXCEPT what is selected).
	Increased the power output of "User/DFHack/Stirling Engine" (to 6), 5 was not enough to power anything.
		The smallest workshops require 5 power, but you need a axle between the engine and the workshop, so 6 is the minimum.
	Fixed a major bug in "User/DFHack/Powered/Cart Launcher" (improper vehicle lookup).
	The cart launcher can now optionally forbid a cart when it launches it (reduces road kill potential).
	The cart launcher should now be capable of launching off track stops and in other adverse situations.
	Allowed the cart loader and cart launcher to use outputs to restrict which carts they should effect.
	Added a bunch more stuff (and a few diagrams) to the "User/DFHack/Powered" addon documentation.
	Added a new "incompatible" key to addon.meta, this key allows specifying addon incompatibilities.
	Addon information pages (both web UI and static) now list incompatibilities right after dependencies.
	Added a new "format" key to addon.meta, set to "html" to allow "header" and "description" to be HTML formated.
	Addon information pages (both web UI and static) now handle the new "format" addon.meta key.
	Removed most of the addon dependency management commands from the standard library, use addon.meta!
	Made the rubble:activateaddon standard user command recursive (it now activates dependencies as well).
		This command is still available (against my better judgment) because activating addons is easy to get wrong,
		and there are cases where it needs to be done from a load script, but use addon.meta if at all possible!
	The addon activation stage will now print warnings when a non-library addon is automatically activated.
		This is so users can see why they have addons active that they did not activate.
	Nonexistent templates are now errors again.
	Updated Rex to version 1.7, this fixes a minor GenII bug that one of this version's changes triggered.
		Several other bugs (that did not effect Rubble) were fixed, plus the sort:array, int:++ and int:-- commands were added.
	DFHack onload.init reaction and command output is now sorted (AFAIK this is cosmetic only).
	You can now clear your selection in the web UI on the "Select Addons" page during normal generation.
	The glass "metal" names in the "User/Glass Forge" addon are now prefixed with "tempered" ("tempered green glass", etc).
	It is now possible for the web UI to load user defaults for configuration variables from "rubble:userconfig.ini".
		This uses the same format as "out:genconfig.ini", so it should be easy to setup.
		Obviously this only takes effect if you choose "Edit Configuration Variables" after selecting addons.
	Added "User/Speluncaphobia", makes selected creatures get dizzy and slow after too long underground.
		By default this effects saurians and humans (this is changeable via config vars).
	Added "User/Speluncaphobia/DFHack", allows suppressing the speluncaphobia effect with a torch "shield".
		Which entities can make torches is controlled via a config var, by default it is saurians and humans.
		Currently non-working due to (I think) a DFHack bug (see the addon description).
	Saurians are now better diggers (from 1% to 25% learn rate), there is now a better way to keep you above ground.
	Saurians may now make coke from coal (this was an oversight).
	Added "User/Shovel", craft shovels from wood or bone that break after a while.
	The Rubble DFHack psudo-module function "rubble.require" now generates an error if you try to load a non-existent module.
	Added three new base templates to deal with entity nobles (#ENTITY_NOBLES, ENTITY_ADD_NOBLE, and ENTITY_REPLACE_NOBLES).
	Added #ENTITY_NOBLES to all entities in "Base" except kobolds and the layer-linked animal people.
	Added #ENTITY_NOBLES to "User/Saurian".
	Added "User/Human/Nobles", static nobles for the human entity.
	Changed "User/Fix/Human Diplomats" to use the new templates.
	Added "User/Human/Playable" rudimentary playable humans.

v5.5 (for DF 40.15)
	Updated to 40.15 (no actual changes needed).
	Changed the "User/DFHack/Powered/Kiln" addon to work like the other powered production buildings.
		You can now fire clay directly to items (this change *just* missed the last version).
	Added the following new powered workshops (each in it's own addon):
		Spinner, process plants to thread (outputs seeds if needed).
		Loom, weave thread to cloth.
		Mill, grind plants to dye powder or flour (outputs seeds if needed).
	Added "User/DFHack/Powered/Logic Gates", an experimental powered logic gate addon (untested!).
		Once this is perfected I plan to add many gates and sensors to make complex (automatic!) world interaction possible.
	Fixed an issue with glass item production in "User/DFHack/Powered/Factory".
		If you ran out of raw glass you could end up with a bunch of unusable raw glass items piling up at the input, ouch.
	Fixed an issue with powered sorters sorting by material (like an idiot I had only tested sorting by item).
	The powered sorter now allows you to set an "output item limit".
		If there are more items on the sorter's output tile(s) than the limit, no more will be output until some are used.
		This should help stop "bag starvation" and other problems caused by over-eager sorters.
	Added an option to make powered sorters require power.
		Making them require power by default was deemed to be too annoying, but if you want things to be a little harder...
	Fixed several other small issues in various powered workshops (none that could cause real problems).
	Rebalanced all powered workshop time and power requirements, in general they are slower and take more power.
	Extended the powered workshop support API a little ("User/DFHack/Powered").
	Finally gave up and based my persistent data APIs on the default "dfhack.persistent", things should stay in sync now.
	Fixed a major bug with automatic addon activation, this is now carried out recursively.
		If an addon activated an addon that was supposed to activate a third addon, the third addon was never activated, oops.
	Fixed "User/Metallurgy/Bloomery", I had forgotten to register the buildings with an addon hook.
	Extended "Dev/Cheat" (again), more stuff to test, more raw materials required.
		I also reordered the reactions to make them "flow" better.
	"addon.meta" files are now treated as being "parseable" for the purpose of determining if a directory is an addon.
	A "generation report" is now written to the output directory when generation is finished.
		This report is basically a static version of the web UI documentation, but only for addons, etc, that were active.
	The "clear script" now removes old generation reports.
	Fixed some more file permission errors, stupid Linux.
	There are two erroneous entries in the previous version's changelog.
		I listed two changes to "User/DFHack/Powered/Cart Loader" and "User/DFHack/Powered/Cart Launcher" that were incorrect.
			Neither of these addons used "Libs/DFHack/Upgrade Building" (directly, it was a dependency of a dependency).
			The cart launcher supports a lot more than just four launch thresholds.
		(These addons were rewritten two or three times over the course of getting 5.4 ready, and I forgot to fix the changelog)

v5.4 (for DF 40.14)
	Updated "Base" to DF 40.14.
	Added stepladders and removed LIKES_FIGHTING from "User/Saurian".
	Added "Libs/DFHack/Upgrade Building", allows you to change the type of already constructed workshops.
	Changed "User/DFHack/Cart Loader" and "User/DFHack/Cart Launcher" to use "Libs/DFHack/Upgrade Building".
	Changed "User/DFHack/Cart Launcher" to support 100%, 75%, 50%, and 25% launch thresholds.
	It is now possible to disable reactions in "User/Factory" by category (all wood reactions, all stone reactions, etc.).
		This is done via configuration variables, if you use the web UI just click "Edit Config Variables" when generating.
	Removed all metal working reactions from "User/Factory".
	Added reactions to make tools to "User/Factory".
	Removed block creation reaction from "User/Factory".
	Added a BUNCH of new powered workshops and moved the existing ones to the "User/DFHack/Powered" addon group.
		If you activate all the powered workshops you can make complicated automatic manufacturing setups, have fun!
		Read the addon descriptions! They have important information about how to use the machines.
		In particular the description for the "User/DFHack/Powered" addon has an overview of sorts.
		The existing powered workshops have been changed to better integrate with the new ones.
	Added "User/DFHack/Prospect", allows your miners to examine the area looking for ores and gems.
	Fixed a Rex bug that stopped sort:map type indexables from working sometimes, this did not effect any existing addons.
	Extended "Dev/Cheat" some (obsidian mechanisms and some metal ores).
	Reworked the persistence support in "Libs/DFHack/Timeout" to better support script system reloads.
		Now creatures can't get "stuck" when you reload the scripting system (not sure if they ever did, but better safe...).

v5.3 (for DF 40.13)
	Updated Rex (v1.6).
		array indexables should be a little faster in most cases.
		Some errors (specifically those dealing with indexables) should now be properly handled.
		A pseudo-random number generator is now included (the "rand" module).
		debug:value now prints the raw value data with a much nicer (clearer) format.
		New (sorta, Raptor had breakif) commands breakif and breakloopif (short versions of if + break or breakloop).
		New command modval, enables some interesting pointer-like behavior (for advanced users only!).
	Script command documentation is now generated for ALL included Rex command packages (I had forgotten a few...).
	Added "User/Deadly Bones", makes creatures extra vulnerable vs. their own bone material (10x damage bonus!).
		Not 100% sure this works, if you want to test effectiveness some ideas are listed in the addon description.
	The "shell" and "bone" materials in "User/Warcrafter" now have the same colors as vanilla bone and shell.
	Saurians may now make socks (because it's just not DF without piles of XX(rope reed fiber sock)XX).
	Added "User/DFHack/Pop Cap", set the population cap(s) on a per-world basis (overriding the init settings).
	Added "User/DFHack/Cart Loader", load items into minecarts without dwarf power.
	Added "User/DFHack/Cart Launcher", launch carts that are more than 75% full (automatically!).
	Added "User/DFHack/Stirling Engine", generate a (very) small amount of power for free.
	Added "User/DFHack/Bulletin Board", a bulletin board, where your idle dwarves will post thoughts, preferences, etc...
	The web UI now passes the full URL of the menu page (including the "http://" part) to the browser start script.
	Fixed some file permission bugs (writing files with the wrong mode), this should only effect Linux users.
	Added default browser startup scripts for the web UI (for Windows and Linux).
	Added a check to the Web UI server to make sure startup errors for the HTTP server are reported (yes, I am an idiot).
	Changed the web UI default port to 2120 (from 1010, which Linux doesn't like).
		The port was picked arbitrarily, if it collides with something common let me know and I'll change it again.
	Added missing dependency in the "Dwarf/Castes/DFHack/Transform" addon.meta file (so it should work now).
	Replaced the old (way out of date) "Rubble Addon List" file with a new automatically generated one.
	Changed the DFHack Lua pseudo module syntax so it matches the syntax for regular modules.
		(just use rubble.require and rubble.mkmodule)
	Fixed an AXIS bug that made the web UI crash if the "Go Directly To Regeneration" button was clicked.
		This bug effected other things as well, but with a default setup this was the only thing that triggered it.
	Revised some of the templates in "Libs/Base" (mostly removing unneeded block declarations).
	All templates that had a command version for use in prep mode now have a command version for normal use as well.
	Changed the default tileset addons to use a different system for prep support.
		They generate their own prep scripts now, so the whole process is less prone to error.

v5.2 (for DF 40.13)
	Changed "User/Saurian" entity values a little, they should now dream about things other than becoming a legendary warrior.
		Saurians are still properly warlike, they just have more other interests now.
	Fixed Reaction permissions for "User/Saurian", saurians should now be able to brew alcohol and process plants to bags.
	Fixed the saurian creature's body (it had a few "issues"), this should fix saurians having 5 pairs of gloves and shoes.
	Added a few more scale color options for saurians.
	Made a much needed update to the Notepad++ UDL file.
	Reordered the entries in the "Rubble Base Templates" documentation file to group templates by theme.
		I also made minor corrections and edits to most (if not all) of the other documentation.
	Updated to the latest Rex version.
		This version brings two major changes:
			Better variadic command/block parameters (to keep old code going replace "..." with "params=...").
			The ability to specify optional command parameters out of order (by name).
		And several relatively minor ones:
			Failed local variable lookups will now fail over to looking the variable up in the global module.
			ALL global data is now stored in the global module.
			New bool:sand and bool:sor short-circuit evaluation boolean logic commands.
			New "struct" type indexable, basically a map that cannot have new keys added that is created from a template.
		See "Rex Docs/Rex Basics" for details.
	Added the missing "thread" script commands, I could have sworn I had already included them...
	Fixed all script code to use the new variadic parameter syntax.
	Fixed a bug that caused "User/DFHack/Steam" to not work when a research addon was activated.
	Fixed several bugs in "Libs/Research", I *thought* it worked. I thought wrong.
	Removed the old (depreciated) GUI, use the web UI.
	Errors are now tied to a position (which may be WAY off, so user beware).
		In general error positions should be more-or-less correct, but sadly unless I rewrite large parts of Rubble
		they will never be as good as the positions reported by Rex (Rex errors are almost always spot on).
	Rubble now uses rex.Value and rex.Position internally rather than use it's custom (less flexible) versions.
	Removed the only remaining native template (!TEMPLATE), it is now a script template (backed by a native Rex command).
	Added rubble:usertemplate script command, allows you to make user templates (like what !TEMPLATE makes) from scripts.
	Removed support for native templates.
	The lexer now resets and continues when it finds unmatched curly brackets (this only applies to extra closing brackets).
	The stage parser now enforces correct template syntax ALL the time, not just when the template name is valid.
		All opening curly brackets MUST have a matching closing bracket (the old way could cause weird bugs).
		Opening brackets in single-quote character entities are not effected by this change (they are still ignored).
	The above two changes fix the "panic attack" Rubble would have if you had unmatched curly brackets.
	Rubble should now be able to load zipped addons from the Internet by URL (using .webload files), see "Rubble Basics".
		This is mostly untested (I don't have Internet), but it works fine with a simple HTTP server running on localhost.
	The "addons:zip:" AXIS DataSource is now created (or recreated!) just before addons are loaded.
		This allows zip file addons that were added after Rubble has started to be loaded when the other addons are.
		This change allows web load addons and makes the web UI more consistent.
	Updated to the latest version of AXIS VFS (major internal changes, but the script API is unchanged).
	Changed the syntax of the basic tileset templates a little (they now take a third parameter).
	Changed the "Dev/Tileset/Insert" and "Base" addons to match the tileset template changes.
		No other addons needed to be updated.
	All interfaces now have slightly more user-friendly handling of unrecovered errors.
		The error and program stack trace is now written to the log file (please upload your log if you get an error!).
		(It is my fondest wish that users never see this change in action!)
	Moved the shared file "wood_additions.rbl" to it's own addon "Libs/Wood Additions".
	All addons that used the "wood_additions.rbl" file now use the "Libs/Wood Additions" addon.
	Removed many addons from the addon list file, to get the real up-to-date list see the web UI addon list page.
		This file is mostly for posting to the forum, it is not updated except in the most general way.
	Both parts of the addon description in the web UI use the same font now.

v5.1 (for DF 40.13)
	Updated the base to DF 40.13
		40.12's change log lists certain changes that are listed in 40.13's change log as being made in 40.13.
		Did Toady jump the gun? Anyway the base is now up to date.
	Fixed "User/Saurian/DFHack/Transform", forgot to update the creature body for 40.x, oops.
	Added default values to the saurian entity.
		They are really big on martial prowess and value other things that seemed appropriate at least a little.
	Fixed a minor JS error that stopped the web UI's regeneration mode addon list page from working.
	When running a regeneration cycle in the web UI configuration variable default values are now read from genconfig.ini.
		This was the case before when you chose to go straight to generation, but now it works for when you
		choose to edit the variables manually as well.
	Fixed the bone shafts reaction in "User/Archery" to use the proper number of bones.
	Fixed the default "User/Saurian" bone item reactions to use the proper number of bones.
	Added "User/DFHack/Spatter" and "User/DFHack/Steam" from BD (based on DFHack example mods).
		Somehow I forgot to included these earlier, anyway here they are.
	Added "User/Metallurgy/Aluminum Bronze", a better bronze for those with extra aluminum.
	Fixed a bunch of outdated documentation.

v5.0 (for DF 40.12)
	Rubble has been restructured as a library.
		This change makes it possible (even easy) to make multiple different user interfaces.
		In addition many parts of Rubble were refactored and error handling is much improved (in some ways).
		Basic functionality should be (mostly) unaffected, some minor features may not work quite right yet,
		but all the standard addons generate without errors.
	New UI: A CLI front end.
		This is mostly like the old interface, you can even use it with the old windows GUI.
		Rubble comes with 32 bit Windows, Linux, and OSX versions of this UI.
	New UI: A simple web UI.
		This simple UI supports all normal operations, it's kinda ugly, but it gets the job done.
		Rubble comes with 32 bit Windows, Linux, and OSX versions of this UI.
	Updated the "Base" addon for DF 40.12 (no actual changes made).
		Despite what the changelog says my diff program could find no raw changes from 40.11 to 40.12, odd...
	Removed install mode, maybe later.
	Rubble can now load addons from more than one directory at once! 
		Duplicate addons will cause Rubble to abort, so don't do that.
	Rubble now simply removes nonexistent templates and prints a warning (as opposed to aborting).
	The old "After Blast comes Rubble." line in the log header is replaced with much cooler random lines.
		Samples:
			"Why do all of these end in exclamation points!"
			"Unintentionally Ironic!"
			"Why did I add this feature?"
			(as well as many more)
	Addon meta data support: Addons may now have a file named "addon.meta" that has information about the addon.
		This information includes descriptions and information about how the addon should be loaded.
	All default addons have complete addon.meta files:
		Addons will now automatically activate ALL their dependencies.
		Addons that do not directly effect the user (such as libraries) can no longer be directly activated.
			If you look at the description page for such an addon it will note that it is "An automatically managed library".
		Addons have at least some description.
		Any config variables that the user should know about are listed with appropriate defaults (the web UI uses this).
	Removed almost all the loader scripts that were made obsolete by addon.meta files (which was most of them).
	Removed a bunch of obsolete pre scripts (another casualty of addon.meta).
	Rex has been updated: It is possible to redeclare global data, and GenII now allows direct replacements for type index.
	Removed "User/Tilesets/MLC/Tracks", the track tile mapping looked really, really, bad with the new trees.
	Removed the GUI tutorial as the GUI will be removed shortly.

v4.7 (for DF 40.11)
	This is the first version of Rubble for DF 40.x that has all the DFHack stuff working!
		Or at least everything *should* work, some stuff has less testing than I would like...
	Updated the "Base" addon to DF 40.11
	Fixed a major Rex bug with by-value command calls, this should have only effected "Libs/Macro".
	Rewrote the "User/Pottery/Glaze" addon and made it no longer dependent on "User/Pottery".
	Vastly simplified glazing in "User/Pottery", the reaction is now more-or-less like vanilla.
	Added "Libs/DFHack/Timeout", adds support for persistent time delays.
	Rewrote "Libs/Castes/DFHack/Transform" so it no longer requires autoSyndrome (so it should work with the new DFHack).
		The API has changed significantly! Luckily the new system is vastly simpler than the old.
	Renamed "User/Saurian/Basic Castes/DFHack/Transform" and "User/Saurian/Basic Castes" (removed the "Basic" part).
	Rewrote "User/Saurian/Castes/DFHack/Transform" to use the new system.
	Added "User/Dwarf/Castes/DFHack/Transform", change your dwarves castes for a nominal fee.
	Added "Libs/DFHack/Fix Handedness", automatic fix for making gloves in reactions.
	Added "User/Archery", a workshop that allows you to make all types of ranged weapons and ammo.
	Removed "User/Dwarf/Ranged Weapons", made obsolete by "User/Archery".
	Fixed bad craft reaction in "User/Bonecarver".
	Permitted the ALCHEMIST labor for the dwarven entity in the "Base" addon.
		This fixes an issue with DFHack's manipulator plugin, namely that you can only toggle labors that are permitted.
	Added a new system for loading DFHack Lua libraries from the region raw directory (I call it "pseudo modules").
		See "Rubble Basics" for details about the new system.
	DFHack scripts and pseudo modules are now loaded dynamically in alphabetical order.
	Changed "Libs/DFHack/Fluids" to use the new pseudo module system.
	Changed all clients of "Libs/DFHack/Fluids" to work with the new system.
	Fixed meat names for "User/Smoked Meats".

v4.6 (for DF 40.10)
	Updated the "Base" addon to DF 40.10
	Changed the way the encoder is handled internally, this should not effect the user in any way.
	It is now legal to have embedded spaces and tabs in base 64 encoded text.
	Added new templates: !DECOMPRESS and !PATCH (simple template interfaces to preexisting script commands).
	The GUI now hides the Rubble console window at all times and displays the log in its own tab.
		This makes it a little easier to use, also makes finding startup errors simpler.
	Updated Rex: By-value command calls are now possible and the error position bug is fixed.
	Somehow I forgot to add the caste templates to the dwarves, so obviously none of the dwarf caste stuff worked.
	Added several new templates for setting entity playability (with script support for advanced users of course).
	Removed the #ADVENTURE_TIER template as the new templates make it obsolete.
	Significantly revised ADDON_HOOKS to interface with the new entity playability templates.
	Expanded the tech class templates a good bit, you can now remove and add classes from tech after registration.
	Expanded the shared item templates to (mostly) match the tech class templates.
	Changed all vanilla reactions and buildings in the "Base" addon to use addon hooks.
	Revised many of the User addons to use the new templates.
	It is now possible to use type command script values as template bodies in addition to type code script values.

v4.5 (for DF 40.8)
	Fixed DFHACK_REACTION, it should now work as expected.
	Added the ability to add custom stuff to the init.lua and onLoad.init files (via scripts).
	The script indexable rubble:raws now properly returns nil for non-existent keys.
	Fixed a bug that stopped "User/DFHack/Fill Barrel" from working.
	Fixed non-working block reaction in "User/Bonecarver".
	Fixed the rubble:activateaddon script command.
	Added script commands rubble:compress and rubble:decompress, compress and decompress strings.
		The result is still plain text, so it can be posted on the forums, great for small addons and bug fixes!
	It is now possible to load base 64 encoded zip files, use the extension ".zip.b64".
		This makes it possible to post zipped addons on the forum as text.
	Rubble now comes with a simple command line utility for decoding and encoding files.
		This utility will decide how to process the file based on extension,
		just drag a single file onto encoder.exe to run.
	Added "User/Dwarf/Castes", an expanded port of the old BD addon, adds profession castes for the dwarves.
	Added "User/Dwarf/Ranged Weapons", allows dwarves to use all ranged weapons (and make ammo for them!).

v4.4 (for DF 40.8)
	Updated the "Base" addon to DF 40.8
	All addons in the "Libs" group are now automatically activated if needed.
	Changed all addons that use bones to use more realistic amounts.
	Added "User/DFHack/Fill Barrel", fill barrels and buckets with water at the still.
	Added "User/Metallurgy", better access to alloys.
	Added "User/Metallurgy/Mithril", a new high performance weapons grade alloy (great for blades and armor).
	Added "User/Metallurgy/DFHack/Volcanic", a new high performance weapons grade alloy (great for hammers and maces).
	Added "User/Metallurgy/Bloomery", more realistic iron and steel production.
	Added "User/Metallurgy/Smelter", more realistic (and more streamlined) smelting.
	Added "User/Metallurgy/Smelter/Remove Alloys", remove the alloy reactions from the smelter.
	Added "User/Fix/Gays", by default this only fixes humans and dwarves.
		The list of creatures to fix can be modified via config variable, see the script header for details.
	Added "User/Smoked Meat", dwarven emergency rations!
	Removed "User/Block/Bone".
	Added "User/Bonecarver", finally, bone furniture (as well as blocks, crafts, and a few other items).
	Saurians are now strictly "straight".
	Removed "User/Add CHILD Tags" as it is no longer needed.
	AXIS VFS has been updated to the latest version, this change is fully backwards compatible.
	Rex (specifically the GenII sub-system) has been updated, this allows extended capabilities for a few script variables.
	The scripting subsystem initialization now takes place in two stages.
		This will allow scripts to be run (with slightly reduced functionality) immediately after Rubble starts.
	The way file type is handled internally has changed significantly (the new system is more flexible).
	Added a new type of script: loader scripts run immediately after the addon activation finishes.
		Loader scripts allow you to do things like load patches and modify addon state before the file list is generated.
		Like init scripts loader scripts are run even if their containing addon is not active or loaded.
	The rubble:activate_addon and rubble:add_file script commands have been removed.
		Loader scripts are more flexible and less error-prone, so this hack is no longer required.
	The "Libs/Base" and "Libs/Castes" addons are now forced active via loader script.
		This is far better than the old way as the addons in question are truly active.
	The way addons are activated (internally) has changed slightly to allow loader scripts to run at the most useful time.
	Some of the standard user script commands have been moved to a loader script (dependency checks mostly).
	Added a new standard user script command to fetch addons by name for use with loader scripts: rubble:fetchaddon.
	Added a new standard user script command to force an addon active: rubble:activateaddon.
		If the addon does not exist an nicely formated error message is printed.
	Added a new native script command: rubble:patch, allows applying patches to strings.
		Patch support is almost certainly buggy! What I tested worked, but I only tested the simplest uses.
	Added two new naive commands for making "virtual addons": rubble:newaddon and rubble:newfile.
		This commands allow you to create new addons in loader scripts.
	Added a new native script command: rubble:filetag, allows reading and editing file tags.

v4.3 (for DF 40.5)
	Fixed "User/Fix/Human Diplomats", two small bugs bugs in it's script kept it from working.
	Fixed "User/Cat Damper", I really should test these things better.
	Split the readme into several parts.
	Split "User/Generic Animal Mats" into many little addons, change only the stuff you want to!
	Made the "User/Block" addons have less redundant content, this change should not affect users.
	Fixed an error that kept rubble:groupactive from working.

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
