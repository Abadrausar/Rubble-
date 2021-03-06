
There are some tutorials in the file "How To Rubble.txt", I suggest you read at least the first two as they explain the basics in some detail.

Despite the name this file deals with some very advanced topics! Basically anything that needed to be documented but didn't deserve it's own file got a section here. If you do not understand something just skip it, chances are it is written for advanced users about a feature added to solve a rare problem you have no need to solve :p

That said you should at least skim over everything so you have some idea what is possible.

Note that some of the things described here are added by the "Libs/Base" addon and are not hard-coded parts of Rubble. "Libs/Base" is commonly assumed to be an integral part of Rubble, it is safe to assume that it's features are always available.

==============================================
Rubble OS Specific Information:
==============================================

On Windows:

On windows Rubble should "just work", everything was originally written and tested on a 64 bit Windows 7 system, so if Rubble works anywhere it should work there :)

No specific issues are known.

The default browser startup script should start IE, if you don't like IE (like 90% of the right-thinking world) you will have to modify the script.

----------------------------------------------
On Linux:

I recently setup a a Linux VM, did some Rubble testing, and found a bunch of issues *sigh*.

To get anything working at all you need to:
	All the binaries need to be marked "executable" (chmod +x) and moved to the same directory as the Windows binaries.
	The browser startup script needs to be marked "executable" (chmod +x)

After that things *should* work fine.

If the web UI does not automatically start your browser you will have to modify the browser startup script (comments in the default script list some options).

----------------------------------------------
On OSX:

It is impossible for me to do any OSX testing (and I have never used OSX for anything but the most trivial tasks on a borrowed computer).
If you want to know how Rubble works on this platform either try it yourself or donate some hardware :)

At a minimum you will have to move the binaries to the same directory as the Windows binaries (do you have to mark them executable like on Linux? OSX and Linux are both UNIX clones).

There is no default browser startup script for the web UI on OSX. It will try to use the Linux script, but this is unlikely to work.

==============================================
AXIS VFS:
==============================================

For a variety of reasons Rubble uses AXIS VFS for (almost) all file access.
Most places where a path is given the AXIS syntax is used, so you will need to know at least a little so you can read these paths.

AXIS uses special (OS-independent) multi-part paths. The first part is one or more colon separated location IDs. You can think of location IDs as drive letters, they each specify a different location on the file system. The second part is the traditional slash separated path.

Examples:
	df:data/init/d_init.txt
	out:objects/entity_default.txt

In some cases it may be possible to specify multiple location IDs on a path, in that case just put them one right after another.

Example:
	addons:dir:addonlist.ini

One thing to keep in mind: Relative paths (paths containing "." or ".." elements) are utterly illegal in AXIS, so don't do that.

The only paths that are not handled by AXIS are:
	dfdir command-line/config option
	outputdir command-line/config option
	addonsdir command-line/config option

And even these paths may use the AXIS syntax (with a reduced set of location IDs):
	dfdir may use the "rubble" location
	outputdir may use the "df" and "rubble" locations
	addonsdir may use the "out", "df", and "rubble" locations

You may still use OS paths (which may be relative) for these three settings, but you do not have to if using the AXIS syntax would be easier.

Location IDs:
rubble				The current working directory, most of the time this is where the Rubble binary is.
df					The dwarf fortress directory, defaults to ".."
out					The output directory, by default "df:raw"
addons:dir			The addons directory, by default "rubble:addons"
addons:zip			A composite, read-only, directory made up of the contents of all the zip files in the addons directory

==============================================
File Extensions Used by Rubble (and Their Meanings):
==============================================

Rubble generally uses file extensions to determine file type, these associations are hardcoded (but may be overridden by clever use of file tags in loader scripts).
Most of these file types have a dedicated section in this document.

Extension	Tags				Use
.txt		RawFile				Assumed to be parseable raws that ere to be written to "out:objects".
.rbl		RawFile, NoWrite	Same as ".txt", but not written out.
.load.rex	LoadScript			A loader script, run just after addons are activated, but before incompatibilities are checked.
.init.rex	InitScript			An init script, run just before generation (even if their addon is not active).
.pre.rex	PreScript			A pre script, run just after init scripts.
.post.rex	PostScript			A post script, run just after generation.
.tset.rex	TSetScript			A tileset install script, run whenever a tileset is applied (just after post scripts).
.tset		TSetFile			A tileset definition file, not parsed!
.zip		<no tags>			A zipped addon, treated like a directory, see the section about the addon loader.
.zip.b64	<no tags>			A base64 encoded zipped addon, same as ".zip".
.webload	<no tags>			A webload file, contains the URL of a zip file to download.

There are some files that have hardcoded names. You may only have one of each of these for each addon (obviously).

Name			Use
addon.meta		Addon meta-data file, contains information about the addon.

==============================================
Template Prefixes Explained:
==============================================

When a template is run is determined by a one character prefix (or lack of a prefix).
There are three parse stages and four prefixes.
Prefix	Meaning
!		Preparse, a template with this prefix will run in the preparse stage.
none	Parse, a template with no prefix will run in the parse stage.
#		Postparse, a template with this prefix will run in the postparse stage.
@		Any, a template with this prefix will run in the earliest possible parse stage.

@ templates are usually run in preparse but could be run in any stage depending on when it is parsed (parsing can be delayed by nesting in a later parse stage template for example). This prefix is used for templates that may be run in any stage, if you need to control when these templates are parsed use !ECHO, ECHO, or #ECHO.

The template prefix is part of the template's name, you cannot change a template's parse stage by calling it with a different prefix (unless you rename the template)

==============================================
User Template Parse Sequence:
==============================================

The following details how templates created with the !TEMPLATE template are parsed.

	If the last param is "..." then append the params of the last template call to this call
	All leading and trailing whitespace is striped from the params
	All variables are expanded in the params
	The template's body is loaded and these transforms are run on it:
		Text of the forms $<NAME> or ${<NAME>} that is not in a child template is replaced with the value of the variable <NAME>
		Text of the forms %<NAME> or %{<NAME>} is replaced with the value of the param named <NAME> or the default value specified for that param
	The template body is passed to the stage parser and the result is returned

==============================================
The Configuration File:
==============================================

Rubble allows you to change its settings via command line options or a config file. To see these options and their defaults run "rubble -h" and "rubble_web -h" (each one has a different list).

ALL command line options may also be specified in the config file.

Rubble tries to read the file "./rubble.ini", if this does not fail Rubble will load settings from here before processing command line options (command line options always take precedence).
Example "rubble.ini" (using some of the defaults):
	[rubble]
	dfdir = ..
	outputdir = df:raw
	addonsdir = rubble:addons
Duplicate keys are just fine, they will act pretty much exactly like duplicate options on the command line.

Most users will never need the configuration file.

==============================================
Rubble Web UI Customization:
==============================================

The Browser Startup Script:

To make things easy you can create a batch file or shell script named "./other/webUI/browser" that starts your web browser, the server will try to run this file when it starts. This file will be passed the URL for the main menu as an argument, the script should pass this argument through to the browser (so the page opens automatically).

The browser start script file extension MUST be listed in your PATHEXT environment variable (on windows) or have no extension (Linux, OSX). Windows users who do not know what "PATHEXT" is can assume ".com", ".exe", ".bat", or ".cmd" are usable.

My browser file is named "browser.bat" and contains:
	@"C:\Program Files (x86)\SRWare Iron\iron.exe" "%1"
You will obviously need to use something different if you are not using a default SRWare Iron install on Windows x64 :)

Rubble comes with some default browser files, so things *may* work just fine right out of the box. The included files should work on most Linux or Windows systems (Linux users will have to mark the script executable).
OSX users are more or less on their own, but under the hood OSX is a lot like Linux (they are both UNIX clones), so if you provide the proper browser path the Linux script should work with minimal modification.

Linux users will have to do a bunch of extra stuff to get Rubble working, see the "Rubble OS Specific Information" section of this document.

----------------------------------------------
User Configuration Variable Defaults:

During normal generation the configuration variable editor will get default values for the variables from the addon.meta files. This works great except for one thing: While the defaults are relatively sane most advanced users will quickly get tired of changing them to their favorite settings :)

The solution to this is simple: Create a file named "rubble:userconfig.ini" and put values to override the defaults there. This file uses the same format as "out:genconfig.ini", so you can just copy the values you want from that file.

(The above two paths are in the AXIS syntax, see the section on AXIS VFS for details)

----------------------------------------------
HTML Template Customization:

When the web UI server is run it tries to load HTML templates and CSS files from "./other/webUI/", if it does not find the file it is looking for it tries to write the default version of the file out to that directory.

If you want to customize the look of the UI modify these files, from then on the server will use your files. (templates for the generation report are also here, you can tell which are which because the report templates end in "_static")

Each page template receives raw data about whatever it is supposed to do, in the default templates this data is formatted as JSON and most of the page content is generated client-side by JavaScript, but it is possible to generate the page directly on the server by using the template system for everything. In some cases doing so will cut the page size considerably as some pages receive the parsed version of the entire addons directory! (that's over 7mb of data as encoded!)

Keep in mind that the way the server interacts with these files is hardcoded! Each template receives a specific set of data, and must act a certain way (as observed by the server).

If you wish to use files beyond the default templates (for example images or extra HTML and CSS) the contents of the "./other/webUI/extras" directory are served as "/extras/<filename>". None of the default templates make use of this capability.

For more details on how these templates work see "golang.org/pkg/html/template/".
For details about exactly what Rubble passes each page you will have to look at the source code ("./other/src/rubble" and "./other/src/rubble/interface/web", specifically the files "meta.go", "addon.go", and "main.go" will be the most helpful).

If anybody comes up with a nicer set of templates (or some better CSS) I want to see them!
The current look isn't all that great, so this is your chance to get listed on the about page :p

==============================================
.webload files:
==============================================

It is possible to order Rubble to download a zipped addon from a web server when it is run, these addons are downloaded only if their size has changed from the last time they were downloaded (so auto-update for addon packs is a thing now!).

AFAIK if you use the correct URL it should be possible for Rubble to read addons directly from DFFD!

To make use of this ability simply create a file with the extension ".webload" (the file name is used as the addon name), the file's contents are stripped of whitespace and used as a download URL.

URLs used for web loading MUST be complete, eg "www.example.com/somefile.zip" will not work, use "http://www.example.com/somefile.zip".

.webload files are only allowed in the root of the addon tree (basically exactly the same rules as normal zipped addons).

==============================================
addon.meta:
==============================================

Addons may include a special file named addon.meta that contains extra information about the addon. This file is optional but HIGHLY recommended.

The meta file is actually a Rex script, but please don't abuse this by running arbitrary code as this script MUST return a special value (and load scripts are more flexible). The only reason a script is used is so that I can get complex parsing for "free".

Example addon.meta file:
	# Missing fields are set to their defaults.
	# If the addon does not have a meta file then all fields will be default.
	
	# This file MUST return a rubble:addonmeta indexable.
	
	# String fields generally have leading and trailing white space stripped.
	
	<rubble:addonmeta
		lib=false
		
		format="txt"
		
		header=`An example addon, nothing special about it.`
		
		description=`
This addon does little of interest,
Actually it is totally made up.
`
		
		activates="Libs/Example;Some/Addon"
		
		incompatible="Other/Addon"
		
		vars=<map
			CONFIG_VAR="default value"
			OTHER_CONFIGVAR=<rubble:addonmetavar name="Use some option?" val="YES" choices="YES" choices="NO">
		>
	>

Most of the fields are not used directly by Rubble, but may be used by various interfaces.

The most important field by far is the "activates" field. This field specifies addons that are required for this addon to function. Addons listed in the activates field are activated automatically if this addon is active, basically an easier to use version of rubble:activateaddon.

Closely related to activates is "incompatible". This field allows specifying addon that are incompatible with the current one. In most cases it is better to insert special case code to make your addon compatible, but sometimes this is not possible.
If an addon is active that is listed in another addon's incompatible field Rubble will abort (with a message telling you which addon caused the problem).

The "libs" field is also very important. Any addon that has lib set to true will not allow activation in the normal manner eg it will ignore activation by the -addons command line option, the addons config file key, and any entries in addon list files (when an addon list is generated lib addons will not even be listed). Such "lib" addons may only be activated by a load script or by having an entry in the "activates" meta key of an active addon.
Lib addons are good for libraries and other addons that do not directly effect the user, as they are generally hidden from view and managed automatically.

The "vars" field is a special case in many ways, first it needs to be an indexable mapping variable names to default values, second it may take a special kind of default value that gives you more control over possible values.
Most of the time the default value is interpreted as a string, but if you pass an indexable with the special rubble:addonmetavar type it will allow you to specify how the var is used. For now you can limit the possible choices to a specific set and/or specify a user friendly name.

"header" and "description" are descriptions and documentation for the edification of the user. "header" is a one line description of the addon, it should be short and to the point. "description" may be as long as you wish, put everything a user will need to know to use an addon here, if header already provides sufficient description leave this field empty.

"format" is used to specify how the value of "header" and "description" are formatted. Currently anything other than "html" is assumed to mean pre-formated plain text.

"activates", "incompatible" and "lib" are the only fields that can effect Rubble's behavior, all of the other fields are provided for use by interfaces and the generation report. Current the only interface that uses the meta data is the web UI.

==============================================
Tileset Addons:
==============================================

Tilsets are handled as a special extra generation step. After the raws are generated and just before they are written to disk Rubble collates all ".tset" files and uses the information in them to apply tileset information to the generated raws.

.tset files are normal raws with some minor differences:
	No "OBJECT" tags are needed, all objects can go in the same file.
	Any tags that do not deal with tile or color data can be left out.

For example a .tset file for "creature_amphibians.txt" would be:
	[CREATURE:TOAD]
		[CREATURE_TILE:249]
		[COLOR:2:0:0]
Simple.

Note that while .tset files use raw syntax they are NOT parsed by Rubble! This is due to the fact that they cannot be parsed when applying a tileset to existing raws, so to keep things consistent they cannot be parsed at all, ever.

In general all tileset information will be in a single file (as there is no point in having loads of little files). The .tset file for vanilla DF comes to ~95 KB, so not too large to pack into one file without issues.

.tset files are generally automatically generated, to make things easy every time you generate the raws Rubble will create one in "out:current.tset" using the information from the current set of addons. If you already have raws for a tileset simply place them in a temporary addon and generate this addon by itself.

If you want to change the tileset of an already generated raw set rubble provides two ways of doing this:
	Via the web UI, this is very easy to use, just select your region and addons to load the tileset from.
	Via the command line: `rubble -tileset="<region name>" -addons="<the name of the tileset addon(s)>"`

Scripts with the extension ".tset.rex" are run as part of generation AND when a tileset is added to a pregenerated raw set. These scripts are functionally identical to post scripts, and provide an easy way to install tileset images or make any needed init changes your tileset may need.

At some point I may extend this system to cover certain settings in the init files (the tile numbers in d_init), but for now if you need to modify these settings (and almost every tileset does) you will have to use tileset install scripts.

Due to the way init changes must be handled most (if not all) tileset addons will still require "Libs/Base", just do this in the addon.meta file as normal.

If you need (or want) to use a template in a .tset file you will have to parse it by hand. To parse these files you will need to use rubble:stageparse (with an explicit parse stage) from the tileset install script. Multiple calls to the stage parser may be required depending on what templates are used.
Please note that the base templates will NOT be loaded when applying a tileset to existing raws! (and many of these templates depend on thing that would not be present anyway) If you don't want unhappy users test your tileset addon both during normal generation and in independent application!
It would also be possible to make Rubble parse these files by playing with the file tags some, but DO NOT DO THIS! It will not work when applying a tileset to existing raws!

One interesting little "feature" of the way this system works is that it will work on any raw set, even ones that were not generated by Rubble! (also, when applying to existing raws the headers Rubble adds to generated files will be corrupted, but that is only cosmetic and so therefore nothing to worry about)

==============================================
The Rubble DFHack Lua Startup Script System:
==============================================

I have noticed that many mods that require a DFHack script use the script in one of two ways: either it is run once when the save is loaded or it is run when a reaction is completed. Reaction scripts are already handled by the DFHACK_REACTION and DFHACK_REACTION_BIND templates, but there was no good way to handle startup scripts.

Rubble's solution to the startup script problem is simple: it loads and evaluates scripts from raw/dfhack, once, when the save is loaded. These scripts are not commands (they cannot take parameters), and they do not clutter up the command list with things that will only be run once. Each of these startup scripts is isolated, if one misbehaves and dies the loader just goes on to the next with no interruption.

To load a script with this system simple put any .lua file into out:dfhack, it's that simple. To make it even easier you can just use the rubble:dfhack_loadscript script command, it installs a script by name automatically.

For convenience the loader provides a special value (rubble.savedir) that contains the path to the current region.

When working on an addon you may wish to reload the scripts and modules, to do this simply run `:lua rubble.reload_scripts()`.

==============================================
The Rubble DFHack Lua Pseudo Module System:
==============================================

DFHack can now load commands from the save directory, but if you need to load a module you still need to install it globally, "pseudo modules" are my solution for this issue.

Pseudo modules consist of three parts:
	The module file
	The loader
	Client scripts

The module mostly the same as with normal modules, the only difference is instead of:
	local _ENV = mkmodule("example")
You have this instead:
	local _ENV = rubble.mkmodule("example")
Simple.
(remember to `return _ENV` at the end of the module!)

The loader is simple and automatic, just put any .mod.lua file into out:dfhack. This works exactly like startup scripts (just with a different file extension) so you can even use rubble:dfhack_loadscript if you are feeling lazy.

Clients are mostly the same as with normal modules, the only difference is instead of:
	local example = require "example"
You use:
	local example = rubble.require "example"
Or you can use:
	local example = rubble.example
Or, of course, you could just use the full name everywhere :)
	rubble.example.someFunc()
The first method (rubble.require) is recommended to ensure compatibility with possible future changes.

Obviously, due to the way they are loaded pseudo modules will only work when a save is loaded, but since they are only for use with region-local scripts this is not an issue.

Pseudo modules have a special value normal modules lack:
	"_name" contains the module name as passed to rubble.mkmodule, this is mostly for internal use.

When working on a module you may wish to reload the scripts and modules, to do this simply run `:lua rubble.reload_scripts()`.

For convenience the loader provides a special value (rubble.savedir) that contains the path to the current region.

The following names are reserved and should not be used as a module name:
	savedir
	load_script
	load_module
	mkmodule
	require
	reload_scripts

==============================================
Patches:
==============================================

Rubble provides (very) basic support for applying diff patches to files (or any arbitrary chunk of text really).

This ability is primarily for use by bugfixes and minor updates and is not meant as a general modding tool, but there may be a few case where it will come in handy. The primary problem with Rubble's patch support is that it is fragile, eg even relatively small changes can stop the patch from applying. This behavior is intended.

In general it is better not to use patches, save them for use in loader scripts and even there only use them to apply small changes.

Keep in mind that when applying a patch from a loader script you need to make sure the script only runs once! See the section on loader scripts for details.

(see the rubble:patch script command and the !PATCH base template)

==============================================
Tweak Scripts:
==============================================

Tweak scripts are the most important thing to happen to Rubble since version 1. They allow you to run arbitrary scripts before or after generation. These scripts can change the raws, allowing you to "tweak" the generated output.

Pre tweak scripts have the extension ".pre.rex".
Post tweak scripts have the extension ".post.rex".

Tweak scripts are written in Rex not Rubble code, so Rubble templates are not available except through rubble:calltemplate or rubble:stageparse.

Most tweak scripts will make use of the commands in the df:raw namespace to walk the raw files and make changes as needed.

As for writing your own tweak scripts... It would be a very good idea to read the tweak scripts included with Rubble, if you do not understand those scripts then you will have trouble writing your own.

==============================================
Initialization Scripts:
==============================================

Init scripts are very limited, but also very useful. Basically they are Rex scripts that ALWAYS run, even when their containing addon is not active.
Init scripts have the extension ".init.rex".
Use init scripts to do setup that may be shared by multiple addons and HAS to be done. For example an excellent use for a init script is adding namespaces and global variables shared by a bunch of addons in an addon pack.

Another fine use of init scripts is "always active" templates, basically you add a template stub in an init script an the override it with a full version in a pre script. This way clients of your addon can just assume that your addon's template are present allowing you to skip checking to see if it is active. 
(please note that you really only want to do this within an addon pack so that you can be sure the other addon is installed)

Init scripts have a very narrow use-case and should be used with care.

Init scripts are still loaded as "user data" if their containing addon contains any parseables!

==============================================
Loader Scripts:
==============================================

Loader scripts are a lot like initialization scripts, just that they run earlier.

Loader scripts only do two things that init scripts can't do better: they can modify the activation state of addons and they can modify the files of a specific addon without effecting other addons. Since loader scripts run just after addons have been activated but before the active file list has been generated you can examine what addons are active and change this list if need be. This is particularly useful for automatically activating addon dependencies.

Like init scripts, loader scripts are global and pay no attention to their containing addon.
Loader scripts are only run during normal generation cycles.

Loader scripts make a great place to apply patches and force activate dependencies, for most other uses init scripts are better (for a better way to force dependencies see the section on addon.meta).

One potential problem is that loader scripts can (and will) be called more than once per session (depending on the Rubble interface used).
This can be a BIG problem if you are doing something that may only be done once. If this is a problem put all such items into a separate file and add `(rubble:gfiletag (rubble:curentfile) Skip true)` somewhere in the file.

Remember: The loader scripts are run more than once for a reason! Addon state may change significantly between runs, so do not skip any script that you do not absolutely have to.

(skipping loader scripts may not work anymore!)

Luckily these problems are really rare in practice (none of the default addons had any trouble), as there are very few things that NEED to be done from a loader script that cannot be done multiple times. About the only thing that can cause problems is applying patches.

==============================================
File Tags:
==============================================

Rubble keeps track of file type and other attributes via "file tags", simple boolean flags that can be set at will.

File tags default to false and can be changed at any time, any string can be used as a tag ID.

To modify and/or read file tags the rubble:filetag command is used:
	To read a tag:
		(rubble:filetag "file.name" "TagID")
	To set or clear a tag:
		(rubble:filetag "file.name" "TagID" true)
		(rubble:filetag "file.name" "TagID" false)

It may also be possible to use GenII  to modify these tags:
	[rubble:addonstbl "Addon/Name" Files "file.name" Tags "TagID" true]
This may or may not work after generation has started, all the normal GenII cautions apply!
GenII is the Rex Generic Indexable Interface, it gives low level access to some of Rubble's internals but has a high foot-shot potential, be careful!

Hard Coded File Tags:
	Tag ID			Use
	LoadScript		Marks the file as a loader script.
	InitScript		Marks the file as a init script.
	PreScript		Marks the file as a pre script.
	PostScript		Marks the file as a post script.
	TSetScript		Marks the file as a tileset setup script.
	TSetFile		Marks the file as a tileset information file, these files are raw text but are NOT parsed!
	RawFile			Marks the file as parseable raw text.
	GraphicsFile	Used with RawFile, specifies that the file should be written to "out:graphics" instead of "out:objects".
	TextFile		Used with RawFile, specifies that the file should be written to "out:objects/text" instead of "out:objects".
	NoWrite			Used with RawFile, marks the file as not to be written to the objects folder.
	Skip			The file is to be ignored in any future passes, this tag is automatically cleared at addon activation!

LoadScript and InitScript are only used with certain global files, you will never see these tags when using rubble:filetag, only rubble:gfiletag will show these tags.

==============================================
The Rubble File Encoder:
==============================================

Rubble comes with a simple command line utility for encoding files for posting on the forum. This utility compresses and encodes files so that they are forum postable (eg they contain no unprintable characters).

To encode or decode a file just drag and drop it on the encoder binary (if you want to use the command prompt just run "encoder filename").

The encoder will examine the file extension to decide what to do:
	Extension		Action
	.zip.b64		Base 64 decode
	.zip			Base 64 encode
	.b64			Base 64 decode and DEFLATE uncompress
	anything else	DEFLATE compress and base 64 encode
The output file will have ".b64" added/removed to/from it's extension as appropriate.

Encoded zip files can be decoded by any base 64 decoder (base 64 is used for many things so it is not hard to find decoders), but other encoded files will also need to be decompressed (they are raw DEFLATE compressed blocks, so it may be a little hard to find a ready made tool for processing them). While it is possible to decode one of these files using third party utilities, there is no real need as the encoder application can both decode and encode files without issues.

In any case Rubble can read these files natively, so there is no need for users to decode them before use (but there is no reason why you can't, after all this utility could be very useful for non-Rubble applications).

Please note that for files that are already compressed (for example zip files and most images) the encoder will inflate the file size, this is because base 64 encoding adds 2 bits per byte. In most cases this is fine as this encoding is designed for use with ASCII text and other uncompressed data, but don't be surprised when your zip file gets ~25% larger (it would be exactly 25%, but I add a newline every 80 characters to enhance clarity, so it's a bit more).

There are two basic ways this encoded text can be used:
	As encoded zip files.
	As encoded files uncompressed by a load script or the !DECOMPRESS template.

Encoded zip files (.zip.b64) can be read directly by Rubble, use them exactly like normal zip files (just put them in your addons folder).

Other encoded files need to be uncompressed by a load script before use, most of the time the encoded file is packaged into a simple script as a string. To decode a string use the rubble:decompress script command, this will decode and uncompress the string, if the file was not stored in the script as a string then you will need to load the file first.

There is another, easier, way to decompress a encoded file, just use the !DECOMPRESS template. Sadly this template is far less flexible that using load scripts, but if you would rather not use scripts see the base template docs.

Loading then decompressing a file can be done like this:
	var content = (genii:bytes_string [(rubble:fetchaddon "Addon/Name") Files "file.name.b64" Content])
	(rubble:newfile "Addon/Name" "file.name" (rubble:decompress [content]))
This will result in the addon having two files, one compressed (that was there before) and one uncompressed (which the above code added)

The following would replace the file's contents with the uncompressed version:
	var addon = (rubble:fetchaddon "Addon/Name")
	var content = (genii:bytes_string [addon Files "test compressed.txt" Content])
	[addon Files "test compressed.txt" Content = (genii:string_bytes (rubble:decompress [content]))]
Note that if you choose to do things this way you cannot use the .b64 extension unless you want to fix the file tags!

All in all it is much easier to store the file in the script as a string, as in that case you only need to do half as much.

==============================================
The Addon Loader:
==============================================

Rubble has a recursive addon loader that can read addons from zip files or directories.
What this means is that you can group addons in directories to have something like an addon tree.

The following rules determine how addons are loaded:
	Zip files may only be used in a root addon folder.
	Addons in a zip always override addons in a directory.
	Addons are loaded in alphabetical order except as otherwise specified.
	Zip files and directories are identical as far as the loader is concerned except as otherwise specified.
	If an addon contains one or more directories, the directories are loaded as child addons.
	An addon without parseable files (.txt, .rbl, .rex) is not loaded (but any child addons are).
		An init script (.init.rex) is not considered a paresable file (but such a script will still run), and the same goes for load scripts (.load.rex).

The following rules determine how addon files override each other (if there is a name collision):
	Addons in zip files will always override addons in directories.
	Addons are loaded in alphabetical order, addon "aaa/zzz" will override "aaa/hhh" but not "bbb/aaa"

If you have readmes or the like in an otherwise empty addon it is recommended that you do not give them a file extension. In any case you do not want such an addon appearing in the addon list, so do not use any of the "parseable" extensions (.txt, .rbl, .rex).
A good alternative extension is .text or maybe .me (as in "read.me"), but I prefer no extension at all.

All addons are loaded into memory, but only the active addons are processed.
An addon is active if:
	It's name is listed on the command line in the -addons option
	If it's entry in addonlist.ini equals "true"

If any addon names are specified via the -addons option then addonlist.ini is never read (but it is updated)

So why do you (the modder) need to know all this? To put it simply you need to know what is permissible for structuring your addon(s) so that you can choose an optimal method that makes it easy on both you and your users. 

Do you want to use multiple small addons or just one large one? Do you want all of your mod's addons in one large directory, or do you want to group related things together? These decisions are important part of deciding how you want to package your mod and how you want the addon list presented to users.

==============================================
Multi-Threaded Scripts:
==============================================

It is possible to write multi-threaded Rex scripts, but this capability is not generally useful, as the vast majority of scripts do not run long enough to make it worth while.

That said if you do have a task that would benefit from being run in parallel the ability to do so is there.

I will not try to tell you how to write concurrent code, I will just point out pitfalls and problems.

Unless you included the -threads command line parameter or put "threads=true" in "rubble.ini" Rubble will (most likely) only have a single OS thread to work with, so using the script threading commands will actually slow things down!

Some data exported to scripts in not thread safe! Standard indexables (map, smap, array, sarray) should work fine, as will normal global values, but GenII values and most custom indexables will not. Beware race conditions!

If a script command requires some kind of custom value as a parameter (for example the rand commands or the axis commands) then the value in question is probably not thread safe.

Many of the commands in the "rubble" module use the "rubble:state" variable internally, this value is not thread safe!

Basically you should only do calculations on static data, any IO or interaction with Rubble should be done from a single thread.

I have only ever used the Rex thread API once. The case I used it for was processing the DF language files (specifically I was expanding them so there was a different word for each form (verb, noun, etc)). Building the symbol file was taking far too long, so I made a new thread to generate each symbol. This worked fine and was just under 4x as fast (on a quad core processor).

If you have a data processing task that is far too slow and could be done in parallel, take a look at the commands in the "thread" module, they may be just what you need.

==============================================
Other Random Information:
==============================================

If your addon name starts with "__CONFIG_VAR_" or it is named "__EditConfig__", "__Generate__", or "__REGION__" then it will break the web UI. Do not use these names! (stupid HTML forms...)

If you want to allow an item, reaction, or building with a minimum of fuss, register it as class "ADDON_HOOK_MOUNTAIN" or "ADDON_HOOK_PLAYABLE". Specific addon hooks also exist for the other entities as well as an ADDON_HOOK_GENERIC that every entity has.

If you are making a template that takes a variable number of params and you want to implement it via a different template you may need the ... param.
Example:
	{!TEMPLATE;E;
		{ECHO;...}
	}
If Rubble finds a template call where the last param is ... it appends the params from the previous template call to the current template's params. This trick is good for aliases and the like.

The Rex script command documentation is automatically generated every time I build Rubble, this make it easy for me to keep everything up to date, but the formatting leaves much to be desired.
The problem is that the app that generates the docs is made for go programmers, not Rex programmers, so a lot of information that has nothing to do with the commands makes it's way to these docs.
Sometime when I'm really bored I need to write a program to automatically reformat these files ;)
