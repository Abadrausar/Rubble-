
GoBlast

GoBlast is an easy to use (but not extensible) Blast-compatible raw generator

==============================================
Overview:
==============================================
GoBlast supports all of the standard Blast templates but it is not extendable via the normal Blast methods. GoBlast makes up for this by being easy to setup and distribute (and you don't need perl :p).
If you already have perl on your system you may just want to use Blast instead. 

GoBlast comes prebuilt for 32 bit Windows, Linux, and OSX

==============================================
Compatability Notes:
==============================================
To run GoBlast in Blast compatible mode you need to pass the command line option "-compat=true". In Blast compatibility mode GoBlast acts just like Blast with the exception of ignoring any "addon.pm" and "base.pm" files.

In simple mode (the default) GoBlast reads from "./source" and writes to "./objects". Addons are not supported in this mode and output file names are not postfixed with a namespace. The configuration is read from "./config.txt".

Not all templates produce output exactly like Blast. Any differences should be formatting only, anything else is a bug.

GoBlast is more permissive then Blast, for example GoBlast does not require ANY names to be "localized" with "@" and some templates (notably SHARED_INORGANIC and SHARED_ITEM) do no checks to see if the object being defined is well formed.

!GLOBAL_TEMPLATE is an alias for !TEMPLATE as GoBlast has no concept of a local template.

The replacement for "@@@" is an eight digit HEXADECIMAL number as opposed to Blast's seven digit decimal number, this should cause no problems (unless you use more than 10 million of those and are counting on the wrap-around :p).

Blast has a few inconsistent spots in its documentation, anyplace where usage counterdicts documents GoBlast sticks to usage.

What follows is a full list of docs/usage conflicts I could find:
	Docs for TECH_CLASS say that it only takes one class but there is at least one place in the standard base where it takes two
	
	Docs for #USES_TECH_CLASSES list possible tags as a parameter but no calls to this template have tags specified, also TECH_CLASS does not seem to have provisions for specifing tags

In anycase GoBlast should produce valid raws when run (in compatibility mode) on the standard base.

==============================================
Builtin Blast Templates:
==============================================
	!TEMPLATE, !T
	!GLOBAL_TEMPLATE
	ONCE
	STATIC, S
	LOG, L
	COMMENT, C
	VOID, V
	GETVAR
	SETVAR
	IFVAR
	#ADVENTURE_TIER
	ITEM_CLASS
	#USES_ITEM_CLASSES
	TECH_CLASS
	#USES_TECH_CLASSES
	#AQUIFER
	#ADV_TIME
	#FORT_TIME
	SHARED_INORGANIC
	REGISTER_ORE
	REGISTER_REACTION_CLASS
	SHARED_ITEM

All of the above templates work just like they do in Blast

==============================================
Builtin Non-Blast Templates:
==============================================
	{REGISTER_REACTION_PRODUCT;<INORGANIC>;<CLASS>;<MAT>}
        Adds a [MATERIAL_REACTION_PRODUCT:<CLASS>:<MAT>] tag to the definition of <INORGANIC>. For this to work, <INORGANIC> must be registered as a SHARED_INORGANIC. If this template call is the first to define this reaction product, then the given <MAT> is the one used in the finished raws.
	
	{CONFIG;<NAME>;<VALUE>}
		Creates or sets a Blast configuration variable. Configuration variables have names like @CONFIG_<NAME>. Equivelent to <NAME>=<VALUE> in config.txt.

==============================================
Changelog:
==============================================
v1.1
	Added all the templates from the default Blast base
	Added Blast compatibility mode
	Added REGISTER_REACTION_PRODUCT
	Added support for config.txt
	Added support for @, @@, and @@@

v1.0
	First version