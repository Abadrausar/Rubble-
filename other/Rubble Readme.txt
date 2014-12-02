
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
	Rubble Basics - Lots of stuff about addons and Rubble in general
	Rubble Base Templates - The template documentation for the "Base/Templates" addon
	Rubble Native Templates - The template documentation for the built-in templates
	Rubble Libs Templates - The template documentation for the library templates
	The included addons - A little short on comments but still a great resource
	Everything in "Rex Docs" - the Rex documentation, NEEDED if you plan to do anything with Rex

If you do not plan on doing any modding the only thing you need to read is this readme.

==============================================
Install:
==============================================

Backup your "raw/objects" folder! Rubble will delete all your existing raw files!
Extract Rubble to "<DF Directory>/rubble".
Install your addons to "<DF Directory>/rubble/addons"
Run "rubble -addonlist" OR run "Rubble GUI.exe"

Now you are good to go! Documentation is in the "rubble/other" folder as is source code and OSX/Linux binaries.

To activate or deactivate a Rubble addon you may set it's entry in addons/addonlist.ini to "true" or "false". If you just added an addon it will not have an entry until Rubble has run at least once (after the addon was added). 
If you want to run Rubble without generating anything so as to update the addon list file just run 'rubble -addonlist'

If you can (eg. you are running on windows) it is HIGHLY recommended to use the GUI. The GUI automates the process of updating and editing addonlist.ini and is generally much faster then doing everything by hand (plus you don't have to mess around with the command prompt, if you dislike that kind of thing)

If you want to run multiple worlds with radically different addon sets it is a good idea to run "rubble -prep=<region name>" before switching worlds (this is mostly only for tilesets, other addons should be good to go unless they do something weird and non-standard).

If you use OSX or Linux, 32 bit binaries for these OSes can be found in the "rubble/other" directory. If you want 64 bit binaries you can compile them yourself, source code is in other/src (along with basic build instructions).

==============================================
Configure:
==============================================

Rubble allows you to change its settings via command line options or a config file. To see these options and their defaults run "rubble -h".

ALL command line options may also be specified in the config file.

Rubble tries to read the file "./rubble.ini", if this does not fail Rubble will load settings from here before processing command line options (command line options always take precedence).
Example "rubble.ini" (using some of the defaults):
	[rubble]
	dfdir = ..
	outputdir = df:raw
	addonsdir = rubble:addons
Duplicate keys are just fine, they will act pretty much exactly like duplicate options on the command line.

Rubble supports external launchers via the -config and -addons command line options
-config allows you to set/create rubble variables.
	Usage: -config="key1=value1;key2=value2;keyN=valueN"
-addons allows you to override the default rules for loading addons by explicitly listing which ones you want to load.
	Usage: -addons="addon1;addon2;addon3;ect"
These two options are mostly for use by external launchers.

If you want to regenerate the raws for a save just run 'rubble -outputdir="../data/save/<savename>/raw/objects"'

All directories used by Rubble must exist (you should get an error message if not).

==============================================
Included Addons:
==============================================

I include addons in the base Rubble install that fix bugs, demo something useful, or are just too useful to leave out. 
Needless to say these addons make a good place to look for examples.

If you want to generate the Rubble version of "vanilla DF" you will need to activate the "Base/Files", "Base/Templates", and "Base/Clear" addons, these addons are already active in the default addonlist.ini.

Add CHILD Tags:
	Adds a CHILD tag to every creature that does not already have one, thereby allowing them to breed in fortress mode.
	(Note that sexless creatures will still not breed)
	DOES_NOT_EXIST, NOT_LIVING, and EQUIPMENT_WAGON count as CHILD tags for the purpose of this script, vermin are also skipped.
	This addon is stand-alone, as it uses no outside templates.

Base:
	The standard base addon. Contains Rubblized versions of the vanilla raw files. 
	Do not disable unless you have a replacement.

Building Dimensions:
	Append all building names with their dimensions, really simple.

Dev/Cheat:
	A bunch of cheat reactions to allow for fast setup of testing forts.
	Build a "Heavy Supply Wagon" workshop, both the building and all reactions are free.
	Uses animal caretaker skill (the most useless skill I could think of :p)
	This addon includes basic support for use with the BD3 addon pack.

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

Fix/Butcher Inorganic:
	Allows you to butcher creatures made from (some) inorganics, FB stone walls anyone?
	(Not 100% sure this works, copied from Modest Mod)

Fix/Fish Pops:
	Raises fish vermin populations to the max.
	This should stop you from fishing out a lake or ocean ever again! (unless you catch ~30,000 fish...)

Fix/Undead Melt:
	Fixes bad temperature values in material_template_default.
	Replaces material_template_default with a new copy, but the only changes are the temperature values.
	In any case this should make it possible to melt undead with magma.

Fix/Usable Chitin:
	Allows you to use chitin as shell.

Fix/Usable Feathers:
	Allows you to use feathers as "soft" pearls (you can use them to make crafts).

Fix/Usable Scale:
	Allows you to use scale like leather.

Fix/Vermin Variations:
	This adds two fixes:
		Forces all giant animals to have a life span of at least 10 years and
		Replaces all pool biome tags in giant and animal-man creatures with the correct lake biome.
	This is useful to make variations on vermin appear in game and to help them live long enough to be interesting.

Libs/Base:
	All the base templates and other such stuff, critical to normal operation.
	If you do not activate this addon it will be automatically activated for you during normal generation cycles.
	(in other words: ignore this addon in most cases)

Libs/Castes:
	Adds support for dynamically adding castes to a creature.

Libs/Castes/DFHack/Transform:
	Adds support for transforming castes via autoSyndrome and syndromeTrigger.
	Requires the "Libs/Castes" addon.

Libs/Crates:
	Adds infrastructure for "crates", eg bars that can be bought from caravans and then "unpacked" via special reactions into items.
	The templates added by this addon are always available (even when the addon is not active). If the addon is not active then the templates will not actually do anything.

Libs/DFHack/Add Reaction:
	Allows you to add reactions to buildings that are normally not possible to add reactions to, also allows you to remove all hardcoded reactions from a building.
	Requires the "Base/Templates" addon.
	WARNING! Largely untested!

Libs/DFHack/Announcement:
	This addon is a modders resource for displaying announcements to the player.
	Requires "Libs/DFHack/Command".
	This addon also provides a script command for advanced usage.
	The templates added by this addon are always available (even when the addon is not active). If the addon is not active then the templates will not actually do anything.

Libs/DFHack/Command:
	This addon is a modders resource for running DFHack commands from reactions.
	Activating this addon adds a template that generates reaction product lines for autoSyndrome boiling rocks (the boiling rocks in question are automatically defined).
	This addon also provides a script command for advanced usage.

Libs/DFHack/Fluids:
	Adds the extremely useful rubble_fluids DFHack Lua library and command.
	Also adds three convenience templates for filling minecarts with fluids, eating fluids, and spawning fluids.
	Requires "Libs/DFHack/Command" (for the templates).

Libs/DFHack/Spawn:
	"Libs/DFHack/Spawn" is a modders resource for spawning creatures via DFHack.
	Activating this addon installs a creature spawner script into the DFHack scripts directory.
	Requires "Libs/DFHack/Command".
	This addon also provides a script command for advanced usage.

Libs/Macro:
	A modders resource for generating macros from scripts.
	Adds no templates, only script commands.

Libs/Macro/Example:
	A simple addon that generates a large complicated macro from images, makes a good demo/test of "Libs/Macro".
	This addon takes ~7-8 seconds to generate (the cursor movement generation command is kinda slow, it's a work in progress), in any case activating this addon more than triples normal run time.

Nerf/Ranged:
	A simple addon that nerfs bows and crossbows, uses the stats from the popular Broken Arrow mod.

Nerf/Whips:
	Nerfs whips and scourges so they are not deadly anti-armor weapons.

Tileset/MLC/Normal:
	A simple addon that installs my custom ASCII-like tileset.
	This version is very close to default ASCII.
	This addon makes no raw changes, so it should be usable with all ASCII compatible mods.

Tileset/MLC/Tracks:
	A simple addon that installs my custom ASCII-like tileset.
	This version has really nice track graphics. There are a bunch of things that will look a little funny, but tracks look great.
	This addon makes no raw changes, so it should be usable with all ASCII compatible mods.

Tileset/MLC/DFHack/TWBT:
	A simple addon that installs my custom ASCII-like tileset.
	This version uses the "Text Will Be Text" DFHack plugin (you will need to install that plugin separately).
	Unlike most ASCII tilesets this one has special tiles for most buildings and items, making it (almost) fully graphical.
	This addon makes no raw changes, so it should be usable with all ASCII compatible mods.

Tileset/Vanilla:
	Restores the vanilla tileset related init options.
	It is a very good idea to use this addon whenever you are using the vanilla tileset, as otherwise prepping worlds generated from those raws may not work quite right.
	This addon makes no raw changes, so it should be usable with all ASCII compatible mods.

Tileset/Vanilla/Graphics:
	The vanilla creature graphics, mostly just a demo.

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
	
	It may be possible to see FBs made from "dfhack command" or "caste transformation", there is nothing I can do about this as autoSyndrome requires the use of inorganic materials.

When making an error report please post the FULL log file! Posting just a few lines tells me almost nothing about the problem, after all I made Rubble write all that stuff for a reason :)

If any of the documentation is not 100% clear just ask. 
I know everything there is to know about how Rubble works and so I tend to forget to put 
"obvious" stuff in the docs. A reminder that I forgot something is more a help than a hindrance.

==============================================
Changelog:
==============================================
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
