
Rubble: After Blast comes Rubble

==============================================
Overview:
==============================================
Rubble is a raw generator eg. it takes "preraws" and generates valid Dwarf Fortress raw files from them.
Rubble is based on Blast but is not compatible with that utility.

Pros:
	No need for an external runtime, Rubble is a native application
	Addons may override vanilla files
	Full parser/lexer, not just a bunch of regexes
	Templates are designed to make formatting easy, both output and input files should be easy to read
	Support for embeding scripts in raw files, do advanced setup right in the raws
	Variable expansion, no more GET_VAR, works in any parse stage
	Many templates for registering objects and the like are replacements for vanilla raw tags allowing better formating
	Using scripting some files may be parsed but not included in the generated raws
	Support for easily and quickly installing tilesets as addons

Cons:
	Not compatible with any version of Blast
	No support for handling file name collisions in any way other than as an override
	No support for Blast namespaces (@, @@, and @@@)
	Variables are simple key:value pairs instead of Blasts more powerful system

Rubble comes prebuilt for 32 bit Windows, Linux, and OSX

==============================================
Why Another Raw Generator?
==============================================

Blast has a lot of power and advanced fetures that let you do all kinds of things while still remaining compatible with most if not all other Blast addons. Rubble on the other hand is much more minimalist, Rubble allows you to make many things with a minimum of fuss but has far less in the way of compatibility hand holding.

Which is the best depends on exactly what you want, Rubble for simple and fast, Blast for maximum compatibility and power.

Basicly I just wanted something that allowed me to do most of the things I could do with Blast, but with less fuss.

==============================================
Install:
==============================================

Backup your "raw/objects" folder!
Copy all files/folders to your raw folder. DO NOT COPY ANYTHING to "raw/objects"!
Install your base (if any) to "raw/source/base" and addons to "raw/source/addons"

Now you are good to go! Documentation (as you have obviously discovered) is in the "raw/rubble" folder as is source code.

To activate or deactivate a rubble addon you may remove or add double leading underscores to the addon's directory, eg. __test_addon is disabled while test_addon is enabled.

If you use OSX or linux, 32 bit binarys for these OSes can be found in the "raw/rubble" directory. If you want 64 bit binarys you can compile them yourself.

==============================================
Configure:
==============================================

Rubble allows you to change its directory settings via command line options. To see these options and their defaults run "rubble -h". 

Rubble tries to read the file "./rubble.cfg", if this does not fail Rubble will load directory setting from here before processing command line options (command line options always take presidence).
Example "rubble.cfg":
	dfdir = ..
	outputdir = ./objects
	configdir = ./source
	basedir = ./source/base
	addonsdir = ./source/addons

If you want to regenerate the raws for a save just run 'rubble -dfdir="./../../../.."' in the save's raw directory.
	Keep in mind that using rubble will increase your save size by ~10MB (removing the two extra binarys will drop that amount by about half)

All directorys used by Rubble must exist (if they do not exist nothing bad will happen, Rubble will just quietly fail).

==============================================
BUGS:
==============================================

None known.

==============================================
Changelog:
==============================================
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
	Notepad++ UDF updated
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
	Changed ITEM_CLASS to be more flexable
	Changed SHARED_INORGANIC and SHARED_MATERIAL_TEMPLATE to automaticly support tilesets
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
	Stoped Rubble from mangling special chars

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
