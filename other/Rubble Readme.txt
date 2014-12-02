
Rubble: After Blast comes Rubble

==============================================
Overview:
==============================================

Rubble is a raw generator, eg. it takes "preraws" and generates Dwarf Fortress raw files from them.

By using "templates" (either predefined or defined yourself) you can greatly simplify the task of writing a Dwarf Fortress mod.
For example, there are templates that automate registering workshops and reactions, which will not only save you time but also reduce the chance of bugs as well.

For non-modders Rubble makes it easy to install mini-mods so you can create your perfect custom Dwarf Fortress.

Required Reading for Modders (in most-important-first order):
	This readme
	HowTo Rubble - Tutorials
	Rubble Addons - lots of stuff about addons in general
	Rubble Base Templates - the template documentation for the base_templates addon
	Rubble Native Templates - the template documentation for the builtin templates
	The included addons - A little short on comments but still a great resource
	Raptor Basics - Raptor overview
	Everything in "raptor command docs" - the Raptor command documentation, NEEDED if you plan to do anything with Raptor

If you do not plan on doing any modding the only thing you need to read is this readme.

==============================================
Install:
==============================================

Backup your "raw/objects" folder! Rubble will delete all your existing raw files!
Extract Rubble to "<DF Directory>/rubble".
Install your addons to "<DF Directory>/rubble/addons"
(If you use DFHack) add "script rubble.dfhack" to your dfhack.init file
Run "rubble -addonlist" OR run "Rubble GUI.exe"

Now you are good to go! Documentation (as you have obviously discovered) is in the "rubble/other" folder as is source code and OSX/Linux binaries.

To activate or deactivate a Rubble addon you may set it's entry in addons/addonlist.ini to "true" or "false". If you just added an addon it will not have an entry until Rubble has run at least once (after the addon was added). 
If you want to run Rubble without generating anything so as to update the addon list file just run 'rubble -addonlist'

If you can (eg. you are running on windows) it is HIGHLY recommended to use the GUI. The GUI automates the process of updating and editing addonlist.ini and is generally much faster then doing everything by hand (plus you don't have to mess around with the command prompt, if you dislike that kind of thing)

If you use OSX or Linux, 32 bit binaries for these OSes can be found in the "rubble/other" directory. If you want 64 bit binaries you can compile them yourself, source code is in other/src.

==============================================
Configure:
==============================================

Rubble allows you to change its settings via command line options or a config file. To see these options and their defaults run "rubble -h".

ALL command line options may also be specified in the config file.

Rubble tries to read the file "./rubble.ini", if this does not fail Rubble will load settings from here before processing command line options (command line options always take precedence).
Example "rubble.ini" (using some of the defaults):
	[path]
	dfdir = ..
	outputdir = ../raw/objects
	addonsdir = ./addons

Rubble supports external launchers via the -config and -addons command line options
-config allows you to set/create rubble variables.
	Usage: -config="key1=value1;key2=value2;keyN=valueN" 
		(Note that the ';' may need to be a ':' on non-windows systems)
-addons allows you to override the default rules for loading addons by explicitly listing which ones you want to load.
	Usage: -addons="addon1;addon2;addon3;ect"
		(the same note about ';' applies here as well)
These two options are mostly for use by external launchers.

If you want to regenerate the raws for a save just run 'rubble -outputdir="../data/save/<savename>/raw/objects"'

All directories used by Rubble must exist (you should get an error message if not).

==============================================
Included Addons:
==============================================

I include addons in the base Rubble install that fix bugs, demo something useful, or are just too useful to leave out. 
Needless to say these addons make a good place to look for info on how to do something.

If you want to generate the Rubble version of "vanilla DF" you will need to activate the "Base/Files", "Base/Templates", and "Base/Clear" addons, these addons are already active in the default addonlist.ini.

Add CHILD Tags:
	Adds a CHILD tag to every creature that does not already have one.
	DOES_NOT_EXIST, NOT_LIVING, and EQUIPMENT_WAGON count as CHILD tags for the purpose of this script, 
	vermin are also skipped.
	This is an advanced tweak script that processes files in multiple passes.
	This addon is stand-alone, as it uses no outside templates.

Base/Clear:
	This must-have addon clears the output folder before generation.
	This addon is completely stand-alone.
	You should never disable this addon without a good reason!

Base/Files:
	The standard base addon. Contains Rubblized versions of the vanilla raw files. 
	Do not disable unless you have a replacement.

Base/Templates:
	This contains all kinds of useful templates.
	Required by most of the other addons.
	Do not disable unless you like error messages :)

Dev/Cheat:
	A bunch of cheat reactions to allow for fast setup of testing forts.
	Build a "Heavy Supply Wagon" workshop, both the building and all reactions are free.
	Uses animal caretaker skill (the most useless skill I could think of :) )

Dev/Dummy Reactions:
	This addon adds "dummy reactions" eg. empty reactions registered to "ADDON_HOOK_PLAYABLE".
	By default 15 reactions are created, if you want a different number of reactions to be generated pass '-config="DEV_DUMMY_REACTION_COUNT=count"' to Rubble when generating.
	These reactions are empty by default, use if you want to add new content after worldgen.

Fix/Butcher Inorganic:
	Allows you to butcher creatures made from (some) inorganics, FB stone walls anyone?

Fix/Undead Melt:
	Fixes bad temperature values in material_template_default.
	Replaces material_template_default with a new copy, but the only changes are the temperature values.
	In any case this should make it possible to melt undead with magma.

Fix/Usable Chitin:
	Allows you to use chitin as shell.

Fix/Usable Feathers:
	Allows you to use feathers as "soft" pearls (you can use them to make crafts).

Fix/Usable Scale:
	Allows you to use scales like leather.

Fix/Vermin Variations:
	This adds two fixes:
		It forces all giant animals to have a life span of at least 10 years and
		It replaces all pool biome tags in giant and animal-man creatures with the correct lake biome.
	This is useful to make variations on vermin appear in game and to help them live long enough to be interesting.
	
Generic Animal Mats:
	Make animal mats such as meat and leather generic. This is mostly for those who have FPS issues.
	This addon is packed as a zip to demo how that is done.

Libs/Castes:
	Adds support for dynamically adding castes to a creature.

Libs/Castes/DFHack/Transform:
	Adds support for transforming castes via autosyndrome and true transform.
	Requires the "Libs/Castes" addon.

Libs/Crates:
	Adds infrastructure for "crates", eg bars that can be bought from caravans and then "unpacked" via special reactions into items.
	This templates added by this addon are always available (even when the addon is not active). If the addon is not active then the templates will not actually do anything.
	The main reason for having always active templates is so that adding crate support in many addons is very easy.

Libs/DFHack/Add Reaction:
	Allows you to add reactions to buildings that are normally not possible to add reactions to, also allows you to remove all hardcoded reactions from a building.
	Requires the "Base/Templates" addon.
	WARNING! Largely untested!

Libs/DFHack/Announcement:
	This addon is a modders resource for displaying announcements to the player.
	Requires "Libs/DFHack/Command".
	This addon also provides a Raptor command for advanced usage.
	This templates added by this addon are always available (even when the addon is not active). If the addon is not active then the templates will not actually do anything.

Libs/DFHack/Command:
	This addon is a modders resource for running DFHack commands from reactions.
	Activating this addon adds a template that generates reaction product lines for autosyndrome boiling rocks (the boiling rocks in question are automatically defined).
	This addon also provides a Raptor command for advanced usage.

Libs/DFHack/Spawn:
	"Libs/DFHack/Spawn" is a modders resource for spawning creatures via DFHack.
	Activating this addon installs a creature spawner script into the DFHack scripts directory.
	Requires "Libs/DFHack/Command".
	This addon also provides a Raptor command for advanced usage.

Nerf/Ranged:
	A simple addon that nerfs bows and crossbows, uses the stats from the popular Broken Arrow mod.

Nerf/Whips:
	Nerfs whips and scourges so they are not deadly anti-armor weapons.

Tileset/MLC:
	A simple addon that installs my custom ASCII-like tileset.
	This addon makes no raw changes, so it should be usable with all ASCII compatible mods.
	Demos tileset addons.

Tileset/Vanilla:
	Restores the vanilla tileset related init options.
	This addon makes no raw changes, so it should be usable with all ASCII compatible mods.

Zap Aquifers:
	Disables all AQUIFER tags.
	A simple replacement for #AQUIFER that does not need the config file or require modders to use a template.
	This addon is stand-alone, as it uses no outside templates.

==============================================
BUGS (and other issues):
==============================================

Known Issues (non-bug):
	
	You should not use the extension .txt for readme files, as this will cause Rubble to parse those files.
		See "Rubble Addons.txt" for more details.
	
	It may be possible to see FBs made from "dfhack command" or "caste transformation", there is nothing I can do about this as autosyndrome requires the use of inorganic materials.

When making an error report please post the FULL log file! Posting just a few line tells me almost nothing about the problem, after all I made Rubble write all that stuff for a reason :)

If any of the documentation is not 100% clear just ask. 
I know everything there is to know about how Rubble works and so I tend to forget to put 
"obvious" stuff in the docs. A reminder that I forgot something is more a help than a hindrance.

==============================================
Changelog:
==============================================
v3.10
	Updated to Raptor 3.2
		A (very) slight performance boost, a minor bug fix, a possible string corruption fix,
		and some other minor changes that don't effect Rubble.
	Added profiling (writes data in the format expected by pprof), not that you care.
	Added rubble:ingroup, check if an addon is part of a specified group.
	Added rubble:groupactive, check if any addon in a specified group is active.
	Added native script commands rubble:nowritefile, rubble:graphicsfile, and rubble:currentfile.
	rubble:getvar now returns "" for non-existent variables.
	The behavior of some of the dependency check commands has changed (to make them more accurate)
	The behavior of "Dev/Dummy Reactions" has changed slightly (should not effect users).
	Added some caste related library addons.
	Added the time templates (@ADV_TIME and @FORT_TIME) back to Rubble, this time they are in the "Base/Templates" addon.
		(These templates were removed in 3.0 in order to remove all non-essential native templates)
	Added @GROWDUR, a replacement for the vanilla GROWDUR tag that uses the same mechanism as the other time templates.
	Execution of init scripts moved one step later (to just after built-ins are loaded)
		This should make it so that init scripts can do literally everything that a pre tweak script can.
	Added #INSTALL_GRAPHICS_FILE and @GRAPHICS_FILE, support for creature graphics.
	Updated the Notepad++ UDL (lots of changes from previous versions never got added, now they are).

v3.9
	Fixed lots of cases where Raptor error handlers still used panic, oops.
	Updated some of the docs.
	Moved "Libs/Crates" from BD2 to Rubble.
	Changed "Libs/Crates" to use "always active" templates.
	Changed "Libs/DFHack/Announcement" to use "always active" templates.
		"always active" templates are always available, they just return a placeholder string when the addon is not active.
	"Libs/DFHack/Command" now adds [SPECIAL] to it's generated inorganics.
		This should stop command stone FBs from generating.
	Added rubble:placeholder to help with creating always active templates.
	Updated to Raptor 3.1, this includes only minor changes and bug fixes.

v3.8
	Updated to Raptor 3.0
		This update is strictly an internal refactoring, no commands have changed.
		The internal API is much better now.
		Error messages (on the Raptor side of things) should be much better.
	Changed LOTS of stuff internally to get better error messages.
		(I REALLY need to refactor the Rubble internals sometime...)
	Fixed Rubble reporting that reading addonlist.ini failed when in shell mode even if the read succeeded.
		(it would also update the addon list twice when this happened)
	Init scripts may now override each other.
		The only reason they couldn't before was oversight, init scripts are handled separate from 
		other addon files, and I forgot to allow overriding.
	Added rubble:compat, an indexable of version compatibility information.
		Index with a version string, returns true if the current version is backwards compatible with the requested version.
		If a version is not listed it will return nil (which will act like false in most cases).
		Obviously always reports nil for versions newer than the current version.
		Versions start with "3.8" (so [rubble:compat "3.7"] will return nil)
	Removed command "panic", use rubble:abort.
	Added template DFHACKINIT_ADD to "Base/Templates", allows automatic customization of dfhack.init by addons.
		To activate this feature you need to add "script rubble.dfhack" to your dfhack.init file.
	Added experimental "Libs/DFHack/Add Reaction" addon, manipulate hardcoded buildings/reactions.

v3.7
	Made some template changes to "Base/Templates":
		Added @IF_ACTIVE
		Added @IF_CODE
		Changed IF to @IF
		Changed SET to @SET
		(use the various versions of ECHO and VOID to force a parse stage with @ templates)
	There is now a builtin version variable exposed to scripts: rubble:version
		(type string, set to "3.7" this version)
	ALL command line parameters may now be listed in the config file.
	The way variables in template parameters are handled has changed a little:
		All variables in parameters are now expanded before the template is called.
	Fixed lots of little errors in the docs.
	Some very minor bugfixes.
	Updated the Notepad++ UDL

v3.6
	Updated to Raptor 2.2
		Changes:
			Added a new type (nil)
			Added a new mode for the syntax checker
			Removed obsolete command: debug:clrret, use nil.
			Added new command: copy, allows you to copy a value.
			Added a simple threading library (threading is not useful in Rubble by default).
			Raw lexer commands now require/return a lexer handle.
		This update (particularly the raw parser changes) will break most post (and a few pre) tweak scripts!
	To go with the Raptor threading library you can now run Rubble in "threaded" mode via "-threads".
		This is only useful for multi-threaded tweak scripts, default off.
	Fixed bug in "Libs/DFHack/Announcement", message is now written to gamelog as expected.
	The filename key for rubble:raws no longer has it's extension stripped
		This should resolve issues with files overriding each other when they shouldn't.
		WARNING! This change will break some tweak scripts! (it's an easy fix)
	Updated the tweak scripts to handle the changes to rubble:raws and the raw parser.
	Forced init scripts are now called "init scripts" and no longer need to be named "forced_init".
		Now they use the extension ".init.rsf", functionality is unchanged.
	Renamed the Base init script to "base.init.rsf"
	The "Add CHILD Tags" addon is a little more discriminating.

v3.5
	Updated to Raptor 2.1.2
		2.1.2 adds two new commands, sort:new and sort:map as well as fixing an issue in the shell.
	Rubble is more flexible about it's environment when running in shell mode.
	Added "Libs/DFHack/Announcement", a library addon for displaying announcements to the player.
	Added "Libs/DFHack/Command" (from BD), a library addon for running DFHack commands.
	Added "Libs/DFHack/Spawn" (from BD), a library addon for spawning creatures.

v3.4
	Updated to Raptor 2.1.1
		2.1.1 adds a new command, trap.
	Removed ITEM template (use SHARED_ITEM)
	Added rubble:activeaddons, a params array indexable with a list of all active addon names
	Added rubble:grouprequires script command to the "Base" addon forced init script
		Allows for specifying a dependency for a group of addons.
	Added rubble:groupincompatible script command to the "Base" addon forced init script
		Allows for specifying an incompatibility for a group of addons.
	Rewrote addon hooks to make output format better (still a little ugly)

v3.3
	Added the following templates to help streamline tileset support:
		INSTALL_TILESHEET
		SET_FULLSCREEN_TILESHEET
		SET_WINDOWED_TILESHEET
		OPEN_D_INIT
		EDIT_D_INIT
		CLOSE_D_INIT
	Changed the default tileset addons to use the new templates.
	Renamed ITEM to SHARED_ITEM (to be consistent with the other SHARED_OBJECT templates)
	Added ITEM as an alias for SHARED_ITEM, this will be removed in a future version
	The SHARED_ITEM template is now tileset aware (for tools)
	SHARED_ITEM now properly handles type DIGGER
	Updated SHARED_ITEM to handle type FOOD
	Updated the "Base/Files" addon to reflect the changes to SHARED_ITEM
	Removed REGISTER_ORE, I never used it and you can just use SHARED_OBJECT_ADD instead
	Made SHARED_OBJECT_ADD able to add to shared objects that are not yet defined as well as those already defined
		Code added by SHARED_OBJECT_ADD will always come at the end of the definition
	Removed ONCE and STATIC, as they were just Blast holdovers that have no real use as far as I could see
	Added !ECHO, #ECHO, !VOID, #VOID, and V (all are aliases of existing templates)
	Rubble now includes an embedded copy of the full-up Raptor shell!
		Use "-shell" to enter shell mode, all normal Raptor shell command-line arguments work :)
		The shell is entered just before generation would have started.
		See "Rubble Addons.txt" for details.
	Fixed bug in script command rubble:calltemplate.
	Updated the Notepad++ UDL

v3.2
	Added forced initialization scripts, special scripts named "forced_init.rsf" (or .rbf) that ALWAYS run
		Forced init scripts will run for any addon, even addons that are not active or addons that have no parseables
		See "Rubble Addons.txt" for more details.
	Added rubble:requireaddon to the "Base" addon (as a forced init script),
		checks if an addon is active and aborts (with a nice message) if not
	Added rubble:incompatibleaddon to the "Base" addon (as a forced init script),
		checks if an addon is active and aborts (with a nice message) if so
	Moved the Raptor shell predefs to the new forced init script in "Base"
	Updated to Raptor 2.1
		Raptor 2.1 fixes a few bugs, one of which is fairly major

v3.1
	The addon loader is now recursive
		See "Rubble Addons.txt" for the new loading rules.
	Grouped a bunch of addons
		base -> Base/Files
		base_templates -> Base/Templates
		clear_raws -> Base/Clear
		broken_arrow -> Nerf/Ranged
		nerf_whips -> Nerf/whips
		fix_vermin_variations -> Fix/Vermin Variations
		mlc_tileset -> Tileset/MLC
	Updated the "Add CHILD Tags" addon to handle DOES_NOT_EXIST and EQUIPMENT_WAGON creatures
	Added "Fix/Undead Melt", "Fix/Usable Chitin", "Fix/Usable Feathers", and "Fix/Usable Scale" addons
	Added "Dev/Dummy Reactions" addon
	Added "Tileset/Vanilla" addon
	Added SHARED_OBJECT_ADD template to "Base/Templates" addon.
		SHARED_OBJECT_ADD allows you to append tags to an existing SHARED_OBJECT, very useful.
	Added new script command rubble:template, works just like !SCRIPT_TEMPLATE
	Added new script command rubble:addonactive, returns true if the named addon is active
	Added new script command rubble:abort, prints a message and quits, use instead of panic for most errors
	Replaced !PANIC, PANIC, and #PANIC with !ABORT, ABORT, and #ABORT
	Converted "Base/Templates" into a pre script
		This ensures that all base templates are defined before raw parsing starts 
		(This also improves error messages tremendously)
	Added new template prefix character, @, means "run in first possible stage"
		See "Rubble Addons.txt" for more details
	Script templates are now stored as pre-lexed code, this provides a minor performance boost
	I forgot to update the version number in the header (again), fixed
	Updated to Raptor 2.0, this should make absolutely no difference, if any behavior has changed it's a bug
		The only reason I updated is so that I would not have to maintain two versions of Raptor
	Squashed some stupid script bugs
		Looks like the new Raptor script validation tool was worth it :) too bad it only catches syntax errors...
	Rubble now reports errors encountered when writing files out
	Fixed "Zap Aquifers", I broke it when updating from NCA to Raptor
	The GUI is now resizeable and allows more configuration
		If you have a rubble.ini the GUI will now read the addons directory path from it
		You can set the path to the directory containing rubble.exe in gui.ini
	Updated the Notepad++ UDL
	
v3.0
	Added ability to load an addon from a zip file
	Added ability to make addon "packs", eg. multiple addons grouped in one zip file
		Use the extension .pack.zip, one addon per subdirectory, no nesting
	The generic animal mats addon is now packed as a zip to demo how that is done
	Replaced addon enabling mechanism with an automatically updated ini file, addonlist.ini
	Updated the Rubble GUI to read addon names and active status from addonlist.ini
	Rewrote the lexer and parser from scratch
		The new versions are much better, certainly easier for me to work with (and possibly a little faster)
	All addon files (not just parseable files!) are now available in rubble:raws
	Added an addon to add a CHILD tag to all creatures that do not already have one
	Added an addon that forces all giant animals to have a life span of at least 10 years and 
		replaces all pool biome tags in giant and animal-man creatures with the correct lake biome
	Added an addon to nerf whips (just like broken arrow, but for whips)
	Removed the #AQUIFER template from the base_templates addon
		It was simple to replace this with a tweak script, and not having to use a template results in 
		less work for modders (and me) in the long run.
	Added an addon to strip AQUIFER tags
	Rubble now duplicates it's output to the file rubble.log
	Removed the #ADV_TIME and #FORT_TIME native templates
		I never used them and having non-critical native templates rubs me the wrong way.
		As Raptor has floating point commands it should be possible to reimplement these, so if anyone 
		actually used them please tell me.
	Updated from NCA7 to Raptor v1.0
		Raptor is based on NCA, but there are (a lot of) differences, see Porting.txt for details.
	Pre and Post tweak script extensions have changed from .nca to .rsf (Raptor Script File) or .rbf (Raptor Binary File)
	Fixed some script errors that I discovered as I was porting everything
	Removed support for the base
		The old base is now just another addon (named "base").
	The base templates are now in their own addon named "base_templates"
	Stripped support for addon config files completely
		If you really need a config file it is easy to use a pre script to load one.
	Removed the rubble:basedir script variable as it is obsolete
	Added an optional parameter to rubble:stageparse to allow specifying the parse stage to use
	Fixed some minor mistakes in the docs
	Rewrote the Notepad++ UDL
	
v2.0
	Added pre and post tweak scripts
		Pre tweak scripts have the file extension ".pre.nca" and post scripts have the extension ".post.nca"
	Added external launcher support via -config and -addons
	Added a simple GUI launcher, sorry it's windows only :(
	Fixed incorrect version info in header
	Rubble will now only parse files with the extensions .txt or .rbl 
		.rbl files will not be written out. Saves on calls to rubble:skipfile!
	Default directory settings changed, Rubble no longer bloats your save by ~10MB! (unless you want it to)
	Changed the way addon and base config files are handled
		Config files are now in the same directory as the base or addon and are named "config.ini"
	I think I fixed all the documents I broke... If you see something wrong tell me
	Fixed #ADVENTURE_TIER
	The raw files are now available as unparsed strings to NCA scripts via the indexable rubble:raws
		Skipped files are not available
		This is mostly useful only for pre and post tweak scripts
	Added simple NCA raw parser, currently you can disable, replace or append tags, see command docs for "raw"
	Fixed fire clay giving earthenware items, not sure how this came to be
	Changed the generic_animal_mats addon to use tweak scripts, it should be much more compatible now
	Fixed the generic_animal_mats addon clobbering animal venom, webs, and the like
	The usual little bits added to NCA, just some new commands this time
		See the command docs for base, str, bool and raw for details.
	Updated the Notepad++ UDL
	Stopped ITEM_CLASS from writing junk to the raws, another case of forgetting to clear the return register
	Added template SHARED_OBJECT_EXISTS to base

v1.9
	Fixed ITEM_CLASS, this template is (AFAIK) not used in the base, but it is used in the Broken Arrow addon.
		On line 652 of file "base_templates.txt" there was an extra ')', Oops.
	Rubble now has an icon on windows!
	There may be some small bug fixes to the script engine, I worked on it some and I don't remember if I fixed anything other than adding hooks for (very, very) low level debugging.
	Made a few small fixes to the docs, nothing to worry about.
	
v1.8
	Updated to NCA7, this includes many small (and a few large) script changes
		If you wrote any script code check the command docs, 
		a bunch of commands are now namespaced and/or have had their name changed.
		For example add is now int:add and append is now str:add.
	Notepad++ UDF updated
	Added ADDON_HOOKS template to the base
	Added optional config file, not sure why I didn't do it earlier after all GoBlast had one...
	
v1.7
	Slight changes to NCA Indexables, nothing that effects user script code
	The (much expanded) mess of debugging commands have all been moved to the debug namespace, some have new names
		Old commands with new names are:
			valueinspect	->	debug:value
			ncash			->	debug:shell
		New commands are:
			debug:list		->	Lists all global data, including variables in the root environment
			debug:namespace	->	Lists all variables, commands, and namespaces in a namespace
			debug:env		->	Lists all variables in all environments, sorted by environment
			debug:clrret	->	Clears the return value, useful for cleaning up after playing around in the shell
	Notepad++ UDL updated
	Error reporting now should display a line number for Rubble errors
		Script error line numbers are still an offset from the start of the script
		Errors in templates will refer to the place where the template was called

v1.6
	Updated script runtime to NCA6, this brings major improvements to the way script values are handled
	Added new NCA command len, the old pre-v6 way of fetching an Indexable's element count will no longer work
	Removed source for NCASH, this is now replaced with the ncash NCA command
	Re-exported all of the NCA command documents
	Fixed all the templates to work with NCA6
	Removed the old item templates (the native ones from 1.1), this was supposed to happen in 1.5 but I forgot
	Changed ITEM_CLASS to be more flexible
	Changed SHARED_INORGANIC and SHARED_MATERIAL_TEMPLATE to automatically support tilesets
	Added template SHARED_PLANT, just like SHARED_INORGANIC but for plants
	Added template #WRITE_TILESET to help modders make tileset addons

v1.5
	NCA variable dereference syntax now allows index dereferencing a value (eg. [value index])
	NCA base commands map and array changed, they no longer create a variable, they only return the new map or array
	NCA base commands exists and set now have multiple meanings, check NCA base docs for details
	NCA base command foreach now takes a map or array VALUE instead of a NAME
	Added the "NCASH predefs" to the base
	Added valueinspect NCA command, very useful for debugging
	Added new NCA base command evalinnew, works like run but without param support
	Added new NCA string command trimspace, trims leading and trailing whitespace
	Fixed major bug with NCA arrays, appending to an array did not work
	Ported the following templates to Rubble+NCA:
		BUILDING_WORKSHOP
		BUILDING_FURNACE
		#USES_BUILDINGS
		REACTION
		#USES_REACTIONS
		#USES_TECH
		REGISTER_ORE
		#_REGISTERED_ORES
		REGISTER_REACTION_CLASS
		#_REGISTERED_REACTION_CLASSES
		REGISTER_REACTION_PRODUCT
		#_REGISTERED_REACTION_PRODUCTS
	Rewrote all the item templates
		ITEM no longer takes a class and rarity
		ITEM_CLASS takes the place of ITEM's removed params and ITEM_RARITY
		#USES_ITEMS now only takes one class
	Added Notepad++ user defined language file for NCA5 and Rubble code
	Added source code for NCASH5, a useful debugging tool for scripts
	Added documentation for the NCA base language

v1.4
	Updated script runtime to NCA5, this may break some code as the way maps are handled is changed
	Updated all NCA docs to describe NCA5
	Changed the way the lexer handles the char literals (';', '{', and '}') to remove possible infinite recursion
	Fixed bug in some of the base templates, the return value should have been run through the stage parser
	Fixed #AQUIFER not working at all
	Fixed SET returning whatever junk was in the NCA return register when it was called
	Fixed some variables expanding too early, variables in nested template calls are not expanded until the nested template runs
	Removed SHARED_ITEM (alias for SHARED_OBJECT)
	Added native template !SCRIPT_TEMPLATE for declaring templates consisting of NCA code
	Added base templates SET_TILE and #TILE for tileset support
	Wrapped every (non-creature) tile number with a call to #TILE
	Added rubble:dfdir variable
	Added regex:replace NCA command
	Added an example tileset addon mlc_tileset based on an ASCII-like tileset I made
	Split generic animal mats out of the base and into an addon
	Removed ANIMAL_MAT template from base
	Made clear_raws addon active by default

v1.3
	Added rubble:stageparse NCA command
	Added rubble:calltemplate NCA command
	Added rubble:expandvars NCA command
	Removed some templates that were easy to convert to Rubble+NCA code
		Removed templates are:
			COMMENT
			C
			VOID
			PANIC
			IF
			ONCE
			STATIC
			SET
			#ADVENTURE_TIER
			SHARED_INORGANIC
			SHARED_MATERIAL_TEMPLATE
		The above templates are now in the base template file ("raw/source/base/base_templates.txt")
	Added template ECHO (alias E) to the base
	Added templates !PRINT, PRINT and #PRINT to the base
	Added templates !PANIC and #PANIC to the base
	Replaced SHARED_ITEM with SHARED_OBJECT, moved SHARED_OBJECT to the base
	Added the ability for a template to take the params from the previous template call via a special ... param
	Added more info to the NCA docs (still very spotty)
	Fixed up formatting for SHARED_INORGANIC and SHARED_MATERIAL_TEMPLATE a little
	Made Rubble expand vars in all files as a final (additional) step after postparsing
	Fixed some more NCA bugs, as much as I use NCA you would think they would all be squashed by now.

v1.2
	Added rubble:skipfile NCA command
	Added rubble:getvar NCA command
	Added rubble:setvar NCA command
	Made REGISTER_REACTION_PRODUCT parse it's material
	Added addon tech and item hooks to entities
	Made the lexer handle some char literals (';', '{', and '}')
	Updated the base and addons to use the new abilities
	Fixed some minor bugs in NCA4 and updated docs to match
	Fixed bug that made files process in the wrong order
	Stopped Rubble from mangling special chars

v1.1
	Added the NCA4 file system commands
	Wrote a huge amount of docs
	Added PANIC template to allow aborting
	Added panic NCA command
	Added NCA variables for each directory setting
	Configuration files are now optional
	Rewrote item and tech class templates, things should be much cleaner than before
	Added a base and two example addons

v1.0
	First version
