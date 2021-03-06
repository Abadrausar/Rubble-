
This file contains the source for my planned "HowTo Rubble" tutorial series.
The tutorials in this file are marked up with BBCode, as it makes them easer to post on the bay12 forums.

========================================================

Tutorial 1: Basics

This tutorial introduces a few of the basic Rubble templates and explains how to create a simple addon.

In this tutorial I will be recreating the "User/Bonfire" addon.

The bonfire is a simple furnace that lets you create two kinds of burning objects, one that burns out fairly quickly and one that burns for quite a while.

First install Rubble as directed in it's readme.
Now create a new directory in the addons directory, call it something like "Bonfire Tutorial", this is your addon directory.

Let's start with the building. Create a new file in your addon directory and name it something like "building_bonfire_tutorial.txt", then paste the following code into that file.
[code]
[OBJECT:BUILDING]

{BUILDING_FURNACE;BONFIRE;ADDON_HOOK_PLAYABLE}
	[NAME:Bonfire]
	[NAME_COLOR:6:0:1]
	[DIM:1:1]
	{@BUILD_KEY;SHIFT_F;true}
	[WORK_LOCATION:1:1]
	[BLOCK:1:O]
	[TILE:0:1:240]
	[COLOR:0:1:6:0:0]
	[TILE:1:1:'#']
	[COLOR:1:1:6:0:0]
	[TILE:2:1:19]
	[COLOR:2:1:6:0:0]
	[TILE:3:1:30]
	[COLOR:3:1:6:0:0]
[/code]
There are three things of note in this file.
First there is no file header, Rubble adds this automatically, so renaming a file goes from an error-prone two step process to just a simple file rename.
Second you will note that the building tag is replaced by some weird thing that uses curly-brackets and has an extra parameter, this is a Rubble template. This particular template is a direct replacement for the raw tag of the same name, but it also has one other very important function. This template (as well as lots of other related Rubble templates) registers a building with an object class. A Rubble object class is a list of items/buildings/reactions/whatever listed under a name. By using other templates you can extract a list of a specified type of items from a class and do things with them, in this case automatically add building permissions to entities.
ADDON_HOOK_PLAYABLE is a special class that can be used for items, reactions, and buildings, it is part of a special group of object classes created by the ADDON_HOOKS template (which I will explain in more detail later).
Third the BUILD_KEY tag has also been replaced with a template, this template will automatically resolve key conflicts, so it another building already uses SHIFT_F it will choose the next open key (ALT_F in this case, unless it is in use as well). The second parameter specifies whether the key is for a furnace, so in this case it should be true (this parameter defaults to false, so it can be left off for workshops).

Now that we have a building we need some reactions. Create a new file and name it "reaction_bonfire_tutorial.txt", then paste the following code into that file.
[code]
[OBJECT:REACTION]

{REACTION;BONFIRE_BIG_START;ADDON_HOOK_PLAYABLE}
	[NAME:ignite large fire]
	[BUILDING:BONFIRE:CUSTOM_L]
	[REAGENT:A:5:WOOD:NONE:NONE:NONE]
		[PRESERVE_REAGENT]
	[PRODUCT:100:1:BOULDER:NO_SUBTYPE:INORGANIC:BIG_BONFIRE]
	[FUEL]
	[SKILL:SMELT]

{REACTION;BONFIRE_SMALL_START;ADDON_HOOK_PLAYABLE}
	[NAME:ignite small fire]
	[BUILDING:BONFIRE:CUSTOM_S]
	[REAGENT:A:1:WOOD:NONE:NONE:NONE]
		[PRESERVE_REAGENT]
	[PRODUCT:100:1:WOOD:NO_SUBTYPE:INORGANIC:SMALL_BONFIRE]
	[SKILL:SMELT]
[/code]
The REACTION template is almost exactly like the BUILDING_FURNACE template (except, of course, it replaces the REACTION tag instead of the BUILDING_FURNACE tag)

Of course those reactions won't work without the inorganic materials BIG_BONFIRE and SMALL_BONFIRE.
Create a new file (name it "inorganic_bonfire_tutorial.txt") and paste the following code into it.
[code]
[OBJECT:INORGANIC]

{SHARED_INORGANIC;BIG_BONFIRE;
	[STATE_NAME_ADJ:ALL_SOLID:bonfire]
	[DISPLAY_COLOR:0:0:1]
	[TILE:15]
	[ITEM_SYMBOL:15]
	[IGNITE_POINT:11000]
	[MAT_FIXED_TEMP:20000]
	[MELTING_POINT:NONE]
	[BOILING_POINT:NONE]
	[SPEC_HEAT:10000]
	[SOLID_DENSITY:10000]
	[IS_STONE]
	[NO_STONE_STOCKPILE]
}

{SHARED_INORGANIC;SMALL_BONFIRE;
	[STATE_NAME_ADJ:ALL_SOLID:firewood]
	[DISPLAY_COLOR:0:0:1]
	[TILE:15]
	[ITEM_SYMBOL:15]
	[IGNITE_POINT:11000]
	[MAT_FIXED_TEMP:20000]
	[HEATDAM_POINT:11000]
	[MELTING_POINT:NONE]
	[BOILING_POINT:NONE]
	[SPEC_HEAT:10000]
	[SOLID_DENSITY:10000]
	[IS_STONE]
	[NO_STONE_STOCKPILE]
}
[/code]
Wow, what's this SHARED_INORGANIC thing? That is one of the SHARED_OBJECT templates. SHARED_OBJECT (and the templates based off of it, like SHARED_INORGANIC) provides a mechanism for replacing or modifying objects from other addons. In this case it's not terribly useful, but using it every time you can is a good habit to get into. I will not go into detail about exactly what SHARED_OBJECT (or it's children) does, as that would take far too much time. I suggest you go read the documentation for the "Base Templates" for the whole story.

Now there is one last detail, "How exactly do I make my building/reactions usable?". To put it simply, you already have.
Earlier I mentioned the ADDON_HOOKS template, this very useful template is already in every entity in the "Base" addon. Basically ADDON_HOOKS adds calls to all the required templates to insert entity permissions for items, reactions, and buildings for the following classes:
ADDON_HOOK_GENERIC
ADDON_HOOK_<entity name> (For example ADDON_HOOK_MOUNTAIN or ADDON_HOOK_PLAINS)
and if the entity is playable (dwarves only by default) ADDON_HOOK_PLAYABLE
This means that your new addon should be immediately usable by the dwarves and any other playable races that may be added by other addons.
Of course for ADDON_HOOKS to work you need to have declared the item/building/reaction with the proper Rubble template, but as you have seen that is a simple matter for buildings and reactions. For items it is a little more complicated, but that is for a later tutorial.

========================================================

Tutorial 2: Custom Templates

This tutorial introduces the simplest aspects of custom templates.

I will not be making a full addon in this tutorial, rather I will just make a template and a few reactions.

Say you want to make a workshop that allows you to produce metal weapons without needing an anvil. The first problem you will encounter is that writing enough reactions to make such a workshop useful is a tremendous amount of work. Using a Rubble template to cut out most of the redundant parts of the reaction allows you to make such reactions in as little as one line, allowing you to make dozens of reactions in very little time.

Anyway enough talk, time for the template:
[code]
{!TEMPLATE;FORGE_ITEM;id;name;type;mat;matname;count=1;
{REACTION;FORGE_%{id}_%{mat};ADDON_HOOK_SWAMP}
	[NAME:forge %matname %name]
	[BUILDING:WARCRAFTER_SAURIAN:NONE]
	[REAGENT:metal:150:BAR:NONE:INORGANIC:%mat]
	[PRODUCT:100:%count:%type:%id:INORGANIC:%mat]
	[FUEL]
	[SKILL:METALCRAFT]
}
[/code]
The first thing to note is that the template body consists mostly of normal Rubble code with a few odd little bits prefixed with a '%' symbol. Each of these "Replacement Tokens" is defined just before the template body and just after the template name. For example the first of these tokens (also called "Template Parameters", or just "params") is named "id" and can be accessed via "%id" or "%{id}". The value of a template parameter is determined when the template is called.

To really describe how a user defined template like this works we need values for the parameters, so here is an example call:
[code]
{FORGE_ITEM;ITEM_WEAPON_SPEAR_SAURIAN;small spear;WEAPON;STEEL_WATER;water steel}
[/code]

That call will expand into the following:
[code]
{REACTION;FORGE_ITEM_WEAPON_SPEAR_SAURIAN_STEEL_WATER;ADDON_HOOK_SWAMP}
	[NAME:forge water steel small spear]
	[BUILDING:WARCRAFTER_SAURIAN:NONE]
	[REAGENT:metal:150:BAR:NONE:INORGANIC:STEEL_WATER]
	[PRODUCT:100:1:WEAPON:ITEM_WEAPON_SPEAR_SAURIAN:INORGANIC:STEEL_WATER]
	[FUEL]
	[SKILL:METALCRAFT]
[/code]
As you can see it is just a straight forward variable expansion. 
One thing sharp-eyed readers may have noticed is that there was never a value defined for "count", if you will look up to the template definition you will see that count is followed by an equals sign (=). If a parameter name is followed by an equals sign whatever is after the equals sign is used as the default value for that parameter, so in this case count defaults to "1".

There are a few more interesting things that standard templates like this can do, but most of them have been replaced by the much more flexible rubble:template script command (which is too advanced to detail here). For more on how templates work, and how they interact with other parts of Rubble, see "Rubble Basics".

(The templates/items/materials used in this tutorial come from the (no longer available) "Better Dorfs/Saurian/Warcrafter" addon)

========================================================

Tutorial 3: Some Basic Script Commands and Concepts

Rubble includes powerful scripting support, almost everything can be controlled via scripting at some level.
The only "issue" is that rather than a popular scripting language like Lua or Python, Rubble uses a custom language named Rex.

Rex is a general purpose embedded scripting language that I use in most of my major projects. Rubble uses Rex extensively behind the scenes, almost every part of Rubble integrates with Rex at some level, for example EVERY template is either created by a Rex command or implemented directly by one.

There are three main reasons I use Rex:
	It is really easy to embed with the language I wrote Rubble in.
	I know it very well (after all I wrote it :p).
	Most of the integration work was already done as part of other projects.

The only real drawback to using Rex is that very few people know how to use it, this tutorial should help fix that.

It is a very good idea to read "Rex Docs/Rex Basics.txt" along side this file, as that file has all the basic documentation for the Rex language, this tutorial is only to help fill in holes and expand upon the things in that document.

The first thing you have to understand is that Rex is what is know as a "command based" language, this means that most of the language revolves around "commands". Commands are basically what most languages would call functions, but they do things that most languages would use keywords for.
A good example is if. In Rex if is a command like any other, it even returns a value (which can be very useful).
At first using commands for almost everything seems a little odd, but it allows you to do some very weird and wonderful things :)

Return Values:
	In Rex all return values are stored in a single "register" (it isn't actually a hardware register) that is shared with things like variable dereferences and object literals, basically every time a value is retrieved it is stuck in the the return register. This means that every command "returns a value" whether it wants to or not. If it doesn't return something explicitly it is said to return "unchanged", eg. whatever was in the register before. If you want to insert the last retrieved value into a command call or the like you can use the nop command. nop (no operation) is a simple command that does nothing, and so it obviously returns unchanged, perfect for getting the current value of the return register. 

Some Important Commands:

	if condition then [else]
		(if true {(ret a)}{(ret b)})
		(if true {(ret a)})
		
		if is the main (only!) conditional command in Rex. It is very simple and used all over, just like the if statements in most languages.
		
		One important thing to keep in mind is that, like all commands, if returns a value! This means that cool things like the following can be done:
		(str:add "[NAME:Make Something" (if (int:eq [count] 0){""}{(str:add " (" [count] ")")}) "]")
		
		Which will return the following:
		(if count == 0) "[NAME:Make Something]"
		(if count == 5) "[NAME:Make Something (5)]"
		
		Generally when you need to invert the condition the convention is not to use bool:not and instead just include an empty "then" block:
		(if true {}{
			do something...
		})
		This is shorter, but maybe not as clear, look for empty blocks when reading the default addon scripts!
		
	loop code
		(loop {
			false
		})
		
		loop is the basic looping command, it will run it's code block until it returns false. This means that loop and if are close companions, the vast majority of loop commands will consists of a single if command that surrounds the actual loop body.
		
		loop is actually really rare, foreach is FAR more common.
	
	foreach indexable code
		(foreach <array a b c> block key value {
			(console:print [key] " = " [value] "\n")
			(break true)
		})
		
		foreach will run it's code block for every item in an indexable. The code block needs to have at least two variables, the first of which foreach will fill in with the current indexable key and the second will be filled with the value.
		Each indexable type will determine it's own iteration order, for example most maps will iterate randomly and arrays will go from low keys to high keys.
		
		This example is much more complicated than the previous ones. 
		First of all right after the command name you will see a "object literal", I will not try to explain them here, see the "Rex Docs/Rex Basics.txt" file that comes with Rubble.
		After that comes a block declaration, this is basically a way to add variables to a block before it is created, in this case the variables "key" and "value", which are then filled in by the foreach command. These variables are filled in by position, not name, so you could name them anything you want.
		The body uses the console:print command to write the values of key and value to the console, nothing too fancy. The important bit is the `(break true)` command at the bottom, this ensures that the block returns true, as if the block returns false the foreach loop will abort early.

	console:print items...
		(console:print "Hello!\n")
		
		Write stuff to the log, this command is really nice for printing status in long tweak scripts and/or helping with debugging.
	
	str:add strings...
		(str:add "s1" "s2") # s1s2
		
		Add two or more strings together, this command is probably the most common one after if and foreach, after all Rubble is all about string processing.
	
	int:add a b, int:sub a b
		(int:add 2 2) # 4
		(int:sub 5 4) # 1
		
		These two commands (and their brothers int:mul, int:div, and int:mod) make up the basic math commands. int:add adds two numbers and int:sub subtracts them, simple.
		There also exist counterparts in the float module for working with floating point numbers (there is obviously no float:mod :p).
	
	bool:and a b, bool:or a b, bool:not a
		(bool:and true false) # false
		(bool:or false true) # true
		(bool:not false) # true
		
		The basic boolean operations, I don't use these much, but when you need them you need them.
		In any case not much to say, if you don't know what these are then you need to go read a beginners programming book (any one will do).
		
		These are fairly uncommon, as since these commands do not support short-circuit evaluation most times you will need to use nested if commands or the short-circuit versions (and bool:not is avoided by the empty-then convention).
	
	bool:sand a b, bool:sor a b
		(bool:sand {true} {false}) # false
		(bool:sor {false} {true}) # true
		
		These are sorta like bool:and and bool:or, but they take code values instead of bool values and do short-circuit evaluation.
		The code values are required so that evaluation may be delayed until during the command call.
		
		Short-circuit evaluation means that after the first parameter is evaluated if it is obvious what the result will be the second parameter will be ignored, this means that the following may be done:
			(if (bool:sand {(exists [a] [x])} {(exists [a [x]] [y])}) {
				# do something to [a [x] [y]]
			})
	
	expr expression values...
		(expr "a+b" 1 3) # 4
		
		expr (short for "expression") provides a basic mathematical and logical expression parser.
		The supported operators are: (in precedence order)
			( )
			/ * %
			+ -
			== != > < >= <=
			!
			&& ||
		The types int, float, and bool are supported, any other type is converted to int or bool depending on context.
		Any non-operator character in the expression string is treated as a placeholder for one of the values, syntax errors will cause a crash (which will cause Rubble to abort, so don't do that!)
		If you only need to do something simple like add two numbers it is a far better idea to use the normal commands in the int, float, and bool modules, as expr has a relatively large amount of overhead compared to doing it by hand.
	
Rubble comes with basic documentation for each command, and it is a very good idea to read some (or all!) of the scripts in the default addons.

========================================================

Tutorial 4: Tileset Addons

This tutorial will show you how to build tileset support into your mod, as well as explain some simple procedures for porting existing tilesets to Rubble.

So you want your cool new mod to support tilesets other than the default? Good thing Rubble goes out of its way to make that easy.

Rubble has two special kinds of files to help deal with tilesets: ".tset" files and ".tset.rex" files. ".tset" files contain definitions that are used for any tile and color information the tileset needs. These files use normal raw syntax, but simplified and with only the tags important to tilesets. ".tset.rex" files are scripts that are more-or-less functionally identical to post scripts, but they are also run if a user adds the tileset to existing raws.

.tset files are normally automatically generated, the esiest way to do this is to create a temporary addon with just the tilesets raws, then generate with ONLY that addon active. In the output directory there will be a file named "current.tset", this should contain what you need.

Of course you can also make a .tset file by hand, but unless you are making a tileset from scratch generation is faster.

Once you have your .tset file ready to go (however you got it) create a new addon and add this file as well as your tileset's font image.

Now for your tileset install script (the .tset.rex file). This script should install your tileset image and make any changes to the init file that may be needed.

To install you tileset image all you need to do is add a call to "rubble:install_tilesheet", for example:
[code]
(rubble:install_tilesheet "test_16x16.png")
[/code]

Now for editing the init files. Rubble provides a set of commands specifically for editing the tileset related init settings found in "init.txt" and "d_init.txt". 
To set the "font" (eg the tilesheet) you use the "rubble:set_fullscreen_tilesheet" and "rubble:set_windowed_tilesheet" commands. Both commands also set the font to be used in graphical mode, if you need to set these separately (if you are using Text Will Be Text for example), then use "rubble:set_fullscreen_font", "rubble:set_fullscreen_font_graphics", "rubble:set_windowed_font", and "rubble:set_windowed_font_graphics".
To set things like track tiles you use "rubble:open_d_init", "rubble:edit_d_init", "rubble:close_d_init", and "rubble:d_init_to_defaults". The following example makes pillars look like smooth floors ;)
[code]
(rubble:open_d_init)
(rubble:edit_d_init PILLAR_TILE "'+'")
(rubble:close_d_init)
[/code]
Note that "rubble:close_d_init" will set ALL tileset related init settings to their defaults unless a call to "rubble:edit_d_init" has explicitly set that value! This is so that it is easy to reset values changed by other tilesets you may have installed before. 
If you use all default init settings all you need to do is call "rubble:d_init_to_defaults" to make sure a different tileset didn't clobber your settings.

========================================================

Tutorial 5: Basic Tweak Scripts

This tutorial will introduce tweak scripts and attempt to show some uses for them.

I personally use tweak scripts all over the place in my addons, as (by design) tweak scripts are incredibly flexible.

The original use I had in mind for tweak scripts was allowing you to "tweak" the raws after generation in ways that are hard for templates to do.
For example, the addon "User/Zap Aquifers", this addon consists of a single "post script" (a tweak script that runs after generation) that disables all AQUIFER tags in inorganic materials.
Here is the whole script that I will be dissecting:
[code]

var aquifercount = 0
var aquifercounttotal = 0

(foreach [rubble:raws] block name content {
	(if (str:cmp "inorganic_" (str:left [name] 10)) {
		(console:print "    " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "AQUIFER") {
				[tag disable = true]
				[aquifercount = (int:add [aquifercount] 1)]
			})
			
			(break true)
		})]
		(console:print "      Found Aquifers : " [aquifercount] "\n")
		[aquifercounttotal = (int:add [aquifercounttotal] [aquifercount])]
		[aquifercount = 0]
	})
	(break true)
})
(console:print "    Found Aquifers Total : " [aquifercounttotal] "\n")
[/code]

First we will start with the top:
[code]
var aquifercount = 0
var aquifercounttotal = 0

(foreach [rubble:raws] block name content {
	...
})
[/code]
The two variables are for logging, they keep track of how many aquifer tags we found in the current file and in all the files together.
rubble:raws is an indexable that contains all the files in all the (active) addons indexable by name.
The block declaration is needed so foreach has some variables to populate, two of them to be exact (a block declaration is basically a anonymous command, they are used where you need variables without having formal variable declarations).

Now for the body of the foreach loop:
[code]
...

(foreach [rubble:raws] block name content {
	(if (str:cmp "inorganic_" (str:left [name] 10)) {
		(console:print "    " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			...
			
			(break true)
		})]
		(console:print "      Found Aquifers : " [aquifercount] "\n")
		[aquifercounttotal = (int:add [aquifercounttotal] [aquifercount])]
		[aquifercount = 0]
	})
	(break true)
})
...
[/code]
First we have an if command to check if the file contains inorganics, if not we just go on to the next file, simple.
If the file does contain inorganics we setup a few temporary variables, print the file name, and start parsing. Raw parsing is provided by commands in the "df:raw" module, and is rather simple, as all you can do is step through the file one tag at a time modifying the current tag if you wish (unless you use the slower advanced parser, which we are not doing here). The raw walker automatically returns a string version of the file as modified, so we assign that to the correct entry in the rubble:raws indexable.
After all tags have been processed it prints a message to say how many aquifers it has found, adds the current count to the total, and resets the current count.

Now for the meat of the script, the few lines that actually DO something.
[code]
...

(foreach [rubble:raws] block name content {
	(if (str:cmp "inorganic_" (str:left [name] 10)) {
		...
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "AQUIFER") {
				[tag disable = true]
				[aquifercount = (int:add [aquifercount] 1)]
			})
			
			(break true)
		})]
		...
	})
	(break true)
})
(console:print "    Found Aquifers Total : " [aquifercounttotal] "\n")
[/code]
First we check if the tag id equals "AQUIFER" and disable the tag if so (disabling a tag just replaces it's square brackets with dashes). Also we increment the current aquifer count.
Finally at the end, after all files are processed we print a message saying how many aquifers we found total in all files.

One of the other big uses for tweak scripts is declaring templates.
Script templates are far more flexible than normal templates, and the entire standard library is implemented this way.
Most script templates are declared in pre scripts (scripts that run before generation).
[code]
(rubble:template @TEST block a b c {
	(str:add "You passed: " [a] ", " [b] ", and " [c] " to the @TEST template!")
})
[/code]

"So how does Rubble know what is a postscript and what is a prescript?" you ask. Simple, the file extension. Rubble "knows" a file ending in .pre.rex is a prescript and a file ending in .post.rex is a post script.

There are many, many things that you can do with tweak scripts, for more examples see the addons that come with Rubble.

========================================================

Tutorial 6: The raw Parsers

This tutorial will teach you how to use the script raw parsers.

Rubble comes with two, slightly different, raw parsers that are made available to tweak scripts to use for whatever strikes your fancy.

The first parser is provided by the df:raw:walk command. This is the parser that is generally used, as it is much faster and a little easier to use at the cost of being slightly less flexible.
The primary issue with using this parser is the fact that you can only look at one tag at a time, and you can only move forward from the first tag to the last (it is possible to abort early if you wish), still this is just fine for the vast majority of tasks.

To use this raw parser you pass in a string containing the raws you wish to parse as well as a block of code to run for each tag, the command then returns a string with the raws plus any changes you may have made.

Example:
[code]
(console:print "Raws: " (df:raw:walk "[FOO][BAR][BAZ]" block tag {
	(console:print "The tag ID is: " [tag id] "\n")
	(break true)
}) "\n")
[/code]

Result:
[code]
The tag ID is: FOO
The tag ID is: BAR
The tag ID is: BAZ
Raws: [FOO][BAR][BAZ]
[/code]

Internally df:raw:walk is a lot like foreach, so like foreach you have to remember to "break true" unless you want to abort early (breakloop also works as expected).

The other parser is much more flexible, but it can also be much slower. Rather than running code for each tag this parser converts the raws to an indexable that you can use however you wish. This means that you can store single tags and iterate over the raws in any order that you wish, but converting the indexable back to a string is rather slow.
Another cool thing you can do with this parser is reorder the tags, or even construct new raw files, as the file is just a standard array indexable filled with custom tag indexables.

As this parser just generates a standard array (of special tag values) you just use the normal methods you would use with any indexable when working with it.

Example:
[code]
var raws = (df:raw:parse "[FOO][BAR][BAZ]")
(foreach [raws] block i tag {
	(console:print "The tag ID is: " [tag id] "\n")
	(break true)
})
(console:print "Raws: " (df:raw:dump [raws]) "\n")
[/code]

Result:
[code]
The tag ID is: FOO
The tag ID is: BAR
The tag ID is: BAZ
Raws: [FOO][BAR][BAZ]
[/code]

Although in most cases the tag values generated by the two parsers look and act much the same they are very different!

The "fast" parser (df:raw:walk) supports the following fields:
	Write-Only:
		replace			Replace the entire tag with the given text.
		append			Add text after the tag.
		prepend			Add text before the tag.
		disable			Replace the tag's square brackets with dashes.
	Read-Only:
		id				The tag "ID", eg the first (or only) item.
		<an integer>	A tag's "parameters", eg every item after the first, the len command returns how many of these there is.

The "slow" parser (df:raw:parse) supports the following fields:
	replace			Replace the entire tag with the given text, if this is set id and params is ignored.
	append			Add text after the tag.
	prepend			Add text before the tag.
	disable			Replace the tag's square brackets with dashes.
	comments		Any non-raw text that existed between this tag and the next (or the end of the file). Includes whitespace!
	id				The tag's "ID", eg the first (or only) item.
	params			An indexable (by default a standard array) containing the tag parameters, eg every item after the first.

All of the slow parser fields can be both read and written, so changing the value of the second parameter of a tag with the slow parser is just a matter of "[tag params 1 = value]" where as with the fast parser you have to replace the whole tag (with preformatted text!).

The slow parser does support the "replace" index, but it is a better idea to use the df:raw:tag command if you want to replace a tag with another (better yet modify the existing tag). The replace index is for when you want to replace a tag with arbitrary text, for example if you are inserting a Rubble template.

Example (df:raw:tag):
[code]
var raws = (df:raw:parse "[FOO][BAR][BAZ]")
[raws 1 = (df:raw:tag "NEW" "TAG")]
(console:print "Raws: " (df:raw:dump [raws]) "\n")
[/code]

Example (modify old tag):
[code]
var raws = (df:raw:parse "[FOO][BAR][BAZ]")
[raws 1 id = "NEW"]
[raws 1 params = <array "TAG">]
(console:print "Raws: " (df:raw:dump [raws]) "\n")
[/code]

Result (both):
[code]
Raws: [FOO][NEW:TAG][BAZ]
[/code]

While I have been calling them the "fast" and "slow" parsers, this is not necessarily true, most of what makes the "slow" parser slow lies in df:raw:dump, so if you are not going to modify the raws and then write them out it may actually be faster, testing is needed.
