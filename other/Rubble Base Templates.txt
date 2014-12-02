	
Rubble Template Documentation

Script Templates defined in the Rubble binary.

==============================================
{!TEMPLATE;<NAME>;[<PARAM>[=<DEFAULT_VALUE>]...];<CODE>}
	Creates a new template definition with name <NAME>, subsequently calling {<NAME>} will return a new parse of <CODE>.

	Calling {<NAME>;<ARGUMENT_1>;<ARGUMENT_2>} will replace instances of %<PARAM_1> or %{<PARAM_1>} with <ARGUMENT_1>, instances of %<PARAM_2> or %{<PARAM_2>} with <ARGUMENT_2>, etcetera. Default values may be specified for parameters.
	
	Variables ($<VAR_NAME> or ${<VAR_NAME>}) will be expanded in both the parameters and the <CODE>.
	(variables will NOT be expanded if they are in the body of a nested child template!)
	
	Example:
		{!TEMPLATE;FOO;bar}
		{FOO}
		{FOO}
		{!TEMPLATE;GREET;thing;Hello %{thing}!}
		{GREET;World}
		{!TEMPLATE;GREET_DWARF;dwarf=Urist;Hello %{dwarf}!}
		{GREET_DWARF}
		{GREET_DWARF;Led}
		{@SET;TEST;‼Fun‼}
		{GREET;$TEST}
	
	evaluates to
		bar
		bar
		Hello World!
		Hello Urist!
		Hello Led!
		Hello ‼Fun‼!

User Templates defined in the "Libs/Base" addon.

==============================================
{COMMENT;<STUFF>...}
{C;<STUFF>...}
	Doesn't parse or return anything.

==============================================
{!VOID;<PRERAWS>...}
{VOID;<PRERAWS>...}
{#VOID;<PRERAWS>...}
{V;<PRERAWS>...}
	Parses <PRERAWS>, but doesn't return anything. Useful for suppressing the normal output of a template.

==============================================
{!ECHO;<PRERAWS>...}
{ECHO;<PRERAWS>...}
{#ECHO;<PRERAWS>...}
{E;<PRERAWS>...}
	Returns <PRERAWS>. Used to strip leading and trailing whitespace for better formatting of output and to control variable expansion.

==============================================
{!ABORT;<MESSAGE>}
{ABORT;<MESSAGE>}
{#ABORT;<MESSAGE>}
	Forces Rubble to exit, <MESSAGE> is displayed.

==============================================
{!PRINT;<MSG>...}
{PRINT;<MSG>...}
{#PRINT;<MSG>...}
	Writes <MSG> to the console. Each param gets it's own line.

==============================================
{@SCRIPT;<CODE>}
	Runs script code and returns the result (as a string).

==============================================
{@SET;<NAME>;<VALUE>}
	Sets a variable of name <NAME> to value <VALUE>. Returns nothing.

==============================================
{@IF;<STRING1>;<STRING2>;<THEN_PRERAWS>;[<ELSE_PRERAWS>=""]}
	If <STRING1> and <STRING2> are equal, then <THEN_PRERAWS> are parsed and returned. Else, <ELSE_PRERAWS> are parsed and returned. This is very useful with variables.
	
	Example:
		{IF;$TEST_VAR;YES;[FOO];[BAR]}

==============================================
{@IF_ACTIVE;<ADDON>;<THEN_PRERAWS>;[<ELSE_PRERAWS>=""]}
	If <ADDON> is active, then <THEN_PRERAWS> are parsed and returned. Else, <ELSE_PRERAWS> are parsed and returned.

==============================================
{@IF_CODE;<CODE>;<THEN_PRERAWS>;[<ELSE_PRERAWS>=""]}
	If <CODE> returns true, then <THEN_PRERAWS> are parsed and returned. Else, <ELSE_PRERAWS> are parsed and returned.
	(<CODE> may be any valid script code)
	
==============================================
{@IF_SKIP;<STRING1>;<STRING2>}
	If <STRING1> and <STRING2> are equal then skip the current file.

==============================================
{SHARED_OBJECT;<ID>;<DEFINITION>}
	Adds a common object with the id <ID> to the dictionary. <DEFINITION> may be any raws. If this template call is the first with this <ID>, then the given <DEFINITION> will be used in the finished raws.
	
	Note that the contents of <DEFINITION> are always parsed, whether or not the results will appear in the raws.

==============================================
{SHARED_PLANT;<ID>;<DEFINITION>}
	Adds a common plant to the dictionary. <DEFINITION> should be a complete plant entry, EXCLUDING the [PLANT:<ID>] tag. 
	
	SHARED_PLANT is simply a specialized version of SHARED_OBJECT.

==============================================
{SHARED_INORGANIC;<ID>;<DEFINITION>}
	Adds a common inorganic to the dictionary. <DEFINITION> should be a complete inorganic stone entry, EXCLUDING the [INORGANIC:<ID>] tag. If this template call is the first to define this inorganic, then the given <DEFINITION> will be used in the finished raws. 
	
	SHARED_INORGANIC is simply a specialized version of SHARED_OBJECT.

==============================================
{SHARED_MATERIAL_TEMPLATE;<ID>;<DEFINITION>}
	Adds a common material template to the dictionary. <DEFINITION> should be a complete template entry, EXCLUDING the [MATERIAL_TEMPLATE:<ID>] tag. If this template call is the first to define this material template, then the given <DEFINITION> will be used in the finished raws. 
	
	SHARED_MATERIAL_TEMPLATE is simply a specialized version of SHARED_OBJECT.

==============================================
{SHARED_OBJECT_EXISTS;<ID>;<THEN_PRERAWS>;[<ELSE_PRERAWS>=""]}
	If a SHARED_OBJECT with the id <ID> exists then <THEN_PRERAWS> are parsed and returned. Else, <ELSE_PRERAWS> are parsed and returned.
	This is very useful for making addons with additional behavior that depends on things (items, materials, ect) from other addons.

==============================================
{SHARED_OBJECT_ADD;<ID>;<PRERAWS>}
	Appends the result of parsing <PRERAWS> to the end of SHARED_OBJECT <ID>. If that shared object does not exist it saves the result to append to any definition that may come later.

==============================================
{REGISTER_REACTION_CLASS;<ID>;<CLASS>}
	Adds a [REACTION_CLASS:<CLASS>] tag to a SHARED_OBJECT.
	Just a specialized version of SHARED_OBJECT_ADD.

==============================================
{REGISTER_REACTION_PRODUCT;<ID>;<CLASS>;<PRODUCT>}
	Adds a [MATERIAL_REACTION_PRODUCT:<CLASS>:<PRODUCT>] tag to a SHARED_OBJECT.
	Just a specialized version of SHARED_OBJECT_ADD.

==============================================
{SHARED_OBJECT_KILL_TAG;<ID>;<TAG>}
	Disable all occurrences of <TAG> in shared object <ID>.
	<TAG> is the ID of the tag you wish to disable, do not add square brackets or specify it's parameters.
	
	Example:
		{SHARED_OBJECT_KILL_TAG;SLADE;UNDIGGABLE}

==============================================
{SHARED_OBJECT_REPLACE_TAG;<ID>;<TAG>;<REPLACEMENT>}
	Replace all occurrences of <TAG> in shared object <ID> with <REPLACEMENT>.
	<TAG> is the ID of the tag you wish to replace, do not add square brackets or specify it's parameters.
	
	Example:
		{SHARED_OBJECT_REPLACE_TAG;SLADE;UNDIGGABLE;[AQUIFER]}

==============================================
{@INSTALL_GRAPHICS_FILE;<FILE>}
	Copies the specified file from your addon into the raw/graphics directory.
	Use to install images.

==============================================
{@INSTALL_IMAGES_AS_GRAPHICS;<ADDON>}
	Runs @INSTALL_GRAPHICS_FILE for each PNG image in <ADDON>.
	<ADDON> MUST exist and MUST be active!

==============================================
{@GRAPHICS_FILE}
	Changes the type of the current file to a graphics file.
	Use to install graphics raws.

==============================================
{@TEXT_FILE}
	Changes the type of the current file to an auxiliary text file.
	Use to install files to "out:objects/text" (without adding a standard header).

==============================================
{BUILDING_WORKSHOP;<ID>;<CLASS>}
	Similar to ITEM_CLASS, this registers a workshop to class <CLASS>.
	Used with ADDON_HOOKS, #USES_TECH, or #USES_BUILDINGS.
	
	Returns [BUILDING_WORKSHOP:<ID>]

==============================================
{BUILDING_FURNACE;<ID>;<CLASS>}
	Similar to ITEM_CLASS, this registers a furnace to class <CLASS>.
	Used with ADDON_HOOKS, #USES_TECH, or #USES_BUILDINGS.
	
	Returns [BUILDING_FURNACE:<ID>]

==============================================
{BUILDING_ADD_CLASS;<CLASS>;[<ID>=nil]}
	Adds a class to an existing building.
	Used with ADDON_HOOKS, #USES_TECH, or #USES_BUILDINGS.
	
	If <ID> is not present the name of the last building defined by a BUILDING_WORKSHOP or BUILDING_FURNACE template is used.

==============================================
{REMOVE_BUILDING;<ID>;<CLASS>}
	Removes a building from a class.
	The building does not need to exist yet, this template will work regardless of evaluation order.

==============================================
{REMOVE_BUILDING_FROM_PLAYABLES;<ID>}
	Removes a building from all addon hooks that describe playable races.
	
	The playability information comes from the !ENTITY_PLAYABLE template group.
	Make sure this is not evaluated until playability information is in it's final state!

==============================================
{#USES_BUILDINGS;<CLASSES>...}
	Usable in entity definitions. Expands to a list of building permissions of all <CLASS>es combined.
	
	It is a very good idea to use ADDON_HOOKS instead of this template!

==============================================
{REACTION;<ID>;<CLASS>}
	Similar to ITEM_CLASS, this registers a reaction to class <CLASS>.
	Used with ADDON_HOOKS, #USES_TECH, or #USES_REACTIONS.
	
	Returns [REACTION:<ID>]

==============================================
{DFHACK_REACTION;<ID>;<COMMAND>;<CLASS>}
	Exactly like REACTION, except that <COMMAND> is run whenever the reaction is completed (provided that DFHack is installed of course).

==============================================
{DFHACK_REACTION_BIND;<COMMAND>;[<ID>=nil]}
	Binds a command to a reaction, like DFHACK_REACTION except for an already existing reaction.
	
	If <ID> is not present the name of the last reaction defined by a REACTION template is used.
	(Specialized variants of REACTION, such as DFHACK_REACTION, work as well)

==============================================
{REACTION_ADD_CLASS;<CLASS>;[<ID>=nil]}
	Adds a class to an existing reaction.
	Used with ADDON_HOOKS, #USES_TECH, or #USES_REACTIONS.
	
	If <ID> is not present the name of the last reaction defined by a REACTION template is used.
	(Specialized variants of REACTION, such as DFHACK_REACTION, work as well)

==============================================
{REMOVE_REACTION;<ID>;<CLASS>}
	Removes a reaction from a class.
	The reaction does not need to exist yet, this template will work regardless of evaluation order.

==============================================
{REMOVE_REACTION_FROM_PLAYABLES;<ID>}
	Removes a reaction from all addon hooks that describe playable races.
	The playability information comes from the !ENTITY_PLAYABLE template group.
	Make sure this is not evaluated until playability information is in it's final state!

==============================================
{#USES_REACTIONS;<CLASSES>...}
	Usable in entity definitions. Expands to a list of reaction permissions of all <CLASS>es combined.
	
	It is a very good idea to use ADDON_HOOKS instead of this template!

==============================================
{#USES_TECH;<CLASSES>...}
	Combo of #USES_BUILDINGS and #USES_REACTIONS
	
	It is a very good idea to use ADDON_HOOKS instead of this template!

==============================================
{SHARED_ITEM;<TYPE>;<ITEM>;<DEFINITION>}
	Registers an item <ITEM> of type <TYPE>. Used later with ADDON_HOOKS, #USES_ITEMS, and ITEM_CLASS.
	
	Returns [ITEM_<TYPE>:<ITEM>]
	
	<TYPE> must be one of AMMO, ARMOR, DIGGER, FOOD, GLOVES, HELM, INSTRUMENT, PANTS, SHIELD, SHOES, SIEGEAMMO, TOOL, TOY, TRAPCOMP or WEAPON.
	
	Type DIGGER is treated like type WEAPON except in entities, there it translates into a DIGGER tag.
	Type FOOD just translates directly to a call to SHARED_OBJECT, no support for items classes is initiated (as foods do not need to be registered). The only reason SHARED_ITEM supports FOOD at all is for consistency.
	
	SHARED_ITEM is just a (very) specialized version of SHARED_OBJECT.

==============================================
{ITEM_CLASS;<CLASS>;[<RARITY>=COMMON]}
{ITEM_CLASS;<TYPE>;<ITEM>;<CLASS>;[<RARITY>=COMMON]}
	
	Sets an items class and rarity. 
	
	The first form of an ITEM_CLASS template always refers to the last ITEM template before it.
	The second form is for use in addons or other places where the call cannot follow the item definition.
	
	<RARITY> can be RARE, UNCOMMON, COMMON and FORCED.
	
	Example:
	
	{ITEM;WEAPON;ITEM_WEAPON_TEST;
		The weapon definition...
	}{ITEM_CLASS;TEST_WEAPONS}
	
	{#USES_ITEMS;TEST_WEAPONS} -> [WEAPON:ITEM_WEAPON_TEST]

==============================================
{REMOVE_ITEM;<ID>;<CLASS>}
	Removes an item from a class.
	The item does not need to exist yet, this template will work regardless of evaluation order.

==============================================
{REMOVE_ITEM_FROM_PLAYABLES;<ID>}
	Removes an item from all addon hooks that describe playable races.
	The playability information comes from the !ENTITY_PLAYABLE template group.
	Make sure this is not evaluated until playability information is in it's final state!

==============================================
{#USES_ITEMS;<CLASS>}
	Usable in entity definitions. Expands to a list of item permissions of the <CLASS>. When using multiple #USES_ITEMS calls make sure every item is returned by at most one #USES_ITEMS call.
	
	It is a very good idea to use ADDON_HOOKS instead of this template!

==============================================
{ADDON_HOOKS;<ID>}
	Usable in entity definitions, ADDON_HOOKS expands to a list of addon hooks for the entity (Addon hooks are tech/item classes).
	
	By default the following hooks are installed:
		ADDON_HOOK_GENERIC
		ADDON_HOOK_<ID>
	
	<ID> should be the same ID you pass to !ENTITY_PLAYABLE.
	
	If the entity is playable in fortress mode (according to !ENTITY_PLAYABLE and friends) then the ADDON_HOOK_PLAYABLE hook is also installed.
	
	This template interacts at a low level with !ENTITY_PLAYABLE, SHARED_ITEM, BUILDING_WORKSHOP, BUILDING_FURNACE, and REACTION as well as their many variants and supporting templates. Use this template and it's supports in all your entities for maximum compatibility with the standard addons.

==============================================
{!ENTITY_PLAYABLE;<ID>;<FORT>;<ADV>;<INDIV>}
	Sets an entity's playability state.
	
	<FORT> is whether or not the entity is playable in fortress mode.
	<ADV> controls adventure mode playability.
	<INDIV> is for playing outsiders in adventure mode.
	Acceptable values are true and false (0 and -1 will work, as will a few other things).
	
	<ID> should be the same ID you pass to ADDON_HOOKS.
	
	This template interacts at a low level with ADDON_HOOKS, SHARED_ITEM, BUILDING_WORKSHOP, BUILDING_FURNACE, and REACTION as well as their many variants and supporting templates. Use this template and it's supports in all your entities for maximum compatibility with the standard addons.

==============================================
{ENTITY_PLAYABLE_EDIT;<ID>;<KEY>;<VALUE>}
	Edits an entity's playability state (as set by !ENTITY_PLAYABLE).
	
	<KEY> may be one of FORT, ADV, or INDIV (not case sensitive).
	Acceptable values are true and false (0 and -1 will work, as will a few other things).

==============================================
{@IF_ENTITY_PLAYABLE;<ID>;<KEY>;<THEN>;[<ELSE>=""]}
	If an entity is playable for a specific mode parse <THEN> else parse <ELSE>.
	
	<KEY> may be one of FORT, ADV, or INDIV (not case sensitive).

==============================================
{#ENTITY_NOBLES;<ID>;<DEFAULT>}
	Write the nobles to the entity.
	
	If ENTITY_REPLACE_NOBLES was called use that value, else use <DEFAULT> + anything added by ENTITY_ADD_NOBLE.

==============================================
{ENTITY_ADD_NOBLE;<ID>;<NOBLE>}
	Add nobles to the end of the existing ones for the specified entity.
	This can be called more than once, as the nobles are simply added together.
	
	If ENTITY_REPLACE_NOBLES is used this template will do nothing.

==============================================
{ENTITY_REPLACE_NOBLES;<ID>;<NOBLES>}
	Replace ALL the nobles for the specified entity.
	
	If this is called for an entity any changes made by ENTITY_ADD_NOBLE are ignored.

==============================================
{DFHACK_LOADSCRIPT;<SCRIPTNAME>}
	Installs and adds for automatic loading the script <SCRIPTNAME>.
	This only works for scripts added by an addon, already installed scripts should use DFHACK_RUNCOMMAND.

==============================================
{DFHACK_RUNCOMMAND;<COMMAND>}
	Sets a DFHack command up to run when the world is loaded.
	
	Example:
		{DFHACK_RUNCOMMAND;workflow enable drybuckets auto-melt}

==============================================
{@ADV_TIME;<COUNT>;<UNIT>}
	Generates a time value for use in interactions and the like.
	This version is adventure mode centric.
	
	Using values like ".5" for the <COUNT> value should work but "1/2" will not.
	
	Valid values for <UNIT> are:
		SECOND
		SECONDS
		MINUTE
		MINUTES
		HOUR
		HOURS
		DAY
		DAYS
		WEEK
		WEEKS
		MONTH
		MONTHS
		SEASON
		SEASONS
		YEAR
		YEARS
	
	Example:
		{@ADV_TIME;5;DAYS} -> 432000

==============================================
{@FORT_TIME;<COUNT>;<UNIT>}
	Generates a time value for use in interactions and the like.
	Exactly like @ADV_TIME except for fortress mode.
	
	Units below MINUTE may be less than useful, as fortress mode time units lack precision.
	
	Example:
		{@FORT_TIME;5;DAYS} -> 6000

==============================================
{@GROWDUR;<COUNT>;<UNIT>}
	A direct replacement for the GROWDUR tag that lets you specify time in real-world units.
	Uses the same unit types as @ADV_TIME.
	
	Units below DAY may be less than useful, as growth duration time units lack precision.
	
	Example:
		{@GROWDUR;5;DAYS} -> [GROWDUR:60]

==============================================
{@BUILD_KEY;<KEY>;[<FURNACE>=false]}
	A replacement for the BUILD_KEY tag that automatically works out key conflicts.
	
	<KEY> should be of the form: X, SHIFT_X, CTRL_X or ALT_X.
	<FURNACE> should be true or false, set to true if the building is a furnace.
	
	If the requested key is already used it chooses the next open key in this order:
		X
		SHIFT_X
		CTRL_X
		ALT_X
	For example if you try to set a workshop key to "M" (which is used by the Mason's Workshop)
	you will get [BUILD_KEY:CUSTOM_SHIFT_M].
	"ALT_Z" will wrap around to "A", so no worries about running out (unless you use more than 104 buildings per category).
	
	Please note that while this template "knows" about the hard coded vanilla build keys it will not magically work with other build keys that do not use this template.

==============================================
{!DECOMPRESS;<ENCODED_TEXT>}
	Decompresses Rubble encoded text.
	
	This template is mostly for reducing file size when posting addons to the forum, but it is generally better to use encoded zip files for that (See the section on "Rubble Encoder" in "Rubble Basics").

==============================================
{!PATCH;<FILE>;<PATCH>}
	Applies a patch to the specified file.
	
	Warning! Make sure the file being patched is not processed before patching (unless you know what you are doing), as this can cause parts of the file to not be processed correctly and/or the patch will not apply completely.
	
	If a patch cannot cleanly apply then Rubble will abort!
