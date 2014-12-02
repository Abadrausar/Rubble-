
Rubble: After Blast comes Rubble

==============================================
Overview:
==============================================
Rubble is a raw generator eg. it takes "preraws" and generates valid Dwarf Fortress raw files from them.
Rubble is based on Blast but is not compatible with that utility.

Pros:
	No need for an external runtime, Rubble is a native application
	Configurable directory structure via command line options
	Addons may override vanilla files
	Full parser/lexer, not just a bunch of regexes
	Minimal impact, touch only the minimum amount of files
	Templates are designed to make formatting easy, both output and input files should be easy to read
	Support for embeding scripts in raw files, do advanced setup right in the raws
	Simple variable expansion, no more GET_VAR, works in any parse stage
	Many templates for registering objects and the like are replacements for vanilla raw tags allowing better formating
	Using scripting some files may be parsed but not included in the generated raws

Cons:
	Not compatible with any version of Blast
	No support for handling file name collisions in any way other than as an override
	No support for Blast namespaces (@, @@, and @@@)
	Variables are simple key:value pairs instead of Blasts more powerful system
	NCA (the scripting language) only has minimal documentation

Rubble comes prebuilt for 32 bit Windows, Linux, and OSX

==============================================
Why Another Raw Generator?
==============================================

I really liked Blast but it was far to inflexable, hard to install, and (IMHO) generated ugly output. This led me to write GoBlast, a easier to install, Blast compatible, raw generator. That worked for a while but I soon realized that while GoBlast got the job done the design/templates forced by Blast compatibility were making my job harder than it needed to be. Rubble is the result of redesigning GoBlast from the ground up to better fit the way I do things.

==============================================
Install:
==============================================

Backup your "raw/objects" folder!
Copy all files/folders to your raw folder. DO NOT COPY ANYTHING to "raw/objects"!
Install your base (if any) to "raw/source/base" and addons to "raw/source/addons"

Now you are good to go! Documentation (as you have obviously discovered) is in the "raw/rubble" folder as is source code.

If you use OSX or linux, 32 bit binarys for these OSes can be found in the "raw/rubble" directory. If you want 64 bit binarys you can compile them yourself.

Documentation for all NCA4 script commands are in "raw/rubble/nca docs", the file "NCA4.nca" is syntax/code examples.

==============================================
Configure:
==============================================

Rubble allows you to change its directory settings via command line options. To see these options and their defaults run "rubble -h". 

Rubble does not support a configuration file, if you want to use non-default settings more than once the best thing to do is make a batch file or shortcut. If you do not know how to use command line options, ask google.

All directorys used by Rubble must exist (if they do not exist nothing bad will happen, Rubble will just quietly fail).


==============================================
A Note About SCRIPT:
==============================================

The SCRIPT template (and its pre and post brothers) use a custom scripting language named NCA. The primary reason I used NCA in this program is that it is really easy to use, the main drawback is that I never got around to writing good documentation for the base language. I did write some documentation later but none of it is very complete or particularly good. On the bright side all of the commands are fully documented, and there is a good bit of code in "source/base/base_templates.txt" that you can use as an example.

All the command documents are in the "raw/rubble/nca docs" directory as well as what language documentation exists. The command documentation is in the .txt files and the language documentation is in the .nca file.

==============================================
BUGS:
==============================================

The generic animal mats (in the base) do not appear to work 100%, some things like meat are still not generic. Sometime I'll look into it, until then it may be better not to use that setting.

==============================================
Changelog:
==============================================
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
