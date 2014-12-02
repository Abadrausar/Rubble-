
module rubble:research

var rubble:research:topics = <map>
var rubble:research:names = <map>
var rubble:research:handled = <map>

# Add a research topic.
(rubble:template "!RESEARCH_ADD" block id name class reagents {
	(if (isnil [rubble:research:topics [class]]) {
		[rubble:research:topics [class] = (sort:new)]
	})
	
	[rubble:research:topics  [class] [id] = [reagents]]
	[rubble:research:names [id] = [name]]
	(ret "")
})

# Use as a build item.
(rubble:template "#RESEARCH_DISCOVERY" block id class {
	(if (isnil [rubble:research:topics [class]]) {}{
		(ret "Research class does not exist (may not be an error).")
	})
	
	(if (isnil [rubble:research:handled [class]]) {}{
		(ret "Research class is not handled (may not be an error).")
	})
	
	(if (exists [rubble:research:topics [class]] [id]) {
		(ret (str:add "[BUILD_ITEM:1:TOOL:ITEM_TOOL_RESEARCH_DISCOVERY:INORGANIC:RESEARCH_PAPER_" [id] "_DISCOVERY]"))
	}{
		(ret "Research topic does not exist (may not be an error).")
	})
})

# Use as a reagent.
(rubble:template "RESEARCH_DISCOVERY_REAGENT" block id class {
	(if (isnil [rubble:research:topics [class]]) {}{
		(ret "Research class does not exist (may not be an error).")
	})
	
	(if (exists [rubble:research:topics] [id]) {
		(ret (str:add "[REAGENT:" [id] ":1:TOOL:ITEM_TOOL_RESEARCH_DISCOVERY:INORGANIC:RESEARCH_PAPER_" [id] "_DISCOVERY]"))
	}{
		(ret "Research topic does not exist (may not be an error).")
	})
})

# Convenience template, use as a build item.
(rubble:template "!RESEARCH" block id name class reagents {
	(if (isnil [rubble:research:topics [class]]) {
		[rubble:research:topics [class] = (sort:new)]
	})
	
	[rubble:research:topics [class] [id] = [reagents]]
	[rubble:research:names [id] = [name]]
	(ret (str:add "{#RESEARCH_DISCOVERY;" [id] ";" [class] "}"))
})

(rubble:template "RESEARCH_BASE_REACTIONS" block building skill hook="ADDON_HOOK_PLAYABLE" {
	(rubble:stageparse (str:add "{REACTION;RESEARCH_STEP_TWO;" [hook] "} 64% chance
	[NAME:continue research (II)]
	[BUILDING:" [building] ":NONE]
	[REAGENT:A:1:TOOL:ITEM_TOOL_RESEARCH_ONE:NONE:NONE]
	[PRODUCT:80:1:TOOL:ITEM_TOOL_RESEARCH_TWO:GET_MATERIAL_FROM_REAGENT:A:NONE]
	[SKILL:" [skill] "]
	[AUTOMATIC]

{REACTION;RESEARCH_STEP_THREE;" [hook] "} 51.2% chance
	[NAME:continue research (III)]
	[BUILDING:" [building] ":NONE]
	[REAGENT:A:1:TOOL:ITEM_TOOL_RESEARCH_TWO:NONE:NONE]
	[PRODUCT:80:1:TOOL:ITEM_TOOL_RESEARCH_THREE:GET_MATERIAL_FROM_REAGENT:A:NONE]
	[SKILL:" [skill] "]
	[AUTOMATIC]

{REACTION;RESEARCH_STEP_FOUR;" [hook] "} 40.96% chance
	[NAME:continue research (IV)]
	[BUILDING:" [building] ":NONE]
	[REAGENT:A:1:TOOL:ITEM_TOOL_RESEARCH_THREE:NONE:NONE]
	[PRODUCT:80:1:TOOL:ITEM_TOOL_RESEARCH_FOUR:GET_MATERIAL_FROM_REAGENT:A:NONE]
	[SKILL:" [skill] "]
	[AUTOMATIC]

{REACTION;RESEARCH_STEP_FIVE;" [hook] "} 32% chance
	[NAME:continue research (V)]
	[BUILDING:" [building] ":NONE]
	[REAGENT:A:1:TOOL:ITEM_TOOL_RESEARCH_FOUR:NONE:NONE]
	[PRODUCT:80:1:TOOL:ITEM_TOOL_RESEARCH_FIVE:GET_MATERIAL_FROM_REAGENT:A:NONE]
	[SKILL:" [skill] "]
	[AUTOMATIC]

{REACTION;RESEARCH_STEP_SIX;" [hook] "} 20.96% chance
	[NAME:continue research (VI)]
	[BUILDING:" [building] ":NONE]
	[REAGENT:A:1:TOOL:ITEM_TOOL_RESEARCH_FIVE:NONE:NONE]
	[PRODUCT:80:1:TOOL:ITEM_TOOL_RESEARCH_SIX:GET_MATERIAL_FROM_REAGENT:A:NONE]
	[SKILL:" [skill] "]
	[AUTOMATIC]

{DFHACK_REACTION_ANNOUNCE;RESEARCH_STEP_SEVEN;" [hook] ";Your dwarves have finished some research!} 100% chance. (20.96% in total)
	[NAME:finish research (VII)]
	[BUILDING:" [building] ":NONE]
	[REAGENT:A:1:TOOL:ITEM_TOOL_RESEARCH_SIX:NONE:NONE]
		[HAS_MATERIAL_REACTION_PRODUCT:DISCOVERY_MAT]
	[PRODUCT:100:1:TOOL:ITEM_TOOL_RESEARCH_DISCOVERY:GET_MATERIAL_FROM_REAGENT:A:DISCOVERY_MAT]
	[SKILL:" [skill] "]
	[AUTOMATIC]
"))
})

(rubble:template "RESEARCH_REACTIONS" block class building skill hook="ADDON_HOOK_PLAYABLE" {
	(if (isnil [rubble:research:topics [class]]) {
		(rubble:abort (str:add "Research class: " [class] " does not exist."))
	})
	[rubble:research:handled [class] = true]
	
	var out = ""
	(foreach [rubble:research:topics [class]] block id reagents {
		[out = (str:add [out] "\n\n"
			"{REACTION;RESEARCH_" [id] ";" [hook] "}\n"
			"\t[NAME:research " [rubble:research:names [id]] " (I)]\n"
			"\t[BUILDING:" [building] ":NONE]\n"
			"\t" [reagents] "\n"
			"\t[PRODUCT:80:1:TOOL:ITEM_TOOL_RESEARCH_ONE:INORGANIC:RESEARCH_PAPER_" [id] "]\n"
			"\t[SKILL:" [skill] "]"
		)]
	})
	(rubble:stageparse (str:trimspace [out]))
})

(rubble:template "RESEARCH_MATS" {
	var out = ""
	(foreach [rubble:research:topics] block class _ {
		(foreach [rubble:research:topics [class]] block id _ {
			[out = (str:add [out] "\n"
				"[INORGANIC:RESEARCH_PAPER_" [id] "]\n"
				"\t[USE_MATERIAL_TEMPLATE:RESEARCH_TEMPLATE]\n"
				"\t[SPECIAL] This tag is invalid in a template?\n"
				"\t[STATE_NAME_ADJ:ALL_SOLID:" [rubble:research:names [id]] " research]\n"
				"\t[MATERIAL_REACTION_PRODUCT:DISCOVERY_MAT:INORGANIC:RESEARCH_PAPER_" [id] "_DISCOVERY]\n"
				"\n"
				"[INORGANIC:RESEARCH_PAPER_" [id] "_DISCOVERY]\n"
				"\t[USE_MATERIAL_TEMPLATE:RESEARCH_TEMPLATE]\n"
				"\t[SPECIAL] This tag is invalid in a template?\n"
				"\t[STATE_NAME_ADJ:ALL_SOLID:" [rubble:research:names [id]] " discovery]\n"
			)]
		})
	})
	(str:trimspace [out])
})

# adding research that is complete before stage 7 would be easy, just use reactions classes to select different
# stage reactions.
