
Rubble: After Blast comes Rubble

==============================================
Overview:
==============================================
Rubble is a raw generator like Blast but better ;)

Pros:
	No need for an external runtime
	Configurable
	More operating modes to better support different setups
	Addons may override vanilla files
	Full parser/lexer, not just a bunch of regexes
	Minimal impact, touch only the minimum amount of files
	Support for embeding scripts in raw files, do advanced setup right in the raws

Cons:
	Not compatible with any version of Blast
	No support for handling file name collisions
	No support for Blast namespaces (@, @@, and @@@)
	Variables are simple key:value pairs instead of Blasts more powerful system

Rubble comes prebuilt for 32 bit Windows, Linux, and OSX

==============================================
Why Another Raw Generator?
==============================================

Why another raw generator? Well I really liked Blast but it was far to inflexable. This led me to write GoBlast, a slightly more flexable, Blast compatible, raw generator. That worked for a while but I soon realized that while GoBlast got the job done the design/templates forced by Blast compatibility was making my job harder than it needed to be. Rubble is the result of redesigning GoBlast from the ground up to better fit that way I do things.

At this time Rubble is not finished, only a small fraction of the planed templates are ready and many old ones left over from GoBlast are still in use. This makes Rubble a little hard to use, look for this to improve with newer versions.

==============================================
Install:
==============================================

Backup your "raw/objects" folder!
Copy all files/folders to your raw folder. DO NOT COPY ANYTHING to raw/objects!
Install your base (if any) to "raw/source/base" and addons to "raw/source/addons"

Now you are good to go! Documentation (as you have obviously discovered) is in the "raw/rubble" folder as is source code.

If you use OSX or linux, 32 bit binarys for these OSes can be found in the "raw/rubble" directory. If you want 64 bit binarys you can compile them yourself (or you can ask very, very nicly and I may do it for you ;) ).

==============================================
Configure:
==============================================

Rubble allows you to change its directory settings via command line options. To see these options and their defaults run "rubble -h". 

If you want to use a non-default setting most of the time the best thing to do is make a batch file something like this:
	@rubble -outputdir="./objects" -configdir="./source" -basedir="./source/base" -addonsdir="./source/addons"
	@pause

All directorys used by Rubble must exist. If Rubble tries to read from a non-existant directory it may crash, in such a case no changes will be made to anything. If Rubble tries to write to a non-existant directory it will silently fail.
	
==============================================
Changelog:
==============================================
v1.0
	First version