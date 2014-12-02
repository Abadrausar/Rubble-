
# Patch and compression templates.

(rubble:template "!PATCH" block file patch {
	[rubble:raws [file] = (rubble:patch [rubble:raws [file]] [patch])]
	(ret "")
})

(rubble:template "!DECOMPRESS" block text {
	(rubble:stageparse (rubble:decompress [text]))
})
