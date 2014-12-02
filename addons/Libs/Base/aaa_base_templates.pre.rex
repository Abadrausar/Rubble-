
# Templates used in the base

# The comment templates, very simple, just do nothing.
(rubble:template "COMMENT" block ... {(ret "")})
(rubble:template "C" block ... {(ret "")})

command rubble:void params {
	(foreach [params] block _ raws {
		(rubble:stageparse [raws])
		(break true)
	})
	(ret "")
}
(rubble:template "!VOID" block ... {
	(rubble:void [params])
})
(rubble:template "VOID" block ... {
	(rubble:void [params])
})
(rubble:template "#VOID" block ... {
	(rubble:void [params])
})
(rubble:template "V" block ... {
	(rubble:void [params])
})

# ABORT, cause Rubble to exit with an error, not used in the base.
(rubble:template "!ABORT" block msg {(rubble:abort [msg])})
(rubble:template "ABORT" block msg {(rubble:abort [msg])})
(rubble:template "#ABORT" block msg {(rubble:abort [msg])})

# The next templates are for stripping leading/trailing whitespace from a string.
# A formatting tool mostly, helps keep whitespace under control in generated files.
# May also be used to help control variable expansion.
command rubble:echo params {
	var out = ""
	(foreach [params] block key msg {
		[out = (str:add [out] [msg])]
	})
	(rubble:stageparse [out])
}
(rubble:template "!ECHO" block ... {
	(rubble:echo [params])
})
(rubble:template "ECHO" block ... {
	(rubble:echo [params])
})
(rubble:template "#ECHO" block ... {
	(rubble:echo [params])
})
(rubble:template "E" block ... {
	(rubble:echo [params])
})

(rubble:template "@IF" block a b then else="" {
	(if (str:cmp [a] [b]) {
		(rubble:stageparse [then])
	}{
		(rubble:stageparse [else])
	})
})

(rubble:template "@IF_ACTIVE" block addon then else="" {
	(if (rubble:addonactive [addon]){
		(rubble:stageparse [then])
	}{
		(rubble:stageparse [else])
	})
})

(rubble:template "@IF_CODE" block code then else="" {
	(if (eval [code]) {
		(rubble:stageparse [then])
	}{
		(rubble:stageparse [else])
	})
})

(rubble:template "@IF_SKIP" block a b {
	(if (str:cmp [a] [b]) {
		(rubble:skipfile (rubble:currentfile))
	})
	(ret "")
})

(rubble:template "@SCRIPT" block code {
	(eval [code])
})

(rubble:template "@SET" block name value {
	(rubble:setvar [name] [value])
	(ret "")
})

var rubble:adventure_tier_data = 0
(rubble:template "#ADVENTURE_TIER" {
	[rubble:adventure_tier_data = (int:add [rubble:adventure_tier_data] 1)]
	(str:add "[ADVENTURE_TIER:" [rubble:adventure_tier_data] "]")
})

command rubble:print params {
	(foreach [params] block key msg {
		(console:print "    " (rubble:expandvars [msg]) "\n")
		(break true)
	})
	(ret "")
}
(rubble:template "!PRINT" block ... {(rubble:print [params])})
(rubble:template "PRINT" block ... {(rubble:print [params])})
(rubble:template "#PRINT" block ... {(rubble:print [params])})
