
Rubble: Modding Made Easy!

==============================================
Overview:
==============================================

In Rubble many little (and not so little) mods (called addons) are combined and run through a special parser to generate standard Dwarf Fortress raw files.

These addons are, by construction, able to be assembled in many configurations with minimal direct internal dependencies. Most of the time an addon is completely standalone or dependent on only a few common items/materials present in the vanilla raws.

The beauty of this system is two fold: on one hand users have many choices and can construct their own private version of DF with minimal effort, on the other hand modders can make use of the power of Rubble to make installation of their mod automatic and use the template/scripting system to automate most if not all of the more repetitive parts of modding.

Not all addons are made for users, whole groups of addons are made specially for modders to help with testing and automating certain common tasks (not to mention library addons full of ready to use templates and scripts).

Rubble has been in continuous development since mid 2013, and many changes and improvements have been made since the first version (which kinda sucked :p). All of my mods have been made to use it, so I have extensive experience with Rubble modding and along the way fixed most of the bugs and streamlined things as much as possible both for modders and users.

I hope Rubble fills your needs for a general content installer and modding tool (and if not post your suggestions!)

.. contents::

==============================================
Where to Find Help:
==============================================

All documentation is in the "other" directory, this includes tutorials, template documentation, script documentation, and basic documentation of Rubble's internals.

General users have no need to read ANYTHING except this readme! Everything most users need to know (aside from things like how to install Rubble, which can be found in this document) can be found in the web UI addon description pages.

Most of the documentation assumes you know what you are doing as far as general DF modding goes, please do not read documentation that is above your level and then complain that "Rubble is too hard".
I do assume that modders will have at least some knowledge of programming (I apologize for this, as I know that not everyone does). If you do not know anything about programming most of the scripting stuff will be hard (if not impossible) to follow, but if you stick to using templates you will be fine (scripting is only needed for advanced stuff anyway, most people will never need it).

For examples of what all is possible with Rubble see the included addons, they cover a broad range of possibilities. Sadly comments explaining what is being done and why are often lacking, so if you cannot figure out why something was done a certain way just ask.

Modders will want to read the following (in most-important-first order):
 [Texto del enlace aqu�](direcci�n.enlace.aqu� "t�tulo del enlace aqu�")
 
	1.	[HowTo Rubble - Tutorials](https://github.com/Abadrausar/Rubble-/blob/master/other/HowTo%20Rubble.md "Here we have some Tutorials, read at least the first two!") Feel free to skip the others.
	[Rubble Base Templates - The template documentation for the "Libs/Base" addon, this stuff is VERY important!
	[Rubble Basics - Lots of stuff about addons and Rubble in general, some of this stuff is kinda advanced.
	[Rubble Libs Templates - The template documentation for the library addons
	[Everything in "Rex Docs" - The Rex documentation, feel free to skip this if you are not interested in scripting
	[The included addons - A little short on comments but full of great examples

==============================================
Install:
==============================================

If you have anything in your Dwarf Fortress directory that you want to keep back it up! Rubble will delete all your existing raw files, including creature graphics and the DFHack onLoad.init and init.lua files (and many other things)!
Anything outside of your "data/init" and "raw" directories should be safe, but just in case...

Delete or otherwise remove any old Rubble version you may have.
Extract the new Rubble to "<DF Directory>/rubble".
Install any custom addons you may have to "<DF Directory>/rubble/addons".
See the appropriate "Running Rubble" section below.

Now you are good to go!

DO NOT extract a new Rubble version over an old one! This can cause all sorts of hard to find problems!

If you use OSX or Linux, 32 bit binaries for these OSes are included. If you want 64 bit binaries you can compile them yourself, source code is in "./other/src" (along with basic build instructions).
To use non-Windows binaries they must be placed in the same directory as the Windows ones (you can delete the Windows binaries if you don't need them). OSX binaries are in "./other/darwin_386", Linux binaries may be found in "./other/linux_386".
For more help see the "Rubble OS Specific Information" section in "Rubble Basics".

==============================================
Running Rubble (Web UI, all OSes):
==============================================

The web UI is kinda ugly, but simple, functional, and easy to use.

This interface is the recommended way of running Rubble, as it always has full support for the latest ways of doing things.

To use the web UI all you need to do is start the server (rubble_web) and then point your browser to "http://127.0.0.1:2120" (by default), from there just follow the menus.

To make things easy you can create a batch file or shell script in "./other/webUI" named "browser" that starts your web browser, the server will try to run this file when it starts. This file will be passed the URL for the main menu as an argument, the script should pass this argument through to the browser (so the page opens automatically).
Rubble comes with default scripts for Windows and Linux that should "just work" on most systems.
For more help with the browser startup script see the "Rubble Web UI Customization" section in "Rubble Basics".

The web UI is the most advanced (functionality-wise) UI available. The command line UI can do most of the things the web UI can do, but it is limited by being non-interactive.

So far this UI is the only one that fully supports addon meta data.
A particularly visible example of this meta data in action is addon descriptions, if you click on an addon name (anywhere addon names are shown) it will take you to a small page with information about the addon! No more digging through a text file for addon details, you can even get a list of dependencies or look at the contents of individual files!

The server supports the rubble config file, run "rubble_web -h" for a full list of supported keys and command line options.

The UI's look and feel can be customized, see the "Rubble Web UI Customization" section in "Rubble Basics".

==============================================
Running Rubble (CLI, all OSes):
==============================================

It is HIGHLY recommended that you use the Web UI. It is generally much faster then doing everything by hand (plus you don't have to mess around with the command prompt, if you dislike that kind of thing)

The Rubble command line interface is fairly simple, doing non-interactive batch processing only.

For basic documentation on each command line option run "rubble -h".

Common Tasks:

To activate or deactivate a Rubble addon manually you may set it's entry in "addons:dir:addonlist.ini" to "true" or "false". If you just installed an addon it will not have an entry until Rubble has run at least once (after the addon was installed). 
If you want to run Rubble without generating anything (so as to update the addon list file) just run 'rubble -addonlist'.
It is possible to specify addons manually as well, just use -addons.

To generate with the last used addons (aka the settings in "addons:dir:addonlist.ini") just run Rubble with no arguments.

Here are a few example command lines to generate various configurations:
	Vanilla (ignoring the addon list, but not rubble.ini):
		rubble -addons="Base"
	
	Vanilla (ignoring all external settings):
		rubble -zapaddons -zapconfig -addons="Base"
	
	Vanilla with Phoebus tile mappings:
		rubble -addons="Base;User/Tilesets/Phoebus"

Some addons may allow additional configuration via "config variables", these are generally for advanced users and may be specified with the -config command line option.

Example (vanilla with 10 dummy reactions):
	rubble -addons="Base;Dev/Dummy Reactions" -config="DEV_DUMMY_REACTION_COUNT=10"

To switch tilesets on a world in progress use `rubble -tileset="<region name>" -addons="<name of tileset addon>"`.

For example to switch "region1" to ASCII mappings:
	rubble -tileset="region1" -addons="User/Tilesets/ASCII"

It is generally possible to regenerate the raws for a world, this is a task fraught with peril, do not attempt unless you know what you are doing!
The first step is to change the "addonlist.ini" file that is in the world's raw directory, make sure not to change it too much or you can mess up your world!
If you wish you can change "genconfig.ini" as well, but that is usually a bad idea.
Now for regenerating the raws, what follows is an example command line to do that for the save "region1":
	rubble -zapaddons -zapconfig -addons="df:data/save/region1/raw/addonlist.ini" -config="df:data/save/region1/raw/genconfig.ini" -outputdir="df:data/save/region1/raw"

==============================================
BUGS (and other issues):
==============================================

Known Issues (non-bug):
	
	You should not use the extension .txt for readme files, as this will cause Rubble to parse those files.
		See "Rubble Basics.txt" for more details.
	
	The scalemail armor from the "User/Warcrafter" addon is always marked "Foreign" in the equipment screen, you will have to explicitly assign it.
		IMHO this is better than sometimes not having all the armor pieces (which is what happened before).

Bugs:
	
	The Web UI server does not always shutdown immediately when told to.
		Sometimes when you click "Quit" the server does not shut down, but as soon as you close the exit page/tab/whatever the server closes. It acts almost like it was waiting for browser (which should not be the case). I have no idea what causes this, but it is intermittent and doesn't really cause any problems, so I am unlikely to dig too much either.
	
	Several vanilla DF bugs make themselves known:
		Adventure mode reactions do not always work the same as fortress mode reactions:
			"User/Warcrafter/Adventure" does not work quite right with "User/Tanning" because the adventure mode tanning reaction uses too much skin. This cannot be fixed by me.
		
		Reactions produce unusable gloves:
			The glove reactions are left in only because DFHack can be used to fix this bug, if DFHack is available then Rubble will automatically load a fix.

When making an error report please post the FULL log file! Posting just a few lines tells me almost nothing about the problem, after all I made Rubble write all that stuff for a reason :)
I cannot stress this enough! With Rubble the actual error message is only a small part of the information I need for tracking an error down! In particular the list of active addons is VERY important (and it tends to get cut from most reports because it is near the start of the log).

If any of the documentation is not 100% clear just ask. 
I know everything there is to know about how Rubble works and so I tend to forget to put 
"obvious" stuff in the docs. A reminder that I forgot something is more a help than a hindrance.

In the event that I cannot be contacted on the Bay12 forums (user name "milo christiansen"), my email address is:
	milo.christiansen (at) gmail (dot) com
Please wait 1-4 weeks before giving up hope, as my Internet access is VERY irregular (and I check my email less often then I check the forums).
