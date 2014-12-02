
# This just supports basic castes, for more advanced abilities you can further specialize these templates.

module rubble:libs_castes
var rubble:libs_castes:castes = <map>
var rubble:libs_castes:bonuses = <map>
var rubble:libs_castes:names = <map>
var rubble:libs_castes:names_plur = <map>

var rubble:libs_castes:creature_castes = (sort:new)
var rubble:libs_castes:creature_desc = <map>
var rubble:libs_castes:creature_male = <map>
var rubble:libs_castes:creature_female = <map>
var rubble:libs_castes:creature_adj = <map>

command rubble:libs_castes:registercreature id desc male female adj {
	[rubble:libs_castes:creature_castes [id] = (sort:new)]
	[rubble:libs_castes:creature_desc [id] = [desc]]
	[rubble:libs_castes:creature_male [id] = [male]]
	[rubble:libs_castes:creature_female [id] = [female]]
	[rubble:libs_castes:creature_adj [id] = [adj]]
	
	[rubble:libs_castes:castes [id] = <map>]
	[rubble:libs_castes:bonuses [id] = <map>]
	[rubble:libs_castes:names [id] = <map>]
	[rubble:libs_castes:names_plur [id] = <map>]
	
	(ret "")
}

command rubble:libs_castes:newcaste creature id desc_name name name_plur popm popf desc bonus {
	(if (str:cmp [desc_name] "") {}{
		[desc_name = (str:add "*" [desc_name] "* ")]
	})
	(if (str:cmp [desc] "") {}{
		[desc = (str:add " " [desc])]
	})
	
	[rubble:libs_castes:castes [creature] [id] = (str:add
		"\t[CASTE:MALE_" [id] "]\n"
		"\t\t[DESCRIPTION:" [desc_name] [rubble:libs_castes:creature_desc [creature]] [desc] "]\n"
		"\t\t[POP_RATIO:" [popm] "]\n"
		"\t[CASTE:FEMALE_" [id] "]\n"
		"\t\t[DESCRIPTION:" [desc_name] [rubble:libs_castes:creature_desc [creature]] [desc] "]\n"
		"\t\t[POP_RATIO:" [popf] "]\n"
	)]
	
	[rubble:libs_castes:creature_castes [creature] [id] = true]
	[rubble:libs_castes:bonuses [creature] [id] = [bonus]]
	[rubble:libs_castes:names [creature] [id] = [name]]
	[rubble:libs_castes:names_plur [creature] [id] = [name_plur]]
	(ret "")
}

(rubble:template "!REGISTER_CREATURE" block id desc male female adj {
	(rubble:libs_castes:registercreature [id] [desc] [male] [female] [adj])
	(ret "")
})

(rubble:template "!DEFAULT_CASTE" block creature id desc_name name name_plur popm popf desc bonus {
	(rubble:libs_castes:newcaste [creature] [id] [desc_name] [name] [name_plur] [popm] [popf] [desc] [bonus])
	(ret "")
})

(rubble:template "CASTE" block creature id desc_name name name_plur popm popf desc bonus {
	(rubble:libs_castes:newcaste [creature] [id] [desc_name] [name] [name_plur] [popm] [popf] [desc] [bonus])
	(ret "")
})

(rubble:template "#CASTE_INSERT" block creature {
	var out = ""
	
	# Caste Declarations
	[out = (str:add [out] "# Generated Castes\n")]
	(foreach [rubble:libs_castes:creature_castes [creature]] block caste _ {
		[out = (str:add [out] [rubble:libs_castes:castes [creature] [caste]])]
		(break true)
	})
	
	# Set male/female info
	var m_castes = ""
	var f_castes = ""
	var once = true
	(foreach [rubble:libs_castes:creature_castes [creature]] block id _ {
		(if [once] {
			[m_castes = (str:add [m_castes] "\t[SELECT_CASTE:MALE_" [id] "]\n")]
			[f_castes = (str:add [f_castes] "\t[SELECT_CASTE:FEMALE_" [id] "]\n")]
			[once = false]
		}{
			[m_castes = (str:add [m_castes] "\t[SELECT_ADDITIONAL_CASTE:MALE_" [id] "]\n")]
			[f_castes = (str:add [f_castes] "\t[SELECT_ADDITIONAL_CASTE:FEMALE_" [id] "]\n")]
		})
		(break true)
	})
	[out = (str:add [out] "\n" [m_castes] "\t\t" [rubble:libs_castes:creature_male [creature]] "\n")]
	[out = (str:add [out] "\n" [f_castes] "\t\t" [rubble:libs_castes:creature_female [creature]])]
	
	# Set bonuses
	(foreach [rubble:libs_castes:creature_castes [creature]] block id _ {
		[out = (str:add [out] 
			"\n\n"
			"\t[SELECT_CASTE:MALE_" [id] "]\n"
			"\t[SELECT_ADDITIONAL_CASTE:FEMALE_" [id] "]\n"
			"\t\t[CASTE_NAME:" [rubble:libs_castes:names [creature] [id]] ":" 
				[rubble:libs_castes:names_plur [creature] [id]] ":" [rubble:libs_castes:creature_adj [creature]] "]\n"
			"\t\t" [rubble:libs_castes:bonuses [creature] [id]]
		)]
		(break true)
	})
	
	[out]
})

command rubble:libs_castes:tabs tabs {
	var out = ""
	(for 1 [tabs] 1 block _ {
		[out = (str:add [out] "\t")]
	})
	(ret [out])
}

(rubble:template "#SYN_AFFECTED_MALES" block creature tabs=2 {
	var out = ""
	(foreach [rubble:libs_castes:creature_castes [creature]] block id _ {
		[out = (str:add [out] 
			(rubble:libs_castes:tabs [tabs]) "[SYN_AFFECTED_CREATURE:" [creature] ":MALE_" [id] "]\n")]
		(break true)
	})
	(str:trimspace [out])
})

(rubble:template "#SYN_AFFECTED_FEMALES" block creature tabs=2 {
	var out = ""
	(foreach [rubble:libs_castes:creature_castes [creature]] block id _ {
		[out = (str:add [out] 
			(rubble:libs_castes:tabs [tabs]) "[SYN_AFFECTED_CREATURE:" [creature] ":FEMALE_" [id] "]\n")]
		(break true)
	})
	(str:trimspace [out])
})
