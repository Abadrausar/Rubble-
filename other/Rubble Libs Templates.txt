	
Rubble Template Documentation
Templates added by the addons in the "Libs" addon group.

**********************************************
Libs/Castes

==============================================
{!REGISTER_CREATURE;<ID>;<DESCRIPTION>;<MALE>;<FEMALE>;<ADJ>}
	Register a creature with the caste library, this needs to be done before any other templates are called for this creature.

==============================================
{!DEFAULT_CASTE;<CREATURE>;<ID>;<DESC_NAME>;<NAME>;<NAME_PLUR>;<POPM>;<POPF>;<DESC>;<BONUS>}
	Create a new caste for the specified creature using the provided information.
	As it's name suggests this template is generally used for the default caste.

==============================================
{CASTE;<CREATURE>;<ID>;<DESC_NAME>;<NAME>;<NAME_PLUR>;<POPM>;<POPF>;<DESC>;<BONUS>}
	Create a new caste for the specified creature using the provided information.
	Exactly like !DEFAULT_CASTE but for all the other castes you may want to add.

==============================================
{#CASTE_INSERT;<CREATURE>}
	Write the generated castes here.
	Place this in your creature file where you want the castes to be added.

==============================================
{#SYN_AFFECTED_MALES;<CREATURE>;[<TABS>=2]}
{#SYN_AFFECTED_FEMALES;<CREATURE>;[<TABS>=2]}
	Write [SYN_AFFECTED_CREATURE:<CREATURE>:CASTE]tags for every male/female caste in the specified creature.
	VERY useful for transformations.
	<TABS> specifies the number of leading tabs to add to each new line, for the formatting freaks out there :)

**********************************************
Libs/Castes/DFHack/Transform

==============================================
{DFHACK_REACTION_CASTE_TRANSFORM;<ID>;<CLASS>;<CREATURE>;<CASTE>;[<DELAY>=0]}
	Transforms a creature to the specified type when the reaction is run.
	
	<CASTE> should be the ID passed to the CASTE or !DEFAULT_CASTE template used to create the caste you want to transform into. DO NOT use a gender prefix! The transformation script adds this automatically as needed.
	
	<DELAY> is a time in ticks to wait before transforming, this allows you to sequence several transformations.

==============================================
{DFHACK_CASTE_TRANSFORM;<CREATURE>;<CASTE>;[<DELAY>=0;[<ID>=nil]]}
	Transforms a creature to the specified type when the reaction is run.
	
	<CASTE> should be the ID passed to the CASTE or !DEFAULT_CASTE template used to create the caste you want to transform into. DO NOT use a gender prefix! The transformation script adds this automatically as needed.
	
	<DELAY> is a time in ticks to wait before transforming, this allows you to sequence several transformations.
	
	If <ID> is nil then the ID of the last reaction defined by a REACTION template is used.
	(Specialized variants of REACTION, such as DFHACK_REACTION, work as well)

**********************************************
Libs/Crates

==============================================
{!CRATE;<ID>;<NAME>;<VALUE>;<PRODUCT>}
	Define a new crate.
	<PRODUCT> is written directly into the unpack reaction.
	
	This template is always active even when it's addon is not enabled.
	If the addon is not active then this template does nothing.

==============================================
{!CRATE_BARS;<ID>;<NAME>;<VALUE>;<PRODUCT>}
	Define a new crate containing 10 bar items.
	<PRODUCT> is a material token.
	
	This template is always active even when it's addon is not enabled.
	If the addon is not active then this template does nothing.

==============================================
{!CRATE_CLASS;[<ID>];<CLASS>}
	Add a crate to a class.
	If <ID> is not specified it defaults to the id of the last crate defined.
	
	This template is always active even when it's addon is not enabled.
	If the addon is not active then this template does nothing.

==============================================
{CRATE_WORLDGEN_REACTION_PRODUCTS}
	This generates a list of product lines for ALL crates, use in world gen reactions.
	
	This template is always active even when it's addon is not enabled.
	If the addon is not active then this template does nothing.

==============================================
{CRATE_WORLDGEN_REACTION_PRODUCTS_CLASSED;<CLASS>}
	This generates a list of product lines for all crates in a class, use in world gen reactions.
	
	This template is always active even when it's addon is not enabled.
	If the addon is not active then this template does nothing.

==============================================
{CRATE_UNPACK_REACTIONS;<BUILDING>;<SKILL>;<TECHCLASS>;[<AUTO>=true]}
	Generate unpack reactions for all crates.

==============================================
{CRATE_UNPACK_REACTIONS_CLASSED;<BUILDING>;<SKILL>;<TECHCLASS>;<CRATECLASS>;[<AUTO>=true]}
	Same as CRATE_UNPACK_REACTIONS, except for only a single crate class.

**********************************************
Libs/DFHack/Add Reaction

==============================================
{REACTION_ADD;<REACTION>;<BUILDING>}
	Add a reaction to a hardcoded building.
	Building List:
		CARPENTERS
		FARMERS
		MASONS
		CRAFTSDWARFS
		JEWELERS
		METALSMITHSFORGE
		MAGMAFORGE
		BOWYERS
		MECHANICS
		SIEGE
		BUTCHERS
		LEATHERWORKS
		TANNERS
		CLOTHIERS
		FISHERY
		STILL
		LOOM
		QUERN
		KENNELS
		ASHERY
		KITCHEN
		DYERS
		TOOL
		MILLSTONE
		
		WOOD_FURNACE
		SMELTER
		GLASS_FURNACE
		MAGMA_SMELTER
		MAGMA_GLASS_FURNACE
		MAGMA_KILN
		KILN

==============================================
{REACTION_CLEAR_LIST;<BUILDING>}
	Clear the reaction list of a hardcoded building.
	Uses the same building list as REACTION_ADD.

**********************************************
Libs/DFHack/Announcement

==============================================
{DFHACK_REACTION_ANNOUNCE;<ID>;<CLASS>;<TEXT>;[<COLOR>=COLOR_WHITE]}
	Write an announcement to the screen and gamelog when the reaction completes, works exactly like the REACTION base template.
	
	This template is always active even when it's addon is not enabled.
	If the addon is not active then this template acts exactly like REACTION.

==============================================
{DFHACK_ANNOUNCE;<TEXT>;[<COLOR>=COLOR_WHITE];[<ID>=nil]}
	Write an announcement to the screen and gamelog when the reaction <ID> completes, works exactly like the DFHACK_REACTION_BIND base template.
	If <ID> is nil then the ID of the last reaction defined by a REACTION template is used.
	(Specialized variants of REACTION, such as DFHACK_REACTION, work as well)
	
	This template is always active even when it's addon is not enabled.
	If the addon is not active then this template does nothing.

**********************************************
Libs/DFHack/Upgrade Building

==============================================
{DFHACK_REACTION_UPGRADE_BUILDING;<SOURCE>;<DEST>;<NAME>;<CLASS>}
	Change a custom workshop/furnace from one type to another (via a reaction).
	
	<SOURCE> and <DEST> are the IDs of the workshop to upgrade and what it should be changed to.
	<NAME> is the user visible name for the upgrade reaction
	<CLASS> is the addon hook to register the reaction with.
	
	This template provides a full featured reaction with no products or reagents, if you need either just add them after the template.
	May be used with all the normal templates that are designed to work with REACTION.
	Basically this expands to a REACTION, NAME, and BUILDING tag and otherwise acts much like the DFHACK_REACTION template.
	
	You should only change workshop -> workshop or furnace -> furnace!
	Blocked and work tiles should match up and size should be the same.

**********************************************
Libs/Research

==============================================
{!RESEARCH_ADD;<ID>;<NAME>;<CLASS>;<REAGENTS>}
	Add a new research topic.
	
	This template is always active even when it's addon is not enabled.
	If the addon is not active then this template does nothing.

==============================================
{RESEARCH_DISCOVERY;<ID>;<CLASS>}
	Returns a build item for the specified research topic.
	
	This template is always active even when it's addon is not enabled.
	If the addon is not active then this template does nothing.

==============================================
{RESEARCH_DISCOVERY_REAGENT;<ID>;<CLASS>}
	Returns a reagent for the specified research topic.
	
	This template is always active even when it's addon is not enabled.
	If the addon is not active then this template does nothing.

==============================================
{!RESEARCH;<ID>;<NAME>;<CLASS>;<REAGENTS>}
	Add a new research topic and returns a build item for it.
	Basically !RESEARCH_ADD + RESEARCH_DISCOVERY.
	
	This template is always active even when it's addon is not enabled.
	If the addon is not active then this template does nothing.

==============================================
{RESEARCH_BASE_REACTIONS;<BUILDING>;<SKILL>;[<HOOK>="ADDON_HOOK_PLAYABLE"]}
	Returns the base reactions that make the research system work.

==============================================
{RESEARCH_REACTIONS;<CLASS>;<BUILDING>;<SKILL>;[<HOOK>="ADDON_HOOK_PLAYABLE"]}
	Returns the generated reaction for a specified class.
