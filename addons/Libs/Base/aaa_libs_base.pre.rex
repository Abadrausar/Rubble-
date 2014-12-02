
# Templates used in the base

# The comment templates, very simple, just do nothing.
(rubble:template "COMMENT" block a=... {(ret "")})
(rubble:template "C" block a=... {(ret "")})

command rubble:void params=... {
	(foreach [params] block _ raws {
		(rubble:stageparse [raws])
		(break true)
	})
	(ret "")
}
(rubble:template "!VOID" [rubble:void])
(rubble:template "VOID" [rubble:void])
(rubble:template "#VOID" [rubble:void])
(rubble:template "V" [rubble:void])

# ABORT, cause Rubble to exit with an error, not used in the base.
(rubble:template "!ABORT" [rubble:abort])
(rubble:template "ABORT" [rubble:abort])
(rubble:template "#ABORT" [rubble:abort])

# The next templates are for stripping leading/trailing whitespace from a string.
# A formatting tool mostly, helps keep whitespace under control in generated files.
# May also be used to help control variable expansion.
command rubble:echo params=... {
	var out = ""
	(foreach [params] block _ msg {
		[out = (str:add [out] [msg])]
	})
	(rubble:stageparse [out])
}
(rubble:template "!ECHO" [rubble:echo])
(rubble:template "ECHO" [rubble:echo])
(rubble:template "#ECHO" [rubble:echo])
(rubble:template "E" [rubble:echo])

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

(rubble:template "@SCRIPT" [global:eval])

(rubble:template "@SET" block name value {
	(rubble:configvar [name] [value])
	(ret "")
})

command rubble:print params=... {
	(foreach [params] block _ msg {
		(console:print "    " (rubble:expandvars [msg]) "\n")
		(break true)
	})
	(ret "")
}
(rubble:template "!PRINT" [rubble:print])
(rubble:template "PRINT" [rubble:print])
(rubble:template "#PRINT" [rubble:print])
