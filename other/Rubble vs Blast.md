
Early versions of Rubble were heavily based on a similar utility, Blast, but recent versions resemble that utility only in the basic syntax.

==============================================
Compared to Blast
==============================================

Pros:
	Comes with an extensive standard library (FAR more so than Blast)
	Extensive support for DFHack
	No need for an external runtime, Rubble is a native application
	Addons may override each other's files (without a fragile diff system)
	Full parser/lexer, not just a bunch of regular expressions
	It's really easy to mix script code and raw text in a template, use whichever one is easier
	Variable expansion, no more GET_VAR, works in any parse stage
	Many templates for registering objects and the like are replacements for vanilla raw tags allowing better formatting
	Some files may be parsed but not included in the generated raws
	Support for easily and quickly installing tilesets from addons (during AND after generation!)
	Allows you to run "tweak scripts" before or after generation for fine-tuning the results
	Easy to use HTML-based UI that supports all common actions
	Support for detailed addon meta data, including dependencies and descriptions.
	Addons can be grouped into directories to make relations clear (and to make addon packs easy to distribute)
	Support for loading addons (and groups of addons) directly from zip files
	Has an extra template prefix for greater flexibility
	It is possible to automatically load addons from the internet
	Faster (not that it really matters)

Cons:
	No support for handling file name collisions in any way other than as an override
	No support for Blast namespaces (eg. @, @@, and @@@)
	Variables are simple key:value pairs instead of Blast's more powerful system
	No built-in pretty printer (unless, like me, you think of this as a pro :p)
	The scripting language isn't exactly mainstream

==============================================
Why Another Raw Generator?
==============================================

Blast and Rubble have completely different visions, Blast is focused on single mods with a small collection of optional features available as addons where as Rubble is focused on many individual mods that are not part of a unified whole. To put it simply Blast was designed with making life easy for a single modder, where as Rubble tries to be more flexible, it can work (like Blast) as a modder time-saver or it can be used (with a little more effort) as a way to tie many dissimilar mods together in a way that "just works" from the users point of view. Of course Rubble simply provides the tools, it is up to the modder to use them, but Rubble goes out of it's way to make that easy.

Basically I just wanted something that allowed me to do most of the things I could do with Blast, but with less fuss and more power (with a great deal of that power being focused on mod interaction).

I have put a huge amount of effort into making Rubble as flexible and powerful as possible, so newer features like tweak scripts, the addon loader, the web UI, and addon.meta make Blast look like a toy.
A cool, powerful toy, but a toy nonetheless.

Blast is still a good utility, but Rubble is far more powerful.

==============================================
Porting From Blast
==============================================

Rubble does not support the split base/addons architecture Blast uses, everything is an addon. While there is a "Base" addon it is just like any other addon as far as Rubble is concerned. If you need your own "base" make sure it is marked as being incompatible with the default "Base" addon in it's addon.meta file! (more on addon.meta later)

You should take care that you do not reuse file names, as while Blast namespaces files automatically Rubble does not. The convention is to name your files like so:
	"optional_prefix_<addon_name_with_spaces_and_slashes_to_underscores>.<extension>"
For example:
	building_user_dfhack_powered.txt
	user_tilesets_mlc_normal.rbl
	aaa_dfhack_libs_base.pre.rex

Rubble uses a significantly different set of standard templates. There are VERY few templates that are the same for both utilities!

The "!TEMPLATE" template is more-or-less the same but even there you need to be careful. While Blast uses "@", "@@", and "@@@" as "namespace" placeholders, Rubble uses "@" as a template prefix that means "run whenever", so you will need to strip these characters from your template names (and optionally manually namespace them).

Blast also has a relatively complex variable system, Rubble does not support this, but instead encourages you to use scripts (and the much more flexible system they use).

Instead of ITEM_CLASS or TECH_CLASS Rubble provides simpler SHARED_ITEM, REACTION, BUILDING_WORKSHOP, and BUILDING_FURNACE templates. These templates are not simpler in a "reduced functionality" way, they are simpler in a "less typing" way :)
When porting it may be easier to use the ITEM_CLASS, REACTION_ADD_CLASS, and BUILDING_ADD_CLASS templates.

Speaking of item/tech classes, Rubble supports a tightly interlocked system generally called "addon hooks" that mostly takes the place of the more generic classes Blast uses. Addon hooks allow you to add items/tech to an entity without needing to know anything about it ahead of time (via ADDON_HOOK_PLAYABLE and ADDON_HOOK_GENERIC), plus part of the system is a set of templates for managing entity playability. See the documentation for the ADDON_HOOKS and !ENTITY_PLAYABLE templates.

Rubble also supports something called "addon.meta" rather than Blast's "config.txt". addon.meta is actually used for a lot of things and config variables is only a small part of what it is for. Generally you should not use config variables except for minor tweaks that most users will want to ignore, instead split each option into a separate addon.

One main use of addon.meta is specifying dependencies and incompatibilities (due to the way Rubble is constructed these are MUCH more important than in Blast). Rubble will automatically enable addons if you request it here (making library addon simple and easy).

Rubble is a completely different utility, porting is possible but the differences outnumber the commonalities these days.
