
Rubble: After Blast comes Rubble

==============================================
Overview:
==============================================
Rubble is a raw generator.

Pros:
	No need for an external runtime
	Configurable directory structure
	Addons may override vanilla files
	Full parser/lexer, not just a bunch of regexes
	Minimal impact, touch only the minimum amount of files
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

Why another raw generator? Well I really liked Blast but it was far to inflexable and (IMHO) generated ugly output. This led me to write GoBlast, a slightly more flexable, Blast compatible, raw generator. That worked for a while but I soon realized that while GoBlast got the job done the design/templates forced by Blast compatibility was making my job harder than it needed to be. Rubble is the result of redesigning GoBlast from the ground up to better fit that way I do things.

==============================================
Install:
==============================================

Backup your "raw/objects" folder!
Copy all files/folders to your raw folder. DO NOT COPY ANYTHING to raw/objects!
Install your base (if any) to "raw/source/base" and addons to "raw/source/addons"

Now you are good to go! Documentation (as you have obviously discovered) is in the "raw/rubble" folder as is source code.

If you use OSX or linux, 32 bit binarys for these OSes can be found in the "raw/rubble" directory. If you want 64 bit binarys you can compile them yourself (or you can ask very, very nicly and I may do it for you ;) ).

Docs for all NCA4 script commands are in "raw/rubble/nca docs", the file "NCA4.nca" is syntax/code examples.

==============================================
Configure:
==============================================

Rubble allows you to change its directory settings via command line options. To see these options and their defaults run "rubble -h". 

If you want to use non-default settings more than once, the best thing to do is make a batch file something like this:
	@rubble -outputdir="./objects" -configdir="./source" -basedir="./source/base" -addonsdir="./source/addons"
	@pause

All directorys used by Rubble must exist (if they do not exist nothing bad will happen, Rubble will just quietly fail). SCRIPT templates may be used to check for (and possibly fix-up) proper directory structure.

==============================================
Changelog:
==============================================
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
