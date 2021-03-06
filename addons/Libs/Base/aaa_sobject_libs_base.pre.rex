
# SHARED_OBJECT, the single most important template in all of Rubble,
# also a bunch of related templates.

var rubble:shared_object_data = <map>
var rubble:shared_object_add_data = <map>
var rubble:shared_object_add_doonce = <map>
(rubble:template "SHARED_OBJECT" block id raws {
	(if (exists [rubble:shared_object_data] [id]){
		(rubble:stageparse [raws])
		(ret "")
	}{
		[rubble:shared_object_data [id] = (rubble:stageparse [raws])]
		(ret (str:add "{#_INSERT_SHARED_OBJECT;" [id] "}"))
	})
})
(rubble:template "SHARED_OBJECT_EXISTS" block id then else="" {
	(if (exists [rubble:shared_object_data] [id]) {
		(rubble:stageparse [then])
	}{
		(rubble:stageparse [else])
	})
})
command rubble:sharedobject_add id raws {
	# WARNING: Major issue with line endings detected!
	# Let's just strip any CR characters for now...
	# To keep the output consistent we only use the stripped version for the do-once key.
	var tmp = (str:replace [raws] "\r" "" -1)
	
	# Ouch, that could be a very expensive hash lookup...
	(if (exists [rubble:shared_object_add_doonce] [id]) {
		(if (exists [rubble:shared_object_add_doonce [id]] [tmp]) {
			(ret "")
		}{
			[rubble:shared_object_add_doonce [id] [tmp] = true]
		})
	}{
		[rubble:shared_object_add_doonce [id] = <map [tmp]=true>]
	})
	
	(if (exists [rubble:shared_object_add_data] [id]) {
		[raws = (str:add [rubble:shared_object_add_data [id]] "\n\t" (rubble:stageparse [raws]))]
		[rubble:shared_object_add_data [id] = [raws]]
	}{
		[raws = (str:add "\n\t" (rubble:stageparse [raws]))]
		[rubble:shared_object_add_data [id] = [raws]]
	})
}
(rubble:template "SHARED_OBJECT_ADD" block id raws {
	(rubble:sharedobject_add [id] [raws])
	(ret "")
})
(rubble:template "REGISTER_REACTION_CLASS" block id class {
	(rubble:sharedobject_add [id] (str:add "[REACTION_CLASS:" [class] "]"))
	(ret "")
})
(rubble:template "REGISTER_REACTION_PRODUCT" block id class mat {
	(rubble:sharedobject_add [id] (str:add "[MATERIAL_REACTION_PRODUCT:" [class] ":" [mat] "]"))
	(ret "")
})
(rubble:template "SHARED_OBJECT_KILL_TAG" block id tag {
	(if (exists [rubble:shared_object_data] [id]) {
	}{
		(ret "")
	})
	var file = (df:raw:parse [rubble:shared_object_data [id]])
	(foreach [file] block _ tag {
		(if (str:cmp [tag id] [tag]) {
			[tag disable = true]
		})
	})
	
	[rubble:shared_object_data [id] = (df:raw:dump [file])]
	(ret "")
})
(rubble:template "SHARED_OBJECT_REPLACE_TAG" block id tag replacement {
	(if (exists [rubble:shared_object_data] [id]) {
	}{
		(ret "")
	})
	var file = (df:raw:parse [rubble:shared_object_data [id]])
	(foreach [file] block _ tag {
		(if (str:cmp [tag id] [tag]) {
			[tag replace = [replacement]]
		})
	})
	
	[rubble:shared_object_data [id] = (df:raw:dump [file])]
	(ret "")
})
(rubble:template "#_INSERT_SHARED_OBJECT" block id {
	(if (exists [rubble:shared_object_data] [id]) {
		var out = (rubble:stageparse [rubble:shared_object_data [id]])
		(if (exists [rubble:shared_object_add_data] [id]) {
			[out = (str:add [out] (rubble:stageparse [rubble:shared_object_add_data [id]]))]
		})
		(ret [out])
	}{
		(ret "")
	})
})

# Specialized versions of SHARED_OBJECT

(rubble:template "SHARED_PLANT" block id raws {
	(rubble:stageparse (str:add
		"{SHARED_OBJECT;" [id] ";\n"
		"[PLANT:" [id] "]\n\t"
		[raws]
		"\n}"
	))
})

(rubble:template "SHARED_INORGANIC" block id raws {
	(rubble:stageparse (str:add
		"{SHARED_OBJECT;" [id] ";\n"
		"[INORGANIC:" [id] "]\n\t"
		[raws]
		"\n}"
	))
})

(rubble:template "SHARED_MATERIAL_TEMPLATE" block id raws {
	(rubble:stageparse (str:add
		"{SHARED_OBJECT;" [id] ";\n"
		"[MATERIAL_TEMPLATE:" [id] "]\n\t"
		[raws]
		"\n}"
	))
})
