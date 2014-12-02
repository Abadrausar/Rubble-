
(rubble:requireaddon "Libs/Castes/DFHack/Transform" "Libs/Castes")

module rubble:libs_castes_transform

# creature -> array of transform objects
var rubble:libs_castes_transform:transforms = (sort:new)

(rubble:template "CASTE_TRANSFORM" block creature to intermediate time {
	(if (exists [rubble:libs_castes_transform:transforms] [creature]) {
	}{
		[rubble:libs_castes_transform:transforms [creature] = (sort:new)]
	})
	
	[rubble:libs_castes_transform:transforms [creature] (str:add [creature] "_" [to]) = <map
		"to"=[to]
		"intermediate"=[intermediate]
		"time"=[time]
	>]
	(ret "")
})

(rubble:template "#CASTE_TRANSFORM_PRODUCT" block creature to {
	(str:add "[PRODUCT:100:1:BOULDER:NONE:INORGANIC:CASTE_TRANSFORM_" [creature] "_" [to] "]")
})

(rubble:template "#CASTE_TRANSFORM_MATS" {
	var out = ""
	(foreach [rubble:libs_castes_transform:transforms] block creature transforms {
		(foreach [transforms] block _ transform {
			[out = (str:add [out]
				"\n"
				"[INORGANIC:CASTE_TRANSFORM_" [creature] "_" [transform "to"] "]\n"
				"\t[STATE_COLOR:ALL:ORANGE]\n"
				"\t[STATE_NAME_ADJ:ALL:caste transformation]\n"
				"\t[SPECIAL]\n"
				"\t[NO_STONE_STOCKPILE]\n"
				"\t[DISPLAY_COLOR:7:0:0]\n"
				"\t[MATERIAL_VALUE:0]\n"
				"\t[SPEC_HEAT:800]\n"
				"\t[IGNITE_POINT:NONE]\n"
				"\t[MELTING_POINT:8998]\n"
				"\t[BOILING_POINT:8999]\n"
				"\t[HEATDAM_POINT:NONE]\n"
				"\t[COLDDAM_POINT:NONE]\n"
				"\t[MAT_FIXED_TEMP:8999]\n"
				"\t[SOLID_DENSITY:2000000]\n"
				"\t[LIQUID_DENSITY:2000000]\n"
				"\t[SYNDROME]\n"
				"\t\t{#SYN_AFFECTED_MALES;" [creature] "}\n"
				"\t\t[SYN_CLASS:\\WORKER_ONLY]\n"
				"\t\t[CE_CAN_DO_INTERACTION:START:0:END:299]\n"
				"\t\t\t[CDI:INTERACTION:CASTE_TRANSFORM_MALE_" [creature] "_" [transform "to"] "]\n"
				"\t\t\t[CDI:TARGET:A:SELF_ONLY]\n"
				"\t\t\t[CDI:ADV_NAME:change caste]\n"
				"\t\t\t[CDI:TARGET_RANGE:A:1]\n"
				"\t\t\t[CDI:MAX_TARGET_NUMBER:A:1]\n"
				"\t\t\t[CDI:WAIT_PERIOD:300]\n"
				"\t[SYNDROME]\n"
				"\t\t{#SYN_AFFECTED_FEMALES;" [creature] "}\n"
				"\t\t[SYN_CLASS:\\WORKER_ONLY]\n"
				"\t\t[CE_CAN_DO_INTERACTION:START:0:END:299]\n"
				"\t\t\t[CDI:INTERACTION:CASTE_TRANSFORM_FEMALE_" [creature] "_" [transform "to"] "]\n"
				"\t\t\t[CDI:TARGET:A:SELF_ONLY]\n"
				"\t\t\t[CDI:ADV_NAME:change caste]\n"
				"\t\t\t[CDI:TARGET_RANGE:A:1]\n"
				"\t\t\t[CDI:MAX_TARGET_NUMBER:A:1]\n"
				"\t\t\t[CDI:WAIT_PERIOD:300]\n"
			)]
			(break true)
		})
		(break true)
	})
	
	(str:trimspace (rubble:stageparse [out]))
})

(rubble:template "#CASTE_TRANSFORM_INTERACTIONS" {
	var out = ""
	(foreach [rubble:libs_castes_transform:transforms] block creature transforms {
		(foreach [transforms] block _ transform {
			[out = (str:add [out]
				"\n"
				"[INTERACTION:CASTE_TRANSFORM_MALE_" [creature] "_" [transform "to"] "]\n"
				"\t[I_TARGET:A:CREATURE]\n"
				"\t\t[IT_LOCATION:CONTEXT_CREATURE]\n"
				"\t\t[IT_MANUAL_INPUT:self]\n"
				"\t[I_EFFECT:ADD_SYNDROME]\n"
				"\t\t[IE_TARGET:A]\n"
				"\t\t[IE_IMMEDIATE]\n"
				"\t\t[SYNDROME]\n"
				"\t\t\t[CE_BODY_TRANSFORMATION:PROB:100:START:0:END:" [transform "time"] "]\n"
				"\t\t\t[CE:CREATURE:" [transform "intermediate"] "]\n"
				"\t\t\t[SYN_CLASS:\\PERMANENT]\n"
				"\t\t\t[SYN_CLASS:" [creature] "]\n"
				"\t\t\t[SYN_CLASS:MALE_" [transform "to"] "]\n"
				"\n"
				"[INTERACTION:CASTE_TRANSFORM_FEMALE_" [creature] "_" [transform "to"] "]\n"
				"\t[I_TARGET:A:CREATURE]\n"
				"\t\t[IT_LOCATION:CONTEXT_CREATURE]\n"
				"\t\t[IT_MANUAL_INPUT:self]\n"
				"\t[I_EFFECT:ADD_SYNDROME]\n"
				"\t\t[IE_TARGET:A]\n"
				"\t\t[IE_IMMEDIATE]\n"
				"\t\t[SYNDROME]\n"
				"\t\t\t[CE_BODY_TRANSFORMATION:PROB:100:START:0:END:" [transform "time"] "]\n"
				"\t\t\t[CE:CREATURE:" [transform "intermediate"] "]\n"
				"\t\t\t[SYN_CLASS:\\PERMANENT]\n"
				"\t\t\t[SYN_CLASS:" [creature] "]\n"
				"\t\t\t[SYN_CLASS:FEMALE_" [transform "to"] "]\n"
			)]
			(break true)
		})
		(break true)
	})
	
	(str:trimspace [out])
})
